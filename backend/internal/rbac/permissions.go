package rbac

// Permissions are a fixed catalog defined here in code. Admins can assign them to roles
// but cannot create new ones, so there is no permissions table in the DB — these constants
// are the single source of truth.
const (
	PermissionRolesManage = "roles.manage"
	// PermissionEventsDeleteAny allows deleting events created by other users.
	// Admins bypass this check; event owners can always delete their own events.
	PermissionEventsDeleteAny = "events.delete_any"
)
