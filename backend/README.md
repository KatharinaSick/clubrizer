# Server

Contains the Clubrizer backend.

## Known Limitations

- **First admin must be set manually**: There is no automated bootstrap for the initial admin user. Grant the admin role directly in the database:
  ```sql
  INSERT INTO user_roles (user_id, role_id)
  SELECT u.id, r.id FROM users u, roles r
  WHERE u.email = 'your@email.com' AND r.name = 'admin';
  ```

## Useful Resources

- Structure: 
  - https://go.dev/doc/modules/layout#server-project
  - https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
- Server: 
  - net/http intro: https://www.youtube.com/watch?v=H7tbjKFSg58
  - generic handlers: https://www.willem.dev/articles/generic-http-handlers/
