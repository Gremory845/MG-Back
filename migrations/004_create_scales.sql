CREATE TABLE scales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scale_definition_id UUID NOT NULL,
    root_pitch_class SMALLINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    difficulty SMALLINT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_scales_definition FOREIGN KEY (scale_definition_id) REFERENCES scale_definitions(id),
    CONSTRAINT scales_root_pitch_class CHECK (root_pitch_class BETWEEN 0 AND 11),
    CONSTRAINT scales_difficulty CHECK (difficulty BETWEEN 1 AND 10),
    CONSTRAINT scales_unique_root_definition UNIQUE (scale_definition_id, root_pitch_class)
);

CREATE INDEX idx_scales_definition ON scales(scale_definition_id);
