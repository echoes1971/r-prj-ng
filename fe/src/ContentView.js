import React from 'react';
import { useState, useEffect } from 'react';
import { Card, Container, Spinner } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

// View for DBFolder
function FolderView({ data, metadata, dark }) {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language; // 'it', 'en', 'de', 'fr'
    
    const [indexContent, setIndexContent] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadIndexContent();
    }, [data.id, currentLanguage]);

    const loadIndexContent = async () => {
        try {
            setLoading(true);
            const response = await fetch(`/api/nav/${data.id}/indexes`);
            const indexData = await response.json();
            // setIndexContent(indexData);
            // alert("currentLanguage: " + currentLanguage);
            // alert(JSON.stringify(indexData));
            // const test = indexData.indexes.filter(index => index.data.language.indexOf(currentLanguage)>=0)
            // alert("Filtered indexes: " + JSON.stringify(test));
            // IF indexData.indexes has more than one language, filter by currentLanguage in data.language array
            if (indexData.indexes && indexData.indexes.length > 0) {
                const filteredIndexes = indexData.indexes.filter(index => index.data.language.indexOf(currentLanguage)>=0);
                if (filteredIndexes.length == 1) {
                    setIndexContent(filteredIndexes[0].data);
                } else {
                    setIndexContent(indexData.indexes[0].data);
                }
            } else {
                setIndexContent({html:data.description});
            }
        } catch (err) {
            console.error('Error loading index content:', err);
        } finally {
            setLoading(false);
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
    return (
        <div>
            {/* {data.name && (
            <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            )}
            {data.description && (
            <small style={{ opacity: 0.7 }}>{data.description}</small>
            )} */}
            <div dangerouslySetInnerHTML={{ __html: data.html }}></div>
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
function NoteView({ data, metadata, dark }) {
    return (
        <Card className="mb-3 border-warning" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                <small style={{ opacity: 0.7 }}>Note · ID: {data.id}</small>
            </Card.Header>
            <Card.Body>
                {data.description && (
                    <Card.Text>{data.description}</Card.Text>
                )}
                {data.content && (
                    <div className="content">
                        {data.content}
                    </div>
                )}
                <div className="text-muted mt-3">
                    <small>Owner: {data.owner} | Group: {data.group_id}</small>
                    <br />
                    <small>Permissions: {data.permissions}</small>
                </div>
            </Card.Body>
        </Card>
    );
}

// Generic view for DBObject
function ObjectView({ data, metadata, dark }) {
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                <small style={{ opacity: 0.7 }}>{metadata.classname || 'Object'} · ID: {data.id}</small>
            </Card.Header>
            <Card.Body>
                {data.description && (
                    <Card.Text>{data.description}</Card.Text>
                )}
                <div className="text-muted">
                    <small>Owner: {data.owner} | Group: {data.group_id}</small>
                    <br />
                    <small>Permissions: {data.permissions}</small>
                </div>
            </Card.Body>
        </Card>
    );
}

// Main ContentView component - switches based on classname
function ContentView({ data, metadata, dark }) {
    if (!data || !metadata) {
        return null;
    }

    const classname = metadata.classname;

    switch (classname) {
        case 'DBFolder':
            return <FolderView data={data} metadata={metadata} dark={dark} />;
        case 'DBPage':
        case 'DBNews':
            return <PageView data={data} metadata={metadata} dark={dark} />;
        case 'DBNote':
            return <NoteView data={data} metadata={metadata} dark={dark} />;
        default:
            return <ObjectView data={data} metadata={metadata} dark={dark} />;
    }
}

export default ContentView;
