import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { ShieldAlert, RefreshCw, ClipboardCheck, AlertTriangle, CheckCircle2, Package, Search, ListFilter, Trash2, Users, Settings, Mail, Globe, Save, Lock, ToggleLeft, ToggleRight } from 'lucide-react';
import { GlassCard } from '../components/Shared';

const AdminCenter = () => {
    const [activeTab, setActiveTab] = useState('recall'); // 'recall' | 'recon'
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState(null);

    // Bulk Recall State
    const [itemsToRecall, setItemsToRecall] = useState([]);
    const [selectedIds, setSelectedIds] = useState([]);
    const [recallFilter, setRecallFilter] = useState('');

    // Recon State
    const [reconLocation, setReconLocation] = useState('Central Ops Warehouse');
    const [scannedTags, setScannedTags] = useState('');
    const [reconReport, setReconReport] = useState(null);

    // Inspection Templates State
    const [templates, setTemplates] = useState([]);

    // User Management State
    const [users, setUsers] = useState([]);

    // Settings State
    const [settings, setSettings] = useState({
        company_identity: { name: '', logo_url: '', support_email: '' },
        logistics_policies: { default_return_window_days: 14, late_fee_per_day: 0, currency: 'USD' },
        feature_flags: { enable_auto_alerts: true, enable_ai_forecasting: false }
    });

    useEffect(() => {
        if (activeTab === 'recall') {
            fetchRecallableItems();
        } else if (activeTab === 'inspections') {
            fetchTemplates();
        } else if (activeTab === 'users') {
            fetchUsers();
        } else if (activeTab === 'settings') {
            fetchSettings();
        }
    }, [activeTab]);

    const fetchTemplates = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/v1/catalog/inspection-templates');
            setTemplates(res.data || []);
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to fetch inspection templates.' });
        } finally {
            setLoading(false);
        }
    };

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/v1/admin/users');
            setUsers(res.data || []);
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to fetch users.' });
        } finally {
            setLoading(false);
        }
    };

    const fetchSettings = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/v1/admin/settings');
            // Merge into defaults
            setSettings(prev => ({
                ...prev,
                ...res.data
            }));
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to fetch system settings.' });
        } finally {
            setLoading(false);
        }
    };

    const handleUpdateUser = async (user) => {
        try {
            await axios.put(`/v1/admin/users/${user.id}`, user);
            setMessage({ type: 'success', text: `User ${user.username} updated.` });
            fetchUsers();
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to update user.' });
        }
    };

    const handleDeleteUser = async (id) => {
        if (!window.confirm("Are you sure you want to delete this user?")) return;
        try {
            await axios.delete(`/v1/admin/users/${id}`);
            setMessage({ type: 'success', text: 'User deleted.' });
            fetchUsers();
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to delete user.' });
        }
    };

    const handleSaveSetting = async (key, value) => {
        try {
            await axios.put('/v1/admin/settings', { key, value });
            setMessage({ type: 'success', text: `${key.replace('_', ' ')} saved.` });
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to save setting.' });
        }
    };

    const fetchRecallableItems = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/v1/inventory/assets');
            // Show only deployed/available assets for recall
            setItemsToRecall(res.data.filter(a => a.status === 'deployed' || a.status === 'available'));
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to fetch inventory.' });
        } finally {
            setLoading(false);
        }
    };

    const handleBulkRecall = async () => {
        if (selectedIds.length === 0) return;
        setLoading(true);
        try {
            await axios.post('/v1/inventory/assets/bulk-recall', { asset_ids: selectedIds });
            setMessage({ type: 'success', text: `Successfully recalled ${selectedIds.length} assets.` });
            setSelectedIds([]);
            fetchRecallableItems();
        } catch (err) {
            setMessage({ type: 'error', text: 'Bulk recall failed.' });
        } finally {
            setLoading(false);
        }
    };

    const handleRecon = async () => {
        if (!scannedTags.trim()) return;
        setLoading(true);
        try {
            const tags = scannedTags.split('\n').map(t => t.trim()).filter(t => t);
            const res = await axios.post('/v1/inventory/reconcile', {
                location: reconLocation,
                scanned_tags: tags
            });
            setReconReport(res.data);
            setMessage({ type: 'success', text: 'Reconciliation report generated.' });
        } catch (err) {
            setMessage({ type: 'error', text: 'Reconciliation failed.' });
        } finally {
            setLoading(false);
        }
    };

    const toggleId = (id) => {
        setSelectedIds(prev => prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]);
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <header style={{ marginBottom: '3rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.75rem' }}>
                    <div style={{ background: 'var(--error)20', padding: '0.75rem', borderRadius: '0.75rem' }}>
                        <ShieldAlert size={24} color="var(--error)" />
                    </div>
                    <h1 style={{ fontSize: '1.75rem', fontWeight: 800 }}>System Administration</h1>
                </div>
                <p style={{ color: 'var(--text-muted)' }}>Management of high-level fleet operations and data integrity.</p>
            </header>

            <div style={{ display: 'flex', gap: '1rem', marginBottom: '2rem', borderBottom: '1px solid var(--border)' }}>
                <TabButton active={activeTab === 'recall'} onClick={() => setActiveTab('recall')} icon={RefreshCw} label="Bulk Recall" />
                <TabButton active={activeTab === 'recon'} onClick={() => setActiveTab('recon')} icon={ClipboardCheck} label="Inventory Recon" />
                <TabButton active={activeTab === 'inspections'} onClick={() => setActiveTab('inspections')} icon={ClipboardCheck} label="Inspection Templates" />
                <TabButton active={activeTab === 'users'} onClick={() => setActiveTab('users')} icon={Users} label="Users" />
                <TabButton active={activeTab === 'settings'} onClick={() => setActiveTab('settings')} icon={Settings} label="Global Settings" />
            </div>

            {message && (
                <div className="glass" style={{ padding: '1rem', borderRadius: '0.75rem', marginBottom: '2rem', display: 'flex', gap: '0.5rem', alignItems: 'center', color: message.type === 'error' ? 'var(--error)' : 'var(--success)', background: message.type === 'error' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(16, 185, 129, 0.1)' }}>
                    {message.type === 'error' ? <AlertTriangle size={18} /> : <CheckCircle2 size={18} />}
                    {message.text}
                </div>
            )}

            {activeTab === 'recall' && (
                <div className="animate-in fade-in slide-in-from-bottom-4">
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.5rem', alignItems: 'center' }}>
                        <div style={{ position: 'relative', flex: 1, maxWidth: '400px' }}>
                            <input
                                type="text"
                                className="glass"
                                placeholder="Filter items to recall..."
                                style={{ width: '100%', padding: '0.625rem 1rem 0.625rem 2.5rem', borderRadius: '0.5rem', color: 'white' }}
                                value={recallFilter}
                                onChange={(e) => setRecallFilter(e.target.value)}
                            />
                            <Search size={16} style={{ position: 'absolute', left: '0.75rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.4 }} />
                        </div>
                        <button
                            className="btn-primary"
                            style={{ background: 'var(--error)' }}
                            disabled={selectedIds.length === 0 || loading}
                            onClick={handleBulkRecall}
                        >
                            Perform Recall ({selectedIds.length})
                        </button>
                    </div>

                    <div className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                            <thead>
                                <tr style={{ textAlign: 'left', background: 'rgba(255,255,255,0.03)', color: 'var(--text-muted)' }}>
                                    <th style={{ padding: '1rem' }}><input type="checkbox" onChange={(e) => setSelectedIds(e.target.checked ? itemsToRecall.map(i => i.id) : [])} checked={selectedIds.length === itemsToRecall.length && itemsToRecall.length > 0} /></th>
                                    <th style={{ padding: '1rem' }}>Asset Tag</th>
                                    <th style={{ padding: '1rem' }}>Serial</th>
                                    <th style={{ padding: '1rem' }}>Current Status</th>
                                    <th style={{ padding: '1rem' }}>Location</th>
                                </tr>
                            </thead>
                            <tbody>
                                {itemsToRecall.filter(i => (i.asset_tag || '').toLowerCase().includes(recallFilter.toLowerCase()) || (i.serial_number || '').toLowerCase().includes(recallFilter.toLowerCase())).map(item => (
                                    <tr key={item.id} style={{ borderTop: '1px solid var(--border)', background: selectedIds.includes(item.id) ? 'rgba(99, 102, 241, 0.05)' : 'transparent' }}>
                                        <td style={{ padding: '1rem' }}>
                                            <input type="checkbox" checked={selectedIds.includes(item.id)} onChange={() => toggleId(item.id)} />
                                        </td>
                                        <td style={{ padding: '1rem', fontWeight: 600 }}>{item.asset_tag}</td>
                                        <td style={{ padding: '1rem', fontFamily: 'monospace' }}>{item.serial_number}</td>
                                        <td style={{ padding: '1rem' }}>
                                            <span style={{ fontSize: '0.75rem', padding: '0.125rem 0.5rem', borderRadius: '1rem', background: item.status === 'available' ? 'var(--success)20' : 'var(--primary)20', color: item.status === 'available' ? 'var(--success)' : 'var(--primary)' }}>
                                                {item.status.toUpperCase()}
                                            </span>
                                        </td>
                                        <td style={{ padding: '1rem', color: 'var(--text-muted)' }}>{item.location || 'Unknown'}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {activeTab === 'inspections' && (
                <div className="animate-in fade-in slide-in-from-bottom-4">
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1.5rem', alignItems: 'center' }}>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>Active Templates</h2>
                        <a href="/admin/inspections/new" className="btn-primary" style={{ textDecoration: 'none' }}>
                            Create Template
                        </a>
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                        {templates.map(t => (
                            <GlassCard key={t.id} style={{ padding: '1.5rem' }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem' }}>
                                    <h4 style={{ fontWeight: 800 }}>{t.name}</h4>
                                    <button
                                        onClick={async () => {
                                            if (window.confirm("Delete this template?")) {
                                                await axios.delete(`/v1/catalog/inspection-templates/${t.id}`);
                                                fetchTemplates();
                                            }
                                        }}
                                        style={{ background: 'transparent', color: 'var(--text-muted)' }}
                                    >
                                        <Trash2 size={16} />
                                    </button>
                                </div>
                                <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)', marginBottom: '1.5rem', minHeight: '3rem' }}>{t.description || 'No description provided.'}</p>
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Created: {new Date(t.created_at).toLocaleDateString()}</span>
                                    <a href={`/admin/inspections/${t.id}`} style={{ fontSize: '0.875rem', fontWeight: 600, color: 'var(--primary)', textDecoration: 'none' }}>Edit Details</a>
                                </div>
                            </GlassCard>
                        ))}
                    </div>
                    {templates.length === 0 && !loading && (
                        <div style={{ textAlign: 'center', padding: '4rem', color: 'var(--text-muted)' }}>
                            No templates found. Click "Create Template" to get started.
                        </div>
                    )}
                </div>
            )}

            {activeTab === 'users' && (
                <div className="animate-in fade-in slide-in-from-bottom-4">
                    <div className="flex-between" style={{ marginBottom: '1.5rem' }}>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>User Accounts</h2>
                        <button className="btn-primary" onClick={() => alert("Invite system not implemented. Use /register for now.")}>
                            Manage Invites
                        </button>
                    </div>

                    <div className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                            <thead>
                                <tr style={{ textAlign: 'left', background: 'rgba(255,255,255,0.03)', color: 'var(--text-muted)' }}>
                                    <th style={{ padding: '1rem' }}>Username</th>
                                    <th style={{ padding: '1rem' }}>Role</th>
                                    <th style={{ padding: '1rem' }}>Status</th>
                                    <th style={{ padding: '1rem' }}>Last Login</th>
                                    <th style={{ padding: '1rem' }}>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map(u => (
                                    <tr key={u.id} style={{ borderTop: '1px solid var(--border)' }}>
                                        <td style={{ padding: '1rem' }}>
                                            <div style={{ fontWeight: 600 }}>{u.username}</div>
                                            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{u.email}</div>
                                        </td>
                                        <td style={{ padding: '1rem' }}>
                                            <select
                                                className="glass"
                                                style={{ padding: '0.25rem 0.5rem', borderRadius: '0.4rem', color: 'white', background: 'var(--surface)' }}
                                                value={u.role}
                                                onChange={(e) => handleUpdateUser({ ...u, role: e.target.value })}
                                            >
                                                <option value="admin">Admin</option>
                                                <option value="fleet_manager">Fleet Manager</option>
                                                <option value="technician">Technician</option>
                                                <option value="viewer">Viewer</option>
                                            </select>
                                        </td>
                                        <td style={{ padding: '1rem' }}>
                                            <button
                                                onClick={() => handleUpdateUser({ ...u, is_enabled: !u.is_enabled })}
                                                style={{ background: 'transparent', display: 'flex', alignItems: 'center', gap: '0.5rem', color: u.is_enabled ? 'var(--success)' : 'var(--error)' }}
                                            >
                                                {u.is_enabled ? <ToggleRight size={20} /> : <ToggleLeft size={20} />}
                                                {u.is_enabled ? 'Enabled' : 'Disabled'}
                                            </button>
                                        </td>
                                        <td style={{ padding: '1rem', color: 'var(--text-muted)' }}>
                                            {u.last_login_at ? new Date(u.last_login_at).toLocaleString() : 'Never'}
                                        </td>
                                        <td style={{ padding: '1rem' }}>
                                            <button onClick={() => handleDeleteUser(u.id)} className="glass" style={{ padding: '0.4rem', borderRadius: '0.4rem', color: 'var(--error)' }}>
                                                <Trash2 size={16} />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {activeTab === 'settings' && (
                <div className="animate-in fade-in slide-in-from-bottom-4" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '2rem' }}>
                    <GlassCard>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Globe size={18} color="var(--primary)" /> Company Identity
                        </h3>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            <div>
                                <label className="form-label">Legal Name</label>
                                <input
                                    className="glass" style={{ width: '100%', padding: '0.75rem', color: 'white' }}
                                    value={settings.company_identity?.name || ''}
                                    onChange={e => setSettings({ ...settings, company_identity: { ...settings.company_identity, name: e.target.value } })}
                                />
                            </div>
                            <div>
                                <label className="form-label">Support Email</label>
                                <input
                                    className="glass" style={{ width: '100%', padding: '0.75rem', color: 'white' }}
                                    value={settings.company_identity?.support_email || ''}
                                    onChange={e => setSettings({ ...settings, company_identity: { ...settings.company_identity, support_email: e.target.value } })}
                                />
                            </div>
                            <button className="btn-primary" onClick={() => handleSaveSetting('company_identity', settings.company_identity)}>
                                Save Identity
                            </button>
                        </div>
                    </GlassCard>

                    <GlassCard>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Lock size={18} color="var(--primary)" /> Logistics Policies
                        </h3>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            <div>
                                <label className="form-label">Default Return Window (Days)</label>
                                <input
                                    type="number" className="glass" style={{ width: '100%', padding: '0.75rem', color: 'white' }}
                                    value={settings.logistics_policies?.default_return_window_days || 0}
                                    onChange={e => setSettings({ ...settings, logistics_policies: { ...settings.logistics_policies, default_return_window_days: parseInt(e.target.value) } })}
                                />
                            </div>
                            <div>
                                <label className="form-label">Currency Code</label>
                                <input
                                    className="glass" style={{ width: '100%', padding: '0.75rem', color: 'white' }}
                                    value={settings.logistics_policies?.currency || 'USD'}
                                    onChange={e => setSettings({ ...settings, logistics_policies: { ...settings.logistics_policies, currency: e.target.value } })}
                                />
                            </div>
                            <button className="btn-primary" onClick={() => handleSaveSetting('logistics_policies', settings.logistics_policies)}>
                                Save Policies
                            </button>
                        </div>
                    </GlassCard>

                    <GlassCard>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <ShieldAlert size={18} color="var(--primary)" /> Feature Flags
                        </h3>
                        <div style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: '1.5rem' }}>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
                                <input
                                    type="checkbox"
                                    checked={settings.feature_flags?.enable_auto_alerts}
                                    onChange={e => {
                                        const newVal = { ...settings.feature_flags, enable_auto_alerts: e.target.checked };
                                        setSettings({ ...settings, feature_flags: newVal });
                                        handleSaveSetting('feature_flags', newVal);
                                    }}
                                />
                                Auto Alerts
                            </label>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
                                <input
                                    type="checkbox"
                                    checked={settings.feature_flags?.enable_ai_forecasting}
                                    onChange={e => {
                                        const newVal = { ...settings.feature_flags, enable_ai_forecasting: e.target.checked };
                                        setSettings({ ...settings, feature_flags: newVal });
                                        handleSaveSetting('feature_flags', newVal);
                                    }}
                                />
                                AI Forecasting
                            </label>
                        </div>
                    </GlassCard>
                </div>
            )}

            {activeTab === 'recon' && (
                <div className="animate-in fade-in slide-in-from-bottom-4" style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem' }}>
                    <div>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Location to Reconcile</label>
                            <input
                                type="text"
                                className="glass"
                                value={reconLocation}
                                onChange={(e) => setReconLocation(e.target.value)}
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                            />
                        </div>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Scanned Asset Tags (one per line)</label>
                            <textarea
                                className="glass"
                                rows="12"
                                placeholder={"Tag 001\nTag 002..."}
                                value={scannedTags}
                                onChange={(e) => setScannedTags(e.target.value)}
                                style={{ width: '100%', padding: '1rem', borderRadius: '1rem', color: 'white', fontFamily: 'monospace', resize: 'none' }}
                            />
                        </div>
                        <button className="btn-primary" style={{ width: '100%' }} onClick={handleRecon} disabled={loading}>
                            {loading ? 'Analyzing Fleet State...' : 'Start Reconciliation Analysis'}
                        </button>
                    </div>

                    <div>
                        <div className="glass" style={{ padding: '1.5rem', borderRadius: '1.5rem', minHeight: '400px' }}>
                            <h3 style={{ fontWeight: 600, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                <ListFilter size={18} color="var(--primary)" /> Reconciliation Report
                            </h3>

                            {!reconReport ? (
                                <div style={{ height: '300px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontSize: '0.875rem' }}>
                                    Scan tags and run analysis to see discrepancies.
                                </div>
                            ) : (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                                    <ReportSection title="Verified" items={reconReport.verified_tags} color="var(--success)" icon={CheckCircle2} />
                                    <ReportSection title="Missing Tags" items={reconReport.missing_tags} color="var(--error)" icon={AlertTriangle} desc="Expected at location but not scanned." />
                                    <ReportSection title="Unexpected Scans" items={reconReport.unexpected_tags} color="var(--warning)" icon={Package} desc="Scanned but not found at this location." />
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

const TabButton = ({ active, onClick, icon: Icon, label }) => (
    <button
        onClick={onClick}
        style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem',
            padding: '1rem',
            background: 'transparent',
            color: active ? 'var(--primary)' : 'var(--text-muted)',
            borderBottom: active ? '2px solid var(--primary)' : '2px solid transparent',
            fontWeight: active ? 700 : 500,
            transition: 'all 0.2s'
        }}
    >
        <Icon size={18} /> {label}
    </button>
);

const ReportSection = ({ title, items, color, icon: Icon, desc }) => (
    <div>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color, fontWeight: 700, fontSize: '0.875rem' }}>
                <Icon size={16} /> {title}
            </div>
            <span style={{ fontSize: '0.75rem', background: 'rgba(255,255,255,0.05)', padding: '0.125rem 0.375rem', borderRadius: '0.25rem' }}>{items?.length || 0}</span>
        </div>
        {desc && <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>{desc}</p>}
        {items && items.length > 0 ? (
            <div style={{ maxHeight: '100px', overflowY: 'auto', display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                {items.map(tag => (
                    <span key={tag} style={{ fontSize: '0.75rem', background: 'var(--surface)', padding: '0.25rem 0.5rem', borderRadius: '0.375rem', border: `1px solid ${color}40` }}>{tag}</span>
                ))}
            </div>
        ) : (
            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', opacity: 0.5 }}>None</div>
        )}
    </div>
);

export default AdminCenter;
