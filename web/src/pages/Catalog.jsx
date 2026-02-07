import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import { Box, Search, Filter, ChevronRight, CheckCircle2, Clock, Plus } from 'lucide-react';
import { GlassCard, PageHeader, StatusBadge } from '../components/Shared';

const Catalog = () => {
    const [itemTypes, setItemTypes] = useState([]);
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [includeInactive, setIncludeInactive] = useState(false);

    // Form State
    const [formData, setFormData] = useState({
        name: '',
        code: '',
        kind: 'serialized',
        is_active: true,
        supported_features: {
            remote_management: false,
            provisioning: false,
            refurbishment: false,
            build_spec_tracking: false
        }
    });

    const fetchCatalog = async () => {
        try {
            setLoading(true);
            const response = await axios.get(`/v1/catalog/item-types${includeInactive ? '?include_inactive=true' : ''}`);
            setItemTypes(response.data || []);
        } catch (error) {
            console.error("Failed to fetch catalog", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCatalog();
    }, [includeInactive]);

    const filteredItems = itemTypes.filter(item =>
        item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.code.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const handleCreateItem = async (e) => {
        e.preventDefault();
        try {
            // Convert features to match struct if needed, but direct mapping should work based on struct JSON tags
            await axios.post('/v1/catalog/item-types', formData);
            alert("Item Type created successfully!");
            setShowCreateModal(false);
            setFormData({
                name: '',
                code: '',
                kind: 'serialized',
                is_active: true,
                supported_features: {
                    remote_management: false,
                    provisioning: false,
                    refurbishment: false,
                    build_spec_tracking: false
                }
            });
            fetchCatalog();
        } catch (error) {
            alert("Failed to create item: " + (error.response?.data || error.message));
        }
    };

    const handleFeatureChange = (feature) => {
        setFormData(prev => ({
            ...prev,
            supported_features: {
                ...prev.supported_features,
                [feature]: !prev.supported_features[feature]
            }
        }));
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <PageHeader
                title="Equipment Catalog"
                subtitle="Browse and manage fleet assets."
                actions={
                    <button onClick={() => setShowCreateModal(true)} className="btn-primary">
                        <Plus size={18} /> New Item
                    </button>
                }
            />

            {/* Search Bar */}
            <GlassCard className="flex-between" style={{ padding: '1rem', marginBottom: '2rem', display: 'flex', gap: '1rem' }}>
                <div style={{ flex: 1, position: 'relative' }}>
                    <Search size={18} style={{ position: 'absolute', left: '1rem', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                    <input
                        type="text"
                        placeholder="Search equipment by name or code..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '0.75rem 1rem 0.75rem 3rem',
                            background: 'var(--surface)',
                            border: '1px solid var(--border)',
                            borderRadius: '0.5rem',
                            color: 'var(--text)',
                            fontSize: '0.875rem'
                        }}
                    />
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', cursor: 'pointer', color: 'var(--text-muted)' }}>
                        <input
                            type="checkbox"
                            checked={includeInactive}
                            onChange={(e) => setIncludeInactive(e.target.checked)}
                        />
                        Show Archived
                    </label>
                </div>
            </GlassCard>

            {loading ? (
                <div style={{ textAlign: 'center', padding: '4rem' }}>
                    <div className="animate-spin" style={{ marginBottom: '1rem' }}><Box /></div>
                    <p>Loading catalog items...</p>
                </div>
            ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                    {filteredItems.map(item => (
                        <GlassCard key={item.id} style={{ padding: '0', overflow: 'hidden', cursor: 'default' }}>
                            <div style={{ padding: '1.5rem' }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                                    <code style={{ fontSize: '0.75rem', color: 'var(--primary)', fontWeight: 700 }}>{item.code}</code>
                                    <span style={{
                                        fontSize: '0.65rem',
                                        textTransform: 'uppercase',
                                        fontWeight: 800,
                                        background: 'var(--surface)',
                                        padding: '0.25rem 0.5rem',
                                        borderRadius: '0.25rem'
                                    }}>
                                        {item.kind}
                                    </span>
                                </div>
                                <h3 style={{ fontSize: '1.25rem', fontWeight: 700, marginBottom: '0.5rem' }}>{item.name}</h3>

                                <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem', fontSize: '0.875rem' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem', color: 'var(--success)' }}>
                                        <CheckCircle2 size={16} /> <span>Available</span>
                                    </div>
                                </div>
                            </div>

                            <div style={{
                                padding: '1rem 1.5rem',
                                background: 'rgba(255,255,255,0.03)',
                                borderTop: '1px solid var(--border)',
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center'
                            }}>
                                <Link to={`/catalog/${item.id}`} style={{ fontSize: '0.875rem', color: 'var(--text-muted)', textDecoration: 'none', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                    View Details <ChevronRight size={14} />
                                </Link>
                                <button
                                    className="btn-primary"
                                    style={{ fontSize: '0.75rem', padding: '0.4rem 0.75rem' }}
                                    onClick={() => navigate(`/reserve?item_type_id=${item.id}`)}
                                >
                                    Add to Cart
                                </button>
                            </div>
                        </GlassCard>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            {
                showCreateModal && (
                    <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100 }}>
                        <GlassCard style={{ width: '500px', padding: '2rem' }}>
                            <h2 style={{ marginBottom: '1.5rem', fontSize: '1.25rem', fontWeight: 700 }}>Create New Item Type</h2>
                            <form onSubmit={handleCreateItem}>
                                <div style={{ marginBottom: '1rem' }}>
                                    <label style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.5rem', color: 'var(--text-muted)' }}>Name</label>
                                    <input
                                        type="text"
                                        required
                                        value={formData.name}
                                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                                        style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                    />
                                </div>
                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
                                    <div>
                                        <label style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.5rem', color: 'var(--text-muted)' }}>Code (SKU)</label>
                                        <input
                                            type="text"
                                            required
                                            value={formData.code}
                                            onChange={e => setFormData({ ...formData, code: e.target.value })}
                                            style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                        />
                                    </div>
                                    <div>
                                        <label style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.5rem', color: 'var(--text-muted)' }}>Kind</label>
                                        <select
                                            value={formData.kind}
                                            onChange={e => setFormData({ ...formData, kind: e.target.value })}
                                            style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', background: 'var(--surface)', border: '1px solid var(--border)', color: 'var(--text)' }}
                                        >
                                            <option value="serialized">Serialized</option>
                                            <option value="fungible">Fungible</option>
                                            <option value="kit">Kit</option>
                                        </select>
                                    </div>
                                </div>

                                <div style={{ marginBottom: '1.5rem' }}>
                                    <label style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.5rem', color: 'var(--text-muted)' }}>Supported Features</label>
                                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
                                        {Object.keys(formData.supported_features).map(feature => (
                                            <label key={feature} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', cursor: 'pointer' }}>
                                                <input
                                                    type="checkbox"
                                                    checked={formData.supported_features[feature]}
                                                    onChange={() => handleFeatureChange(feature)}
                                                />
                                                <span style={{ textTransform: 'capitalize' }}>{feature.replace(/_/g, ' ')}</span>
                                            </label>
                                        ))}
                                    </div>
                                </div>

                                <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                                    <button type="button" onClick={() => setShowCreateModal(false)} className="glass" style={{ padding: '0.5rem 1rem', borderRadius: '0.5rem' }}>Cancel</button>
                                    <button type="submit" className="btn-primary">Create Item</button>
                                </div>
                            </form>
                        </GlassCard>
                    </div>
                )
            }
        </div >
    );
};

export default Catalog;
