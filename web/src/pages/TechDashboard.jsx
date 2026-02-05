import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import { Wrench, Activity, Zap, ClipboardCheck, ArrowLeftRight, Clock, ChevronRight, AlertCircle, Box } from 'lucide-react';

const TechDashboard = () => {
    const [tasks, setTasks] = useState({
        inspections: [],
        provisioning: [],
        returns: []
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const res = await axios.get('/v1/inventory/assets');
                const assets = res.data || [];

                // Logical filtering for "To-Do" items
                setTasks({
                    inspections: assets.filter(a => !a.last_inspection_at || a.status === 'maintenance'),
                    provisioning: assets.filter(a => a.provisioning_status && a.provisioning_status !== 'complete'),
                    returns: assets.filter(a => a.status === 'deployed') // In a real app, maybe filter by "upcoming returns"
                });
            } catch (err) {
                console.error("Failed to fetch tech tasks", err);
            } finally {
                setLoading(false);
            }
        };
        fetchTasks();
    }, []);

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2.5rem' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '0.5rem' }}>Maintenance Station</h1>
                <p style={{ color: 'var(--text-muted)' }}>Technician operations and fleet lifecycle management.</p>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem', marginBottom: '3rem' }}>
                <TaskStatCard icon={ClipboardCheck} label="Inspections Due" count={tasks.inspections.length} color="var(--warning)" />
                <TaskStatCard icon={Zap} label="Provisioning Tasks" count={tasks.provisioning.length} color="var(--primary)" />
                <TaskStatCard icon={ArrowLeftRight} label="Returns to Process" count={tasks.returns.length} color="var(--success)" />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: '2rem' }}>
                {/* Inspection Queue */}
                <section className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <h3 style={{ fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <ClipboardCheck size={20} color="var(--warning)" /> Inspection Queue
                        </h3>
                    </div>
                    <div style={{ padding: '0.5rem' }}>
                        {tasks.inspections.length === 0 ? (
                            <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>All assets inspected!</div>
                        ) : (
                            tasks.inspections.slice(0, 8).map(asset => (
                                <div key={asset.id} style={{
                                    padding: '1rem',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'space-between',
                                    borderBottom: '1px solid var(--border)'
                                }}>
                                    <div>
                                        <div style={{ fontWeight: 600 }}>{asset.asset_tag || asset.serial_number}</div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Last Inspected: {asset.last_inspection_at ? new Date(asset.last_inspection_at).toLocaleDateString() : 'Never'}</div>
                                    </div>
                                    <Link to={`/tech/inspect/${asset.id}`} className="btn-primary" style={{ fontSize: '0.75rem', padding: '0.4rem 0.8rem', textDecoration: 'none' }}>
                                        Start Inspection
                                    </Link>
                                </div>
                            ))
                        )}
                    </div>
                </section>

                {/* Provisioning & System Status */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                    <section className="glass" style={{ borderRadius: '1rem', padding: '1.5rem' }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Zap size={20} color="var(--primary)" /> Active Provisioning
                        </h3>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            {tasks.provisioning.map(asset => (
                                <div key={asset.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <div style={{ fontSize: '0.875rem' }}>
                                        <div style={{ fontWeight: 600 }}>{asset.asset_tag}</div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{asset.provisioning_status}</div>
                                    </div>
                                    <Link to={`/tech/provision/${asset.id}`} style={{ color: 'var(--primary)' }}><ChevronRight size={18} /></Link>
                                </div>
                            ))}
                            {tasks.provisioning.length === 0 && <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>No active provisioning.</p>}
                        </div>
                    </section>

                    <section className="glass" style={{ borderRadius: '1rem', padding: '1.5rem', background: 'var(--error)05' }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '1rem', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--error)' }}>
                            <AlertCircle size={20} /> High Severity Issues
                        </h3>
                        <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>No critical failures detected in fleet telemetry.</p>
                    </section>
                </div>
            </div>
        </div>
    );
};

const TaskStatCard = ({ icon: Icon, label, count, color }) => (
    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
        <div style={{ background: `${color}15`, padding: '1rem', borderRadius: '0.75rem' }}>
            <Icon size={24} color={color} />
        </div>
        <div>
            <div style={{ fontSize: '1.5rem', fontWeight: 800 }}>{count}</div>
            <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>{label}</div>
        </div>
    </div>
);

export default TechDashboard;
