CREATE TABLE service_galleries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  service service NOT NULL,
  media_type gallery_media NOT NULL,
  src text NOT NULL,
  thumbnail text,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_service_galleries_service ON service_galleries(service);
CREATE INDEX idx_service_galleries_active ON service_galleries(is_active);
CREATE INDEX idx_service_galleries_media_type ON service_galleries(media_type);
