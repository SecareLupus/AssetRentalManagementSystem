import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Building2, MapPin, Calendar, Plus, Users, Search, Edit2, ChevronRight } from 'lucide-react';
import { GlassCard, PageHeader, StatusBadge } from '../components/Shared';

const EntityManager = () => {
    const [activeTab, setActiveTab] = useState('companies');
    const [companies, setCompanies] = useState([]);
    const [sites, setSites] = useState([]);
    const [events, setEvents] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

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
            }
        } catch (error) {
            console.error(`Failed to fetch ${activeTab}`, error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [activeTab]);

    const renderTabs = () => (
        <div className="flex bg-slate-800/50 p-1 rounded-lg mb-6 w-fit border border-slate-700">
            {[
                { id: 'companies', label: 'Companies', icon: <Building2 size={16} /> },
                { id: 'sites', label: 'Sites & Locations', icon: <MapPin size={16} /> },
                { id: 'events', label: 'Events & Projects', icon: <Calendar size={16} /> }
            ].map(tab => (
                <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-md transition-all ${
                        activeTab === tab.id 
                        ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20' 
                        : 'text-slate-400 hover:text-white hover:bg-slate-700/50'
                    }`}
                >
                    {tab.icon}
                    {tab.label}
                </button>
            ))}
        </div>
    );

    return (
        <div className="p-8 max-w-7xl mx-auto">
            <PageHeader 
                title="Entity Management" 
                subtitle="Manage your organizational structure, physical locations, and project calendar."
                actions={
                    <button className="btn-primary flex items-center gap-2">
                        <Plus size={18} /> New {activeTab.slice(0, -1)}
                    </button>
                }
            />

            {renderTabs()}

            <GlassCard className="p-4 mb-6">
                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" size={18} />
                    <input 
                        type="text"
                        placeholder={`Search ${activeTab}...`}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full bg-slate-900/50 border border-slate-700 rounded-lg py-2 pl-10 pr-4 text-white focus:outline-none focus:border-blue-500 transition-colors"
                    />
                </div>
            </GlassCard>

            <div className="grid gap-4">
                {loading ? (
                    <div className="text-center py-20 text-slate-500">Loading entities...</div>
                ) : activeTab === 'companies' ? (
                    companies.length === 0 ? (
                        <div className="text-center py-20 text-slate-500 border border-dashed border-slate-800 rounded-xl">No companies configured.</div>
                    ) : (
                        companies.map(c => (
                            <GlassCard key={c.id} className="p-4 flex items-center justify-between hover:border-slate-600 transition-colors cursor-pointer group">
                                <div className="flex items-center gap-4">
                                    <div className="w-12 h-12 bg-blue-900/30 rounded-lg flex items-center justify-center text-blue-400">
                                        <Building2 size={24} />
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-lg">{c.name}</h3>
                                        <p className="text-slate-400 text-sm">{c.legal_name || 'No legal name'}</p>
                                    </div>
                                </div>
                                <div className="flex items-center gap-4 text-slate-500 group-hover:text-blue-400 transition-colors">
                                    <ChevronRight size={20} />
                                </div>
                            </GlassCard>
                        ))
                    )
                ) : activeTab === 'sites' ? (
                    sites.length === 0 ? (
                        <div className="text-center py-20 text-slate-500 border border-dashed border-slate-800 rounded-xl">No sites configured.</div>
                    ) : (
                        sites.map(s => (
                            <GlassCard key={s.id} className="p-4 flex items-center justify-between">
                                <div className="flex items-center gap-4">
                                    <div className="w-12 h-12 bg-emerald-900/30 rounded-lg flex items-center justify-center text-emerald-400">
                                        <MapPin size={24} />
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-lg">{s.name}</h3>
                                        <p className="text-slate-400 text-sm">{s.address_city}, {s.address_country}</p>
                                    </div>
                                </div>
                                <div className="flex items-center gap-2">
                                    <button className="p-2 hover:bg-slate-800 rounded-md text-slate-400 hover:text-white transition-colors">
                                        <Edit2 size={18} />
                                    </button>
                                </div>
                            </GlassCard>
                        ))
                    )
                ) : (
                    events.length === 0 ? (
                        <div className="text-center py-20 text-slate-500 border border-dashed border-slate-800 rounded-xl">No events scheduled.</div>
                    ) : (
                        events.map(e => (
                            <GlassCard key={e.id} className="p-4 flex items-center justify-between">
                                <div className="flex items-center gap-4">
                                    <div className="w-12 h-12 bg-purple-900/30 rounded-lg flex items-center justify-center text-purple-400">
                                        <Calendar size={24} />
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-lg">{e.name}</h3>
                                        <p className="text-slate-400 text-sm">
                                            {new Date(e.start_time).toLocaleDateString()} - {e.status.toUpperCase()}
                                        </p>
                                    </div>
                                </div>
                                <StatusBadge status={e.status === 'confirmed' ? 'active' : 'warn'} label={e.status} />
                            </GlassCard>
                        ))
                    )
                )}
            </div>
        </div>
    );
};

export default EntityManager;
