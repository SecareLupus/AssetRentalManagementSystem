import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';
import { Box, ArrowLeft, Shield, Cpu, Wifi, Settings, Activity, History, Info } from 'lucide-react';

const ItemTypeDetails = () => {
    const { id } = useParams();
    const [item, setItem] = useState(null);
    const [assets, setAssets] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const itemRes = await axios.get(`/v1/fleet/item-types/${id}`);
                setItem(itemRes.data);

                const assetsRes = await axios.get(`/v1/inventory/assets?item_type_id=${id}`);
                setAssets(assetsRes.data || []);
            } catch (error) {
                console.error("Error fetching item details", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Loading details...</div>;
    if (!item) return <div style={{ padding: '4rem', textAlign: 'center' }}>Item not found.</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <Link to="/catalog" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', textDecoration: 'none', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Back to Catalog
            </Link>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 350px', gap: '2rem' }}>
                {/* Main Content */}
                <div>
                    <header style={{ marginBottom: '2rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.5rem' }}>
                            <h1 style={{ fontSize: '2.5rem', fontWeight: 800 }}>{item.name}</h1>
                            <span style={{ background: 'var(--primary)20', color: 'var(--primary)', padding: '0.25rem 0.75rem', borderRadius: '2rem', fontSize: '0.875rem', fontWeight: 600 }}>
                                {item.kind}
                            </span>
                        </div>
                        <p style={{ fontSize: '1.25rem', color: 'var(--text-muted)' }}>Code: <code>{item.code}</code></p>
                    </header>

                    <section className="glass" style={{ padding: '2rem', borderRadius: '1rem', marginBottom: '2rem' }}>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Shield size={20} color="var(--primary)" /> Supported Features
                        </h2>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.5rem' }}>
                            <FeatureItem icon={Cpu} label="Remote Management" active={item.supported_features?.remote_management} />
                            <FeatureItem icon={Activity} label="Provisioning" active={item.supported_features?.provisioning} />
                            <FeatureItem icon={History} label="Build Spec Tracking" active={item.supported_features?.build_spec_tracking} />
                            <FeatureItem icon={Wifi} label="Telemetry" active={item.supported_features?.telemetry} />
                        </div>
                    </section>

                    <section className="glass" style={{ padding: '2rem', borderRadius: '1rem' }}>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem' }}>Specific Assets ({assets.length})</h2>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            {assets.map(asset => (
                                <div key={asset.id} style={{
                                    padding: '1rem',
                                    borderRadius: '0.75rem',
                                    background: 'var(--surface)',
                                    border: '1px solid var(--border)',
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center'
                                }}>
                                    <div>
                                        <div style={{ fontWeight: 600 }}>{asset.asset_tag || asset.serial_number || `Asset #${asset.id}`}</div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{asset.location || 'Unknown Location'}</div>
                                    </div>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
                                        <StatusBadge status={asset.status} />
                                        <button className="glass" style={{ padding: '0.4rem', borderRadius: '0.4rem' }}>
                                            <Settings size={14} />
                                        </button>
                                    </div>
                                </div>
                            ))}
                            {assets.length === 0 && <p style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '1rem' }}>No individual assets found for this type.</p>}
                        </div>
                    </section>
                </div>

                {/* Sidebar / Quick Actions */}
                <aside>
                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', position: 'sticky', top: '2rem' }}>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem' }}>Quick Reserve</h3>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Quantity</label>
                            <input type="number" defaultValue="1" style={{ width: '100%', padding: '0.75rem', background: 'var(--surface)', border: '1px solid var(--border)', borderRadius: '0.5rem', color: 'var(--text)' }} />
                        </div>
                        <button className="btn-primary" style={{ width: '100%', marginBottom: '1rem' }}>Request Reservation</button>
                        <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', textAlign: 'center' }}>
                            Expected availability: <span style={{ color: 'var(--success)' }}>Immediate</span>
                        </p>

                        <hr style={{ margin: '1.5rem 0', border: 'none', borderTop: '1px solid var(--border)' }} />

                        <div style={{ display: 'flex', alignItems: 'flex-start', gap: '0.75rem' }}>
                            <Info size={16} color="var(--primary)" style={{ marginTop: '0.125rem' }} />
                            <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', lineHeight: 1.4 }}>
                                This item requires a technician inspection before it can be deployed to a client site.
                            </p>
                        </div>
                    </div>
                </aside>
            </div>
        </div>
    );
};

const FeatureItem = ({ icon: Icon, label, active }) => (
    <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', opacity: active ? 1 : 0.3 }}>
        <div style={{ background: active ? 'var(--primary)20' : 'var(--surface)', padding: '0.5rem', borderRadius: '0.5rem' }}>
            <Icon size={18} color={active ? 'var(--primary)' : 'var(--text-muted)'} />
        </div>
        <span style={{ fontSize: '0.875rem', fontWeight: 500 }}>{label}</span>
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
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.375rem' }}>
            <div style={{ width: '6px', height: '6px', borderRadius: '50%', background: color }} />
            <span style={{ fontSize: '0.75rem', textTransform: 'capitalize', fontWeight: 600, color }}>{status}</span>
        </div>
    );
};

export default ItemTypeDetails;
