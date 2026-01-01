CREATE TABLE client_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id uuid NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
  language language NOT NULL,
  title text,
  alt text,
  description text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (client_id, language)
);

CREATE INDEX idx_client_translations_language ON client_translations(language);
CREATE INDEX idx_client_translations_client ON client_translations(client_id);
