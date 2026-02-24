import React, { useState } from 'react';
import { GlassCard } from '../components/Shared';
import PredictiveLoadout from '../components/PredictiveLoadout';

const SeasonPlanner = () => {
    // Scaffold state for the new Season Planner dashboard
    const [selectedCompanyId, setSelectedCompanyId] = useState('');
    const [selectedSeasonId, setSelectedSeasonId] = useState('');
    const [selectedHistoricalShowId, setSelectedHistoricalShowId] = useState('');

    const handlePredictionComplete = async (finalLoadout) => {
        try {
            // Assume we are creating a new show here, and then applying the loadout
            const newShowId = 999; // Mock ID for demonstration

            const response = await fetch(`/v1/seasons/shows/${newShowId}/loadout`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(finalLoadout)
            });

            if (!response.ok) {
                throw new Error('Failed to save predicted loadout to the new show');
            }
            alert('Success! New show loadout has been created.');
        } catch (err) {
            alert(err.message);
        }
    };

    return (
        <div style={{ padding: '1.5rem' }}>
            <h1 style={{ fontSize: '1.5rem', fontWeight: 'bold', marginBottom: '1.5rem', color: 'white' }}>Season & Show Planner</h1>
            <p style={{ color: 'var(--text-muted)', marginBottom: '2rem' }}>
                Manage the logistical hierarchy of Show Companies, Seasons, Shows, and their associated Rings.
            </p>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
                <GlassCard>
                    <div style={{ marginBottom: '1.5rem' }}>
                        <h3 style={{ fontSize: '1.125rem', fontWeight: 600 }}>Hierarchy Setup</h3>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
                        <div>
                            <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Select Show Company</label>
                            <select
                                style={{ width: '100%', background: 'var(--surface)', border: '1px solid var(--border)', borderRadius: '0.5rem', padding: '0.5rem', color: 'white' }}
                                value={selectedCompanyId}
                                onChange={(e) => setSelectedCompanyId(e.target.value)}
                            >
                                <option value="">-- Choose Company --</option>
                                <option value="1">Alpha Show Management</option>
                                <option value="2">Beta Equestrian Events</option>
                            </select>
                        </div>

                        <div>
                            <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Select Season</label>
                            <select
                                style={{ width: '100%', background: 'var(--surface)', border: '1px solid var(--border)', borderRadius: '0.5rem', padding: '0.5rem', color: 'white' }}
                                value={selectedSeasonId}
                                onChange={(e) => setSelectedSeasonId(e.target.value)}
                                disabled={!selectedCompanyId}
                            >
                                <option value="">-- Choose Season --</option>
                                <option value="1">2026 Winter Circuit</option>
                                <option value="2">2026 Summer Series</option>
                            </select>
                        </div>

                        <div>
                            <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Historical Show Reference (For Prediction)</label>
                            <select
                                style={{ width: '100%', background: 'var(--surface)', border: '1px solid var(--border)', borderRadius: '0.5rem', padding: '0.5rem', color: 'white' }}
                                value={selectedHistoricalShowId}
                                onChange={(e) => setSelectedHistoricalShowId(e.target.value)}
                                disabled={!selectedSeasonId}
                            >
                                <option value="">-- Choose Historical Show --</option>
                                <option value="1">2025 Winter Circuit - Week 1</option>
                                <option value="2">2025 Winter Circuit - Week 2</option>
                            </select>
                        </div>
                    </div>
                </GlassCard>

                <div>
                    {selectedHistoricalShowId && (
                        <PredictiveLoadout
                            historicalShowId={selectedHistoricalShowId}
                            onPredictionComplete={handlePredictionComplete}
                        />
                    )}
                </div>
            </div>
        </div>
    );
};

export default SeasonPlanner;

