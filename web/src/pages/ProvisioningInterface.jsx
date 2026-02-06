import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { Zap, ArrowLeft, Play, CheckCircle, Save, Settings, Database, Activity } from 'lucide-react';

const ProvisioningInterface = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useAuth();
    const [asset, setAsset] = useState(null);
    const [buildSpecs, setBuildSpecs] = useState([]);
    const [selectedSpec, setSelectedSpec] = useState('');
    const [loading, setLoading] = useState(true);
    const [provisioning, setProvisioning] = useState(false);
    const [log, setLog] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const assetRes = await axios.get(`/v1/inventory/assets/${id}`);
                setAsset(assetRes.data);

                const specsRes = await axios.get('/v1/fleet/build-specs');
                setBuildSpecs(specsRes.data || []);
                if (specsRes.data?.length > 0) setSelectedSpec(specsRes.data[0].id);
            } catch (err) {
                console.error("Failed to load provisioning context", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    const startProvisioning = async () => {
        setProvisioning(true);
        addLog("Initializing provisioning context...");
        try {
            await axios.post(`/v1/inventory/assets/${id}/provision`, {
                build_spec_id: parseInt(selectedSpec),
                performed_by: user?.username || 'Tech Station #1'
            });
            addLog("Build spec assigned. Status polling active.");
        } catch (err) {
            addLog("ERROR: API rejected provisioning start.");
            setProvisioning(false);
        }
    };

    // Polling for real-time status
    useEffect(() => {
        let interval;
        if (provisioning) {
            interval = setInterval(async () => {
                try {
                    const res = await axios.get(`/v1/inventory/assets/${id}`);
                    const newAsset = res.data;
                    setAsset(newAsset);

                    if (newAsset.status === 'available') {
                        addLog("System ready for final check.");
                        setProvisioning(false);
                        clearInterval(interval);
                    } else if (newAsset.status === 'maintenance') {
                         addLog("Provisioning halted: Manual intervention required.");
                         setProvisioning(false);
                         clearInterval(interval);
                    } else {
                         addLog(`Device Status: ${newAsset.status}...`);
                    }
                } catch (err) {
                    addLog("Communication error while polling status.");
                }
            }, 3000);
        }
        return () => clearInterval(interval);
    }, [provisioning, id]);

    const completeProvisioning = async () => {
        try {
            await axios.post(`/v1/inventory/assets/${id}/complete-provisioning`, {
                notes: "Provisioned successfully via Tech Station Console."
            });
            navigate('/tech');
        } catch (err) {
            alert("Failed to complete provisioning.");
        }
    };

    const addLog = (msg) => {
        setLog(prev => [...prev, { time: new Date().toLocaleTimeString(), msg }]);
    };

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Loading...</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '900px', margin: '0 auto' }}>
            <button onClick={() => navigate(-1)} style={{ background: 'transparent', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Back
            </button>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem' }}>
                {/* Configuration Side */}
                <div>
                    <header style={{ marginBottom: '2rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.75rem' }}>
                            <Zap size={24} color="var(--primary)" />
                            <h1 style={{ fontSize: '1.75rem', fontWeight: 800 }}>Asset Provisioning</h1>
                        </div>
                        <p style={{ color: 'var(--text-muted)' }}>Configure and flash device: <strong>{asset.asset_tag}</strong></p>
                    </header>

                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', marginBottom: '1.5rem' }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '1rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Settings size={18} /> Configuration
                        </h3>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Build Specification</label>
                            <select
                                className="glass"
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                value={selectedSpec}
                                onChange={(e) => setSelectedSpec(e.target.value)}
                                disabled={provisioning}
                            >
                                {buildSpecs.map(spec => (
                                    <option key={spec.id} value={spec.id}>{spec.name} ({spec.version})</option>
                                ))}
                            </select>
                        </div>

                        <button
                            className="btn-primary"
                            style={{ width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '0.5rem' }}
                            onClick={startProvisioning}
                            disabled={provisioning}
                        >
                            <Play size={18} /> {provisioning ? 'Provisioning In Progress...' : 'Start Flashing Sequence'}
                        </button>
                    </div>

                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', opacity: provisioning ? 1 : 0.4 }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '1rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <CheckCircle size={18} /> Quality Sign-off
                        </h3>
                        <button
                            className="btn-primary"
                            style={{ width: '100%', background: 'var(--success)' }}
                            disabled={!provisioning}
                            onClick={completeProvisioning}
                        >
                            Finalize & Set Available
                        </button>
                    </div>
                </div>

                {/* Status Side */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <div className="glass" style={{ flex: 1, borderRadius: '1rem', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
                        <div style={{ padding: '1rem', background: 'rgba(255,255,255,0.05)', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between' }}>
                            <span style={{ fontSize: '0.75rem', fontWeight: 700, textTransform: 'uppercase' }}>Console Log</span>
                            <Activity size={14} color="var(--primary)" />
                        </div>
                        <div style={{ flex: 1, padding: '1rem', background: '#0a0a0a', fontFamily: 'monospace', fontSize: '0.75rem', overflowY: 'auto' }}>
                            {log.length === 0 && <span style={{ color: '#444' }}>Waiting for sequence start...</span>}
                            {log.map((l, i) => (
                                <div key={i} style={{ marginBottom: '0.25rem' }}>
                                    <span style={{ color: 'var(--primary)' }}>[{l.time}]</span> {l.msg}
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
                            <Database size={18} color="var(--text-muted)" />
                            <span style={{ fontWeight: 600 }}>Device Identity</span>
                        </div>
                        <div style={{ fontSize: '0.875rem' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                                <span style={{ color: 'var(--text-muted)' }}>Serial:</span>
                                <span>{asset.serial_number}</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                                <span style={{ color: 'var(--text-muted)' }}>Remote ID:</span>
                                <span>{asset.remote_management_id || 'unassigned'}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProvisioningInterface;
