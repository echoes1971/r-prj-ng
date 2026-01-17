import React, { useState, useEffect } from 'react';
import { Spinner } from 'react-bootstrap';
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
        case 'DBEvent':
            return 'calendar-event';
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

        case 'DBUser':
            return 'person-fill';
        case 'DBGroup':
            return 'people-fill';
        default:
            return 'question-circle-fill';
    }
}
export function formatDescription(description) {
    if (!description) return '';
    // replace \n with <br/>

    // escape HTML special characters
    const escapeHtml = (text) => {
        return text.replace(/&/g, "&amp;")
                   .replace(/</g, "&lt;")
                   .replace(/>/g, "&gt;")
                   .replace(/"/g, "&quot;")
                   .replace(/'/g, "&#039;");
    };

    return escapeHtml(description).replace(/\n/g, '<br/>');
}
export function formateDateTimeString(dateTimeString) {
    if (!dateTimeString) return '';
    const date = new Date(dateTimeString);
    return date.toLocaleString();
}


export function isTokenValid() {
    const token = localStorage.getItem("token");
    if (!token) return false;
    const expiresAt = localStorage.getItem("expires_at");
    if (!expiresAt) return false;
    const nowInSeconds = Math.floor(Date.now() / 1000);
    console.log(`Token expiry check: now=${nowInSeconds}, expires_at=${expiresAt}`);
    if (nowInSeconds >= Number(expiresAt)) return false;

    // Optionally, you can implement further validation, e.g., check expiration
    return true;
}

export function isAdminUser() {
    const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
    return groups.includes("-2");
}

export function isWebmasterUser() {
    const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
    const webmasterGroupId = process.env.REACT_APP_WEBMASTER_GROUP_ID || "-3";
    return groups.includes(webmasterGroupId);
}

export function isGuestUser() {
    const token = localStorage.getItem("token");
    const groupIDs = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
    return !token || token === "" || (groupIDs.length <= 2 && groupIDs.includes("-4"));
}
