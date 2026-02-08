import React, { useState, useEffect } from 'react';
import { Modal } from '../Shared';
import axios from 'axios';
import { Search, Zap, CheckCircle, AlertTriangle, ChevronRight, Fingerprint } from 'lucide-react';

const TARGET_FIELDS = {
    item_type: ['code', 'name', 'kind', 'is_active'],
    asset: ['item_type_id', 'asset_tag', 'serial_number', 'status', 'place_id'],
    company: ['name', 'legal_name', 'description'],
    person: ['given_name', 'family_name', 'company_id'],
    place: ['name', 'description', 'category', 'is_internal']
};

const IngestMappingModal = ({ isOpen, onClose, source, onSave }) => {
    const [mappings, setMappings] = useState([]);
    const [previewData, setPreviewData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        if (source && isOpen) {
            setMappings(source.mappings || []);
            handleDiscovery();
        }
    }, [source, isOpen]);

    const handleDiscovery = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await axios.post(`/v1/admin/ingest/sources/${source.id}/preview`);
            setPreviewData(res.data);

            // Auto-detect mappings if none exist
            if ((!source.mappings || source.mappings.length === 0) && res.data.inferred_fields) {
                const inferred = res.data.inferred_fields.map(f => ({
                    json_path: f.path,
                    target_field: f.suggested_mapping,
                    is_identity: f.suggested_mapping === 'code' || f.suggested_mapping === 'asset_tag'
                }));
                // Filter out those that don't match target fields if we want to be strict
                setMappings(inferred.filter(m => m.target_field));
            }
        } catch (err) {
            setError("Discovery failed: " + (err.response?.data || err.message));
        } finally {
            setLoading(false);
        }
    };

    const handleAddMapping = () => {
        setMappings([...mappings, { json_path: '$', target_field: '', is_identity: false }]);
    };

    const handleRemoveMapping = (index) => {
        setMappings(mappings.filter((_, i) => i !== index));
    };

    const handleUpdateMapping = (index, field, value) => {
        const next = [...mappings];
        next[index][field] = value;
        setMappings(next);
    };

    const handleSubmit = async () => {
        try {
            await axios.put(`/v1/admin/ingest/sources/${source.id}/mappings`, { mappings });
            onSave();
            onClose();
        } catch (err) {
            alert("Failed to save mappings.");
        }
    };

    const targetModelFields = TARGET_FIELDS[source?.target_model] || [];

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={`Configure Field Mappings: ${source?.name}`}
            actions={(
                <button onClick={handleSubmit} className="btn-primary">
                    Save Mappings
                </button>
            )}
        >
            <div className="space-y-6 max-h-[70vh] overflow-y-auto pr-2">
                <div className="flex items-center justify-between">
                    <div>
                        <h4 className="font-bold flex items-center gap-2">
                            <Search size={16} />
                            Schema Discovery
                        </h4>
                        <p className="text-xs text-text-muted">Detected fields from the API response</p>
                    </div>
                    <button onClick={handleDiscovery} disabled={loading} className="text-xs text-primary hover:underline font-medium flex items-center gap-1">
                        <Zap size={12} className={loading ? 'animate-spin' : ''} />
                        Refresh Schema
                    </button>
                </div>

                {error && (
                    <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-500 text-xs flex items-center gap-2">
                        <AlertTriangle size={14} />
                        {error}
                    </div>
                )}

                <div className="grid gap-2">
                    {mappings.map((m, idx) => (
                        <div key={idx} className="flex gap-2 items-center p-3 bg-surface border border-border rounded-lg group animate-in slide-in-from-right-2" style={{ animationDelay: `${idx * 0.05}s` }}>
                            <div className="flex-1 space-y-2">
                                <div className="flex items-center justify-between">
                                    <span className="text-[10px] uppercase font-bold text-text-muted tracking-wide flex items-center gap-1">
                                        JSONPath
                                        {m.is_identity && <Fingerprint size={10} className="text-primary" />}
                                    </span>
                                    <label className="flex items-center gap-1 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={m.is_identity}
                                            onChange={e => handleUpdateMapping(idx, 'is_identity', e.target.checked)}
                                        />
                                        <span className="text-[10px]">Identity</span>
                                    </label>
                                </div>
                                <input
                                    className="glass w-full p-2 text-xs font-mono"
                                    value={m.json_path}
                                    onChange={e => handleUpdateMapping(idx, 'json_path', e.target.value)}
                                />
                            </div>
                            <ChevronRight size={14} className="text-text-muted mt-4" />
                            <div className="flex-1 space-y-2">
                                <span className="text-[10px] uppercase font-bold text-text-muted tracking-wide">Internal Field</span>
                                <select
                                    className="glass w-full p-2 text-xs"
                                    value={m.target_field}
                                    onChange={e => handleUpdateMapping(idx, 'target_field', e.target.value)}
                                >
                                    <option value="">(Ignore)</option>
                                    {targetModelFields.map(f => (
                                        <option key={f} value={f}>{f}</option>
                                    ))}
                                </select>
                            </div>
                            <button
                                onClick={() => handleRemoveMapping(idx)}
                                className="p-2 hover:bg-surface-light rounded text-text-muted hover:text-red-500 transition-colors mt-4"
                            >
                                <XCircle size={14} />
                            </button>
                        </div>
                    ))}
                </div>

                <button
                    onClick={handleAddMapping}
                    className="w-full py-3 border-2 border-dashed border-border rounded-lg text-xs font-medium text-text-muted hover:text-primary hover:border-primary transition-all flex items-center justify-center gap-2"
                >
                    <Plus size={14} />
                    Add Manual Mapping
                </button>

                {previewData && (
                    <div className="space-y-2 border-t border-border pt-4">
                        <span className="text-[10px] uppercase font-bold text-text-muted tracking-wide">Sample Response Preview</span>
                        <pre className="text-[10px] font-mono p-3 bg-surface rounded-lg overflow-x-auto border border-border">
                            {JSON.stringify(previewData.sample_items?.[0] || previewData.raw_response, null, 2)}
                        </pre>
                    </div>
                )}
            </div>
        </Modal>
    );
};

const Plus = ({ size }) => <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>;
const XCircle = ({ size }) => <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>;

export default IngestMappingModal;
