CREATE TABLE hero_slide_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  hero_slide_id uuid NOT NULL REFERENCES hero_slides(id) ON DELETE CASCADE,

  language language NOT NULL,

  alt text NOT NULL,
  title text NOT NULL,
  cta_label text,

  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),

  UNIQUE(hero_slide_id, language)
);

CREATE INDEX idx_hero_slide_translations_slide ON hero_slide_translations(hero_slide_id);
CREATE INDEX idx_hero_slide_translations_language ON hero_slide_translations(language);
