CREATE TABLE service_certification_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  certification_id uuid NOT NULL REFERENCES service_certifications(id) ON DELETE CASCADE,
  language language NOT NULL,
  title text,
  alt text,
  caption text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (certification_id, language)
);

CREATE INDEX idx_service_certification_translations_certification ON service_certification_translations(certification_id);
CREATE INDEX idx_service_certification_translations_language ON service_certification_translations(language);
