CREATE TABLE service_matrix_columns (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  matrix_id uuid NOT NULL REFERENCES service_matrices(id) ON DELETE CASCADE,
  column_key text NOT NULL,
  popular boolean DEFAULT false,
  order_index integer DEFAULT 1,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (matrix_id, column_key)
);

CREATE INDEX idx_service_matrix_columns_matrix ON service_matrix_columns(matrix_id);
CREATE INDEX idx_service_matrix_columns_order ON service_matrix_columns(matrix_id, order_index);
CREATE INDEX idx_service_matrix_columns_popular ON service_matrix_columns(popular);
