package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/gorilla/mux"
)

type fileStruct struct {
	Name  string
	IsDir bool
}

var projectPath = rootDir() + "\\golang12\\"

func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	dir := viewDirectory(projectPath)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dir)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)

	r.ParseForm()
	fmt.Println(r.PostForm)
	path := r.PostForm.Get("dir")
	fmt.Println(path)
	dir := viewDirectory(path)

	json.NewEncoder(w).Encode(dir)
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
		dir = append(dir, fileStruct{Name: name, IsDir: file.IsDir()})
	}

	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name < dir[j].Name
	})

	fmt.Println(dir)

	return dir
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST", "OPTIONS")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
