import React, { useState, useEffect, useContext } from 'react';
import { Button, Container, Row, Col, Breadcrumb, Spinner, Alert } from 'react-bootstrap';
import { useParams, useNavigate, Link } from 'react-router-dom';
import axiosInstance from './axios';
import { useTranslation } from 'react-i18next';
import { app_cfg } from './app.cfg';
import ContentView from './ContentView';
import NewObjectButton from './NewObjectButton';
import ObjectList from './ObjectList';
import { getErrorMessage } from './errorHandler';
import { ThemeContext } from './ThemeContext';
import { formatObjectId } from './sitenavigation_utils';

function SiteNavigation() {
    const { objectId } = useParams();
    const navigate = useNavigate();
    const { dark, themeClass } = useContext(ThemeContext);
    const { t } = useTranslation();
    
    // Use home object ID from config if no objectId in URL
    const currentObjectId = objectId || app_cfg.app_home_object_id;
    
    const [content, setContent] = useState(null);
    const [breadcrumb, setBreadcrumb] = useState([]);
    const [children, setChildren] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        loadContent();
        loadBreadcrumb();
        loadChildren();
    }, [currentObjectId]);

    const loadContent = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await axiosInstance.get(`/content/${formatObjectId(currentObjectId)}`);
            setContent(response.data);
        } catch (err) {
            console.error('Error loading content:', err);
            setError(getErrorMessage(err));
        } finally {
            setLoading(false);
        }
    };

    const loadBreadcrumb = async () => {
        try {
            const response = await axiosInstance.get(`/nav/breadcrumb/${formatObjectId(currentObjectId)}`);
            setBreadcrumb(response.data.breadcrumb || []);
        } catch (err) {
            console.error('Error loading breadcrumb:', err);
        }
    };

    const loadChildren = async () => {
        try {
            const response = await axiosInstance.get(`/nav/children/${formatObjectId(currentObjectId)}`);
            setChildren(response.data.children || []);
        } catch (err) {
            console.error('Error loading children:', err);
        }
    };

    const handleNavigate = (objId) => {
        navigate(`/c/${formatObjectId(objId)}`);
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

    if (error) {
        return (
            <Container className="mt-4">
                <Alert variant="danger">
                    {error}
                </Alert>
            </Container>
        );
    }

    if (!content) {
        return (
            <Container className="mt-4">
                <Alert variant="warning">
                    Content not found
                </Alert>
            </Container>
        );
    }

    return (
        <Container className={`mt-4 ${themeClass}`}>
            <Row>
                <Col>
                    {/* Breadcrumb with Edit and New buttons */}
                    {/* {breadcrumb.length > 0 && ( */}
                        <div className="d-flex justify-content-between align-items-center mb-3">
                                <Breadcrumb className="mb-0" data-bs-theme={dark ? 'dark' : 'light'}>
                                    {breadcrumb.length > 1 && breadcrumb.map((item, index) => (
                                        <Breadcrumb.Item
                                            key={item.data.id}
                                            active={index === breadcrumb.length - 1}
                                            linkAs={Link}
                                            linkProps={{ to: `/c/${formatObjectId(item.data.id)}` }}
                                        >
                                            {item.data.name}
                                        </Breadcrumb.Item>
                                    ))}
                                </Breadcrumb>
                            {!loading && content && content.metadata && content.metadata.can_edit && (
                                <div className="d-flex gap-2">
                                    <Button 
                                        variant="outline-primary" 
                                        size="sm"
                                        onClick={() => navigate(`/e/${currentObjectId}`)}
                                    >
                                        <i className="bi bi-pencil me-1"></i>
                                        {t('common.edit')}
                                    </Button>
                                    <NewObjectButton 
                                        fatherId={currentObjectId}
                                        onObjectCreated={() => {
                                            loadChildren(); // Refresh children list
                                        }}
                                    />
                                </div>
                            )}
                        </div>
                    {/* )} */}

                    {/* Main Content */}
                    <ContentView 
                        data={content.data} 
                        metadata={content.metadata}
                        dark={dark}
                    />
                    {/* Children List*/}
                    {children.length > 0 && (
                        <div className="mt-4">
                            <ObjectList
                                items={children.map(child => ({
                                    id: child.data.id,
                                    name: child.data.name,
                                    description: child.metadata.classname !== 'DBNote' ? child.data.description : '',
                                    classname: child.metadata.classname
                                }))}
                                showViewToggle={true}
                                storageKey="siteNavigationChildrenViewMode"
                                defaultView="cards"
                                dark={dark}
                            />
                        </div>
                    )}
                </Col>
            </Row>
        </Container>
    );
}

export default SiteNavigation;
