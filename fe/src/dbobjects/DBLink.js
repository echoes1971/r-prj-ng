import React, { useContext, useState } from 'react';
import { Form, Button, Spinner } from 'react-bootstrap';
import { ThemeContext } from "../ThemeContext";
import { useTranslation } from 'react-i18next';
import ObjectLinkSelector from '../ObjectLinkSelector';
import PermissionsEditor from '../PermissionsEditor';
import { formatDescription } from '../sitenavigation_utils';
import { ObjectSearch } from "./DBObject";


/**
 * LinkView - Display component for DBLink objects
 * Shows a link with name, description, and external link icon
 */
export function LinkView({ data, metadata, objectData, dark }) {
    const { t } = useTranslation();

    const isDeleted = data && data.deleted_date;

    // Default target if not specified
    const target = data.target || '_blank';
    const isExternal = target === '_blank';

    return (
        <div style={isDeleted ? { opacity: 0.5 } : {}}>
            <h2 className={dark ? 'text-light' : 'text-dark'}>
                <i className="bi bi-link-45deg me-2"></i>
                {data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}
            </h2>
            
            {data.description && (
                <>
                    <hr />
                    <div dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></div>
                </>
            )}

            <hr />

            <div className="d-grid gap-2">
                <a 
                    href={data.href} 
                    target={target}
                    rel={isExternal ? "noopener noreferrer" : undefined}
                    className="btn btn-primary btn-lg"
                >
                    <i className={`bi bi-${isExternal ? 'box-arrow-up-right' : 'arrow-right'} me-2`}></i>
                    {t('link.open') || 'Open Link'}
                    {isExternal && (
                        <small className="ms-2 opacity-75">
                            ({t('link.new_tab') || 'opens in new tab'})
                        </small>
                    )}
                </a>
            </div>

            {data.href && (
                <div className="mt-3">
                    <small className="text-muted">
                        <i className="bi bi-link me-1"></i>
                        {data.href}
                    </small>
                </div>
            )}
        </div>
    );
}


/**
 * LinkEdit - Edit form component for DBLink objects
 * Form fields: name, description, href (URL), target, fk_obj_id, permissions, father_id
 */
export function LinkEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        href: data.href || '',
        target: data.target || '_blank',
        fk_obj_id: data.fk_obj_id || '0',
        permissions: data.permissions || 'rwxr-x---',
        father_id: data.father_id || '0',
        owner: data.owner || null,
        group_id: data.group_id || null,
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
        onSave(formData);
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
                <Form.Label>{t('common.name') || 'Name'}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    placeholder={t('link.name_placeholder') || 'Link title'}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description') || 'Description'}</Form.Label>
                <Form.Control
                    as="textarea"
                    rows={3}
                    name="description"
                    value={formData.description}
                    onChange={handleChange}
                    placeholder={t('link.description_placeholder') || 'Optional description'}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>
                    <i className="bi bi-link-45deg me-1"></i>
                    {t('link.url') || 'URL'}
                </Form.Label>
                <Form.Control
                    // Can be https:// or relative URL
                    type="text"

                    // type="url"
                    name="href"
                    value={formData.href}
                    onChange={handleChange}
                    required
                    placeholder="https://example.com"
                />
                <Form.Text className="text-muted">
                    {t('link.url_help') || 'Full URL including http:// or https://'}
                </Form.Text>
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('link.target') || 'Open In'}</Form.Label>
                <Form.Select
                    name="target"
                    value={formData.target}
                    onChange={handleChange}
                >
                    <option value="_blank">
                        {t('link.target_blank') || 'New tab (_blank)'}
                    </option>
                    <option value="_self">
                        {t('link.target_self') || 'Same tab (_self)'}
                    </option>
                    <option value="_parent">
                        {t('link.target_parent') || 'Parent frame (_parent)'}
                    </option>
                    <option value="_top">
                        {t('link.target_top') || 'Top frame (_top)'}
                    </option>
                </Form.Select>
            </Form.Group>

            <ObjectLinkSelector
                value={formData.fk_obj_id || '0'}
                onChange={handleChange}
                classname="DBObject"
                fieldName="fk_obj_id"
                label={t('dbobjects.linked_to') || 'Linked To'}
            />

            {error && (
                <div className="alert alert-danger">
                    {error}
                </div>
            )}

            <div className="d-flex gap-2">
                <Button 
                    variant="primary" 
                    type="submit" 
                    disabled={saving || isDeleted}
                >
                    {saving ? (
                        <><Spinner animation="border" size="sm" className="me-2" />{t('common.saving') || 'Saving...'}</>
                    ) : (
                        <>
                        <i className="bi bi-check-lg me-1"></i>
                        {t('common.save') || 'Save'}
                        </>
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
                    <i className="bi bi-x-lg me-1"></i>
                    {t('common.cancel') || 'Cancel'}
                </Button>
                {onDelete && (
                    <Button 
                        variant="outline-danger" 
                        onClick={onDelete}
                        disabled={saving}
                        className="ms-auto"
                    >
                        <i className="bi bi-trash me-1"></i>
                        {t('common.delete') || 'Delete'}
                    </Button>
                )}
            </div>
        </Form>
    );
}

export function Links() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);

  const searchClassname = "DBLink";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
  ];

  const resultsColumns = [
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    // { name: t("files.preview") || "File", attribute: "id", type: "imageView", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
