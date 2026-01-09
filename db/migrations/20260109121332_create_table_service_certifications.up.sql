CREATE TABLE service_certifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  service service NOT NULL,
  order_index integer DEFAULT 1,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_service_certifications_service ON service_certifications(service);
CREATE INDEX idx_service_certifications_active ON service_certifications(is_active);
CREATE INDEX idx_service_certifications_order ON service_certifications(order_index);
