CREATE TABLE IF NOT EXISTS roles (
    id         UUID      DEFAULT gen_random_uuid()  PRIMARY KEY,
    name       VARCHAR(100) UNIQUE                  NOT NULL,
    protected  BOOLEAN                              NOT NULL DEFAULT false,
    created_by UUID      REFERENCES users(id),
    created_at TIMESTAMP DEFAULT current_timestamp  NOT NULL
);

-- permission_key is a plain VARCHAR rather than a FK into a permissions table because
-- permissions are a fixed catalog defined in Go constants — admins can assign them but
-- cannot create new ones, so a separate table would just mirror the code with no benefit.
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id        UUID         NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_key VARCHAR(100) NOT NULL,
    PRIMARY KEY (role_id, permission_key)
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS category_creator_roles (
    category_id UUID NOT NULL REFERENCES event_categories(id) ON DELETE CASCADE,
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, role_id)
);

INSERT INTO roles (name, protected) VALUES ('admin', true), ('member', true) ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_key)
SELECT id, 'roles.manage' FROM roles WHERE name = 'admin'
ON CONFLICT (role_id, permission_key) DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r WHERE r.name = 'member'
ON CONFLICT (user_id, role_id) DO NOTHING;
