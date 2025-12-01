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
import FileSelector from './FileSelector';
import { ObjectEdit } from './DBObject';
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

// Helper functions for DBFile token management in HTML content

/**
 * Extract all file IDs from HTML content that have data-dbfile-id attribute
 */
function extractFileIDs(html) {
    if (!html) return [];
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    const elements = doc.querySelectorAll('[data-dbfile-id]');
    const fileIDs = new Set();
    elements.forEach(el => {
        const fileId = el.getAttribute('data-dbfile-id');
        if (fileId && fileId !== '0') {
            fileIDs.add(fileId);
        }
    });
    return Array.from(fileIDs);
}

/**
 * Request temporary tokens for multiple file IDs
 */
async function requestFileTokens(fileIDs) {
    if (!fileIDs || fileIDs.length === 0) return {};
    
    try {
        const response = await axiosInstance.post('/files/preview-tokens', {
            file_ids: fileIDs
        });
        return response.data.tokens || {};
    } catch (error) {
        console.error('Failed to request file tokens:', error);
        return {};
    }
}

/**
 * Inject tokens into HTML for WYSIWYG editing
 * Adds ?token=... to src/href attributes of elements with data-dbfile-id
 */
function injectTokensForEditing(html, tokens) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    
    doc.querySelectorAll('[data-dbfile-id]').forEach(el => {
        const fileId = el.getAttribute('data-dbfile-id');
        const token = tokens[fileId];
        
        if (token) {
            if (el.tagName === 'IMG') {
                // Remove existing token if any
                const currentSrc = el.src || el.getAttribute('src') || '';
                const baseUrl = currentSrc.split('?')[0];
                el.setAttribute('src', `${baseUrl}?token=${token}`);
            } else if (el.tagName === 'A') {
                const currentHref = el.href || el.getAttribute('href') || '';
                const baseUrl = currentHref.split('?')[0];
                el.setAttribute('href', `${baseUrl}?token=${token}`);
            }
        }
    });
    
    return doc.body.innerHTML;
}

/**
 * Clean tokens from HTML before saving
 * Removes ?token=... from src/href but keeps data-dbfile-id
 */
function cleanTokensBeforeSave(html) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    
    doc.querySelectorAll('[data-dbfile-id]').forEach(el => {
        const fileId = el.getAttribute('data-dbfile-id');
        
        if (el.tagName === 'IMG') {
            el.setAttribute('src', `/api/files/${fileId}/download`);
        } else if (el.tagName === 'A') {
            el.setAttribute('href', `/api/files/${fileId}/download`);
        }
    });
    
    return doc.body.innerHTML;
}

// Edit form for DBNote
function NoteEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        fk_obj_id: data.fk_obj_id || '0',
        father_id: data.father_id || '0',
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

            <ObjectLinkSelector
                value={formData.father_id || '0'}
                onChange={handleChange}
                classname="DBObject"
                fieldName="father_id"
                label={t('dbobjects.parent')}
            />

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

            <ObjectLinkSelector
                value={formData.fk_obj_id}
                onChange={handleChange}
                classname="DBObject"
                fieldName="fk_obj_id"
                label={t('files.linked_object') || 'Linked Object'}
                required={false}
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

// Edit form for DBPage
function PageEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [htmlMode, setHtmlMode] = useState('wysiwyg'); // 'wysiwyg' or 'source'
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        html: data.html || '',
        language: data.language || 'en',
        father_id: data.father_id || '0',
        fk_obj_id: data.fk_obj_id || '0',
    });
    const [htmlWithTokens, setHtmlWithTokens] = useState(data.html || '');
    const [loadingTokens, setLoadingTokens] = useState(false);
    const [showFileSelector, setShowFileSelector] = useState(false);
    const [fileSelectorType, setFileSelectorType] = useState('file'); // 'file' or 'image'
    const [quillRef, setQuillRef] = useState(null);

    // Load tokens for embedded files when component mounts or HTML changes
    useEffect(() => {
        const loadTokens = async () => {
            const fileIDs = extractFileIDs(formData.html);
            if (fileIDs.length === 0) {
                setHtmlWithTokens(formData.html);
                return;
            }

            setLoadingTokens(true);
            try {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = injectTokensForEditing(formData.html, tokens);
                setHtmlWithTokens(htmlWithTokens);
            } catch (error) {
                console.error('Failed to load tokens for embedded files:', error);
                setHtmlWithTokens(formData.html);
            } finally {
                setLoadingTokens(false);
            }
        };

        loadTokens();
    }, [data.id]); // Only reload when page ID changes

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleHtmlChange = (value) => {
        setFormData(prev => ({...prev, html: value}));
        setHtmlWithTokens(value);
    };

    // Handle file selection from modal
    const handleFileSelect = (file) => {
        if (!quillRef) return;

        const quill = quillRef.getEditor();
        const range = quill.getSelection(true);
        
        if (fileSelectorType === 'image') {
            // Insert image tag
            const imageHtml = `<img src="/api/files/${file.id}/download" data-dbfile-id="${file.id}" alt="${file.name}" style="max-width: 100%;" />`;
            quill.clipboard.dangerouslyPasteHTML(range.index, imageHtml);
        } else {
            // Insert link
            const linkHtml = `<a href="/api/files/${file.id}/download" data-dbfile-id="${file.id}">${file.name}</a>`;
            quill.clipboard.dangerouslyPasteHTML(range.index, linkHtml);
        }
        
        // Update state
        handleHtmlChange(quill.root.innerHTML);
    };

    // Open file selector modal
    const handleInsertFile = () => {
        setFileSelectorType('file');
        setShowFileSelector(true);
    };

    const handleInsertImage = () => {
        setFileSelectorType('image');
        setShowFileSelector(true);
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        // Clean tokens before saving
        const cleanedHtml = cleanTokensBeforeSave(formData.html);
        onSave({
            ...formData,
            html: cleanedHtml
        });
    };

    return (
        <Form onSubmit={handleSubmit}>

            <ObjectLinkSelector
                value={formData.father_id || '0'}
                onChange={handleChange}
                classname="DBObject"
                fieldName="father_id"
                label={t('dbobjects.parent')}
            />

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
                
                {/* Custom buttons for inserting files/images */}
                {htmlMode === 'wysiwyg' && !loadingTokens && (
                    <div className="mb-2">
                        <ButtonGroup size="sm">
                            <Button 
                                variant="outline-secondary"
                                onClick={handleInsertImage}
                                title={t('files.insert_image') || 'Insert Image'}
                            >
                                <i className="bi bi-image me-1"></i>
                                {t('files.insert_image') || 'Insert Image'}
                            </Button>
                            <Button 
                                variant="outline-secondary"
                                onClick={handleInsertFile}
                                title={t('files.insert_file') || 'Insert File'}
                            >
                                <i className="bi bi-paperclip me-1"></i>
                                {t('files.insert_file') || 'Insert File'}
                            </Button>
                        </ButtonGroup>
                        <Form.Text className="text-muted ms-2">
                            Insert files/images you have permission to edit
                        </Form.Text>
                    </div>
                )}
                
                {loadingTokens && (
                    <div className="text-center py-3">
                        <Spinner animation="border" size="sm" className="me-2" />
                        <span>Loading file tokens...</span>
                    </div>
                )}
                {!loadingTokens && htmlMode === 'wysiwyg' ? (
                    <ReactQuill 
                        ref={setQuillRef}
                        value={htmlWithTokens}
                        onChange={handleHtmlChange}
                        theme="snow"
                        modules={{
                            toolbar: [
                                [{ 'header': [1, 2, 3, false] }],
                                ['bold', 'italic', 'underline', 'strike'],
                                [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                                [{ 'indent': '-1'}, { 'indent': '+1' }],
                                ['link'], // Removed 'image' - use custom buttons instead
                                ['clean']
                            ]
                        }}
                    />
                ) : !loadingTokens ? (
                    <Form.Control
                        as="textarea"
                        name="html"
                        value={formData.html}
                        onChange={(e) => handleHtmlChange(e.target.value)}
                        rows={15}
                        style={{ fontFamily: 'monospace', fontSize: '0.9em' }}
                    />
                ) : null}
                <Form.Text className="text-muted">
                    HTML content for the page. Use data-dbfile-id attribute to embed files (e.g., &lt;img src="/api/files/ID/download" data-dbfile-id="ID" /&gt;)
                </Form.Text>
            </Form.Group>

            <ObjectLinkSelector
                value={formData.fk_obj_id}
                onChange={handleChange}
                classname="DBObject"
                fieldName="fk_obj_id"
                label={t('files.linked_object') || 'Linked Object'}
                required={false}
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
            
            {/* File Selector Modal */}
            <FileSelector
                show={showFileSelector}
                onHide={() => setShowFileSelector(false)}
                onSelect={handleFileSelect}
                fileType={fileSelectorType}
            />
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
        father_id: data.father_id || '0',
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

            <ObjectLinkSelector
                value={formData.father_id || '0'}
                onChange={handleChange}
                classname="DBObject"
                fieldName="father_id"
                label={t('dbobjects.parent')}
            />

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
        father_id: data.father_id || '0',
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

            <ObjectLinkSelector
                value={formData.father_id || '0'}
                onChange={handleChange}
                classname="DBObject"
                fieldName="father_id"
                label={t('dbobjects.parent')}
            />

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
        father_id: data.father_id || '0',
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
                {/* <Form.Label>{t('dbobjects.parent')}</Form.Label> */}
                <ObjectLinkSelector
                    value={formData.father_id || '0'}
                    onChange={handleChange}
                    classname="DBObject"
                    fieldName="father_id"
                    label={t('dbobjects.parent')}
                />
            </Form.Group>

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
                classname="DBObject"
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

// Edit form for DBFolder
function FolderEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        father_id: data.father_id || '0',
        name: data.name || '',
        description: data.description || '',
        fk_obj_id: data.fk_obj_id || '0',
        permissions: data.permissions || 'rwxr-x---',
        childs_sort_order: data.childs_sort_order || '',
        father_id: data.father_id || '0',
    });
    const [children, setChildren] = useState([]);
    const [loadingChildren, setLoadingChildren] = useState(false);
    const [sortedChildrenIds, setSortedChildrenIds] = useState([]);
    const [draggedIndex, setDraggedIndex] = useState(null);

    // Load children on mount
    useEffect(() => {
        if (data.id) {
            loadChildren();
        }
    }, [data.id]);

    const loadChildren = async () => {
        setLoadingChildren(true);
        try {
            const response = await axiosInstance.get(`/nav/children/${data.id}`);
            const childrenData = response.data.children || [];
            setChildren(childrenData);
            // console.log('Children data:', childrenData);
            
            // Initialize sorted order from childs_sort_order or use current order
            if (formData.childs_sort_order) {
                const orderIds = formData.childs_sort_order.split(',').filter(id => id);
                setSortedChildrenIds(orderIds);
            // } else {
            //     setSortedChildrenIds(childrenData.map(child => child.data.id));
            }
            console.log('Initial sortedChildrenIds:', sortedChildrenIds);
        } catch (error) {
            console.error('Failed to load children:', error);
        } finally {
            setLoadingChildren(false);
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleDragStart = (e, index) => {
        setDraggedIndex(index);
        e.dataTransfer.effectAllowed = 'move';
    };

    const handleDragOver = (e, index) => {
        e.preventDefault();
        if (draggedIndex === null || draggedIndex === index) return;

        const newOrder = [...sortedChildrenIds];
        const draggedItem = newOrder[draggedIndex];
        newOrder.splice(draggedIndex, 1);
        newOrder.splice(index, 0, draggedItem);

        setSortedChildrenIds(newOrder);
        setDraggedIndex(index);
    };

    const handleDragEnd = () => {
        setDraggedIndex(null);
        // Update formData with new order
        setFormData(prev => ({
            ...prev,
            childs_sort_order: sortedChildrenIds.join(',')
        }));
    };

    const toggleChildInOrder = (childId) => {
        const newOrder = sortedChildrenIds.includes(childId)
            ? sortedChildrenIds.filter(id => id !== childId)
            : [...sortedChildrenIds, childId];
        
        setSortedChildrenIds(newOrder);
        setFormData(prev => ({
            ...prev,
            childs_sort_order: newOrder.join(',')
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(formData);
    };

    // Get child name by ID
    const getChildName = (childId) => {
        const child = children.find(c => c.data.id === childId);
        return child ? child.data.name : childId;
    };

    return (
        <Form onSubmit={handleSubmit}>

            <Form.Group className="mb-3">
                {/* <Form.Label>{t('dbobjects.parent')}</Form.Label> */}
                <ObjectLinkSelector
                    value={formData.father_id || '0'}
                    onChange={handleChange}
                    classname="DBObject"
                    fieldName="father_id"
                    name="father_id"
                    label={t('dbobjects.parent')}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.name')}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    disabled={saving}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description')}</Form.Label>
                <Form.Control
                    as="textarea"
                    rows={3}
                    name="description"
                    value={formData.description}
                    onChange={handleChange}
                    disabled={saving}
                />
            </Form.Group>
<h2 className='text-danger'>TODO: Edit the index page, automatically create if not existing for the current language.</h2>
            <Form.Group className="mb-3">
                <Form.Label>{t('files.linked_object')}</Form.Label>
                <ObjectLinkSelector
                    value={formData.fk_obj_id || '0'}
                    // onChange={(value) => setFormData(prev => ({ ...prev, fk_obj_id: value }))}
                    onChange={handleChange}
                    name="fk_obj_id"
                    fieldName="fk_obj_id"
                    // disabled={saving}
                    classname="DBObject"
                    // allowedTypes={['DBPage', 'DBNews']}
                />
            </Form.Group>

            {/* Children Sort Order */}
            {children.length > 0 && (
                <Form.Group className="mb-3">
                    <Form.Label>
                        {t('folder.children_order')}
                        <small className="ms-2 text-muted">
                            ({t('folder.drag_to_reorder')})
                        </small>
                    </Form.Label>
                    
                    {loadingChildren ? (
                        <div className="text-center p-3">
                            <Spinner animation="border" size="sm" />
                        </div>
                    ) : (
                        <>
                            {/* List of sorted children (draggable) */}
                            <div className={`border rounded p-2 mb-2 ${dark ? 'border-secondary' : ''}`}>
                                {sortedChildrenIds.length === 0 ? (
                                    <div className="text-muted text-center p-2">
                                        {t('folder.no_children_selected')}
                                    </div>
                                ) : (
                                    sortedChildrenIds.map((childId, index) => (
                                        <div
                                            key={childId}
                                            draggable
                                            onDragStart={(e) => handleDragStart(e, index)}
                                            onDragOver={(e) => handleDragOver(e, index)}
                                            onDragEnd={handleDragEnd}
                                            className={`d-flex align-items-center p-2 mb-1 rounded ${
                                                dark ? 'bg-dark' : 'bg-light'
                                            } ${draggedIndex === index ? 'opacity-50' : ''}`}
                                            style={{ cursor: 'move' }}
                                        >
                                            <i className="bi bi-grip-vertical me-2"></i>
                                            <span className="flex-grow-1">{getChildName(childId)}</span>
                                            <Button
                                                variant="outline-danger"
                                                size="sm"
                                                onClick={() => toggleChildInOrder(childId)}
                                                disabled={saving}
                                            >
                                                <i className="bi bi-x"></i>
                                            </Button>
                                        </div>
                                    ))
                                )}
                            </div>

                            {/* List of available children (not in sort order) */}
                            {children.filter(child => !sortedChildrenIds.includes(child.data.id)).length > 0 && (
                                <>
                                    <Form.Label className="mt-2">
                                        {t('folder.available_children')}
                                    </Form.Label>
                                    <div className={`border rounded p-2 ${dark ? 'border-secondary' : ''}`}>
                                        {children
                                            .filter(child => !sortedChildrenIds.includes(child.data.id))
                                            .map(child => (
                                                <div
                                                    key={child.data.id}
                                                    className={`d-flex align-items-center p-2 mb-1 rounded ${
                                                        dark ? 'bg-secondary bg-opacity-25' : 'bg-light'
                                                    }`}
                                                >
                                                    <span className="flex-grow-1">{child.data.name}</span>
                                                    <Button
                                                        variant="outline-primary"
                                                        size="sm"
                                                        onClick={() => toggleChildInOrder(child.data.id)}
                                                        disabled={saving}
                                                    >
                                                        <i className="bi bi-plus"></i>
                                                    </Button>
                                                </div>
                                            ))}
                                    </div>
                                </>
                            )}
                        </>
                    )}
                </Form.Group>
            )}

            <Form.Group className="mb-3">
                <Form.Label>{t('common.permissions')}</Form.Label>
                <PermissionsEditor
                    value={formData.permissions}
                    onChange={(value) => setFormData(prev => ({ ...prev, permissions: value }))}
                    disabled={saving}
                    dark={dark}
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
                            <Spinner animation="border" size="sm" className="me-2" />
                            {t('common.saving')}
                        </>
                    ) : (
                        t('common.save')
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
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
        case 'DBNews':
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
        case 'DBFolder':
            EditComponent = FolderEdit;
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
