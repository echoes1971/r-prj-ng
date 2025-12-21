import React, { use, useContext } from 'react';
import { useState, useEffect, useRef } from 'react';
import { ButtonGroup, Form, Spinner, Button, Overlay, Popover, Modal } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import ReactDOM from 'react-dom';
import ReactQuill, { Quill } from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { useTranslation } from 'react-i18next';
import EmojiPicker from 'emoji-picker-react';
import FileSelector from './FileSelector';
import ObjectLinkSelector from './ObjectLinkSelector';
import axiosInstance from './axios';


// **** RRA - Start: is this still needed?

// Configure Quill to preserve data-dbfile-id attribute and support width/height
const Image = Quill.import('formats/image');
class CustomImage extends Image {
    static formats(domNode) {
        const formats = super.formats(domNode);
        formats['data-dbfile-id'] = domNode.getAttribute('data-dbfile-id');
        
        // Preserve width and height from style attribute
        const style = domNode.getAttribute('style') || '';
        const widthMatch = style.match(/width:\s*([^;]+)/);
        const heightMatch = style.match(/height:\s*([^;]+)/);
        if (widthMatch) formats['width'] = widthMatch[1].trim();
        if (heightMatch) formats['height'] = heightMatch[1].trim();
        
        return formats;
    }
    
    format(name, value) {
        if (name === 'data-dbfile-id') {
            if (value) {
                this.domNode.setAttribute('data-dbfile-id', value);
            } else {
                this.domNode.removeAttribute('data-dbfile-id');
            }
        } else if (name === 'width' || name === 'height') {
            if (value) {
                this.domNode.style[name] = value;
            } else {
                this.domNode.style[name] = '';
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

// **** RRA - End





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
 * Converts Quill classes to inline styles for proper display outside editor
 */
export function injectTokensForViewing(html, tokens) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    
    // Inject file tokens
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

export function convertQuillClassesToStyles(html) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');

    // Convert Quill classes to inline styles for display
    doc.querySelectorAll('[class*="ql-"]').forEach(el => {
        const classList = el.classList;
        const styles = [];
        const classesToRemove = [];
        
        // Alignment: ql-align-center, ql-align-right, ql-align-justify
        if (classList.contains('ql-align-center')) {
            styles.push('text-align: center');
            classesToRemove.push('ql-align-center');
        } else if (classList.contains('ql-align-right')) {
            styles.push('text-align: right');
            classesToRemove.push('ql-align-right');
        } else if (classList.contains('ql-align-justify')) {
            styles.push('text-align: justify');
            classesToRemove.push('ql-align-justify');
        }
        
        // Indentation: ql-indent-1 to ql-indent-8
        for (let i = 1; i <= 8; i++) {
            if (classList.contains(`ql-indent-${i}`)) {
                styles.push(`margin-left: ${i * 3}em`);
                classesToRemove.push(`ql-indent-${i}`);
                break;
            }
        }
        
        // Direction: ql-direction-rtl
        if (classList.contains('ql-direction-rtl')) {
            styles.push('direction: rtl');
            styles.push('text-align: right');
            classesToRemove.push('ql-direction-rtl');
        }
        
        // Apply styles if any
        if (styles.length > 0) {
            const existingStyle = el.getAttribute('style') || '';
            const newStyle = existingStyle + (existingStyle ? '; ' : '') + styles.join('; ');
            el.setAttribute('style', newStyle);
            
            // Remove Quill classes
            classesToRemove.forEach(cls => el.classList.remove(cls));
        }
    });
    return doc.body.innerHTML;
}

/**
 * Clean tokens from HTML before saving
 * Removes ?token=... from src/href but keeps data-dbfile-id
 * IMPORTANT: Keeps Quill classes (ql-align-*, ql-indent-*, etc.) so they work when re-opened in editor
 */
export function cleanTokensBeforeSave(html) {
    if (!html) return html;
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');
    
    // Clean file tokens
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
                const htmlWithTokens = convertQuillClassesToStyles(html);
                setHtmlWithTokens(htmlWithTokens);
                return;
            }

            setLoadingTokens(true);
            try {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = convertQuillClassesToStyles(injectTokensForViewing(html, tokens));
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


/**
 * HTML Editor Component with WYSIWYG and Source modes
 * 
 * Reference: https://quilljs.com/docs/configuration
 */
export function HtmlEdit({objID, htmlContent, onHtmlContentChange, dark}) {
    const { t } = useTranslation();

    const [htmlMode, setHtmlMode] = useState('wysiwyg'); // 'wysiwyg' or 'source'
    const [showFileSelector, setShowFileSelector] = useState(false);
    const [fileSelectorType, setFileSelectorType] = useState('file'); // 'file' or 'image'
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const [htmlWithTokens, setHtmlWithTokens] = useState(htmlContent || '');
    const [loadingTokens, setLoadingTokens] = useState(false);
    // const quillRef = useRef(null);
    const [quillRef, setQuillRef] = useState(null);
    const emojiButtonRef = useRef(null);

    const [currentFileIDs, setCurrentFileIDs] = useState([]);

    // Resize modal states
    const [showResizeModal, setShowResizeModal] = useState(false);
    const [selectedImage, setSelectedImage] = useState(null);
    const [resizeWidth, setResizeWidth] = useState('');
    const [resizeHeight, setResizeHeight] = useState('');

    // Load tokens for embedded files when component mounts or HTML changes
    useEffect(() => {
        const loadTokens = async () => {
            const fileIDs = extractFileIDs(htmlContent);

            // Compare with currentFileIDs to avoid unnecessary reloads
            if (JSON.stringify(fileIDs) === JSON.stringify(currentFileIDs)) {
                // No change in file IDs, skip reload
                setHtmlWithTokens(htmlContent);
                return;
            }
            setCurrentFileIDs(fileIDs);


            if (fileIDs.length === 0) {
                // No embedded files, use original HTML
                setHtmlWithTokens(htmlContent);
                return;
            }

            setLoadingTokens(true);
            try {
                const tokens = await requestFileTokens(fileIDs);
                const htmlWithTokens = injectTokensForEditing(htmlContent, tokens);
                setHtmlWithTokens(htmlWithTokens);
            } catch (error) {
                console.error('Failed to load tokens for embedded files:', error);
                setHtmlWithTokens(htmlContent);
            } finally {
                setLoadingTokens(false);
            }
        };

        loadTokens();
    }, [htmlContent]); // Only reload when page HTML changes

    const handleHtmlChange = async (value) => {
        onHtmlContentChange(value);
        
        // Extract file IDs and reload tokens for immediate preview
        const fileIDs = extractFileIDs(value);

        // Compare with currentFileIDs to avoid unnecessary reloads
        if (JSON.stringify(fileIDs) === JSON.stringify(currentFileIDs)) {
            console.log('File IDs unchanged (' + fileIDs.join(', ') + '), skipping token reload');
            setHtmlWithTokens(value);
            return;
        }
        setCurrentFileIDs(fileIDs);

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

    // Handle image resize
    const handleResizeImage = () => {
        if (!quillRef) return;

        const quill = quillRef.getEditor();
        const range = quill.getSelection();
        if (!range) {
            alert('Please select an image first.');
            return;
        }

        // Check if the current position contains an image
        const contents = quill.getContents(range.index, 1);
        const isImage = contents.ops && contents.ops[0] && contents.ops[0].insert && typeof contents.ops[0].insert === 'object' && contents.ops[0].insert.image;

        if (!isImage) {
            alert('Please select an image first.');
            return;
        }

        // Find the image DOM node
        const [leaf] = quill.getLeaf(range.index);
        if (leaf && leaf.domNode && leaf.domNode.tagName === 'IMG') {
            setSelectedImage(leaf.domNode);
            // Pre-fill with current dimensions
            setResizeWidth(leaf.domNode.style.width || '');
            setResizeHeight(leaf.domNode.style.height || '');
            setShowResizeModal(true);
        } else {
            alert('Unable to find the selected image.');
        }
    };

    // Apply resize to selected image
    const applyResize = () => {
        if (!selectedImage) return;

        if (resizeWidth) selectedImage.style.width = resizeWidth;
        if (resizeHeight) selectedImage.style.height = resizeHeight;

        // Update the editor content
        const quill = quillRef.getEditor();
        handleHtmlChange(quill.root.innerHTML);

        setShowResizeModal(false);
        setSelectedImage(null);
        setResizeWidth('');
        setResizeHeight('');
    };


    return (
        <>
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
                            variant="outline-secondary"
                            onClick={handleResizeImage}
                            title="Resize Image"
                        >
                            <i className="bi bi-arrows-angle-expand me-1"></i>
                            Resize Image
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
                    theme="snow"    // "snow" or "bubble"
                    style={{ height: '40rem', marginBottom: '3rem' }}
                    modules={{
                        toolbar: [
                            [{ 'header': [1, 2, 3, false] }],
                            ['bold', 'italic', 'underline', 'strike'],
                            [{ 'align': [] }],
                            [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                            [{ 'indent': '-1'}, { 'indent': '+1' }],
                            [{ 'color': [] }, { 'background': [] }],
                            ['blockquote', 'code-block'],
                            ['clean']
                        ]
                    }}
                />
            ) : !loadingTokens ? (
                <Form.Control
                    as="textarea"
                    name="html"
                    value={htmlWithTokens}
                    onChange={(e) => handleHtmlChange(e.target.value)}
                    rows={15}
                    style={{ fontFamily: 'monospace', fontSize: '0.9em' }}
                />
            ) : null}
            <Form.Text className="text-secondary">
                HTML content for the page. Use data-dbfile-id attribute to embed files (e.g., &lt;img src="/api/files/ID/download" data-dbfile-id="ID" /&gt;)
            </Form.Text>

            {/* File Selector Modal */}
            <FileSelector
                show={showFileSelector}
                onHide={() => setShowFileSelector(false)}
                onSelect={handleFileSelect}
                fileType={fileSelectorType}
            />

            {/* Image Resize Modal */}
            <Modal show={showResizeModal} onHide={() => setShowResizeModal(false)}>
                <Modal.Header closeButton>
                    <Modal.Title>Resize Image</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form.Group className="mb-3">
                        <Form.Label>Width (e.g., 200px or 50%)</Form.Label>
                        <Form.Control
                            type="text"
                            value={resizeWidth}
                            onChange={(e) => setResizeWidth(e.target.value)}
                            placeholder="e.g., 200px"
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Height (e.g., 150px or auto)</Form.Label>
                        <Form.Control
                            type="text"
                            value={resizeHeight}
                            onChange={(e) => setResizeHeight(e.target.value)}
                            placeholder="e.g., 150px"
                        />
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={() => setShowResizeModal(false)}>
                        Cancel
                    </Button>
                    <Button variant="primary" onClick={applyResize}>
                        Apply
                    </Button>
                </Modal.Footer>
            </Modal>

        </>
    );
}
