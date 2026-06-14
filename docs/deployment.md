# Deployment

## Architecture

| Layer    | Service                         |
|----------|---------------------------------|
| Frontend | Netlify (`lisc-2010.clubrizer.com`) |
| Backend  | Google Cloud Run (`clubrizer-backend`, `europe-west1`) |
| Database | CockroachDB Serverless (`clubrizer-16993.jxf.gcp-europe-west3`) |
| Storage  | Google Cloud Storage (`clubrizer_profile-pictures`) |
| Email    | Resend (`noreply@clubrizer.com`) |

`lisc-2010.clubrizer.com/api/*` is proxied by Netlify to Cloud Run (stripping the `/api` prefix), so the browser always talks to a single origin — no CORS required.

## CI/CD

Every push to `main` triggers `.github/workflows/deploy.yml`, which:

1. **Runs database migrations** against CockroachDB using `golang-migrate`. Migrations must always be backwards-compatible with the previous backend version — see `AGENTS.md` for the expand/contract rule.
2. **Builds a Docker image** of the Go backend and pushes it to Artifact Registry (`europe-west1-docker.pkg.dev/clubrizer-gcloud/clubrizer/backend`).
3. **Deploys to Cloud Run** with the new image. Secrets are read at runtime from Google Secret Manager — GitHub Actions never sees their values.

Authentication from GitHub Actions to GCP uses **Workload Identity Federation** (no long-lived keys).

## GCP Resources

| Resource | Name |
|---|---|
| Project | `clubrizer-gcloud` |
| Artifact Registry repo | `europe-west1-docker.pkg.dev/clubrizer-gcloud/clubrizer` |
| Cloud Run service | `clubrizer-backend` (region: `europe-west1`) |
| GitHub Actions SA | `github-actions@clubrizer-gcloud.iam.gserviceaccount.com` |
| Cloud Run runtime SA | `clubrizer-backend@clubrizer-gcloud.iam.gserviceaccount.com` |
| WIF pool | `projects/281221302650/locations/global/workloadIdentityPools/github` |
| WIF provider | `projects/281221302650/locations/global/workloadIdentityPools/github/providers/github-actions` |

## Secrets & Config

### Google Secret Manager
Sensitive runtime secrets — only the Cloud Run service account can read these:

| Secret name | Description |
|---|---|
| `DB_URL` | CockroachDB connection string |
| `JWT_ACCESS_TOKEN_SECRET_KEY` | JWT signing key for access tokens |
| `JWT_REFRESH_TOKEN_SECRET_KEY` | JWT signing key for refresh tokens |
| `RESEND_API_KEY` | Resend API key |

To add a new version of a secret:
```bash
echo -n "new-value" | gcloud secrets versions add SECRET_NAME --data-file=- --project=clubrizer-gcloud
```

### GitHub Secrets
Only what GitHub Actions needs to authenticate and run migrations:

| Secret | Description |
|---|---|
| `WIF_PROVIDER` | Workload Identity Federation provider resource name |
| `WIF_SERVICE_ACCOUNT` | GitHub Actions service account email |
| `DB_URL` | CockroachDB connection string (needed for migrations) |

### Non-sensitive config
Hardcoded in `backend/config.yaml` — no secrets needed:
- `RESEND_FROM_ADDRESS`: `noreply@clubrizer.com`
- `GCS_PROFILE_PICTURES_BUCKET_NAME`: `clubrizer_profile-pictures`

## One-time Setup (already done)

Documented here for reference if the environment ever needs to be recreated.

### GCP
```bash
# Enable APIs
gcloud services enable run.googleapis.com artifactregistry.googleapis.com iamcredentials.googleapis.com secretmanager.googleapis.com --project=clubrizer-gcloud

# Artifact Registry
gcloud artifacts repositories create clubrizer --repository-format=docker --location=europe-west1 --project=clubrizer-gcloud

# Service accounts
gcloud iam service-accounts create github-actions --display-name="GitHub Actions" --project=clubrizer-gcloud
gcloud iam service-accounts create clubrizer-backend --display-name="Clubrizer Backend" --project=clubrizer-gcloud

# Workload Identity Federation
gcloud iam workload-identity-pools create github --location=global --display-name="GitHub Actions" --project=clubrizer-gcloud
gcloud iam workload-identity-pools providers create-oidc github-actions --location=global --workload-identity-pool=github --issuer-uri="https://token.actions.githubusercontent.com" --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" --attribute-condition="assertion.repository=='KatharinaSick/clubrizer'" --project=clubrizer-gcloud

# IAM bindings
gcloud iam service-accounts add-iam-policy-binding github-actions@clubrizer-gcloud.iam.gserviceaccount.com --role="roles/iam.workloadIdentityUser" --member="principalSet://iam.googleapis.com/projects/281221302650/locations/global/workloadIdentityPools/github/attribute.repository/KatharinaSick/clubrizer" --project=clubrizer-gcloud
gcloud projects add-iam-policy-binding clubrizer-gcloud --member="serviceAccount:github-actions@clubrizer-gcloud.iam.gserviceaccount.com" --role="roles/run.admin"
gcloud artifacts repositories add-iam-policy-binding clubrizer --location=europe-west1 --member="serviceAccount:github-actions@clubrizer-gcloud.iam.gserviceaccount.com" --role="roles/artifactregistry.writer" --project=clubrizer-gcloud
gcloud iam service-accounts add-iam-policy-binding clubrizer-backend@clubrizer-gcloud.iam.gserviceaccount.com --role="roles/iam.serviceAccountUser" --member="serviceAccount:github-actions@clubrizer-gcloud.iam.gserviceaccount.com" --project=clubrizer-gcloud
gcloud storage buckets add-iam-policy-binding gs://clubrizer_profile-pictures --member="serviceAccount:clubrizer-backend@clubrizer-gcloud.iam.gserviceaccount.com" --role="roles/storage.objectAdmin"
gcloud projects add-iam-policy-binding clubrizer-gcloud --member="serviceAccount:clubrizer-backend@clubrizer-gcloud.iam.gserviceaccount.com" --role="roles/secretmanager.secretAccessor"
```

### CockroachDB
Cluster connection details are in the `DB_URL` GitHub secret and Google Secret Manager entry.
Database: `clubrizer`. Cluster region: `gcp-europe-west3`.

Migrations are run automatically by CI. To run them manually:
```bash
devbox run migratedb
```
