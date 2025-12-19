import React, { useState, useEffect } from 'react';
import { Button, Dropdown, DropdownButton, Spinner } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';

function NewObjectButton({ fatherId, onObjectCreated }) {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [creatableTypes, setCreatableTypes] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadCreatableTypes();
    }, [fatherId]);

    const loadCreatableTypes = async () => {
        try {
            setLoading(true);
            const params = fatherId ? `?father_id=${fatherId}` : '';
            const response = await axiosInstance.get(`/objects/creatable-types${params}`);
            setCreatableTypes(response.data.types || []);
        } catch (err) {
            console.error('Error loading creatable types:', err);
            setCreatableTypes([]);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateObject = async (classname) => {
        try {
            // Create minimal object
            const payload = {
                classname: classname,
                father_id: fatherId || "0",
                name: `New ${classname.replace('DB', '')}`,
                description: ""
            };

            switch (classname) {
                case 'DBLink':
                    payload.href = "";
                    break;
                default:
                    break;
            }

            const response = await axiosInstance.post('/objects', payload);
            const newObjectId = response.data.data.id;

            // Notify parent (for refreshing children list)
            if (onObjectCreated) {
                onObjectCreated(newObjectId);
            }

            // Navigate to edit page
            navigate(`/e/${newObjectId}`);
        } catch (err) {
            console.error('Error creating object:', err);
            alert('Failed to create object: ' + (err.response?.data?.error || err.message));
        }
    };

    if (loading) {
        return (
            <Button variant="success" size="sm" disabled>
                <Spinner animation="border" size="sm" className="me-1" />
                {t('common.loading')}
            </Button>
        );
    }

    if (creatableTypes.length === 0) {
        return null; // No types to create
    }

    return (
        <DropdownButton
            id="new-object-dropdown"
            variant="success"
            size="sm"
            title={<><i className="bi bi-plus-lg me-1"></i>{t('common.new')}</>}
            align="end"
        >
            {creatableTypes.map((classname) => (
                <Dropdown.Item
                    key={classname}
                    onClick={() => handleCreateObject(classname)}
                >
                    <i className={`bi bi-file-earmark-plus me-2`}></i>
                    {classname.replace('DB', '')}
                </Dropdown.Item>
            ))}
        </DropdownButton>
    );
}

export default NewObjectButton;
