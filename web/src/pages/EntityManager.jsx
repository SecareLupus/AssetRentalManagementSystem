import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {
    Building2, MapPin, Calendar, Plus, Search,
    Edit2, ChevronRight, X, ArrowLeft,
    Box, Layers, Globe, Clock, Users, Mail, Phone, UserPlus
} from 'lucide-react';
import { GlassCard, PageHeader, StatusBadge, Modal } from '../components/Shared';

const EntityManager = () => {
    const [activeTab, setActiveTab] = useState('companies');
    const [companies, setCompanies] = useState([]);
    const [places, setPlaces] = useState([]);
    const [events, setEvents] = useState([]);
    const [people, setPeople] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    // Hierarchical Navigation State for Places
    const [viewingEntity, setViewingEntity] = useState(null); // { type: 'place'|'event', item: Object }
    const [navigationStack, setNavigationStack] = useState([]);
    const [childEntities, setChildEntities] = useState([]); // Sub-places or Event Needs
    const [childLoading, setChildLoading] = useState(false);

    // Modal states
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [modalMode, setModalMode] = useState('create'); // 'create' or 'edit'
    const [modalType, setModalType] = useState(''); // 'company', 'site', 'event', 'location', 'need', 'contact'
    const [selectedItem, setSelectedItem] = useState(null);
    const [formData, setFormData] = useState({});

    // Catalog for asset needs
    const [itemTypes, setItemTypes] = useState([]);

    const fetchData = async () => {
        setLoading(true);
        try {
            if (activeTab === 'companies') {
                const res = await axios.get('/v1/entities/companies');
                setCompanies(res.data || []);
            } else if (activeTab === 'places') {
                const res = await axios.get('/v1/entities/places');
                // Filter for top-level places in main view if needed, or show all
                setPlaces(res.data || []);
            } else if (activeTab === 'events') {
                const res = await axios.get('/v1/entities/events');
                setEvents(res.data || []);
            } else if (activeTab === 'people') {
                const res = await axios.get('/v1/entities/people');
                setPeople(res.data || []);
            }
        } catch (error) {
            console.error(`Failed to fetch ${activeTab}`, error);
        } finally {
            setLoading(false);
        }
    };

    const fetchChildEntities = async (type, parentId) => {
        setChildLoading(true);
        try {
            if (type === 'place') {
                const res = await axios.get(`/v1/entities/places?parent_id=${parentId}`);
                setChildEntities(res.data || []);
            } else if (type === 'event') {
                const res = await axios.get(`/v1/entities/events/${parentId}/needs`);
                setChildEntities(res.data || []);
            }
        } catch (error) {
            console.error(`Failed to fetch children for ${type}`, error);
        } finally {
            setChildLoading(false);
        }
    };

    const fetchItemTypes = async () => {
        try {
            const res = await axios.get('/v1/catalog/item-types');
            setItemTypes(res.data || []);
        } catch (error) {
            console.error('Failed to fetch item types', error);
        }
    };

    useEffect(() => {
        setIsModalOpen(false);
        setViewingEntity(null);
        setChildEntities([]);
        fetchData();
        if (activeTab === 'events') fetchItemTypes();
        // Always fetch companies if not on companies tab for dropdowns
        if (activeTab !== 'companies' && companies.length === 0) {
            axios.get('/v1/entities/companies').then(res => setCompanies(res.data || []));
        }
    }, [activeTab]);

    const handleOpenModal = (type, mode, item = null) => {
        setModalType(type);
        setModalMode(mode);
        setSelectedItem(item);
        if (mode === 'edit' && item) {
            setFormData({ ...item });
        } else {
            const defaults = {
                metadata: {},
                address: {},
                contact_points: [{ email: '', phone: '', type: 'general' }]
            };
            if (type === 'place' && viewingEntity?.type === 'place') {
                defaults.contained_in_place_id = viewingEntity.item.id;
                defaults.owner_id = viewingEntity.item.owner_id;
            }
            if (type === 'person') {
                if (activeTab === 'companies' && selectedItem) {
                    defaults.company_id = selectedItem.id;
                    defaults.role_name = 'Contact';
                }
            }
            if (type === 'event') {
                defaults.status = 'assumed';
                if (activeTab === 'companies' && selectedItem) defaults.company_id = selectedItem.id;
            }
            if (type === 'need' && viewingEntity) defaults.event_id = viewingEntity.item.id;
            if (type === 'place' && activeTab === 'companies' && selectedItem) defaults.owner_id = selectedItem.id;
            setFormData(defaults);
        }
        setIsModalOpen(true);
    };

    const getModalTitle = () => {
        const typeMap = {
            company: 'Corporate Entity',
            place: 'Operational Place',
            event: 'Project Timeline',
            need: 'Asset Requirement',
            person: 'Directory Individual',
            role: 'Organizational Role'
        };
        const title = typeMap[modalType] || modalType;
        return modalMode === 'create' ? `Assemble ${title}` : `Modify ${title}`;
    };

    const handleSave = async () => {
        try {
            let url = '';
            if (modalType === 'company') url = '/v1/entities/companies';
            if (modalType === 'place') url = '/v1/entities/places';
            if (modalType === 'event') url = '/v1/entities/events';
            if (modalType === 'person') url = '/v1/entities/people';
            if (modalType === 'role') url = '/v1/entities/roles';
            if (modalType === 'need') url = `/v1/entities/events/${viewingEntity.item.id}/needs`;

            if (modalMode === 'edit') {
                if (modalType === 'need') {
                    await axios.post(url, [formData]);
                } else {
                    await axios.put(`${url}/${selectedItem.id}`, formData);
                }
            } else {
                if (modalType === 'need') {
                    await axios.post(url, [formData]);
                } else if (modalType === 'person') {
                    // Create person, then create role if specified
                    const personRes = await axios.post(url, formData);
                    if (formData.company_id && formData.role_name) {
                        await axios.post('/v1/entities/roles', {
                            person_id: personRes.data.id,
                            organization_id: formData.company_id,
                            role_name: formData.role_name
                        });
                    }
                } else {
                    await axios.post(url, formData);
                }
            }
            setIsModalOpen(false);
            if (viewingEntity) {
                fetchChildEntities(viewingEntity.type, viewingEntity.item.id);
            } else {
                fetchData();
            }
        } catch (error) {
            console.error('Failed to save entity', error);
            alert(`Error: ${error.response?.data || error.message}`);
        }
    };

    const handleBack = () => {
        if (navigationStack.length > 0) {
            const previous = navigationStack[navigationStack.length - 1];
            const newStack = navigationStack.slice(0, -1);
            setNavigationStack(newStack);
            setViewingEntity(previous);
            fetchChildEntities(previous.type, previous.item.id);
        } else {
            setViewingEntity(null);
            setChildEntities([]);
        }
    };

    const filteredData = () => {
        const data = activeTab === 'companies' ? companies :
            activeTab === 'places' ? places.filter(p => !p.contained_in_place_id) :
                activeTab === 'events' ? events : people;
        if (!searchTerm) return data;
        return data.filter(item => {
            const name = item.name || `${item.given_name} ${item.family_name}`;
            const nameMatch = name.toLowerCase().includes(searchTerm.toLowerCase());
            const legalMatch = item.legal_name && item.legal_name.toLowerCase().includes(searchTerm.toLowerCase());
            return nameMatch || legalMatch;
        });
    };

    return (
        <div className="animate-fade-in" style={{ padding: '2rem', maxWidth: '1400px', margin: '0 auto' }}>
            <PageHeader
                title="Entity Management"
                subtitle="Centralized registry for corporate nodes, logistic hubs, and scheduled project timelines."
                actions={
                    !viewingEntity && (
                        <button
                            onClick={() => {
                                const typeMap = {
                                    companies: 'company',
                                    places: 'place',
                                    events: 'event',
                                    people: 'person'
                                };
                                handleOpenModal(typeMap[activeTab] || activeTab.slice(0, -1), 'create');
                            }}
                            className="btn-primary"
                        >
                            <Plus size={18} /> Add {
                                activeTab === 'companies' ? 'Company' :
                                    activeTab === 'places' ? 'Operational Place' :
                                        activeTab === 'events' ? 'Project/Event' : 'Individual'
                            }
                        </button>
                    )
                }
            />

            {!viewingEntity && (
                <div className="tab-nav">
                    {[
                        { id: 'companies', label: 'Companies', icon: <Building2 size={16} /> },
                        { id: 'people', label: 'Personnel', icon: <Users size={16} /> },
                        { id: 'places', label: 'Operational Places', icon: <MapPin size={16} /> },
                        { id: 'events', label: 'Events & Projects', icon: <Calendar size={16} /> }
                    ].map(tab => (
                        <button
                            key={tab.id}
                            className={`tab-button ${activeTab === tab.id ? 'active' : ''}`}
                            onClick={() => setActiveTab(tab.id)}
                        >
                            {tab.icon} {tab.label}
                        </button>
                    ))}
                </div>
            )}

            {!viewingEntity && (
                <div className="glass-card glass-surface" style={{ padding: '1rem', marginBottom: '2.5rem', borderRadius: 'var(--radius-lg)' }}>
                    <div style={{ position: 'relative' }}>
                        <Search size={18} style={{ position: 'absolute', left: '1rem', top: '50%', transform: 'translateY(-50%)', opacity: 0.4 }} />
                        <input
                            type="text"
                            placeholder={`Filter ${activeTab}...`}
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="form-input"
                            style={{ width: '100%', paddingLeft: '3rem', background: 'transparent', border: 'none' }}
                        />
                    </div>
                </div>
            )}

            {viewingEntity ? renderDetailView() : (
                <div className="entity-grid">
                    {loading ? (
                        <div className="flex-center flex-column" style={{ gridColumn: '1/-1', padding: '5rem 0' }}>
                            <div className="w-16 h-16 border-4 border-blue-600/10 border-t-blue-600 rounded-full animate-spin"></div>
                            <p style={{ marginTop: '1.5rem', color: 'var(--text-muted)' }}>Syncing Enterprise Directory...</p>
                        </div>
                    ) : filteredData().length === 0 ? (
                        <div className="flex-center flex-column" style={{ gridColumn: '1/-1', padding: '5rem 0', border: '2px dashed var(--border)', borderRadius: 'var(--radius-xl)' }}>
                            <p style={{ color: 'var(--text-muted)', fontSize: '1.25rem' }}>No records found.</p>
                        </div>
                    ) : renderCards()}
                </div>
            )}

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={getModalTitle()}
                actions={
                    <>
                        <button onClick={() => setIsModalOpen(false)} className="btn-secondary">Cancel</button>
                        <button onClick={handleSave} className="btn-primary">Save Changes</button>
                    </>
                }
            >
                {renderModalContent()}
            </Modal>
        </div>
    );

    function renderCards() {
        return filteredData().map(item => (
            <GlassCard key={item.id} className="glass-interactive entity-card animate-slide-up">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <div className="entity-icon-wrapper" style={{
                        background: activeTab === 'companies' ? 'rgba(99, 102, 241, 0.1)' :
                            activeTab === 'places' ? 'rgba(16, 185, 129, 0.1)' :
                                activeTab === 'events' ? 'rgba(168, 85, 247, 0.1)' : 'rgba(236, 72, 153, 0.1)',
                        color: activeTab === 'companies' ? 'var(--primary)' :
                            activeTab === 'places' ? 'var(--success)' :
                                activeTab === 'events' ? '#a855f7' : '#ec4899'
                    }}>
                        {activeTab === 'companies' ? <Building2 size={28} /> :
                            activeTab === 'places' ? <MapPin size={28} /> :
                                activeTab === 'events' ? <Calendar size={28} /> : <Users size={28} />}
                    </div>
                    <button
                        onClick={() => {
                            const typeMap = {
                                companies: 'company',
                                places: 'place',
                                events: 'event',
                                people: 'person'
                            };
                            handleOpenModal(typeMap[activeTab] || activeTab.slice(0, -1), 'edit', item);
                        }}
                        className="flex-center"
                        style={{ background: 'rgba(255,255,255,0.05)', width: '32px', height: '32px', borderRadius: '8px', color: 'var(--text-muted)' }}
                    >
                        <Edit2 size={14} />
                    </button>
                </div>

                <h3 className="entity-title">{activeTab === 'people' ? `${item.given_name} ${item.family_name}` : item.name}</h3>
                <p className="entity-subtitle">
                    {activeTab === 'companies' ? (item.legal_name || 'Generic Subsidiary') :
                        activeTab === 'places' ? `${item.category.toUpperCase()} // ${item.address?.address_locality || 'On-Site'}` :
                            activeTab === 'events' ? `${new Date(item.start_time).toLocaleDateString()} - ${item.status.toUpperCase()}` :
                                `${item.metadata?.title || 'Team Member'}`}
                </p>

                <div className="badge-group">
                    <span style={{ fontSize: '0.7rem', fontWeight: 800, opacity: 0.3 }}>UID_{String(item.id).padStart(4, '0')}</span>
                    <button
                        onClick={() => {
                            if (activeTab === 'places') {
                                setViewingEntity({ type: 'place', item });
                                fetchChildEntities('place', item.id);
                            } else if (activeTab === 'events') {
                                setViewingEntity({ type: 'event', item });
                                fetchChildEntities('event', item.id);
                            } else if (activeTab === 'companies') {
                                // Maybe show company places?
                                setActiveTab('places');
                                setSearchTerm(item.name);
                            } else if (activeTab === 'people') {
                                // Link to roles?
                            }
                        }}
                        style={{ background: 'transparent', color: 'var(--primary)', fontWeight: 700, fontSize: '0.875rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
                    >
                        {activeTab === 'companies' ? 'Explore Facilities' : activeTab === 'people' ? 'View Profile' : 'Manage Details'} <ChevronRight size={16} />
                    </button>
                </div>
            </GlassCard>
        ));
    }

    function renderDetailView() {
        const { type, item } = viewingEntity;
        return (
            <div className="animate-slide-up">
                <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', marginBottom: '3rem' }}>
                    <button onClick={handleBack} className="flex-center glass-surface" style={{ width: '48px', height: '48px', borderRadius: '12px', color: 'var(--text-muted)' }}>
                        <ArrowLeft size={24} />
                    </button>
                    <div style={{ flex: 1 }}>
                        <h2 style={{ fontSize: '2.5rem', fontWeight: 900, letterSpacing: '-0.04em' }}>{item.name}</h2>
                        <p style={{ color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', fontSize: '0.75rem', letterSpacing: '0.1em', marginTop: '0.25rem' }}>
                            {type === 'place' ? item.category : type} Configuration Node <span style={{ opacity: 0.3 }}>//</span> ID: {item.id}
                        </p>
                    </div>
                    <button
                        onClick={() => handleOpenModal(type === 'place' ? 'place' : 'need', 'create')}
                        className="btn-primary"
                        style={{ padding: '0.75rem 2rem', borderRadius: 'var(--radius-lg)' }}
                    >
                        <Plus size={18} /> Add {type === 'place' ? 'Sub-Place' : 'Asset Need'}
                    </button>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '3rem' }}>
                    <div className="glass-card glass-surface" style={{ padding: '2rem' }}>
                        <h4 className="form-label" style={{ marginBottom: '1.5rem' }}>Core Attributes</h4>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                            {type === 'place' ? (
                                <div className="flex-column gap-sm">
                                    <span className="form-label" style={{ fontSize: '0.65rem' }}>Registry Address</span>
                                    <p style={{ fontSize: '1.125rem', fontWeight: 600 }}>{item.address?.street_address || 'Internal Operations'}</p>
                                    <p style={{ color: 'var(--text-muted)' }}>{item.address?.address_locality}, {item.address?.address_country}</p>
                                    <div style={{ marginTop: '1rem', paddingTop: '1rem', borderTop: '1px solid var(--border)' }}>
                                        <span className="form-label" style={{ fontSize: '0.65rem' }}>Hierarchy Level</span>
                                        <p style={{ fontSize: '0.875rem', opacity: 0.6 }}>{navigationStack.map(n => n.item.name).join(' > ') || 'Root Level'}</p>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="flex-column gap-sm">
                                        <span className="form-label" style={{ fontSize: '0.65rem' }}>Project Window</span>
                                        <p style={{ fontSize: '1.125rem', fontWeight: 600 }}>{new Date(item.start_time).toLocaleString()}</p>
                                        <p style={{ color: 'var(--text-muted)' }}>Ends: {new Date(item.end_time).toLocaleString()}</p>
                                    </div>
                                    <div className="flex-column gap-sm">
                                        <span className="form-label" style={{ fontSize: '0.65rem' }}>Operations Mode</span>
                                        <StatusBadge status={item.status === 'confirmed' ? 'approved' : item.status === 'cancelled' ? 'cancelled' : 'pending'} />
                                    </div>
                                </>
                            )}
                        </div>
                    </div>

                    <div className="flex-column gap-md">
                        {childLoading ? <div className="flex-center" style={{ padding: '5rem' }}><div className="w-8 h-8 border-2 border-slate-800 border-t-primary rounded-full animate-spin"></div></div> :
                            childEntities.length === 0 ? <div className="flex-center glass-surface" style={{ padding: '5rem', borderRadius: 'var(--radius-xl)', border: '1px dashed var(--border)', color: 'var(--text-muted)' }}>Empty Manifest</div> :
                                childEntities.map(child => (
                                    <div key={child.id} className="glass-card glass-surface" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '1.25rem 2rem' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
                                            <div className="flex-center" style={{ width: '40px', height: '40px', background: 'rgba(255,255,255,0.03)', borderRadius: '10px', color: 'var(--text-muted)' }}>
                                                {type === 'place' ? <Layers size={18} /> : <Box size={18} />}
                                            </div>
                                            <div>
                                                <h5 style={{ fontWeight: 800, fontSize: '1.125rem' }}>
                                                    {type === 'place' ? child.name : (itemTypes.find(it => it.id === child.item_type_id)?.name || child.item_type_id)}
                                                </h5>
                                                <p style={{ color: 'var(--text-muted)', fontSize: '0.75rem', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                                                    {type === 'place' ? (child.category || 'unmapped area') : `${child.quantity} units requested`}
                                                </p>
                                            </div>
                                        </div>
                                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                                            {type === 'place' && (
                                                <button
                                                    onClick={() => {
                                                        setNavigationStack([...navigationStack, viewingEntity]);
                                                        setViewingEntity({ type: 'place', item: child });
                                                        fetchChildEntities('place', child.id);
                                                    }}
                                                    className="flex-center"
                                                    style={{ width: '32px', height: '32px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)', color: 'var(--primary)' }}
                                                >
                                                    <ChevronRight size={14} />
                                                </button>
                                            )}
                                            <button
                                                onClick={() => handleOpenModal(type === 'place' ? 'place' : 'need', 'edit', child)}
                                                className="flex-center"
                                                style={{ width: '32px', height: '32px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)', color: 'var(--text-muted)' }}
                                            >
                                                <Edit2 size={14} />
                                            </button>
                                        </div>
                                    </div>
                                ))}
                    </div>
                </div>
            </div>
        );
    }

    function renderModalContent() {
        if (modalType === 'company') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Registry Name</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} placeholder="Acme Corp" />
                </div>
                <div className="form-group">
                    <label className="form-label">Legal Designation</label>
                    <input className="form-input" value={formData.legal_name || ''} onChange={e => setFormData({ ...formData, legal_name: e.target.value })} placeholder="Acme Corporation Ltd." />
                </div>
                <div className="form-group">
                    <label className="form-label">Contextual Overview</label>
                    <textarea className="form-input" style={{ height: '100px', resize: 'none' }} value={formData.description || ''} onChange={e => setFormData({ ...formData, description: e.target.value })} />
                </div>
            </div>
        );
        if (modalType === 'person') return (
            <div className="flex-column">
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">Given Name</label>
                        <input className="form-input" value={formData.given_name || ''} onChange={e => setFormData({ ...formData, given_name: e.target.value })} />
                    </div>
                    <div className="form-group">
                        <label className="form-label">Family Name</label>
                        <input className="form-input" value={formData.family_name || ''} onChange={e => setFormData({ ...formData, family_name: e.target.value })} />
                    </div>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">Primary Email</label>
                        <input className="form-input" value={formData.contact_points?.[0]?.email || ''} onChange={e => {
                            const cp = [...(formData.contact_points || [{ email: '', phone: '', type: 'general' }])];
                            cp[0].email = e.target.value;
                            setFormData({ ...formData, contact_points: cp });
                        }} placeholder="email@example.com" />
                    </div>
                    <div className="form-group">
                        <label className="form-label">Phone Line</label>
                        <input className="form-input" value={formData.contact_points?.[0]?.phone || ''} onChange={e => {
                            const cp = [...(formData.contact_points || [{ email: '', phone: '', type: 'general' }])];
                            cp[0].phone = e.target.value;
                            setFormData({ ...formData, contact_points: cp });
                        }} placeholder="+1..." />
                    </div>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">Assigned Organization</label>
                        <select className="form-input form-select" value={formData.company_id || ''} onChange={e => setFormData({ ...formData, company_id: parseInt(e.target.value) })}>
                            <option value="">None / External</option>
                            {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                        </select>
                    </div>
                    <div className="form-group">
                        <label className="form-label">Professional Role</label>
                        <input className="form-input" value={formData.role_name || ''} onChange={e => setFormData({ ...formData, role_name: e.target.value })} placeholder="Technician, Manager..." />
                    </div>
                </div>
            </div>
        );
        if (modalType === 'place') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Owning Entity</label>
                    <select className="form-input form-select" value={formData.owner_id || ''} onChange={e => setFormData({ ...formData, owner_id: parseInt(e.target.value) })}>
                        <option value="">Select Company...</option>
                        {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                    </select>
                </div>
                <div className="form-group">
                    <label className="form-label">Place Designation</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} placeholder="e.g. North Wing, Rack 04, Warehouse B" />
                </div>
                <div className="form-group">
                    <label className="form-label">Place Category</label>
                    <select className="form-input form-select" value={formData.category || ''} onChange={e => setFormData({ ...formData, category: e.target.value })}>
                        <option value="">Select Category...</option>
                        <option value="site">Logistics Site (Facility)</option>
                        <option value="building">Managed Building</option>
                        <option value="floor">Strategic Floor</option>
                        <option value="room">Operations Room</option>
                        <option value="zone">Security Zone</option>
                        <option value="cabinet">Hardware Cabinet</option>
                    </select>
                </div>
                <div className="form-group" style={{ flexDirection: 'row', alignItems: 'center', gap: '0.5rem', display: 'flex' }}>
                    <input
                        type="checkbox"
                        id="is_internal"
                        checked={formData.is_internal || false}
                        onChange={e => setFormData({ ...formData, is_internal: e.target.checked })}
                    />
                    <label htmlFor="is_internal" className="form-label" style={{ marginBottom: 0, cursor: 'pointer' }}>Internal Logistics Hub (Warehouse/HQ)</label>
                </div>
                <div className="form-group">
                    <label className="form-label">Contextual Description</label>
                    <textarea className="form-input" style={{ height: '60px', resize: 'none' }} value={formData.description || ''} onChange={e => setFormData({ ...formData, description: e.target.value })} />
                </div>
                {(formData.category === 'site' || formData.category === 'building') && (
                    <>
                        <div className="form-group">
                            <label className="form-label">Street Address</label>
                            <input className="form-input" value={formData.address?.street_address || ''} onChange={e => setFormData({ ...formData, address: { ...formData.address, street_address: e.target.value } })} />
                        </div>
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                            <div className="form-group">
                                <label className="form-label">City Node</label>
                                <input className="form-input" value={formData.address?.address_locality || ''} onChange={e => setFormData({ ...formData, address: { ...formData.address, address_locality: e.target.value } })} />
                            </div>
                            <div className="form-group">
                                <label className="form-label">State / Region</label>
                                <input className="form-input" value={formData.address?.address_region || ''} onChange={e => setFormData({ ...formData, address: { ...formData.address, address_region: e.target.value } })} />
                            </div>
                        </div>
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                            <div className="form-group">
                                <label className="form-label">Postal Code</label>
                                <input className="form-input" value={formData.address?.postal_code || ''} onChange={e => setFormData({ ...formData, address: { ...formData.address, postal_code: e.target.value } })} />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Nation State</label>
                                <input className="form-input" value={formData.address?.address_country || ''} onChange={e => setFormData({ ...formData, address: { ...formData.address, address_country: e.target.value } })} />
                            </div>
                        </div>
                    </>
                )}
            </div>
        );
        if (modalType === 'event') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Project Designation</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} />
                </div>
                <div className="form-group">
                    <label className="form-label">Operational Status</label>
                    <select className="form-input form-select" value={formData.status || 'assumed'} onChange={e => setFormData({ ...formData, status: e.target.value })}>
                        <option value="assumed">Speculative (Assumed)</option>
                        <option value="confirmed">Operational (Confirmed)</option>
                        <option value="cancelled">Decommissioned (Cancelled)</option>
                    </select>
                </div>
                <div className="form-group">
                    <label className="form-label">Client / Sponsor</label>
                    <select className="form-input form-select" value={formData.company_id || ''} onChange={e => setFormData({ ...formData, company_id: parseInt(e.target.value) })}>
                        <option value="">Select Company...</option>
                        {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                    </select>
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">Commencement</label>
                        <input type="datetime-local" style={{ colorScheme: 'dark' }} className="form-input" value={formData.start_time ? new Date(formData.start_time).toISOString().slice(0, 16) : ''} onChange={e => setFormData({ ...formData, start_time: e.target.value })} />
                    </div>
                    <div className="form-group">
                        <label className="form-label">Conclusion</label>
                        <input type="datetime-local" style={{ colorScheme: 'dark' }} className="form-input" value={formData.end_time ? new Date(formData.end_time).toISOString().slice(0, 16) : ''} onChange={e => setFormData({ ...formData, end_time: e.target.value })} />
                    </div>
                </div>
                <div className="form-group">
                    <label className="form-label">Mission Description</label>
                    <textarea className="form-input" style={{ height: '60px', resize: 'none' }} value={formData.description || ''} onChange={e => setFormData({ ...formData, description: e.target.value })} />
                </div>
            </div>
        );

        if (modalType === 'need') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Requested Asset Class</label>
                    <select className="form-input form-select" value={formData.item_type_id || ''} onChange={e => setFormData({ ...formData, item_type_id: parseInt(e.target.value) })}>
                        <option value="">Select Catalog Item...</option>
                        {itemTypes.map(it => <option key={it.id} value={it.id}>{it.name} [{it.code}]</option>)}
                    </select>
                </div>
                <div className="form-group">
                    <label className="form-label">Target Location (Place)</label>
                    <select className="form-input form-select" value={formData.place_id || ''} onChange={e => setFormData({ ...formData, place_id: parseInt(e.target.value) })}>
                        <option value="">Unspecified location</option>
                        {places.map(p => <option key={p.id} value={p.id}>{p.name} ({p.category})</option>)}
                    </select>
                    <p style={{ fontSize: '0.65rem', color: 'var(--text-muted)', marginTop: '0.25rem' }}>Optional: Specify where this asset is needed within the site hierarchy.</p>
                </div>
                <div className="form-group">
                    <label className="form-label">Quantity Projection</label>
                    <input type="number" className="form-input" value={formData.quantity || 0} onChange={e => setFormData({ ...formData, quantity: parseInt(e.target.value) })} />
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginTop: '1rem' }}>
                    <input type="checkbox" checked={formData.is_assumed || false} onChange={e => setFormData({ ...formData, is_assumed: e.target.checked })} />
                    <span className="form-label" style={{ fontSize: '0.65rem' }}>Mark as theoretical requirement</span>
                </div>
            </div>
        );
        return null;
    }
};

export default EntityManager;
