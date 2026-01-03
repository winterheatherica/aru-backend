CREATE TABLE service_matrices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  service service NOT NULL,
  compact boolean DEFAULT false,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_service_matrices_service ON service_matrices(service);
CREATE INDEX idx_service_matrices_active ON service_matrices(is_active);
CREATE INDEX idx_service_matrices_compact ON service_matrices(compact);
