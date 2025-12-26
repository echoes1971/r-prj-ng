import React, { useState, useEffect } from 'react';
import { Form, Spinner } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';
import { classname2bootstrapIcon } from './sitenavigation_utils';


// Convert ISO 3166-1 alpha-2 code to flag emoji
export function getFlagEmoji(countryCode) {
    if (!countryCode || countryCode.length !== 2) return '';
    const codePoints = countryCode
        .toUpperCase()
        .split('')
        .map(char => 127397 + char.charCodeAt());
    return String.fromCodePoint(...codePoints);
}

// Component: Display country with flag emoji
export function CountryView({ country_id, dark }) {
    const [country, setCountry] = useState(null);

    useEffect(() => {
        const fetchCountry = async () => {
            // Read from localStorage every time (avoid stale state)
            const stored = localStorage.getItem('countries_cache');
            const cacheData = stored ? JSON.parse(stored) : null;

            // Check if cache exists and is not expired (24 hours)
            const now = Date.now();
            const CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours in milliseconds
            
            let countries = {};
            if (cacheData && cacheData.expires_at && cacheData.expires_at > now) {
                // Cache is valid
                countries = cacheData.countries || {};
                const remainingSeconds = Math.floor((cacheData.expires_at - now) / 1000);
                console.log(`Countries cache VALID - expires in ${remainingSeconds}s at:`, new Date(cacheData.expires_at).toLocaleTimeString());
            } else {
                // Cache expired or doesn't exist - will be recreated
                if (cacheData?.expires_at) {
                    console.log('Countries cache EXPIRED at:', new Date(cacheData.expires_at).toLocaleTimeString(), 'now:', new Date(now).toLocaleTimeString());
                } else {
                    console.log('Countries cache MISSING, will rebuild');
                }
            }
            
            if (countries[country_id]) {
                setCountry(countries[country_id]);
                console.log('Loaded country from cache: ', country_id, "=", countries[country_id].Common_Name);
                return;
            }
            
            // Country not in cache, fetch from backend
            try {
                const response = await axiosInstance.get(`/content/country/${country_id}`);
                setCountry(response.data);
                
                // Update cache - re-read to avoid race conditions
                const currentStored = localStorage.getItem('countries_cache');
                const currentCache = currentStored ? JSON.parse(currentStored) : null;
                
                // Preserve expiry if cache is still valid, otherwise create new expiry
                let expiresAt = now + CACHE_DURATION;
                // TESTING: Force new expiry every time (comment out to preserve existing expiry)
                /*
                if (currentCache && currentCache.expires_at && currentCache.expires_at > now) {
                    expiresAt = currentCache.expires_at; // Keep existing expiry
                }
                */
                
                const updatedCache = {
                    expires_at: expiresAt,
                    countries: {
                        ...(currentCache?.countries || {}),
                        [country_id]: response.data
                    }
                };
                
                localStorage.setItem('countries_cache', JSON.stringify(updatedCache));
                console.log('Fetched and cached country: ', country_id, "=", response.data);
            } catch (error) {
                console.error('Error fetching country:', error);
            }
        }

        if (country_id && country_id !== "0") {
            fetchCountry();
        }
    }, [country_id]);

    if (!country_id || country_id === "0") {
        return null;
    }

    if (!country) {
        return <>Loading...</>;
    }

    const flag = getFlagEmoji(country.ISO_3166_1_2_Letter_Code);
    
    return (
        <>
            {flag && <span style={{ fontSize: '1.2em', marginRight: '0.3em' }}>{flag}</span>}
            {country.Common_Name}
        </>
    );
}
export function CountrySelector({ value, onChange, name, required }) {
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

// Component: Link to user profile
export function UserLinkView({ user_id, dark }) {
    const [user, setUser] = useState(null);

    useEffect(() => {
        const fetchUser = async () => {
            // Read from localStorage every time (avoid stale state)
            const stored = localStorage.getItem('users_cache');
            const cacheData = stored ? JSON.parse(stored) : null;

            // Check if cache exists and is not expired (1 hour)
            // Do localStorage.removeItem('users_cache'); // TESTING: clear cache
            const now = Date.now();
            const CACHE_DURATION = 1 * 60 * 60 * 1000 //24 * 60 * 60 * 1000; // 24 hours in milliseconds
            
            let users = {};
            if (cacheData && cacheData.expires_at && cacheData.expires_at > now) {
                // Cache is valid
                users = cacheData.users || {};
                const remainingSeconds = Math.floor((cacheData.expires_at - now) / 1000);
                console.log(`Users cache VALID - expires in ${remainingSeconds}s at:`, new Date(cacheData.expires_at).toLocaleTimeString());
            } else {
                // Cache expired or doesn't exist - will be recreated
                if (cacheData?.expires_at) {
                    console.log('Users cache EXPIRED at:', new Date(cacheData.expires_at).toLocaleTimeString(), 'now:', new Date(now).toLocaleTimeString());
                } else {
                    console.log('Users cache MISSING, will rebuild');
                }
            }
            
            if (users[user_id]) {
                setUser(users[user_id]);
                console.log('Loaded user from cache: ', user_id, "=", users[user_id].fullname);
                return;
            }
            
            // User not in cache, fetch from backend
            try {
                const response = await axiosInstance.get(`/users/${user_id}`);
                setUser(response.data);

                // Update cache - re-read to avoid race conditions
                const currentStored = localStorage.getItem('users_cache');
                const currentCache = currentStored ? JSON.parse(currentStored) : null;
                
                // Preserve expiry if cache is still valid, otherwise create new expiry
                let expiresAt = now + CACHE_DURATION;
                // TESTING: Force new expiry every time (comment out to preserve existing expiry)
                /*
                if (currentCache && currentCache.expires_at && currentCache.expires_at > now) {
                    expiresAt = currentCache.expires_at; // Keep existing expiry
                }
                */
                
                const updatedCache = {
                    expires_at: expiresAt,
                    users: {
                        ...(currentCache?.users || {}),
                        [user_id]: response.data
                    }
                };
                
                localStorage.setItem('users_cache', JSON.stringify(updatedCache));
                console.log('Fetched and cached user: ', user_id, "=", response.data.fullname);
            } catch (error) {
                console.error('Error fetching user:', error);
            }
        }

        if (user_id && user_id !== "0") {
            fetchUser();
        }
    }, [user_id]);

    if (!user_id || user_id === "0") {
        return null;
    }

    return (
        <a href={'/users/'+user_id} rel="noopener noreferrer">
            {/* <i className="bi bi-person-circle" title={user ? user.fullname : ''}></i> */}
            {user ? user.fullname : user_id}
        </a>
    );
}

export function GroupLinkView({ group_id, dark }) {
    const [group, setGroup] = useState(null);

    useEffect(() => {
        const fetchGroup = async () => {
            // Read from localStorage every time (avoid stale state)
            const stored = localStorage.getItem('groups_cache');
            const cacheData = stored ? JSON.parse(stored) : null;

            // Check if cache exists and is not expired (1 hour)
            // Do localStorage.removeItem('users_cache'); // TESTING: clear cache
            const now = Date.now();
            const CACHE_DURATION = 1 * 60 * 60 * 1000 //24 * 60 * 60 * 1000; // 24 hours in milliseconds
            
            let groups = {};
            if (cacheData && cacheData.expires_at && cacheData.expires_at > now) {
                // Cache is valid
                groups = cacheData.groups || {};
                const remainingSeconds = Math.floor((cacheData.expires_at - now) / 1000);
                console.log(`Groups cache VALID - expires in ${remainingSeconds}s at:`, new Date(cacheData.expires_at).toLocaleTimeString());
            } else {
                // Cache expired or doesn't exist - will be recreated
                if (cacheData?.expires_at) {
                    console.log('Groups cache EXPIRED at:', new Date(cacheData.expires_at).toLocaleTimeString(), 'now:', new Date(now).toLocaleTimeString());
                } else {
                    console.log('Groups cache MISSING, will rebuild');
                }
            }
            
            if (groups[group_id]) {
                setGroup(groups[group_id]);
                console.log('Loaded group from cache: ', group_id, "=", groups[group_id].name);
                return;
            }
            
            // Group not in cache, fetch from backend
            try {
                const response = await axiosInstance.get(`/groups/${group_id}`);
                setGroup(response.data);

                // Update cache - re-read to avoid race conditions
                const currentStored = localStorage.getItem('groups_cache');
                const currentCache = currentStored ? JSON.parse(currentStored) : null;
                
                // Preserve expiry if cache is still valid, otherwise create new expiry
                let expiresAt = now + CACHE_DURATION;
                // TESTING: Force new expiry every time (comment out to preserve existing expiry)
                /*
                if (currentCache && currentCache.expires_at && currentCache.expires_at > now) {
                    expiresAt = currentCache.expires_at; // Keep existing expiry
                }
                */
                
                const updatedCache = {
                    expires_at: expiresAt,
                    groups: {
                        ...(currentCache?.groups || {}),
                        [group_id]: response.data
                    }
                };
                
                localStorage.setItem('groups_cache', JSON.stringify(updatedCache));
                console.log('Fetched and cached group: ', group_id, "=", response.data.name);
            } catch (error) {
                console.error('Error fetching group:', error);
            }
        }

        if (group_id && group_id !== "0") {
            fetchGroup();
        }
    }, [group_id]);

    if (!group_id || group_id === "0") {
        return null;
    }

    return (
        <a href={'/groups/'+group_id} rel="noopener noreferrer">
            <i className="bi bi-person-circle" title={group ? group.name : ''}></i> {group ? group.name : group_id}
        </a>
    );
}

/* Image Viewer Component

Params:
- id: file ID
- title: alt/title text
- thumbnail: boolean, whether to load thumbnail version
- style: CSS styles for the image
*/
export function ImageView({id, title, thumbnail, style}) {
    const [preview, setPreview] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // console.log('ImageView useEffect:', { id });
        const loadPreview = async () => {
            try {
                // console.log('Loading image preview for:', id);
                const url = thumbnail ? `/files/${id}/download?preview=yes` : `/files/${id}/download`;
                const response = await axiosInstance.get(url, {
                    responseType: 'blob'
                });
                // console.log('Image loaded, blob size:', response.data.size, 'type:', response.data.type);
                // IF an image, create blob URL
                if (response.data.type.startsWith('image/')) {
                    const blobUrl = URL.createObjectURL(response.data);
                    // console.log('Blob URL created:', blobUrl);
                    setPreview(blobUrl);
                } else {
                    setPreview(null);
                }
            } catch (error) {
                console.error('Failed to load image preview:', error);
                setPreview(null);
            }
            finally {
                setLoading(false);
            }
        };
        loadPreview();

        // Cleanup blob URL on unmount
        return () => {
            if (preview && preview.startsWith('blob:')) {
                console.log('Revoking blob URL:', preview);
                URL.revokeObjectURL(preview);
            }
        };
    }, [id]);

    return (
        <>
        {preview && (
            <img 
                src={preview}
                alt={title || 'Preview'}
                title={title || 'Preview'}
                style={style || { maxWidth: '100%', maxHeight: '300px' }}
            />
        )}
        { !preview && loading && (
            // Show a spinner or placeholder while loading
            <div style={{style }}>
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </div>
        )}
        { !preview && !loading && (
            // Show a placeholder if no preview is available
            <i 
                className={`bi bi-${classname2bootstrapIcon('DBFile')}`}
                style={{ ...style }}
            ></i>
            // <div style={{...style, display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: '#f0f0f0', color: '#888' }}>
            //     No Preview Available
            // </div>
        )}
        </>
    );
}

// Component: Link to object
export function ObjectLinkView({ obj_id, dark }) {
    const [myObject, setMyObject] = useState(null);

    useEffect(() => {
        const fetchObject = async () => {
            // Read from localStorage every time (avoid stale state)
            const stored = localStorage.getItem('objects_cache');
            const cacheData = stored ? JSON.parse(stored) : null;

            // Check if cache exists and is not expired (1 minute)
            // Do localStorage.removeItem('objects_cache'); // TESTING: clear cache
            const now = Date.now();
            const CACHE_DURATION = 1 * 60 * 1000 //24 * 60 * 60 * 1000; // 24 hours in milliseconds
            
            let objects = {};
            if (cacheData && cacheData.expires_at && cacheData.expires_at > now) {
                // Cache is valid
                objects = cacheData.objects || {};
                const remainingSeconds = Math.floor((cacheData.expires_at - now) / 1000);
                // console.log(`Objects cache VALID - expires in ${remainingSeconds}s at:`, new Date(cacheData.expires_at).toLocaleTimeString());
            } else {
                // Cache expired or doesn't exist - will be recreated
                if (cacheData?.expires_at) {
                    console.log('Objects cache EXPIRED at:', new Date(cacheData.expires_at).toLocaleTimeString(), 'now:', new Date(now).toLocaleTimeString());
                } else {
                    console.log('Objects cache MISSING, will rebuild');
                }
            }
            
            if (objects[obj_id]) {
                setMyObject(objects[obj_id]);
                // console.log('Loaded object from cache: ', obj_id, "=", objects[obj_id].data.name);
                return;
            }
            
            // Object not in cache, fetch from backend
            try {
                const response = await axiosInstance.get(`/content/${obj_id}`);
                setMyObject(response.data);
                // Update cache - re-read to avoid race conditions
                const currentStored = localStorage.getItem('objects_cache');
                const currentCache = currentStored ? JSON.parse(currentStored) : null;
                
                // Preserve expiry if cache is still valid, otherwise create new expiry
                let expiresAt = now + CACHE_DURATION;
                // TESTING: Force new expiry every time (comment out to preserve existing expiry)
                /*
                if (currentCache && currentCache.expires_at && currentCache.expires_at > now) {
                    expiresAt = currentCache.expires_at; // Keep existing expiry
                }
                */
                
                const updatedCache = {
                    expires_at: expiresAt,
                    objects: {
                        ...(currentCache?.objects || {}),
                        [obj_id]: response.data
                    }
                };
                
                localStorage.setItem('objects_cache', JSON.stringify(updatedCache));
                // console.log('Fetched and cached object: ', obj_id, "=", response.data.data.name);
            } catch (error) {
                console.error('Error fetching object:', error);
            }
        }

        if (obj_id && obj_id !== "0") {
            fetchObject();
        }
    }, [obj_id]);

    if (!obj_id || obj_id === "0") {
        return null;
    }

    return (
        <a href={'/c/'+obj_id} rel="noopener noreferrer">
            <i className={`bi bi-${classname2bootstrapIcon(myObject ? myObject.metadata.classname : '')}`} title={myObject ? myObject.metadata.classname : ''}></i> {myObject ? myObject.data.name : obj_id}
        </a>
    );
}


export function LanguageSelector({ fieldName, value, onChange, dark }) {
    const { t } = useTranslation();

    return (
        <Form.Group className="mb-3">
            <Form.Label>{t('common.language')}</Form.Label>
            <Form.Select
                name={fieldName}
                value={value}
                onChange={onChange}
            >
                <option value="">{t('common.select')}</option>
                <option value="en">🇬🇧 English</option>
                <option value="it">🇮🇹 Italiano</option>
                <option value="de">🇩🇪 Deutsch</option>
                <option value="fr">🇫🇷 Français</option>
            </Form.Select>
        </Form.Group>

    );
}

export function LanguageView({ language, short }) {
    const languagePrefix = language ? language.split('_')[0] : language;
    const languageMap = {
        'en': '🇬🇧 English',
        'it': '🇮🇹 Italiano',
        'de': '🇩🇪 Deutsch',
        'fr': '🇫🇷 Français',
    };
    const languageShortMap = {
        'en': '🇬🇧',
        'it': '🇮🇹',
        'de': '🇩🇪',
        'fr': '🇫🇷',
    };

    return (
        <span>{short ? languageShortMap[languagePrefix] || languagePrefix : languageMap[languagePrefix] || languagePrefix}</span>
    );
}
