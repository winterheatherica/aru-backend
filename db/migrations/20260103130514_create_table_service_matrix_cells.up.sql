CREATE TABLE service_matrix_cells (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  row_id uuid NOT NULL REFERENCES service_matrix_rows(id) ON DELETE CASCADE,
  column_id uuid NOT NULL REFERENCES service_matrix_columns(id) ON DELETE CASCADE,
  value_boolean boolean,
  value_number numeric,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (row_id, column_id)
);

CREATE INDEX idx_service_matrix_cells_row ON service_matrix_cells(row_id);
CREATE INDEX idx_service_matrix_cells_column ON service_matrix_cells(column_id);
