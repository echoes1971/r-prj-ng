import React, { useContext, useEffect, useState } from "react";
import { Alert, Container, Form, Button, Spinner } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import axiosInstance from './axios';
import ObjectLinkSelector from './ObjectLinkSelector';
import ObjectList from "./ObjectList";
import { ObjectSearch } from "./DBObject";
import { cleanTokensBeforeSave, HtmlEdit, HtmlView } from "./ContentHtml";
import PermissionsEditor from './PermissionsEditor';



// View for DBFolder
export function FolderView({ data, metadata, dark, onFilesUploaded }) {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language; // 'it', 'en', 'de', 'fr'

    const isDeleted = data && data.deleted_date;


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
                formData.append('permissions', 'rwxr-----'); // Default permissions

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
                opacity: isDeleted ? 0.5 : 1,
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
                    <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}</h2>
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

// Edit form for DBFolder
export function FolderEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        father_id: data.father_id || '0',
        name: data.name || '',
        description: data.description || '',
        fk_obj_id: data.fk_obj_id || '0',
        permissions: data.permissions || 'rwxr-x---',
        childs_sort_order: data.childs_sort_order || '',
    });
    const [children, setChildren] = useState([]);
    const [loadingChildren, setLoadingChildren] = useState(false);
    const [sortedChildrenIds, setSortedChildrenIds] = useState([]);
    const [draggedIndex, setDraggedIndex] = useState(null);
    
    // Index page editor states
    const [indexPages, setIndexPages] = useState([]);
    const [selectedIndexLanguage, setSelectedIndexLanguage] = useState('en');
    const [indexHtml, setIndexHtml] = useState('');
    const [savingIndex, setSavingIndex] = useState(false);

    const isDeleted = data && data.deleted_date;

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
            
            // Initialize sorted order from childs_sort_order or use current order
            if (formData.childs_sort_order) {
                const orderIds = formData.childs_sort_order.split(',').filter(id => id);
                setSortedChildrenIds(orderIds);
            // } else {
            //     setSortedChildrenIds(childrenData.map(child => child.data.id));
            }
            if(sortedChildrenIds.length !== 0) {
                console.log('Initial sortedChildrenIds:', sortedChildrenIds);
            }
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
            } else {
                setIndexHtml('');
            }
        } catch (error) {
            console.error('Failed to load index pages:', error);
            setIndexPages([]);
        }
    };

    // Reload index HTML when language changes
    useEffect(() => {
        const currentPage = indexPages.find(p => p.data.language && p.data.language.indexOf(selectedIndexLanguage) >= 0);
        if (currentPage) {
            setIndexHtml(currentPage.data.html || '');
        } else {
            setIndexHtml('');
        }
    }, [selectedIndexLanguage, indexPages]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
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
                <ObjectLinkSelector
                    value={formData.father_id || '0'}
                    onChange={handleChange}
                    classname="DBObject"
                    fieldName="father_id"
                    name="father_id"
                    label={t('dbobjects.parent')}
                />
            </Form.Group>

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

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

            {/* Index Page Editor */}
            <div className="mb-4 p-3 border rounded">
                <h5>{t('folder.index_page_editor')}</h5>
                <p className="text-secondary small">{t('folder.index_page_hint')}</p>
                
                <Form.Group className="mb-3">
                    <Form.Label>{t('common.language')}</Form.Label>
                    <Form.Select
                        value={selectedIndexLanguage}
                        onChange={(e) => setSelectedIndexLanguage(e.target.value)}
                        disabled={savingIndex}
                    >
                        <option value="en">English</option>
                        <option value="it">Italiano</option>
                        <option value="de">Deutsch</option>
                        <option value="fr">Français</option>
                    </Form.Select>
                </Form.Group>

                <HtmlEdit
                    htmlContent={indexHtml}
                    onHtmlContentChange={setIndexHtml}
                    dark={dark}
                    disabled={savingIndex}
                />

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
            </div>

            <Form.Group className="mb-3">
                <Form.Label>{t('files.linked_object')}</Form.Label>
                <ObjectLinkSelector
                    value={formData.fk_obj_id || '0'}
                    onChange={handleChange}
                    name="fk_obj_id"
                    fieldName="fk_obj_id"
                    classname="DBObject"
                />
            </Form.Group>

            {/* Children Sort Order */}
            {children.length > 0 && (
                <Form.Group className="mb-3">
                    <Form.Label>
                        {t('folder.children_order')}
                        <small className="ms-2 text-secondary">
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
                                    <div className="text-secondary text-center p-2">
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
                                    <Form.Label className="mt-3 mb-2">
                                        {t('folder.available_children')}
                                    </Form.Label>
                                    
                                    <ObjectList
                                        items={children
                                            .filter(child => !sortedChildrenIds.includes(child.data.id))
                                            .map(child => ({
                                                id: child.data.id,
                                                name: child.data.name,
                                                description: child.data.description,
                                                classname: child.metadata?.classname
                                            }))
                                        }
                                        onItemClick={(item) => toggleChildInOrder(item.id)}
                                        showViewToggle={true}
                                        storageKey="folderChildrenViewMode"
                                        defaultView="list"
                                    />
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
                    disabled={saving || isDeleted}
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

export function Folders() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);

  const searchClassname = "DBFolder";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
  ];

  const resultsColumns = [
    { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
