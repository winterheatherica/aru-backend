CREATE TABLE space_room_image_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_id uuid NOT NULL REFERENCES space_room_images(id) ON DELETE CASCADE,
  language language NOT NULL,
  alt text,
  title text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (image_id, language)
);

CREATE INDEX idx_space_room_image_translations_image ON space_room_image_translations(image_id);
CREATE INDEX idx_space_room_image_translations_language ON space_room_image_translations(language);
