import React from 'react';
import { Check, X, Clock, AlertCircle, Info } from 'lucide-react';

/**
 * Standard glassmorphism card wrapper.
 */
export const GlassCard = ({ children, className = '', style = {}, ...props }) => (
  <div className={`glass-card ${className}`} style={style} {...props}>
    {children}
  </div>
);

/**
 * Standard page header with title, subtitle, and optional action buttons.
 */
export const PageHeader = ({ title, subtitle, actions }) => (
  <header className="page-header">
    <div>
      <h1>{title}</h1>
      {subtitle && <p>{subtitle}</p>}
    </div>
    {actions && <div style={{ display: 'flex', gap: '0.75rem' }}>{actions}</div>}
  </header>
);

/**
 * Unified Status Badge component.
 * Supports different semantic statuses: 'draft', 'pending', 'approved', 'rejected', 'cancelled', 'available', 'reserved', 'maintenance', 'deployed'.
 */
export const StatusBadge = ({ status }) => {
  // Map status to visual styles
  // We combine the logic from ReservationsList and ItemTypeDetails
  const config = {
    // Reservation statuses
    draft: { bg: 'rgba(148, 163, 184, 0.2)', text: 'var(--text-muted)', icon: Clock },
    pending: { bg: 'rgba(245, 158, 11, 0.2)', text: 'var(--warning)', icon: AlertCircle },
    approved: { bg: 'rgba(16, 185, 129, 0.2)', text: 'var(--success)', icon: Check },
    rejected: { bg: 'rgba(239, 104, 104, 0.2)', text: 'var(--error)', icon: X },
    cancelled: { bg: 'rgba(148, 163, 184, 0.1)', text: 'var(--text-muted)', icon: X },

    // Asset statuses
    available: { bg: 'rgba(16, 185, 129, 0.2)', text: 'var(--success)', icon: Check },
    reserved: { bg: 'rgba(245, 158, 11, 0.2)', text: 'var(--warning)', icon: Clock },
    maintenance: { bg: 'rgba(239, 68, 68, 0.2)', text: 'var(--error)', icon: AlertCircle },
    deployed: { bg: 'rgba(99, 102, 241, 0.2)', text: 'var(--primary)', icon: Info },

    // Fallback
    default: { bg: 'var(--surface)', text: 'var(--text)', icon: null }
  };

  const style = config[status?.toLowerCase()] || config.default;
  const Icon = style.icon;

  return (
    <div style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: '0.35rem',
      padding: '0.25rem 0.6rem',
      borderRadius: '2rem',
      background: style.bg,
      color: style.text,
      fontSize: '0.75rem',
      fontWeight: 700,
      textTransform: 'uppercase',
      whiteSpace: 'nowrap'
    }}>
      {Icon && <Icon size={12} />}
      <span>{status}</span>
    </div>
  );
};

/**
 * Standard Modal component.
 */
export const Modal = ({ isOpen, onClose, title, children, actions, size = 'md' }) => {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className={`modal-container modal-${size}`}>
        <div className="modal-header">
          <h2 className="modal-title">{title}</h2>
          <button onClick={onClose} className="modal-close">
            <X size={20} />
          </button>
        </div>
        <div className="modal-body">
          {children}
        </div>
        {actions && (
          <div className="modal-footer">
            {actions}
          </div>
        )}
      </div>
    </div>
  );
};

