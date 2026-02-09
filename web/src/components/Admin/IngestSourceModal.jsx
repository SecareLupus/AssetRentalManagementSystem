import React, { useState, useEffect } from 'react';
import { Modal, GlassCard } from '../Shared';
import axios from 'axios';
import {
    CheckCircle2, AlertCircle, Loader2, Plus, Trash2,
    Settings, Lock, Globe, Activity, ArrowRight, ArrowLeft, Play
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
                { path: '', method: 'GET', resp_strategy: 'auto', is_active: true }
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
        <div className="flex items-center justify-between mb-8 px-4">
            {[1, 2, 3].map(i => (
                <React.Fragment key={i}>
                    <div className="flex flex-col items-center gap-2">
                        <div className={`w-10 h-10 rounded-full flex items-center justify-center border-2 transition-all ${step === i ? 'border-primary bg-primary/20 text-text' :
                            step > i ? 'border-green-500 bg-green-500/20 text-green-500' :
                                'border-border text-text-muted'
                            }`}>
                            {step > i ? <CheckCircle2 size={20} /> : i}
                        </div>
                        <span className={`text-[10px] uppercase tracking-wider font-bold ${step === i ? 'text-primary' : 'text-text-muted'}`}>
                            {i === 1 ? 'Connection' : i === 2 ? 'Security' : 'Endpoints'}
                        </span>
                    </div>
                    {i < 3 && <div className={`flex-1 h-0.5 mx-4 ${step > i ? 'bg-green-500' : 'bg-border'}`} />}
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

            <div className="min-h-[300px]">
                {step === 1 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="grid grid-cols-2 gap-6">
                            <div className="space-y-4">
                                <div>
                                    <label className="form-label text-primary flex items-center gap-2">
                                        <Settings size={14} /> Source Name
                                    </label>
                                    <input
                                        className="glass w-full p-3 text-sm"
                                        placeholder="e.g. Siemens ERP Extension"
                                        value={formData.name}
                                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                                    />
                                </div>
                                <div>
                                    <label className="form-label text-primary flex items-center gap-2">
                                        <Globe size={14} /> Base API URL
                                    </label>
                                    <input
                                        className="glass w-full p-3 text-sm"
                                        placeholder="https://api.factory.internal/v2"
                                        value={formData.base_url}
                                        onChange={e => setFormData({ ...formData, base_url: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="space-y-4">
                                <div>
                                    <label className="form-label text-primary flex items-center gap-2">
                                        <Activity size={14} /> Sync Interval
                                    </label>
                                    <select
                                        className="glass w-full p-3 text-sm"
                                        value={formData.sync_interval_seconds}
                                        onChange={e => setFormData({ ...formData, sync_interval_seconds: parseInt(e.target.value) })}
                                    >
                                        <option value={60}>Every Minute</option>
                                        <option value={3600}>Every Hour</option>
                                        <option value={86400}>Once Daily</option>
                                        <option value={604800}>Weekly</option>
                                    </select>
                                </div>
                                <div className="pt-8">
                                    <label className="flex items-center gap-3 cursor-pointer p-4 rounded-xl border border-white/5 bg-white/5 hover:bg-white/10 transition-colors">
                                        <input
                                            type="checkbox"
                                            className="w-5 h-5 accent-primary"
                                            checked={formData.is_active}
                                            onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                                        />
                                        <div>
                                            <div className="font-bold text-sm">Automated Collection</div>
                                            <div className="text-[10px] text-text-muted">Enable background sync worker for this source</div>
                                        </div>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {step === 2 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="grid grid-cols-2 gap-6">
                            <div className="space-y-4">
                                <div>
                                    <label className="form-label text-primary flex items-center gap-2">
                                        <Lock size={14} /> Auth Architecture
                                    </label>
                                    <div className="grid grid-cols-2 gap-2">
                                        {['none', 'bearer'].map(type => (
                                            <button
                                                key={type}
                                                onClick={() => setFormData({ ...formData, auth_type: type })}
                                                className={`p-3 rounded-xl border text-xs font-bold transition-all ${formData.auth_type === type
                                                    ? 'border-primary bg-primary/20 text-primary'
                                                    : 'border-white/10 bg-white/5 text-text-muted'
                                                    }`}
                                            >
                                                {type.toUpperCase()}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                                {formData.auth_type === 'bearer' && (
                                    <div className="space-y-4">
                                        <div>
                                            <label className="form-label text-primary">Initial Auth Endpoint</label>
                                            <input
                                                className="glass w-full p-3 text-sm"
                                                placeholder="/auth/login"
                                                value={formData.auth_endpoint}
                                                onChange={e => setFormData({ ...formData, auth_endpoint: e.target.value })}
                                            />
                                        </div>
                                        <div>
                                            <label className="form-label text-primary flex items-center gap-2">
                                                Verification Endpoint
                                                <span className="text-[10px] text-text-muted font-normal">(Optional)</span>
                                            </label>
                                            <input
                                                className="glass w-full p-3 text-sm"
                                                placeholder="/auth/verify"
                                                value={formData.verify_endpoint}
                                                onChange={e => setFormData({ ...formData, verify_endpoint: e.target.value })}
                                            />
                                        </div>
                                        <div>
                                            <label className="form-label text-primary flex items-center gap-2">
                                                Refresh Token Endpoint
                                                <span className="text-[10px] text-text-muted font-normal">(Optional)</span>
                                            </label>
                                            <input
                                                className="glass w-full p-3 text-sm"
                                                placeholder="/auth/refresh"
                                                value={formData.refresh_endpoint}
                                                onChange={e => setFormData({ ...formData, refresh_endpoint: e.target.value })}
                                            />
                                        </div>
                                    </div>
                                )}
                            </div>

                            {formData.auth_type !== 'none' && (
                                <div className="space-y-4">
                                    <div>
                                        <label className="form-label text-primary">Credentials (JSON)</label>
                                        <textarea
                                            className="glass w-full p-3 font-mono text-xs"
                                            rows={5}
                                            placeholder='{\n  "client_id": "...",\n  "client_secret": "..."\n}'
                                            value={formData.auth_credentials}
                                            onChange={e => setFormData({ ...formData, auth_credentials: e.target.value })}
                                        />
                                    </div>
                                    <button
                                        onClick={handleTestAuth}
                                        disabled={testStatus === 'testing'}
                                        className="w-full flex items-center justify-center gap-2 p-3 rounded-xl bg-primary/10 border border-primary/20 text-primary font-bold hover:bg-primary/20 transition-all disabled:opacity-50"
                                    >
                                        {testStatus === 'testing' ? <Loader2 className="animate-spin" size={16} /> : <Play size={16} />}
                                        Test Authentication Lifecycle
                                    </button>

                                    {testStatus === 'success' && (
                                        <div className="flex items-center gap-2 text-green-500 text-[10px] font-bold">
                                            <CheckCircle2 size={12} /> Authentication Successful & Token Cached
                                        </div>
                                    )}
                                    {testStatus === 'error' && (
                                        <div className="flex items-center gap-2 text-red-500 text-[10px] font-bold">
                                            <AlertCircle size={12} /> {testError}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {step === 3 && (
                    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="flex items-center justify-between">
                            <h4 className="text-sm font-bold uppercase tracking-widest text-text-muted">Service Endpoints</h4>
                            <button onClick={handleAddEndpoint} className="btn-secondary flex items-center gap-2 py-1 h-auto text-[10px]">
                                <Plus size={14} /> Add Endpoint
                            </button>
                        </div>

                        <div className="space-y-3 max-h-[400px] overflow-y-auto pr-2">
                            {formData.endpoints.map((ep, idx) => (
                                <GlassCard key={idx} className="p-4 relative">
                                    <button
                                        onClick={() => handleRemoveEndpoint(idx)}
                                        className="absolute top-4 right-4 text-text-muted hover:text-red-500 transition-colors"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                    <div className="grid grid-cols-12 gap-4">
                                        <div className="col-span-3">
                                            <label className="text-[10px] uppercase font-bold text-text-muted mb-1 block">Method</label>
                                            <select
                                                className="glass w-full p-2 text-xs"
                                                value={ep.method}
                                                onChange={e => handleEndpointChange(idx, 'method', e.target.value)}
                                            >
                                                <option value="GET">GET</option>
                                                <option value="POST">POST</option>
                                            </select>
                                        </div>
                                        <div className="col-span-6">
                                            <label className="text-[10px] uppercase font-bold text-text-muted mb-1 block">Resource Path</label>
                                            <input
                                                className="glass w-full p-2 text-xs"
                                                placeholder="/products/inventory"
                                                value={ep.path}
                                                onChange={e => handleEndpointChange(idx, 'path', e.target.value)}
                                            />
                                        </div>
                                        <div className="col-span-3">
                                            <label className="text-[10px] uppercase font-bold text-text-muted mb-1 block">Strategy</label>
                                            <select
                                                className="glass w-full p-2 text-xs"
                                                value={ep.resp_strategy}
                                                onChange={e => handleEndpointChange(idx, 'resp_strategy', e.target.value)}
                                            >
                                                <option value="auto">Auto</option>
                                                <option value="list">List</option>
                                                <option value="single">Single</option>
                                            </select>
                                        </div>
                                        {ep.method === 'POST' && (
                                            <div className="col-span-12 mt-2">
                                                <label className="text-[10px] uppercase font-bold text-text-muted mb-1 block text-primary">Request Body (JSON)</label>
                                                <textarea
                                                    className="glass w-full p-2 text-xs font-mono"
                                                    rows={3}
                                                    placeholder='{ "id": "123" }'
                                                    value={typeof ep.request_body === 'string' ? ep.request_body : JSON.stringify(ep.request_body, null, 2)}
                                                    onChange={e => handleEndpointChange(idx, 'request_body', e.target.value)}
                                                />
                                            </div>
                                        )}
                                    </div>
                                </GlassCard>
                            ))}
                            {formData.endpoints.length === 0 && (
                                <div className="text-center py-12 border-2 border-dashed border-white/5 rounded-3xl">
                                    <div className="text-text-muted text-sm italic mb-2">No endpoints defined yet</div>
                                    <button onClick={handleAddEndpoint} className="text-primary text-[10px] font-bold uppercase hover:underline">
                                        Click to add your first API call
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
