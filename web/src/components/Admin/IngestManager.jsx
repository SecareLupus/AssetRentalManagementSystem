import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    Plus,
    RefreshCw,
    Trash2,
    Edit,
    ExternalLink,
    CheckCircle,
    XCircle,
    Clock,
    Search,
    Database,
    Fingerprint,
    Zap,
    ChevronRight,
    Play,
    Globe,
    Activity
} from 'lucide-react';
import { GlassCard } from '../Shared';
import IngestSourceModal from './IngestSourceModal';
import IngestMappingModal from './IngestMappingModal';

const IngestManager = () => {
    const [sources, setSources] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [showSourceModal, setShowSourceModal] = useState(false);
    const [showMappingModal, setShowMappingModal] = useState(false);
    const [currentSource, setCurrentSource] = useState(null);

    useEffect(() => {
        fetchSources();
    }, []);

    const fetchSources = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/v1/admin/ingest/sources');
            setSources(res.data || []);
        } catch (err) {
            setError("Failed to load ingest sources.");
        } finally {
            setLoading(false);
        }
    };

    const handleSyncNow = async (id) => {
        try {
            await axios.post(`/v1/admin/ingest/sources/${id}/sync`);
            fetchSources();
        } catch (err) {
            alert("Failed to trigger sync.");
        }
    };

    const handleDeleteSource = async (id) => {
        if (!window.confirm("Delete this source and all its endpoints?")) return;
        try {
            await axios.delete(`/v1/admin/ingest/sources/${id}`);
            fetchSources();
        } catch (err) {
            alert("Delete failed.");
        }
    };

    const formatDate = (dateStr) => {
        if (!dateStr) return 'Never';
        try {
            return new Date(dateStr).toLocaleString();
        } catch {
            return 'Invalid Date';
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-xl font-bold flex items-center gap-2">
                        <Database className="text-primary" />
                        Universal Ingestion Engine
                    </h2>
                    <p className="text-sm text-text-muted">Harvest and map data from external REST APIs</p>
                </div>
                <button
                    onClick={() => { setCurrentSource(null); setShowSourceModal(true); }}
                    className="btn-primary flex items-center gap-2"
                >
                    <Plus size={18} />
                    New Source
                </button>
            </div>

            {loading && sources.length === 0 ? (
                <div className="flex justify-center py-12">
                    <RefreshCw className="animate-spin text-primary" size={32} />
                </div>
            ) : (
                <div className="grid gap-6">
                    {sources.map(source => (
                        <GlassCard key={source.id} className="p-5 overflow-hidden group">
                            <div className="flex justify-between items-start">
                                <div className="space-y-1">
                                    <div className="flex items-center gap-3">
                                        <h3 className="font-bold text-lg text-text group-hover:text-primary transition-colors">{source.name}</h3>
                                        <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-white/5 border border-white/10">
                                            <Globe size={10} className="text-text-muted" />
                                            <span className="text-[10px] font-bold text-text-muted uppercase tracking-wider">
                                                {source.auth_type}
                                            </span>
                                        </div>
                                        {source.is_active ?
                                            <span className="text-emerald-500 flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest bg-emerald-500/10 px-2 py-0.5 rounded-full">
                                                <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                                                Live
                                            </span> :
                                            <span className="text-text-muted text-[10px] font-bold uppercase tracking-widest bg-white/5 px-2 py-0.5 rounded-full">Paused</span>
                                        }
                                    </div>
                                    <p className="text-xs font-mono text-text-muted flex items-center gap-1">
                                        {source.base_url}
                                    </p>
                                </div>
                                <div className="flex gap-1">
                                    <button
                                        onClick={() => handleSyncNow(source.id)}
                                        className="p-2.5 hover:bg-primary/20 rounded-xl text-primary transition-all active:scale-95"
                                        title="Sync All Endpoints"
                                    >
                                        <Play size={20} fill="currentColor" />
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowMappingModal(true); }}
                                        className="p-2.5 hover:bg-white/10 rounded-xl text-text-muted hover:text-text transition-all"
                                        title="Configure Mappings"
                                    >
                                        <Fingerprint size={20} />
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowSourceModal(true); }}
                                        className="p-2.5 hover:bg-white/10 rounded-xl text-text-muted hover:text-text transition-all"
                                        title="Edit Connection"
                                    >
                                        <Edit size={20} />
                                    </button>
                                    <button
                                        onClick={() => handleDeleteSource(source.id)}
                                        className="p-2.5 hover:bg-red-500/20 rounded-xl text-red-500 transition-all opacity-40 hover:opacity-100"
                                        title="Delete Source"
                                    >
                                        <Trash2 size={20} />
                                    </button>
                                </div>
                            </div>

                            <div className="mt-6 grid grid-cols-4 gap-6 pt-5 border-t border-white/5">
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-widest text-text-muted font-bold block">Integrity</span>
                                    <div className={`text-xs font-bold flex items-center gap-1.5 ${(!source.last_error) ? 'text-emerald-400' : 'text-red-400'}`}>
                                        {(!source.last_error) ? <CheckCircle size={14} /> : <AlertTriangle size={14} />}
                                        {source.last_status || 'Checking...'}
                                    </div>
                                    {source.last_error && <div className="text-[8px] opacity-60 truncate">{source.last_error}</div>}
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-widest text-text-muted font-bold block">Endpoints</span>
                                    <div className="text-xs font-bold text-text flex items-center gap-1.5">
                                        <Activity size={14} className="text-primary" />
                                        {source.endpoints?.length || 0} Registered
                                    </div>
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-widest text-text-muted font-bold block">Last Harvest</span>
                                    <div className="text-xs font-bold text-text flex items-center gap-1.5">
                                        <Clock size={14} className="text-text-muted" />
                                        {formatDate(source.last_success_at)}
                                    </div>
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-widest text-text-muted font-bold block">Frequency</span>
                                    <div className="text-xs font-bold text-text opacity-70">
                                        Every {source.sync_interval_seconds / 60}m
                                    </div>
                                    <div className="text-[8px] text-text-muted">
                                        Next: {source.next_sync_at ? new Date(source.next_sync_at).toLocaleTimeString() : 'N/A'}
                                    </div>
                                </div>
                            </div>
                        </GlassCard>
                    ))}
                    {sources.length === 0 && !loading && (
                        <div className="text-center py-24 border-2 border-dashed border-white/5 rounded-[40px] bg-white/[0.02]">
                            <Database className="mx-auto text-text-muted/20 mb-6" size={64} />
                            <h3 className="text-xl font-bold text-text mb-2">Passive Ingestion Idle</h3>
                            <p className="text-text-muted max-w-sm mx-auto text-sm">Connect your external ecosystem to enable automated asset and item synchronization.</p>
                            <button
                                onClick={() => { setCurrentSource(null); setShowSourceModal(true); }}
                                className="mt-8 btn-primary px-8"
                            >
                                Add Discovery Source
                            </button>
                        </div>
                    )}
                </div>
            )}

            <IngestSourceModal
                isOpen={showSourceModal}
                onClose={() => setShowSourceModal(false)}
                source={currentSource}
                onSave={fetchSources}
            />

            <IngestMappingModal
                isOpen={showMappingModal}
                onClose={() => setShowMappingModal(false)}
                source={currentSource}
                onSave={fetchSources}
            />
        </div>
    );
};

const AlertTriangle = ({ size }) => <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m11.73 3 8 14a2 2 0 0 1-1.73 3H4a2 2 0 0 1-1.73-3l8-14a2 2 0 0 1 3.46 0Z" /><line x1="12" y1="9" x2="12" y2="13" /><line x1="12" y1="17" x2="12.01" y2="17" /></svg>;

export default IngestManager;
