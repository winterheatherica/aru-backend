CREATE TABLE partners (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_path text NOT NULL,
  order_index integer DEFAULT 0,
  is_active_partner_grid boolean DEFAULT true,
  is_active_partner_scroller boolean DEFAULT false,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);


CREATE INDEX idx_partners_order_index ON partners(order_index);
CREATE INDEX idx_partners_active_grid ON partners(is_active_partner_grid);
CREATE INDEX idx_partners_active_scroller ON partners(is_active_partner_scroller);