CREATE TABLE space_room_translations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id uuid NOT NULL REFERENCES space_rooms(id) ON DELETE CASCADE,
  language language NOT NULL,
  title text NOT NULL,
  description text,
  facilities text[],
  main_image_alt text,
  main_image_title text,
  slug text NOT NULL,
  meta_title text,
  meta_description text,
  meta_keywords text[],
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE (room_id, language),
  UNIQUE (language, slug)
);


CREATE INDEX idx_space_room_translations_room ON space_room_translations(room_id);
CREATE INDEX idx_space_room_translations_language ON space_room_translations(language);
