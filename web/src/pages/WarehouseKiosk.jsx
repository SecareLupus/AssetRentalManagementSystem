import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Scan, ArrowUpRight, ArrowDownLeft, Box, Search, CheckCircle2, AlertTriangle, ShieldCheck, List, Trash2, Zap } from 'lucide-react';

const WarehouseKiosk = () => {
    const [scanTerm, setScanTerm] = useState('');
    const [scannedAssets, setScannedAssets] = useState([]); // List of assets in the current batch
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState(null);
    const [inventory, setInventory] = useState([]);
    const [transactionContext, setTransactionContext] = useState({
        location: '',
        estimated_return_at: ''
    });

    useEffect(() => {
        const fetchInventory = async () => {
            try {
                const res = await axios.get('/v1/inventory/assets');
                setInventory(res.data || []);
            } catch (err) {
                console.error("Failed to load inventory for kiosk lookup", err);
            }
        };
        fetchInventory();
    }, []);

    const handleScan = (e) => {
        e.preventDefault();
        setMessage(null);
        
        const found = inventory.find(a => a.asset_tag === scanTerm || a.serial_number === scanTerm);
        if (found) {
            if (scannedAssets.find(a => a.id === found.id)) {
                setMessage({ type: 'error', text: 'Asset already in batch.' });
            } else {
                setScannedAssets([found, ...scannedAssets]);
                setScanTerm('');
            }
        } else {
            setMessage({ type: 'error', text: 'Asset not found in registry.' });
        }
    };

    const removeFromBatch = (id) => {
        setScannedAssets(scannedAssets.filter(a => a.id !== id));
    };

    const executeBatchUpdate = async (newStatus) => {
        setLoading(true);
        try {
            // Sequential updates as the backend doesn't have a generic bulk status endpoint yet
            // (Bulk Recall exists but handles a specific flow)
            await Promise.all(scannedAssets.map(a => 
                axios.patch(`/v1/inventory/assets/${a.id}/status`, { 
                    status: newStatus,
                    location: transactionContext.location || undefined,
                    metadata: { 
                        ...a.metadata, 
                        estimated_return_at: transactionContext.estimated_return_at || undefined,
                        transaction_source: 'warehouse_kiosk'
                    }
                })
            ));
            
            setMessage({ type: 'success', text: `Successfully updated ${scannedAssets.length} assets to ${newStatus}.` });
            setScannedAssets([]);
            
            // Refresh local inventory to reflect changes
            const res = await axios.get('/v1/inventory/assets');
            setInventory(res.data || []);
        } catch (err) {
            setMessage({ type: 'error', text: 'One or more updates failed.' });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <header style={{ textAlign: 'center', marginBottom: '3rem' }}>
                <div style={{ display: 'inline-flex', background: 'var(--primary)', padding: '1rem', borderRadius: '1rem', marginBottom: '1.5rem' }}>
                    <Scan size={32} color="white" />
                </div>
                <h1 style={{ fontSize: '2.5rem', fontWeight: 900 }}>Warehouse Kiosk</h1>
                <p style={{ color: 'var(--text-muted)' }}>Scan asset tags or serials for high-efficiency batch processing.</p>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '1.25fr 1fr', gap: '3rem' }}>
                {/* Left Side: Scanning & Actions */}
                <div>
                    <form onSubmit={handleScan} style={{ marginBottom: '2rem' }}>
                        <div style={{ position: 'relative' }}>
                            <input
                                type="text"
                                placeholder="Scan Tag or Serial..."
                                className="glass"
                                style={{ width: '100%', padding: '1.5rem 1.5rem 1.5rem 3.5rem', borderRadius: '1rem', fontSize: '1.25rem', color: 'white' }}
                                value={scanTerm}
                                onChange={(e) => setScanTerm(e.target.value)}
                                autoFocus
                            />
                            <Search size={24} style={{ position: 'absolute', left: '1.25rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.3 }} />
                        </div>
                    </form>
                    
                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', marginBottom: '2rem', display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <div>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Destination / Target Location</label>
                            <input 
                                type="text" 
                                className="glass" 
                                placeholder="e.g. Rack A-1, Client Site X" 
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                value={transactionContext.location}
                                onChange={e => setTransactionContext({ ...transactionContext, location: e.target.value })}
                            />
                        </div>
                        <div>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Est. Return Date (Optional)</label>
                            <input 
                                type="datetime-local" 
                                className="glass" 
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                value={transactionContext.estimated_return_at}
                                onChange={e => setTransactionContext({ ...transactionContext, estimated_return_at: e.target.value })}
                            />
                        </div>
                    </div>

                    {message && (
                        <div className="glass" style={{
                            padding: '1rem',
                            borderRadius: '0.75rem',
                            marginBottom: '2rem',
                            background: message.type === 'error' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                            color: message.type === 'error' ? 'var(--error)' : 'var(--success)',
                            display: 'flex', alignItems: 'center', gap: '0.5rem'
                        }}>
                            {message.type === 'error' ? <AlertTriangle size={18} /> : <CheckCircle2 size={18} />}
                            {message.text}
                        </div>
                    )}

                    <div className="glass" style={{ padding: '2rem', borderRadius: '1.5rem' }}>
                        <h3 style={{ fontWeight: 800, marginBottom: '2rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Zap size={20} color="var(--primary)" /> Batch Operations
                        </h3>
                        
                        <div style={{ display: 'grid', gridTemplateRows: 'repeat(3, 1fr)', gap: '1rem' }}>
                            <KioskAction
                                icon={ArrowUpRight}
                                label="Bulk Check-out (Deploy)"
                                desc={`Transition ${scannedAssets.length} items to Deployed`}
                                color="var(--primary)"
                                onClick={() => executeBatchUpdate('deployed')}
                                disabled={scannedAssets.length === 0 || loading}
                            />
                            <KioskAction
                                icon={ArrowDownLeft}
                                label="Bulk Check-in (Receive)"
                                desc={`Transition ${scannedAssets.length} items to Available`}
                                color="var(--success)"
                                onClick={() => executeBatchUpdate('available')}
                                disabled={scannedAssets.length === 0 || loading}
                            />
                            <KioskAction
                                icon={ShieldCheck}
                                label="Bulk Maintenance / QC"
                                desc={`Transition ${scannedAssets.length} items to Maintenance`}
                                color="var(--error)"
                                onClick={() => executeBatchUpdate('maintenance')}
                                disabled={scannedAssets.length === 0 || loading}
                            />
                        </div>
                    </div>
                </div>

                {/* Right Side: Batch Queue */}
                <div className="glass" style={{ borderRadius: '1.5rem', display: 'flex', flexDirection: 'column', height: '600px' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <List size={20} color="var(--text-muted)" />
                            <h3 style={{ fontWeight: 700 }}>Scanned Batch</h3>
                        </div>
                        <span style={{ fontSize: '0.75rem', fontWeight: 800, background: 'var(--surface)', padding: '0.25rem 0.5rem', borderRadius: '0.5rem' }}>
                            {scannedAssets.length} Items
                        </span>
                    </div>

                    <div style={{ flex: 1, overflowY: 'auto', padding: '1rem' }}>
                        {scannedAssets.length === 0 ? (
                            <div style={{ height: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', opacity: 0.5 }}>
                                <Scan size={48} style={{ marginBottom: '1rem' }} />
                                <p>Queue is empty. Start scanning.</p>
                            </div>
                        ) : (
                            scannedAssets.map(a => (
                                <div key={a.id} className="glass" style={{ padding: '0.75rem 1rem', borderRadius: '0.75rem', marginBottom: '0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', border: '1px solid rgba(255,255,255,0.05)' }}>
                                    <div>
                                        <div style={{ fontWeight: 700, fontSize: '0.875rem' }}>{a.asset_tag}</div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{a.serial_number} • <span style={{ textTransform: 'uppercase', color: 'var(--primary)' }}>{a.status}</span></div>
                                    </div>
                                    <button onClick={() => removeFromBatch(a.id)} style={{ color: 'var(--text-muted)', background: 'transparent' }} onMouseOver={(e) => e.target.style.color='var(--error)'} onMouseOut={(e) => e.target.style.color='var(--text-muted)'}>
                                        <Trash2 size={16} />
                                    </button>
                                </div>
                            ))
                        )}
                    </div>

                    <div style={{ padding: '1rem', borderTop: '1px solid var(--border)' }}>
                        <button 
                            className="glass" 
                            style={{ width: '100%', padding: '0.75rem', fontSize: '0.75rem', color: 'var(--text-muted)' }} 
                            onClick={() => setScannedAssets([])}
                            disabled={scannedAssets.length === 0}
                        >
                            Clear Batch
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const KioskAction = ({ icon: Icon, label, desc, color, onClick, disabled }) => (
    <button
        onClick={onClick}
        disabled={disabled}
        style={{
            display: 'flex',
            alignItems: 'center',
            gap: '1.25rem',
            padding: '1.25rem',
            borderRadius: '1rem',
            background: disabled ? 'rgba(255,255,255,0.02)' : `${color}10`,
            border: `1px solid ${disabled ? 'transparent' : color}30`,
            opacity: disabled ? 0.3 : 1,
            textAlign: 'left',
            transition: 'all 0.2s',
            cursor: disabled ? 'not-allowed' : 'pointer'
        }}
    >
        <div style={{ background: disabled ? 'var(--surface)' : color, padding: '0.75rem', borderRadius: '0.75rem' }}>
            <Icon size={24} color="white" />
        </div>
        <div style={{ flex: 1 }}>
            <div style={{ fontWeight: 800, fontSize: '1rem', color: disabled ? 'var(--text-muted)' : 'var(--text)' }}>{label}</div>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{desc}</div>
        </div>
    </button>
);

export default WarehouseKiosk;
