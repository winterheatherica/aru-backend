CREATE TABLE promo_slide_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_slide_id uuid NOT NULL REFERENCES promo_slides(id) ON DELETE CASCADE,
  language language_type NOT NULL,
  title text,
  alt text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (promo_slide_id, language)
);

CREATE INDEX idx_promo_slide_translations_language ON promo_slide_translations(language);
CREATE INDEX idx_promo_slide_translations_slide ON promo_slide_translations(promo_slide_id);
