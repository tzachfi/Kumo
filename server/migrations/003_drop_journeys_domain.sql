-- +goose Up
ALTER TABLE journeys DROP COLUMN IF EXISTS domain;

-- +goose Down
ALTER TABLE journeys
    ADD COLUMN domain TEXT NOT NULL DEFAULT 'FITNESS'
        CHECK (domain IN ('FITNESS','TECHNICAL','CREATIVE'));
