import React, { useState, useEffect } from 'react';
import { Modal, GlassCard } from '../Shared';
import axios from 'axios';
import {
    CheckCircle2, AlertCircle, Loader2, Plus, Trash2,
    Settings, Lock, Globe, Activity, ArrowRight, ArrowLeft, Play, Check
} from 'lucide-react';

const IngestSourceModal = ({ isOpen, onClose, source, onSave }) => {
    const [step, setStep] = useState(1);
    const [formData, setFormData] = useState({
        name: '',
        base_url: '',
        auth_type: 'none',
        auth_endpoint: '',
        verify_endpoint: '',
        refresh_endpoint: '',
        auth_credentials: '',
        sync_interval_seconds: 3600,
        is_active: true,
        endpoints: []
    });

    const [testStatus, setTestStatus] = useState(null); // 'testing', 'success', 'error'
    const [testError, setTestError] = useState('');

    useEffect(() => {
        if (source) {
            setFormData({
                ...source,
                auth_credentials: source.auth_credentials ? (typeof source.auth_credentials === 'string' ? JSON.stringify(JSON.parse(source.auth_credentials), null, 2) : JSON.stringify(source.auth_credentials, null, 2)) : '',
                endpoints: source.endpoints || []
            });
            setStep(1);
        } else {
            setFormData({
                name: '',
                base_url: '',
                auth_type: 'none',
                auth_endpoint: '',
                verify_endpoint: '',
                refresh_endpoint: '',
                auth_credentials: '',
                sync_interval_seconds: 3600,
                is_active: true,
                endpoints: []
            });
            setStep(1);
        }
    }, [source, isOpen]);

    const handleTestAuth = async () => {
        setTestStatus('testing');
        setTestError('');
        try {
            let activeId = source?.id || formData.id;

            // If it's a new source, we must save it first to have an ID in the DB
            if (!activeId) {
                const data = { ...formData, auth_credentials: formatCreds(formData.auth_credentials) };
                const res = await axios.post('/v1/admin/ingest/sources', data);
                activeId = res.data.id;
                setFormData(prev => ({ ...prev, id: activeId }));
            }

            const res = await axios.post('/v1/admin/ingest/test-auth', { source_id: activeId });
            setTestStatus('success');
        } catch (err) {
            setTestStatus('error');
            setTestError(err.response?.data || err.message);
        }
    };

    const formatCreds = (creds) => {
        if (!creds) return null;
        try {
            return JSON.parse(creds);
        } catch {
            return creds;
        }
    };

    const handleAddEndpoint = () => {
        setFormData(prev => ({
            ...prev,
            endpoints: [
                ...prev.endpoints,
                { path: '', method: 'GET', resp_strategy: 'auto', items_path: '$', is_active: true }
            ]
        }));
    };

    const handleRemoveEndpoint = async (index) => {
        const ep = formData.endpoints[index];
        if (ep.id) {
            if (!window.confirm("Are you sure you want to delete this endpoint? This cannot be undone.")) return;
            try {
                await axios.delete(`/v1/admin/ingest/endpoints/${ep.id}`);
            } catch (err) {
                alert("Failed to delete endpoint: " + (err.response?.data || err.message));
                return;
            }
        }
        setFormData(prev => ({
            ...prev,
            endpoints: prev.endpoints.filter((_, i) => i !== index)
        }));
    };

    const handleEndpointChange = (index, field, value) => {
        const newEndpoints = [...formData.endpoints];
        newEndpoints[index] = { ...newEndpoints[index], [field]: value };
        setFormData(prev => ({ ...prev, endpoints: newEndpoints }));
    };

    const handleSubmit = async () => {
        try {
            const data = {
                ...formData,
                auth_credentials: formatCreds(formData.auth_credentials)
            };

            let sourceId = source?.id || formData.id;
            if (sourceId) {
                await axios.put(`/v1/admin/ingest/sources/${sourceId}`, data);
            } else {
                const res = await axios.post('/v1/admin/ingest/sources', data);
                sourceId = res.data.id;
            }

            // Endpoints are handled separately or by updated repository
            // In my SqlRepository, GetIngestSource includes endpoints, 
            // but UpdateSource doesn't automatically upsert them.
            // I should handle endpoints CRUD here if needed, 
            // or assume backend handles it if I send them in the source object.
            // My backend currently doesn't auto-upsert endpoints from UpdateSource.

            // For each endpoint, create or update
            for (const ep of formData.endpoints) {
                const epData = {
                    ...ep,
                    source_id: sourceId,
                    request_body: formatCreds(ep.request_body) // Reuse same logic for body
                };
                if (ep.id) {
                    await axios.put(`/v1/admin/ingest/endpoints/${ep.id}`, epData);
                } else {
                    await axios.post('/v1/admin/ingest/endpoints', epData);
                }
            }

            onSave();
            onClose();
        } catch (err) {
            alert("Failed to save: " + (err.response?.data || err.message));
        }
    };

    const StepIndicator = () => (
        <div className="flex flex-row items-center justify-center gap-4 mb-12 relative px-4 w-full max-w-2xl mx-auto">
            {[1, 2, 3].map(i => (
                <React.Fragment key={i}>
                    <div className="flex flex-col items-center gap-3 relative z-10 min-w-[100px] shrink-0">
                        <div className={`w-12 h-12 rounded-full flex items-center justify-center font-bold transition-all duration-300 ${step === i ? 'bg-primary text-white scale-110 ring-4 ring-primary/20 shadow-lg shadow-primary/20' :
                            step > i ? 'bg-green-500 text-white' :
                                'bg-surface border-2 border-white/10 text-text-muted'
                            }`}>
                            {step > i ? <CheckCircle2 size={24} /> : i}
                        </div>
                        <span className={`text-[10px] uppercase tracking-widest font-black transition-colors ${step === i ? 'text-primary' : 'text-text-muted'}`}>
                            {i === 1 ? 'Source' : i === 2 ? 'Security' : 'Endpoints'}
                        </span>
                    </div>
                    {i < 3 && (
                        <div className="flex-1 h-[2px] bg-white/5 relative min-w-[30px] -mt-6">
                            <div className={`absolute inset-0 bg-primary transition-all duration-500 ${step > i ? 'w-full' : 'w-0'}`} />
                        </div>
                    )}
                </React.Fragment>
            ))}
        </div>
    );

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="lg"
            title={source ? `Edit ${source.name}` : "Connect Data Source"}
            actions={(
                <div className="flex gap-3">
                    {step > 1 && (
                        <button onClick={() => setStep(step - 1)} className="btn-secondary flex items-center gap-2">
                            <ArrowLeft size={16} /> Back
                        </button>
                    )}
                    {step < 3 ? (
                        <button onClick={() => setStep(step + 1)} className="btn-primary flex items-center gap-2">
                            Next <ArrowRight size={16} />
                        </button>
                    ) : (
                        <button onClick={handleSubmit} className="btn-primary flex items-center gap-2">
                            Finish Setup <CheckCircle2 size={16} />
                        </button>
                    )}
                </div>
            )}
        >
            <StepIndicator />

            <div className="min-h-[350px]">
                {step === 1 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="grid grid-cols-1 gap-6">
                            <div className="space-y-5">
                                <div className="space-y-2">
                                    <label className="form-label text-primary flex items-center gap-2 px-1">
                                        <Settings size={14} /> Source Name
                                    </label>
                                    <input
                                        className="glass w-full p-3.5 text-sm"
                                        placeholder="e.g. Siemens ERP Extension"
                                        value={formData.name}
                                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                                    />
                                </div>
                                <div className="space-y-2">
                                    <label className="form-label text-primary flex items-center gap-2 px-1">
                                        <Globe size={14} /> Base API URL
                                    </label>
                                    <input
                                        className="glass w-full p-3.5 text-sm"
                                        placeholder="https://api.factory.internal/v2"
                                        value={formData.base_url}
                                        onChange={e => setFormData({ ...formData, base_url: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="space-y-5">
                                <div className="space-y-2">
                                    <label className="form-label text-primary flex items-center gap-2 px-1">
                                        <Activity size={14} /> Sync Interval
                                    </label>
                                    <select
                                        className="glass w-full p-3.5 text-sm"
                                        value={formData.sync_interval_seconds}
                                        onChange={e => setFormData({ ...formData, sync_interval_seconds: parseInt(e.target.value) })}
                                    >
                                        <option value={60}>Every Minute</option>
                                        <option value={3600}>Every Hour</option>
                                        <option value={86400}>Once Daily</option>
                                        <option value={604800}>Weekly</option>
                                    </select>
                                </div>
                                <div className="pt-6">
                                    <label className="flex flex-row items-center gap-6 cursor-pointer p-6 rounded-[24px] border border-white/5 bg-white/5 hover:bg-white/10 hover:border-primary/30 transition-all group relative w-full">
                                        <div className="relative flex flex-row items-center justify-center shrink-0">
                                            <input
                                                type="checkbox"
                                                className="w-7 h-7 accent-primary opacity-0 absolute cursor-pointer z-10"
                                                checked={formData.is_active}
                                                onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                                            />
                                            <div className={`w-7 h-7 rounded-lg border-2 flex items-center justify-center transition-all ${formData.is_active ? 'bg-primary border-primary scale-110 shadow-lg shadow-primary/20' : 'border-white/20'}`}>
                                                {formData.is_active && <Check size={18} className="text-white" />}
                                            </div>
                                        </div>
                                        <div className="flex flex-col">
                                            <div className="flex flex-row items-center gap-2">
                                                <div className={`font-black text-sm uppercase tracking-tight transition-colors ${formData.is_active ? 'text-primary' : 'text-text'}`}>
                                                    Enable Automatic Background Sync
                                                </div>
                                                {formData.is_active && <Activity size={14} className="text-primary animate-pulse" />}
                                            </div>
                                            <div className="text-[11px] text-text-muted font-medium mt-1 leading-relaxed">
                                                Active background worker will poll this source at the specified interval
                                            </div>
                                        </div>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {step === 2 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="grid grid-cols-1 gap-8">
                            <div className="space-y-5">
                                <div className="space-y-2">
                                    <label className="form-label text-primary flex items-center gap-2 px-1">
                                        <Lock size={14} /> Auth Architecture
                                    </label>
                                    <div className="grid grid-cols-2 gap-3">
                                        {['none', 'bearer'].map(type => (
                                            <button
                                                key={type}
                                                type="button"
                                                onClick={() => setFormData({ ...formData, auth_type: type })}
                                                className={`p-4 rounded-2xl border text-xs font-black uppercase tracking-widest transition-all ${formData.auth_type === type
                                                    ? 'border-primary bg-primary/20 text-text ring-2 ring-primary/40 shadow-lg shadow-primary/10'
                                                    : 'border-white/5 bg-white/5 text-text-muted hover:border-white/10 hover:text-text'
                                                    }`}
                                            >
                                                {type === 'bearer' ? 'Bearer Token' : 'No Auth'}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                                {formData.auth_type === 'bearer' && (
                                    <div className="space-y-4 animate-in fade-in zoom-in-95 duration-300">
                                        <div className="space-y-2">
                                            <label className="form-label text-primary px-1">Initial Auth Endpoint</label>
                                            <input
                                                className="glass w-full p-3.5 text-sm"
                                                placeholder="/auth/login"
                                                value={formData.auth_endpoint}
                                                onChange={e => setFormData({ ...formData, auth_endpoint: e.target.value })}
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="form-label text-primary flex items-center gap-2 px-1">
                                                Verification Endpoint
                                                <span className="text-[10px] text-text-muted font-normal lowercase tracking-normal">(optional)</span>
                                            </label>
                                            <input
                                                className="glass w-full p-3.5 text-sm"
                                                placeholder="/auth/verify"
                                                value={formData.verify_endpoint}
                                                onChange={e => setFormData({ ...formData, verify_endpoint: e.target.value })}
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="form-label text-primary flex items-center gap-2 px-1">
                                                Refresh Token Endpoint
                                                <span className="text-[10px] text-text-muted font-normal lowercase tracking-normal">(optional)</span>
                                            </label>
                                            <input
                                                className="glass w-full p-3.5 text-sm"
                                                placeholder="/auth/refresh"
                                                value={formData.refresh_endpoint}
                                                onChange={e => setFormData({ ...formData, refresh_endpoint: e.target.value })}
                                            />
                                        </div>
                                    </div>
                                )}
                            </div>

                            {formData.auth_type !== 'none' && (
                                <div className="space-y-5">
                                    <div className="space-y-2">
                                        <label className="form-label text-primary px-1 text-xs">Credentials (JSON Secret)</label>
                                        <textarea
                                            className="glass w-full p-4 font-mono text-xs leading-relaxed"
                                            rows={6}
                                            placeholder='{\n  "client_id": "...",\n  "client_secret": "..."\n}'
                                            value={formData.auth_credentials}
                                            onChange={e => setFormData({ ...formData, auth_credentials: e.target.value })}
                                        />
                                    </div>
                                    <button
                                        onClick={handleTestAuth}
                                        disabled={testStatus === 'testing'}
                                        className="w-full flex items-center justify-center gap-3 p-4 rounded-2xl bg-primary border-2 border-primary/20 text-white font-black uppercase tracking-widest text-[11px] hover:bg-primary/90 hover:scale-[1.02] shadow-xl shadow-primary/20 transition-all disabled:opacity-50 disabled:scale-100 active:scale-95 group"
                                    >
                                        {testStatus === 'testing' ? <Loader2 className="animate-spin" size={18} /> :
                                            <Play size={18} fill="currentColor" className="group-hover:translate-x-1 transition-transform" />
                                        }
                                        Run Connection Test
                                    </button>

                                    {testStatus === 'success' && (
                                        <div className="flex items-center gap-3 px-4 py-3 rounded-xl bg-green-500/10 border border-green-500/20 text-green-500 text-[11px] font-bold animate-in bounce-in duration-300">
                                            <CheckCircle2 size={16} /> Authentication Successful & Token Cached
                                        </div>
                                    )}
                                    {testStatus === 'error' && (
                                        <div className="flex items-center gap-3 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-500 text-[11px] font-bold animate-in bounce-in duration-300">
                                            <AlertCircle size={16} /> {testError}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {step === 3 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="flex items-center justify-between px-1">
                            <div>
                                <h4 className="text-sm font-bold uppercase tracking-widest text-text">Service Endpoints</h4>
                                <p className="text-[10px] text-text-muted">Define the REST paths to be polled</p>
                            </div>
                            <button onClick={handleAddEndpoint} className="btn-secondary flex items-center gap-2 py-2 px-4 h-auto text-[10px] border-primary/30 text-primary hover:bg-primary/5">
                                <Plus size={14} /> Add New API Call
                            </button>
                        </div>

                        <div className="space-y-4 max-h-[400px] overflow-y-auto pr-3 custom-scrollbar">
                            {formData.endpoints.map((ep, idx) => (
                                <GlassCard key={idx} className="p-5 relative border-white/5 hover:border-white/10 transition-colors group/ep">
                                    <button
                                        onClick={() => handleRemoveEndpoint(idx)}
                                        className="absolute top-4 right-4 text-text-muted hover:text-red-500 transition-all opacity-0 group-hover/ep:opacity-100"
                                    >
                                        <Trash2 size={16} />
                                    </button>
                                    <div className="grid grid-cols-12 gap-5">
                                        <div className="col-span-3 space-y-2">
                                            <label className="text-[10px] uppercase font-bold text-text-muted px-1">Method</label>
                                            <select
                                                className="glass w-full p-2.5 text-xs font-bold min-w-[100px]"
                                                value={ep.method}
                                                onChange={e => handleEndpointChange(idx, 'method', e.target.value)}
                                            >
                                                <option value="GET">GET</option>
                                                <option value="POST">POST</option>
                                            </select>
                                        </div>
                                        <div className="col-span-6 space-y-2">
                                            <label className="text-[10px] uppercase font-bold text-text-muted px-1">Resource Path</label>
                                            <input
                                                className="glass w-full p-2.5 text-xs font-medium"
                                                placeholder="/products/inventory"
                                                value={ep.path}
                                                onChange={e => handleEndpointChange(idx, 'path', e.target.value)}
                                            />
                                        </div>
                                        <div className="col-span-3 space-y-2">
                                            <label className="text-[10px] uppercase font-bold text-text-muted px-1">Processing</label>
                                            <select
                                                className="glass w-full p-2.5 text-xs min-w-[120px]"
                                                value={ep.resp_strategy}
                                                onChange={e => handleEndpointChange(idx, 'resp_strategy', e.target.value)}
                                            >
                                                <option value="auto">Auto-detect</option>
                                                <option value="list">List Payload</option>
                                                <option value="single">Single Object</option>
                                            </select>
                                        </div>
                                        <div className="col-span-12 space-y-2">
                                            <label className="text-[10px] uppercase font-bold text-text-muted px-1">Items List Path (JSONPath Expression)</label>
                                            <input
                                                className="glass w-full p-2.5 text-xs font-mono text-primary"
                                                placeholder="$.data or $"
                                                value={ep.items_path || '$'}
                                                onChange={e => handleEndpointChange(idx, 'items_path', e.target.value)}
                                            />
                                        </div>
                                        {ep.method === 'POST' && (
                                            <div className="col-span-12 space-y-2 animate-in fade-in slide-in-from-top-2 duration-300">
                                                <label className="text-[10px] uppercase font-bold text-primary px-1">Payload Template (JSON Body)</label>
                                                <textarea
                                                    className="glass w-full p-3 text-xs font-mono leading-relaxed"
                                                    rows={4}
                                                    placeholder='{ "filter": { "category": "assets" } }'
                                                    value={typeof ep.request_body === 'string' ? ep.request_body : JSON.stringify(ep.request_body, null, 2)}
                                                    onChange={e => handleEndpointChange(idx, 'request_body', e.target.value)}
                                                />
                                            </div>
                                        )}
                                    </div>
                                </GlassCard>
                            ))}
                            {formData.endpoints.length === 0 && (
                                <div className="text-center py-16 border-2 border-dashed border-white/5 rounded-[32px] bg-white/[0.02] animate-in fade-in duration-500">
                                    <Activity className="mx-auto text-text-muted/20 mb-4" size={48} />
                                    <div className="text-text-muted text-sm font-medium mb-1">No API endpoints registered</div>
                                    <div className="text-[10px] text-text-muted/60 mb-6">You need at least one endpoint to fetch data</div>
                                    <button onClick={handleAddEndpoint} className="text-primary text-[10px] font-bold uppercase tracking-widest px-6 py-2 rounded-full border border-primary/20 hover:bg-primary/5 transition-all">
                                        Setup Discovery Endpoint
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </Modal>
    );
};

export default IngestSourceModal;
