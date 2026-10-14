CREATE TABLE level_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_id UUID NOT NULL,
    position SMALLINT NOT NULL,
    pitch_class SMALLINT NOT NULL,
    octave SMALLINT NOT NULL,
    tolerance_cents NUMERIC(6,2) NOT NULL DEFAULT 35.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_level_notes_level FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE,
    CONSTRAINT level_notes_position_positive CHECK (position > 0),
    CONSTRAINT level_notes_pitch_class CHECK (pitch_class BETWEEN 0 AND 11),
    CONSTRAINT level_notes_octave CHECK (octave BETWEEN 0 AND 10),
    CONSTRAINT level_notes_tolerance CHECK (tolerance_cents > 0 AND tolerance_cents <= 100),
    CONSTRAINT level_notes_unique_position UNIQUE (level_id, position)
);

CREATE INDEX idx_level_notes_level ON level_notes(level_id);
