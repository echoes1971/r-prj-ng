import React, { useState, useEffect } from 'react';
import { Form, Row, Col } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

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
function PermissionsEditor({ value = 'rwxr-x---', onChange, name = 'permissions', label, disabled = false, dark = false }) {
    const { t } = useTranslation();
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
        <Form.Group className="mb-3">
            {label && <Form.Label>{label}</Form.Label>}
            
            <div className={`border rounded p-3 ${dark ? 'border-secondary' : 'bg-light'}`}>
                <div className="mb-3 pb-3 border-bottom">
                    <small className="text-secondary">
                        <strong>{t('permissions.current') || 'Current permissions'}:</strong>{' '}
                        <code className={`px-2 py-1 rounded ${dark ? 'bg-secondary text-light' : 'bg-white'}`}>
                            {permissionsToString(permissions)}
                        </code>
                    </small>
                </div>

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
            </div>
            
            <Form.Text className="text-secondary">
                {t('permissions.hint') || 'Set read, write, and execute permissions for owner, group, and others'}
            </Form.Text>
        </Form.Group>
    );
}

export default PermissionsEditor;
