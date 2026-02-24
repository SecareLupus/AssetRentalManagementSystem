import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    X,
    Search,
    Scan,
    Package,
    CheckCircle2,
    AlertCircle,
    ArrowRight,
    ChevronRight,
    Loader2
} from 'lucide-react';
import { GlassCard, StatusBadge } from './Shared';

const ShipmentAllocationUI = ({ shipmentId, onClose, onAllocationComplete }) => {
    const [shipment, setShipment] = useState(null);
    const [delivery, setDelivery] = useState(null);
    const [reservation, setReservation] = useState(null);
    const [demands, setDemands] = useState([]);
    const [assets, setAssets] = useState([]); // Currently allocated assets

    const [searchTerm, setSearchTerm] = useState('');
    const [availableAssets, setAvailableAssets] = useState([]);
    const [loading, setLoading] = useState(true);
    const [allocating, setAllocating] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchShipmentData();
    }, [shipmentId]);

    const fetchShipmentData = async () => {
        setLoading(true);
        try {
            const shipRes = await axios.get(`/v1/logistics/shipments/${shipmentId}`);
            setShipment(shipRes.data);

            if (shipRes.data.scheduledDeliveryId) {
                const delRes = await axios.get(`/v1/logistics/deliveries/${shipRes.data.scheduledDeliveryId}`);
                setDelivery(delRes.data);

                // Fetch reservation and its demands
                const resRes = await axios.get(`/v1/logistics/reservations/${delRes.data.eventId}`);
                setReservation(resRes.data);

                // Fetch demands
                const demandsRes = await axios.get(`/v1/entities/events/${delRes.data.eventId}/demands`);
                setDemands(demandsRes.data || []);

                // Fetch already allocated assets for this shipment (via CheckOutActions)
                const actionsRes = await axios.get(`/v1/logistics/reservations/${delRes.data.eventId}/fulfillment`);
                // Filter actions for this shipment
                const allocated = (actionsRes.data?.details || []).filter(d => d.shipment_id === parseInt(shipmentId));
                setAssets(allocated);
            }
        } catch (err) {
            setError("Failed to load shipment details");
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const searchAvailableAssets = async (term) => {
        if (!term || term.length < 2) {
            setAvailableAssets([]);
            return;
        }
        try {
            // Search for available assets (optionally filter by item types in demands)
            const res = await axios.get(`/v1/inventory/assets?asset_tag=${term}&status=available`);
            setAvailableAssets(res.data || []);
        } catch (err) {
            console.error("Search failed", err);
        }
    };

    useEffect(() => {
        const delaySearch = setTimeout(() => {
            searchAvailableAssets(searchTerm);
        }, 300);
        return () => clearTimeout(delaySearch);
    }, [searchTerm]);

    const handleAllocate = async (assetId) => {
        setAllocating(true);
        setError(null);
        try {
            await axios.post(`/v1/logistics/shipments/${shipmentId}/allocate`, {
                asset_ids: [parseInt(assetId)]
            });
            setSearchTerm('');
            setAvailableAssets([]);
            fetchShipmentData(); // Refresh fulfillment
        } catch (err) {
            setError(err.response?.data || "Allocation failed");
        } finally {
            setAllocating(false);
        }
    };

    if (loading) return (
        <div className="modal-overlay">
            <GlassCard style={{ width: '400px', padding: '3rem', textAlign: 'center' }}>
                <Loader2 className="animate-spin" style={{ margin: '0 auto 1rem' }} size={32} />
                <p>Loading Shipment Context...</p>
            </GlassCard>
        </div>
    );

    return (
        <div className="modal-overlay">
            <GlassCard style={{ width: '900px', maxHeight: '90vh', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
                <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <div>
                        <h2 style={{ fontSize: '1.5rem', fontWeight: 800 }}>Asset Allocation</h2>
                        <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
                            Fulfilling Shipment #{shipmentId} for {reservation?.reservationName || 'Unknown Event'}
                        </p>
                    </div>
                    <button onClick={onClose} className="glass" style={{ padding: '0.5rem' }}><X size={20} /></button>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 350px', gap: '1px', background: 'var(--border)', flex: 1, overflow: 'hidden' }}>
                    {/* Left: Search and Available Assets */}
                    <div style={{ background: 'var(--surface)', padding: '1.5rem', overflowY: 'auto' }}>
                        <div style={{ marginBottom: '2rem' }}>
                            <label style={{ display: 'block', fontSize: '0.75rem', fontWeight: 900, marginBottom: '0.75rem', color: 'var(--text-muted)' }}>
                                SCAN OR SEARCH ASSET TAG
                            </label>
                            <div className="glass" style={{ display: 'flex', alignItems: 'center', padding: '0.75rem 1rem', borderRadius: '0.75rem', gap: '0.75rem' }}>
                                <Scan size={20} color="var(--primary)" />
                                <input
                                    autoFocus
                                    placeholder="Enter asset tag..."
                                    style={{ background: 'transparent', border: 'none', color: 'white', flex: 1, outline: 'none', fontSize: '1rem' }}
                                    value={searchTerm}
                                    onChange={e => setSearchTerm(e.target.value)}
                                />
                                {searchTerm && <button onClick={() => setSearchTerm('')}><X size={16} /></button>}
                            </div>
                        </div>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                            {availableAssets.length > 0 ? (
                                availableAssets.map(asset => (
                                    <div key={asset.id} className="glass" style={{ padding: '1rem', borderRadius: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <div>
                                            <div style={{ fontWeight: 700 }}>{asset.asset_tag}</div>
                                            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{asset.item_type_name}</div>
                                        </div>
                                        <button
                                            className="btn-primary"
                                            style={{ padding: '0.4rem 0.8rem', fontSize: '0.75rem' }}
                                            onClick={() => handleAllocate(asset.id)}
                                            disabled={allocating}
                                        >
                                            Allocate <ArrowRight size={14} />
                                        </button>
                                    </div>
                                ))
                            ) : searchTerm.length >= 2 ? (
                                <div style={{ textAlign: 'center', padding: '2rem', opacity: 0.5 }}>
                                    <p>No available assets found matching "{searchTerm}"</p>
                                </div>
                            ) : (
                                <div style={{ textAlign: 'center', padding: '4rem', opacity: 0.3 }}>
                                    <Package size={48} style={{ margin: '0 auto 1rem' }} />
                                    <p>Scan an asset tag to begin allocation</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Right: Shipment Requirements / Fulfillment */}
                    <aside style={{ background: 'rgba(0,0,0,0.1)', padding: '1.5rem', overflowY: 'auto' }}>
                        <h3 style={{ fontSize: '1rem', fontWeight: 800, marginBottom: '1.5rem' }}>Shipment Requirements</h3>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                            {demands.map(demand => {
                                const allocatedCount = assets.filter(a => a.item_type_id === demand.item_type_id).length;
                                const isFulfilled = allocatedCount >= demand.quantity;

                                return (
                                    <div key={demand.id}>
                                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                                            <span style={{ fontSize: '0.85rem', fontWeight: 600 }}>{demand.item_type_name || `Type ${demand.item_type_id}`}</span>
                                            <span style={{ fontSize: '0.85rem', fontWeight: 800, color: isFulfilled ? 'var(--success)' : 'white' }}>
                                                {allocatedCount} / {demand.quantity}
                                            </span>
                                        </div>
                                        <div style={{ height: '6px', background: 'rgba(255,255,255,0.05)', borderRadius: '3px', overflow: 'hidden' }}>
                                            <div style={{
                                                width: `${Math.min(100, (allocatedCount / demand.quantity) * 100)}%`,
                                                height: '100%',
                                                background: isFulfilled ? 'var(--success)' : 'var(--primary)',
                                                transition: 'width 0.3s'
                                            }} />
                                        </div>
                                    </div>
                                );
                            })}
                        </div>

                        <hr style={{ margin: '2rem 0', border: 'none', borderTop: '1px solid var(--border)' }} />

                        <h3 style={{ fontSize: '0.85rem', fontWeight: 800, marginBottom: '1rem', color: 'var(--text-muted)' }}>ALLOCATED ASSETS</h3>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                            {assets.map(asset => (
                                <div key={asset.id} style={{ fontSize: '0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0.5rem', background: 'rgba(255,255,255,0.02)', borderRadius: '0.5rem' }}>
                                    <span>{asset.asset_tag}</span>
                                    <CheckCircle2 size={14} color="var(--success)" />
                                </div>
                            ))}
                            {assets.length === 0 && <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', textAlign: 'center' }}>No assets allocated yet.</p>}
                        </div>

                        {error && (
                            <div style={{ marginTop: '2rem', padding: '0.75rem', borderRadius: '0.5rem', background: 'rgba(239, 68, 68, 0.1)', color: 'var(--error)', fontSize: '0.75rem', display: 'flex', gap: '0.5rem', alignItems: 'start' }}>
                                <AlertCircle size={14} style={{ marginTop: '0.1rem' }} />
                                <span>{error}</span>
                            </div>
                        )}
                    </aside>
                </div>

                <div style={{ padding: '1.25rem', borderTop: '1px solid var(--border)', display: 'flex', justifyContent: 'flex-end', background: 'var(--surface)' }}>
                    <button className="btn-primary" onClick={onClose} style={{ padding: '0.6rem 2rem' }}>
                        Done
                    </button>
                </div>
            </GlassCard>
        </div>
    );
};

export default ShipmentAllocationUI;
