import React, { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { ArrowLeft, Shield, Cpu, Wifi, Settings, Activity, History, Info, Plus, Package } from 'lucide-react';
import { GlassCard, PageHeader, StatusBadge } from '../components/Shared';

const ItemTypeDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [item, setItem] = useState(null);
    const [assets, setAssets] = useState([]);
    const [places, setPlaces] = useState([]);
    const [loading, setLoading] = useState(true);

    // Modals
    const [showEditModal, setShowEditModal] = useState(false);
    const [showAddAssetModal, setShowAddAssetModal] = useState(false);

    // Form States
    const [editForm, setEditForm] = useState({});
    const [assetForm, setAssetForm] = useState({
        asset_tag: '',
        serial_number: '',
        place_id: '',
        status: 'available'
    });

    const fetchData = async () => {
        try {
            const itemRes = await axios.get(`/v1/catalog/item-types/${id}`);
            setItem(itemRes.data);
            setEditForm(itemRes.data); // Initialize edit form

            const assetsRes = await axios.get(`/v1/inventory/assets?item_type_id=${id}`);
            setAssets(assetsRes.data || []);

            const placesRes = await axios.get('/v1/entities/places');
            setPlaces(placesRes.data || []);
        } catch (error) {
            console.error("Error fetching item details", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [id]);

    const handleEditSubmit = async (e) => {
        e.preventDefault();
        try {
            await axios.put(`/v1/catalog/item-types/${id}`, editForm);
            alert("Item Type updated!");
            setShowEditModal(false);
            fetchData();
        } catch (error) {
            alert("Update failed: " + (error.response?.data || error.message));
        }
    };

    const handleAddAsset = async (e) => {
        e.preventDefault();
        try {
            await axios.post('/v1/inventory/assets', {
                ...assetForm,
                item_type_id: parseInt(id)
            });
            alert("Asset created!");
            setShowAddAssetModal(false);
            setAssetForm({ asset_tag: '', serial_number: '', place_id: '', status: 'available' });
            fetchData();
        } catch (error) {
            alert("Asset creation failed: " + (error.response?.data || error.message));
        }
    };

    const toggleArchive = async () => {
        const action = item.is_active ? "archive" : "restore";
        if (window.confirm(`Are you sure you want to ${action} this item type?`)) {
            try {
                if (item.is_active) {
                    await axios.delete(`/v1/catalog/item-types/${id}`);
                } else {
                    await axios.put(`/v1/catalog/item-types/${id}`, { ...item, is_active: true });
                }
                fetchData();
            } catch (err) {
                alert(`Failed to ${action} item type.`);
            }
        }
    };

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Loading...</div>;
    if (!item) return <div style={{ padding: '4rem', textAlign: 'center' }}>Item type not found.</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <Link to="/catalog" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', textDecoration: 'none', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Back to Catalog
            </Link>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 350px', gap: '2rem' }}>
                {/* Main Content */}
                <div>
                    <PageHeader
                        title={item.name}
                        subtitle={<span>Code: <code>{item.code}</code> | Kind: <span style={{ textTransform: 'capitalize' }}>{item.kind}</span></span>}
                        actions={
                            <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                                <span style={{ background: 'var(--primary)20', color: 'var(--primary)', padding: '0.25rem 0.75rem', borderRadius: '2rem', fontSize: '0.875rem', fontWeight: 600 }}>
                                    {item.kind}
                                </span>
                                <button onClick={() => setShowEditModal(true)} className="glass" style={{ padding: '0.4rem 0.75rem', fontSize: '0.875rem' }}>
                                    Edit
                                </button>
                                <button
                                    onClick={toggleArchive}
                                    className="glass"
                                    style={{
                                        padding: '0.4rem 0.75rem',
                                        fontSize: '0.875rem',
                                        color: item.is_active ? 'var(--error)' : 'var(--success)',
                                        borderColor: item.is_active ? 'var(--error)30' : 'var(--success)30'
                                    }}
                                >
                                    <Package size={16} /> {item.is_active ? 'Archive Type' : 'Restore Type'}
                                </button>
                            </div>
                        }
                    />

                    <GlassCard style={{ marginBottom: '2rem' }}>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Shield size={20} color="var(--primary)" /> Supported Features
                        </h2>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.5rem' }}>
                            <FeatureItem icon={Cpu} label="Remote Management" active={item.supported_features?.remote_management} />
                            <FeatureItem icon={Activity} label="Provisioning" active={item.supported_features?.provisioning} />
                            <FeatureItem icon={History} label="Build Spec Tracking" active={item.supported_features?.build_spec_tracking} />
                            <FeatureItem icon={Wifi} label="Telemetry" active={item.supported_features?.telemetry} />
                        </div>
                    </GlassCard>

                    <GlassCard>
                        <div className="flex-between" style={{ marginBottom: '1.5rem' }}>
                            <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>Specific Assets ({assets.length})</h2>
                            <button onClick={() => setShowAddAssetModal(true)} className="btn-primary" style={{ fontSize: '0.875rem', padding: '0.4rem 0.75rem' }}>
                                <Plus size={16} /> Add Asset
                            </button>
                        </div>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            {assets.map(asset => (
                                <div
                                    key={asset.id}
                                    style={{
                                        padding: '1rem',
                                        borderRadius: '0.75rem',
                                        background: 'var(--surface)',
                                        border: '1px solid var(--border)',
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        cursor: 'pointer'
                                    }}
                                    onClick={() => navigate(`/assets/${asset.id}`)}
                                    onMouseEnter={(e) => e.currentTarget.style.borderColor = 'var(--primary)50'}
                                    onMouseLeave={(e) => e.currentTarget.style.borderColor = 'var(--border)'}
                                >
                                    <div>
                                        <div style={{ fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                            {asset.asset_tag || asset.serial_number || `Asset #${asset.id}`}
                                            {item.supported_features?.remote_management && asset.remote_management_id && (
                                                <div className="pulse" style={{ width: '8px', height: '8px', background: 'var(--success)', borderRadius: '50%', boxShadow: '0 0 8px var(--success)' }} title="Live Connectivity Active" />
                                            )}
                                        </div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{asset.location || 'Unknown Location'}</div>
                                    </div>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
                                        <StatusBadge status={asset.status} />
                                        <button
                                            className="glass"
                                            style={{ padding: '0.4rem', borderRadius: '0.4rem' }}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                // Potential future settings navigation
                                            }}
                                        >
                                            <Settings size={14} />
                                        </button>
                                    </div>
                                </div>
                            ))}
                            {assets.length === 0 && <p style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '1rem' }}>No individual assets found for this type.</p>}
                        </div>
                    </GlassCard>
                </div>

                {/* Sidebar / Quick Actions */}
                <aside>
                    <GlassCard style={{ position: 'sticky', top: '2rem' }}>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem' }}>Quick Reserve</h3>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Quantity</label>
                            <input type="number" defaultValue="1" style={{ width: '100%', padding: '0.75rem', background: 'var(--surface)', border: '1px solid var(--border)', borderRadius: '0.5rem', color: 'var(--text)' }} />
                        </div>
                        <button
                            className="btn-primary"
                            style={{ width: '100%', marginBottom: '1rem', justifyContent: 'center' }}
                            onClick={() => navigate(`/reserve?item_type_id=${id}`)}
                        >
                            Request Reservation
                        </button>
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
                    </GlassCard>
                </aside>
            </div>

            {/* Edit Modal */}
            {showEditModal && (
                <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100 }}>
                    <GlassCard style={{ width: '500px', padding: '2rem' }}>
                        <h2 style={{ marginBottom: '1.5rem' }}>Edit Item Type</h2>
                        <form onSubmit={handleEditSubmit}>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Name</label>
                                <input value={editForm.name} onChange={e => setEditForm({ ...editForm, name: e.target.value })} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
                            </div>
                            <div style={{ marginBottom: '1.5rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.75rem', fontSize: '0.875rem', fontWeight: 600 }}>Enabled Features</label>
                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' }}>
                                    {[
                                        { key: 'remote_management', label: 'Remote Mgmt' },
                                        { key: 'provisioning', label: 'Provisioning' },
                                        { key: 'refurbishment', label: 'Refurbishment' },
                                        { key: 'build_spec_tracking', label: 'Build Specs' }
                                    ].map(f => (
                                        <label key={f.key} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.75rem', cursor: 'pointer' }}>
                                            <input
                                                type="checkbox"
                                                checked={editForm.supported_features?.[f.key] || false}
                                                onChange={e => setEditForm({
                                                    ...editForm,
                                                    supported_features: {
                                                        ...(editForm.supported_features || {}),
                                                        [f.key]: e.target.checked
                                                    }
                                                })}
                                            />
                                            {f.label}
                                        </label>
                                    ))}
                                </div>
                            </div>

                            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
                                <div>
                                    <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.75rem' }}>Critical Shortage Threshold</label>
                                    <input
                                        type="number"
                                        value={editForm.metadata?.critical_shortage_threshold || 0}
                                        onChange={e => setEditForm({
                                            ...editForm,
                                            metadata: { ...editForm.metadata, critical_shortage_threshold: parseInt(e.target.value) }
                                        })}
                                        style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                    />
                                </div>
                                <div>
                                    <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.75rem' }}>Forecast Horizon (Days)</label>
                                    <input
                                        type="number"
                                        value={editForm.metadata?.forecast_horizon_days || 14}
                                        onChange={e => setEditForm({
                                            ...editForm,
                                            metadata: { ...editForm.metadata, forecast_horizon_days: parseInt(e.target.value) }
                                        })}
                                        style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                    />
                                </div>
                            </div>

                            <div style={{ marginBottom: '1.5rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.75rem' }}>Tag Auto-Mapping (Regex)</label>
                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                                    <input
                                        placeholder="Display Pattern (e.g. SN:XXXXX)"
                                        value={editForm.metadata?.tag_display_pattern || ''}
                                        onChange={e => setEditForm({
                                            ...editForm,
                                            metadata: { ...editForm.metadata, tag_display_pattern: e.target.value }
                                        })}
                                        style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                    />
                                    <input
                                        placeholder="Regex (e.g. SN:([0-9A-Z]+))"
                                        value={editForm.metadata?.tag_extract_regex || ''}
                                        onChange={e => setEditForm({
                                            ...editForm,
                                            metadata: { ...editForm.metadata, tag_extract_regex: e.target.value }
                                        })}
                                        style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                    />
                                </div>
                            </div>
                            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                                <button type="button" onClick={() => setShowEditModal(false)} className="glass" style={{ padding: '0.5rem 1rem' }}>Cancel</button>
                                <button type="submit" className="btn-primary">Save Changes</button>
                            </div>
                        </form>
                    </GlassCard>
                </div>
            )}

            {/* Add Asset Modal */}
            {showAddAssetModal && (
                <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100 }}>
                    <GlassCard style={{ width: '500px', padding: '2rem' }}>
                        <h2 style={{ marginBottom: '1.5rem' }}>Add New Asset</h2>
                        <form onSubmit={handleAddAsset}>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Asset Tag</label>
                                <input value={assetForm.asset_tag} onChange={e => setAssetForm({ ...assetForm, asset_tag: e.target.value })} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
                            </div>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Serial Number</label>
                                <input value={assetForm.serial_number} onChange={e => setAssetForm({ ...assetForm, serial_number: e.target.value })} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
                            </div>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Initial Deployment Location</label>
                                <select
                                    value={assetForm.place_id}
                                    onChange={e => setAssetForm({ ...assetForm, place_id: e.target.value ? parseInt(e.target.value) : '' })}
                                    style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                >
                                    <option value="">Status will be assigned based on location...</option>
                                    {places.map(p => (
                                        <option key={p.id} value={p.id}>{p.name} {p.is_internal ? '(Internal)' : '(External)'}</option>
                                    ))}
                                </select>
                                <p style={{ fontSize: '0.65rem', color: 'var(--text-muted)', marginTop: '0.4rem' }}>
                                    System will automatically set status to 'Available' for internal sites or 'Deployed' for client/external sites.
                                </p>
                            </div>
                            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                                <button type="button" onClick={() => setShowAddAssetModal(false)} className="glass" style={{ padding: '0.5rem 1rem' }}>Cancel</button>
                                <button type="submit" className="btn-primary">Create Asset</button>
                            </div>
                        </form>
                    </GlassCard>
                </div>
            )}
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

export default ItemTypeDetails;
