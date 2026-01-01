CREATE TABLE partner_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  partner_id uuid NOT NULL REFERENCES partners(id) ON DELETE CASCADE,
  language language_type NOT NULL,
  title text,
  alt text,
  description text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (partner_id, language)
);

CREATE INDEX idx_partner_translations_language ON partner_translations(language);
CREATE INDEX idx_partner_translations_partner ON partner_translations(partner_id);
