import React, { use } from 'react';
import { useState, useEffect } from 'react';
import { Card, Container, Spinner, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
    CountryView,
    UserLinkView,
    ObjectLinkView,
    HtmlFieldView
} from './sitenavigation_utils';
import axiosInstance from './axios';

// View for DBFolder
function FolderView({ data, metadata, dark }) {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language; // 'it', 'en', 'de', 'fr'

    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;
    
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
                    <div dangerouslySetInnerHTML={{ __html: indexContent.html }}>
                        {/* <h3>Indexes in this folder:</h3>
                        <ul>
                            {indexContent.indexes.map((index) => (
                                <li key={index.data.id}>{index.data.name} (Language: {index.data.language}) (ID: {index.data.id})</li>
                            ))}
                        </ul> */}
                    </div>
                    {canEdit && (
                            <Button 
                                variant="primary" 
                                size="sm" 
                                className="mt-2"
                                onClick={() => navigate(`/e/${data.id}`)}
                            >
                                <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                            </Button>
                    )}
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

// View for DBPage
function PageView({ data, metadata, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;

    return (
        <div>
            {data.name && (
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            )}
            {data.description && (
                <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
            )}
            {data.html && (
                <div dangerouslySetInnerHTML={{ __html: data.html }}></div>
            )}
            {canEdit && (
                <>
                    <br />
                    <Button 
                        variant="primary" 
                        size="sm" 
                        className="mt-2"
                        onClick={() => navigate(`/e/${data.id}`)}
                    >
                        <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                    </Button>
                </>
            )}
        </div>
    );
    // return (
    //     <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
    //         <Card.Header>
    //             <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
    //             <small style={{ opacity: 0.7 }}>Page · ID: {data.id}</small>
    //         </Card.Header>
    //         <Card.Body>
    //             {data.description && (
    //                 <Card.Text className="lead">{data.description}</Card.Text>
    //             )}
    //             {data.content && (
    //                 <div 
    //                     className="content"
    //                     dangerouslySetInnerHTML={{ __html: data.content }}
    //                 />
    //             )}
    //             <div className="text-muted mt-3">
    //                 <small>Owner: {data.owner} | Group: {data.group_id}</small>
    //                 <br />
    //                 <small>Permissions: {data.permissions}</small>
    //                 {data.last_modify_date && (
    //                     <>
    //                         <br />
    //                         <small>Last modified: {data.last_modify_date}</small>
    //                     </>
    //                 )}
    //             </div>
    //         </Card.Body>
    //     </Card>
    // );
}

// View for DBNote
function NoteView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;

    return (
        <Card className="mb-3 border-warning" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                {/* <small style={{ opacity: 0.7 }}><i className={`bi bi-${classname2bootstrapIcon(metadata.classname)}`} title={metadata.classname}></i> Note · ID: {data.id}</small> */}
            </Card.Header>
            <Card.Body>
                {data.description && (
                    // <Card.Text>{data.description}</Card.Text>
                    <div className="content">
                        <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
                    </div>
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                {canEdit && (
                    <>
                        <br />
                        <Button 
                            variant="primary" 
                            size="sm" 
                            className="mt-2"
                            onClick={() => navigate(`/e/${data.id}`)}
                        >
                            <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                        </Button>
                    </>
                )}
            </Card.Footer>
        </Card>
    );
}

function PersonView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                {data.father_id && data.father_id!=="0" && <small style={{ opacity: 0.7 }}>Parent: <ObjectLinkView obj_id={data.father_id} dark={dark} /></small>}
                {data.father_id && data.father_id!=="0" && <br />}
                <small style={{ opacity: 0.7 }}><i className={`bi bi-${classname2bootstrapIcon(metadata.classname)}`} title={metadata.classname}></i> ID: {data.id}</small>
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <br />}
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <small style={{ opacity: 0.7 }}>Linked to: <ObjectLinkView obj_id={data.fk_obj_id} dark={dark} /></small>}
            </Card.Header>
            <Card.Body>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                {!data.html && data.description && <hr />}
                {data.description && (
                    <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
                )}
                {data.html && <hr />}
                {data.html && (
                    <HtmlFieldView htmlContent={data.html} dark={dark} />
                )}
                <hr />
                {data.fk_users_id && data.fk_users_id !== "0" && (
                    <p>👤 User: <UserLinkView user_id={data.fk_users_id} dark={dark} /></p>
                )}
                <p>
                {data.street}<br/>
                {data.zip} {data.city} ({data.state})<br/>
                <CountryView country_id={data.fk_countrylist_id} dark={dark} />
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
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <small style={{ opacity: 0.7 }}>Owner: {objectData && objectData.owner_name} | Group: {objectData && objectData.group_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Permissions: {data.permissions}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Created: {formateDateTimeString(data.creation_date)} - {objectData && objectData.creator_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Last update: {formateDateTimeString(data.last_modify_date)} - {objectData && objectData.last_modifier_name}</small>
                {data.deleted_date && <br />}
                {data.deleted_date && <small style={{ opacity: 0.7 }}>Deleted: {formateDateTimeString(data.deleted_date)} - {objectData && objectData.deleted_by_name}</small>}
                {canEdit && (
                    <>
                        <br />
                        <Button 
                            variant="primary" 
                            size="sm" 
                            className="mt-2"
                            onClick={() => navigate(`/e/${data.id}`)}
                        >
                            <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                        </Button>
                    </>
                )}
            </Card.Footer>
        </Card>
    );
}

function CompanyView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                {data.father_id && data.father_id!=="0" && <small style={{ opacity: 0.7 }}>Parent: <ObjectLinkView obj_id={data.father_id} dark={dark} /></small>}
                {data.father_id && data.father_id!=="0" && <br />}
                <small style={{ opacity: 0.7 }}><i className={`bi bi-${classname2bootstrapIcon(metadata.classname)}`} title={metadata.classname}></i> ID: {data.id}</small>
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <br />}
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <small style={{ opacity: 0.7 }}>Linked to: <ObjectLinkView obj_id={data.fk_obj_id} dark={dark} /></small>}
            </Card.Header>
            <Card.Body>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                {!data.html && data.description && <hr />}
                {data.description && (
                    <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
                )}
                {data.html && <hr />}
                {data.html && (
                    <HtmlFieldView htmlContent={data.html} dark={dark} />
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
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <small style={{ opacity: 0.7 }}>Owner: {objectData && objectData.owner_name} | Group: {objectData && objectData.group_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Permissions: {data.permissions}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Created: {formateDateTimeString(data.creation_date)} - {objectData && objectData.creator_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Last update: {formateDateTimeString(data.last_modify_date)} - {objectData && objectData.last_modifier_name}</small>
                {data.deleted_date && <br />}
                {data.deleted_date && <small style={{ opacity: 0.7 }}>Deleted: {formateDateTimeString(data.deleted_date)} - {objectData && objectData.deleted_by_name}</small>}
                {canEdit && (
                    <>
                        <br />
                        <Button 
                            variant="primary" 
                            size="sm" 
                            className="mt-2"
                            onClick={() => navigate(`/e/${data.id}`)}
                        >
                            <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                        </Button>
                    </>
                )}
            </Card.Footer>
        </Card>
    );
}

// Generic view for DBObject
function ObjectView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const canEdit = metadata.can_edit;
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                {data.father_id && data.father_id!=="0" && <small style={{ opacity: 0.7 }}>Parent: <ObjectLinkView obj_id={data.father_id} dark={dark} /></small>}
                {data.father_id && data.father_id!=="0" && <br />}
                <small style={{ opacity: 0.7 }}><i className={`bi bi-${classname2bootstrapIcon(metadata.classname)}`} title={metadata.classname}></i> ID: {data.id}</small>
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <br />}
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <small style={{ opacity: 0.7 }}>Linked to: <ObjectLinkView obj_id={data.fk_obj_id} dark={dark} /></small>}
            </Card.Header>
            <Card.Body>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                {!data.html && data.description && <hr />}
                {data.description && (
                    <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
                )}
                {data.html && <hr />}
                {data.html && (
                    <HtmlFieldView htmlContent={data.html} dark={dark} />
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <small style={{ opacity: 0.7 }}>Owner: {objectData && objectData.owner_name} | Group: {objectData && objectData.group_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Permissions: {data.permissions}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Created: {formateDateTimeString(data.creation_date)} - {objectData && objectData.creator_name}</small>
                <br />
                <small style={{ opacity: 0.7 }}>Last update: {formateDateTimeString(data.last_modify_date)} - {objectData && objectData.last_modifier_name}</small>
                {data.deleted_date && <br />}
                {data.deleted_date && <small style={{ opacity: 0.7 }}>Deleted: {formateDateTimeString(data.deleted_date)} - {objectData && objectData.deleted_by_name}</small>}
                {canEdit && (
                    <>
                        <br />
                        <Button 
                            variant="primary" 
                            size="sm" 
                            className="mt-2"
                            onClick={() => navigate(`/e/${data.id}`)}
                        >
                            <i className="bi bi-pencil me-1"></i>{t('common.edit')}
                        </Button>
                    </>
                )}
            </Card.Footer>
        </Card>
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
        // case 'DBFile':
        //     return <FileView data={data} metadata={metadata} dark={dark} />;
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
