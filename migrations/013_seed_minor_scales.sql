-- Seed: escalas menores (natural, armónica y melódica) en las tónicas C, D, E, F, G, A y B.
-- Misma estructura que las mayores: 3 niveles por escala
--   Nivel 1: grados 1-3 | Nivel 2: grados 1-4 | Nivel 3: escala completa (8 notas)
-- Las notas se calculan desde los intervalos de cada definición.
-- Octava inicial: A y B arrancan en octava 3; el resto en octava 4.
-- La melódica usa la forma ascendente (6º y 7º grado elevados).

INSERT INTO scale_definitions (id, code, name, intervals, note_count, category, family) VALUES
  ('a0000000-0000-0000-0000-000000000002', 'NATURAL_MINOR',  'Natural Minor',  ARRAY[0,2,3,5,7,8,10,12], 8, 'MINOR', 'NATURAL'),
  ('a0000000-0000-0000-0000-000000000003', 'HARMONIC_MINOR', 'Harmonic Minor', ARRAY[0,2,3,5,7,8,11,12], 8, 'MINOR', 'HARMONIC'),
  ('a0000000-0000-0000-0000-000000000004', 'MELODIC_MINOR',  'Melodic Minor',  ARRAY[0,2,3,5,7,9,11,12], 8, 'MINOR', 'MELODIC');

-- Dificultad = dificultad de la tónica (igual que en las mayores) + bonus por familia
INSERT INTO scales (scale_definition_id, root_pitch_class, name, difficulty)
SELECT d.id, r.pc, r.letter || ' ' || d.name,
       LEAST(10, r.base + CASE d.family WHEN 'NATURAL' THEN 1 WHEN 'MELODIC' THEN 2 ELSE 3 END)
FROM scale_definitions d
CROSS JOIN (VALUES (0, 'C', 1), (7, 'G', 2), (5, 'F', 2), (2, 'D', 3), (9, 'A', 4), (4, 'E', 5), (11, 'B', 6)) AS r(pc, letter, base)
WHERE d.category = 'MINOR';

INSERT INTO levels (scale_id, level_number, name, required_accuracy)
SELECT s.id, lv.n, s.name || ' - Nivel ' || lv.n, 70.00
FROM scales s
JOIN scale_definitions d ON d.id = s.scale_definition_id
CROSS JOIN generate_series(1, 3) AS lv(n)
WHERE d.category = 'MINOR';

INSERT INTO level_notes (level_id, position, pitch_class, octave, tolerance_cents)
SELECT
  l.id,
  p.pos,
  (s.root_pitch_class + d.intervals[p.pos]) % 12,
  (12 * ((CASE WHEN s.root_pitch_class >= 9 THEN 3 ELSE 4 END) + 1)
     + s.root_pitch_class + d.intervals[p.pos]) / 12 - 1,
  35.00
FROM levels l
JOIN scales s ON s.id = l.scale_id
JOIN scale_definitions d ON d.id = s.scale_definition_id
CROSS JOIN LATERAL generate_series(
  1,
  CASE l.level_number WHEN 1 THEN 3 WHEN 2 THEN 4 ELSE 8 END
) AS p(pos)
WHERE d.category = 'MINOR';
