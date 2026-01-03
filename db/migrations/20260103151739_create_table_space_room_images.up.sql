CREATE TABLE space_room_images (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id uuid NOT NULL REFERENCES space_rooms(id) ON DELETE CASCADE,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  image_url text NOT NULL,
  is_active boolean DEFAULT true,
  order_index integer DEFAULT 1,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_space_room_images_room ON space_room_images(room_id);
CREATE INDEX idx_space_room_images_active ON space_room_images(is_active);
CREATE INDEX idx_space_room_images_order ON space_room_images(room_id, order_index);
