-- Seed: escalas mayores D, E, F, G, A y B (Do Mayor ya existe en 010).
-- Misma estructura que Do Mayor: 3 niveles por escala
--   Nivel 1: grados 1-3 | Nivel 2: grados 1-4 | Nivel 3: escala completa (8 notas)
-- Las notas se calculan desde los intervalos de la definición MAJOR, no se escriben a mano.
-- Octava inicial: A y B arrancan en octava 3 (A3, B3); el resto en octava 4.
-- Así todas las escalas quedan centradas cerca del Do central.

INSERT INTO scales (id, scale_definition_id, root_pitch_class, name, difficulty) VALUES
  ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 7,  'G Major', 2),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 5,  'F Major', 2),
  ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 2,  'D Major', 3),
  ('b0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 9,  'A Major', 4),
  ('b0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 4,  'E Major', 5),
  ('b0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 11, 'B Major', 6);

INSERT INTO levels (scale_id, level_number, name, required_accuracy)
SELECT s.id, lv.n, s.name || ' - Nivel ' || lv.n, 70.00
FROM scales s
CROSS JOIN generate_series(1, 3) AS lv(n)
WHERE s.id IN (
  'b0000000-0000-0000-0000-000000000002',
  'b0000000-0000-0000-0000-000000000003',
  'b0000000-0000-0000-0000-000000000004',
  'b0000000-0000-0000-0000-000000000005',
  'b0000000-0000-0000-0000-000000000006',
  'b0000000-0000-0000-0000-000000000007'
);

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
WHERE d.code = 'MAJOR'
  AND s.id IN (
    'b0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000004',
    'b0000000-0000-0000-0000-000000000005',
    'b0000000-0000-0000-0000-000000000006',
    'b0000000-0000-0000-0000-000000000007'
  );
