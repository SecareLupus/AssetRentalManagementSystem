import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';
import { ArrowLeft, Cpu, Wifi, Activity, History, Wrench, ShieldCheck, MapPin, Gauge } from 'lucide-react';

const AssetDetails = () => {
    const { id } = useParams();
    const [asset, setAsset] = useState(null);
    const [logs, setLogs] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const assetRes = await axios.get(`/v1/inventory/assets/${id}`);
                setAsset(assetRes.data);

                // Use the maintenance logs endpoint
                const logsRes = await axios.get(`/v1/inventory/assets/${id}/maintenance-logs`);
                setLogs(logsRes.data || []);
            } catch (error) {
                console.error("Error fetching asset details", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Loading asset details...</div>;
    if (!asset) return <div style={{ padding: '4rem', textAlign: 'center' }}>Asset not found.</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <Link to="/catalog" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', textDecoration: 'none', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Dashboard
            </Link>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 300px', gap: '2rem' }}>
                <div>
                    <header style={{ marginBottom: '2.5rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.5rem' }}>
                            <h1 style={{ fontSize: '2rem', fontWeight: 800 }}>{asset.asset_tag || asset.serial_number || `Asset #${asset.id}`}</h1>
                            <StatusBadge status={asset.status} />
                        </div>
                        <div style={{ display: 'flex', gap: '1.5rem', color: 'var(--text-muted)', fontSize: '0.875rem' }}>
                            <span style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}><MapPin size={14} /> {asset.location || 'Warehouse A'}</span>
                            <span style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}><ShieldCheck size={14} /> ID: {asset.id}</span>
                        </div>
                    </header>

                    <section style={{ marginBottom: '3rem' }}>
                        <h3 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem' }}>System Specifications</h3>
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                            <SpecBox icon={Cpu} label="Firmware" value={asset.firmware_version || 'v1.4.2-stable'} />
                            <SpecBox icon={Gauge} label="Build Spec" value={asset.build_spec_version || 'Default'} />
                            <SpecBox icon={Wifi} label="Remote ID" value={asset.remote_management_id || 'Not Enrolled'} />
                            <SpecBox icon={Activity} label="Health" value="Optimal" />
                        </div>
                    </section>

                    <section>
                        <h3 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <History size={20} /> Maintenance History
                        </h3>
                        <div className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                            {logs.length === 0 ? (
                                <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                                    No maintenance logs recorded.
                                </div>
                            ) : (
                                <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                                    <thead>
                                        <tr style={{ textAlign: 'left', background: 'rgba(255,255,255,0.02)', color: 'var(--text-muted)' }}>
                                            <th style={{ padding: '1rem' }}>Date</th>
                                            <th style={{ padding: '1rem' }}>Action</th>
                                            <th style={{ padding: '1rem' }}>Performed By</th>
                                            <th style={{ padding: '1rem' }}>Notes</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {logs.map(log => (
                                            <tr key={log.id} style={{ borderTop: '1px solid var(--border)' }}>
                                                <td style={{ padding: '1rem' }}>{new Date(log.created_at).toLocaleDateString()}</td>
                                                <td style={{ padding: '1rem' }}><span style={{ textTransform: 'capitalize', fontWeight: 600 }}>{log.action_type}</span></td>
                                                <td style={{ padding: '1rem' }}>{log.performed_by || 'System'}</td>
                                                <td style={{ padding: '1rem', color: 'var(--text-muted)' }}>{log.notes}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    </section>
                </div>

                <aside>
                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', marginBottom: '1.5rem' }}>
                        <h4 style={{ fontWeight: 700, marginBottom: '1rem' }}>Asset Actions</h4>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                            <button onClick={() => alert("Inspection Templates not yet configured in backend.")} className="glass" style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', fontSize: '0.875rem', textAlign: 'left', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                <ShieldCheck size={16} /> Manual Inspection
                            </button>
                            
                            {asset.status !== 'maintenance' ? (
                                <button onClick={async () => {
                                    if(!window.confirm("Mark this asset for repair? It will be unavailable.")) return;
                                    try {
                                        await axios.post(`/v1/inventory/assets/${id}/repair`);
                                        setAsset({...asset, status: 'maintenance'});
                                        alert("Asset marked for repair.");
                                    } catch(e) { alert("Action failed: " + e.message); }
                                }} className="glass" style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', fontSize: '0.875rem', textAlign: 'left', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--warning)' }}>
                                    <Wrench size={16} /> Mark for Repair
                                </button>
                            ) : (
                                <button onClick={async () => {
                                    if(!window.confirm("Mark this asset as Available?")) return;
                                    try {
                                        await axios.patch(`/v1/inventory/assets/${id}/status`, { status: "available" });
                                        setAsset({...asset, status: 'available'});
                                        alert("Asset marked as Available.");
                                    } catch(e) { alert("Action failed: " + e.message); }
                                }} className="btn-primary" style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', fontSize: '0.875rem', textAlign: 'left', display: 'flex', alignItems: 'center', gap: '0.5rem', justifyContent: 'center' }}>
                                    <Activity size={16} /> Return to Service
                                </button>
                            )}
                        </div>
                    </div>
                </aside>
            </div>
        </div>
    );
};

const SpecBox = ({ icon: Icon, label, value }) => (
    <div className="glass" style={{ padding: '1rem', borderRadius: '0.75rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
        <Icon size={20} color="var(--primary)" />
        <div>
            <div style={{ fontSize: '0.65rem', color: 'var(--text-muted)', textTransform: 'uppercase', fontWeight: 700 }}>{label}</div>
            <div style={{ fontSize: '0.875rem', fontWeight: 600 }}>{value}</div>
        </div>
    </div>
);

const StatusBadge = ({ status }) => {
    const colors = {
        available: 'var(--success)',
        reserved: 'var(--warning)',
        maintenance: 'var(--error)',
        deployed: 'var(--primary)'
    };
    const color = colors[status] || 'var(--text-muted)';

    return (
        <div style={{ padding: '0.25rem 0.75rem', borderRadius: '1rem', background: `${color}20`, color: color, fontSize: '0.75rem', fontWeight: 800, textTransform: 'uppercase' }}>
            {status}
        </div>
    );
};

export default AssetDetails;
