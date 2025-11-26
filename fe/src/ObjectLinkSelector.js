import React, { useState, useEffect } from 'react';
import { Form, Button, Spinner, InputGroup, ListGroup } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';

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
function ObjectLinkSelector({ value, onChange, classname, fieldName, label, required = false }) {
    const { t } = useTranslation();
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(false);
    const [results, setResults] = useState([]);
    const [selectedObject, setSelectedObject] = useState(null);
    const [showResults, setShowResults] = useState(false);

    // Load selected object details on mount or when value changes
    useEffect(() => {
        if (value && value !== '0' && value !== '') {
            loadObjectDetails(value);
        } else {
            setSelectedObject(null);
        }
    }, [value]);

    // Load object details by ID
    const loadObjectDetails = async (objectId) => {
        try {
            const response = await axiosInstance.get(`/content/${objectId}`);
            // if (response.data.success && response.data.data) { if you check .success it fails
            if (response.data.data) {
                setSelectedObject({
                    id: response.data.data.id,
                    name: response.data.data.name || 'Unnamed',
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

        setLoading(true);
        try {
            const response = await axiosInstance.get('/objects/search', {
                params: {
                    classname: classname,
                    name: term,
                    limit: 20
                }
            });

            if (response.data.success && response.data.objects) {
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
                value: '0'
            }
        });
    };

    return (
        <Form.Group className="mb-3">
            <Form.Label>{label}</Form.Label>
            
            {selectedObject ? (
                <InputGroup>
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
                                            <strong>{obj.name}</strong>
                                            {obj.description && (
                                                <div className="text-muted small">
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
                            <ListGroup.Item className="text-muted">
                                {t('common.no_results') || 'No results found'}
                            </ListGroup.Item>
                        </ListGroup>
                    )}
                </div>
            )}
            
            <Form.Text className="text-muted">
                {t('common.search_hint') || 'Type at least 2 characters to search'}
            </Form.Text>
        </Form.Group>
    );
}

export default ObjectLinkSelector;
