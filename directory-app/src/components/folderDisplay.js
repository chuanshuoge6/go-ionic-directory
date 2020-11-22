import React, { useState } from 'react';
import { folder, folderOpenOutline, documentOutline, cloudUploadOutline, cloudDownloadOutline } from 'ionicons/icons';
import { IonIcon } from '@ionic/react';
import axios from 'axios'
import '../pages/Home.css'
import saveAs from 'file-saver';
const qs = require('querystring')

export default function FolderDisplay(props) {
    const [open, setOpen] = useState(false)
    const [subDir, setSubDir] = useState([])
    const [showLoad, setShowLoad] = useState(false)
    const [downlaodProgress, setDownloadProgress] = useState(0)

    function post_dir(d) {
        axios({
            method: 'post',
            url: 'http://192.168.0.18:8000/',
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

    function download_dir(d, name) {
        axios({
            method: 'post',
            url: 'http://192.168.0.18:8000/download/',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            data: qs.stringify({ dir: d }),
            responseType: 'blob',  // important

            onDownloadProgress: progressEvent => {
                const p = parseInt(progressEvent.loaded / props.size * 100);
                setDownloadProgress(p)
            },
        })
            .then(response => {
                console.log(response)
                saveAs(response.data, name);
            })
            .catch(function (error) {
                alert(error);
            });
    }

    function clickEvent() {
        setShowLoad(!showLoad)

        if (!props.isDir) { return }

        if (!open) {
            post_dir(props.dir)
        }

        setOpen(!open)
    }

    function downloadEvent() {
        download_dir(props.dir, props.name)
    }

    const folder_icon = open ? <IonIcon icon={folderOpenOutline} /> : <IonIcon icon={folder} />

    const sub_folder_list = subDir.map((item, index) => {
        const dir_split = item.Name.split("\\")
        const folder_name = dir_split[dir_split.length - 1].replace(props.name, '')
        const dir_name = item.Name.replace(folder_name, '') + '\\' + folder_name
        return <li>
            <FolderDisplay name={folder_name} dir={dir_name} isDir={item.IsDir} size={item.Size} />
        </li>
    })

    return (
        <li >
            <span onClick={() => clickEvent()}>
                {props.isDir ? folder_icon : <IonIcon icon={documentOutline} />}{' '}
                {props.name}{' '}
                {showLoad && props.isDir ? <IonIcon icon={cloudUploadOutline} /> : null}{' '}
                {showLoad && !props.isDir ? <IonIcon icon={cloudDownloadOutline} onClick={() => downloadEvent()} /> : null}{' '}
                {downlaodProgress != 0 ? downlaodProgress + '%' : null}
            </span>
            {open ? <ul class="no-bullets">
                {sub_folder_list}
            </ul> : null}
        </li>
    );
};

