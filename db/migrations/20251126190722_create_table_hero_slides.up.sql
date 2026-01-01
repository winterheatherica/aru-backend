CREATE TABLE hero_slides (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

  image_path text NOT NULL,
  order_index integer DEFAULT 0,
  is_active boolean DEFAULT true,

  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,

  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_hero_slides_order_index ON hero_slides(order_index);
CREATE INDEX idx_hero_slides_is_active ON hero_slides(is_active);
