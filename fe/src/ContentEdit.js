import React, { useState, useEffect, useContext } from 'react';
import { Card, Container, Form, Button, Spinner, Alert, ButtonGroup } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import axiosInstance from './axios';
import CountrySelector from './CountrySelector';
import ObjectLinkSelector from './ObjectLinkSelector';
import PermissionsEditor from './PermissionsEditor';
import FileSelector from './FileSelector';
import { FileEdit } from './DBFile';
import { NoteEdit } from './DBNote';
import { ObjectEdit } from './DBObject';
import { 
    cleanTokensBeforeSave,
    extractFileIDs,
    injectTokensForEditing,
    requestFileTokens,
    PageEdit } from './DBPage';
import { ThemeContext } from './ThemeContext';


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

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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
    const [htmlMode, setHtmlMode] = useState('wysiwyg'); // 'wysiwyg' or 'source'
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
    
    // Index page editor states
    const [indexPages, setIndexPages] = useState([]);
    const [selectedIndexLanguage, setSelectedIndexLanguage] = useState('en');
    const [indexHtml, setIndexHtml] = useState('');
    const [indexHtmlWithTokens, setIndexHtmlWithTokens] = useState('');
    const [loadingIndexTokens, setLoadingIndexTokens] = useState(false);
    const [savingIndex, setSavingIndex] = useState(false);
    const [quillRefIndex, setQuillRefIndex] = useState(null);
    const [showFileSelectorIndex, setShowFileSelectorIndex] = useState(false);
    const [fileSelectorTypeIndex, setFileSelectorTypeIndex] = useState('file');

    // Load children and index pages on mount
    useEffect(() => {
        if (data.id) {
            loadChildren();
            loadIndexPages();
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

    const loadIndexPages = async () => {
        try {
            const response = await axiosInstance.get(`/nav/${data.id}/indexes`);
            const pages = response.data.indexes || [];
            setIndexPages(pages);
            // Load HTML for current language if exists
            const currentPage = pages.find(p => p.data.language && p.data.language.indexOf(selectedIndexLanguage) >= 0);
            if (currentPage) {
                setIndexHtml(currentPage.data.html || '');
                await loadIndexTokens(currentPage.data.html || '');
            } else {
                setIndexHtml('');
                setIndexHtmlWithTokens('');
            }
        } catch (error) {
            console.error('Failed to load index pages:', error);
            setIndexPages([]);
        }
    };

    const loadIndexTokens = async (html) => {
        if (!html) {
            setIndexHtmlWithTokens('');
            return;
        }
        setLoadingIndexTokens(true);
        try {
            const fileIDs = extractFileIDs(html);
            if (fileIDs.length > 0) {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = injectTokensForEditing(html, tokens);
                setIndexHtmlWithTokens(htmlWithTokens);
            } else {
                setIndexHtmlWithTokens(html);
            }
        } catch (error) {
            console.error('Failed to load tokens for index:', error);
            setIndexHtmlWithTokens(html);
        } finally {
            setLoadingIndexTokens(false);
        }
    };

    // Reload index HTML when language changes
    useEffect(() => {
        const currentPage = indexPages.find(p => p.data.language && p.data.language.indexOf(selectedIndexLanguage) >= 0);
        if (currentPage) {
            setIndexHtml(currentPage.data.html || '');
            loadIndexTokens(currentPage.data.html || '');
        } else {
            setIndexHtml('');
            setIndexHtmlWithTokens('');
        }
    }, [selectedIndexLanguage, indexPages]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    // const handleIndexHtmlChange = (content) => {
    const handleIndexHtmlChange = async (content) => {
        // RRA: start.
        setIndexHtml(content);

        // Extract file IDs and reload tokens for immediate preview
        const fileIDs = extractFileIDs(content);
        if (fileIDs.length === 0) {
            setIndexHtmlWithTokens(content);
            return;
        }

        try {
            const tokens = await requestFileTokens(fileIDs);
            const htmlWithTokens = injectTokensForEditing(content, tokens);
            setIndexHtmlWithTokens(htmlWithTokens);
        } catch (error) {
            console.error('Failed to reload tokens after HTML change:', error);
            setIndexHtmlWithTokens(content);
        }

        // setIndexHtml(content);
        // setIndexHtmlWithTokens(content);
        // RRA: end.
    };

    const handleFileSelectIndex = (file) => {
        if (!quillRefIndex) return;
        
        const quill = quillRefIndex.getEditor();
        const range = quill.getSelection(true);
        
        if (fileSelectorTypeIndex === 'image') {
            const imgHtml = `<img src="/api/files/${file.id}/download" data-dbfile-id="${file.id}" alt="${file.name}" style="max-width: 100%;" />`;
            quill.clipboard.dangerouslyPasteHTML(range.index, imgHtml);
        } else {
            const linkHtml = `<a href="/api/files/${file.id}/download" data-dbfile-id="${file.id}">${file.name}</a>`;
            quill.clipboard.dangerouslyPasteHTML(range.index, linkHtml);
        }
        
        // Update state
        handleIndexHtmlChange(quill.root.innerHTML);
    };

    const handleInsertFileIndex = () => {
        setFileSelectorTypeIndex('file');
        setShowFileSelectorIndex(true);
    };

    const handleInsertImageIndex = () => {
        setFileSelectorTypeIndex('image');
        setShowFileSelectorIndex(true);
    };

    const handleSaveIndex = async () => {
        setSavingIndex(true);
        try {
            const cleanedHtml = cleanTokensBeforeSave(indexHtml);
            const currentPage = indexPages.find(p => p.data.language && p.data.language.indexOf(selectedIndexLanguage) >= 0);
            
            if (currentPage) {
                // Update existing page
                await axiosInstance.put(`/objects/${currentPage.data.id}`, {
                    ...currentPage.data,
                    html: cleanedHtml
                });
            } else {
                // Create new index page
                await axiosInstance.post('/objects', {
                    classname: 'DBPage',
                    name: 'index',
                    // description: `Index page for ${selectedIndexLanguage}`,
                    description: '',
                    language: selectedIndexLanguage,
                    html: cleanedHtml,
                    father_id: data.id,
                    permissions: data.permissions || ''
                });
            }
            
            // Reload index pages
            await loadIndexPages();
            alert(t('common.saved'));
        } catch (error) {
            console.error('Failed to save index:', error);
            alert(t('errors.save_failed'));
        } finally {
            setSavingIndex(false);
        }
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

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

            {/* Index Page Editor */}
            <div className="mb-4 p-3 border rounded">
                <h5>{t('folder.index_page_editor')}</h5>
                <p className="text-muted small">{t('folder.index_page_hint')}</p>
                
                <Form.Group className="mb-3">
                    <Form.Label>{t('common.language')}</Form.Label>
                    <Form.Select
                        value={selectedIndexLanguage}
                        onChange={(e) => setSelectedIndexLanguage(e.target.value)}
                        disabled={savingIndex || loadingIndexTokens}
                    >
                        <option value="en">English</option>
                        <option value="it">Italiano</option>
                        <option value="de">Deutsch</option>
                        <option value="fr">Français</option>
                    </Form.Select>
                </Form.Group>

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

                {htmlMode === 'wysiwyg' && !loadingIndexTokens && (
                    <div className="mb-2">
                        <ButtonGroup size="sm">
                            <Button
                                variant="outline-secondary"
                                size="sm"
                                onClick={handleInsertImageIndex}
                                disabled={savingIndex}
                            >
                                <i className="bi bi-image"></i> {t('files.insert_image')}
                            </Button>
                            <Button
                                variant="outline-secondary"
                                size="sm"
                                onClick={handleInsertFileIndex}
                                disabled={savingIndex}
                                className="me-2"
                            >
                                <i className="bi bi-file-earmark"></i> {t('files.insert_file')}
                            </Button>
                        </ButtonGroup>
                    </div>
                )}

                {loadingIndexTokens && (
                    <div className="text-center p-3">
                        <Spinner animation="border" />
                    </div>
                )}
                {!loadingIndexTokens && htmlMode === 'wysiwyg' ? (
                    <ReactQuill
                        ref={setQuillRefIndex}
                        theme="snow"
                        value={indexHtmlWithTokens}
                        onChange={handleIndexHtmlChange}
                        modules={{
                            toolbar: [
                                [{ 'header': [1, 2, 3, false] }],
                                ['bold', 'italic', 'underline', 'strike'],
                                [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                                [{ 'indent': '-1'}, { 'indent': '+1' }],
                                [{ 'color': [] }, { 'background': [] }],
                                ['link', 'blockquote', 'code-block'],
                                ['clean']
                            ]
                        }}
                    />
                ) : !loadingIndexTokens ? (
                    <Form.Control
                        as="textarea"
                        name="html"
                        value={indexHtmlWithTokens}
                        onChange={(e) => handleIndexHtmlChange(e.target.value)}
                        rows={15}
                        style={{ fontFamily: 'monospace', fontSize: '0.9em' }}
                    />
                ) : null}
                {!loadingIndexTokens && (
                    <div className="mt-2">
                        <Button
                            variant="primary"
                            onClick={handleSaveIndex}
                            disabled={savingIndex}
                        >
                            {savingIndex ? (
                                <>
                                    <Spinner animation="border" size="sm" className="me-2" />
                                    {t('common.saving')}
                                </>
                            ) : (
                                <>
                                    <i className="bi bi-save"></i> {t('folder.save_index')}
                                </>
                            )}
                        </Button>
                    </div>
                )}
            </div>

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
            
            <FileSelector
                show={showFileSelectorIndex}
                onHide={() => setShowFileSelectorIndex(false)}
                onSelect={handleFileSelectIndex}
                fileType={fileSelectorTypeIndex}
            />
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
