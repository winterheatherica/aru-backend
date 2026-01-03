CREATE TABLE service_matrix_row_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  row_id uuid NOT NULL REFERENCES service_matrix_rows(id) ON DELETE CASCADE,
  language language NOT NULL,
  feature text NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (row_id, language)
);

CREATE INDEX idx_service_matrix_row_translations_row ON service_matrix_row_translations(row_id);
CREATE INDEX idx_service_matrix_row_translations_language ON service_matrix_row_translations(language);
