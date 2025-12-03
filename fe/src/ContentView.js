import React, { use } from 'react';
import { useState, useEffect } from 'react';
import { Card, Container, Spinner, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { FileView } from './DBFile';
import { NoteView } from './DBNote';
import { ObjectHeaderView, ObjectFooterView, ObjectView } from './DBObject';
import { HtmlView, PageView } from './DBPage';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
    CountryView,
    UserLinkView,
    ObjectLinkView
} from './sitenavigation_utils';
import axiosInstance from './axios';


// View for DBFolder
function FolderView({ data, metadata, dark }) {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language; // 'it', 'en', 'de', 'fr'

    const navigate = useNavigate();
    const { t } = useTranslation();
    
    const [indexContent, setIndexContent] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadIndexContent = async () => {
            try {
                setLoading(true);
                const response = await axiosInstance.get(`/nav/${data.id}/indexes`);
                const indexData = response.data;
                
                // IF indexData.indexes has more than one language, filter by currentLanguage in data.language array
                if (indexData.indexes && indexData.indexes.length > 0) {
                    const filteredIndexes = indexData.indexes.filter(index => index.data.language.indexOf(currentLanguage) >= 0);
                    if (filteredIndexes.length === 1) {
                        setIndexContent(filteredIndexes[0].data);
                    } else {
                        setIndexContent(indexData.indexes[0].data);
                    }
                } else {
                    setIndexContent({html: data.description});
                }
            } catch (err) {
                console.error('Error loading index content:', err);
            } finally {
                setLoading(false);
            }
        };
        
        loadIndexContent();
    }, [data.id, currentLanguage, data.description]);

    if (loading) {
        return (
            <Container className="mt-4 text-center">
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </Container>
        );
    }

    return (
        <div>
            {indexContent === null ? (
                <p>No indexes found in this folder.</p>
            ) : (
                <div>
                    {data.name && (
                    <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                    )}
                    {indexContent.description && (
                    <small style={{ opacity: 0.7 }}>-{indexContent.description}</small>
                    )}
                    <HtmlView html={indexContent.html} dark={dark} />
                </div>
            )}
        </div>
    );
    
    // return (
    //     <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
    //         <Card.Header>
    //             <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
    //             <small style={{ opacity: 0.7 }}>Folder · ID: {data.id}</small>
    //         </Card.Header>
    //         <Card.Body>
    //             {data.description && (
    //                 <Card.Text>{data.description}</Card.Text>
    //             )}
    //             <div style={{ opacity: 0.7 }}>
    //                 <small>Owner: {data.owner} | Group: {data.group_id}</small>
    //                 <br />
    //                 <small>Permissions: {data.permissions}</small>
    //                 {data.creation_date && (
    //                     <>
    //                         <br />
    //                         <small>Created: {data.creation_date}</small>
    //                     </>
    //                 )}
    //             </div>
    //         </Card.Body>
    //     </Card>
    // );
}

function PersonView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    return (
        <div>
            <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            {!data.html && data.description && <hr />}
            {data.description && (
                <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
            )}
            {data.html && <hr />}
            {data.html && (
                <HtmlView htmlContent={data.html} dark={dark} />
            )}
            <hr />
            {data.fk_users_id && data.fk_users_id !== "0" && (
                <p>👤 User: <UserLinkView user_id={data.fk_users_id} dark={dark} /></p>
            )}
            <p>
            {data.street}{data.street ? <br/> : ""}
            {data.zip} {data.city} {data.state ? `(${data.state})` : ''}{data.street || data.zip || data.city || data.state ? <br/> : ""}
            {data.fk_countrylist_id && (
                <CountryView country_id={data.fk_countrylist_id} dark={dark} />
            )}
            </p>
            {data.fk_companies_id && data.fk_companies_id !== "0" && (
                <p><ObjectLinkView obj_id={data.fk_companies_id} dark={dark} /></p>
            )}
            {data.phone && <p>📞 {data.phone}</p>}
            {data.office_phone && <p>🏢 {data.office_phone}</p>}
            {data.mobile && <p>📱 {data.mobile}</p>}
            {data.fax && <p>📠 {data.fax}</p>}
            {data.email && <p>✉️ <a href={`mailto:${data.email}`}>{data.email}</a></p>}
            {data.url && <p>🔗 <a href={data.url} target="_blank" rel="noopener noreferrer">{data.url}</a></p>}
            {data.codice_fiscale && <p>🆔 {data.codice_fiscale}</p>}
            {data.p_iva && <p>💰 {data.p_iva}</p>}
        </div>
    );
}

function CompanyView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    return (
        <div>
            <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            {!data.html && data.description && <hr />}
            {data.description && (
                <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
            )}
            {data.html && <hr />}
            {data.html && (
                <HtmlView htmlContent={data.html} dark={dark} />
            )}
            <hr />
            <p>
            {data.street}<br/>
            {data.zip} {data.city} ({data.state})<br/>
            <CountryView country_id={data.fk_countrylist_id} dark={dark} />
            </p>
            {data.phone && <p>📞 {data.phone}</p>}
            {data.office_phone && <p>🏢 {data.office_phone}</p>}
            {data.mobile && <p>📱 {data.mobile}</p>}
            {data.fax && <p>📠 {data.fax}</p>}
            {data.email && <p>✉️ <a href={`mailto:${data.email}`}>{data.email}</a></p>}
            {data.url && <p>🔗 <a href={data.url} target="_blank" rel="noopener noreferrer">{data.url}</a></p>}
            {data.p_iva && <p>💰 {data.p_iva}</p>}
        </div>
    );
}

// Main ContentView component - switches based on classname
function ContentView({ data, metadata, dark }) {
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
            return <FolderView data={data} metadata={metadata} dark={dark} />;
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

export default ContentView;
