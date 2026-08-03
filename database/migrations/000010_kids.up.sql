-- Initial 'onboarding' account status that precedes 'pending'. New signups start here
-- and only become 'pending' (approvable) once they've chosen their account type and
-- added their kids — so an admin never approves an account before knowing whether it's
-- a participating member or a guardian. Existing users are unaffected. (Safe to add in
-- this migration: the value is not used here, only later by application code.)
ALTER TYPE user_status ADD VALUE IF NOT EXISTS 'onboarding';

-- Account-level flag: does the account holder participate themselves, or are they
-- only a guardian managing kids? Defaults true so every existing account is a
-- participant (backwards-compatible). Guardian ("only my kids") accounts set false.
ALTER TABLE users ADD COLUMN self_participates BOOLEAN NOT NULL DEFAULT true;

-- Kids are participants managed by a parent account. They never authenticate — no
-- email, no roles, no auth of their own. Reuses the user_status enum for the
-- approval lifecycle. given_name/family_name are nullable, mirroring users (name is
-- captured first, richer details filled in later). deleted_at supports soft-delete:
-- removing a kid stamps deleted_at instead of dropping the row, so the kid's past event
-- responses survive and stay visible in the attendee lists of events they already
-- responded to; a removed kid is hidden from management and can't receive new responses.
CREATE TABLE IF NOT EXISTS kids
(
    id          UUID        DEFAULT gen_random_uuid()               PRIMARY KEY,
    user_id     UUID        REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    given_name  VARCHAR(100),
    family_name VARCHAR(100),
    picture     VARCHAR(500),
    status      user_status DEFAULT 'pending'                       NOT NULL,
    created_at  TIMESTAMP   DEFAULT current_timestamp               NOT NULL,
    deleted_at  TIMESTAMP
);

-- A response is now about exactly one subject: the account holder (user_id) OR a
-- kid (kid_id), never both, never neither.
ALTER TABLE event_responses ADD COLUMN kid_id UUID REFERENCES kids (id) ON DELETE CASCADE;
ALTER TABLE event_responses ALTER COLUMN user_id DROP NOT NULL;

-- Exactly one of user_id / kid_id is set (portable across Postgres and CockroachDB;
-- equivalent to num_nonnulls(user_id, kid_id) = 1 for two columns).
ALTER TABLE event_responses ADD CONSTRAINT event_responses_one_subject
    CHECK ((user_id IS NULL) <> (kid_id IS NULL));

-- Uniqueness: one response per person per event. The existing plain UNIQUE(event_id,
-- user_id) from 000003 is left untouched — it already enforces one own-RSVP per account,
-- and it never constrained kid rows (their user_id is NULL and NULLs are distinct). We
-- only ADD one partial unique index for kids, so the migration stays a clean,
-- backwards-compatible expand: the old backend's own-RSVP upsert
-- (`ON CONFLICT (event_id, user_id)`) keeps matching the same constraint during the
-- rolling deploy window, and only the new kid upsert uses the new arbiter.
CREATE UNIQUE INDEX event_responses_event_kid_key
    ON event_responses (event_id, kid_id) WHERE kid_id IS NOT NULL;
