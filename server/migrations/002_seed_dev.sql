-- +goose Up
INSERT INTO users (id, email)
VALUES (
    '00000000-0000-4000-8000-000000000001',
    'dev@localhost'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO ai_mentors (id, name, expertise_area, system_prompt, model_provider, model_version)
VALUES (
    '00000000-0000-4000-8000-000000000010',
    'Kumo',
    'GENERAL',
    'You are Kumo, a supportive learning mentor.',
    'mock',
    'v1'
)
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM ai_mentors WHERE id = '00000000-0000-4000-8000-000000000010';
DELETE FROM users WHERE id = '00000000-0000-4000-8000-000000000001';
