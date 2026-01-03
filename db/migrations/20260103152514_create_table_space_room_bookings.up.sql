CREATE TABLE space_room_bookings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id uuid NOT NULL REFERENCES space_rooms(id) ON DELETE CASCADE,
  start_time timestamptz NOT NULL,
  end_time   timestamptz NOT NULL,
  status room_booking_status NOT NULL DEFAULT 'PENDING',
  booked_by uuid REFERENCES users(id) ON DELETE SET NULL,
  added_by  uuid REFERENCES users(id) ON DELETE SET NULL,
  purpose text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CHECK (end_time > start_time)
);

CREATE INDEX idx_space_room_bookings_room_start ON space_room_bookings (room_id, start_time);
CREATE INDEX idx_space_room_bookings_room_end ON space_room_bookings (room_id, end_time);
CREATE INDEX idx_space_room_bookings_status ON space_room_bookings (status);
