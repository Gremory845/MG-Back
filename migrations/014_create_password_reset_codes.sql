-- Códigos de recuperación de contraseña.
-- El código de 6 dígitos NUNCA se guarda en claro: solo su HMAC.
-- Flujo: forgot-password crea la fila -> verify-reset-code la marca verificada y
-- guarda el hash del token de restablecimiento -> reset-password la consume.
CREATE TABLE password_reset_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    attempts SMALLINT NOT NULL DEFAULT 0,
    verified_at TIMESTAMPTZ,
    reset_token_hash TEXT,
    reset_expires_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_password_reset_codes_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT password_reset_codes_attempts CHECK (attempts >= 0),
    CONSTRAINT password_reset_codes_token_unique UNIQUE (reset_token_hash)
);

CREATE INDEX idx_password_reset_codes_user_created
    ON password_reset_codes(user_id, created_at DESC);
