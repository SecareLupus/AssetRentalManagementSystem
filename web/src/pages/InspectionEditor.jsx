import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import axios from 'axios';
import { ArrowLeft, Plus, Trash2, Save, ClipboardCheck, Info } from 'lucide-react';
import { GlassCard, PageHeader } from '../components/Shared';

const InspectionEditor = () => {
    const { id } = useParams(); // undefined if creating new
    const navigate = useNavigate();
    const [template, setTemplate] = useState({
        name: '',
        description: '',
        fields: []
    });
    const [itemTypes, setItemTypes] = useState([]);
    const [assignedTypes, setAssignedTypes] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const typesRes = await axios.get('/v1/catalog/item-types');
                setItemTypes(typesRes.data || []);

                if (id) {
                    const detailRes = await axios.get(`/v1/catalog/inspection-templates/${id}`);
                    setTemplate(detailRes.data);
                }
            } catch (err) {
                console.error("Fetch failed", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    const addField = () => {
        setTemplate({
            ...template,
            fields: [...template.fields, { label: '', field_type: 'text', required: false, display_order: template.fields.length }]
        });
    };

    const removeField = (index) => {
        const fields = template.fields.filter((_, i) => i !== index);
        setTemplate({ ...template, fields });
    };

    const handleSave = async () => {
        try {
            if (id) {
                await axios.put(`/v1/catalog/inspection-templates/${id}`, template);
            } else {
                const res = await axios.post('/v1/catalog/inspection-templates', template);
                navigate(`/admin/inspections/${res.data.id}`);
            }
            alert("Template saved!");
        } catch (err) {
            alert("Save failed");
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <Link to="/admin" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-muted)', textDecoration: 'none', marginBottom: '2rem' }}>
                <ArrowLeft size={16} /> Back to Admin
            </Link>

            <PageHeader
                title={id ? "Edit Template" : "New Inspection Template"}
                actions={<button onClick={handleSave} className="btn-primary"><Save size={18} /> Save Template</button>}
            />

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 300px', gap: '2rem' }}>
                <div>
                    <GlassCard style={{ marginBottom: '2rem' }}>
                        <div style={{ marginBottom: '1.5rem' }}>
                            <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Template Name</label>
                            <input
                                className="glass"
                                style={{ width: '100%', padding: '1rem', borderRadius: '0.75rem', color: 'white' }}
                                value={template.name}
                                onChange={e => setTemplate({ ...template, name: e.target.value })}
                            />
                        </div>
                        <div>
                            <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.5rem' }}>Description</label>
                            <textarea
                                className="glass"
                                rows="3"
                                style={{ width: '100%', padding: '1rem', borderRadius: '0.75rem', color: 'white', resize: 'none' }}
                                value={template.description}
                                onChange={e => setTemplate({ ...template, description: e.target.value })}
                            />
                        </div>
                    </GlassCard>

                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
                        <h3 style={{ fontWeight: 700, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <ClipboardCheck size={20} color="var(--primary)" /> Form Fields
                        </h3>
                        <button onClick={addField} className="glass" style={{ padding: '0.4rem 0.75rem', fontSize: '0.875rem' }}>
                            <Plus size={16} /> Add Field
                        </button>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                        {template.fields.map((field, i) => (
                            <GlassCard key={i} style={{ padding: '1rem', display: 'grid', gridTemplateColumns: '2fr 1fr 100px 40px', gap: '1rem', alignItems: 'end' }}>
                                <div>
                                    <label style={{ fontSize: '0.65rem', color: 'var(--text-muted)' }}>Label</label>
                                    <input
                                        className="glass"
                                        style={{ width: '100%', padding: '0.5rem', fontSize: '0.875rem', color: 'white' }}
                                        value={field.label}
                                        onChange={e => {
                                            const fields = [...template.fields];
                                            fields[i].label = e.target.value;
                                            setTemplate({ ...template, fields });
                                        }}
                                    />
                                </div>
                                <div>
                                    <label style={{ fontSize: '0.65rem', color: 'var(--text-muted)' }}>Type</label>
                                    <select
                                        className="glass"
                                        style={{ width: '100%', padding: '0.5rem', fontSize: '0.875rem', color: 'white', background: 'var(--surface)' }}
                                        value={field.field_type}
                                        onChange={e => {
                                            const fields = [...template.fields];
                                            fields[i].field_type = e.target.value;
                                            setTemplate({ ...template, fields });
                                        }}
                                    >
                                        <option value="text">Text / Comment</option>
                                        <option value="boolean">Pass / Fail</option>
                                        <option value="number">Numeric Measure</option>
                                    </select>
                                </div>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', height: '40px' }}>
                                    <input
                                        type="checkbox"
                                        checked={field.required}
                                        onChange={e => {
                                            const fields = [...template.fields];
                                            fields[i].required = e.target.checked;
                                            setTemplate({ ...template, fields });
                                        }}
                                    />
                                    <span style={{ fontSize: '0.75rem' }}>Required</span>
                                </div>
                                <button onClick={() => removeField(i)} style={{ background: 'transparent', color: 'var(--text-muted)' }}>
                                    <Trash2 size={18} />
                                </button>
                            </GlassCard>
                        ))}
                    </div>
                </div>

                <aside>
                    <GlassCard>
                        <h4 style={{ fontWeight: 700, marginBottom: '1rem' }}>Template Usage</h4>
                        <div style={{ display: 'flex', alignItems: 'flex-start', gap: '0.5rem', marginBottom: '1.5rem' }}>
                            <Info size={14} color="var(--primary)" style={{ marginTop: '0.2rem' }} />
                            <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Assign this template to Item Types to enforce it during check-in/out.</p>
                        </div>
                        {/* Assignment UI - TBD */}
                        <p style={{ fontSize: '0.75rem', opacity: 0.5 }}>Manage assignments via Item Type settings.</p>
                    </GlassCard>
                </aside>
            </div>
        </div>
    );
};

export default InspectionEditor;
