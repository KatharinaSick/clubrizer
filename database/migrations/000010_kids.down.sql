DROP INDEX IF EXISTS event_responses_event_kid_key;
ALTER TABLE event_responses DROP CONSTRAINT IF EXISTS event_responses_one_subject;
-- Remove kid responses (user_id NULL) before restoring the NOT NULL constraint.
DELETE FROM event_responses WHERE user_id IS NULL;
ALTER TABLE event_responses DROP COLUMN kid_id;
ALTER TABLE event_responses ALTER COLUMN user_id SET NOT NULL;
-- The original UNIQUE(event_id, user_id) from 000003 was never touched by the up
-- migration, so there is nothing to restore here.

DROP TABLE IF EXISTS kids;

ALTER TABLE users DROP COLUMN self_participates;
