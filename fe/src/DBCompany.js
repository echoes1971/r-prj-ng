import React, { useContext, useEffect, useState } from "react";
import { Accordion, Card, Container, Form, Button, Spinner, Alert, ButtonGroup, Overlay, Popover } from 'react-bootstrap';
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import ObjectLinkSelector from './ObjectLinkSelector';
import { ObjectSearch } from "./DBObject";
import PermissionsEditor from './PermissionsEditor';
import { CountrySelector, CountryView } from './ContentWidgets'


export function CompanyView({ data, metadata, objectData, dark }) {
    // const navigate = useNavigate();
    // const { t } = useTranslation();
    
    return (
        <div>
            <h2 className={dark ? 'text-light' : 'text-dark'}>{data.name}</h2>
            {!data.html && data.description && <hr />}
            {data.description && (
                <Card.Text dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></Card.Text>
            )}
            {data.html && <hr />}
            {data.html && (
                <HtmlView htmlContent={data.html} dark={dark} />
            )}
            <hr />
            <p>
            {data.street}<br/>
            {data.zip} {data.city} ({data.state})<br/>
            <CountryView country_id={data.fk_countrylist_id} dark={dark} />
            </p>
            {data.phone && <p>📞 {data.phone}</p>}
            {data.office_phone && <p>🏢 {data.office_phone}</p>}
            {data.mobile && <p>📱 {data.mobile}</p>}
            {data.fax && <p>📠 {data.fax}</p>}
            {data.email && <p>✉️ <a href={`mailto:${data.email}`}>{data.email}</a></p>}
            {data.url && <p>🔗 <a href={data.url} target="_blank" rel="noopener noreferrer">{data.url}</a></p>}
            {data.p_iva && <p>💰 {data.p_iva}</p>}
        </div>
    );
}

// Edit form for DBCompany
export function CompanyEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
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
        fax: data.fax || '',
        email: data.email || '',
        url: data.url || '',
        p_iva: data.p_iva || '',
        permissions: data.permissions || 'rwxr-x---',
        father_id: data.father_id || '0',
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
                        <Form.Group className="mb-3">
                            <Form.Label>{t('common.p_iva')}</Form.Label>
                            <Form.Control
                                type="text"
                                name="p_iva"
                                value={formData.p_iva}
                                onChange={handleChange}
                            />
                        </Form.Group>
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


export function Companies() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);
  // const [query, setQuery] = useState("");
  // const [editingFolder, setEditingFolder] = useState(null); // folder in editing

  const searchClassname = "DBCompany";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.fk_countrylist_id") || "Country", attribute: "fk_countrylist_id", type: "countrySelector" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },

    // { name: t("dbobjects.name") || "Name", attribute: "name2", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name3", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name4", type: "string" },
  ];

  const resultsColumns = [
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    // { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    // { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
    { name: t("dbobjects.phone") || "Phone", attribute: "phone", type: "string", hideOnSmall: false },
    { name: t("dbobjects.url") || "URL", attribute: "url", type: "urlView", hideOnSmall: true },
    { name: t("dbobjects.fk_countrylist_id") || "Country", attribute: "fk_countrylist_id", type: "countryView", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
