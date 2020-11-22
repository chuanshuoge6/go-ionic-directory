import React, { useState, useEffect } from 'react';
import './Home.css';
import axios from 'axios'
import FolderDisplay from '../components/folderDisplay'

export default function Home() {
  const [dir, setDir] = useState([])

  useEffect(() => {
    document.title = 'directory'
    get_dir()
  }, [])

  function get_dir() {
    axios({
      method: 'get',
      url: 'http://192.168.0.18:8000/',
    })
      .then(response => {
        setDir(response.data)
      })
      .catch(function (error) {
        alert(error);
        setTimeout(() => {
          get_dir()
        }, 5000);
      });
  }



  return (
    <div style={{ overflowY: 'scroll', height: '100%' }}>
      <h1>Server</h1>
      <ul class="no-bullets">
        {
          dir.map((item, index) => {
            const dir_split = item.Name.split("\\")
            const folder_name = dir_split[dir_split.length - 1]
            return <FolderDisplay name={folder_name}
              dir={item.Name}
              isDir={item.IsDir}
              size={item.Size}
            />
          })
        }
      </ul>
    </div>
  );
};

