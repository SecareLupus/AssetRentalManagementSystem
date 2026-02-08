CREATE TABLE IF NOT EXISTS system_settings (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed initial configuration
INSERT INTO system_settings (key, value) VALUES 
('company_identity', '{"name": "Secare Lupus Logistics", "logo_url": "", "support_email": "ops@secarelupus.com"}'),
('logistics_policies', '{"default_return_window_days": 14, "late_fee_per_day": 5.00, "currency": "USD"}'),
('feature_flags', '{"enable_auto_alerts": true, "enable_ai_forecasting": false}')
ON CONFLICT (key) DO NOTHING;
