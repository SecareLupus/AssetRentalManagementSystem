import React, { useState, useEffect } from 'react';
import { Modal, GlassCard } from '../Shared';
import axios from 'axios';
import {
    Search, Zap, CheckCircle, AlertTriangle, ChevronRight,
    Fingerprint, Plus, Trash2, Database, Target, Activity,
    Layout, Globe, Check
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
            // FIX: Backend expects a raw array of mappings, not { mappings: [...] }
            // Ensure we are sending ONLY the array.
            const payload = Array.isArray(mappings) ? mappings : [];
            await axios.post(`/v1/admin/ingest/endpoints/${selectedEndpointId}/mappings`, payload);

            if (source.endpoints) {
                const ep = source.endpoints.find(e => e.id === selectedEndpointId);
                if (ep) ep.mappings = payload;
            }
            onSave();
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
            <div className="flex-row gap-8 min-h-[600px] border-t border-white/5 pt-6">
                {/* Left Sidebar: Endpoints */}
                <div className="w-60 pr-4 border-r border-white/5 space-y-3 shrink-0">
                    <div className="px-1 mb-4">
                        <label className="text-[10px] uppercase font-black text-primary tracking-[0.2em] mb-1 block">API Endpoints</label>
                        <div className="text-[9px] text-text-muted font-medium">Select an endpoint to map models</div>
                    </div>
                    {source.endpoints.map(ep => (
                        <button
                            key={ep.id}
                            onClick={() => { setSelectedEndpointId(ep.id); loadEndpointMappings(ep.id); }}
                            className={`w-full text-left p-3 rounded-2xl transition-all border group relative overflow-hidden ${selectedEndpointId === ep.id
                                ? 'bg-primary/10 border-primary/40 ring-1 ring-primary/40'
                                : 'bg-white/5 border-transparent hover:border-white/10'
                                }`}
                        >
                            {selectedEndpointId === ep.id && (
                                <div className="absolute left-0 top-0 bottom-0 w-1 bg-primary" />
                            )}
                            <div className="flex items-center gap-3 mb-2">
                                <span className={`text-[9px] font-black px-2 py-0.5 rounded-md ${ep.method === 'GET' ? 'bg-blue-500/20 text-blue-400' : 'bg-green-500/20 text-green-400'
                                    }`}>
                                    {ep.method}
                                </span>
                                <span className={`text-[10px] font-bold truncate uppercase tracking-tight ${selectedEndpointId === ep.id ? 'text-primary' : 'text-white'}`}>
                                    {ep.path}
                                </span>
                            </div>
                            <div className="flex items-center justify-between">
                                <div className={`text-[9px] flex items-center gap-1.5 font-bold uppercase tracking-wider ${selectedEndpointId === ep.id ? 'text-primary/60' : 'text-text-muted opacity-60'}`}>
                                    <Activity size={10} /> {ep.mappings?.length || 0} mapped
                                </div>
                                {ep.mappings?.length > 0 && <CheckCircle size={10} className="text-primary" />}
                            </div>
                        </button>
                    ))}
                </div>

                {/* Main: Mapping UI */}
                <div className="flex-1 space-y-6 max-h-[70vh] overflow-y-auto pr-3 custom-scrollbar">
                    <div className="flex items-center justify-between px-1">
                        <div>
                            <h4 className="font-bold flex items-center gap-2 text-primary">
                                <Database size={16} /> Data Field Mapping
                            </h4>
                            <p className="text-[10px] text-text-muted font-medium">Transform external JSON attributes to internal fields</p>
                        </div>
                        <button
                            onClick={handleDiscovery}
                            disabled={loading || !selectedEndpointId}
                            className="btn-secondary flex items-center gap-2 py-2 px-4 h-auto text-[10px] border-yellow-500/30 text-yellow-500 hover:bg-yellow-500/5 transition-all active:scale-95 disabled:opacity-50"
                        >
                            <Zap size={14} className={loading ? 'animate-spin' : ''} />
                            Run Discovery
                        </button>
                    </div>

                    {error && (
                        <div className="p-4 bg-red-500/10 border border-red-500/20 rounded-2xl text-red-500 text-xs flex items-center gap-3 animate-in fade-in zoom-in-95 duration-300">
                            <AlertTriangle size={18} />
                            <span className="font-medium">{error}</span>
                        </div>
                    )}

                    <div className="space-y-4">
                        {mappings.map((m, idx) => (
                            <GlassCard key={idx} className="p-4 relative group hover:ring-2 hover:ring-primary/20 transition-all border-white/5 bg-white/[0.03]">
                                <div className="grid grid-cols-12 gap-6 items-end">
                                    {/* Attribute Path & ID Key (4 columns) */}
                                    <div className="col-span-4 space-y-2">
                                        <div className="flex flex-row items-center justify-between px-1 h-5">
                                            <label className="text-[10px] uppercase font-black text-text-muted flex flex-row items-center gap-1.5 tracking-[0.15em]">
                                                <Layout size={12} className="text-primary" /> Attribute Path
                                            </label>
                                            <label className="flex flex-row items-center gap-2 cursor-pointer group/id relative hover:opacity-80 transition-opacity">
                                                <input
                                                    type="checkbox"
                                                    className="w-4 h-4 accent-primary opacity-0 absolute cursor-pointer z-10"
                                                    checked={!!m.is_identity}
                                                    onChange={e => handleUpdateMapping(idx, 'is_identity', e.target.checked)}
                                                />
                                                <div className={`w-3.5 h-3.5 rounded-md border-2 flex items-center justify-center transition-all ${m.is_identity ? 'bg-primary border-primary scale-110 shadow-lg shadow-primary/20' : 'border-white/20'}`}>
                                                    {m.is_identity && <Check size={8} className="text-white" />}
                                                </div>
                                                <span className={`text-[8px] font-black uppercase tracking-[0.2em] ${m.is_identity ? 'text-primary' : 'text-text-muted'}`}>ID Key</span>
                                            </label>
                                        </div>
                                        <input
                                            className="glass w-full h-form-input px-2.5 text-sm font-mono text-blue-300 border-white/5 focus:border-primary/40 focus:bg-primary/5 transition-all outline-none"
                                            value={m.json_path}
                                            placeholder="$.data.id"
                                            onChange={e => handleUpdateMapping(idx, 'json_path', e.target.value)}
                                        />
                                    </div>

                                    {/* Model Selection (3 columns) */}
                                    <div className="col-span-3 space-y-2">
                                        <div className="h-5 flex items-center px-1">
                                            <span className="text-[10px] uppercase font-black text-text-muted flex items-center gap-1.5 tracking-[0.15em]">
                                                <Target size={12} className="text-primary" /> Model
                                            </span>
                                        </div>
                                        <select
                                            className="glass w-full h-form-input px-2.5 text-[10px] font-bold uppercase tracking-wider border-white/5"
                                            value={m.target_model}
                                            onChange={e => handleUpdateMapping(idx, 'target_model', e.target.value)}
                                        >
                                            {TARGET_MODELS.map(tm => (
                                                <option key={tm.value} value={tm.value}>{tm.label}</option>
                                            ))}
                                        </select>
                                    </div>

                                    {/* Field Selection (4 columns) */}
                                    <div className="col-span-4 space-y-2">
                                        <div className="h-5 flex items-center px-1">
                                            <span className="text-[10px] uppercase font-black text-text-muted tracking-[0.15em] block">Field</span>
                                        </div>
                                        <select
                                            className="glass w-full h-form-input px-2.5 text-[10px] font-bold uppercase tracking-wider text-primary border-white/5"
                                            value={m.target_field}
                                            onChange={e => handleUpdateMapping(idx, 'target_field', e.target.value)}
                                        >
                                            <option value="">-- Ignore --</option>
                                            {(TARGET_FIELDS[m.target_model] || []).map(f => (
                                                <option key={f} value={f}>{f}</option>
                                            ))}
                                        </select>
                                    </div>

                                    {/* Delete Action (1 column) */}
                                    <div className="col-span-1 space-y-2">
                                        <div className="h-5" />
                                        <button
                                            onClick={() => handleRemoveMapping(idx)}
                                            className="w-full h-form-input rounded-xl bg-red-500/5 border border-red-500/10 text-red-500/40 hover:text-red-500 hover:bg-red-500/10 hover:border-red-500/40 transition-all active:scale-90 flex items-center justify-center"
                                            title="Remove Mapping"
                                        >
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </div>
                            </GlassCard>
                        ))}

                        <button
                            onClick={handleAddMapping}
                            className="w-full py-4 border-2 border-dashed border-white/5 rounded-2xl text-[10px] uppercase font-black tracking-[0.2em] text-text-muted hover:text-primary hover:border-primary/40 hover:bg-primary/5 transition-all flex items-center justify-center gap-3 bg-white/[0.01] active:scale-[0.99] group shadow-inner"
                        >
                            <div className="w-8 h-8 rounded-full bg-white/5 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                                <Plus size={16} className="group-hover:rotate-180 transition-transform duration-500" />
                            </div>
                            Add Custom Property Mapping
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
        </Modal >
    );
};

export default IngestMappingModal;
