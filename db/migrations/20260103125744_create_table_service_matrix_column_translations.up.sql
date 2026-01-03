CREATE TABLE service_matrix_column_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  column_id uuid NOT NULL REFERENCES service_matrix_columns(id) ON DELETE CASCADE,
  language language NOT NULL,
  label text NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (column_id, language)
);

CREATE INDEX idx_service_matrix_column_translations_column ON service_matrix_column_translations(column_id);
CREATE INDEX idx_service_matrix_column_translations_language ON service_matrix_column_translations(language);
