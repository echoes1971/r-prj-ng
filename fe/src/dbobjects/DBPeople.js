import React, { useContext, useEffect, useState } from "react";
import { Accordion, Form, Button, Spinner, Alert } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from "react-i18next";
import { ThemeContext } from "../ThemeContext";
import { CountrySelector, CountryView, UserLinkView, ObjectLinkView } from '../ContentWidgets';
import ObjectLinkSelector from '../ObjectLinkSelector';
import { ObjectSearch } from "../DBObject";
import PermissionsEditor from '../PermissionsEditor';


export function PersonView({ data, metadata, objectData, dark }) {
    const navigate = useNavigate();
    const { t } = useTranslation();
    
    const isDeleted = data && data.deleted_date;

    return (
        <div style={isDeleted ? { opacity: 0.5 } : {}}>
            <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}</h2>
            {!data.html && data.description && <hr />}
            {data.description && (
                <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
            )}
            {data.html && <hr />}
            {data.html && (
                <HtmlView htmlContent={data.html} dark={dark} />
            )}
            <hr />
            {data.fk_users_id && data.fk_users_id !== "0" && (
                <p>👤 User: <UserLinkView user_id={data.fk_users_id} dark={dark} /></p>
            )}
            <p>
            {data.street}{data.street ? <br/> : ""}
            {data.zip} {data.city} {data.state ? `(${data.state})` : ''}{data.street || data.zip || data.city || data.state ? <br/> : ""}
            {data.fk_countrylist_id && (
                <CountryView country_id={data.fk_countrylist_id} dark={dark} />
            )}
            </p>
            {data.fk_companies_id && data.fk_companies_id !== "0" && (
                <p><ObjectLinkView obj_id={data.fk_companies_id} dark={dark} /></p>
            )}
            {data.phone && <p>📞 {data.phone}</p>}
            {data.office_phone && <p>🏢 {data.office_phone}</p>}
            {data.mobile && <p>📱 {data.mobile}</p>}
            {data.fax && <p>📠 {data.fax}</p>}
            {data.email && <p>✉️ <a href={`mailto:${data.email}`}>{data.email}</a></p>}
            {data.url && <p>🔗 <a href={data.url} target="_blank" rel="noopener noreferrer">{data.url}</a></p>}
            {data.codice_fiscale && <p>🆔 {data.codice_fiscale}</p>}
            {data.p_iva && <p>💰 {data.p_iva}</p>}
        </div>
    );
}

// Edit form for DBPerson
export function PersonEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    const { themeClass } = useContext(ThemeContext);
    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        street: data.street || '',
        zip: data.zip || '',
        city: data.city || '',
        state: data.state || '',
        fk_countrylist_id: data.fk_countrylist_id || '0',
        phone: data.phone || '',
        mobile: data.mobile || '',
        office_phone: data.office_phone || '',
        fax: data.fax || '',
        email: data.email || '',
        url: data.url || '',
        codice_fiscale: data.codice_fiscale || '',
        p_iva: data.p_iva || '',
        fk_users_id: data.fk_users_id || '0',
        fk_companies_id: data.fk_companies_id || '0',
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

            <Accordion className='mb-3' alwaysOpen>
                <Accordion.Item eventKey="0" className={themeClass}>
                    <Accordion.Header>{t('common.address')}</Accordion.Header>
                    <Accordion.Body>
                        <Form.Group className="mb-3">
                            <Form.Label>{t('common.street')}</Form.Label>
                            <Form.Control
                                type="text"
                                name="street"
                                value={formData.street}
                                onChange={handleChange}
                            />
                        </Form.Group>

                        <div className="row">
                            <div className="col-md-4">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.zip')}</Form.Label>
                                    <Form.Control
                                        type="text"
                                        name="zip"
                                        value={formData.zip}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                            <div className="col-md-8">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.city')}</Form.Label>
                                    <Form.Control
                                        type="text"
                                        name="city"
                                        value={formData.city}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                        </div>

                        <div className="row">
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.state')}</Form.Label>
                                    <Form.Control
                                        type="text"
                                        name="state"
                                        value={formData.state}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                            <div className="col-md-6">
                                <CountrySelector
                                    value={formData.fk_countrylist_id}
                                    onChange={handleChange}
                                    name="fk_countrylist_id"
                                />
                            </div>
                        </div>
                    </Accordion.Body>
                </Accordion.Item>
                <Accordion.Item eventKey="1" className={themeClass}>
                    <Accordion.Header>{t('common.contact_info')}</Accordion.Header>
                    <Accordion.Body>
                        <div className="row">
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.phone')}</Form.Label>
                                    <Form.Control
                                        type="tel"
                                        name="phone"
                                        value={formData.phone}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.mobile')}</Form.Label>
                                    <Form.Control
                                        type="tel"
                                        name="mobile"
                                        value={formData.mobile}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                        </div>

                        <div className="row">
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.office_phone')}</Form.Label>
                                    <Form.Control
                                        type="tel"
                                        name="office_phone"
                                        value={formData.office_phone}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.fax')}</Form.Label>
                                    <Form.Control
                                        type="tel"
                                        name="fax"
                                        value={formData.fax}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                        </div>

                        <Form.Group className="mb-3">
                            <Form.Label>{t('common.email')}</Form.Label>
                            <Form.Control
                                type="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                            />
                        </Form.Group>

                        <Form.Group className="mb-3">
                            <Form.Label>{t('common.website')}</Form.Label>
                            <Form.Control
                                type="url"
                                name="url"
                                value={formData.url}
                                onChange={handleChange}
                            />
                        </Form.Group>
                    </Accordion.Body>
                </Accordion.Item>
                <Accordion.Item eventKey="2" className={themeClass}>
                    <Accordion.Header>{t('common.tax_info')}</Accordion.Header>
                    <Accordion.Body>
                        <div className="row">
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.codice_fiscale')}</Form.Label>
                                    <Form.Control
                                        type="text"
                                        name="codice_fiscale"
                                        value={formData.codice_fiscale}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                            <div className="col-md-6">
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('common.p_iva')}</Form.Label>
                                    <Form.Control
                                        type="text"
                                        name="p_iva"
                                        value={formData.p_iva}
                                        onChange={handleChange}
                                    />
                                </Form.Group>
                            </div>
                        </div>
                    </Accordion.Body>
                </Accordion.Item>
                <Accordion.Item eventKey="3" className={themeClass}>
                    <Accordion.Header>{t('common.links')}</Accordion.Header>
                    <Accordion.Body>
                        <ObjectLinkSelector
                            value={formData.fk_users_id}
                            onChange={handleChange}
                            classname="DBUser"
                            fieldName="fk_users_id"
                            label={t('common.user_id')}
                            required={false}
                        />

                        <ObjectLinkSelector
                            value={formData.fk_companies_id}
                            onChange={handleChange}
                            classname="DBCompany"
                            fieldName="fk_companies_id"
                            label={t('common.company_id')}
                            required={false}
                        />
                    </Accordion.Body>
                </Accordion.Item>
            </Accordion>

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

export function People() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);
  // const [query, setQuery] = useState("");
  // const [editingFolder, setEditingFolder] = useState(null); // folder in editing

  const searchClassname = "DBPerson";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    // { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
    { name: t("dbobjects.fk_companies_id") || "Company", attribute: "fk_companies_id", type: "objectLink" },
    { name: t("dbobjects.fk_users_id") || "User", attribute: "fk_users_id", type: "userSelector" },
    { name: t("dbobjects.fk_countrylist_id") || "Country", attribute: "fk_countrylist_id", type: "countrySelector" },

    // { name: t("dbobjects.name") || "Name", attribute: "name2", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name3", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name4", type: "string" },
  ];

  const resultsColumns = [
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    // { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.fk_companies_id") || "Company", attribute: "fk_companies_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.fk_users_id") || "User", attribute: "fk_users_id", type: "userLink", hideOnSmall: true },
    { name: t("dbobjects.fk_countrylist_id") || "Country", attribute: "fk_countrylist_id", type: "countryView", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
