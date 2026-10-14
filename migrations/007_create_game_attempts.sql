CREATE TABLE game_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    level_id UUID NOT NULL,
    score INTEGER NOT NULL DEFAULT 0,
    accuracy NUMERIC(5,2) NOT NULL DEFAULT 0,
    mistakes SMALLINT NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    duration_ms INTEGER,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_game_attempts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_game_attempts_level FOREIGN KEY (level_id) REFERENCES levels(id),
    CONSTRAINT game_attempts_score CHECK (score >= 0),
    CONSTRAINT game_attempts_accuracy CHECK (accuracy >= 0 AND accuracy <= 100),
    CONSTRAINT game_attempts_mistakes CHECK (mistakes >= 0),
    CONSTRAINT game_attempts_duration CHECK (duration_ms IS NULL OR duration_ms >= 0),
    CONSTRAINT game_attempts_finished CHECK ((completed = FALSE) OR finished_at IS NOT NULL)
);

CREATE INDEX idx_game_attempts_user ON game_attempts(user_id);
CREATE INDEX idx_game_attempts_user_level ON game_attempts(user_id, level_id);
CREATE INDEX idx_game_attempts_leaderboard
    ON game_attempts(level_id, score DESC, accuracy DESC)
    WHERE completed = TRUE;
