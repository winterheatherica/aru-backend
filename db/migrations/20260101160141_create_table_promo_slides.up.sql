CREATE TABLE promo_slides (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_path text NOT NULL,
  order_index integer DEFAULT 0,
  is_active boolean DEFAULT true,
  uploaded_by uuid REFERENCES users(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX idx_promo_slides_is_active ON promo_slides(is_active);
CREATE INDEX idx_promo_slides_order_index ON promo_slides(order_index);
