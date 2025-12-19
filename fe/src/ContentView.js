import React, { use } from 'react';
import { useState, useEffect } from 'react';
import { Card, Container, Spinner, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { FileView } from './DBFile';
import { LinkView } from './DBLink';
import { NoteView } from './DBNote';
import { ObjectHeaderView, ObjectFooterView, ObjectView } from './DBObject';
import { HtmlView, PageView } from './DBPage';
import {
    CountryView,
    UserLinkView,
    ObjectLinkView
} from './ContentWidgets'
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon
} from './sitenavigation_utils';
import axiosInstance from './axios';


// View for DBFolder
function FolderView({ data, metadata, dark, onFilesUploaded }) {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language; // 'it', 'en', 'de', 'fr'

    const navigate = useNavigate();
    const { t } = useTranslation();
    
    const [indexContent, setIndexContent] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isDragging, setIsDragging] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState({ current: 0, total: 0 });
    const [uploadError, setUploadError] = useState(null);

    const MAX_FILES = 30;

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
                    // setIndexContent({html: data.description});
                }
            } catch (err) {
                console.error('Error loading index content:', err);
            } finally {
                setLoading(false);
            }
        };
        
        loadIndexContent();
    }, [data.id, currentLanguage, data.description]);

    const handleDragEnter = (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (metadata.can_edit && !uploading) {
            setIsDragging(true);
        }
    };

    const handleDragLeave = (e) => {
        e.preventDefault();
        e.stopPropagation();
        // Only set to false if leaving the main container
        if (e.currentTarget === e.target) {
            setIsDragging(false);
        }
    };

    const handleDragOver = (e) => {
        e.preventDefault();
        e.stopPropagation();
    };

    const handleDrop = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);

        if (!metadata.can_edit || uploading) return;

        const files = Array.from(e.dataTransfer.files);
        
        if (files.length === 0) return;

        if (files.length > MAX_FILES) {
            alert(t('files.too_many_files', { max: MAX_FILES, count: files.length }) || 
                  `Too many files! Maximum ${MAX_FILES} files at once. You tried to upload ${files.length} files.`);
            return;
        }

        await uploadFiles(files);
    };

    const uploadFiles = async (files) => {
        setUploading(true);
        setUploadProgress({ current: 0, total: files.length });
        setUploadError(null);

        const errors = [];

        for (let i = 0; i < files.length; i++) {
            const file = files[i];
            setUploadProgress({ current: i + 1, total: files.length });

            try {
                const formData = new FormData();
                formData.append('file', file);
                formData.append('name', file.name);
                formData.append('father_id', data.id);
                formData.append('permissions', 'rw-r-----'); // Default permissions

                await axiosInstance.post('/objects', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                });
            } catch (error) {
                console.error(`Error uploading ${file.name}:`, error);
                errors.push(`${file.name}: ${error.response?.data?.error || 'Upload failed'}`);
            }
        }

        setUploading(false);

        if (errors.length > 0) {
            setUploadError(errors.join('\n'));
        }

        // Notify parent to refresh children list
        if (onFilesUploaded) {
            onFilesUploaded();
        }
    };

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
        <div
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            style={{
                position: 'relative',
                border: isDragging ? '3px dashed #0d6efd' : 'none',
                borderRadius: '8px',
                padding: isDragging ? '10px' : '0',
                backgroundColor: isDragging ? (dark ? 'rgba(13, 110, 253, 0.1)' : 'rgba(13, 110, 253, 0.05)') : 'transparent',
                transition: 'all 0.2s ease',
            }}
        >
            {isDragging && metadata.can_edit && (
                <div
                    style={{
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        right: 0,
                        bottom: 0,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        backgroundColor: dark ? 'rgba(0, 0, 0, 0.7)' : 'rgba(255, 255, 255, 0.7)',
                        borderRadius: '8px',
                        zIndex: 1000,
                        pointerEvents: 'none',
                    }}
                >
                    <div className="text-center">
                        <i className="bi bi-cloud-upload" style={{ fontSize: '3rem', color: '#0d6efd' }}></i>
                        <h4 className={dark ? 'text-light mt-2' : 'text-dark mt-2'}>
                            {t('files.drop_files_here') || 'Drop files here to upload'}
                        </h4>
                        <p className="text-secondary">
                            {t('files.max_files', { max: MAX_FILES }) || `Maximum ${MAX_FILES} files at once`}
                        </p>
                    </div>
                </div>
            )}

            {uploading && (
                <div className="alert alert-info mb-3">
                    <Spinner animation="border" size="sm" className="me-2" />
                    {t('files.uploading') || 'Uploading'} {uploadProgress.current} / {uploadProgress.total}...
                </div>
            )}

            {uploadError && (
                <div className="alert alert-warning mb-3">
                    <strong>{t('files.upload_errors') || 'Upload errors'}:</strong>
                    <pre style={{ whiteSpace: 'pre-wrap', fontSize: '0.9em', marginTop: '0.5rem' }}>
                        {uploadError}
                    </pre>
                </div>
            )}

            {indexContent === null ? (
                <div>
                    {data.name && (
                    <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                    )}
                    {data.description && (
                    <small style={{ opacity: 0.7 }}>{data.description}</small>
                    )}
                </div>
            ) : (
                <div>
                    {/* {data.name && (
                    <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                    )}
                    {indexContent.description && (
                    <small style={{ opacity: 0.7 }}>-{indexContent.description}</small>
                    )} */}
                    <HtmlView html={indexContent.html} dark={dark} />
                </div>
            )}
            
            {metadata.can_edit && !uploading && (
                // Invisible on small screens
                <div 
                    className="alert alert-info mt-3 d-none d-md-block" 
                    style={{ 
                        borderStyle: 'dashed',
                        cursor: 'default'
                    }}
                >
                    <i className="bi bi-cloud-upload me-2"></i>
                    {t('files.drag_drop_hint') || 'Drag & drop files here to upload them to this folder'}
                    <small className="d-block mt-1 text-secondary">
                        ({t('files.max_files', { max: MAX_FILES }) || `Maximum ${MAX_FILES} files at once`})
                    </small>
                </div>
            )}
        </div>
    );
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
function ContentView({ data, metadata, dark, onFilesUploaded }) {
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

export default ContentView;
