CREATE TABLE service_matrix_cell_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cell_id uuid NOT NULL REFERENCES service_matrix_cells(id) ON DELETE CASCADE,
  language language NOT NULL,
  value_text text NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (cell_id, language)
);

CREATE INDEX idx_service_matrix_cell_translations_cell ON service_matrix_cell_translations(cell_id);
CREATE INDEX idx_service_matrix_cell_translations_language ON service_matrix_cell_translations(language);
