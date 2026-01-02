CREATE TABLE service_pricing_tier_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tier_id uuid NOT NULL REFERENCES service_pricing_tiers(id) ON DELETE CASCADE,
  language language NOT NULL,
  name text NOT NULL,
  description text,
  features text[] NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),

  UNIQUE (tier_id, language)
);

CREATE INDEX idx_pricing_tier_translations_tier ON service_pricing_tier_translations(tier_id);
CREATE INDEX idx_pricing_tier_translations_language ON service_pricing_tier_translations(language);
