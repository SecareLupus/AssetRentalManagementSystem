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
    Play
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
            // Toast notification would be nice here
            fetchSources();
        } catch (err) {
            alert("Failed to trigger sync.");
        }
    };

    const handleDelete = async (id) => {
        if (!window.confirm("Delete this source?")) return;
        try {
            await axios.delete(`/v1/admin/ingest/sources/${id}`);
            fetchSources();
        } catch (err) {
            alert("Delete failed.");
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
                <div className="grid gap-4">
                    {sources.map(source => (
                        <GlassCard key={source.id} className="p-4">
                            <div className="flex justify-between items-start">
                                <div className="space-y-1">
                                    <div className="flex items-center gap-3">
                                        <h3 className="font-bold text-lg">{source.name}</h3>
                                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium bg-surface border border-border`}>
                                            {source.target_model}
                                        </span>
                                        {source.is_active ?
                                            <span className="text-emerald-500 flex items-center gap-1 text-xs">
                                                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                                                Active
                                            </span> :
                                            <span className="text-text-muted text-xs">Paused</span>
                                        }
                                    </div>
                                    <p className="text-xs font-mono text-text-muted truncate max-w-md">{source.api_url}</p>
                                </div>
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => handleSyncNow(source.id)}
                                        className="p-2 hover:bg-surface rounded-lg text-primary transition-colors"
                                        title="Sync Now"
                                    >
                                        <Play size={18} />
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowMappingModal(true); }}
                                        className="p-2 hover:bg-surface rounded-lg text-text-muted transition-colors"
                                        title="Configure Mappings"
                                    >
                                        <Fingerprint size={18} />
                                    </button>
                                    <button
                                        onClick={() => { setCurrentSource(source); setShowSourceModal(true); }}
                                        className="p-2 hover:bg-surface rounded-lg text-text-muted transition-colors"
                                    >
                                        <Edit size={18} />
                                    </button>
                                    <button
                                        onClick={() => handleDelete(source.id)}
                                        className="p-2 hover:bg-surface rounded-lg text-red-500 transition-colors"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            </div>

                            <div className="mt-4 grid grid-cols-4 gap-4 pt-4 border-t border-border">
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-wider text-text-muted font-bold">Last Sync</span>
                                    <div className="text-sm flex items-center gap-2">
                                        <Clock size={14} className="text-text-muted" />
                                        {source.last_sync_at ? new time.Time(source.last_sync_at).toLocaleString() : 'Never'}
                                    </div>
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-wider text-text-muted font-bold">Status</span>
                                    <div className={`text-sm flex items-center gap-2 ${source.last_status?.includes('Success') ? 'text-emerald-500' : 'text-amber-500'}`}>
                                        {source.last_status?.includes('Success') ? <CheckCircle size={14} /> : <XCircle size={14} />}
                                        {source.last_status || 'Pending'}
                                    </div>
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-wider text-text-muted font-bold">Mappings</span>
                                    <div className="text-sm flex items-center gap-2">
                                        <Zap size={14} className="text-primary" />
                                        {source.mappings?.length || 0} fields
                                    </div>
                                </div>
                                <div className="space-y-1">
                                    <span className="text-[10px] uppercase tracking-wider text-text-muted font-bold">Next Run</span>
                                    <div className="text-sm opacity-60">
                                        {source.next_sync_at ? new Date(source.next_sync_at).toLocaleTimeString() : 'N/A'}
                                    </div>
                                </div>
                            </div>
                        </GlassCard>
                    ))}
                    {sources.length === 0 && !loading && (
                        <div className="text-center py-20 border-2 border-dashed border-border rounded-xl">
                            <Database className="mx-auto text-text-muted mb-4" size={48} />
                            <p className="text-text-muted">No ingestion sources configured.</p>
                            <button
                                onClick={() => setShowSourceModal(true)}
                                className="mt-4 text-primary hover:underline font-medium"
                            >
                                Get started by adding your first source
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

export default IngestManager;
