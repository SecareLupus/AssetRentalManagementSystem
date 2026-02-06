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
    const [sites, setSites] = useState([]);
    const [events, setEvents] = useState([]);
    const [contacts, setContacts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    // Detailed View State
    const [viewingEntity, setViewingEntity] = useState(null); // { type: 'site'|'event', item: Object }
    const [childEntities, setChildEntities] = useState([]); // Locations for site, Needs for event
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
            } else if (activeTab === 'sites') {
                const res = await axios.get('/v1/entities/sites');
                setSites(res.data || []);
            } else if (activeTab === 'events') {
                const res = await axios.get('/v1/entities/events');
                setEvents(res.data || []);
            } else if (activeTab === 'contacts') {
                const res = await axios.get('/v1/entities/contacts');
                setContacts(res.data || []);
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
            if (type === 'site') {
                const res = await axios.get(`/v1/entities/locations?site_id=${parentId}`);
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
            const defaults = {};
            if (type === 'location' && viewingEntity) defaults.site_id = viewingEntity.item.id;
            if (type === 'need' && viewingEntity) defaults.event_id = viewingEntity.item.id;
            if (type === 'site' && activeTab === 'companies' && selectedItem) defaults.company_id = selectedItem.id;
            if (type === 'contact' && activeTab === 'companies' && selectedItem) defaults.company_id = selectedItem.id;
            setFormData(defaults);
        }
        setIsModalOpen(true);
    };

    const getModalTitle = () => {
        const typeMap = {
            company: 'Corporate Entity',
            site: 'Logistics Site',
            event: 'Project Timeline',
            location: 'Site Location',
            need: 'Asset Requirement',
            contact: 'Directory Contact'
        };
        const title = typeMap[modalType] || modalType;
        return modalMode === 'create' ? `Assemble ${title}` : `Modify ${title}`;
    };

    const handleSave = async () => {
        try {
            let url = '';
            if (modalType === 'company') url = '/v1/entities/companies';
            if (modalType === 'site') url = '/v1/entities/sites';
            if (modalType === 'event') url = '/v1/entities/events';
            if (modalType === 'location') url = '/v1/entities/locations';
            if (modalType === 'contact') url = '/v1/entities/contacts';
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
        }
    };

    const handleBack = () => {
        setViewingEntity(null);
        setChildEntities([]);
    };

    const filteredData = () => {
        const data = activeTab === 'companies' ? companies :
            activeTab === 'sites' ? sites :
                activeTab === 'events' ? events : contacts;
        if (!searchTerm) return data;
        return data.filter(item => {
            const nameMatch = (item.name || `${item.first_name} ${item.last_name}`).toLowerCase().includes(searchTerm.toLowerCase());
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
                                let type = activeTab.slice(0, -1);
                                if (activeTab === 'companies') type = 'company';
                                handleOpenModal(type, 'create');
                            }}
                            className="btn-primary"
                        >
                            <Plus size={18} /> Add {activeTab === 'companies' ? 'Company' : activeTab === 'sites' ? 'Site' : activeTab === 'events' ? 'Event' : 'Contact'}
                        </button>
                    )
                }
            />

            {!viewingEntity && (
                <div className="tab-nav">
                    {[
                        { id: 'companies', label: 'Companies', icon: <Building2 size={16} /> },
                        { id: 'contacts', label: 'Contacts', icon: <Users size={16} /> },
                        { id: 'sites', label: 'Sites & Locations', icon: <MapPin size={16} /> },
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
                            activeTab === 'sites' ? 'rgba(16, 185, 129, 0.1)' :
                                activeTab === 'events' ? 'rgba(168, 85, 247, 0.1)' : 'rgba(236, 72, 153, 0.1)',
                        color: activeTab === 'companies' ? 'var(--primary)' :
                            activeTab === 'sites' ? 'var(--success)' :
                                activeTab === 'events' ? '#a855f7' : '#ec4899'
                    }}>
                        {activeTab === 'companies' ? <Building2 size={28} /> :
                            activeTab === 'sites' ? <MapPin size={28} /> :
                                activeTab === 'events' ? <Calendar size={28} /> : <Users size={28} />}
                    </div>
                    <button
                        onClick={() => {
                            let type = activeTab.slice(0, -1);
                            if (activeTab === 'companies') type = 'company';
                            handleOpenModal(type, 'edit', item);
                        }}
                        className="flex-center"
                        style={{ background: 'rgba(255,255,255,0.05)', width: '32px', height: '32px', borderRadius: '8px', color: 'var(--text-muted)' }}
                    >
                        <Edit2 size={14} />
                    </button>
                </div>

                <h3 className="entity-title">{activeTab === 'contacts' ? `${item.first_name} ${item.last_name}` : item.name}</h3>
                <p className="entity-subtitle">
                    {activeTab === 'companies' ? (item.legal_name || 'Generic Subsidiary') :
                        activeTab === 'sites' ? `${item.address_city}, ${item.address_country}` :
                            activeTab === 'events' ? `${new Date(item.start_time).toLocaleDateString()} - ${item.status.toUpperCase()}` :
                                `${item.role || 'Unspecified Role'} @ ${companies.find(c => c.id === item.company_id)?.name || 'External Organization'}`}
                </p>

                <div className="badge-group">
                    <span style={{ fontSize: '0.7rem', fontWeight: 800, opacity: 0.3 }}>UID_{String(item.id).padStart(4, '0')}</span>
                    <button
                        onClick={() => {
                            if (activeTab === 'sites') {
                                setViewingEntity({ type: 'site', item });
                                fetchChildEntities('site', item.id);
                            } else if (activeTab === 'events') {
                                setViewingEntity({ type: 'event', item });
                                fetchChildEntities('event', item.id);
                            } else if (activeTab === 'companies') {
                                setActiveTab('contacts');
                                setSearchTerm(item.name);
                            } else {
                                // For contacts, maybe link to company?
                                if (item.company_id) {
                                    setActiveTab('companies');
                                    setSearchTerm(companies.find(c => c.id === item.company_id)?.name || '');
                                }
                            }
                        }}
                        style={{ background: 'transparent', color: 'var(--primary)', fontWeight: 700, fontSize: '0.875rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
                    >
                        {activeTab === 'companies' ? 'Explore Contacts' : activeTab === 'contacts' ? 'View Company' : 'Manage Details'} <ChevronRight size={16} />
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
                            {type} Configuration Node <span style={{ opacity: 0.3 }}>//</span> ID: {item.id}
                        </p>
                    </div>
                    <button
                        onClick={() => handleOpenModal(type === 'site' ? 'location' : 'need', 'create')}
                        className="btn-primary"
                        style={{ padding: '0.75rem 2rem', borderRadius: 'var(--radius-lg)' }}
                    >
                        <Plus size={18} /> Add {type === 'site' ? 'Site Location' : 'Asset Need'}
                    </button>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '3rem' }}>
                    <div className="glass-card glass-surface" style={{ padding: '2rem' }}>
                        <h4 className="form-label" style={{ marginBottom: '1.5rem' }}>Core Attributes</h4>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                            {type === 'site' ? (
                                <div className="flex-column gap-sm">
                                    <span className="form-label" style={{ fontSize: '0.65rem' }}>Registry Address</span>
                                    <p style={{ fontSize: '1.125rem', fontWeight: 600 }}>{item.address_street || 'No Street Record'}</p>
                                    <p style={{ color: 'var(--text-muted)' }}>{item.address_city}, {item.address_country}</p>
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
                                                {type === 'site' ? <Layers size={18} /> : <Box size={18} />}
                                            </div>
                                            <div>
                                                <h5 style={{ fontWeight: 800, fontSize: '1.125rem' }}>
                                                    {type === 'site' ? child.name : (itemTypes.find(it => it.id === child.item_type_id)?.name || child.item_type_id)}
                                                </h5>
                                                <p style={{ color: 'var(--text-muted)', fontSize: '0.75rem', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                                                    {type === 'site' ? (child.location_type || 'unmapped area') : `${child.quantity} units requested`}
                                                </p>
                                            </div>
                                        </div>
                                        <button
                                            onClick={() => handleOpenModal(type === 'site' ? 'location' : 'need', 'edit', child)}
                                            className="flex-center"
                                            style={{ width: '32px', height: '32px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)', color: 'var(--text-muted)' }}
                                        >
                                            <Edit2 size={14} />
                                        </button>
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
        if (modalType === 'contact') return (
            <div className="flex-column">
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">First Name</label>
                        <input className="form-input" value={formData.first_name || ''} onChange={e => setFormData({ ...formData, first_name: e.target.value })} />
                    </div>
                    <div className="form-group">
                        <label className="form-label">Last Name</label>
                        <input className="form-input" value={formData.last_name || ''} onChange={e => setFormData({ ...formData, last_name: e.target.value })} />
                    </div>
                </div>
                <div className="form-group">
                    <label className="form-label">Affiliated Company</label>
                    <select className="form-input form-select" value={formData.company_id || ''} onChange={e => setFormData({ ...formData, company_id: parseInt(e.target.value) })}>
                        <option value="">Select Company...</option>
                        {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                    </select>
                </div>
                <div className="form-group">
                    <label className="form-label">Professional Role</label>
                    <input className="form-input" value={formData.role || ''} onChange={e => setFormData({ ...formData, role: e.target.value })} placeholder="Technical Lead" />
                </div>
                <div className="form-group">
                    <label className="form-label">Electronic Mail</label>
                    <input className="form-input" value={formData.email || ''} onChange={e => setFormData({ ...formData, email: e.target.value })} placeholder="user@example.com" />
                </div>
                <div className="form-group">
                    <label className="form-label">Communication Line</label>
                    <input className="form-input" value={formData.phone || ''} onChange={e => setFormData({ ...formData, phone: e.target.value })} placeholder="+1 (555) 000-0000" />
                </div>
            </div>
        );
        if (modalType === 'site') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Owning Entity</label>
                    <select className="form-input form-select" value={formData.company_id || ''} onChange={e => setFormData({ ...formData, company_id: parseInt(e.target.value) })}>
                        <option value="">Select Company...</option>
                        {companies.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                    </select>
                </div>
                <div className="form-group">
                    <label className="form-label">Operational Site Name</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div className="form-group">
                        <label className="form-label">City Node</label>
                        <input className="form-input" value={formData.address_city || ''} onChange={e => setFormData({ ...formData, address_city: e.target.value })} />
                    </div>
                    <div className="form-group">
                        <label className="form-label">Nation State</label>
                        <input className="form-input" value={formData.address_country || ''} onChange={e => setFormData({ ...formData, address_country: e.target.value })} />
                    </div>
                </div>
            </div>
        );
        if (modalType === 'event') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Project Designation</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} />
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
            </div>
        );
        if (modalType === 'location') return (
            <div className="flex-column">
                <div className="form-group">
                    <label className="form-label">Logistics Node Name</label>
                    <input className="form-input" value={formData.name || ''} onChange={e => setFormData({ ...formData, name: e.target.value })} placeholder="VIP Suite A" />
                </div>
                <div className="form-group">
                    <label className="form-label">Categorical Type</label>
                    <select className="form-input form-select" value={formData.location_type || ''} onChange={e => setFormData({ ...formData, location_type: e.target.value })}>
                        <option value="">Select Type...</option>
                        <option value="room">Room / Suite</option>
                        <option value="zone">Operations Zone</option>
                        <option value="floor">Strategic Floor</option>
                    </select>
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
