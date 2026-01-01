CREATE TABLE award_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  award_id uuid NOT NULL REFERENCES awards(id) ON DELETE CASCADE,
  language language NOT NULL,
  alt text,
  title text,
  label text,
  description text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (award_id, language)
);

CREATE INDEX idx_award_translations_language ON award_translations(language);
CREATE INDEX idx_award_translations_award ON award_translations(award_id);
