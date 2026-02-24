import React, { useState, useEffect } from 'react';
import { GlassCard } from '../components/Shared';

const PredictiveLoadout = ({ historicalShowId, onPredictionComplete }) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [predictedRings, setPredictedRings] = useState([]);

    const handlePredict = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(`/v1/seasons/predict/${historicalShowId}`);
            if (!response.ok) {
                throw new Error('Failed to fetch prediction');
            }
            const data = await response.json();
            setPredictedRings(data || []);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleQuantityChange = (ringIndex, itemIndex, newQuantity) => {
        const updated = [...predictedRings];
        updated[ringIndex].loadout_items[itemIndex].quantity = parseInt(newQuantity, 10) || 0;
        setPredictedRings(updated);
    };

    const handleConfirm = () => {
        if (onPredictionComplete) {
            onPredictionComplete(predictedRings);
        }
    };

    return (
        <GlassCard className="mt-4" style={{ marginTop: '1rem' }}>
            <div style={{ marginBottom: '1.5rem' }}>
                <h3 style={{ fontSize: '1.25rem', fontWeight: 700 }}>Predictive Equipment Loadout</h3>
            </div>
            <div>
                {predictedRings.length === 0 ? (
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Predict next year's ring loadouts based on historical data.</p>
                        <button
                            onClick={handlePredict}
                            disabled={loading || !historicalShowId}
                            className="btn-primary"
                            style={{ fontSize: '0.875rem', padding: '0.5rem 1rem' }}
                        >
                            {loading ? 'Predicting...' : 'Generate Prediction'}
                        </button>
                    </div>
                ) : (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                        {predictedRings.map((ring, rIndex) => (
                            <div key={rIndex} style={{ background: 'var(--surface)', padding: '1rem', borderRadius: '0.75rem', border: '1px solid var(--border)' }}>
                                <h4 style={{ fontWeight: 600, fontSize: '1.1rem', marginBottom: '1rem' }}>{ring.ring?.name || `Ring ${ring.ring_id}`}</h4>
                                {ring.loadout_items?.map((item, iIndex) => (
                                    <div key={iIndex} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid var(--border)', padding: '0.5rem 0' }}>
                                        <span style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>Item Type ID: {item.item_type_id}</span>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Qty:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                style={{ background: 'transparent', border: '1px solid var(--border)', borderRadius: '0.4rem', padding: '0.25rem 0.5rem', width: '60px', textAlign: 'right', color: 'white' }}
                                                value={item.quantity}
                                                onChange={(e) => handleQuantityChange(rIndex, iIndex, e.target.value)}
                                            />
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ))}

                        {error && <p style={{ color: 'var(--error)', fontSize: '0.875rem' }}>{error}</p>}

                        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
                            <button
                                onClick={() => setPredictedRings([])}
                                className="glass"
                                style={{ padding: '0.5rem 1rem' }}
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleConfirm}
                                className="btn-primary"
                                style={{ padding: '0.5rem 1rem' }}
                            >
                                Confirm Loadout
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </GlassCard>
    );
};

export default PredictiveLoadout;

