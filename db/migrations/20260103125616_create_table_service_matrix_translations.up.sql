CREATE TABLE service_matrix_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  matrix_id uuid NOT NULL REFERENCES service_matrices(id) ON DELETE CASCADE,
  language language NOT NULL,
  title text,
  description text,
  footnote text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (matrix_id, language)
);

CREATE INDEX idx_service_matrix_translations_matrix ON service_matrix_translations(matrix_id);
CREATE INDEX idx_service_matrix_translations_language ON service_matrix_translations(language);
