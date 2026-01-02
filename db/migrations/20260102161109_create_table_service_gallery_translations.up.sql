CREATE TABLE service_gallery_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  gallery_id uuid NOT NULL REFERENCES service_galleries(id) ON DELETE CASCADE,
  language language NOT NULL,
  title text,
  alt text,
  caption text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (gallery_id, language)
);

CREATE INDEX idx_service_gallery_translations_gallery ON service_gallery_translations(gallery_id);
CREATE INDEX idx_service_gallery_translations_language ON service_gallery_translations(language);
