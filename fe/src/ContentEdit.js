import React, { useContext, useEffect, useState } from 'react';
import { Card, Container, Button, Spinner, Alert } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import 'react-quill/dist/quill.snow.css';
import axiosInstance from './axios';
import { CompanyEdit } from './DBCompany';
import { EventEdit } from './DBEvent';
import { FileEdit } from './DBFile';
import { FolderEdit } from './DBFolders';
import { LinkEdit } from './DBLink';
import { NoteEdit } from './DBNote';
import { ObjectEdit } from './DBObject';
import { PageEdit } from './DBPage';
import { PersonEdit } from './DBPeople';
import { ThemeContext } from './ThemeContext';


// Main ContentEdit component
function ContentEdit() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { dark } = useContext(ThemeContext);

    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState(null);
    const [data, setData] = useState(null);
    const [metadata, setMetadata] = useState(null);

    useEffect(() => {
        const loadObject = async () => {
            try {
                setLoading(true);
                setError(null);

                // Load object data
                const response = await axiosInstance.get(`/content/${id}`);
                setData(response.data.data);
                setMetadata(response.data.metadata);
            } catch (err) {
                console.error('Error loading object:', err);
                setError(err.response?.data?.message || 'Failed to load object');
            } finally {
                setLoading(false);
            }
        };

        loadObject();
    }, [id]);

    const handleSave = async (formData, isMultipart = false) => {
        try {
            setSaving(true);
            setError(null);

            if (isMultipart) {
                // Upload with multipart/form-data for file uploads
                await axiosInstance.put(`/objects/${id}`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            } else {
                // Regular JSON update
                await axiosInstance.put(`/objects/${id}`, formData);
            }

            // Navigate back to view mode
            navigate(`/c/${id}`);
        } catch (err) {
            console.error('Error saving object:', err);
            setError(err.response?.data?.message || 'Failed to save changes');
            setSaving(false);
        }
    };

    const handleCancel = () => {
        // navigate(`/c/${id}`);
        // Navigate back to previous page
        navigate(-1);
    };

    const handleDelete = async () => {
        if (!window.confirm(t('navigation.delete_confirm'))) {
            return;
        }

        try {
            setSaving(true);
            setError(null);

            // Delete object via API
            await axiosInstance.delete(`/objects/${id}`);

            // Navigate to parent or home
            // alert("father_id=" + data.father_id);
            if (data.father_id) {
                navigate(`/c/${data.father_id}`);
                return;
            }
            navigate(-1);
        } catch (err) {
            console.error('Error deleting object:', err);
            setError(err.response?.data?.message || 'Failed to delete object');
            setSaving(false);
        }
    };

    if (loading) {
        return (
            <Container className="mt-4 text-center">
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </Container>
        );
    }

    if (error && !data) {
        return (
            <Container className="mt-4">
                <Alert variant="danger">
                    <Alert.Heading>Error</Alert.Heading>
                    <p>{error}</p>
                    <Button variant="outline-danger" onClick={() => navigate(-1)}>
                        Go Back
                    </Button>
                </Alert>
            </Container>
        );
    }

    if (!data || !metadata) {
        return null;
    }

    // Check if user can edit
    if (!metadata.can_edit) {
        return (
            <Container className="mt-4">
                <Alert variant="warning">
                    <Alert.Heading>Access Denied</Alert.Heading>
                    <p>You don't have permission to edit this object.</p>
                    <Button variant="outline-warning" onClick={() => navigate(`/c/${id}`)}>
                        View Object
                    </Button>
                </Alert>
            </Container>
        );
    }

    const classname = metadata.classname;

    // Render appropriate edit form based on classname
    let EditComponent;
    switch (classname) {
        case 'DBNote':
            EditComponent = NoteEdit;
            break;
        case 'DBNews':
        case 'DBPage':
            EditComponent = PageEdit;
            break;
        case 'DBPerson':
            EditComponent = PersonEdit;
            break;
        case 'DBCompany':
            EditComponent = CompanyEdit;
            break;
        case 'DBEvent':
            EditComponent = EventEdit;
            break;
        case 'DBFile':
            EditComponent = FileEdit;
            break;
        case 'DBFolder':
            EditComponent = FolderEdit;
            break;
        case 'DBLink':
            EditComponent = LinkEdit;
            break;
        default:
            EditComponent = ObjectEdit;
            break;
    }

    return (
        <Container className="mt-4">
            <Card bg={dark ? 'dark' : 'light'} text={dark ? 'light' : 'dark'}>
                <Card.Header className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
                    <h2 className={dark ? 'text-light' : 'text-dark'}>
                        <i className="bi bi-pencil me-2"></i>
                        {t('common.edit')}: {data.name}
                    </h2>
                    <small style={{ opacity: 0.7 }}>
                        {classname} · ID: {id}
                    </small>
                </Card.Header>
                <Card.Body className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderBottom: '0px solid rgba(255,255,255,0.1)' } : {}}>
                    <EditComponent
                        data={data}
                        metadata={metadata}
                        onSave={handleSave}
                        onCancel={handleCancel}
                        onDelete={handleDelete}
                        saving={saving}
                        error={error}
                        dark={dark}
                    />
                </Card.Body>
            </Card>
        </Container>
    );
}

export default ContentEdit;
