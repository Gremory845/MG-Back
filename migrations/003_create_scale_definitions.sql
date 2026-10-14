-- scale_types + scale_definitions fusionadas en una sola tabla
CREATE TABLE scale_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    intervals SMALLINT[] NOT NULL,
    note_count SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT scale_definitions_intervals_not_empty CHECK (cardinality(intervals) > 0),
    CONSTRAINT scale_definitions_note_count_positive CHECK (note_count > 0),
    CONSTRAINT scale_definitions_note_count_matches CHECK (note_count = cardinality(intervals))
);
