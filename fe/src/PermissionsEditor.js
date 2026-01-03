import React, { useContext, useState, useEffect } from 'react';
import { Form, Row, Col, Accordion } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { ThemeContext } from "./ThemeContext";

/**
 * Permissions editor component for Unix-style permissions (rwxr-x---)
 * 
 * Format: owner(rwx) group(rwx) others(rwx)
 * - r: read permission
 * - w: write permission  
 * - x: execute permission
 * - -: no permission
 * 
 * @param {string} value - Current permissions string (e.g., "rwxr-x---")
 * @param {function} onChange - Callback when permissions change
 * @param {string} name - Name of the field (default: "permissions")
 * @param {string} label - Label to display above the editor
 * @param {boolean} disabled - Whether the editor is disabled
 * @param {boolean} dark - Whether to use dark theme
 */
function PermissionsEditor({ value = 'rwxr-x---', onChange, name = 'permissions', label, disabled = false }) {
    const { t } = useTranslation();
    const { dark, themeClass } = useContext(ThemeContext);
    const [permissions, setPermissions] = useState({
        owner: { r: true, w: true, x: true },
        group: { r: true, w: false, x: true },
        others: { r: false, w: false, x: false }
    });

    // Parse permissions string into object
    useEffect(() => {
        if (value && value.length === 9) {
            setPermissions({
                owner: {
                    r: value[0] === 'r',
                    w: value[1] === 'w',
                    x: value[2] === 'x'
                },
                group: {
                    r: value[3] === 'r',
                    w: value[4] === 'w',
                    x: value[5] === 'x'
                },
                others: {
                    r: value[6] === 'r',
                    w: value[7] === 'w',
                    x: value[8] === 'x'
                }
            });
        }
    }, [value]);

    // Convert permissions object to string
    const permissionsToString = (perms) => {
        const owner = (perms.owner.r ? 'r' : '-') + (perms.owner.w ? 'w' : '-') + (perms.owner.x ? 'x' : '-');
        const group = (perms.group.r ? 'r' : '-') + (perms.group.w ? 'w' : '-') + (perms.group.x ? 'x' : '-');
        const others = (perms.others.r ? 'r' : '-') + (perms.others.w ? 'w' : '-') + (perms.others.x ? 'x' : '-');
        return owner + group + others;
    };

    // Handle checkbox change
    const handleChange = (category, permission) => {
        const newPermissions = {
            ...permissions,
            [category]: {
                ...permissions[category],
                [permission]: !permissions[category][permission]
            }
        };
        setPermissions(newPermissions);
        
        // Notify parent component
        onChange({
            target: {
                name: name,
                value: permissionsToString(newPermissions)
            }
        });
    };

    // Render permission checkboxes for a category
    const renderPermissionCheckboxes = (category, categoryLabel) => (
        <div className="mb-3">
            <strong>{categoryLabel}</strong>
            <div className="d-flex gap-3 mt-2">
                <Form.Check
                    type="checkbox"
                    id={`${name}-${category}-read`}
                    label={t('permissions.read') || 'Read'}
                    checked={permissions[category].r}
                    onChange={() => handleChange(category, 'r')}
                    disabled={disabled}
                />
                <Form.Check
                    type="checkbox"
                    id={`${name}-${category}-write`}
                    label={t('permissions.write') || 'Write'}
                    checked={permissions[category].w}
                    onChange={() => handleChange(category, 'w')}
                    disabled={disabled}
                />
                <Form.Check
                    type="checkbox"
                    id={`${name}-${category}-execute`}
                    label={t('permissions.execute') || 'Execute'}
                    checked={permissions[category].x}
                    onChange={() => handleChange(category, 'x')}
                    disabled={disabled}
                />
            </div>
        </div>
    );

    return (
        // <Form.Group className={`mb-3 ${themeClass}`}>
            // {label && <Form.Label>{label}</Form.Label>}
            <Accordion className="mb-3 rhobee-theme">
                <Accordion.Item eventKey="0" className='rhobee-theme'>
                    <Accordion.Header className='rhobee-theme'>
                        <i className="bi bi-shield-lock me-2"></i>
                        { label ? label :
                            t('permissions.current') || 'Current permissions'
                        }: 
                            <code className={`ms-2 me-2 ${dark ? 'text-light' : 'text-dark'}`} style={{opacity: 0.7}}>{permissionsToString(permissions)}</code>
                    </Accordion.Header>
                    <Accordion.Body className='rhobee-theme'>
                        <Row>
                            <Col md={4}>
                                {renderPermissionCheckboxes('owner', t('permissions.owner') || 'Owner')}
                            </Col>
                            <Col md={4}>
                                {renderPermissionCheckboxes('group', t('permissions.group') || 'Group')}
                            </Col>
                            <Col md={4}>
                                {renderPermissionCheckboxes('others', t('permissions.others') || 'Others')}
                            </Col>
                        </Row>
                        
                        <div className="mt-3 text-secondary small">
                            <i className="bi bi-info-circle me-1"></i>
                            {t('permissions.hint') || 'Set read, write, and execute permissions for owner, group, and others'}
                        </div>
                    </Accordion.Body>
                </Accordion.Item>
            </Accordion>
        // </Form.Group>
    );
}

export default PermissionsEditor;
