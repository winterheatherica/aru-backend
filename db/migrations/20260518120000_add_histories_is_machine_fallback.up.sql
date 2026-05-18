ALTER TABLE histories
  ADD COLUMN is_machine_fallback boolean NOT NULL DEFAULT false;
