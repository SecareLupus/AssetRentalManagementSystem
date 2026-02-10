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
    Database,
    Fingerprint,
    Zap,
    ChevronRight,
    Play,
    Globe,
    Activity,
    AlertTriangle
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
        <div className="space-y-6 uie-manager-page">
            <div className="flex flex-row justify-between items-center mb-8 px-2 mt-4 w-full">
                <div className="flex flex-col">
                    <h2 className="text-2xl font-black text-text flex items-center gap-3">
                        <Activity className="text-primary" size={28} /> Universal Ingestion Engine
                    </h2>
                    <p className="text-xs text-text-muted font-medium mt-1">Manage external data sources and automated polling intervals</p>
                </div>
                <button
                    onClick={() => { setCurrentSource(null); setShowSourceModal(true); }}
                    className="btn-primary flex items-center gap-3 py-3 px-6 shadow-xl shadow-primary/20 hover:scale-[1.02] active:scale-95 transition-all font-bold uppercase tracking-widest text-[11px] ml-auto"
                >
                    <Plus size={18} strokeWidth={3} /> Register New Data Source
                </button>
            </div>

            {loading && sources.length === 0 ? (
                <div className="flex justify-center py-12">
                    <RefreshCw className="animate-spin text-primary" size={32} />
                </div>
            ) : (
                <div className="grid grid-cols-2 gap-8">
                    {sources.map(source => (
                        <GlassCard key={source.id} className="p-5 overflow-hidden group">
                            <div className="flex flex-col space-y-4">
                                <div className="space-y-4 min-w-0">
                                    <div className="flex items-center gap-3 min-w-0">
                                        <h3 className="font-bold text-lg text-text group-hover:text-primary transition-colors truncate min-w-0">{source.name}</h3>
                                        <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-white/5 border border-white/10 flex-shrink-0">
                                            <Globe size={10} className="text-text-muted" />
                                            <span className="text-[10px] font-bold text-text-muted uppercase tracking-wider">
                                                {source.auth_type}
                                            </span>
                                        </div>
                                        {source.is_active ?
                                            <span className="text-emerald-500 flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest bg-emerald-500/10 px-2 py-0.5 rounded-full whitespace-nowrap flex-shrink-0 flex-nowrap min-w-max">
                                                <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse flex-shrink-0" />
                                                Live
                                            </span> :
                                            <span className="text-text-muted text-[10px] font-bold uppercase tracking-widest bg-white/5 px-2 py-0.5 rounded-full whitespace-nowrap flex-shrink-0 flex-nowrap min-w-max">Paused</span>
                                        }
                                    </div>
                                    <p className="text-xs font-mono text-text-muted flex items-center gap-1">
                                        {source.base_url}
                                    </p>
                                </div>
                                <div className="grid grid-cols-4 gap-4 mt-6 relative z-10 w-full">
                                    <button
                                        onClick={() => handleSyncNow(source.id)}
                                        className="flex flex-col items-center gap-2 p-3 w-full bg-primary/5 hover:bg-primary/20 border border-primary/20 rounded-xl text-primary transition-all active:scale-95 group/btn shadow-sm min-w-0"
                                    >
                                        <Play size={18} fill="currentColor" className="group-hover/btn:scale-110 transition-transform" />
                                        <span className="text-[9px] font-bold uppercase tracking-widest truncate w-full text-center">Sync</span>
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowMappingModal(true); }}
                                        className="flex flex-col items-center gap-2 p-3 w-full bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl text-text-muted hover:text-text transition-all active:scale-95 group/btn shadow-sm min-w-0"
                                    >
                                        <Fingerprint size={18} className="group-hover/btn:rotate-12 transition-transform" />
                                        <span className="text-[9px] font-bold uppercase tracking-widest truncate w-full text-center">Mappings</span>
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowSourceModal(true); }}
                                        className="flex flex-col items-center gap-2 p-3 w-full bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl text-text-muted hover:text-text transition-all active:scale-95 group/btn shadow-sm min-w-0"
                                    >
                                        <Edit size={18} className="group-hover/btn:-rotate-12 transition-transform" />
                                        <span className="text-[9px] font-bold uppercase tracking-widest truncate w-full text-center">Edit</span>
                                    </button>
                                    <button
                                        onClick={() => handleDeleteSource(source.id)}
                                        className="flex flex-col items-center gap-2 p-3 w-full bg-red-500/5 hover:bg-red-500/10 border border-red-500/20 rounded-xl text-red-500/60 hover:text-red-500 transition-all active:scale-95 group/btn shadow-sm min-w-0"
                                    >
                                        <Trash2 size={18} className="group-hover/btn:scale-110 transition-transform" />
                                        <span className="text-[9px] font-bold uppercase tracking-widest truncate w-full text-center">Delete</span>
                                    </button>
                                </div>
                            </div>

                            <div className="mt-6 grid grid-cols-2 gap-x-4 gap-y-6 pt-5 border-t border-white/5">
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-widest text-text-muted font-bold block">Integrity</span>
                                    <div className={`text-xs font-bold flex items-center gap-1.5 ${(!source.last_error) ? 'text-emerald-400' : 'text-red-400'}`}>
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

// Local components if needed

export default IngestManager;
