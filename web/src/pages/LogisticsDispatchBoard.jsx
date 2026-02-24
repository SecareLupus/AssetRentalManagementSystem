import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
    Truck,
    Calendar,
    Package,
    ChevronRight,
    Plus,
    ArrowRight,
    Filter,
    Search,
    Clock,
    CheckCircle2,
    Briefcase,
    Box
} from 'lucide-react';
import ShipmentAllocationUI from '../components/ShipmentAllocationUI';
import { StatusBadge } from '../components/Shared';

const LogisticsDispatchBoard = () => {
    const [reservations, setReservations] = useState([]);
    const [deliveries, setDeliveries] = useState([]); // ScheduledDeliveries
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState('unassigned'); // unassigned, scheduled, shipped
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedShipmentId, setSelectedShipmentId] = useState(null);
    const [showAllocationModal, setShowAllocationModal] = useState(false);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [resRes, delRes] = await Promise.all([
                axios.get('/v1/logistics/reservations'),
                axios.get('/v1/logistics/deliveries')
            ]);
            // Only confirmed reservations that need dispatch
            setReservations(resRes.data?.filter(r => r.reservationStatus === 'ReservationConfirmed') || []);
            setDeliveries(delRes.data || []);
        } catch (err) {
            console.error("Failed to load dispatch data", err);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateDelivery = async (reservation) => {
        try {
            const deliveryData = {
                eventId: reservation.id, // Linking to reservation ID for now if no specific event
                availableFrom: reservation.startTime,
                notes: `Consolidated delivery for ${reservation.reservationName || 'Reservation #' + reservation.id}`
            };
            await axios.post('/v1/logistics/deliveries', deliveryData);
            fetchData();
        } catch (err) {
            console.error("Failed to create delivery", err);
        }
    };

    const filteredReservations = reservations.filter(r =>
    (r.reservationName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        r.id.toString().includes(searchTerm))
    );

    return (
        <div style={{ padding: '2rem', maxWidth: '1400px', margin: '0 auto' }}>
            <header style={{ marginBottom: '3rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end' }}>
                <div>
                    <div style={{ display: 'inline-flex', background: 'var(--primary)', padding: '0.75rem', borderRadius: '0.75rem', marginBottom: '1rem' }}>
                        <Truck size={24} color="white" />
                    </div>
                    <h1 style={{ fontSize: '2.5rem', fontWeight: 900, marginBottom: '0.5rem' }}>Logistics Dispatch Board</h1>
                    <p style={{ color: 'var(--text-muted)' }}>Manage outbound fulfillment, group reservations into deliveries, and track shipments.</p>
                </div>

                <div style={{ display: 'flex', gap: '1rem' }}>
                    <div className="glass" style={{ display: 'flex', alignItems: 'center', padding: '0.5rem 1rem', borderRadius: '0.75rem', gap: '0.5rem' }}>
                        <Search size={18} color="var(--text-muted)" />
                        <input
                            type="text"
                            placeholder="Search reservations..."
                            style={{ background: 'transparent', border: 'none', color: 'white', outline: 'none', width: '200px' }}
                            value={searchTerm}
                            onChange={e => setSearchTerm(e.target.value)}
                        />
                    </div>
                </div>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '350px 1fr', gap: '2rem' }}>
                {/* Pending Reservations Queue */}
                <aside className="glass" style={{ borderRadius: '1.5rem', display: 'flex', flexDirection: 'column', height: 'calc(100vh - 250px)' }}>
                    <div style={{ padding: '1.5rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <h3 style={{ fontWeight: 800, fontSize: '1.1rem' }}>Pending Dispatch</h3>
                        <span style={{ fontSize: '0.7rem', fontWeight: 900, background: 'var(--primary)', padding: '0.2rem 0.5rem', borderRadius: '1rem' }}>
                            {filteredReservations.length}
                        </span>
                    </div>

                    <div style={{ flex: 1, overflowY: 'auto', padding: '1rem' }}>
                        {loading ? (
                            <div style={{ textAlign: 'center', padding: '2rem', opacity: 0.5 }}>Loading...</div>
                        ) : filteredReservations.length === 0 ? (
                            <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-muted)', opacity: 0.5 }}>
                                <Clock size={32} style={{ marginBottom: '1rem' }} />
                                <p>No confirmed reservations awaiting dispatch.</p>
                            </div>
                        ) : (
                            filteredReservations.map(res => (
                                <ReservationMiniCard
                                    key={res.id}
                                    reservation={res}
                                    onSchedule={() => handleCreateDelivery(res)}
                                />
                            ))
                        )}
                    </div>
                </aside>

                {/* Planned Deliveries & Shipments */}
                <main style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                    <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
                        <TabButton active={activeTab === 'unassigned'} onClick={() => setActiveTab('unassigned')}>Active Deliveries</TabButton>
                        <TabButton active={activeTab === 'shipped'} onClick={() => setActiveTab('shipped')}>In Transit</TabButton>
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(400px, 1fr))', gap: '1.5rem' }}>
                        {deliveries.length === 0 ? (
                            <div className="glass" style={{ gridColumn: '1/-1', padding: '4rem', textAlign: 'center', borderRadius: '1.5rem', opacity: 0.5 }}>
                                <Package size={48} style={{ marginBottom: '1rem' }} />
                                <h3>No active delivery events.</h3>
                                <p>Schedule a reservation from the sidebar to start.</p>
                            </div>
                        ) : (
                            deliveries.map(del => (
                                <DeliveryCard
                                    key={del.id}
                                    delivery={del}
                                    onAllocate={(shipmentId) => {
                                        setSelectedShipmentId(shipmentId);
                                        setShowAllocationModal(true);
                                    }}
                                />
                            ))
                        )}
                    </div>
                </main>
            </div>

            {showAllocationModal && (
                <ShipmentAllocationUI
                    shipmentId={selectedShipmentId}
                    onClose={() => setShowAllocationModal(false)}
                    onAllocationComplete={() => {
                        setShowAllocationModal(false);
                        fetchData();
                    }}
                />
            )}
        </div>
    );
};

const ReservationMiniCard = ({ reservation, onSchedule }) => (
    <div className="glass" style={{
        padding: '1.25rem',
        borderRadius: '1rem',
        marginBottom: '1rem',
        border: '1px solid rgba(255,255,255,0.05)',
        transition: 'transform 0.2s',
        cursor: 'pointer'
    }} onMouseEnter={e => e.currentTarget.style.transform = 'translateY(-2px)'}
        onMouseLeave={e => e.currentTarget.style.transform = 'translateY(0)'}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '0.75rem' }}>
            <div style={{ fontWeight: 800, fontSize: '0.9rem' }}>{reservation.reservationName || `RES-${reservation.id}`}</div>
            <span style={{ fontSize: '0.7rem', color: 'var(--success)', fontWeight: 800 }}>CONFIRMED</span>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', fontSize: '0.75rem', marginBottom: '1rem' }}>
            <Calendar size={14} />
            <span>{new Date(reservation.startTime).toLocaleDateString()}</span>
        </div>

        <button
            onClick={onSchedule}
            style={{
                width: '100%',
                padding: '0.5rem',
                borderRadius: '0.5rem',
                background: 'var(--surface)',
                fontSize: '0.75rem',
                fontWeight: 700,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '0.5rem',
                border: '1px solid var(--border)'
            }}
        >
            Schedule Delivery <Plus size={14} />
        </button>
    </div>
);

const DeliveryCard = ({ delivery, onAllocate }) => {
    const [shipments, setShipments] = useState([]);

    useEffect(() => {
        if (delivery.id) {
            axios.get(`/v1/logistics/shipments?delivery_id=${delivery.id}`)
                .then(res => setShipments(res.data || []))
                .catch(err => console.error(err));
        }
    }, [delivery.id]);

    return (
        <div className="glass" style={{ borderRadius: '1.5rem', padding: '1.5rem', border: '1px solid rgba(255,255,255,0.1)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                    <div style={{ background: 'rgba(99, 102, 241, 0.1)', padding: '0.5rem', borderRadius: '0.5rem', color: 'var(--primary)' }}>
                        <Calendar size={20} />
                    </div>
                    <div>
                        <div style={{ fontWeight: 800 }}>Delivery #{delivery.id}</div>
                        <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Target: {new Date(delivery.availableFrom).toLocaleDateString()}</div>
                    </div>
                </div>
                <button className="glass" style={{ padding: '0.4rem 0.8rem', borderRadius: '0.5rem', fontSize: '0.7rem', fontWeight: 700 }}>
                    Details
                </button>
            </div>

            <div style={{ background: 'rgba(255,255,255,0.02)', borderRadius: '1rem', padding: '1rem', marginBottom: '1.5rem' }}>
                <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '1rem' }}>{delivery.notes}</p>

                {shipments.length > 0 && (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                        {shipments.map(s => (
                            <div key={s.id} className="glass" style={{ padding: '0.75rem', borderRadius: '0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.75rem' }}>
                                    <Box size={14} color="var(--primary)" />
                                    <span>Shipment #{s.id}</span>
                                    <StatusBadge status={s.status} />
                                </div>
                                <button
                                    className="glass"
                                    style={{ padding: '0.25rem 0.5rem', fontSize: '0.65rem' }}
                                    onClick={() => onAllocate(s.id)}
                                >
                                    Allocate
                                </button>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div style={{ display: 'flex', gap: '0.75rem' }}>
                <button style={{
                    flex: 1,
                    padding: '0.75rem',
                    borderRadius: '0.75rem',
                    background: 'var(--primary)',
                    color: 'white',
                    fontWeight: 700,
                    fontSize: '0.85rem',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '0.5rem'
                }}>
                    Add Shipment <Plus size={16} />
                </button>
                <button className="glass" style={{ flex: 1, padding: '0.75rem', borderRadius: '0.75rem', fontWeight: 700, fontSize: '0.85rem' }}>
                    Assign Truck
                </button>
            </div>
        </div>
    );
};

const TabButton = ({ active, onClick, children }) => (
    <button
        onClick={onClick}
        style={{
            padding: '0.6rem 1.25rem',
            borderRadius: '1rem',
            fontSize: '0.85rem',
            fontWeight: 700,
            background: active ? 'var(--primary)' : 'var(--surface)',
            color: active ? 'white' : 'var(--text-muted)',
            border: 'none',
            transition: 'all 0.2s'
        }}
    >
        {children}
    </button>
);

export default LogisticsDispatchBoard;
