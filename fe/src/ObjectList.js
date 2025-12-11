import React, { useState, useContext } from 'react';
import { ListGroup, Card, Row, Col, Button, ButtonGroup } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import {
  classname2bootstrapIcon,
 } from './sitenavigation_utils';
 import {
  ImageView
 } from './ContentWidgets';

/**
 * ObjectList - Reusable component to display a list of objects
 * Supports both list and card view modes, with persistent preference in localStorage
 * 
 * @param {Array} items - Array of objects to display
 * @param {Function} onItemClick - Optional custom click handler, receives (item)
 * @param {boolean} showViewToggle - Show/hide the view mode toggle buttons (default: true)
 * @param {string} storageKey - localStorage key for view mode preference (default: 'objectListViewMode')
 * @param {string} defaultView - Default view mode: 'list' or 'cards' (default: 'list')
 * @param {boolean} dark - Whether to use dark theme styles
 */
function ObjectList({ 
  items = [], 
  onItemClick = null, 
  showViewToggle = true,
  storageKey = 'objectListViewMode',
  defaultView = 'list',
  dark = false
}) {
  const navigate = useNavigate();
  // const { dark, themeClass } = useContext(ThemeContext);
  const [viewMode, setViewMode] = useState(
    localStorage.getItem(storageKey) || defaultView
  );

  const handleViewModeChange = (mode) => {
    setViewMode(mode);
    localStorage.setItem(storageKey, mode);
  };

  const handleItemClick = (item) => {
    if (onItemClick) {
      onItemClick(item);
    } else {
      // Default behavior: navigate to /c/{id}
      navigate(`/c/${item.id}`);
    }
  };

  if (!items || items.length === 0) {
    return null;
  }

  return (
    <>
      {showViewToggle && (
        <div className="d-flex justify-content-end mb-3">
          <ButtonGroup size="sm">
            <Button 
              variant={viewMode === 'list' ? 'primary' : 'outline-secondary'}
              onClick={() => handleViewModeChange('list')}
            >
              <i className="bi bi-list-ul"></i>
            </Button>
            <Button 
              variant={viewMode === 'cards' ? 'primary' : 'outline-secondary'}
              onClick={() => handleViewModeChange('cards')}
            >
              <i className="bi bi-grid-3x3-gap"></i>
            </Button>
          </ButtonGroup>
        </div>
      )}
      
      {viewMode === 'list' ? (
        <ListGroup variant={dark ? 'dark' : undefined}>
          {items.map((item) => (
            <ListGroup.Item
              key={item.id}
              action
              onClick={() => handleItemClick(item)}
              style={{ cursor: 'pointer' }}
              variant={dark ? 'dark' : undefined}
            >
              <div className="d-flex justify-content-between align-items-center">
                <div>
                  <strong>{item.name || 'Untitled'}</strong>
                  {item.description && (
                    <div className="small" style={{ opacity: 0.7 }}>
                      {item.description.length > 200
                        ? item.description.substring(0, 200) + '...'
                        : item.description}
                    </div>
                  )}
                </div>
                { item.classname === 'DBFile' && (
                  <ImageView id={item.id} title={item.name || 'Image'} thumbnail={true} style={{ fontSize: '2rem', minHeight: '2rem', maxWidth: '50px', maxHeight: '50px', borderRadius: '0.5rem' }} />
                )}
                { item.classname !== 'DBFile' && (
                // <span className="badge bg-secondary">
                  <i 
                    className={`bi bi-${classname2bootstrapIcon(item.classname)}`} 
                    title={item.classname}
                    style={{ fontSize: '2rem' }}
                  ></i>
                // </span>
                )}
              </div>
            </ListGroup.Item>
          ))}
        </ListGroup>
      ) : (
        <Row>
          {items.map((item) => (
            <Col key={item.id} xs={12} md={6} lg={4} className="mb-3">
              <Card 
                className="h-100"
                style={{ cursor: 'pointer' }}
                onClick={() => handleItemClick(item)}
                data-bs-theme={dark? 'dark' : 'light'}
              >
                <Card.Body
                className={dark ? 'bg-secondary bg-opacity-25' : ''}
                >
                  <div className="d-flex justify-content-between align-items-start mb-2">
                    { item.classname !== 'DBFile' && (
                      <i 
                      className={`bi bi-${classname2bootstrapIcon(item.classname)}`}
                      style={{ fontSize: '2rem' }}
                    ></i>
                    )}
                    { item.classname === 'DBFile' && (
                      <ImageView id={item.id} title={item.name || 'Image'} thumbnail={false} style={{ fontSize: '2rem', minHeight: '2rem', maxWidth: '100px', maxHeight: '100px', borderRadius: '0.25rem' }} />
                    )}
                    <span className="badge bg-secondary">
                      {item.classname}
                    </span>
                  </div>
                  
                  <Card.Title className="mb-2">
                    {item.name || 'Untitled'}
                  </Card.Title>
                  
                  {item.description && (
                    <Card.Text className="text-muted small">
                      {item.description.length > 150
                        ? item.description.substring(0, 150) + '...'
                        : item.description}
                    </Card.Text>
                  )}
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </>
  );
}

export default ObjectList;
