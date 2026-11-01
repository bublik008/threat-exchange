-- Таблица угроз. MD5 — естественный первичный ключ: он уникален по определению
-- и по нему идёт вся выборка в API (GET/PUT /info/{md5})
CREATE TABLE IF NOT EXISTS threats (
    md5         CHAR(32) PRIMARY KEY,
    name        TEXT NOT NULL,
    severity    TEXT NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    threat_type TEXT,
    description TEXT,
    tags        TEXT[] NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);