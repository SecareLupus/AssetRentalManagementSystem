import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { BarChart3, Download, Filter, TrendingUp, PieChart, Clock, Wrench, Package } from 'lucide-react';

const FleetReports = () => {
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState(null);
    const [timeframe, setTimeframe] = useState('30d');

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const res = await axios.get('/v1/dashboard/stats');
                setStats(res.data);
            } catch (err) {
                console.error("Failed to load report data", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Generating Intelligence Report...</div>;

    const utilizationRate = stats ? ((stats.active_rentals / (stats.total_assets || 1)) * 100).toFixed(1) : 0;
    const maintenanceRate = stats ? ((stats.assets_by_status?.maintenance / (stats.total_assets || 1)) * 100).toFixed(1) : 0;

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '3rem' }}>
                <div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.75rem' }}>
                        <div style={{ background: 'var(--primary)20', padding: '0.75rem', borderRadius: '0.75rem' }}>
                            <BarChart3 size={24} color="var(--primary)" />
                        </div>
                        <h1 style={{ fontSize: '1.75rem', fontWeight: 800 }}>Fleet Operations Intelligence</h1>
                    </div>
                    <p style={{ color: 'var(--text-muted)' }}>High-level utilization, maintenance, and lifecycle reporting.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <div style={{ position: 'relative' }}>
                        <select
                            className="glass"
                            style={{ padding: '0.625rem 2rem 0.625rem 1rem', borderRadius: '0.5rem', appearance: 'none', color: 'white', fontSize: '0.875rem' }}
                            value={timeframe}
                            onChange={(e) => setTimeframe(e.target.value)}
                        >
                            <option value="7d">Last 7 Days</option>
                            <option value="30d">Last 30 Days</option>
                            <option value="90d">Last quarter</option>
                        </select>
                        <Filter size={14} style={{ position: 'absolute', right: '0.75rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.5 }} />
                    </div>
                    <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem' }}>
                        <Download size={16} /> Export CSV
                    </button>
                </div>
            </header>

            {/* Top Level Metrics */}
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1.5rem', marginBottom: '3rem' }}>
                <ReportCard label="Fleet Utilization" value={`${utilizationRate}%`} trend="+2.4%" color="var(--primary)" icon={TrendingUp} />
                <ReportCard label="Maintenance Ratio" value={`${maintenanceRate}%`} trend="-0.5%" color="var(--error)" icon={Wrench} />
                <ReportCard label="Average Rental Term" value="12.4 days" trend="+1.2d" color="var(--success)" icon={Clock} />
                <ReportCard label="Inventory Velocity" value="84%" trend="Stable" color="var(--warning)" icon={PieChart} />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '2fr 1.2fr', gap: '2rem' }}>
                {/* Status Breakdown Bar chart simulation */}
                <section className="glass" style={{ padding: '2rem', borderRadius: '1.5rem' }}>
                    <h3 style={{ fontWeight: 600, marginBottom: '2rem' }}>Inventory Distribution by Status</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                        {Object.entries(stats?.assets_by_status || {}).map(([status, count]) => (
                            <div key={status}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem', fontSize: '0.875rem' }}>
                                    <span style={{ textTransform: 'capitalize' }}>{status}</span>
                                    <span style={{ fontWeight: 700 }}>{count} ({((count / (stats.total_assets || 1)) * 100).toFixed(0)}%)</span>
                                </div>
                                <div style={{ height: '12px', background: 'rgba(255,255,255,0.05)', borderRadius: '6px', overflow: 'hidden' }}>
                                    <div style={{
                                        width: `${(count / (stats.total_assets || 1)) * 100}%`,
                                        height: '100%',
                                        background: status === 'available' ? 'var(--success)' : status === 'maintenance' ? 'var(--error)' : status === 'deployed' ? 'var(--primary)' : 'var(--text-muted)'
                                    }} />
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                <section className="glass" style={{ padding: '2rem', borderRadius: '1.5rem' }}>
                    <h3 style={{ fontWeight: 600, marginBottom: '1.5rem' }}>Operational Health</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                        <HealthItem label="Database Latency" status="Optimal" value="14ms" color="var(--success)" />
                        <HealthItem label="API Error Rate" status="Nominal" value="0.02%" color="var(--success)" />
                        <HealthItem label="Recall Response" status="Excellent" value="98%" color="var(--success)" />
                        <HealthItem label="Audit Coverage" status="Incomplete" value="72%" color="var(--warning)" />
                    </div>

                    <div style={{ marginTop: '2.5rem', padding: '1.25rem', borderRadius: '1rem', background: 'rgba(255,255,255,0.02)', border: '1px solid var(--border)' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem' }}>
                            <Package size={18} color="var(--primary)" />
                            <span style={{ fontWeight: 700, fontSize: '0.875rem' }}>Warehouse Efficiency</span>
                        </div>
                        <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Technician throughput has increased by 15% since the implementation of high-efficiency batch processing in the kiosk.</p>
                    </div>
                </section>
            </div>
        </div>
    );
};

const ReportCard = ({ label, value, trend, color, icon: Icon }) => (
    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1.25rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
            <div style={{ background: `${color}15`, padding: '0.5rem', borderRadius: '0.5rem' }}>
                <Icon size={20} color={color} />
            </div>
            <span style={{ fontSize: '0.75rem', color: trend.startsWith('+') ? 'var(--success)' : trend.startsWith('-') ? 'var(--error)' : 'var(--text-muted)', fontWeight: 700 }}>{trend}</span>
        </div>
        <div style={{ fontSize: '1.75rem', fontWeight: 800 }}>{value}</div>
        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 600 }}>{label}</div>
    </div>
);

const HealthItem = ({ label, status, value, color }) => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
            <div style={{ fontSize: '0.875rem', fontWeight: 600 }}>{label}</div>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{status}</div>
        </div>
        <div style={{ fontWeight: 700, color, fontSize: '1rem' }}>{value}</div>
    </div>
);

export default FleetReports;
