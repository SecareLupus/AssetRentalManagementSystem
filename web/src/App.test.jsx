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
import IntelligenceOverview from './pages/IntelligenceOverview';
import AvailabilityHeatmap from './pages/AvailabilityHeatmap';
import ApiInspector from './components/ApiInspector';
import { LayoutDashboard, Box, Calendar, Settings, User, Terminal, Tool, Scan, Brain } from 'lucide-react';
import './App.css';

function App() {
  return <h1>Minimal App for Testing</h1>;
}

export default App;
