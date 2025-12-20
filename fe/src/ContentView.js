import React, { use } from 'react';
import { useState, useEffect } from 'react';
import { Card } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { CompanyView } from './DBCompany';
import { FileView } from './DBFile';
import { FolderView } from './DBFolders';
import { LinkView } from './DBLink';
import { NoteView } from './DBNote';
import { ObjectView } from './DBObject';
import { PageView } from './DBPage';
import { PersonView } from './DBPeople';
import axiosInstance from './axios';




// Main ContentView component - switches based on classname
export function ContentView({ data, metadata, dark, onFilesUploaded }) {
    const [objectData, setObjectData] = useState(null);
    
    // const token = localStorage.getItem("token");
    // const username = localStorage.getItem("username");
    // const groups = localStorage.getItem("groups");
      
    if (!data || !metadata) {
        return null;
    }

    useEffect(() => {
        // if (metadata.classname === 'DBFolder' || metadata.classname === 'DBPage' || metadata.classname === 'DBNews') {
        //     return;
        // }
        
        const loadUserData = async () => {
            try {
                // Collect unique user IDs
                const uniqueUserIds = new Set();
                if (data.owner) uniqueUserIds.add(data.owner);
                if (data.creator) uniqueUserIds.add(data.creator);
                if (data.last_modify) uniqueUserIds.add(data.last_modify);
                if (data.deleted_by) uniqueUserIds.add(data.deleted_by);
                
                // Fetch all unique users in parallel
                const userPromises = Array.from(uniqueUserIds).map(userId =>
                    axiosInstance.get(`/users/${userId}`).then(res => ({ id: userId, data: res.data }))
                );
                
                const groupPromise = data.group_id && data.group_id!=="0" ? axiosInstance.get(`/groups/${data.group_id}`) : Promise.resolve({data: { name: '' }});
                
                const [users, groupRes] = await Promise.all([
                    Promise.all(userPromises),
                    groupPromise
                ]);
                
                // Create a map of userId -> user data
                const userMap = {};
                users.forEach(user => {
                    userMap[user.id] = user.data.fullname;
                });
                
                setObjectData({
                    owner_name: userMap[data.owner] || '',
                    group_name: groupRes.data.name,
                    creator_name: userMap[data.creator] || '',
                    last_modifier_name: userMap[data.last_modify] || '',
                    deleted_by_name: userMap[data.deleted_by] || ''
                });
            } catch (error) {
                console.error('Error loading user data:', error);
            }
        };
        
        loadUserData();
    }, [data.owner, data.group_id, data.creator, data.last_modify, data.deleted_by, metadata.classname]);

    const classname = metadata.classname;

    switch (classname) {
        case 'DBCompany':
            return <CompanyView data={data} metadata={metadata} objectData={objectData} dark={dark} />;
        case 'DBPerson':
            return <PersonView data={data} metadata={metadata} objectData={objectData} dark={dark} />;
        // // CMS
        // case 'DBEvent':
        //     return <EventView data={data} metadata={metadata} dark={dark} />;
        case 'DBFile':
            return <FileView data={data} metadata={metadata} dark={dark} />;
        case 'DBFolder':
            return <FolderView data={data} metadata={metadata} dark={dark} onFilesUploaded={onFilesUploaded} />;
        case 'DBLink':
            return <LinkView data={data} metadata={metadata} objectData={objectData} dark={dark} />;
        case 'DBNote':
            return <NoteView data={data} metadata={metadata} objectData={objectData} dark={dark} />;
        case 'DBNews':
        //     return <NewsView data={data} metadata={metadata} dark={dark} />;
        case 'DBPage':
            return <PageView data={data} metadata={metadata} dark={dark} />;
        default:
            return <ObjectView data={data} metadata={metadata} objectData={objectData} dark={dark} />;
    }
}

// export default ContentView;
