import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Calendar, Package, FileText, CheckCircle, ChevronRight, ChevronLeft, Plus, Trash2 } from 'lucide-react';

const ReservationWizard = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [catalog, setCatalog] = useState([]);

    const [formData, setFormData] = useState({
        requester_ref: user?.username || 'Unknown User',
        created_by_ref: user?.username || 'Unknown User',
        priority: 'normal',
        start_time: '',
        end_time: '',
        is_asap: false,
        description: '',
        items: [] // { item_kind: 'item_type', item_id: ID, requested_quantity: 1, name: '' }
    });

    useEffect(() => {
        axios.get('/v1/catalog/item-types').then(res => setCatalog(res.data || []));
    }, []);

    const addItem = (item) => {
        setFormData(prev => ({
            ...prev,
            items: [...prev.items, { item_kind: 'item_type', item_id: item.id, requested_quantity: 1, name: item.name }]
        }));
    };

    const removeItem = (index) => {
        setFormData(prev => ({
            ...prev,
            items: prev.items.filter((_, i) => i !== index)
        }));
    };

    const updateItemQty = (index, qty) => {
        const newItems = [...formData.items];
        newItems[index].requested_quantity = parseInt(qty);
        setFormData({ ...formData, items: newItems });
    };

    const handleSubmit = async () => {
        setLoading(true);
        try {
            // Format dates for Go
            const payload = {
                ...formData,
                start_time: new Date(formData.start_time).toISOString(),
                end_time: new Date(formData.end_time).toISOString(),
                status: 'draft'
            };
            const res = await axios.post('/v1/rent-actions', payload);
            // Automatically submit if user wants
            await axios.post(`/v1/rent-actions/${res.data.id}/submit`);
            navigate('/reservations');
        } catch (err) {
            console.error("Submission failed", err);
            alert("Failed to create reservation. Check the API Inspector for details.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '800px', margin: '0 auto' }}>
            <header style={{ marginBottom: '3rem', textAlign: 'center' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '1.5rem' }}>Reservation Wizard</h1>

                {/* Progress Stepper */}
                <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '1rem' }}>
                    <StepIcon active={step >= 1} current={step === 1} icon={Calendar} label="Dates" />
                    <div style={{ width: '40px', height: '1px', background: 'var(--border)' }} />
                    <StepIcon active={step >= 2} current={step === 2} icon={Package} label="Items" />
                    <div style={{ width: '40px', height: '1px', background: 'var(--border)' }} />
                    <StepIcon active={step >= 3} current={step === 3} icon={FileText} label="Details" />
                    <div style={{ width: '40px', height: '1px', background: 'var(--border)' }} />
                    <StepIcon active={step >= 4} current={step === 4} icon={CheckCircle} label="Review" />
                </div>
            </header>

            <div className="glass" style={{ padding: '2.5rem', borderRadius: '1.5rem', minHeight: '400px', display: 'flex', flexDirection: 'column' }}>
                {step === 1 && (
                    <div className="animate-in fade-in slide-in-from-bottom-4">
                        <h2 style={{ marginBottom: '1.5rem' }}>Select Deployment Window</h2>
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
                            <div>
                                <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem' }}>Start Date</label>
                                <input type="datetime-local" className="glass" style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                    value={formData.start_time} onChange={e => setFormData({ ...formData, start_time: e.target.value })}
                                />
                            </div>
                            <div>
                                <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem' }}>End Date</label>
                                <input type="datetime-local" className="glass" style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                    value={formData.end_time} onChange={e => setFormData({ ...formData, end_time: e.target.value })}
                                />
                            </div>
                        </div>
                        <div style={{ marginTop: '2rem' }}>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', cursor: 'pointer' }}>
                                <input type="checkbox" checked={formData.is_asap} onChange={e => setFormData({ ...formData, is_asap: e.target.checked })} />
                                <span>Mark as ASAP Priority</span>
                            </label>
                        </div>
                    </div>
                )}

                {step === 2 && (
                    <div className="animate-in fade-in slide-in-from-bottom-4">
                        <h2 style={{ marginBottom: '1.5rem' }}>Add Equipment</h2>
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem' }}>
                            <div style={{ borderRight: '1px solid var(--border)', paddingRight: '2rem' }}>
                                <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)', marginBottom: '1rem' }}>Catalog Items</p>
                                <div style={{ maxHeight: '300px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                                    {catalog.map(it => (
                                        <button key={it.id} onClick={() => addItem(it)} style={{ textAlign: 'left', padding: '0.75rem', background: 'var(--surface)', borderRadius: '0.5rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                            <span>{it.name}</span>
                                            <Plus size={16} color="var(--primary)" />
                                        </button>
                                    ))}
                                </div>
                            </div>
                            <div>
                                <p style={{ fontSize: '0.875rem', color: 'var(--text-muted)', marginBottom: '1rem' }}>Selected ({formData.items.length})</p>
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                                    {formData.items.map((item, i) => (
                                        <div key={i} style={{ padding: '0.75rem', background: 'var(--surface)', borderRadius: '0.5rem', border: '1px solid var(--primary)30' }}>
                                            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                                                <span style={{ fontWeight: 600 }}>{item.name}</span>
                                                <button onClick={() => removeItem(i)} style={{ background: 'transparent', color: 'var(--error)' }}><Trash2 size={14} /></button>
                                            </div>
                                            <input type="number" min="1" value={item.requested_quantity} onChange={e => updateItemQty(i, e.target.value)} style={{ width: '60px', padding: '0.25rem', background: 'black', color: 'white', border: 'none', borderRadius: '4px' }} />
                                        </div>
                                    ))}
                                    {formData.items.length === 0 && <p style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-muted)', fontSize: '0.875rem' }}>Cart is empty</p>}
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {step === 3 && (
                    <div className="animate-in fade-in slide-in-from-bottom-4">
                        <h2 style={{ marginBottom: '1.5rem' }}>Logistics & Context</h2>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem' }}>Description / Use Case</label>
                            <textarea
                                rows="4"
                                className="glass"
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white', resize: 'none' }}
                                placeholder="Explain the project or specific requirements..."
                                value={formData.description}
                                onChange={e => setFormData({ ...formData, description: e.target.value })}
                            />
                        </div>
                        <div>
                            <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem' }}>Priority</label>
                            <select
                                className="glass"
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                value={formData.priority}
                                onChange={e => setFormData({ ...formData, priority: e.target.value })}
                            >
                                <option value="low">Low</option>
                                <option value="normal">Normal</option>
                                <option value="high">High</option>
                                <option value="urgent">Urgent</option>
                            </select>
                        </div>
                    </div>
                )}

                {step === 4 && (
                    <div className="animate-in fade-in slide-in-from-bottom-4">
                        <h2 style={{ marginBottom: '1.5rem' }}>Review Summary</h2>
                        <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem', background: 'rgba(255,255,255,0.03)' }}>
                            <div style={{ marginBottom: '1rem' }}><strong>Window:</strong> {new Date(formData.start_time).toLocaleString()} - {new Date(formData.end_time).toLocaleString()}</div>
                            <div style={{ marginBottom: '1rem' }}><strong>Items:</strong> {formData.items.map(i => `${i.requested_quantity}x ${i.name}`).join(', ')}</div>
                            <div style={{ marginBottom: '1rem' }}><strong>Priority:</strong> <span style={{ textTransform: 'capitalize' }}>{formData.priority}</span></div>
                            <div><strong>Notes:</strong> {formData.description || 'None'}</div>
                        </div>
                        <div style={{ marginTop: '2rem', padding: '1rem', borderRadius: '0.5rem', background: 'var(--primary)10', color: 'var(--primary)', fontSize: '0.875rem' }}>
                            By clicking "Complete", your request will be sent to a Fleet Manager for approval.
                        </div>
                    </div>
                )}

                {/* Navigation */}
                <div style={{ marginTop: 'auto', paddingTop: '2rem', borderTop: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between' }}>
                    <button
                        onClick={() => setStep(step - 1)}
                        disabled={step === 1}
                        style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', background: 'transparent', color: 'var(--text)', opacity: step === 1 ? 0 : 1 }}
                    >
                        <ChevronLeft size={20} /> Previous
                    </button>

                    {step < 4 ? (
                        <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }} onClick={() => setStep(step + 1)}>
                            Next Step <ChevronRight size={20} />
                        </button>
                    ) : (
                        <button className="btn-primary" disabled={loading} style={{ padding: '0.75rem 2rem' }} onClick={handleSubmit}>
                            {loading ? 'Processing...' : 'Complete & Request Approval'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

const StepIcon = ({ active, current, icon: Icon, label }) => (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '0.5rem', opacity: active ? 1 : 0.3 }}>
        <div style={{
            width: '40px',
            height: '40px',
            borderRadius: '50%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: current ? 'var(--primary)' : active ? 'var(--primary)40' : 'var(--surface)',
            border: current ? 'none' : '1px solid var(--border)'
        }}>
            <Icon size={18} />
        </div>
        <span style={{ fontSize: '0.75rem', fontWeight: 600 }}>{label}</span>
    </div>
);

export default ReservationWizard;
