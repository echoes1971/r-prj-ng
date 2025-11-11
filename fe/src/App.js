import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import Login from "./Login";
import Users from "./Users";

function App() {
  const token = localStorage.getItem("token");
  return (
    <Router>
      <AppNavbar />
      <Routes>
        <Route path="/" element={<DefaultPage />} />
        <Route path="/login" element={<Login />} />

        {/* Protected route */}
        <Route
          path="/users"
          element={token ? <Users /> : <Navigate to="/" />}
        />

        {/* Default -> redirect to / */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;
