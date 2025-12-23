import React, { useState, useEffect, useContext } from 'react';
import { Card, Container, Form, Button, Spinner, Alert, ButtonGroup } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import ReactDOM from 'react-dom';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import axiosInstance from './axios';
import ObjectLinkSelector from './ObjectLinkSelector';
import PermissionsEditor from './PermissionsEditor';
import { formatDescription, formatObjectId } from './sitenavigation_utils';
import { ThemeContext } from './ThemeContext';
import { ObjectSearch } from './DBObject';


export function FileDownload() {
    const { objectId } = useParams();
    const navigate = useNavigate();

    const currentObjectId = objectId || null;

    const [content, setContent] = useState(null);

    useEffect(() => {
        const loadContent = async () => {
            try {
                const response = await axiosInstance.get(`/content/${formatObjectId(currentObjectId)}`);
                setContent(response.data);
                // alert(JSON.stringify(response.data));

                // Trigger file download
                const downloadUrl = `/api/files/${response.data.data.id}/download?token=${response.data.metadata.download_token}`;
                window.location.href = downloadUrl;

                // Optionally, navigate back or to another page after download
                navigate(-1); // Go back to previous page

                // Close the window after download starts with a slight delay
                setTimeout(() => {
                    window.close();
                }, 500);

            } catch (error) {
                console.error('Error loading content for file download:', error);
                setContent(null);
            }
        }
        loadContent();
    }, [currentObjectId]);
    return (<Container className="mt-4">
        {content ? (
            <Alert variant="success">
                {`Initiated download for file: ${content.name}`}
            </Alert>
        ) : (
            <Spinner animation="border" role="status">
                <span className="visually-hidden">Loading...</span>
            </Spinner>
        )}
    </Container>
    );
}

export function FileView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [preview, setPreview] = useState(null);

    const isDeleted = data && data.deleted_date;

    
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
        <div style={isDeleted ? { opacity: 0.5 } : {}}>
            {data.name && (
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}</h2>
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
                        <small className={dark ? "text-white-50" : "text-secondary"}>{data.name}</small>
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

// Edit form for DBFile
export function FileEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
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

    const isDeleted = data && data.deleted_date;

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
                name: prev.name==='' || prev.name==='New File' ? file.name : prev.name,
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
                                <small className="text-secondary">{formData.filename}</small>
                            </div>
                        </div>
                    ) : (
                        <>
                            <i className="bi bi-cloud-upload fs-1"></i>
                            <p className="mb-0">
                                {selectedFile ? selectedFile.name : (t('files.drop_or_click') || 'Drop file here or click to browse')}
                            </p>
                            {formData.filename && !selectedFile && (
                                <small className="text-secondary d-block mt-2">
                                    {t('files.current') || 'Current'}: {formData.filename}
                                </small>
                            )}
                        </>
                    )}
                </div>
                <Form.Text className="text-secondary">
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
                <Form.Text className="text-secondary">
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

export function Files() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);

  const searchClassname = "DBFile";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
  ];

  const resultsColumns = [
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("files.preview") || "File", attribute: "id", type: "imageView", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
  ];
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
  );
}
