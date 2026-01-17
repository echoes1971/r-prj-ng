import React, { useParams } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import { AdminDashboard } from './admin/Dashboard';
import Login from "./Login";
import Users from "./Users";
import UserProfile from "./UserProfile";
import Groups from './Groups';
import GroupProfile from './GroupProfile';
import { Companies } from './dbobjects/DBCompany';
import { Events } from './dbobjects/DBEvent';
import { FileDownload, Files } from "./dbobjects/DBFile";
import { Folders } from './dbobjects/DBFolders';
import { Links } from './dbobjects/DBLink';
import { News } from './dbobjects/DBNews';
import { Notes } from './dbobjects/DBNote';
import { Objects } from "./dbobjects/DBObject";
import { Pages } from "./dbobjects/DBPage";
import { People } from './dbobjects/DBPeople';
import SiteNavigation from './SiteNavigation';
import ContentEdit from './ContentEdit';
import Search from './Search';
import { AppFooter } from './Footer';
import { isTokenValid } from './sitenavigation_utils';

function App() {
  // const token = localStorage.getItem("token");
  const validToken = isTokenValid();
  // console.log("App: token present: " + (token ? "yes" : "no") + ", valid: " + (validToken ? "yes" : "no"));
  const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
  const isAdmin = groups.includes("-2");
  const isWebmaster = groups.includes(app_cfg.webmaster_group_id);
  
  return (
    <div className="d-flex flex-column min-vh-100">
      <main className="flex-fill">
        <Router>
          <AppNavbar />
          <Routes>
            <Route path="/" element={<SiteNavigation />} />
            <Route path="/default" element={<DefaultPage />} />
            <Route path="/login" element={<Login />} />

            {/* Site Navigation - content by object ID */}
            <Route path="/c/:objectId" element={<SiteNavigation />} />
            <Route path="/c" element={<SiteNavigation />} />

            {/* Global Search */}
            <Route path="/search" element={<Search />} />

            {/* File download by object ID */}
            <Route path="/f/:objectId/download" element={<FileDownload />} />

            {/* Content Edit - edit object by ID (requires authentication) */}
            <Route path="/e/:id" element={validToken ? <ContentEdit /> : <Navigate to={`/c/${window.location.pathname.split('/').pop()}`} replace />} />

            {/* User profile - accessible by the user themselves or admins */}
            <Route path="/users/:userId" element={validToken ? <UserProfile /> : <Navigate to="/login" />} />

            {/* Group profile - only for admins */}
            <Route path="/groups/:groupId" element={validToken ? <GroupProfile /> : <Navigate to="/" />} />

            {/* **** Webmaster **** */}
            <Route path="/folders" element={validToken && isWebmaster ? <Folders /> : <Navigate to="/" />} />
            <Route path="/pages"   element={validToken && isWebmaster ?   <Pages /> : <Navigate to="/" />} />
            <Route path="/news"    element={validToken && isWebmaster ?    <News /> : <Navigate to="/" />} />
            <Route path="/files"   element={validToken && isWebmaster ?   <Files /> : <Navigate to="/" />} />
            <Route path="/links"   element={validToken && isWebmaster ?   <Links /> : <Navigate to="/" />} />
            <Route path="/events"  element={validToken && isWebmaster ?  <Events /> : <Navigate to="/" />} />

            {/* **** Contacts **** */}
            <Route path="/companies" element={validToken ? <Companies /> : <Navigate to="/" />} />
            <Route path="/people"    element={validToken ?    <People /> : <Navigate to="/" />} />
            {/* **** User menu **** */}
            <Route path="/notes"    element={validToken ?    <Notes /> : <Navigate to="/" />} />

            {/* **** Admin **** */}

            {/* Protected routes - only for admins (group -2) */}
            <Route path="/admin/dashboard" element={validToken && isAdmin ? <AdminDashboard /> : <Navigate to="/" />} />

            <Route path="/users" element={validToken && isAdmin ? <Users /> : <Navigate to="/" />} />
            <Route path="/groups" element={validToken && isAdmin ? <Groups /> : <Navigate to="/" />} />

            <Route path="/objects"    element={validToken ?    <Objects /> : <Navigate to="/" />} />

            {/* Default -> redirect to / */}
            <Route path="*" element={<Navigate to="/default" />} />
          </Routes>
        </Router>
      </main>
      <AppFooter />
    </div>
  );
}

export default App;
