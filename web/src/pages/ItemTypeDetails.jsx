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
    const [loading, setLoading] = useState(true);
    
    // Modals
    const [showEditModal, setShowEditModal] = useState(false);
    const [showAddAssetModal, setShowAddAssetModal] = useState(false);

    // Form States
    const [editForm, setEditForm] = useState({});
    const [assetForm, setAssetForm] = useState({
        asset_tag: '',
        serial_number: '',
        location: '',
        status: 'available'
    });

    const fetchData = async () => {
        try {
            const itemRes = await axios.get(`/v1/catalog/item-types/${id}`);
            setItem(itemRes.data);
            setEditForm(itemRes.data); // Initialize edit form

            const assetsRes = await axios.get(`/v1/inventory/assets?item_type_id=${id}`);
            setAssets(assetsRes.data || []);
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
            setAssetForm({ asset_tag: '', serial_number: '', location: '', status: 'available' });
            fetchData();
        } catch (error) {
            alert("Asset creation failed: " + (error.response?.data || error.message));
        }
    };

    const handleDelete = async () => {
        if (window.confirm("Are you sure you want to archive this item type? This will make it unavailable for new reservations.")) {
            try {
                await axios.delete(`/v1/catalog/item-types/${id}`);
                navigate('/catalog');
            } catch (err) {
                alert("Failed to delete item type.");
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
                                <button onClick={handleDelete} className="glass" style={{ padding: '0.4rem 0.75rem', fontSize: '0.875rem', color: 'var(--error)', borderColor: 'var(--error)30' }}>
                                    <Package size={16} /> Archive Type
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
                        <button className="btn-primary" style={{ width: '100%', marginBottom: '1rem', justifyContent: 'center' }}>Request Reservation</button>
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
                                <input value={editForm.name} onChange={e => setEditForm({...editForm, name: e.target.value})} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
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
                                <input value={assetForm.asset_tag} onChange={e => setAssetForm({...assetForm, asset_tag: e.target.value})} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
                            </div>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Serial Number</label>
                                <input value={assetForm.serial_number} onChange={e => setAssetForm({...assetForm, serial_number: e.target.value})} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
                            </div>
                            <div style={{ marginBottom: '1rem' }}>
                                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Location</label>
                                <input value={assetForm.location} onChange={e => setAssetForm({...assetForm, location: e.target.value})} style={{ width: '100%', padding: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }} />
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
