-- Seed: Do Mayor con 3 niveles progresivos (UUIDs fijos para poder referenciarlos)

INSERT INTO scale_definitions (id, code, name, intervals, note_count) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'MAJOR', 'Major', ARRAY[0,2,4,5,7,9,11,12], 8);

INSERT INTO scales (id, scale_definition_id, root_pitch_class, name, difficulty) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 0, 'C Major', 1);

INSERT INTO levels (id, scale_id, level_number, name, required_accuracy) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 1, 'C Major - Nivel 1', 70.00),
  ('c0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', 2, 'C Major - Nivel 2', 70.00),
  ('c0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001', 3, 'C Major - Nivel 3', 70.00);

-- Nivel 1: C D E
INSERT INTO level_notes (level_id, position, pitch_class, octave, tolerance_cents) VALUES
  ('c0000000-0000-0000-0000-000000000001', 1, 0, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000001', 2, 2, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000001', 3, 4, 4, 35.00);

-- Nivel 2: C D E F
INSERT INTO level_notes (level_id, position, pitch_class, octave, tolerance_cents) VALUES
  ('c0000000-0000-0000-0000-000000000002', 1, 0, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000002', 2, 2, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000002', 3, 4, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000002', 4, 5, 4, 35.00);

-- Nivel 3: escala completa C D E F G A B C
INSERT INTO level_notes (level_id, position, pitch_class, octave, tolerance_cents) VALUES
  ('c0000000-0000-0000-0000-000000000003', 1, 0, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 2, 2, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 3, 4, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 4, 5, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 5, 7, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 6, 9, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 7, 11, 4, 35.00),
  ('c0000000-0000-0000-0000-000000000003', 8, 0, 5, 35.00);
