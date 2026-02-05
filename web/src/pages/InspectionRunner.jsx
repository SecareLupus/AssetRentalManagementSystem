import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { ClipboardCheck, ArrowLeft, Save, CheckCircle2, AlertTriangle, Camera } from 'lucide-react';

const InspectionRunner = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [asset, setAsset] = useState(null);
    const [templates, setTemplates] = useState([]);
    const [responses, setResponses] = useState({}); // { templateId: { status: 'pass/fail', notes: '' } }
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const assetRes = await axios.get(`/v1/inventory/assets/${id}`);
                setAsset(assetRes.data);

                const templRes = await axios.get(`/v1/inventory/assets/${id}/required-inspections`);
                setTemplates(templRes.data || []);
            } catch (err) {
                console.error("Failed to load inspection context", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    const handleResponse = (templId, field, value) => {
        setResponses(prev => ({
            ...prev,
            [templId]: { ...prev[templId], [field]: value }
        }));
    };

    const handleSubmit = async () => {
        setSubmitting(true);
        try {
            const submission = {
                performed_by: 'Inspector Tech',
                responses: Object.entries(responses).map(([templId, resp]) => ({
                    template_id: parseInt(templId),
                    value: resp.status,
                    notes: resp.notes
                }))
            };
            await axios.post(`/v1/inventory/assets/${id}/inspections`, submission);
            navigate('/tech');
        } catch (err) {
            alert("Failed to submit inspection.");
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) return <div style={{ padding: '4rem', textAlign: 'center' }}>Loading...</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '800px', margin: '0 auto' }}>
            <button onClick={() => navigate(-1)} style={{ background: 'transparent', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Cancel
            </button>

            <header style={{ marginBottom: '3rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.75rem' }}>
                    <div style={{ background: 'var(--warning)20', padding: '0.75rem', borderRadius: '0.75rem' }}>
                        <ClipboardCheck size={24} color="var(--warning)" />
                    </div>
                    <h1 style={{ fontSize: '1.75rem', fontWeight: 800 }}>Asset Inspection Runner</h1>
                </div>
                <div style={{ display: 'flex', gap: '1.5rem', color: 'var(--text-muted)', fontSize: '0.875rem' }}>
                    <span>Asset: <strong>{asset.asset_tag || asset.serial_number}</strong></span>
                    <span>Type: {asset.item_type_id}</span>
                </div>
            </header>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem', marginBottom: '3rem' }}>
                {templates.length === 0 ? (
                    <div className="glass" style={{ padding: '2rem', textAlign: 'center', borderRadius: '1rem' }}>
                        No specific inspection templates found for this asset type.
                    </div>
                ) : (
                    templates.map(templ => (
                        <div key={templ.id} className="glass" style={{ padding: '1.5rem', borderRadius: '1rem' }}>
                            <div style={{ marginBottom: '1rem' }}>
                                <h3 style={{ fontSize: '1.125rem', fontWeight: 700 }}>{templ.name}</h3>
                                <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>{templ.category || 'General Check'}</p>
                            </div>

                            <div style={{ display: 'flex', gap: '1rem', marginBottom: '1.5rem' }}>
                                <button
                                    onClick={() => handleResponse(templ.id, 'status', 'pass')}
                                    style={{
                                        flex: 1,
                                        padding: '0.75rem',
                                        borderRadius: '0.5rem',
                                        background: responses[templ.id]?.status === 'pass' ? 'var(--success)' : 'var(--surface)',
                                        color: responses[templ.id]?.status === 'pass' ? 'white' : 'var(--text)',
                                        display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '0.5rem'
                                    }}
                                >
                                    <CheckCircle2 size={18} /> Pass
                                </button>
                                <button
                                    onClick={() => handleResponse(templ.id, 'status', 'fail')}
                                    style={{
                                        flex: 1,
                                        padding: '0.75rem',
                                        borderRadius: '0.5rem',
                                        background: responses[templ.id]?.status === 'fail' ? 'var(--error)' : 'var(--surface)',
                                        color: responses[templ.id]?.status === 'fail' ? 'white' : 'var(--text)',
                                        display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '0.5rem'
                                    }}
                                >
                                    <AlertTriangle size={18} /> Fail
                                </button>
                            </div>

                            <div>
                                <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Inspector Notes</label>
                                <textarea
                                    className="glass"
                                    placeholder="Describe findings or issues..."
                                    style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white', resize: 'none' }}
                                    onChange={(e) => handleResponse(templ.id, 'notes', e.target.value)}
                                />
                            </div>

                            <button
                                className="glass"
                                style={{ marginTop: '1rem', width: '100%', padding: '0.5rem', borderRadius: '0.5rem', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '0.5rem', fontSize: '0.875rem', color: 'var(--text-muted)' }}
                            >
                                <Camera size={16} /> Attach Photo Evidence
                            </button>
                        </div>
                    ))
                )}
            </div>

            <div style={{ borderTop: '1px solid var(--border)', paddingTop: '2rem', display: 'flex', justifyContent: 'flex-end' }}>
                <button
                    className="btn-primary"
                    disabled={submitting || Object.keys(responses).length < templates.length}
                    onClick={handleSubmit}
                    style={{ padding: '0.75rem 2.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
                >
                    <Save size={20} /> {submitting ? 'Submitting...' : 'Complete & Sign Off'}
                </button>
            </div>
        </div>
    );
};

export default InspectionRunner;
