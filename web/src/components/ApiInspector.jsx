import React, { useState } from 'react';
import { useDeveloper } from '../context/DeveloperContext';
import { Terminal, X, ChevronRight, ChevronDown, Activity, Clock, Database, Globe } from 'lucide-react';

const ApiInspector = () => {
    const { isDevMode, apiLogs, clearLogs } = useDeveloper();
    const [isOpen, setIsOpen] = useState(false);
    const [expandedLog, setExpandedLog] = useState(null);

    if (!isDevMode) return null;

    return (
        <>
            {/* Trigger Button */}
            <button
                onClick={() => setIsOpen(true)}
                className="btn-primary"
                style={{
                    position: 'fixed',
                    bottom: '1rem',
                    right: '1rem',
                    zIndex: 1000,
                    display: 'flex',
                    alignItems: 'center',
                    gap: '0.5rem',
                    borderRadius: '9999px',
                    boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.4)'
                }}
            >
                <Terminal size={20} />
                <span style={{ fontSize: '0.875rem' }}>API Inspector</span>
                {apiLogs.length > 0 && (
                    <span style={{
                        backgroundColor: 'var(--error)',
                        borderRadius: '9999px',
                        padding: '0 0.4rem',
                        fontSize: '0.75rem',
                        color: 'white'
                    }}>
                        {apiLogs.length}
                    </span>
                )}
            </button>

            {/* Slide-over Panel */}
            {isOpen && (
                <div
                    className="glass"
                    style={{
                        position: 'fixed',
                        top: 0,
                        right: 0,
                        bottom: 0,
                        width: '100%',
                        maxWidth: '600px',
                        zIndex: 1001,
                        display: 'flex',
                        flexDirection: 'column',
                        boxShadow: '-10px 0 30px rgba(0,0,0,0.5)',
                        borderLeft: '1px solid var(--border)'
                    }}
                >
                    {/* Header */}
                    <div style={{ padding: '1rem', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                            <Terminal size={20} color="var(--primary)" />
                            <h2 style={{ fontSize: '1.125rem', fontWeight: 600 }}>API Inspector</h2>
                        </div>
                        <div style={{ display: 'flex', gap: '0.75rem', alignItems: 'center' }}>
                            <button
                                onClick={clearLogs}
                                style={{ fontSize: '0.75rem', color: 'var(--text-muted)', background: 'transparent' }}
                            >
                                Clear
                            </button>
                            <button
                                onClick={() => setIsOpen(false)}
                                style={{ background: 'transparent', color: 'var(--text)', display: 'flex' }}
                            >
                                <X size={20} />
                            </button>
                        </div>
                    </div>

                    {/* Logs List */}
                    <div style={{ flex: 1, overflowY: 'auto', padding: '0.5rem' }}>
                        {apiLogs.length === 0 ? (
                            <div style={{ padding: '3rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                                <Activity size={48} style={{ marginBottom: '1rem', opacity: 0.2, margin: '0 auto' }} />
                                <p>No API calls captured yet.</p>
                                <p style={{ fontSize: '0.75rem' }}>Interaction with the UI will trigger logs.</p>
                            </div>
                        ) : (
                            apiLogs.map((log) => (
                                <div
                                    key={log.id}
                                    style={{
                                        marginBottom: '0.5rem',
                                        borderRadius: '0.5rem',
                                        overflow: 'hidden',
                                        background: expandedLog === log.id ? 'var(--surface-hover)' : 'var(--surface)',
                                        border: '1px solid var(--border)'
                                    }}
                                >
                                    <div
                                        onClick={() => setExpandedLog(expandedLog === log.id ? null : log.id)}
                                        style={{
                                            padding: '0.75rem',
                                            display: 'flex',
                                            alignItems: 'center',
                                            gap: '0.75rem',
                                            cursor: 'pointer'
                                        }}
                                    >
                                        <span style={{
                                            fontSize: '0.75rem',
                                            fontWeight: 700,
                                            padding: '0.125rem 0.375rem',
                                            borderRadius: '0.25rem',
                                            background: log.status < 300 ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)',
                                            color: log.status < 300 ? 'var(--success)' : 'var(--error)'
                                        }}>
                                            {log.method}
                                        </span>
                                        <span style={{
                                            fontSize: '0.875rem',
                                            fontWeight: 500,
                                            flex: 1,
                                            whiteSpace: 'nowrap',
                                            overflow: 'hidden',
                                            textOverflow: 'ellipsis',
                                            color: 'var(--text)'
                                        }}>
                                            {log.url}
                                        </span>
                                        <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                                            <Clock size={12} /> {log.duration}ms
                                        </span>
                                        {expandedLog === log.id ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                                    </div>

                                    {expandedLog === log.id && (
                                        <div style={{ padding: '0 1rem 1rem 1rem', fontSize: '0.75rem', borderTop: '1px solid var(--border)' }}>
                                            <div style={{ marginTop: '1rem' }}>
                                                <div style={{ color: 'var(--primary)', fontWeight: 600, marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                                    <Globe size={14} /> Request
                                                </div>
                                                <pre style={{ background: '#0a0a0a', padding: '0.75rem', borderRadius: '0.375rem', overflowX: 'auto', color: '#10b981' }}>
                                                    {JSON.stringify(log.requestData || 'No data', null, 2)}
                                                </pre>
                                            </div>

                                            <div style={{ marginTop: '1rem' }}>
                                                <div style={{ color: 'var(--success)', fontWeight: 600, marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                                    <Database size={14} /> Response ({log.status})
                                                </div>
                                                <pre style={{ background: '#0a0a0a', padding: '0.75rem', borderRadius: '0.375rem', overflowX: 'auto', color: '#6366f1' }}>
                                                    {JSON.stringify(log.responseData, null, 2)}
                                                </pre>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            ))
                        )}
                    </div>
                </div>
            )}
        </>
    );
};

export default ApiInspector;
