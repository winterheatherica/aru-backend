CREATE TABLE space_rooms (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code text NOT NULL,
  capacity integer,
  floor integer,
  address text,
  main_image_url text,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (code)
);

CREATE INDEX idx_space_rooms_active ON space_rooms(is_active);
CREATE INDEX idx_space_rooms_floor ON space_rooms(floor);