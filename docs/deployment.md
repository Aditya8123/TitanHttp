# Deployment Strategy

> How we ship TitanHTTP to the world.

While the primary focus of Phase 1 through 9 is building the core infrastructure, Phase 10 centers on getting TitanHTTP running in a production-like environment. This document outlines the high-level strategy for containerization, continuous integration, and final deployment.

## Containerization (Docker)

TitanHTTP will be deployed as a minimal Docker container. Because Go compiles to a single, statically linked binary, our final Docker image will be extremely small and secure.

### Multi-Stage Build Strategy
1. **Builder Stage:** Uses a full `golang` Alpine or Debian base image. It pulls dependencies, runs tests, and compiles the `titanhttp` binary with flags (`CGO_ENABLED=0`) to ensure a static build.
2. **Production Stage:** Uses the `scratch` (empty) or `alpine` base image. It copies only the compiled binary from the builder stage.

*Result:* An image size often less than 20MB, containing zero unnecessary OS utilities, drastically reducing the security attack surface.

## Continuous Integration (CI)

We will utilize GitHub Actions to automate our pipeline. Every push or pull request to the `main` branch will trigger the following checks:

1. **Linting:** Runs `golangci-lint` to enforce style and catch potential bugs early.
2. **Formatting:** Ensures all code adheres to `gofumpt`.
3. **Testing:** Executes the full suite of unit and integration tests (see [testing.md](./testing.md)).
4. **Build Verification:** Confirms the binary compiles successfully across target operating systems (Linux, Windows, macOS).

## Continuous Deployment (CD)

Once code is merged to `main` and all CI checks pass, a release process can be triggered:

1. **Tagging:** A semantic version tag (e.g., `v1.0.0`) is pushed.
2. **Image Build & Push:** GitHub Actions builds the Docker image and pushes it to a container registry (e.g., GitHub Container Registry or Docker Hub).
3. **Showcase Deployment:** The premium documentation and showcase website (Phase 9) will be deployed, potentially utilizing a platform-as-a-service (PaaS) or a lightweight VPS.

---
> *Design Note: The deployment pipeline should be as deterministic and clean as the HTTP server code itself.*
