import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Brain, AlertCircle, Calendar, LineChart, ShieldAlert, Zap, ArrowRight, Gauge } from 'lucide-react';
import { Link } from 'react-router-dom';

const IntelligenceOverview = () => {
    const [alerts, setAlerts] = useState([]);
    const [forecasts, setForecasts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [alertsRes, forecastRes] = await Promise.all([
                    axios.get('/v1/intelligence/shortage-alerts'),
                    axios.get('/v1/intelligence/maintenance-forecast')
                ]);
                setAlerts(alertsRes.data || []);
                setForecasts(forecastRes.data || []);
            } catch (err) {
                console.error("Intelligence fetch failed", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '3rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 style={{ fontSize: '2.25rem', fontWeight: 900, marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
                        <Brain size={32} color="var(--primary)" /> Intelligence Hub
                    </h1>
                    <p style={{ color: 'var(--text-muted)' }}>AI-assisted fleet optimization and predictive logistics.</p>
                </div>
                <Link to="/analytics/heatmap" className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', textDecoration: 'none' }}>
                    Go to Heatmap <ArrowRight size={18} />
                </Link>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1.25fr)', gap: '2rem' }}>
                {/* Shortage Alerts */}
                <section className="glass" style={{ borderRadius: '1.5rem', overflow: 'hidden' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <ShieldAlert size={20} color="var(--error)" />
                        <h3 style={{ fontWeight: 700 }}>Critical Shortage Alerts</h3>
                    </div>
                    <div style={{ padding: '1.5rem' }}>
                        {loading ? (
                            <p>Scanning reservations...</p>
                        ) : alerts.length === 0 ? (
                            <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-muted)' }}>
                                <Gauge size={48} style={{ marginBottom: '1rem', opacity: 0.2, margin: '0 auto' }} />
                                <p>No inventory shortages detected in the next 14 days.</p>
                            </div>
                        ) : (
                            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                                {alerts.map((alert, i) => (
                                    <div key={i} style={{ padding: '1rem', borderRadius: '1rem', background: 'rgba(239, 68, 68, 0.05)', border: '1px solid rgba(239, 68, 68, 0.2)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <div>
                                            <div style={{ fontWeight: 700, color: 'var(--error)' }}>{alert.item_type_name}</div>
                                            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Conflict Date: {new Date(alert.date).toLocaleDateString()}</div>
                                        </div>
                                        <div style={{ textAlign: 'right' }}>
                                            <div style={{ fontSize: '1.25rem', fontWeight: 800 }}>-{alert.shortage_count} Units</div>
                                            <div style={{ fontSize: '0.65rem', color: 'var(--text-muted)' }}>Needed: {alert.total_needed} | Owned: {alert.total_owned}</div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </section>

                {/* Predictive Maintenance */}
                <section className="glass" style={{ borderRadius: '1.5rem', overflow: 'hidden' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <Zap size={20} color="var(--warning)" />
                        <h3 style={{ fontWeight: 700 }}>Service Forecast</h3>
                    </div>
                    <div style={{ padding: '1rem' }}>
                        {forecasts.map((f, i) => (
                            <div key={i} style={{ padding: '1rem', borderBottom: i < forecasts.length - 1 ? '1px solid var(--border)' : 'none' }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                        <span style={{ fontWeight: 600 }}>{f.asset_tag}</span>
                                        <button 
                                            className="glass" 
                                            style={{ padding: '0.2rem 0.4rem', fontSize: '0.65rem' }}
                                            onClick={() => alert("Forecast snoozed for 7 days (Logic placeholder)")}
                                        >
                                            Snooze
                                        </button>
                                    </div>
                                    <span style={{ fontSize: '0.75rem', fontWeight: 800, color: f.urgency_score > 0.9 ? 'var(--error)' : 'var(--warning)' }}>
                                        {(f.urgency_score * 100).toFixed(0)}% URGENT
                                    </span>
                                </div>
                                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Reason: {f.reason}</div>
                                <div style={{ fontSize: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                    <Calendar size={12} /> Est. Service: {new Date(f.next_service_date).toLocaleDateString()}
                                </div>
                            </div>
                        ))}
                        {forecasts.length === 0 && !loading && <p style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>Fleet is in optimal condition.</p>}
                    </div>
                </section>
            </div>

            {/* "What-If" Modal Placeholder */}
            <div style={{ marginTop: '3rem' }}>
                <div className="glass" style={{ padding: '2rem', borderRadius: '1.5rem', textAlign: 'center', background: 'linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(16, 185, 129, 0.1))' }}>
                    <Brain size={48} color="var(--primary)" style={{ marginBottom: '1rem', margin: '0 auto' }} />
                    <h3 style={{ fontSize: '1.25rem', fontWeight: 800, marginBottom: '0.5rem' }}>What-If Planning Mode</h3>
                    <p style={{ color: 'var(--text-muted)', maxWidth: '500px', margin: '0 auto 1.5rem' }}>
                        Simulate large-scale deployment requests without committing. See how a new project would impact future fleet health.
                    </p>
                    <Link to="/simulator" className="btn-primary" style={{ textDecoration: 'none' }}>Launch Simulator</Link>
                </div>
            </div>
        </div>
    );
};

export default IntelligenceOverview;
