# Postman collection

| File | What |
|---|---|
| `simple-auth.postman_collection.json` | All endpoints, grouped in folders, with test scripts |
| `local.postman_environment.json` | `baseUrl` (http://localhost:8080), test user name/password, admin email/password |

## Import

1. Postman → **Import** → drop both JSON files.
2. Select the **simple-auth local** environment (top-right dropdown).

## Use

Start the API first:

```bash
make db-up && make migrate-up && make dev
```

**Once**, create the admin the collection logs in as (matches `adminEmail` / `adminPassword`):

```bash
make seed-admin email=admin@example.com password=admin12345
```

- **Manually:** run **01 Auth → Sign up**, **Log in** and **Log in as admin**. Both tokens are saved automatically.
- **Everything at once:** right-click the collection → **Run collection** → **Run**. Requests run top to bottom and each one checks its own result.

## How it's set up

- **Auth is defined once on the collection** (`Bearer {{accessToken}}`, the regular user). Requests inherit it. Admin-only requests (marked *(admin)*) use `Bearer {{adminToken}}`. Public ones (health, signup, login) use *No Auth*.
- **Scripts fill variables for you:**

  | Variable | Set by |
  |---|---|
  | `userEmail` | Sign up (pre-request: unique email per run) |
  | `userId` | Sign up |
  | `accessToken` | Log in |
  | `otherUserEmail` | Auth → Negative cases → Sign up – 'role: admin' is ignored |
  | `managedUserEmail`, `managedUserId` | Users → Create user (admin) |
  | `adminToken`, `adminId` | Log in as admin |

- **Collection-level tests** run after every request: response time < 2s, JSON content type, and the standard error shape `{ "error": { "code", "message" } }` on every 4xx/5xx.
- **Negative cases** folders check that bad input, missing or invalid tokens, **forbidden actions (403)**, duplicates and unknown ids fail with the right status and error code.

## Run from the terminal (optional)

With [Newman](https://www.npmjs.com/package/newman), Postman's CLI runner:

```bash
npx newman run postman/simple-auth.postman_collection.json -e postman/local.postman_environment.json
```
