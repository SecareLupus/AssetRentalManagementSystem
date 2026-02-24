import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    Download,
    Scan,
    Search,
    ShipWheel,
    Box,
    CheckCircle2,
    AlertTriangle,
    ArrowLeft,
    History,
    Wrench,
    QrCode
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const ReturnProcessing = () => {
    const navigate = useNavigate();
    const [searchTerm, setSearchTerm] = useState('');
    const [activeShipment, setActiveShipment] = useState(null);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState(null);
    const [recentReturns, setRecentReturns] = useState([]);

    const handleSearch = async (e) => {
        if (e) e.preventDefault();
        setLoading(true);
        setMessage(null);
        try {
            // In a real app, we'd search by tracking number or customer
            // For now, let's try to find a shipment
            const res = await axios.get('/v1/logistics/shipments');
            const found = res.data?.find(s => s.trackingNumber === searchTerm || s.id.toString() === searchTerm);

            if (found) {
                setActiveShipment(found);
            } else {
                setMessage({ type: 'error', text: 'No matching shipment found.' });
            }
        } catch (err) {
            setMessage({ type: 'error', text: 'Failed to search shipments.' });
        } finally {
            setLoading(false);
        }
    };

    const handleReceiveAsset = async (asset) => {
        try {
            // 1. Create ReturnAction (in a real app, this would be a bulk endpoint)
            // 2. Navigate to InspectionRunner
            navigate(`/tech/inspect/${asset.id}`);
        } catch (err) {
            console.error("Failed to initiate return", err);
        }
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <header style={{ textAlign: 'center', marginBottom: '4rem' }}>
                <div style={{ display: 'inline-flex', background: 'var(--success)', padding: '1.25rem', borderRadius: '1.25rem', marginBottom: '1.5rem', boxShadow: '0 0 20px rgba(16, 185, 129, 0.3)' }}>
                    <Download size={32} color="white" />
                </div>
                <h1 style={{ fontSize: '3rem', fontWeight: 900, letterSpacing: '-0.025em' }}>Inbound Return Center</h1>
                <p style={{ color: 'var(--text-muted)', fontSize: '1.1rem' }}>Scan tracking numbers or asset tags to process returned equipment and trigger inspections.</p>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: activeShipment ? '1fr 1fr' : '1fr', gap: '3rem', transition: 'all 0.4s ease' }}>
                {/* Search / Scan Section */}
                <div>
                    <form onSubmit={handleSearch} style={{ marginBottom: '2.5rem' }}>
                        <div style={{ position: 'relative' }}>
                            <input
                                type="text"
                                placeholder="Scan Tracking # or Shipment ID..."
                                className="glass"
                                style={{
                                    width: '100%',
                                    padding: '1.75rem 1.75rem 1.75rem 4rem',
                                    borderRadius: '1.25rem',
                                    fontSize: '1.5rem',
                                    color: 'white',
                                    fontWeight: 600,
                                    border: '1px solid rgba(255,255,255,0.1)'
                                }}
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                autoFocus
                            />
                            <QrCode size={28} style={{ position: 'absolute', left: '1.5rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.4 }} />
                        </div>
                    </form>

                    <div className="glass" style={{ padding: '2rem', borderRadius: '1.5rem' }}>
                        <h3 style={{ fontWeight: 800, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                            <History size={20} color="var(--primary)" /> Recent Activity
                        </h3>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            <div style={{ padding: '1rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', opacity: 0.5 }}>
                                <span style={{ fontSize: '0.85rem' }}>No recent returns processed in this session.</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Shipment / Asset List Section */}
                {activeShipment && (
                    <div className="glass" style={{ borderRadius: '2rem', display: 'flex', flexDirection: 'column', border: '2px solid var(--primary)', overflow: 'hidden' }}>
                        <div style={{ padding: '1.5rem', background: 'rgba(99, 102, 241, 0.1)', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <div>
                                <h3 style={{ fontWeight: 800, fontSize: '1.25rem' }}>Shipment Found</h3>
                                <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>ID: {activeShipment.id} • Carrier: {activeShipment.carrier || 'Unknown'}</p>
                            </div>
                            <button onClick={() => setActiveShipment(null)} style={{ background: 'transparent', color: 'var(--text-muted)' }}>
                                <ArrowLeft size={20} />
                            </button>
                        </div>

                        <div style={{ flex: 1, padding: '1.5rem', overflowY: 'auto' }}>
                            <h4 style={{ fontSize: '0.9rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--text-muted)', marginBottom: '1.5rem' }}>Expected Assets</h4>

                            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                                {/* Mocking some assets for the UI demonstration */}
                                <AssetReturnItem tag="SGL-4001" name="Video Switcher Pro" onReturn={() => navigate('/tech/inspect/1')} />
                                <AssetReturnItem tag="SGL-4002" name="Video Switcher Pro" onReturn={() => navigate('/tech/inspect/2')} />
                                <AssetReturnItem tag="SGL-9023" name="Router Array B" onReturn={() => navigate('/tech/inspect/3')} />
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {message && (
                <div style={{
                    position: 'fixed',
                    bottom: '2rem',
                    left: '50%',
                    transform: 'translateX(-50%)',
                    padding: '1.25rem 2rem',
                    borderRadius: '1rem',
                    background: message.type === 'error' ? 'rgba(239, 68, 68, 0.9)' : 'rgba(16, 185, 129, 0.9)',
                    backdropFilter: 'blur(10px)',
                    color: 'white',
                    fontWeight: 700,
                    display: 'flex', alignItems: 'center', gap: '1rem',
                    boxShadow: '0 10px 30px rgba(0,0,0,0.3)',
                    zIndex: 1000
                }}>
                    {message.type === 'error' ? <AlertTriangle size={20} /> : <CheckCircle2 size={20} />}
                    {message.text}
                </div>
            )}
        </div>
    );
};

const AssetReturnItem = ({ tag, name, onReturn }) => (
    <div className="glass hover-lift" style={{
        padding: '1.25rem',
        borderRadius: '1.25rem',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        border: '1px solid rgba(255,255,255,0.05)'
    }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div style={{ background: 'var(--surface)', padding: '0.75rem', borderRadius: '0.75rem' }}>
                <Box size={20} color="var(--primary)" />
            </div>
            <div>
                <div style={{ fontWeight: 800, fontSize: '1rem' }}>{tag}</div>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{name}</div>
            </div>
        </div>
        <button
            onClick={onReturn}
            style={{
                padding: '0.6rem 1rem',
                borderRadius: '0.75rem',
                background: 'var(--success)',
                color: 'white',
                fontWeight: 700,
                fontSize: '0.85rem',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem'
            }}
        >
            Receive & QC <Wrench size={16} />
        </button>
    </div>
);

export default ReturnProcessing;
