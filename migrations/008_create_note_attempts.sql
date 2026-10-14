CREATE TABLE note_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_attempt_id UUID NOT NULL,
    position SMALLINT NOT NULL,
    target_pitch_class SMALLINT NOT NULL,
    target_octave SMALLINT NOT NULL,
    detected_pitch_class SMALLINT,
    detected_octave SMALLINT,
    target_frequency NUMERIC(10,4) NOT NULL,
    detected_frequency NUMERIC(10,4),
    cents_difference NUMERIC(8,2),
    accuracy NUMERIC(5,2) NOT NULL,
    result VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_note_attempts_game_attempt FOREIGN KEY (game_attempt_id) REFERENCES game_attempts(id) ON DELETE CASCADE,
    CONSTRAINT note_attempts_position CHECK (position > 0),
    CONSTRAINT note_attempts_target_pitch CHECK (target_pitch_class BETWEEN 0 AND 11),
    CONSTRAINT note_attempts_target_octave CHECK (target_octave BETWEEN 0 AND 10),
    CONSTRAINT note_attempts_detected_pitch CHECK (detected_pitch_class IS NULL OR detected_pitch_class BETWEEN 0 AND 11),
    CONSTRAINT note_attempts_accuracy CHECK (accuracy >= 0 AND accuracy <= 100),
    CONSTRAINT note_attempts_result CHECK (result IN ('HIT', 'MISS'))
);

CREATE INDEX idx_note_attempts_game_attempt ON note_attempts(game_attempt_id);
