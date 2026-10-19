-- Agrupa las definiciones de escala para poder navegarlas por tarjetas:
--   category: MAJOR | MINOR
--   family:   NATURAL | HARMONIC | MELODIC

ALTER TABLE scale_definitions
  ADD COLUMN category VARCHAR(20),
  ADD COLUMN family VARCHAR(20);

UPDATE scale_definitions SET category = 'MAJOR', family = 'NATURAL' WHERE code = 'MAJOR';

ALTER TABLE scale_definitions
  ALTER COLUMN category SET NOT NULL,
  ALTER COLUMN family SET NOT NULL,
  ADD CONSTRAINT scale_definitions_category CHECK (category IN ('MAJOR', 'MINOR')),
  ADD CONSTRAINT scale_definitions_family CHECK (family IN ('NATURAL', 'HARMONIC', 'MELODIC'));
