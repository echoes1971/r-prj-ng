import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import axiosInstance from './axios';

// Format object ID: if 16 chars, format as xxxx-xxxxxxxx-xxxx
export function formatObjectId(objId) {
    if (!objId) return objId;
    if (objId.length === 16) {
        return `${objId.slice(0, 4)}-${objId.slice(4, 12)}-${objId.slice(12, 16)}`;
    }
    return objId;
}
export function classname2bootstrapIcon(classname) {
    switch (classname) {
        case 'DBCompany':
            return 'building';
        case 'DBFile':
            return 'file-earmark-fill';
        case 'DBFolder':
            return 'folder-fill';
        // case 'DBImage':
        //     return 'image-fill';
        case 'DBLink':
            return 'link-45deg';
        case 'DBNews':
            return 'newspaper-fill';
        case 'DBNote':
            return 'file-text-fill';
        case 'DBObject':
            return 'box-fill';
        case 'DBPage':
            return 'file-richtext-fill';
        case 'DBPerson':
            return 'person-fill';
        default:
            return 'question-circle-fill';
    }
}
export function formatDescription(description) {
    if (!description) return '';
    // replace \n with <br/>
    return description.replace(/\n/g, '<br/>');
}
export function formateDateTimeString(dateTimeString) {
    if (!dateTimeString) return '';
    const date = new Date(dateTimeString);
    return date.toLocaleString();
}

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
            try {
                const response = await axiosInstance.get(`/content/country/${country_id}`);
                setCountry(response.data);
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

// Component: Link to user profile
export function UserLinkView({ user_id, dark }) {
    const [user, setUser] = useState(null);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axiosInstance.get(`/users/${user_id}`);
                setUser(response.data);
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
            <i className="bi bi-person-circle" title={user ? user.fullname : ''}></i> {user ? user.fullname : user_id}
        </a>
    );
}

export function GroupLinkView({ group_id, dark }) {
    const [group, setGroup] = useState(null);

    useEffect(() => {
        const fetchGroup = async () => {
            try {
                const response = await axiosInstance.get(`/groups/${group_id}`);
                setGroup(response.data);
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

// Component: Link to object
export function ObjectLinkView({ obj_id, dark }) {
    const [myObject, setMyObject] = useState(null);

    useEffect(() => {
        const fetchObject = async () => {
            try {
                const response = await axiosInstance.get(`/content/${obj_id}`);
                setMyObject(response.data);
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

// Component: Render HTML content safely
export function HtmlFieldView({ htmlContent, dark }) {
    return (
        <div dangerouslySetInnerHTML={{ __html: htmlContent }}></div>
    );
}
