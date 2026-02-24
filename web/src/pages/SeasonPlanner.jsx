import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../components/Shared';
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

            const response = await fetch(`/api/v1/seasons/shows/${newShowId}/loadout`, {
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
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6 text-white">Season & Show Planner</h1>
            <p className="text-gray-400 mb-8">
                Manage the logistical hierarchy of Show Companies, Seasons, Shows, and their associated Rings.
            </p>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <Card>
                    <CardHeader>
                        <CardTitle>Hierarchy Setup</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-1">Select Show Company</label>
                            <select
                                className="w-full bg-gray-900 border border-gray-700 rounded p-2 text-white"
                                value={selectedCompanyId}
                                onChange={(e) => setSelectedCompanyId(e.target.value)}
                            >
                                <option value="">-- Choose Company --</option>
                                <option value="1">Alpha Show Management</option>
                                <option value="2">Beta Equestrian Events</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-1">Select Season</label>
                            <select
                                className="w-full bg-gray-900 border border-gray-700 rounded p-2 text-white"
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
                            <label className="block text-sm font-medium text-gray-300 mb-1">Historical Show Reference (For Prediction)</label>
                            <select
                                className="w-full bg-gray-900 border border-gray-700 rounded p-2 text-white"
                                value={selectedHistoricalShowId}
                                onChange={(e) => setSelectedHistoricalShowId(e.target.value)}
                                disabled={!selectedSeasonId}
                            >
                                <option value="">-- Choose Historical Show --</option>
                                <option value="1">2025 Winter Circuit - Week 1</option>
                                <option value="2">2025 Winter Circuit - Week 2</option>
                            </select>
                        </div>
                    </CardContent>
                </Card>

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
