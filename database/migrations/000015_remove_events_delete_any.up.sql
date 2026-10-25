-- Events can no longer be deleted (they're cancelled instead), so the events.delete_any
-- permission was removed from the catalog in backend/internal/rbac/permissions.go. This
-- removes any grants of it that roles still hold. Safe while the old backend is still
-- running: without the grant, deleting other people's events is limited to admins and
-- event owners for the few minutes until the new backend is live.
DELETE FROM role_permissions WHERE permission_key = 'events.delete_any';
