import React, { useState, useEffect } from 'react';
import { Card, Container, Form, Button, Spinner, Alert, ButtonGroup } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import ReactDOM from 'react-dom';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import axiosInstance from './axios';

// Polyfill for findDOMNode (removed in React 19)
if (!ReactDOM.findDOMNode) {
    ReactDOM.findDOMNode = (node) => {
        if (node == null) return null;
        if (node instanceof HTMLElement) return node;
        if (node._reactInternals?.stateNode instanceof HTMLElement) {
            return node._reactInternals.stateNode;
        }
        console.warn('findDOMNode fallback used');
        return null;
    };
}

// Edit form for DBNote
function NoteEdit({ data, onSave, onCancel, saving, error }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(formData);
    };

    return (
        <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3">
                <Form.Label>{t('common.name')}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description')}</Form.Label>
                <Form.Control
                    as="textarea"
                    name="description"
                    rows={10}
                    value={formData.description}
                    onChange={handleChange}
                />
            </Form.Group>

            {error && (
                <Alert variant="danger" className="mb-3">
                    {error}
                </Alert>
            )}

            <div className="d-flex gap-2">
                <Button 
                    variant="primary" 
                    type="submit"
                    disabled={saving}
                >
                    {saving ? (
                        <>
                            <Spinner
                                as="span"
                                animation="border"
                                size="sm"
                                role="status"
                                aria-hidden="true"
                                className="me-2"
                            />
                            {t('common.saving')}
                        </>
                    ) : (
                        <>
                            <i className="bi bi-check-lg me-1"></i>
                            {t('common.save')}
                        </>
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
                    <i className="bi bi-x-lg me-1"></i>
                    {t('common.cancel')}
                </Button>
            </div>
        </Form>
    );
}

// Edit form for DBPage
function PageEdit({ data, onSave, onCancel, saving, error }) {
    const { t } = useTranslation();
    const [htmlMode, setHtmlMode] = useState('wysiwyg'); // 'wysiwyg' or 'source'
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        html: data.html || '',
        language: data.language || 'en',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(formData);
    };

    return (
        <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3">
                <Form.Label>{t('common.name')}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description')}</Form.Label>
                <Form.Control
                    as="textarea"
                    name="description"
                    rows={3}
                    value={formData.description}
                    onChange={handleChange}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Language</Form.Label>
                <Form.Select
                    name="language"
                    value={formData.language}
                    onChange={handleChange}
                >
                    <option value="en">English</option>
                    <option value="it">Italiano</option>
                    <option value="de">Deutsch</option>
                    <option value="fr">Français</option>
                </Form.Select>
            </Form.Group>

            <Form.Group className="mb-3">
                <div className="d-flex justify-content-between align-items-center mb-2">
                    <Form.Label className="mb-0">HTML Content</Form.Label>
                    <ButtonGroup size="sm">
                        <Button 
                            variant={htmlMode === 'wysiwyg' ? 'primary' : 'outline-primary'}
                            onClick={() => setHtmlMode('wysiwyg')}
                        >
                            <i className="bi bi-eye me-1"></i>WYSIWYG
                        </Button>
                        <Button 
                            variant={htmlMode === 'source' ? 'primary' : 'outline-primary'}
                            onClick={() => setHtmlMode('source')}
                        >
                            <i className="bi bi-code-slash me-1"></i>HTML Source
                        </Button>
                    </ButtonGroup>
                </div>
                {htmlMode === 'wysiwyg' ? (
                    <ReactQuill 
                        value={formData.html}
                        onChange={(value) => setFormData(prev => ({...prev, html: value}))}
                        theme="snow"
                        modules={{
                            toolbar: [
                                [{ 'header': [1, 2, 3, false] }],
                                ['bold', 'italic', 'underline', 'strike'],
                                [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                                [{ 'indent': '-1'}, { 'indent': '+1' }],
                                ['link', 'image'],
                                ['clean']
                            ]
                        }}
                    />
                ) : (
                    <Form.Control
                        as="textarea"
                        name="html"
                        value={formData.html}
                        onChange={handleChange}
                        rows={15}
                        style={{ fontFamily: 'monospace', fontSize: '0.9em' }}
                    />
                )}
                <Form.Text className="text-muted">
                    HTML content for the page
                </Form.Text>
            </Form.Group>

            {error && (
                <Alert variant="danger" className="mb-3">
                    {error}
                </Alert>
            )}

            <div className="d-flex gap-2">
                <Button 
                    variant="primary" 
                    type="submit"
                    disabled={saving}
                >
                    {saving ? (
                        <>
                            <Spinner
                                as="span"
                                animation="border"
                                size="sm"
                                role="status"
                                aria-hidden="true"
                                className="me-2"
                            />
                            {t('common.saving')}
                        </>
                    ) : (
                        <>
                            <i className="bi bi-check-lg me-1"></i>
                            {t('common.save')}
                        </>
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
                    <i className="bi bi-x-lg me-1"></i>
                    {t('common.cancel')}
                </Button>
            </div>
        </Form>
    );
}

// Generic edit form for other DBObjects
function ObjectEdit({ data, metadata, onSave, onCancel, saving, error }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(formData);
    };

    return (
        <Form onSubmit={handleSubmit}>
            <Alert variant="info" className="mb-3">
                <i className="bi bi-info-circle me-2"></i>
                Editing {metadata.classname} - Basic fields only
            </Alert>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.name')}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description')}</Form.Label>
                <Form.Control
                    as="textarea"
                    name="description"
                    rows={10}
                    value={formData.description}
                    onChange={handleChange}
                />
            </Form.Group>

            {error && (
                <Alert variant="danger" className="mb-3">
                    {error}
                </Alert>
            )}

            <div className="d-flex gap-2">
                <Button 
                    variant="primary" 
                    type="submit"
                    disabled={saving}
                >
                    {saving ? (
                        <>
                            <Spinner
                                as="span"
                                animation="border"
                                size="sm"
                                role="status"
                                aria-hidden="true"
                                className="me-2"
                            />
                            {t('common.saving')}
                        </>
                    ) : (
                        <>
                            <i className="bi bi-check-lg me-1"></i>
                            {t('common.save')}
                        </>
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
                    <i className="bi bi-x-lg me-1"></i>
                    {t('common.cancel')}
                </Button>
            </div>
        </Form>
    );
}

// Main ContentEdit component
function ContentEdit({ dark }) {
    const { id } = useParams();
    const navigate = useNavigate();
    const { t } = useTranslation();

    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState(null);
    const [data, setData] = useState(null);
    const [metadata, setMetadata] = useState(null);

    useEffect(() => {
        const loadObject = async () => {
            try {
                setLoading(true);
                setError(null);

                // Load object data
                const response = await axiosInstance.get(`/content/${id}`);
                setData(response.data.data);
                setMetadata(response.data.metadata);
            } catch (err) {
                console.error('Error loading object:', err);
                setError(err.response?.data?.message || 'Failed to load object');
            } finally {
                setLoading(false);
            }
        };

        loadObject();
    }, [id]);

    const handleSave = async (formData) => {
        try {
            setSaving(true);
            setError(null);

            // Update object via API
            await axiosInstance.put(`/objects/${id}`, formData);

            // Navigate back to view mode
            navigate(`/c/${id}`);
        } catch (err) {
            console.error('Error saving object:', err);
            setError(err.response?.data?.message || 'Failed to save changes');
            setSaving(false);
        }
    };

    const handleCancel = () => {
        navigate(`/c/${id}`);
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

    if (error && !data) {
        return (
            <Container className="mt-4">
                <Alert variant="danger">
                    <Alert.Heading>Error</Alert.Heading>
                    <p>{error}</p>
                    <Button variant="outline-danger" onClick={() => navigate(-1)}>
                        Go Back
                    </Button>
                </Alert>
            </Container>
        );
    }

    if (!data || !metadata) {
        return null;
    }

    // Check if user can edit
    if (!metadata.can_edit) {
        return (
            <Container className="mt-4">
                <Alert variant="warning">
                    <Alert.Heading>Access Denied</Alert.Heading>
                    <p>You don't have permission to edit this object.</p>
                    <Button variant="outline-warning" onClick={() => navigate(`/c/${id}`)}>
                        View Object
                    </Button>
                </Alert>
            </Container>
        );
    }

    const classname = metadata.classname;

    // Render appropriate edit form based on classname
    let EditComponent;
    switch (classname) {
        case 'DBNote':
            EditComponent = NoteEdit;
            break;
        case 'DBPage':
            EditComponent = PageEdit;
            break;
        default:
            EditComponent = ObjectEdit;
            break;
    }

    return (
        <Container className="mt-4">
            <Card bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
                <Card.Header>
                    <h2 className={dark ? 'text-light' : 'text-dark'}>
                        <i className="bi bi-pencil me-2"></i>
                        {t('common.edit')}: {data.name}
                    </h2>
                    <small style={{ opacity: 0.7 }}>
                        {classname} · ID: {id}
                    </small>
                </Card.Header>
                <Card.Body>
                    <EditComponent
                        data={data}
                        metadata={metadata}
                        onSave={handleSave}
                        onCancel={handleCancel}
                        saving={saving}
                        error={error}
                    />
                </Card.Body>
            </Card>
        </Container>
    );
}

export default ContentEdit;
