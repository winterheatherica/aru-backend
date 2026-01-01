CREATE TABLE clients (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_path text NOT NULL,
  order_index integer DEFAULT 0,
  is_active_client_scroller boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_clients_order_index ON clients(order_index);
CREATE INDEX idx_clients_active_scroller ON clients(is_active_client_scroller);