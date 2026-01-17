import React, { use, useContext } from 'react';
import { useState, useEffect, useRef } from 'react';
import { ButtonGroup, Form, Spinner, Button, Overlay, Popover } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import ReactDOM from 'react-dom';
import ReactQuill, { Quill } from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { useTranslation } from 'react-i18next';
import ObjectLinkSelector from '../ObjectLinkSelector';
import PermissionsEditor from '../PermissionsEditor';
import { cleanTokensBeforeSave, HtmlEdit, HtmlView } from '../ContentHtml';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
} from '../sitenavigation_utils';
import axiosInstance from '../axios';
import { ObjectSearch } from '../DBObject';
import { ThemeContext } from '../ThemeContext';



// View for DBPage
export function PageView({ data, metadata, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [htmlWithTokens, setHtmlWithTokens] = useState(data.html || '');
    const [loadingTokens, setLoadingTokens] = useState(false);

    const isDeleted = data && data.deleted_date;

    // Load tokens for embedded files when component mounts or HTML changes
    // useEffect(() => {
    //     const loadTokens = async () => {
    //         const fileIDs = extractFileIDs(data.html);
    //         if (fileIDs.length === 0) {
    //             setHtmlWithTokens(data.html);
    //             return;
    //         }
    //         setLoadingTokens(true);
    //         try {
    //             const tokens = await requestFileTokens(fileIDs);
    //             const htmlWithTokens = injectTokensForViewing(data.html, tokens);
    //             setHtmlWithTokens(htmlWithTokens);
    //         } catch (error) {
    //             console.error('Failed to load tokens for embedded files:', error);
    //             setHtmlWithTokens(data.html);
    //         } finally {
    //             setLoadingTokens(false);
    //         }
    //     };
    //     loadTokens();
    // }, [data.id, data.html]);

    return (
        <div style={isDeleted ? { opacity: 0.5 } : {}}>
            {data.name && (
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}</h2>
            )}
            {data.description && (
                <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
            )}
            {loadingTokens && (
                <div className="text-center py-3">
                    <Spinner animation="border" size="sm" className="me-2" />
                    <span>Loading...</span>
                </div>
            )}
            {!loadingTokens && htmlWithTokens && (
                <HtmlView html={htmlWithTokens} dark={dark} />
                // <div dangerouslySetInnerHTML={{ __html: htmlWithTokens }}></div>
            )}
        </div>
    );
}

// Edit form for DBPage
export function PageEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        permissions: data.permissions || 'rwx------',
        html: data.html || '',
        language: data.language || 'en',
        father_id: data.father_id || '0',
        owner: data.owner || null,
        group_id: data.group_id || null,
        fk_obj_id: data.fk_obj_id || '0',
    });

    const isDeleted = data && data.deleted_date;

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
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

            <div className="row">
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.father_id || '0'}
                    onChange={handleChange}
                    classname="DBObject"
                    fieldName="father_id"
                    label={t('dbobjects.parent')}
                />
                </div>
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.owner}
                    onChange={handleChange}
                    classname="DBUser"
                    fieldName="owner"
                    label={t('permissions.owner')}
                    required={false}
                />
                </div>
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.group_id}
                    onChange={handleChange}
                    classname="DBGroup"
                    fieldName="group_id"
                    label={t('permissions.group')}
                    required={false}
                />
                </div>
            </div>

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

            <Form.Group className="mb-3">
                <Form.Label>{t('common.language')}</Form.Label>
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
                <HtmlEdit htmlContent={formData.html} dark={dark}
                    onHtmlContentChange={(newHtml) => setFormData(prev => ({...prev, html: newHtml}))} />
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
            
            {/* File Selector Modal
            <FileSelector
                show={showFileSelector}
                onHide={() => setShowFileSelector(false)}
                onSelect={handleFileSelect}
                fileType={fileSelectorType}
            /> */}
        </Form>
    );
}

export function Pages() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);
  // const [query, setQuery] = useState("");
  // const [editingFolder, setEditingFolder] = useState(null); // folder in editing
  const searchClassname = "DBPage";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    // { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
    { name: t("dbobjects.language") || "Language", attribute: "language", type: "languageSelector" },
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector" },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector" },
  ];

  const resultsColumns = [
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector", hideOnSmall: true },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
    { name: t("common.language") || "Language", attribute: "language", type: "languageView", hideOnSmall: true },
  ];
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
  );
}
