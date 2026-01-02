CREATE TABLE service_pricing_tiers (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  service service NOT NULL,
  price_monthly numeric(12,2) NOT NULL,
  price_yearly  numeric(12,2) NOT NULL,
  popular boolean DEFAULT false,
  order_index integer DEFAULT 1,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_pricing_tiers_service ON service_pricing_tiers(service);
CREATE INDEX idx_pricing_tiers_active ON service_pricing_tiers(is_active);
CREATE INDEX idx_pricing_tiers_order ON service_pricing_tiers(order_index);
