import React, { useState, useEffect } from 'react';
import { Form, Spinner } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';
import { getFlagEmoji } from './sitenavigation_utils';

function CountrySelector({ value, onChange, name, required }) {
    const { t } = useTranslation();
    const [countries, setCountries] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadCountries();
    }, []);

    const loadCountries = async () => {
        try {
            setLoading(true);
            const response = await axiosInstance.get('/countries');
            setCountries(response.data.countries || []);
        } catch (err) {
            console.error('Error loading countries:', err);
            setCountries([]);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <Form.Group className="mb-3">
                <Form.Label>{t('common.country')}</Form.Label>
                <div className="d-flex align-items-center">
                    <Spinner animation="border" size="sm" className="me-2" />
                    <span>{t('common.loading')}</span>
                </div>
            </Form.Group>
        );
    }

    return (
        <Form.Group className="mb-3">
            <Form.Label>{t('common.country')}</Form.Label>
            <Form.Select
                name={name || 'fk_countrylist_id'}
                value={value || '0'}
                onChange={onChange}
                required={required}
            >
                <option value="0">-- {t('common.select')} --</option>
                {countries.map((country) => (
                    <option key={country.id} value={country.id}>
                        {getFlagEmoji(country.ISO_3166_1_2_Letter_Code)} {country.Common_Name}
                    </option>
                ))}
            </Form.Select>
        </Form.Group>
    );
}

export default CountrySelector;
