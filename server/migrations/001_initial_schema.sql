-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE users (
    id         UUID PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE journeys (
    id            UUID PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title         TEXT NOT NULL,
    state         TEXT NOT NULL CHECK (state IN ('INITIALIZING','ACTIVE','PAUSED','COMPLETED')),
    deadline      TIMESTAMPTZ NOT NULL,
    config        JSONB NOT NULL DEFAULT '{}',
    progress_pct  SMALLINT NOT NULL DEFAULT 0 CHECK (progress_pct BETWEEN 0 AND 100),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_journeys_user_id ON journeys (user_id);
CREATE INDEX idx_journeys_user_state ON journeys (user_id, state);

CREATE TRIGGER journeys_set_updated_at
    BEFORE UPDATE ON journeys
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE milestones (
    id              UUID PRIMARY KEY,
    journey_id      UUID NOT NULL REFERENCES journeys(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    sequence_order  INT NOT NULL CHECK (sequence_order >= 0),
    state           TEXT NOT NULL CHECK (state IN ('LOCKED','ACTIVE','FINISHED')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT milestones_journey_sequence_uq
        UNIQUE (journey_id, sequence_order)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_milestones_journey_id ON milestones (journey_id);

CREATE TRIGGER milestones_set_updated_at
    BEFORE UPDATE ON milestones
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE tasks (
    id              UUID PRIMARY KEY,
    milestone_id    UUID NOT NULL REFERENCES milestones(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    scheduled_at    TIMESTAMPTZ NOT NULL,
    state           TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (state IN ('PENDING','EVALUATING','COMPLETED','MISSED','SKIPPED')),
    details         JSONB NOT NULL DEFAULT '{}',
    proof_of_work   JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_milestone_id ON tasks (milestone_id);
CREATE INDEX idx_tasks_scheduled_at ON tasks (scheduled_at);
CREATE INDEX idx_tasks_milestone_state ON tasks (milestone_id, state);

CREATE TRIGGER tasks_set_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE ai_mentors (
    id              UUID PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    expertise_area  VARCHAR(50) NOT NULL,
    system_prompt   TEXT NOT NULL,
    model_provider  VARCHAR(50) NOT NULL,
    model_version   VARCHAR(50) NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_mentors_active ON ai_mentors (is_active) WHERE is_active = true;
CREATE INDEX idx_ai_mentors_expertise ON ai_mentors (expertise_area);

CREATE TABLE conversations (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ai_mentor_id    UUID NOT NULL REFERENCES ai_mentors(id) ON DELETE RESTRICT,
    journey_id      UUID REFERENCES journeys(id) ON DELETE SET NULL,
    title           VARCHAR(255) NOT NULL DEFAULT 'New conversation',
    is_archived     BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_conversations_user_id ON conversations (user_id);
CREATE INDEX idx_conversations_user_updated ON conversations (user_id, updated_at DESC);
CREATE INDEX idx_conversations_journey_id ON conversations (journey_id) WHERE journey_id IS NOT NULL;

CREATE TRIGGER conversations_set_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE messages (
    id                UUID PRIMARY KEY,
    conversation_id   UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role              TEXT NOT NULL CHECK (role IN ('user','assistant','system')),
    content           TEXT NOT NULL,
    tokens_used       INT NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_conversation_created
    ON messages (conversation_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS ai_mentors;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS journeys;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS set_updated_at();
