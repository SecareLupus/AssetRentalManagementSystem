import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import { DeveloperProvider, useDeveloper } from './context/DeveloperContext';
import Dashboard from './pages/Dashboard';
import Catalog from './pages/Catalog';
import ItemTypeDetails from './pages/ItemTypeDetails';
import AssetDetails from './pages/AssetDetails';
import ReservationWizard from './pages/ReservationWizard';
import ReservationsList from './pages/ReservationsList';
import TechDashboard from './pages/TechDashboard';
import InspectionRunner from './pages/InspectionRunner';
import ProvisioningInterface from './pages/ProvisioningInterface';
import WarehouseKiosk from './pages/WarehouseKiosk';
import ApiInspector from './components/ApiInspector';
import { LayoutDashboard, Box, Calendar, Settings, User, Terminal, Tool, Scan } from 'lucide-react';
import './App.css';

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

const NavLink = ({ to, icon: Icon, label }) => {
  const location = useLocation();
  const active = location.pathname === to;

  return (
    <Link to={to} style={{
      display: 'flex',
      alignItems: 'center',
      gap: '0.75rem',
      padding: '0.75rem 1rem',
      borderRadius: '0.5rem',
      background: active ? 'rgba(99, 102, 241, 0.1)' : 'transparent',
      color: active ? 'var(--primary)' : 'var(--text-muted)',
      textDecoration: 'none',
      fontWeight: active ? 600 : 500,
      transition: 'all 0.2s'
    }}>
      <Icon size={20} /> {label}
    </Link>
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
        backdropFilter: 'blur(10px)',
        position: 'sticky',
        top: 0,
        height: '100vh'
      }}>
        <div style={{ padding: '2rem', display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <div style={{ background: 'var(--primary)', padding: '0.5rem', borderRadius: '0.75rem' }}>
            <Box color="white" size={24} />
          </div>
          <span style={{ fontWeight: 800, fontSize: '1.25rem', letterSpacing: '-0.025em' }}>RMS Fleet</span>
        </div>

        <div style={{ padding: '0 1rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <NavLink to="/" icon={LayoutDashboard} label="Dashboard" />
          <NavLink to="/catalog" icon={Box} label="Equipment Catalog" />
          <NavLink to="/reservations" icon={Calendar} label="Reservations" />

          <div style={{ margin: '1rem 0', padding: '0 1rem', height: '1px', background: 'var(--border)' }} />

          <NavLink to="/tech" icon={Tool} label="Maintenance" />
          <NavLink to="/kiosk" icon={Scan} label="Warehouse Kiosk" />

          <div style={{ margin: '1rem 0', padding: '0 1rem', height: '1px', background: 'var(--border)' }} />

          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', padding: '0.75rem 1rem', borderRadius: '0.5rem', color: 'var(--text-muted)', cursor: 'not-allowed', fontSize: '0.875rem' }}>
            <Settings size={20} /> Settings
          </div>
        </div>

        <DevToggle />
      </nav>

      {/* Main Content */}
      <main style={{ flex: 1, minHeight: '100vh', overflowY: 'auto' }}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/catalog" element={<Catalog />} />
          <Route path="/catalog/:id" element={<ItemTypeDetails />} />
          <Route path="/assets/:id" element={<AssetDetails />} />
          <Route path="/reserve" element={<ReservationWizard />} />
          <Route path="/reservations" element={<ReservationsList />} />
          <Route path="/tech" element={<TechDashboard />} />
          <Route path="/tech/inspect/:id" element={<InspectionRunner />} />
          <Route path="/tech/provision/:id" element={<ProvisioningInterface />} />
          <Route path="/kiosk" element={<WarehouseKiosk />} />
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
