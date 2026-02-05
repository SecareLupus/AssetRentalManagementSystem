import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Calendar, Check, X, Clock, AlertCircle, User, ArrowUpRight } from 'lucide-react';

const ReservationsList = () => {
    const [reservations, setReservations] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchReservations = async () => {
        try {
            const response = await axios.get('/v1/rent-actions');
            setReservations(response.data || []);
        } catch (error) {
            console.error("Failed to fetch reservations", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchReservations();
    }, []);

    const handleAction = async (id, action) => {
        try {
            await axios.post(`/v1/rent-actions/${id}/${action}`);
            fetchReservations(); // Refresh
        } catch (error) {
            alert(`Action ${action} failed. ` + (error.response?.data || error.message));
        }
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2.5rem' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '0.5rem' }}>Reservations & Approval</h1>
                <p style={{ color: 'var(--text-muted)' }}>Review and manage equipment assignment requests.</p>
            </header>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                {loading ? (
                    <p style={{ textAlign: 'center', padding: '4rem' }}>Loading reservations...</p>
                ) : reservations.length === 0 ? (
                    <div className="glass" style={{ padding: '4rem', textAlign: 'center', borderRadius: '1rem' }}>
                        <Calendar size={48} style={{ marginBottom: '1rem', opacity: 0.2, margin: '0 auto' }} />
                        <p>No reservations found in the system.</p>
                    </div>
                ) : (
                    reservations.map(res => (
                        <div key={res.id} className="glass" style={{ borderRadius: '1rem', overflow: 'hidden' }}>
                            <div style={{ padding: '1.5rem', display: 'grid', gridTemplateColumns: '1fr 2fr 1.5fr', alignItems: 'center', gap: '2rem' }}>

                                {/* ID & Status */}
                                <div>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}>
                                        <span style={{ fontSize: '0.875rem', fontWeight: 800, color: 'var(--primary)' }}>RES-{res.id}</span>
                                        <StatusBadge status={res.status} />
                                    </div>
                                    <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                        <User size={12} /> {res.requester_ref}
                                    </div>
                                </div>

                                {/* Details */}
                                <div>
                                    <div style={{ fontWeight: 600, marginBottom: '0.25rem' }}>{res.description || 'No description provided'}</div>
                                    <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'flex', gap: '1rem' }}>
                                        <span style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                            <Clock size={12} /> {new Date(res.start_time).toLocaleDateString()}
                                        </span>
                                        <span>&rarr;</span>
                                        <span>{new Date(res.end_time).toLocaleDateString()}</span>
                                        {res.priority === 'high' && <span style={{ color: 'var(--error)', fontWeight: 700 }}>HIGH PRIORITY</span>}
                                    </div>
                                </div>

                                {/* Actions */}
                                <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '0.75rem' }}>
                                    {res.status === 'draft' && (
                                        <button onClick={() => handleAction(res.id, 'submit')} className="btn-primary" style={{ fontSize: '0.75rem' }}>
                                            Submit
                                        </button>
                                    )}
                                    {res.status === 'pending' && (
                                        <>
                                            <button onClick={() => handleAction(res.id, 'reject')} className="glass" style={{ border: '1px solid var(--error)40', color: 'var(--error)', padding: '0.5rem 0.75rem', borderRadius: '0.5rem', fontSize: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                                <X size={14} /> Reject
                                            </button>
                                            <button onClick={() => handleAction(res.id, 'approve')} className="btn-primary" style={{ fontSize: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                                <Check size={14} /> Approve
                                            </button>
                                        </>
                                    )}
                                    {(res.status === 'approved' || res.status === 'pending') && (
                                        <button onClick={() => handleAction(res.id, 'cancel')} style={{ background: 'transparent', color: 'var(--text-muted)', fontSize: '0.75rem' }}>
                                            Cancel
                                        </button>
                                    )}
                                    <button className="glass" style={{ padding: '0.5rem', borderRadius: '0.5rem' }}>
                                        <ArrowUpRight size={14} />
                                    </button>
                                </div>

                            </div>

                            {/* Optional footer for line items if many */}
                            <div style={{ padding: '0.75rem 1.5rem', background: 'rgba(255,255,255,0.02)', borderTop: '1px solid var(--border)', fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                                Items Requested: {res.items?.map(i => `${i.requested_quantity}x Item#${i.item_id}`).join(', ') || 'N/A'}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

const StatusBadge = ({ status }) => {
    const styles = {
        draft: { bg: 'rgba(148, 163, 184, 0.2)', text: 'var(--text-muted)', icon: Clock },
        pending: { bg: 'rgba(245, 158, 11, 0.2)', text: 'var(--warning)', icon: AlertCircle },
        approved: { bg: 'rgba(16, 185, 129, 0.2)', text: 'var(--success)', icon: Check },
        rejected: { bg: 'rgba(239, 104, 104, 0.2)', text: 'var(--error)', icon: X },
        cancelled: { bg: 'rgba(148, 163, 184, 0.1)', text: 'var(--text-muted)', icon: X },
    };
    const style = styles[status] || styles.draft;
    const Icon = style.icon;

    return (
        <div style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: '0.25rem',
            padding: '0.125rem 0.5rem',
            borderRadius: '2rem',
            background: style.bg,
            color: style.text,
            fontSize: '0.65rem',
            fontWeight: 800,
            textTransform: 'uppercase'
        }}>
            <Icon size={10} /> {status}
        </div>
    );
};

export default ReservationsList;
