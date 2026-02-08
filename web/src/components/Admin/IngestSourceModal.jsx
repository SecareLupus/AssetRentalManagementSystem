import React, { useState, useEffect } from 'react';
import { Modal } from '../Shared';
import axios from 'axios';

const INGEST_TARGET_MODELS = [
    { value: 'item_type', label: 'Item Type (SKU)' },
    { value: 'asset', label: 'Asset (Individual)' },
    { value: 'company', label: 'Company' },
    { value: 'person', label: 'Person' },
    { value: 'place', label: 'Place (Location)' }
];

const IngestSourceModal = ({ isOpen, onClose, source, onSave }) => {
    const [formData, setFormData] = useState({
        name: '',
        target_model: 'item_type',
        api_url: '',
        auth_type: 'none',
        auth_credentials: '',
        sync_interval_seconds: 3600,
        is_active: true
    });

    useEffect(() => {
        if (source) {
            setFormData({
                ...source,
                auth_credentials: source.auth_credentials ? JSON.stringify(JSON.parse(source.auth_credentials), null, 2) : ''
            });
        } else {
            setFormData({
                name: '',
                target_model: 'item_type',
                api_url: '',
                auth_type: 'none',
                auth_credentials: '',
                sync_interval_seconds: 3600,
                is_active: true
            });
        }
    }, [source, isOpen]);

    const handleSubmit = async () => {
        try {
            const data = { ...formData };
            if (data.auth_credentials) {
                try {
                    // If it's just a token, wrap it in JSON if needed, or parse as JSON
                    // The backend expects auth_credentials to be JSON.
                    // For Bearer, we'll store it as a stringified JSON string.
                    const parsed = JSON.parse(data.auth_credentials);
                    data.auth_credentials = JSON.stringify(parsed);
                } catch {
                    // Assume it's a raw token string and wrap it
                    data.auth_credentials = JSON.stringify(data.auth_credentials);
                }
            }

            if (source?.id) {
                await axios.put(`/v1/admin/ingest/sources/${source.id}`, data);
            } else {
                await axios.post('/v1/admin/ingest/sources', data);
            }
            onSave();
            onClose();
        } catch (err) {
            alert("Failed to save source: " + (err.response?.data || err.message));
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={source ? "Edit Ingest Source" : "New Ingest Source"}
            actions={(
                <button onClick={handleSubmit} className="btn-primary">
                    {source ? "Update Source" : "Create Source"}
                </button>
            )}
        >
            <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="form-label">Source Name</label>
                        <input
                            className="glass w-full p-2"
                            placeholder="e.g. ERP Product Catalog"
                            value={formData.name}
                            onChange={e => setFormData({ ...formData, name: e.target.value })}
                        />
                    </div>
                    <div>
                        <label className="form-label">Target Model</label>
                        <select
                            className="glass w-full p-2"
                            value={formData.target_model}
                            onChange={e => setFormData({ ...formData, target_model: e.target.value })}
                        >
                            {INGEST_TARGET_MODELS.map(m => (
                                <option key={m.value} value={m.value}>{m.label}</option>
                            ))}
                        </select>
                    </div>
                </div>

                <div>
                    <label className="form-label">API URL</label>
                    <input
                        className="glass w-full p-2"
                        placeholder="https://api.example.com/v1/items"
                        value={formData.api_url}
                        onChange={e => setFormData({ ...formData, api_url: e.target.value })}
                    />
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="form-label">Auth Type</label>
                        <select
                            className="glass w-full p-2"
                            value={formData.auth_type}
                            onChange={e => setFormData({ ...formData, auth_type: e.target.value })}
                        >
                            <option value="none">None</option>
                            <option value="bearer">Bearer Token</option>
                        </select>
                    </div>
                    <div>
                        <label className="form-label">Sync Interval (seconds)</label>
                        <input
                            type="number"
                            className="glass w-full p-2"
                            value={formData.sync_interval_seconds}
                            onChange={e => setFormData({ ...formData, sync_interval_seconds: parseInt(e.target.value) })}
                        />
                    </div>
                </div>

                {formData.auth_type !== 'none' && (
                    <div>
                        <label className="form-label">Auth Credentials (JSON or Token)</label>
                        <textarea
                            className="glass w-full p-2 font-mono text-xs"
                            rows={3}
                            placeholder='e.g. "your-token-here"'
                            value={formData.auth_credentials}
                            onChange={e => setFormData({ ...formData, auth_credentials: e.target.value })}
                        />
                        <p className="text-[10px] text-text-muted mt-1">For Bearer, enter the token string directly (including quotes if raw string).</p>
                    </div>
                )}

                <label className="flex items-center gap-2 cursor-pointer mt-4">
                    <input
                        type="checkbox"
                        checked={formData.is_active}
                        onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                    />
                    <span className="text-sm font-medium">Capture data automatically on schedule</span>
                </label>
            </div>
        </Modal>
    );
};

export default IngestSourceModal;
