import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { LayoutDashboard, Box, Calendar, AlertCircle, ChevronRight, Activity } from 'lucide-react';

const Dashboard = () => {
    const [stats, setStats] = useState({
        totalItems: 0,
        availableAssets: 0,
        pendingReservations: 0,
        activeAlerts: 4
    });
    const [itemTypes, setItemTypes] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get('/v1/catalog/item-types');
                setItemTypes(response.data || []);
                setStats(prev => ({ ...prev, totalItems: response.data?.length || 0 }));
            } catch (error) {
                console.error("Failed to fetch dashboard data", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2.5rem' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '0.5rem' }}>Commander's Dashboard</h1>
                <p style={{ color: 'var(--text-muted)' }}>Real-time overview of fleet operations and inventory.</p>
            </header>

            {/* Stats Grid */}
            <div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
                gap: '1.5rem',
                marginBottom: '3rem'
            }}>
                {[
                    { label: 'Total Item Types', value: stats.totalItems, icon: Box, color: 'var(--primary)' },
                    { label: 'Available Assets', value: '42', icon: Activity, color: 'var(--success)' },
                    { label: 'Pending Requests', value: '12', icon: Calendar, color: 'var(--warning)' },
                    { label: 'System Alerts', value: '0', icon: AlertCircle, color: 'var(--error)' },
                ].map((stat, i) => (
                    <div key={i} className="glass" style={{ padding: '1.5rem', borderRadius: '1rem' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem' }}>
                            <stat.icon size={24} color={stat.color} />
                            <div style={{ background: `${stat.color}20`, color: stat.color, padding: '0.25rem 0.5rem', borderRadius: '0.375rem', fontSize: '0.75rem', fontWeight: 700 }}>
                                Live
                            </div>
                        </div>
                        <div style={{ fontSize: '2rem', fontWeight: 800 }}>{stat.value}</div>
                        <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>{stat.label}</div>
                    </div>
                ))}
            </div>

            {/* Main Content Area */}
            <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '2rem' }}>
                {/* Catalog Preview */}
                <section className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <h3 style={{ fontWeight: 600 }}>Catalog Overview</h3>
                        <button className="btn-primary" style={{ fontSize: '0.75rem' }}>View All</button>
                    </div>
                    <div style={{ padding: '1rem' }}>
                        {loading ? (
                            <p style={{ padding: '1rem', textAlign: 'center' }}>Loading catalog...</p>
                        ) : itemTypes.length === 0 ? (
                            <p style={{ padding: '1rem', textAlign: 'center', color: 'var(--text-muted)' }}>No items found.</p>
                        ) : (
                            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                                <thead>
                                    <tr style={{ textAlign: 'left', color: 'var(--text-muted)' }}>
                                        <th style={{ padding: '0.75rem' }}>Code</th>
                                        <th style={{ padding: '0.75rem' }}>Name</th>
                                        <th style={{ padding: '0.75rem' }}>Kind</th>
                                        <th style={{ padding: '0.75rem', textAlign: 'right' }}>Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {itemTypes.slice(0, 5).map(it => (
                                        <tr key={it.id} style={{ borderBottom: '1px solid var(--border)', transition: 'background 0.2s' }}>
                                            <td style={{ padding: '0.75rem' }}><code>{it.code}</code></td>
                                            <td style={{ padding: '0.75rem', fontWeight: 500 }}>{it.name}</td>
                                            <td style={{ padding: '0.75rem' }}>
                                                <span style={{ textTransform: 'capitalize', background: 'var(--surface)', padding: '0.125rem 0.375rem', borderRadius: '0.25rem' }}>
                                                    {it.kind}
                                                </span>
                                            </td>
                                            <td style={{ padding: '0.75rem', textAlign: 'right' }}>
                                                <ChevronRight size={16} color="var(--text-muted)" />
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        )}
                    </div>
                </section>

                {/* Sidebar / Recent Activity placeholder */}
                <section className="glass" style={{ borderRadius: '1rem', padding: '1.5rem' }}>
                    <h3 style={{ fontWeight: 600, marginBottom: '1.5rem' }}>System Status</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                        {[
                            { label: 'API Gateway', status: 'Online', color: 'var(--success)' },
                            { label: 'Fleet Registry', status: 'Online', color: 'var(--success)' },
                            { label: 'Database', status: 'Syncing', color: 'var(--warning)' },
                            { label: 'Remote Agent', status: 'Standby', color: 'var(--text-muted)' },
                        ].map((s, i) => (
                            <div key={i} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <span style={{ fontSize: '0.875rem' }}>{s.label}</span>
                                <span style={{ fontSize: '0.75rem', fontWeight: 600, color: s.color }}>{s.status}</span>
                            </div>
                        ))}
                    </div>
                </section>
            </div>
        </div>
    );
};

export default Dashboard;
