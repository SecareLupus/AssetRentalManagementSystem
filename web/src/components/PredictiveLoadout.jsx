import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../components/Shared';

const PredictiveLoadout = ({ historicalShowId, onPredictionComplete }) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [predictedRings, setPredictedRings] = useState([]);

    const handlePredict = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(`/api/v1/seasons/predict/${historicalShowId}`);
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
        <Card className="mt-4">
            <CardHeader>
                <CardTitle>Predictive Equipment Loadout</CardTitle>
            </CardHeader>
            <CardContent>
                {predictedRings.length === 0 ? (
                    <div className="flex justify-between items-center">
                        <p className="text-gray-400 text-sm">Predict next year's ring loadouts based on historical data.</p>
                        <button
                            onClick={handlePredict}
                            disabled={loading || !historicalShowId}
                            className="bg-purple-600 hover:bg-purple-500 text-white px-4 py-2 rounded text-sm disabled:opacity-50"
                        >
                            {loading ? 'Predicting...' : 'Generate Prediction'}
                        </button>
                    </div>
                ) : (
                    <div className="space-y-6">
                        {predictedRings.map((ring, rIndex) => (
                            <div key={rIndex} className="bg-gray-800 p-4 rounded border border-gray-700">
                                <h4 className="font-semibold text-lg text-white mb-2">{ring.ring?.name || `Ring ${ring.ring_id}`}</h4>
                                {ring.loadout_items?.map((item, iIndex) => (
                                    <div key={iIndex} className="flex justify-between items-center border-b border-gray-700 py-2">
                                        <span className="text-sm text-gray-300">Item Type ID: {item.item_type_id}</span>
                                        <div className="flex items-center gap-2">
                                            <label className="text-xs text-gray-400">Qty:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                className="bg-gray-900 border border-gray-600 rounded px-2 py-1 text-white w-20 text-right"
                                                value={item.quantity}
                                                onChange={(e) => handleQuantityChange(rIndex, iIndex, e.target.value)}
                                            />
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ))}

                        {error && <p className="text-red-400 text-sm mt-2">{error}</p>}

                        <div className="flex justify-end gap-2 mt-4">
                            <button
                                onClick={() => setPredictedRings([])}
                                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded text-sm"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleConfirm}
                                className="px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded text-sm"
                            >
                                Confirm Loadout
                            </button>
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
};

export default PredictiveLoadout;
