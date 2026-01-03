CREATE TABLE service_matrix_rows (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  matrix_id uuid NOT NULL REFERENCES service_matrices(id) ON DELETE CASCADE,
  row_key text NOT NULL,
  order_index integer DEFAULT 1,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (matrix_id, row_key)
);

CREATE INDEX idx_service_matrix_rows_matrix ON service_matrix_rows(matrix_id);
CREATE INDEX idx_service_matrix_rows_order ON service_matrix_rows(matrix_id, order_index);
