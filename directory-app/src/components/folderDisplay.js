import React, { useState } from 'react';
import { folder, folderOpenOutline, documentOutline } from 'ionicons/icons';
import { IonIcon } from '@ionic/react';
import axios from 'axios'
import '../pages/Home.css'
const qs = require('querystring')

export default function FolderDisplay(props) {
    const [open, setOpen] = useState(false)
    const [subDir, setSubDir] = useState([])

    function post_dir(d) {
        axios({
            method: 'post',
            url: 'http://localhost:8000/',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            data: qs.stringify({ dir: d })
        })
            .then(response => {
                setSubDir(response.data)
            })
            .catch(function (error) {
                alert(error);
            });
    }

    function clickEvent() {
        if (!props.isDir) { return }

        if (!open) {
            post_dir(props.dir)
        }

        setOpen(!open)
    }

    const folder_icon = open ? <IonIcon icon={folderOpenOutline} /> : <IonIcon icon={folder} />

    const sub_folder_list = subDir.map((item, index) => {
        const dir_split = item.Name.split("\\")
        const folder_name = dir_split[dir_split.length - 1].replace(props.name, '')
        const dir_name = item.Name.replace(folder_name, '') + '\\' + folder_name
        return <li>
            <FolderDisplay name={folder_name} dir={dir_name} isDir={item.IsDir} />
        </li>
    })

    return (
        <li >
            <span onClick={() => clickEvent()}>
                {props.isDir ? folder_icon : <IonIcon icon={documentOutline} />}{' '}
                {props.name}
            </span>
            {open ? <ul class="no-bullets">
                {sub_folder_list}
            </ul> : null}
        </li>
    );
};

