import React, { useContext, useEffect, useState } from "react";
import { Form, Row, Col } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import api from "./axios";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import ObjectLinkSelector from "./ObjectLinkSelector";
import { GroupLinkView, ObjectLinkView, UserLinkView } from "./sitenavigation_utils";
import AssociationManager from "./AssociationManager";
import { getErrorMessage } from "./errorHandler";
import axios from './axios';


export function Folders() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [searchFormData, setSearchFormData] = useState({
    father_id: "", //"0", // root folders
    name: "",
    description: "",
  });
  const [query, setQuery] = useState("");
    // const [editingFolder, setEditingFolder] = useState(null); // folder in editing
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

  const searchClassname = "DBFolder";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },

    // { name: t("dbobjects.name") || "Name", attribute: "name2", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name3", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name4", type: "string" },
  ];

  const resultsColumns = [
    { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
  ]

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

  const fetchObjects = async (search) => {
    const token = localStorage.getItem("token");
    setLoading(true);
    setErrorMessage("");

    try {
      const response = await axios.get('/objects/search', {
        params: {
          classname: searchClassname,
          searchJson: JSON.stringify(search)
        },
      });
      console.log('Search response:', response.data);
      // Backend returns array directly, not wrapped in results
      setResults(Array.isArray(response.data) ? response.data : response.data.objects || []);
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
    fetchObjects(searchFormData);
  };


  return (
    <div className={`container mt-3 ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>{t("folder.folders")}</h2>

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
          <button className="btn btn-success mb-3" onClick={() => { navigate('/folders/new'); }} >{t("common.new")}</button>
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
                <th className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={index}>{col.name}</th>
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
                  ) : (
                    <td className={col.hideOnSmall ? "d-none d-md-table-cell" : ""} key={cindex}>{result[col.attribute]}</td>
                  )
                ))}
                <td>
                  <button
                    className="btn btn-sm btn-warning"
                    onClick={() => navigate(`/e/${result.id}`)}
                  >
                    {t("common.edit")}
                  </button>
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

    </div>
  );
}
