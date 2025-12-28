import React, { useState, useEffect } from 'react';
import { Form, Button, Spinner, InputGroup, ListGroup } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';
import { classname2bootstrapIcon } from './sitenavigation_utils';

/**
 * Generic object link selector component
 * Searches for objects by classname and displays results in a dropdown
 * 
 * @param {string} value - Current selected object ID
 * @param {function} onChange - Callback when selection changes
 * @param {string} classname - DBObject classname to search (e.g., "DBCompany", "DBUser")
 * @param {string} fieldName - Name of the foreign key field (e.g., "fk_companies_id")
 * @param {string} label - Label to display above the selector
 * @param {boolean} required - Whether the field is required
 */
function ObjectLinkSelector({ value, onChange, classname, fieldName, label, required = false, _type="link" }) {
    const { t } = useTranslation();
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(false);
    const [results, setResults] = useState([]);
    const [selectedObject, setSelectedObject] = useState(null);
    const [showResults, setShowResults] = useState(false);

    // Load selected object details on mount or when value changes
    useEffect(() => {
        if (value && value !== '0' && value !== '') {
            loadObjectDetails(value, classname);
        } else {
            setSelectedObject(null);
        }
    }, [value]);

    // Load object details by ID
    const loadObjectDetails = async (objectId, classname) => {
        try {
            const uri = classname == 'DBUser' ? `/users/${objectId}` : classname == 'DBGroup' ? `/groups/${objectId}` : `/content/${objectId}`;
            console.log("Loading details for "+uri);
            const response = await axiosInstance.get(uri);
            // if (response.data.success && response.data.data) { if you check .success it fails
            // alert(classname + "::" +JSON.stringify(response.data));
            console.log("Loaded details:", response.data);
            const _type = classname == 'DBUser' ? 'DBUser' : classname == 'DBGroup' ? 'DBGroup' : response.data.metadata.classname;
            const _name = classname == 'DBUser' ? response.data.login || 'Unnamed' : classname == 'DBGroup' ? response.data.name || 'Unnamed' : response.data.data.name || 'Unnamed';
            console.log("Determined type/name:", _type, _name);
            if (response.data.data || response.data.login || response.data.name) {
                setSelectedObject({
                    id: classname == 'DBUser' || classname == 'DBGroup' ? response.data.id : response.data.data.id,
                    name: _name,
                    classname: _type,
                });
            }
        } catch (error) {
            console.error('Failed to load object details:', error);
        }
    };

    // Search for objects
    const handleSearch = async (term) => {
        setSearchTerm(term);
        
        if (term.length < 2) {
            setResults([]);
            setShowResults(false);
            return;
        }

        var searchJson = null;
        // switch (classname) {
        //     case 'DBUser':
        //         // searchJson = { "login": term };
        //         break;
        //     case 'DBGroup':
        //         // Example: searchJson = { "name": term };
        //         break;
        //     default:
        //         // No special filters
        //         break;
        // }

        var searchFields = null;
        var orderBy = 'name';
        switch (classname) {
            case 'DBUser':
                searchFields = ['login', 'fullname', 'email'];
                orderBy = 'login ASC';
                break;
            case 'DBGroup':
                searchFields = ['name', 'description'];
                break;
            default:
                break;
        }

        setLoading(true);
        try {
            if (searchFields!==null) {
                searchJson = {
                    "$or": searchFields.map(field => ({
                        [field]: { "$like": `%${term}%` }
                    }))
                };
            }
            const response = await axiosInstance.get('/objects/search', {
                params: {
                    classname: classname,
                    name: term,
                    searchJson: JSON.stringify(searchJson),
                    orderBy: orderBy,
                    offset: 0,
                    limit: 20,
                    type: _type // "link" = I want only objects that I can attach to (write permission)
                }
            });

            if (response.data.success && response.data.objects) {
                // IF DBUser, then copy fullname into description
                if (classname === 'DBUser') {
                    response.data.objects = response.data.objects.map(obj => ({
                        ...obj,
                        name: obj.login || 'Unnamed',
                        description: obj.fullname || ''
                    }));
                }
                setResults(response.data.objects);
                setShowResults(true);
            }
        } catch (error) {
            console.error('Search failed:', error);
            setResults([]);
        } finally {
            setLoading(false);
        }
    };

    // Select an object from results
    const handleSelect = (object) => {
        setSelectedObject(object);
        setSearchTerm('');
        setResults([]);
        setShowResults(false);
        onChange({
            target: {
                name: fieldName,
                value: object.id
            }
        });
    };

    // Clear selection
    const handleClear = () => {
        setSelectedObject(null);
        setSearchTerm('');
        setResults([]);
        setShowResults(false);
        onChange({
            target: {
                name: fieldName,
                value: ''  // it was '0', but it's a problem when resetting the parent selector and then search for objects again
            }
        });
    };

    return (
        <Form.Group className="mb-3">
            {label && <Form.Label className="mb-1">{label}</Form.Label>}
            
            {selectedObject ? (
                <InputGroup>
                    <Button 
                        variant="outline-secondary"
                        disabled
                    >
                    <i className={`bi bi-${classname2bootstrapIcon(selectedObject.classname)}`} title={selectedObject.classname}></i>&nbsp;
                    </Button>
                    <Form.Control
                        type="text"
                        value={selectedObject.name}
                        readOnly
                    />
                    <Button 
                        variant="outline-secondary"
                        onClick={handleClear}
                    >
                        <i className="bi bi-x-lg"></i>
                    </Button>
                </InputGroup>
            ) : (
                <div className="position-relative">
                    <InputGroup>
                        <Form.Control
                            type="text"
                            value={searchTerm}
                            onChange={(e) => handleSearch(e.target.value)}
                            placeholder={t('common.search') || 'Search...'}
                            required={required}
                        />
                        {loading && (
                            <InputGroup.Text>
                                <Spinner animation="border" size="sm" />
                            </InputGroup.Text>
                        )}
                    </InputGroup>
                    
                    {showResults && results.length > 0 && (
                        <ListGroup className="position-absolute w-100" style={{ zIndex: 1000 }}>
                            {results.map((obj) => (
                                <ListGroup.Item
                                    key={obj.id}
                                    action
                                    onClick={() => handleSelect(obj)}
                                    style={{ cursor: 'pointer' }}
                                >
                                    <div className="d-flex justify-content-between align-items-start">
                                        <div>
                                            <i className={`bi bi-${classname2bootstrapIcon(obj.classname)}`} title={obj.classname}></i>
                                            &nbsp;
                                            <strong>{obj.name}</strong>
                                            {obj.description && (
                                                <div className="text-secondary small">
                                                    {obj.description.substring(0, 100)}
                                                    {obj.description.length > 100 && '...'}
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </ListGroup.Item>
                            ))}
                        </ListGroup>
                    )}
                    
                    {showResults && results.length === 0 && searchTerm.length >= 2 && !loading && (
                        <ListGroup className="position-absolute w-100" style={{ zIndex: 1000 }}>
                            <ListGroup.Item className="text-secondary">
                                {t('common.no_results') || 'No results found'}
                            </ListGroup.Item>
                        </ListGroup>
                    )}
                </div>
            )}
            
            {!selectedObject && <Form.Text className="text-secondary">
                {t('common.search_hint') || 'Type at least 2 characters to search'}
            </Form.Text>}
        </Form.Group>
    );
}

export default ObjectLinkSelector;
