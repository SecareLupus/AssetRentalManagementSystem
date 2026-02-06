import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import { DeveloperProvider, useDeveloper } from './context/DeveloperContext';
import { AuthProvider, useAuth } from './context/AuthContext';
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
import IntelligenceOverview from './pages/IntelligenceOverview';
import AvailabilityHeatmap from './pages/AvailabilityHeatmap';
import Login from './pages/Login';
import AdminCenter from './pages/AdminCenter';
import PlanningSimulator from './pages/PlanningSimulator';
import FleetReports from './pages/FleetReports';
import InspectionEditor from './pages/InspectionEditor';
import EntityManager from './pages/EntityManager';
import ApiInspector from './components/ApiInspector';
import { LayoutDashboard, Box, Calendar, Settings, User, Terminal, Wrench, Scan, Brain, LogOut, ChevronLeft, ChevronRight, Menu, ShieldAlert, Calculator, BarChart3, Building2 } from 'lucide-react';
import './App.css';

const DevToggle = ({ collapsed }) => {
  const { isDevMode, setIsDevMode } = useDeveloper();
  return (
    <div style={{
      marginTop: 'auto',
      padding: '1rem',
      borderTop: '1px solid var(--border)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: collapsed ? 'center' : 'space-between',
      fontSize: '0.75rem'
    }}>
      {!collapsed && (
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <Terminal size={14} color={isDevMode ? 'var(--primary)' : 'var(--text-muted)'} />
          <span>Dev Mode</span>
        </div>
      )}
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
        title="Toggle Dev Mode"
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

const NavLink = ({ to, icon: Icon, label, collapsed }) => {
  const location = useLocation();
  const active = location.pathname === to;

  return (
    <Link to={to} title={collapsed ? label : ''} style={{
      display: 'flex',
      alignItems: 'center',
      gap: '0.75rem',
      justifyContent: collapsed ? 'center' : 'flex-start',
      padding: '0.75rem 1rem',
      borderRadius: '0.5rem',
      background: active ? 'rgba(99, 102, 241, 0.1)' : 'transparent',
      color: active ? 'var(--primary)' : 'var(--text-muted)',
      textDecoration: 'none',
      fontWeight: active ? 600 : 500,
      transition: 'all 0.2s',
      whiteSpace: 'nowrap'
    }}>
      <Icon size={20} />
      {!collapsed && <span>{label}</span>}
    </Link>
  );
};

function AppContent() {
  const { isAuthenticated, user, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);

  if (!isAuthenticated) {
    return <Login />;
  }

  return (
    <div className="app-container">
      {/* Sidebar Nav */}
      <nav className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
        <div style={{ padding: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.75rem', justifyContent: collapsed ? 'center' : 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                <div style={{ background: 'var(--primary)', padding: '0.5rem', borderRadius: '0.5rem', minWidth: '36px', display: 'flex', justifyContent: 'center' }}>
                    <Box color="white" size={20} />
                </div>
                {!collapsed && <span style={{ fontWeight: 800, fontSize: '1.25rem', letterSpacing: '-0.025em', whiteSpace: 'nowrap' }}>RMS Fleet</span>}
            </div>
             {!collapsed && (
               <button onClick={() => setCollapsed(true)} className="glass" style={{ padding: '0.25rem', borderRadius: '0.25rem', color: 'var(--text-muted)' }}>
                 <ChevronLeft size={16} />
               </button>
             )}
        </div>

        {collapsed && (
             <div style={{ display: 'flex', justifyContent: 'center', paddingBottom: '1rem' }}>
               <button onClick={() => setCollapsed(false)} className="glass" style={{ padding: '0.25rem', borderRadius: '0.25rem', color: 'var(--text-muted)' }}>
                 <ChevronRight size={16} />
               </button>
             </div>
        )}

        <div style={{ padding: '0 0.5rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <NavLink to="/" icon={LayoutDashboard} label="Dashboard" collapsed={collapsed} />
          <NavLink to="/catalog" icon={Box} label="Equipment Catalog" collapsed={collapsed} />
          <NavLink to="/reservations" icon={Calendar} label="Reservations" collapsed={collapsed} />

          <div style={{ margin: '1rem 0', padding: '0 1rem', height: '1px', background: 'var(--border)' }} />

          <NavLink to="/tech" icon={Wrench} label="Maintenance" collapsed={collapsed} />
          <NavLink to="/kiosk" icon={Scan} label="Warehouse Kiosk" collapsed={collapsed} />
          <NavLink to="/intelligence" icon={Brain} label="Intelligence Hub" collapsed={collapsed} />
          <NavLink to="/entities" icon={Building2} label="Entity Management" collapsed={collapsed} />
          <NavLink to="/simulator" icon={Calculator} label="Fleet Simulator" collapsed={collapsed} />
          <NavLink to="/reports" icon={BarChart3} label="Fleet Reports" collapsed={collapsed} />

          <div style={{ margin: '1rem 0', padding: '0 1rem', height: '1px', background: 'var(--border)' }} />

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
             <NavLink to="/admin" icon={ShieldAlert} label="System Admin" collapsed={collapsed} />
             <NavLink to="#" icon={Settings} label="Settings" collapsed={collapsed} />
            
             <div style={{ padding: '0.75rem 0.5rem', display: 'flex', alignItems: 'center', gap: '0.75rem', borderTop: '1px solid var(--border)', marginTop: '0.5rem', justifyContent: collapsed ? 'center' : 'flex-start' }}>
                <div style={{ minWidth: '32px', width: '32px', height: '32px', borderRadius: '50%', background: 'var(--surface)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--primary)', fontWeight: 700, fontSize: '0.75rem' }}>
                    {user?.username?.charAt(0).toUpperCase() || 'U'}
                </div>
                {!collapsed && (
                    <div style={{ flex: 1, overflow: 'hidden' }}>
                        <div style={{ fontSize: '0.875rem', fontWeight: 600, color: 'var(--text)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{user?.username || 'User'}</div>
                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'capitalize' }}>{user?.role || 'Viewer'}</div>
                    </div>
                )}
                {!collapsed && (
                    <button 
                    onClick={logout}
                    style={{ background: 'transparent', color: 'var(--text-muted)', padding: '0.25rem' }}
                    title="Logout"
                    >
                        <LogOut size={18} />
                    </button>
                )}
            </div>
          </div>
        </div>

        <DevToggle collapsed={collapsed} />
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
          <Route path="/intelligence" element={<IntelligenceOverview />} />
          <Route path="/analytics/heatmap" element={<AvailabilityHeatmap />} />
          <Route path="/admin" element={<AdminCenter />} />
          <Route path="/entities" element={<EntityManager />} />
          <Route path="/admin/inspections/new" element={<InspectionEditor />} />
          <Route path="/admin/inspections/:id" element={<InspectionEditor />} />
          <Route path="/simulator" element={<PlanningSimulator />} />
          <Route path="/reports" element={<FleetReports />} />
        </Routes>
      </main>

      {/* Global API Inspector overlay */}
      <ApiInspector />
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <DeveloperProvider>
        <Router>
          <AppContent />
        </Router>
      </DeveloperProvider>
    </AuthProvider>
  );
}

export default App;
