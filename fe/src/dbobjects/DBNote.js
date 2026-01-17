import React, { useState, useEffect, useContext } from 'react';
import { Card, Form, Button, Spinner, Alert } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import axiosInstance from '../axios';
import ObjectLinkSelector from '../ObjectLinkSelector';
import PermissionsEditor from '../PermissionsEditor';
import { formatDescription } from '../sitenavigation_utils';
import { ObjectSearch } from './DBObject';
import { ThemeContext } from '../ThemeContext';

// View for DBNote
export function NoteView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();

    const isDeleted = data && data.deleted_date;

    return (
        <Card style={isDeleted ? { opacity: 0.5 } : {}} className="mb-3 border-warning" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <br />
                {/* <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} /> */}
            </Card.Header>
            <Card.Body>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                <hr />
                {data.description && (
                    // <Card.Text>{data.description}</Card.Text>
                    <div className="content">
                        <p style={{ opacity: 0.7 }} dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></p>
                    </div>
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-warning bg-opacity-25' : 'bg-warning bg-opacity-10'}>
                <br />
                {/* <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} /> */}
            </Card.Footer>
        </Card>
    );
}

// Edit form for DBNote
export function NoteEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        permissions: data.permissions || 'rwx------',
        fk_obj_id: data.fk_obj_id || '0',
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

export function Notes() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);
  // const [query, setQuery] = useState("");
  // const [editingFolder, setEditingFolder] = useState(null); // folder in editing

  const searchClassname = "DBNote";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector" },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector" },
  ];

  const resultsColumns = [
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector", hideOnSmall: true },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
