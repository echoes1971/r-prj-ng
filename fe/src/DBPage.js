import React, { use, useContext } from 'react';
import { useState, useEffect, useRef } from 'react';
import { ButtonGroup, Form, Spinner, Button, Overlay, Popover } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import ReactDOM from 'react-dom';
import ReactQuill, { Quill } from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { useTranslation } from 'react-i18next';
import EmojiPicker from 'emoji-picker-react';
import FileSelector from './FileSelector';
import ObjectLinkSelector from './ObjectLinkSelector';
import PermissionsEditor from './PermissionsEditor';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
} from './sitenavigation_utils';
import axiosInstance from './axios';
import { ObjectSearch } from './DBObject';
import { ThemeContext } from './ThemeContext';

// Configure Quill to preserve data-dbfile-id attribute

const Image = Quill.import('formats/image');
class CustomImage extends Image {
    static formats(domNode) {
        const formats = super.formats(domNode);
        formats['data-dbfile-id'] = domNode.getAttribute('data-dbfile-id');
        return formats;
    }
    
    format(name, value) {
        if (name === 'data-dbfile-id') {
            if (value) {
                this.domNode.setAttribute('data-dbfile-id', value);
            } else {
                this.domNode.removeAttribute('data-dbfile-id');
            }
        } else {
            super.format(name, value);
        }
    }
}
Quill.register(CustomImage, true);

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
export function extractFileIDs(html) {
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
export async function requestFileTokens(fileIDs) {
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
 * Inject tokens into HTML for viewing
 * Adds ?token=... to src/href attributes of elements with data-dbfile-id
 */
export function injectTokensForViewing(html, tokens) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    
    doc.querySelectorAll('[data-dbfile-id]').forEach(el => {
        const fileId = el.getAttribute('data-dbfile-id');
        const token = tokens[fileId];
        
        if (token) {
            if (el.tagName === 'IMG') {
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
export function cleanTokensBeforeSave(html) {
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

/**
 * Inject tokens into HTML for WYSIWYG editing
 * Adds ?token=... to src/href attributes of elements with data-dbfile-id
 */
export function injectTokensForEditing(html, tokens) {
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


export function HtmlView({ html, dark }) {
    const [htmlWithTokens, setHtmlWithTokens] = useState(html || '');
    const [loadingTokens, setLoadingTokens] = useState(false);

    // Load tokens for embedded files when component mounts or HTML changes
    useEffect(() => {
        const loadTokens = async () => {
            const fileIDs = extractFileIDs(htmlWithTokens);
            if (fileIDs.length === 0) {
                setHtmlWithTokens(html);
                return;
            }

            setLoadingTokens(true);
            try {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = injectTokensForViewing(html, tokens);
                setHtmlWithTokens(htmlWithTokens);
            } catch (error) {
                console.error('Failed to load tokens for embedded files:', error);
                setHtmlWithTokens(html);
            } finally {
                setLoadingTokens(false);
            }
        };

        loadTokens();
    }, [html]);

    return (
        <>
        {loadingTokens && (
            <div className="text-center py-3">
                <Spinner animation="border" size="sm" className="me-2" />
                <span>Loading...</span>
            </div>
        )}
        {!loadingTokens && htmlWithTokens && (
            <div dangerouslySetInnerHTML={{ __html: htmlWithTokens }}></div>
        )}
        </>
    );
}

// View for DBPage
export function PageView({ data, metadata, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [htmlWithTokens, setHtmlWithTokens] = useState(data.html || '');
    const [loadingTokens, setLoadingTokens] = useState(false);

    // Load tokens for embedded files when component mounts or HTML changes
    useEffect(() => {
        const loadTokens = async () => {
            const fileIDs = extractFileIDs(data.html);
            if (fileIDs.length === 0) {
                setHtmlWithTokens(data.html);
                return;
            }

            setLoadingTokens(true);
            try {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = injectTokensForViewing(data.html, tokens);
                setHtmlWithTokens(htmlWithTokens);
            } catch (error) {
                console.error('Failed to load tokens for embedded files:', error);
                setHtmlWithTokens(data.html);
            } finally {
                setLoadingTokens(false);
            }
        };

        loadTokens();
    }, [data.id, data.html]);

    return (
        <div>
            {data.name && (
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
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
    //             <div className="text-secondary mt-3">
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

// Edit form for DBPage
export function PageEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [htmlMode, setHtmlMode] = useState('source'); // 'wysiwyg' or 'source'
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        permissions: data.permissions || 'rwx------',
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
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const emojiButtonRef = useRef(null);

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

    const handleHtmlChange = async (value) => {
        setFormData(prev => ({...prev, html: value}));
        
        // Extract file IDs and reload tokens for immediate preview
        const fileIDs = extractFileIDs(value);
        if (fileIDs.length === 0) {
            setHtmlWithTokens(value);
            return;
        }
        
        try {
            const tokens = await requestFileTokens(fileIDs);
            const htmlWithTokens = injectTokensForEditing(value, tokens);
            setHtmlWithTokens(htmlWithTokens);
        } catch (error) {
            console.error('Failed to reload tokens after HTML change:', error);
            setHtmlWithTokens(value);
        }
    };

    // Handle file selection from modal
    const handleFileSelect = (file) => {
        if (!quillRef) return;

        const quill = quillRef.getEditor();
        const range = quill.getSelection(true);
        
        if (fileSelectorType === 'image') {
            // Insert image tag - CustomImage preserves data-dbfile-id
            const imageHtml = `<img src="/api/files/${file.id}/download" data-dbfile-id="${file.id}" alt="${file.name}" style="max-width: 100%;" />`;
            quill.clipboard.dangerouslyPasteHTML(range.index, imageHtml);
            // TODO is it not redundant?
            handleHtmlChange(quill.root.innerHTML);
        } else {
            // Insert link
            // const linkHtml = `<a href="/api/files/${file.id}/download" data-dbfile-id="${file.id}">${file.name}</a>`;
            const linkHtml = `<a href="/f/${file.id}/download" data-dbfile-id="${file.id}">${file.name}</a>`;
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

    const handleEmojiClick = (emojiObject) => {
        if (!quillRef) return;
        
        const quill = quillRef.getEditor();
        const range = quill.getSelection(true);
        
        // Insert emoji at cursor position
        quill.insertText(range.index, emojiObject.emoji);
        quill.setSelection(range.index + emojiObject.emoji.length);
        
        // Update state but keep picker open for multiple selections
        handleHtmlChange(quill.root.innerHTML);
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
                            <Button 
                                ref={emojiButtonRef}
                                variant="outline-secondary"
                                onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                                title={t('editor.insert_emoji') || 'Insert Emoji'}
                            >
                                <i className="bi bi-emoji-smile me-1"></i>
                                {t('editor.insert_emoji') || 'Insert Emoji'}
                            </Button>
                        </ButtonGroup>
                        <Overlay
                            show={showEmojiPicker}
                            target={emojiButtonRef.current}
                            // placement="bottom-start"
                            rootClose
                            onHide={() => setShowEmojiPicker(false)}
                        >
                            <Popover id="emoji-picker-popover">
                                <Popover.Body>
                                    <EmojiPicker
                                        onEmojiClick={handleEmojiClick}
                                        width={400}
                                        height={400}
                                        autoFocusSearch={false}
                                    />
                                </Popover.Body>
                            </Popover>
                        </Overlay>
                        {/* <Form.Text className="text-secondary ms-2">
                            Insert files/images you have permission to edit
                        </Form.Text> */}
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
                                [{ 'color': [] }, { 'background': [] }],
                                ['link', 'blockquote', 'code-block'],
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
                <Form.Text className="text-secondary">
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
  ];

  const resultsColumns = [
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
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
