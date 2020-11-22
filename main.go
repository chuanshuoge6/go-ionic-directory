package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type fileStruct struct {
	Name  string
	IsDir bool
	Size  int64
}

var projectPath = rootDir() + "\\golang12\\"

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index get")
	dir := viewDirectory(projectPath)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dir)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index post")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)

	r.ParseForm()
	path := r.PostForm.Get("dir")
	dir := viewDirectory(path)

	json.NewEncoder(w).Encode(dir)
}

func downloadPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("download post")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)

	r.ParseForm()
	dir := r.PostForm.Get("dir")
	dirSlash := strings.Replace(dir, "\\", "/", -1)
	_, filename := path.Split(dirSlash)

	fmt.Println(dir)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))

	file, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	io.Copy(w, file)
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func viewDirectory(dirname string) []fileStruct {
	f, err := os.Open(dirname)
	if err != nil {
		//log.Fatal(err)
		return []fileStruct{}
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		//log.Fatal(err)
		return []fileStruct{}
	}

	var dir = []fileStruct{}
	for _, file := range files {
		name := dirname + file.Name()
		var size int64 = 10000000
		if !file.IsDir() {
			size = file.Size()
		}

		dir = append(dir, fileStruct{
			Name:  name,
			IsDir: file.IsDir(),
			Size:  size,
		})
	}

	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name < dir[j].Name
	})

	return dir
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/download/", downloadPostHandler).Methods("POST", "OPTIONS")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
