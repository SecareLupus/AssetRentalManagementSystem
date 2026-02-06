import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Calendar as CalendarIcon, AlertTriangle, TrendingDown, Info, ChevronLeft, ChevronRight } from 'lucide-react';

const AvailabilityHeatmap = () => {
    const [itemTypes, setItemTypes] = useState([]);
    const [selectedItemType, setSelectedItemType] = useState('');
    const [timeline, setTimeline] = useState([]);
    const [loading, setLoading] = useState(false);

    // View range: 14 days from today
    const [startDate, setStartDate] = useState(new Date().toISOString().split('T')[0]);

    useEffect(() => {
        axios.get('/v1/catalog/item-types').then(res => {
            setItemTypes(res.data || []);
            if (res.data?.length > 0) setSelectedItemType(res.data[0].id);
        });
    }, []);

    useEffect(() => {
        if (!selectedItemType) return;

        const fetchTimeline = async () => {
            setLoading(true);
            try {
                const end = new Date(startDate);
                end.setDate(end.getDate() + 14);
                const res = await axios.get(`/v1/intelligence/availability?item_type_id=${selectedItemType}&start=${startDate}&end=${end.toISOString().split('T')[0]}`);
                setTimeline(res.data || []);
            } catch (err) {
                console.error("Failed to fetch timeline", err);
            } finally {
                setLoading(false);
            }
        };
        fetchTimeline();
    }, [selectedItemType, startDate]);

    const changeStartDate = (days) => {
        const current = new Date(startDate);
        current.setDate(current.getDate() + days);
        setStartDate(current.toISOString().split('T')[0]);
    };

    return (
        <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2.5rem' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 800, marginBottom: '0.5rem' }}>Fleet Intelligence</h1>
                <p style={{ color: 'var(--text-muted)' }}>Proactive availability tracking and shortage detection.</p>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(250px, 1fr) 3fr', gap: '2rem' }}>
                {/* Controls */}
                <aside>
                    <div className="glass" style={{ padding: '1.5rem', borderRadius: '1rem' }}>
                        <h3 style={{ fontWeight: 700, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Info size={18} color="var(--primary)" /> Insights Panel
                        </h3>

                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.5rem' }}>Analyze Equipment</label>
                            <select
                                className="glass"
                                style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                value={selectedItemType}
                                onChange={(e) => setSelectedItemType(e.target.value)}
                            >
                                {itemTypes.map(it => <option key={it.id} value={it.id}>{it.name}</option>)}
                            </select>
                        </div>

                        <div style={{ padding: '1rem', borderRadius: '0.5rem', background: 'var(--surface)', fontSize: '0.875rem' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}>
                                <TrendingDown size={14} color="var(--error)" />
                                <span style={{ fontWeight: 600 }}>Bottleneck Alert</span>
                            </div>
                            <p style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>
                                Inventory for this item is projected to reach critical levels in the next 10 days.
                            </p>
                        </div>
                    </div>
                </aside>

                {/* Heatmap Area */}
                <section>
                    <div className="glass" style={{ padding: '2rem', borderRadius: '1rem' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
                            <h3 style={{ fontWeight: 700, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                <CalendarIcon size={20} /> Availability Heatmap
                            </h3>
                            <div style={{ display: 'flex', gap: '0.5rem' }}>
                                <button onClick={() => changeStartDate(-7)} className="glass" style={{ padding: '0.5rem', borderRadius: '0.5rem' }}><ChevronLeft size={16} /></button>
                                <button onClick={() => changeStartDate(7)} className="glass" style={{ padding: '0.5rem', borderRadius: '0.5rem' }}><ChevronRight size={16} /></button>
                            </div>
                        </div>

                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: '1rem' }}>
                            {loading ? <p>Calculating nodes...</p> : timeline.map((p, i) => (
                                <div key={i} className="glass" style={{
                                    height: '100px',
                                    borderRadius: '0.75rem',
                                    padding: '0.75rem',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    justifyContent: 'space-between',
                                    background: p.available <= 0 ? 'rgba(239, 68, 68, 0.1)' : 'rgba(255,255,255,0.02)',
                                    border: p.available <= 0 ? '1px solid var(--error)40' : '1px solid var(--border)'
                                }}>
                                    <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                                        {new Date(p.date).toLocaleDateString('en-US', { weekday: 'short', day: 'numeric' })}
                                    </span>
                                    <div style={{ textAlign: 'right' }}>
                                        <div style={{ fontSize: '1.25rem', fontWeight: 800, color: p.available <= 0 ? 'var(--error)' : 'var(--text)' }}>
                                            {p.available}
                                        </div>
                                        <div style={{ fontSize: '0.65rem', color: 'var(--text-muted)' }}>of {p.total} available</div>
                                    </div>
                                </div>
                            ))}
                        </div>

                        <div style={{ marginTop: '2rem', display: 'flex', gap: '2rem', fontSize: '0.75rem' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                <div style={{ width: '12px', height: '12px', background: 'var(--border)', borderRadius: '2px' }} />
                                <span>Healthy Inventory</span>
                            </div>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                <div style={{ width: '12px', height: '12px', background: 'var(--error)', opacity: 0.3, borderRadius: '2px' }} />
                                <span>Shortage / Over-reserved</span>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
};

export default AvailabilityHeatmap;
