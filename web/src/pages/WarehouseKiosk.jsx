import React, { useState } from 'react';
import axios from 'axios';
import { Scan, ArrowUpRight, ArrowDownLeft, Box, Search, CheckCircle2, AlertTriangle, ShieldCheck } from 'lucide-react';

const WarehouseKiosk = () => {
    const [scanTerm, setScanTerm] = useState('');
    const [asset, setAsset] = useState(null);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState(null);

    const handleScan = async (e) => {
        e.preventDefault();
        setLoading(true);
        setMessage(null);
        try {
            // Find asset by tag/serial
            const res = await axios.get('/v1/inventory/assets');
            const found = res.data.find(a => a.asset_tag === scanTerm || a.serial_number === scanTerm);
            if (found) {
                setAsset(found);
            } else {
                setMessage({ type: 'error', text: 'Asset not found in registry.' });
                setAsset(null);
            }
        } catch (err) {
            setMessage({ type: 'error', text: 'Search failed.' });
        } finally {
            setLoading(false);
        }
    };

    const updateStatus = async (newStatus) => {
        try {
            await axios.patch(`/v1/inventory/assets/${asset.id}/status`, { status: newStatus });
            setMessage({ type: 'success', text: `Asset status updated to ${newStatus}.` });
            setAsset({ ...asset, status: newStatus });
        } catch (err) {
            setMessage({ type: 'error', text: 'Update failed.' });
        }
    };

    return (
        <div style={{ padding: '4rem 2rem', maxWidth: '600px', margin: '0 auto', textAlign: 'center' }}>
            <header style={{ marginBottom: '3rem' }}>
                <div style={{ display: 'inline-flex', background: 'var(--primary)', padding: '1rem', borderRadius: '1rem', marginBottom: '1.5rem' }}>
                    <Scan size={32} color="white" />
                </div>
                <h1 style={{ fontSize: '2.5rem', fontWeight: 900 }}>Warehouse Kiosk</h1>
                <p style={{ color: 'var(--text-muted)' }}>Scan asset tag or serial to process movement.</p>
            </header>

            <form onSubmit={handleScan} style={{ marginBottom: '3rem' }}>
                <div style={{ position: 'relative' }}>
                    <input
                        type="text"
                        placeholder="Scan ID or enter tag..."
                        className="glass"
                        style={{ width: '100%', padding: '1.5rem 1.5rem 1.5rem 3.5rem', borderRadius: '1rem', fontSize: '1.25rem', color: 'white' }}
                        value={scanTerm}
                        onChange={(e) => setScanTerm(e.target.value)}
                        autoFocus
                    />
                    <Search size={24} style={{ position: 'absolute', left: '1.25rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.3 }} />
                </div>
                {loading && <p style={{ marginTop: '1rem', color: 'var(--primary)' }}>Scanning fleet database...</p>}
            </form>

            {message && (
                <div className="glass" style={{
                    padding: '1rem',
                    borderRadius: '0.75rem',
                    marginBottom: '2rem',
                    background: message.type === 'error' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                    color: message.type === 'error' ? 'var(--error)' : 'var(--success)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.5rem'
                }}>
                    {message.type === 'error' ? <AlertTriangle size={18} /> : <CheckCircle2 size={18} />}
                    {message.text}
                </div>
            )}

            {asset && (
                <div className="glass animate-in zoom-in duration-300" style={{ padding: '2rem', borderRadius: '1.5rem', textAlign: 'left' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.5rem' }}>
                        <div>
                            <h2 style={{ fontSize: '1.5rem', fontWeight: 800 }}>{asset.asset_tag}</h2>
                            <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Serial: {asset.serial_number}</p>
                        </div>
                        <StatusBadge status={asset.status} />
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <KioskAction
                            icon={ArrowUpRight}
                            label="Check-out"
                            desc="Assign to client"
                            color="var(--primary)"
                            onClick={() => updateStatus('deployed')}
                            disabled={asset.status !== 'available'}
                        />
                        <KioskAction
                            icon={ArrowDownLeft}
                            label="Check-in"
                            desc="Register return"
                            color="var(--success)"
                            onClick={() => updateStatus('available')}
                            disabled={asset.status === 'available'}
                        />
                    </div>

                    <button
                        className="glass"
                        style={{ width: '100%', marginTop: '1rem', padding: '1rem', borderRadius: '1rem', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '0.5rem', fontSize: '0.875rem' }}
                        onClick={() => updateStatus('maintenance')}
                    >
                        <ShieldCheck size={18} /> Send to Maintenance / QC
                    </button>
                </div>
            )}
        </div>
    );
};

const KioskAction = ({ icon: Icon, label, desc, color, onClick, disabled }) => (
    <button
        onClick={onClick}
        disabled={disabled}
        style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '0.5rem',
            padding: '1.5rem',
            borderRadius: '1rem',
            background: disabled ? 'rgba(255,255,255,0.02)' : `${color}15`,
            border: `1px solid ${disabled ? 'transparent' : color}40`,
            opacity: disabled ? 0.3 : 1
        }}
    >
        <div style={{ background: disabled ? 'var(--surface)' : color, padding: '0.75rem', borderRadius: '50%' }}>
            <Icon size={24} color="white" />
        </div>
        <div style={{ fontWeight: 800, color: disabled ? 'var(--text-muted)' : 'var(--text)' }}>{label}</div>
        <div style={{ fontSize: '0.65rem', color: 'var(--text-muted)' }}>{desc}</div>
    </button>
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
        <div style={{ background: `${color}20`, color, padding: '0.25rem 0.75rem', borderRadius: '2rem', fontSize: '0.75rem', fontWeight: 800 }}>
            {status.toUpperCase()}
        </div>
    );
};

export default WarehouseKiosk;
