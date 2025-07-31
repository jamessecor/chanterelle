import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginButton from './components/LoginButton';
import VerificationPage from './components/VerificationPage';
import AdminPage from './components/AdminPage';
import LandingPage from './components/LandingPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/verify" element={<VerificationPage />} />
        <Route path="/admin" element={<AdminPage />} />
        <Route path="*" element={
          <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh'
          }}>
            <Navigate to="/" replace />
          </div>
        } />
      </Routes>
    </Router>
  );
}

export default App;
