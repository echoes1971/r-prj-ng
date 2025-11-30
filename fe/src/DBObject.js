import React, { useState, useEffect } from 'react';
import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { 
    formateDateTimeString, 
    formatDescription, 
    classname2bootstrapIcon,
    CountryView,
    UserLinkView,
    ObjectLinkView,
    HtmlFieldView
} from './sitenavigation_utils';
import axiosInstance from './axios';

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
                    <HtmlFieldView htmlContent={data.html} dark={dark} />
                )}
            </Card.Body>
            <Card.Footer className={dark ? 'bg-secondary bg-opacity-10' : ''} style={dark ? { borderTop: '1px solid rgba(255,255,255,0.1)' } : {}}>
                <ObjectFooterView data={data} metadata={metadata} objectData={objectData} dark={dark} />
            </Card.Footer>
        </Card>
    );
}
