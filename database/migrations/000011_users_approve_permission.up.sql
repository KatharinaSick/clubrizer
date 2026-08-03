-- Grants the new users.approve permission (the admin approval queue) to the admin role.
-- Permissions are a fixed catalog defined in Go constants (see
-- backend/internal/rbac/permissions.go); this migration only wires the existing admin
-- role to the new key. Backwards-compatible: the old backend never checks this
-- permission, so adding the grant early is a safe expand.
INSERT INTO role_permissions (role_id, permission_key)
SELECT id, 'users.approve' FROM roles WHERE name = 'admin'
ON CONFLICT (role_id, permission_key) DO NOTHING;
