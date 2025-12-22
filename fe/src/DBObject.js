import React, { useState, useEffect } from 'react';
import { Alert, Button, Card, Col, Form, Row, Spinner } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
} from './sitenavigation_utils';
import {
  CountryView,
  CountrySelector,
  GroupLinkView,
  ImageView,
  LanguageSelector,
  LanguageView,
  ObjectLinkView,
  UserLinkView,
} from './ContentWidgets';
import ObjectLinkSelector from './ObjectLinkSelector'
import PermissionsEditor from './PermissionsEditor';
import { getErrorMessage } from "./errorHandler";
import { HtmlView } from './ContentHtml';
import axios from './axios';


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
                    <HtmlView htmlContent={data.html} dark={dark} />
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

export function CheckWritePermission({objectData}) {
    if (!objectData || !objectData.permissions) {
        return false;
    }
    const permissions = objectData.permissions;
    const user_id = localStorage.getItem("user_id");
    // const user_group_id = localStorage.getItem("group_id");
    const group_ids = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];

    // localStorage.setItem("username", login);
    // localStorage.setItem("user_id", res.data.user_id);
    // localStorage.setItem("groups", JSON.stringify(res.data.groups));
    
    return false
         || (objectData.owner === user_id && permissions[1] === 'w') // Is user owner and has user write permission
         || (group_ids.indexOf(objectData.group_id) !== -1 && permissions[4] === 'w') // Is user in group and has group write permission
         || permissions[7] === 'w';
}

export function ObjectSearch({searchClassname, searchColumns, resultsColumns, orderBy, dark, themeClass}) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [searchFormData, setSearchFormData] = useState({
    father_id: "", //"0", // root folders
    name: "",
    description: "",
  });
  const [searchOrderBy, setSearchOrderBy] = useState(orderBy || "name");
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const [includeDeleted, setIncludeDeleted] = useState(false);
  const [limit, setLimit] = useState(5); // Change to 20 for production
  const [offset, setOffset] = useState(0);

  // Load folders on start
  useEffect(() => {
    // fetchObjects(searchFormData);
  }, []);

  // Auto-dismiss error message after 5 seconds
  useEffect(() => {
    if (errorMessage) {
      const timer = setTimeout(() => {
        setErrorMessage("");
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [errorMessage]);

  const fetchObjects = async (search, orderBy, newOffset) => {
    const token = localStorage.getItem("token");
    setLoading(true);
    setErrorMessage("");

    try {
      console.log('Searching objects with:', search, ' order by:', orderBy);
      const response = await axios.get('/objects/search', {
        params: {
          classname: searchClassname,
          searchJson: JSON.stringify(search),
          limit: limit,
          offset: newOffset,
          orderBy: orderBy,
          includeDeleted: includeDeleted ? "true" : "false"
        },
      });
      // console.log('Search response:', response.data);
      // Backend returns array directly, not wrapped in results
      // IF offset is 0, replace results, else append
      if (newOffset === 0) {
        setResults(Array.isArray(response.data) ? response.data : response.data.objects || []);
      } else {
        setResults(results => [...results, ...(Array.isArray(response.data) ? response.data : response.data.objects || [])]);
      }
      setOffset(newOffset + limit);
    } catch (err) {
      console.error('Search error:', err);
      setErrorMessage(err.response?.data?.error || 'Search failed');
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setSearchFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchObjects(searchFormData, searchOrderBy, 0);
  };

  const handleCreateObject = async (classname) => {
      try {
          // Create minimal object
          const payload = {
              classname: classname,
              // father_id: fatherId || "0",
              name: `New ${classname.replace('DB', '')}`,
              description: ""
          };

          const response = await axios.post('/objects', payload);
          const newObjectId = response.data.data.id;

          // Notify parent (for refreshing children list)
          // if (onObjectCreated) {
          //     onObjectCreated(newObjectId);
          // }

          // Navigate to edit page
          navigate(`/e/${newObjectId}`);
      } catch (err) {
          console.error('Error creating object:', err);
          alert('Failed to create object: ' + (err.response?.data?.error || err.message));
      }
  };

  return (
    <div className={`container mt-3 ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>{t("dbobjects."+searchClassname)}</h2>

      {/* Search form */}
      <Form id="searchForm" onSubmit={handleSearch}>
        <Row>
        { searchColumns.map((col, index) => (
            <Col xs={12} md={6} lg={4} xl={3}>
            <Form.Group className="mb-3">
              {
                col.type === "objectLink" ? (
                  <>
                    <ObjectLinkSelector
                        value={searchFormData[col.attribute] || '0'}
                        onChange={handleInputChange}
                        classname="DBObject"
                        fieldName={col.attribute}
                        label={t('dbobjects.' + col.attribute)}
                        _type="search"
                    />
                  </>
                ) : col.type === "countrySelector" ? (
                  <>
                    <CountrySelector
                        value={searchFormData[col.attribute] || ''}
                        onChange={handleInputChange}
                        fieldName={col.attribute}
                        label={t('dbobjects.' + col.attribute)}
                        _type="search"
                    />
                  </>
                ) : col.type === "languageSelector" ? (
                  <>
                    <LanguageSelector
                        value={searchFormData[col.attribute] || ''}
                        onChange={handleInputChange}
                        fieldName={col.attribute}
                        label={t('dbobjects.' + col.attribute)}
                    />
                  </>
                ) : col.type === "userSelector" ? (
                  <>
                    <ObjectLinkSelector
                        value={searchFormData[col.attribute] || '0'}
                        onChange={handleInputChange}
                        classname="DBUser"
                        fieldName={col.attribute}
                        label={t('dbobjects.' + col.attribute)}
                        _type="search"
                    />
                  </>
                ) : (
                  <>
                  <Form.Label>{t('dbobjects.' + col.attribute)}</Form.Label>
                  <Form.Control
                      type="text"
                      name={col.attribute}
                      value={searchFormData[col.attribute] || ''}
                      onChange={handleInputChange}
                      onSubmit={handleSearch}
                  />
                  </>
                )
              }
            </Form.Group>
            </Col>
        ))}
        </Row>
      </Form>
      <div className="row">
        <div className="col-md-6 text-center text-md-start">
          
          <button className="btn btn-success mb-3" onClick={() => { handleCreateObject(searchClassname); }} >{t("common.new")}</button>
          {/* <button className="btn btn-success mb-3" onClick={() => { navigate('/folders/new'); }} >{t("common.new")}</button> */}
        </div>
        <div className="col-md-6 text-center text-md-end">
          <button type="submit" form="searchForm" className="btn btn-primary">{t("common.search")}</button>
        </div>
      </div>


      {results.length > 0 && (
        <table 
        className={`table ${dark ? "table-dark" : "table-striped"} table-hover`}
        >
          <thead>
            <tr>
              <th className="d-none d-md-table-cell">#</th>
              {resultsColumns.map((col, index) => (
                <th className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={index}
                  onClick={() => {
                    const desc = searchOrderBy === col.attribute ? ' desc' : '';
                    setSearchOrderBy(col.attribute + desc);
                    console.log("Setting order by:", searchOrderBy);
                    // // sort results
                    // results.sort((a, b) => {
                    //   if (a[col.attribute] < b[col.attribute]) return desc ? 1 : -1;
                    //   if (a[col.attribute] > b[col.attribute]) return desc ? -1 : 1;
                    //   return 0;
                    // });
                    // setResults([...results]);
                    fetchObjects(searchFormData, col.attribute + desc, 0);
                  }}
                >{col.name}{searchOrderBy === col.attribute+' desc' ? " ▼" : searchOrderBy === col.attribute ? " ▲" : ""}</th>
              ))}
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {results.map((result, index) => (
              <tr key={result.id}>
                <td className="d-none d-md-table-cell">{index+1}</td>
                {resultsColumns.map((col, cindex) => (
                  col.type === "objectLink" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                        <ObjectLinkView obj_id={result[col.attribute]} dark={dark} />
                        </td>
                    ) : col.type === "userLink" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                        <UserLinkView user_id={result[col.attribute]} dark={dark} />
                        </td>
                    ) : col.type === "groupLink" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                            <GroupLinkView group_id={result[col.attribute]} dark={dark} />
                        </td>
                    ) : col.type === "imageView" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                            <ImageView id={result[col.attribute]} title={t("dbobjects.DBFile")} thumbnail={true}  style={{ fontSize: '2rem', minHeight: '2rem', maxWidth: '50px', maxHeight: '50px', borderRadius: '0.5rem' }} />
                        </td>
                    ) : col.type === "languageView" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                            <LanguageView language={result[col.attribute]} short={true} />
                        </td>
                    ) : col.type === "countryView" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                        <CountryView country_id={result[col.attribute]} dark={dark} />
                        </td>
                    ) : col.type === "dateTime" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>{formateDateTimeString(result[col.attribute])}</td>
                    ) : col.type === "urlView" ? (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>
                          {(() => {
                            const url = result[col.attribute].startsWith('http://') || result[col.attribute].startsWith('https://') ? result[col.attribute] : 'http://' + result[col.attribute];
                            return (
                              <a href={url} target="_blank" rel="noopener noreferrer">
                                {result[col.attribute]}
                              </a>
                            );
                          })()}
                        </td>
                    ) : (
                        <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>{result[col.attribute] ? result[col.attribute].slice(0, 30) + (result[col.attribute].length > 30 ? "..." : "") : ""}</td>
                    )
                ))}
                <td>
                  {CheckWritePermission({objectData: result}) && <button
                    className="btn btn-sm btn-warning"
                    onClick={() => navigate(`/e/${result.id}`)}
                  >
                    {t("common.edit")}
                  </button>
                   }
                  <button
                    className="btn btn-sm btn-warning ms-2"
                    onClick={() => navigate(`/c/${result.id}`)}
                  >
                    {t("common.view")}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
        {/* Show load more button if results length equals limit */}
        <div className="text-center mb-3">
          {results.length > 0 && (results.length % limit) === 0 ? (
            <Button
              variant="outline-primary"
              className="mt-2 ms-auto"
              onClick={() => {
                // Load more results
                fetchObjects(searchFormData, searchOrderBy, offset);
              }}
            >
              {t("common.load_more") || "Load More"}
            </Button>
          ) : null}
      </div>
    </div>
  );
}

