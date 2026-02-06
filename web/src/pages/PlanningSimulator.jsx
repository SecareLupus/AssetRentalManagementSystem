import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Calculator, Plus, Trash2, Calendar, TrendingUp, AlertTriangle, CheckCircle2, Info } from 'lucide-react';

const PlanningSimulator = () => {
    const [scenarios, setScenarios] = useState([]);
    const [itemTypes, setItemTypes] = useState([]);
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState({ totalAssets: 0, availableAssets: 0 });

    const [selectedItem, setSelectedItem] = useState('');
    const [quantity, setQuantity] = useState(1);
    const [startDate, setStartDate] = useState(new Date().toISOString().split('T')[0]);
    const [endDate, setEndDate] = useState(new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [typesRes, statsRes] = await Promise.all([
                    axios.get('/v1/catalog/item-types'),
                    axios.get('/v1/dashboard/stats')
                ]);
                setItemTypes(typesRes.data || []);
                const s = statsRes.data || {};
                setStats({
                    totalAssets: s.total_assets || 0,
                    availableAssets: s.assets_by_status?.available || 0
                });
                if (typesRes.data?.length > 0) setSelectedItem(typesRes.data[0].id);
            } catch (err) {
                console.error("Failed to load simulator data", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    const addScenario = () => {
        const item = itemTypes.find(t => t.id === parseInt(selectedItem));
        const newScenario = {
            id: Date.now(),
            item,
            quantity: parseInt(quantity),
            startDate,
            endDate
        };
        setScenarios([...scenarios, newScenario]);
    };

    const removeScenario = (id) => {
        setScenarios(scenarios.filter(s => s.id !== id));
    };

    const getSimulatedImpact = () => {
        if (scenarios.length === 0) {
            return { impactedAssets: 0, remainingAvailable: stats.availableAssets, utilizationDelta: 0 };
        }

        // To find the peak impact, we check all dates where a scenario starts or ends
        const criticalDates = new Set();
        scenarios.forEach(s => {
            criticalDates.add(s.startDate);
            criticalDates.add(s.endDate);
        });

        const sortedDates = Array.from(criticalDates).sort();
        let peakImpact = 0;

        sortedDates.forEach(date => {
            const currentImpact = scenarios.reduce((acc, s) => {
                if (date >= s.startDate && date <= s.endDate) {
                    return acc + s.quantity;
                }
                return acc;
            }, 0);
            if (currentImpact > peakImpact) peakImpact = currentImpact;
        });

        const remainingAvailable = stats.availableAssets - peakImpact;
        return {
            impactedAssets: peakImpact,
            remainingAvailable,
            utilizationDelta: ((peakImpact / (stats.totalAssets || 1)) * 100).toFixed(1)
        };
    };

    const impact = getSimulatedImpact();

    return (
        <div style={{ padding: '2rem', maxWidth: '1100px', margin: '0 auto' }}>
            <header style={{ marginBottom: '3rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '0.75rem' }}>
                    <div style={{ background: 'var(--primary)20', padding: '0.75rem', borderRadius: '0.75rem' }}>
                        <Calculator size={24} color="var(--primary)" />
                    </div>
                    <h1 style={{ fontSize: '1.75rem', fontWeight: 800 }}>"What-If" Planning Simulator</h1>
                </div>
                <p style={{ color: 'var(--text-muted)' }}>Model high-level fleet movements and see projected impact on availability without committing data.</p>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '3rem' }}>
                {/* Simulator Inputs */}
                <div>
                    <section className="glass" style={{ padding: '2rem', borderRadius: '1.5rem', marginBottom: '2rem' }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Plus size={18} color="var(--primary)" /> Add Movement Scenario
                        </h3>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
                            <div>
                                <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Equipment Type</label>
                                <select
                                    className="glass"
                                    style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                    value={selectedItem}
                                    onChange={(e) => setSelectedItem(e.target.value)}
                                >
                                    {itemTypes.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                                </select>
                            </div>

                            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                                <div>
                                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Quantity</label>
                                    <input
                                        type="number"
                                        className="glass"
                                        min="1"
                                        value={quantity}
                                        onChange={(e) => setQuantity(e.target.value)}
                                        style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }}
                                    />
                                </div>
                                <div style={{ display: 'flex', alignItems: 'flex-end' }}>
                                    <button className="btn-primary" style={{ width: '100%', padding: '0.75rem' }} onClick={addScenario}>
                                        Add to Model
                                    </button>
                                </div>
                            </div>

                            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                                <div>
                                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Start Date</label>
                                    <input type="date" className="glass" value={startDate} onChange={(e) => setStartDate(e.target.value)} style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }} />
                                </div>
                                <div>
                                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>End Date</label>
                                    <input type="date" className="glass" value={endDate} onChange={(e) => setEndDate(e.target.value)} style={{ width: '100%', padding: '0.75rem', borderRadius: '0.5rem', color: 'white' }} />
                                </div>
                            </div>
                        </div>
                    </section>

                    <section className="glass" style={{ borderRadius: '1.5rem', overflow: 'hidden' }}>
                        <div style={{ padding: '1.25rem', borderBottom: '1px solid var(--border)', background: 'rgba(255,255,255,0.02)' }}>
                            <h3 style={{ fontWeight: 600, fontSize: '0.875rem' }}>Active Scenarios</h3>
                        </div>
                        <div style={{ maxHeight: '300px', overflowY: 'auto' }}>
                            {scenarios.length === 0 ? (
                                <p style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.875rem' }}>No scenarios added yet.</p>
                            ) : (
                                scenarios.map(s => (
                                    <div key={s.id} style={{ padding: '1rem 1.25rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <div>
                                            <div style={{ fontWeight: 600, fontSize: '0.875rem' }}>{s.item?.name}</div>
                                            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                                                {s.quantity} units | {s.startDate} to {s.endDate}
                                            </div>
                                        </div>
                                        <button onClick={() => removeScenario(s.id)} style={{ padding: '0.5rem', color: 'var(--text-muted)' }} onMouseOver={(e) => e.target.style.color = 'var(--error)'} onMouseOut={(e) => e.target.style.color = 'var(--text-muted)'}>
                                            <Trash2 size={16} />
                                        </button>
                                    </div>
                                ))
                            )}
                        </div>
                    </section>
                </div>

                {/* Simulation Output */}
                <div>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem', marginBottom: '2rem' }}>
                        <div className="glass" style={{ padding: '1.5rem', borderRadius: '1.5rem' }}>
                            <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Baseline Availability</span>
                            <div style={{ fontSize: '1.5rem', fontWeight: 800 }}>{stats.availableAssets}</div>
                        </div>
                        <div className="glass" style={{ padding: '1.5rem', borderRadius: '1.5rem', border: impact.remainingAvailable < 0 ? '1px solid var(--error)50' : '1px solid var(--border)' }}>
                            <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Simulated Remainder</span>
                            <div style={{ fontSize: '1.5rem', fontWeight: 800, color: impact.remainingAvailable < 0 ? 'var(--error)' : 'inherit' }}>{impact.remainingAvailable}</div>
                        </div>
                    </div>

                    <section className="glass" style={{ padding: '2rem', borderRadius: '1.5rem', minHeight: '300px' }}>
                        <h3 style={{ fontWeight: 600, marginBottom: '2rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <TrendingUp size={18} color="var(--primary)" /> Projection Analysis
                        </h3>

                        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                            <div>
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem', fontSize: '0.875rem' }}>
                                    <span>Simulated Fleet Utilization</span>
                                    <span style={{ fontWeight: 700 }}>{impact.utilizationDelta}% increase</span>
                                </div>
                                <div style={{ height: '8px', background: 'var(--surface)', borderRadius: '4px', overflow: 'hidden' }}>
                                    <div style={{ width: `${Math.min(100, impact.utilizationDelta)}%`, height: '100%', background: 'var(--primary)', boxShadow: '0 0 10px var(--primary)40' }} />
                                </div>
                            </div>

                            {impact.remainingAvailable < 0 ? (
                                <div style={{ background: 'var(--error)10', border: '1px solid var(--error)20', padding: '1rem', borderRadius: '1rem', display: 'flex', gap: '1rem' }}>
                                    <AlertTriangle color="var(--error)" size={24} />
                                    <div>
                                        <div style={{ fontWeight: 700, color: 'var(--error)', fontSize: '0.875rem' }}>Critical Availability Alert</div>
                                        <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '0.25rem' }}>Simulated movements exceed current available inventory by {Math.abs(impact.remainingAvailable)} units.</p>
                                    </div>
                                </div>
                            ) : scenarios.length > 0 ? (
                                <div style={{ background: 'var(--success)10', border: '1px solid var(--success)20', padding: '1rem', borderRadius: '1rem', display: 'flex', gap: '1rem' }}>
                                    <CheckCircle2 color="var(--success)" size={24} />
                                    <div>
                                        <div style={{ fontWeight: 700, color: 'var(--success)', fontSize: '0.875rem' }}>Fleet Capacity Sufficient</div>
                                        <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '0.25rem' }}>Simulated movements are sustainable with current inventory levels.</p>
                                    </div>
                                </div>
                            ) : (
                                <div style={{ background: 'rgba(255,255,255,0.02)', border: '1px solid var(--border)', padding: '1rem', borderRadius: '1rem', display: 'flex', gap: '1rem' }}>
                                    <Info color="var(--text-muted)" size={24} />
                                    <div>
                                        <div style={{ fontWeight: 700, color: 'var(--text-muted)', fontSize: '0.875rem' }}>Waiting for scenarios</div>
                                        <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '0.25rem' }}>Add movement scenarios to see projected fleet impact.</p>
                                    </div>
                                </div>
                            )}

                            <div style={{ borderTop: '1px solid var(--border)', paddingTop: '1.5rem' }}>
                                <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontStyle: 'italic' }}>
                                    * This simulation uses current real-time inventory baseline. It does not account for future returns or existing maintenance queues unless they are explicitly modeled as positive availability offsets.
                                </p>
                            </div>
                        </div>
                    </section>
                </div>
            </div>
        </div>
    );
};

export default PlanningSimulator;
