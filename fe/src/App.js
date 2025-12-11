import React, { useParams } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import Login from "./Login";
import Users from "./Users";
import UserProfile from "./UserProfile";
import Groups from './Groups';
import GroupProfile from './GroupProfile';
import { Companies } from './Companies';
import { Files } from './Files';
import { Folders } from './Folders';
import { News } from './News';
import { Pages } from './Pages';
import { People } from './People';
import SiteNavigation from './SiteNavigation';
import ContentEdit from './ContentEdit';
import { FileDownload } from './DBFile';
import Search from './Search';
import { AppFooter } from './Footer';

function App() {
  const token = localStorage.getItem("token");
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
            <Route path="/e/:id" element={token ? <ContentEdit /> : <Navigate to={`/c/${window.location.pathname.split('/').pop()}`} replace />} />

            {/* User profile - accessible by the user themselves or admins */}
            <Route
              path="/users/:userId"
              element={token ? <UserProfile /> : <Navigate to="/login" />}
            />

            {/* Group profile - only for admins */}
            <Route
              path="/groups/:groupId"
              element={token ? <GroupProfile /> : <Navigate to="/" />}
            />

            {/* **** Webmaster **** */}
            <Route path="/folders" element={token && isWebmaster ? <Folders /> : <Navigate to="/" />} />
            <Route path="/pages"   element={token && isWebmaster ?   <Pages /> : <Navigate to="/" />} />
            <Route path="/news"    element={token && isWebmaster ?    <News /> : <Navigate to="/" />} />
            <Route path="/files"   element={token && isWebmaster ?   <Files /> : <Navigate to="/" />} />

            {/* **** Contacts **** */}
            <Route path="/companies" element={token ? <Companies /> : <Navigate to="/" />} />
            <Route path="/people"    element={token ?    <People /> : <Navigate to="/" />} />

            {/* **** Admin **** */}

            {/* Protected routes - only for admins (group -2) */}
            <Route
              path="/users"
              element={token && isAdmin ? <Users /> : <Navigate to="/" />}
            />

            <Route
              path="/groups"
              element={token && isAdmin ? <Groups /> : <Navigate to="/" />}
            />

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
