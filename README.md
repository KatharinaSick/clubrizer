[![Built with Devbox](https://www.jetify.com/img/devbox/shield_galaxy.svg)](https://www.jetify.com/devbox/docs/contributor-quickstart/)

# Clubrizer

A side/fun project that aims to create a platform for teams. More details will follow.

## Deployment

For production deployment setup and architecture, see [docs/deployment.md](docs/deployment.md).

## Known Limitations

- **Google Cloud dependency**: Profile picture storage uses Google Cloud Storage with Application Default Credentials. This currently assumes a GCP environment (locally via `gcloud auth application-default login`, in production via a Cloud Run service account). Support for other storage backends is not yet implemented.
