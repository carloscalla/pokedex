CREATE TABLE IF NOT EXISTS pokemon (
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    sprite_url TEXT,
    height     INTEGER NOT NULL,
    weight     INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_pokemon_name ON pokemon(name);

CREATE TABLE IF NOT EXISTS pokemon_types (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    type       TEXT NOT NULL,
    slot       INTEGER NOT NULL,
    PRIMARY KEY (pokemon_id, slot)
);
CREATE INDEX IF NOT EXISTS idx_pokemon_types_type ON pokemon_types(type);

CREATE TABLE IF NOT EXISTS pokemon_stats (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    stat_name  TEXT NOT NULL,
    base_stat  INTEGER NOT NULL,
    PRIMARY KEY (pokemon_id, stat_name)
);
CREATE INDEX IF NOT EXISTS idx_pokemon_stats_name_value ON pokemon_stats(stat_name, base_stat);

CREATE TABLE IF NOT EXISTS pokemon_abilities (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    ability    TEXT NOT NULL,
    slot       INTEGER NOT NULL,
    PRIMARY KEY (pokemon_id, slot)
);

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS team_members (
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    position   INTEGER NOT NULL,
    added_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, pokemon_id)
);
CREATE INDEX IF NOT EXISTS idx_team_members_user ON team_members(user_id);
