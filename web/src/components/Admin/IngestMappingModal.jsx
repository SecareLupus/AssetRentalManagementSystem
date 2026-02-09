import React, { useState, useEffect } from 'react';
import { Modal, GlassCard } from '../Shared';
import axios from 'axios';
import {
    Search, Zap, CheckCircle, AlertTriangle, ChevronRight,
    Fingerprint, Plus, Trash2, Database, Target, Activity,
    Layout, Globe
} from 'lucide-react';

const TARGET_MODELS = [
    { value: 'item_type', label: 'Item Type' },
    { value: 'asset', label: 'Asset' },
    { value: 'company', label: 'Company' },
    { value: 'person', label: 'Person' },
    { value: 'place', label: 'Place' }
];

const TARGET_FIELDS = {
    item_type: ['code', 'name', 'kind', 'is_active'],
    asset: ['item_type_id', 'asset_tag', 'serial_number', 'status', 'place_id'],
    company: ['name', 'legal_name', 'description'],
    person: ['given_name', 'family_name', 'company_id'],
    place: ['name', 'description', 'category', 'is_internal']
};

const IngestMappingModal = ({ isOpen, onClose, source, onSave }) => {
    const [selectedEndpointId, setSelectedEndpointId] = useState(null);
    const [mappings, setMappings] = useState([]);
    const [previewData, setPreviewData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        if (source && isOpen) {
            if (source.endpoints?.length > 0) {
                const firstId = source.endpoints[0].id;
                setSelectedEndpointId(firstId);
                loadEndpointMappings(firstId);
            }
        }
    }, [source, isOpen]);

    const loadEndpointMappings = (endpointId) => {
        const ep = source.endpoints.find(e => e.id === endpointId);
        if (ep) {
            setMappings(ep.mappings || []);
            setPreviewData(null);
            setError(null);
        }
    };

    const handleDiscovery = async () => {
        if (!selectedEndpointId) return;
        setLoading(true);
        setError(null);
        try {
            const res = await axios.get(`/v1/admin/ingest/endpoints/${selectedEndpointId}/discovery`);
            setPreviewData(res.data);

            // If discovery suggested an items_path and the endpoint doesn't have a specific one (or is default),
            // we should update it. Note: In a real app, we'd probably want to save this back to the endpoint.
            if (res.data.items_path && mappings.length === 0) {
                // Just for UI feedback
                console.log("Discovery suggested Items Path:", res.data.items_path);
            }

            // Auto-detect mappings if none exist
            if (mappings.length === 0 && res.data.inferred_fields) {
                const inferred = res.data.inferred_fields.map(f => ({
                    json_path: f.path,
                    target_model: f.suggested_model || 'asset',
                    target_field: f.suggest_mapping || '', // backend uses suggest_mapping
                    is_identity: f.is_identity
                }));
                // Set all inferred fields
                setMappings(inferred);
            }
        } catch (err) {
            setError("Discovery failed: " + (err.response?.data || err.message));
        } finally {
            setLoading(false);
        }
    };

    const handleAddMapping = () => {
        setMappings([...mappings, { json_path: '$', target_model: 'asset', target_field: '', is_identity: false }]);
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
        if (!selectedEndpointId) return;
        try {
            await axios.post(`/v1/admin/ingest/endpoints/${selectedEndpointId}/mappings`, { mappings });

            // Update local source object to reflect changes if needed
            if (source.endpoints) {
                const ep = source.endpoints.find(e => e.id === selectedEndpointId);
                if (ep) ep.mappings = mappings;
            }

            // If we have more endpoints, maybe stay open? 
            // For now, let's close or show success
            onSave();
            // Optional: alert("Mappings saved for this endpoint");
        } catch (err) {
            alert("Failed to save mappings: " + (err.response?.data || err.message));
        }
    };

    if (!source || !source.endpoints || source.endpoints.length === 0) {
        return (
            <Modal isOpen={isOpen} onClose={onClose} title="No Endpoints Defined">
                <div className="p-8 text-center text-text-muted italic">
                    You must define at least one endpoint in the Source settings before configuring mappings.
                </div>
            </Modal>
        );
    }

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="xl"
            title={`Configure Mappings: ${source.name}`}
            actions={(
                <div className="flex gap-2">
                    <button onClick={onClose} className="btn-secondary">Close</button>
                    <button onClick={handleSubmit} className="btn-primary flex items-center gap-2">
                        <CheckCircle size={16} /> Save Mappings
                    </button>
                </div>
            )}
        >
            <div className="flex gap-6 min-h-[500px]">
                {/* Left Sidebar: Endpoints */}
                <div className="w-64 border-r border-white/5 pr-4 space-y-2">
                    <label className="text-[10px] uppercase font-bold text-text-muted mb-2 block tracking-widest">Select Endpoint</label>
                    {source.endpoints.map(ep => (
                        <button
                            key={ep.id}
                            onClick={() => { setSelectedEndpointId(ep.id); loadEndpointMappings(ep.id); }}
                            className={`w-full text-left p-3 rounded-xl transition-all border group ${selectedEndpointId === ep.id
                                    ? 'bg-primary/10 border-primary/20 ring-1 ring-primary/20'
                                    : 'bg-white/5 border-transparent hover:border-white/10'
                                }`}
                        >
                            <div className="flex items-center gap-2 mb-1">
                                <span className={`text-[8px] font-bold px-1.5 py-0.5 rounded ${ep.method === 'GET' ? 'bg-blue-500/20 text-blue-400' : 'bg-green-500/20 text-green-400'
                                    }`}>
                                    {ep.method}
                                </span>
                                <span className={`text-xs font-bold truncate ${selectedEndpointId === ep.id ? 'text-primary' : 'text-text'}`}>
                                    {ep.path}
                                </span>
                            </div>
                            <div className="text-[9px] text-text-muted flex items-center gap-1">
                                <Activity size={10} /> {ep.mappings?.length || 0} fields mapped
                            </div>
                        </button>
                    ))}
                </div>

                {/* Main: Mapping UI */}
                <div className="flex-1 space-y-6 max-h-[70vh] overflow-y-auto pr-2 custom-scrollbar">
                    <div className="flex items-center justify-between">
                        <div>
                            <h4 className="font-bold flex items-center gap-2 text-primary">
                                <Database size={16} /> Data Field Mapping
                            </h4>
                            <p className="text-xs text-text-muted">Transform external JSON to internal items</p>
                        </div>
                        <button
                            onClick={handleDiscovery}
                            disabled={loading || !selectedEndpointId}
                            className="bg-white/5 hover:bg-white/10 px-3 py-1.5 rounded-lg text-xs font-bold flex items-center gap-2 transition-all"
                        >
                            <Zap size={14} className={loading ? 'animate-spin text-yellow-500' : 'text-yellow-500'} />
                            Run Discovery
                        </button>
                    </div>

                    {error && (
                        <div className="p-4 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-xs flex items-center gap-3">
                            <AlertTriangle size={18} />
                            {error}
                        </div>
                    )}

                    <div className="space-y-3">
                        {mappings.map((m, idx) => (
                            <GlassCard key={idx} className="p-4 relative group hover:ring-1 hover:ring-primary/30 transition-all">
                                <button
                                    onClick={() => handleRemoveMapping(idx)}
                                    className="absolute top-4 right-4 text-text-muted hover:text-red-500 opacity-0 group-hover:opacity-100 transition-all"
                                >
                                    <Trash2 size={14} />
                                </button>

                                <div className="grid grid-cols-12 gap-4">
                                    <div className="col-span-12 md:col-span-5 space-y-2">
                                        <div className="flex items-center justify-between">
                                            <span className="text-[10px] uppercase font-bold text-text-muted flex items-center gap-1">
                                                <Layout size={10} /> JSONPath Source
                                            </span>
                                            <label className="flex items-center gap-2 cursor-pointer group/id">
                                                <input
                                                    type="checkbox"
                                                    className="w-3 h-3 accent-primary"
                                                    checked={m.is_identity}
                                                    onChange={e => handleUpdateMapping(idx, 'is_identity', e.target.checked)}
                                                />
                                                <span className={`text-[9px] font-bold uppercase ${m.is_identity ? 'text-primary' : 'text-text-muted'}`}>ID Key</span>
                                                <Fingerprint size={10} className={m.is_identity ? 'text-primary' : 'text-text-muted'} />
                                            </label>
                                        </div>
                                        <input
                                            className="glass w-full p-2 text-xs font-mono"
                                            value={m.json_path}
                                            placeholder="$.data.id"
                                            onChange={e => handleUpdateMapping(idx, 'json_path', e.target.value)}
                                        />
                                    </div>

                                    <div className="hidden md:flex col-span-1 items-center justify-center pt-6 text-text-muted">
                                        <ChevronRight size={20} />
                                    </div>

                                    <div className="col-span-12 md:col-span-3 space-y-2">
                                        <span className="text-[10px] uppercase font-bold text-text-muted flex items-center gap-1">
                                            <Target size={10} /> Model
                                        </span>
                                        <select
                                            className="glass w-full p-2 text-xs font-bold"
                                            value={m.target_model}
                                            onChange={e => handleUpdateMapping(idx, 'target_model', e.target.value)}
                                        >
                                            {TARGET_MODELS.map(tm => (
                                                <option key={tm.value} value={tm.value}>{tm.label}</option>
                                            ))}
                                        </select>
                                    </div>

                                    <div className="col-span-12 md:col-span-3 space-y-2">
                                        <span className="text-[10px] uppercase font-bold text-text-muted">Field</span>
                                        <select
                                            className="glass w-full p-2 text-xs font-bold text-primary"
                                            value={m.target_field}
                                            onChange={e => handleUpdateMapping(idx, 'target_field', e.target.value)}
                                        >
                                            <option value="">-- Ignore --</option>
                                            {(TARGET_FIELDS[m.target_model] || []).map(f => (
                                                <option key={f} value={f}>{f}</option>
                                            ))}
                                        </select>
                                    </div>
                                </div>
                            </GlassCard>
                        ))}

                        <button
                            onClick={handleAddMapping}
                            className="w-full py-4 border-2 border-dashed border-white/5 rounded-2xl text-[10px] uppercase font-bold text-text-muted hover:text-primary hover:border-primary/30 transition-all flex items-center justify-center gap-2 bg-white/5"
                        >
                            <Plus size={16} /> Add Custom Mapping
                        </button>
                    </div>

                    {previewData && (
                        <div className="mt-8 space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
                            <div className="flex items-center gap-2">
                                <div className="h-px flex-1 bg-white/5" />
                                <span className="text-[10px] uppercase font-bold text-text-muted tracking-widest px-3 py-1 rounded-full border border-white/5">
                                    Discovery Sample
                                </span>
                                <div className="h-px flex-1 bg-white/5" />
                            </div>
                            <div className="glass p-4 rounded-2xl overflow-hidden">
                                <pre className="text-[10px] font-mono leading-relaxed text-blue-300 overflow-x-auto max-h-[300px] custom-scrollbar">
                                    {JSON.stringify(previewData.sample_items?.[0] || previewData.raw_response, null, 2)}
                                </pre>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </Modal>
    );
};

export default IngestMappingModal;
