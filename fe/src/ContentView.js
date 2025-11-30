import React, { use } from 'react';
import { useState, useEffect } from 'react';
import { Card, Container, Spinner, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ObjectHeaderView, ObjectFooterView, ObjectView } from './DBObject';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
    CountryView,
    UserLinkView,
    // ObjectHeaderView,
    // ObjectFooterView,
    ObjectLinkView,
    HtmlFieldView
} from './sitenavigation_utils';
import axiosInstance from './axios';

function FileView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [preview, setPreview] = useState(null);
    
    useEffect(() => {
        console.log('FileEdit useEffect:', { id: data.id, filename: data.filename });
        // Load preview if file exists
        if (data.id && data.filename && data.mime.indexOf('image') === 0) {
            // Fetch the file with authorization header and create a blob URL
            const loadPreview = async () => {
                try {
                    console.log('Loading file preview for:', data.id, 'filename:', data.filename);
                    const response = await axiosInstance.get(`/files/${data.id}/download`, {
                        responseType: 'blob'
                    });
                    console.log('File loaded, blob size:', response.data.size, 'type:', response.data.type);
                    const blobUrl = URL.createObjectURL(response.data);
                    console.log('Blob URL created:', blobUrl);
                    setPreview(blobUrl);
                } catch (error) {
                    console.error('Failed to load file preview:', error);
                    setPreview(null);
                }
            };
            loadPreview();

            // Cleanup blob URL on unmount
            return () => {
                if (preview && preview.startsWith('blob:')) {
                    console.log('Revoking blob URL:', preview);
                    URL.revokeObjectURL(preview);
                }
            };
        } else {
            console.log('Skipping preview load - condition not met');
        }
    }, [data.id, data.filename]);

    return (
        <div>
            {data.name && (
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            )}
            {data.description && (
                <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
            )}
            {data.filename && preview && (
                <div>
                    <img 
                        src={preview}
                        alt="Preview"
                        title={data.name}
                        style={{ maxWidth: '100%', maxHeight: '300px', marginBottom: '10px' }}
                    />
                    {/* <div>
                        <small className={dark ? "text-white-50" : "text-muted"}>{data.name}</small>
                    </div> */}
                    <div>
                        <a href={preview} download={data.filename}>
                            <Button variant="primary">
                                <i className="bi bi-download"></i> {t('dbobjects.download_file')}
                                {/* ({data.filename}) */}
                            </Button>
                        </a>
                    </div>
                </div>
            )}
            {data.filename && data.mime.indexOf('image') !== 0 && (
                <div>
                    <a href={`../api/files/${data.id}/download?token=${metadata.download_token}`} download={data.filename}>
                        <Button variant="primary">
                            <i className="bi bi-download"></i> {t('dbobjects.download_file')}
                        </Button>
                    </a>
                </div>
            )}
        </div>
    );
}

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
                    <div dangerouslySetInnerHTML={{ __html: indexContent.html }}>
                        {/* <h3>Indexes in this folder:</h3>
                        <ul>
                            {indexContent.indexes.map((index) => (
                                <li key={index.data.id}>{index.data.name} (Language: {index.data.language}) (ID: {index.data.id})</li>
                            ))}
                        </ul> */}
                    </div>
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

    return (
        <Card className="mb-3 border-warning" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Header>
            <Card.Body>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                <hr />
                {data.description && (
                    // <Card.Text>{data.description}</Card.Text>
                    <div className="content">
                        <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
                    </div>
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Footer>
        </Card>
    );
}

function PersonView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Header>
            <Card.Body className={dark ? 'bg-secondary bg-opacity-10' : ''}>
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
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Footer>
        </Card>
    );
}

function CompanyView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Header>
            <Card.Body className={dark ? 'bg-secondary bg-opacity-10' : ''}>
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
                <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Footer>
        </Card>
    );
}

// // Generic view for DBObject
// function ObjectView({ data, metadata, objectData, dark }) {
//     const navigate = useNavigate();
//     const { t } = useTranslation();
    
//     return (
//         <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
//             <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
//                 <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} />
//             </Card.Header>
//             <Card.Body>
//                 <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
//                 {!data.html && data.description && <hr />}
//                 {data.description && (
//                     <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
//                 )}
//                 {data.html && <hr />}
//                 {data.html && (
//                     <HtmlFieldView htmlContent={data.html} dark={dark} />
//                 )}
//             </Card.Body>
//             <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
//                 <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
//             </Card.Footer>
//         </Card>
//     );
// }

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
