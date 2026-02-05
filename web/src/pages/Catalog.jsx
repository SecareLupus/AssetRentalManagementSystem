import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import { Box, Search, Filter, ChevronRight, CheckCircle2, Clock } from 'lucide-react';

const Catalog = () => {
    const [itemTypes, setItemTypes] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        const fetchCatalog = async () => {
            try {
                const response = await axios.get('/v1/catalog/item-types');
                setItemTypes(response.data || []);
            } catch (error) {
                console.error("Failed to fetch catalog", error);
            } finally {
                setLoading(false);
            }
        };
        fetchCatalog();
    }, []);

    const filteredItems = itemTypes.filter(item =>
        item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.code.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end' }}>
                <div>
                    <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '0.5rem' }}>Equipment Catalog</h1>
                    <p style={{ color: 'var(--text-muted)' }}>Browse and reserve fleet assets.</p>
                </div>
                <Link to="/reserve" className="btn-primary" style={{ textDecoration: 'none' }}>
                    New Reservation
                </Link>
            </header>

            {/* Search & Filter Bar */}
            <div className="glass" style={{ padding: '1rem', borderRadius: '0.75rem', marginBottom: '2rem', display: 'flex', gap: '1rem' }}>
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
                <button className="glass" style={{ padding: '0.5rem 1rem', borderRadius: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text)' }}>
                    <Filter size={18} /> Filters
                </button>
            </div>

            {loading ? (
                <div style={{ textAlign: 'center', padding: '4rem' }}>
                    <div className="animate-spin" style={{ marginBottom: '1rem' }}><Box /></div>
                    <p>Loading catalog items...</p>
                </div>
            ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                    {filteredItems.map(item => (
                        <div key={item.id} className="glass" style={{ borderRadius: '1rem', overflow: 'hidden', transition: 'transform 0.2s', cursor: 'default' }}>
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
                                <button className="btn-primary" style={{ fontSize: '0.75rem', padding: '0.4rem 0.75rem' }}>
                                    Add to Cart
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Catalog;
