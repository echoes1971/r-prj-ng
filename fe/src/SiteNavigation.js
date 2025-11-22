import React, { useState, useEffect, useContext } from 'react';
import { Container, Row, Col, Breadcrumb, Spinner, Alert, ListGroup } from 'react-bootstrap';
import { useParams, useNavigate, Link } from 'react-router-dom';
import axiosInstance from './axios';
import { app_cfg } from './app.cfg';
import ContentView from './ContentView';
import { getErrorMessage } from './errorHandler';
import { ThemeContext } from './ThemeContext';

// Format object ID: if 16 chars, format as xxxx-xxxxxxxx-xxxx
function formatObjectId(objId) {
    if (!objId) return objId;
    if (objId.length === 16) {
        return `${objId.slice(0, 4)}-${objId.slice(4, 12)}-${objId.slice(12, 16)}`;
    }
    return objId;
}

function SiteNavigation() {
    const { objectId } = useParams();
    const navigate = useNavigate();
    const { dark, themeClass } = useContext(ThemeContext);
    
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
                    {/* Breadcrumb */}
                    {breadcrumb.length > 0 && (
                        <Breadcrumb className="mb-3">
                            {breadcrumb.map((item, index) => (
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
                    )}

                    {/* Main Content */}
                    <ContentView 
                        data={content.data} 
                        metadata={content.metadata}
                        dark={dark}
                    />

                    {/* Children List (if folder) */}
                    {children.length > 0 && content.metadata.classname === 'DBFolder' && (
                        <div className="mt-4">
                            {/* <h4>Contents</h4> */}
                            <ListGroup variant={dark ? 'dark' : undefined}>
                                {children.map((child) => (
                                    <ListGroup.Item
                                        key={child.data.id}
                                        action
                                        onClick={() => handleNavigate(child.data.id)}
                                        style={{ cursor: 'pointer' }}
                                        variant={dark ? 'dark' : undefined}
                                    >
                                        <div className="d-flex justify-content-between align-items-center">
                                            <div>
                                                <strong>{child.data.name}</strong>
                                                {child.data.description && (
                                                    <div className="small" style={{ opacity: 0.7 }}>
                                                        {child.data.description}
                                                    </div>
                                                )}
                                            </div>
                                            <span className="badge bg-secondary">
                                                {child.metadata.classname}
                                            </span>
                                        </div>
                                    </ListGroup.Item>
                                ))}
                            </ListGroup>
                        </div>
                    )}
                </Col>
            </Row>
        </Container>
    );
}

export default SiteNavigation;
