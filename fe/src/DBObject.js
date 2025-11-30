import React, { useState, useEffect } from 'react';
import { Alert, Button, Card, Form, Spinner } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
    CountryView,
    UserLinkView,
    ObjectLinkView,
    HtmlFieldView
} from './sitenavigation_utils';
import ObjectLinkSelector from './ObjectLinkSelector'
import PermissionsEditor from './PermissionsEditor';

import axiosInstance from './axios';

export function ObjectHeaderView({ data, metadata, objectData, dark }) {
    const { t } = useTranslation();

    return (
        <>
            <div className="row">
                {data.father_id && data.father_id!=="0" && <div className="col-md-2 col-4 text-end"><small style={{ opacity: 0.7 }}>{t('dbobjects.parent')}:</small></div>}
                {data.father_id && data.father_id!=="0"  && 
                    <div className="col-md-3 col-8">
                        <small style={{ opacity: 0.7 }}><ObjectLinkView obj_id={data.father_id} dark={dark} /></small>
                    </div>
                }
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0" && <div className="col-md-2 col-4 text-end"><small style={{ opacity: 0.7 }}>{t('dbobjects.linked_to')}:</small></div>}
                {data.fk_obj_id && data.fk_obj_id!==data.father_id && data.fk_obj_id!=="0"  && 
                    <div className="col-md-3 col-8">
                        <small style={{ opacity: 0.7 }}><ObjectLinkView obj_id={data.fk_obj_id} dark={dark} /></small>
                    </div>
                }
            </div>
            <div className="row">
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}><i className={`bi bi-${classname2bootstrapIcon(metadata.classname)}`} title={metadata.classname}></i> {t('dbobjects.' + metadata.classname)}</small>
                </div>
                <div className="col-md-3 col-8">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.id')}: {data.id}</small>
                </div>
                <div className="col-md-2 col-4 text-end"><small style={{ opacity: 0.7 }}>{t('dbobjects.permissions')}:</small></div>
                <div className="col-md-3 col-8">
                    <small style={{ opacity: 0.7 }}>{data.permissions}</small>
                </div>
            </div>
        </>
    );
}

export function ObjectFooterView({ data, metadata, objectData, dark }) {
    const { t } = useTranslation();

    return (
        <>
            <div className="row">
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.owner')}:</small>
                </div>
                <div className="col-md-3 col-8">
                    <small style={{ opacity: 0.7 }}>{objectData && objectData.owner_name}</small>
                </div>
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.group')}:</small>
                </div>
                <div className="col-md-3 col-8">
                    <small style={{ opacity: 0.7 }}>{objectData && objectData.group_name}</small>
                </div>
            </div>
            <div className="row">
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.created')}:</small>
                </div>
                <div className="col-md-6 col-8">
                    <small style={{ opacity: 0.7 }}>{formateDateTimeString(data.creation_date)} - {objectData && objectData.creator_name}</small>
                </div>
            </div>
            <div className="row">
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.modified')}:</small>
                </div>
                <div className="col-md-6 col-8">
                    <small style={{ opacity: 0.7 }}>{formateDateTimeString(data.last_modify_date)} -{objectData && objectData.last_modifier_name}</small>
                </div>
            </div>
            {data.deleted_date && 
            <div className="row">
                <div className="col-md-2 col-4 text-end">
                    <small style={{ opacity: 0.7 }}>{t('dbobjects.deleted')}:</small>
                </div>
                <div className="col-md-3 col-8">
                    <small style={{ opacity: 0.7 }}>{data && data.deleted_date ? formateDateTimeString(data.deleted_date) : '--'} - {objectData && objectData.deleted_by_name ? objectData.deleted_by_name : '--'}</small>
                </div>
            </div>
            }
        </>
    );
}

// Generic view for DBObject
export function ObjectView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    return (
        <Card className="mb-3" bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
            <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectHeaderView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Header>
            <Card.Body className={dark ? 'bg-secondary bg-opacity-10' : ''}>
                <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
                {!data.html && data.description && <hr />}
                {data.description && (
                    <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
                )}
                {data.html && <hr />}
                {data.html && (
                    <HtmlFieldView htmlContent={data.html} dark={dark} />
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Footer>
        </Card>
    );
}

// Generic edit form for other DBObjects
export function ObjectEdit({ data, metadata, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        permissions: data.permissions || 'rwxr-x---',
        father_id: data.father_id || null,
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
            <Alert variant="info" className="mb-3">
                <i className="bi bi-info-circle me-2"></i>
                Editing {metadata.classname} - Basic fields only
            </Alert>

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

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
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
