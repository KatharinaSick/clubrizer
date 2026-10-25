# Server

Contains the Clubrizer backend.

## Known Limitations

- **First admin must be set manually**: There is no automated bootstrap for the initial admin user. Grant the admin role directly in the database:
  ```sql
  INSERT INTO user_roles (user_id, role_id)
  SELECT u.id, r.id FROM users u, roles r
  WHERE u.email = 'your@email.com' AND r.name = 'admin';
  ```

## Push Notifications

Push notifications use standard Web Push with a VAPID key pair, read from the `VAPID_PUBLIC_KEY` and `VAPID_PRIVATE_KEY` environment variables. Generate a pair with:

```sh
go run ./cmd/vapidkeys
```

- **Locally**: add both lines to the root `.env`.
- **Production**: store them as the `VAPID_PUBLIC_KEY` / `VAPID_PRIVATE_KEY` secrets in Google Secret Manager (readable by the backend's service account). Use a separate pair from your local one.
- **Never change the production keys** once people have turned notifications on — every registered device is tied to the key pair, so new keys mean everyone has to turn notifications on again.

## Useful Resources

- Structure: 
  - https://go.dev/doc/modules/layout#server-project
  - https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
- Server: 
  - net/http intro: https://www.youtube.com/watch?v=H7tbjKFSg58
  - generic handlers: https://www.willem.dev/articles/generic-http-handlers/
