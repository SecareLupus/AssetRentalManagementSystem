import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { DeveloperProvider, useDeveloper } from './context/DeveloperContext';
import Dashboard from './pages/Dashboard';
import ApiInspector from './components/ApiInspector';
import { LayoutDashboard, Box, Settings, User, Terminal } from 'lucide-react';
import './App.css';

// A small sub-component to handle the Dev Mode toggle in the sidebar
const DevToggle = () => {
  const { isDevMode, setIsDevMode } = useDeveloper();
  return (
    <div style={{
      marginTop: 'auto',
      padding: '1rem',
      borderTop: '1px solid var(--border)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      fontSize: '0.75rem'
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
        <Terminal size={14} color={isDevMode ? 'var(--primary)' : 'var(--text-muted)'} />
        <span>Dev Mode</span>
      </div>
      <button
        onClick={() => setIsDevMode(!isDevMode)}
        style={{
          width: '32px',
          height: '18px',
          borderRadius: '9px',
          background: isDevMode ? 'var(--primary)' : 'var(--border)',
          position: 'relative',
          transition: 'background 0.3s'
        }}
      >
        <div style={{
          width: '14px',
          height: '14px',
          borderRadius: '50%',
          background: 'white',
          position: 'absolute',
          top: '2px',
          left: isDevMode ? '16px' : '2px',
          transition: 'left 0.3s'
        }} />
      </button>
    </div>
  );
};

function AppContent() {
  return (
    <div className="app-container" style={{ display: 'flex', minHeight: '100vh', background: 'var(--background)', color: 'var(--text)' }}>
      {/* Sidebar Nav */}
      <nav style={{
        width: '260px',
        borderRight: '1px solid var(--border)',
        display: 'flex',
        flexDirection: 'column',
        background: 'rgba(15, 23, 42, 0.5)',
        backdropFilter: 'blur(10px)'
      }}>
        <div style={{ padding: '2rem', display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <div style={{ background: 'var(--primary)', padding: '0.5rem', borderRadius: '0.75rem' }}>
            <Box color="white" size={24} />
          </div>
          <span style={{ fontWeight: 800, fontSize: '1.25rem', letterSpacing: '-0.025em' }}>RMS Fleet</span>
        </div>

        <div style={{ padding: '0 1rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <Link to="/" style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.75rem',
            padding: '0.75rem 1rem',
            borderRadius: '0.5rem',
            background: 'rgba(99, 102, 241, 0.1)',
            color: 'var(--primary)',
            textDecoration: 'none',
            fontWeight: 600
          }}>
            <LayoutDashboard size={20} /> Dashboard
          </Link>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.75rem',
            padding: '0.75rem 1rem',
            borderRadius: '0.5rem',
            color: 'var(--text-muted)',
            cursor: 'not-allowed'
          }}>
            <Box size={20} /> Inventory
          </div>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.75rem',
            padding: '0.75rem 1rem',
            borderRadius: '0.5rem',
            color: 'var(--text-muted)',
            cursor: 'not-allowed'
          }}>
            <User size={20} /> Customers
          </div>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.75rem',
            padding: '0.75rem 1rem',
            borderRadius: '0.5rem',
            color: 'var(--text-muted)',
            cursor: 'not-allowed'
          }}>
            <Settings size={20} /> Settings
          </div>
        </div>

        <DevToggle />
      </nav>

      {/* Main Content */}
      <main style={{ flex: 1, height: '100vh', overflowY: 'auto' }}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
        </Routes>
      </main>

      {/* Global API Inspector overlay */}
      <ApiInspector />
    </div>
  );
}

function App() {
  return (
    <DeveloperProvider>
      <Router>
        <AppContent />
      </Router>
    </DeveloperProvider>
  );
}

export default App;
