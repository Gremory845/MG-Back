CREATE TABLE levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scale_id UUID NOT NULL,
    level_number SMALLINT NOT NULL,
    name VARCHAR(100),
    required_accuracy NUMERIC(5,2) NOT NULL DEFAULT 70.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_levels_scale FOREIGN KEY (scale_id) REFERENCES scales(id) ON DELETE CASCADE,
    CONSTRAINT levels_number_positive CHECK (level_number > 0),
    CONSTRAINT levels_accuracy CHECK (required_accuracy >= 0 AND required_accuracy <= 100),
    CONSTRAINT levels_unique_scale_number UNIQUE (scale_id, level_number)
);

CREATE INDEX idx_levels_scale ON levels(scale_id);
