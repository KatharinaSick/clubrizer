package rbac

// Permissions are a fixed catalog defined here in code. Admins can assign them to roles
// but cannot create new ones, so there is no permissions table in the DB — these constants
// are the single source of truth.
const (
	PermissionRolesManage = "roles.manage"
	// PermissionEventsEditAny allows editing events created by other users.
	// Admins bypass this check; event owners can always edit their own events.
	PermissionEventsEditAny = "events.edit_any"
	// PermissionEventsCancelAny allows cancelling events created by other users.
	// Admins bypass this check; event owners can always cancel their own events.
	PermissionEventsCancelAny = "events.cancel_any"
	// PermissionUsersApprove allows reviewing the approval queue and approving or
	// rejecting pending accounts and kids. Admins bypass this check.
	PermissionUsersApprove = "users.approve"
)

// AllPermissions is the full catalog. Used to expand an admin — who bypasses every
// individual permission check — to the complete set when computing a user's effective
// permissions. Keep this in sync when adding a permission constant above.
var AllPermissions = []string{
	PermissionRolesManage,
	PermissionEventsEditAny,
	PermissionEventsCancelAny,
	PermissionUsersApprove,
}
