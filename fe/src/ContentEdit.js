import React, { useState, useEffect, useContext } from 'react';
import { Card, Container, Form, Button, Spinner, Alert, ButtonGroup } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import ReactDOM from 'react-dom';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import axiosInstance from './axios';
import CountrySelector from './CountrySelector';
import ObjectLinkSelector from './ObjectLinkSelector';
import PermissionsEditor from './PermissionsEditor';
import { ThemeContext } from './ThemeContext';

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
function NoteEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Edit form for DBPage
function PageEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Edit form for DBPerson
function PersonEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        street: data.street || '',
        zip: data.zip || '',
        city: data.city || '',
        state: data.state || '',
        fk_countrylist_id: data.fk_countrylist_id || '0',
        phone: data.phone || '',
        mobile: data.mobile || '',
        office_phone: data.office_phone || '',
        fax: data.fax || '',
        email: data.email || '',
        url: data.url || '',
        codice_fiscale: data.codice_fiscale || '',
        p_iva: data.p_iva || '',
        fk_users_id: data.fk_users_id || '0',
        fk_companies_id: data.fk_companies_id || '0',
        permissions: data.permissions || 'rwxr-x---',
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

            <h5 className="mt-4 mb-3">Address</h5>

            <Form.Group className="mb-3">
                <Form.Label>Street</Form.Label>
                <Form.Control
                    type="text"
                    name="street"
                    value={formData.street}
                    onChange={handleChange}
                />
            </Form.Group>

            <div className="row">
                <div className="col-md-4">
                    <Form.Group className="mb-3">
                        <Form.Label>ZIP Code</Form.Label>
                        <Form.Control
                            type="text"
                            name="zip"
                            value={formData.zip}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-8">
                    <Form.Group className="mb-3">
                        <Form.Label>City</Form.Label>
                        <Form.Control
                            type="text"
                            name="city"
                            value={formData.city}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.state')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="state"
                            value={formData.state}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <CountrySelector
                        value={formData.fk_countrylist_id}
                        onChange={handleChange}
                        name="fk_countrylist_id"
                    />
                </div>
            </div>

            <h5 className="mt-4 mb-3">{t('common.contact_info')}</h5>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.phone')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="phone"
                            value={formData.phone}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.mobile')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="mobile"
                            value={formData.mobile}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.office_phone')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="office_phone"
                            value={formData.office_phone}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.fax')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="fax"
                            value={formData.fax}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.email')}</Form.Label>
                <Form.Control
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.website')}</Form.Label>
                <Form.Control
                    type="url"
                    name="url"
                    value={formData.url}
                    onChange={handleChange}
                />
            </Form.Group>

            <h5 className="mt-4 mb-3">{t('common.tax_info')}</h5>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.codice_fiscale')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="codice_fiscale"
                            value={formData.codice_fiscale}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.p_iva')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="p_iva"
                            value={formData.p_iva}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <h5 className="mt-4 mb-3">{t('common.links')}</h5>

            <ObjectLinkSelector
                value={formData.fk_users_id}
                onChange={handleChange}
                classname="DBUser"
                fieldName="fk_users_id"
                label={t('common.user_id')}
                required={false}
            />

            <ObjectLinkSelector
                value={formData.fk_companies_id}
                onChange={handleChange}
                classname="DBCompany"
                fieldName="fk_companies_id"
                label={t('common.company_id')}
                required={false}
            />

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Edit form for DBCompany
function CompanyEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        street: data.street || '',
        zip: data.zip || '',
        city: data.city || '',
        state: data.state || '',
        fk_countrylist_id: data.fk_countrylist_id || '0',
        phone: data.phone || '',
        fax: data.fax || '',
        email: data.email || '',
        url: data.url || '',
        p_iva: data.p_iva || '',
        permissions: data.permissions || 'rwxr-x---',
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

            <h5 className="mt-4 mb-3">{t('common.address')}</h5>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.street')}</Form.Label>
                <Form.Control
                    type="text"
                    name="street"
                    value={formData.street}
                    onChange={handleChange}
                />
            </Form.Group>

            <div className="row">
                <div className="col-md-4">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.zip')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="zip"
                            value={formData.zip}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-8">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.city')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="city"
                            value={formData.city}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.state')}</Form.Label>
                        <Form.Control
                            type="text"
                            name="state"
                            value={formData.state}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <CountrySelector
                        value={formData.fk_countrylist_id}
                        onChange={handleChange}
                        name="fk_countrylist_id"
                    />
                </div>
            </div>

            <h5 className="mt-4 mb-3">{t('common.contact_info')}</h5>

            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.phone')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="phone"
                            value={formData.phone}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('common.fax')}</Form.Label>
                        <Form.Control
                            type="tel"
                            name="fax"
                            value={formData.fax}
                            onChange={handleChange}
                        />
                    </Form.Group>
                </div>
            </div>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.email')}</Form.Label>
                <Form.Control
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.website')}</Form.Label>
                <Form.Control
                    type="url"
                    name="url"
                    value={formData.url}
                    onChange={handleChange}
                />
            </Form.Group>

            <h5 className="mt-4 mb-3">{t('common.tax_info')}</h5>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.p_iva')}</Form.Label>
                <Form.Control
                    type="text"
                    name="p_iva"
                    value={formData.p_iva}
                    onChange={handleChange}
                />
            </Form.Group>

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Edit form for DBFile
function FileEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        filename: data.filename || '',
        mime: data.mime || '',
        alt_link: data.alt_link || '',
        fk_obj_id: data.fk_obj_id || '0',
        permissions: data.permissions || 'rwxr-x---',
    });
    const [selectedFile, setSelectedFile] = useState(null);
    const [dragActive, setDragActive] = useState(false);
    const [preview, setPreview] = useState(null);

    useEffect(() => {
        console.log('FileEdit useEffect:', { id: data.id, filename: data.filename, selectedFile });
        // Load preview if file exists
        if (data.id && data.filename && !selectedFile) {
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
    }, [data.id, data.filename, selectedFile]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setSelectedFile(file);
            setFormData(prev => ({
                ...prev,
                filename: file.name,
                mime: file.type
            }));

            // Create preview for images
            if (file.type.startsWith('image/')) {
                const reader = new FileReader();
                reader.onloadend = () => {
                    setPreview(reader.result);
                };
                reader.readAsDataURL(file);
            } else {
                setPreview(null);
            }
        }
    };

    const handleDrag = (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === 'dragenter' || e.type === 'dragover') {
            setDragActive(true);
        } else if (e.type === 'dragleave') {
            setDragActive(false);
        }
    };

    const handleDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);

        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            const file = e.dataTransfer.files[0];
            setSelectedFile(file);
            setFormData(prev => ({
                ...prev,
                filename: file.name,
                mime: file.type
            }));

            // Create preview for images
            if (file.type.startsWith('image/')) {
                const reader = new FileReader();
                reader.onloadend = () => {
                    setPreview(reader.result);
                };
                reader.readAsDataURL(file);
            } else {
                setPreview(null);
            }
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (selectedFile) {
            // Upload new file
            const uploadFormData = new FormData();
            uploadFormData.append('file', selectedFile);
            uploadFormData.append('name', formData.name);
            uploadFormData.append('description', formData.description);
            uploadFormData.append('alt_link', formData.alt_link);
            uploadFormData.append('fk_obj_id', formData.fk_obj_id);
            uploadFormData.append('permissions', formData.permissions);

            onSave(uploadFormData, true); // Pass true to indicate multipart upload
        } else {
            // Update metadata only
            onSave(formData);
        }
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

            {/* File Upload */}
            <Form.Group className="mb-3">
                <Form.Label>{t('files.upload') || 'Upload File'}</Form.Label>
                <div
                    className={`border rounded p-4 text-center ${dragActive ? 'border-primary bg-primary bg-opacity-10' : ''} ${dark ? 'bg-dark border-secondary' : 'bg-light'}`}
                    onDragEnter={handleDrag}
                    onDragLeave={handleDrag}
                    onDragOver={handleDrag}
                    onDrop={handleDrop}
                    style={{ cursor: 'pointer' }}
                    onClick={() => document.getElementById('fileInput').click()}
                >
                    <input
                        id="fileInput"
                        type="file"
                        onChange={handleFileChange}
                        style={{ display: 'none' }}
                    />
                    {preview ? (
                        <div>
                            <img 
                                src={preview} 
                                alt="Preview" 
                                style={{ maxWidth: '100%', maxHeight: '300px', marginBottom: '10px' }}
                            />
                            <div>
                                <small className="text-muted">{formData.filename}</small>
                            </div>
                        </div>
                    ) : (
                        <>
                            <i className="bi bi-cloud-upload fs-1"></i>
                            <p className="mb-0">
                                {selectedFile ? selectedFile.name : (t('files.drop_or_click') || 'Drop file here or click to browse')}
                            </p>
                            {formData.filename && !selectedFile && (
                                <small className="text-muted d-block mt-2">
                                    {t('files.current') || 'Current'}: {formData.filename}
                                </small>
                            )}
                        </>
                    )}
                </div>
                <Form.Text className="text-muted">
                    {t('files.hint') || 'Drag and drop a file or click to browse'}
                </Form.Text>
            </Form.Group>

            {/* File Metadata */}
            <div className="row">
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('files.filename') || 'Filename'}</Form.Label>
                        <Form.Control
                            type="text"
                            name="filename"
                            value={formData.filename}
                            onChange={handleChange}
                            readOnly
                        />
                    </Form.Group>
                </div>
                <div className="col-md-6">
                    <Form.Group className="mb-3">
                        <Form.Label>{t('files.mime_type') || 'MIME Type'}</Form.Label>
                        <Form.Control
                            type="text"
                            name="mime"
                            value={formData.mime}
                            onChange={handleChange}
                            readOnly
                        />
                    </Form.Group>
                </div>
            </div>

            <Form.Group className="mb-3">
                <Form.Label>{t('files.alt_link') || 'Alternative Link'}</Form.Label>
                <Form.Control
                    type="url"
                    name="alt_link"
                    value={formData.alt_link}
                    onChange={handleChange}
                />
                <Form.Text className="text-muted">
                    {t('files.alt_link_hint') || 'External URL if file is hosted elsewhere'}
                </Form.Text>
            </Form.Group>

            <ObjectLinkSelector
                value={formData.fk_obj_id}
                onChange={handleChange}
                classname="DBPage"
                fieldName="fk_obj_id"
                label={t('files.linked_object') || 'Linked Object'}
                required={false}
            />

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Generic edit form for other DBObjects
function ObjectEdit({ data, metadata, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        permissions: data.permissions || 'rwxr-x---',
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

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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
                <Button 
                    variant="outline-danger" 
                    onClick={onDelete}
                    disabled={saving}
                    className="ms-auto"
                >
                    <i className="bi bi-trash me-1"></i>
                    {t('common.delete')}
                </Button>
            </div>
        </Form>
    );
}

// Main ContentEdit component
function ContentEdit() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { dark } = useContext(ThemeContext);

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

    const handleSave = async (formData, isMultipart = false) => {
        try {
            setSaving(true);
            setError(null);

            if (isMultipart) {
                // Upload with multipart/form-data for file uploads
                await axiosInstance.put(`/objects/${id}`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            } else {
                // Regular JSON update
                await axiosInstance.put(`/objects/${id}`, formData);
            }

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

    const handleDelete = async () => {
        if (!window.confirm(t('navigation.delete_confirm'))) {
            return;
        }

        try {
            setSaving(true);
            setError(null);

            // Delete object via API
            await axiosInstance.delete(`/objects/${id}`);

            // Navigate to parent or home
            // alert("father_id=" + data.father_id);
            if (data.father_id) {
                navigate(`/c/${data.father_id}`);
                return;
            }
            navigate(-1);
        } catch (err) {
            console.error('Error deleting object:', err);
            setError(err.response?.data?.message || 'Failed to delete object');
            setSaving(false);
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
        case 'DBPerson':
            EditComponent = PersonEdit;
            break;
        case 'DBCompany':
            EditComponent = CompanyEdit;
            break;
        case 'DBFile':
            EditComponent = FileEdit;
            break;
        default:
            EditComponent = ObjectEdit;
            break;
    }

    return (
        <Container className="mt-4">
            <Card bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
                <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                    <h2 className={dark ? 'text-light' : 'text-dark'}>
                        <i className="bi bi-pencil me-2"></i>
                        {t('common.edit')}: {data.name}
                    </h2>
                    <small style={{ opacity: 0.7 }}>
                        {classname} · ID: {id}
                    </small>
                </Card.Header>
                <Card.Body className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '0px solid rgba(255,255,255,0.1)' } : {}}>
                    <EditComponent
                        data={data}
                        metadata={metadata}
                        onSave={handleSave}
                        onCancel={handleCancel}
                        onDelete={handleDelete}
                        saving={saving}
                        error={error}
                        dark={dark}
                    />
                </Card.Body>
            </Card>
        </Container>
    );
}

export default ContentEdit;
