# AI Agent Instructions for Clubrizer

## About
Clubrizer is a platform for sports clubs to manage their club life — initially focused on events, with plans to expand to features like member management, records, protocols, and a forum.

## Repository
GitHub: `KatharinaSick/clubrizer`

## General Guidelines
- **Approval**: Always ask for approval before implementing any changes.
- **Task Breakdown**: Always break tasks into small, manageable chunks.
- **Clarification**: Never guess. If anything is unclear, ask for clarification until the task is fully understood.
- **Implementation**: Do not implement new features or significant changes without explicit approval. Planning should precede implementation.
- **Users**: The end users of this app are not tech-savvy. Keep UX simple, friendly, and avoid technical language in any user-facing text.
- **Full-stack security**: Always think about the full picture — this app has both a frontend and a backend. Hiding or restricting something only on the frontend is not sufficient security. Any access control must be properly enforced in the backend as well.

## Git Commits
- **Style**: Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat:`, `fix:`, `chore:`, `docs:`).
- **Closing issues**: Add `Closes #<issue-number>` in the commit body when a commit resolves a GitHub issue.

## GitHub Issues
- **Labels**: Every issue must have at least one label. Never create an issue without labels. Always suggest labels and ask for confirmation before applying them — do not assume.
- **Milestones**: Not every issue needs a milestone, but always suggest one if applicable and ask for confirmation before assigning — do not assume.

## Frontend Tech Stack
- **Framework**: Vue 3 (Composition API, `<script setup lang="ts">`), TypeScript, Vite.
- **Routing**: Vue Router.
- **HTTP Client**: Axios.
- **Internationalization**: vue-i18n.

## Backend Tech Stack
- **Language**: Go (Golang) 1.26.3.
- **Database Driver**: pgx/v5 (PostgreSQL).
- **Validation**: go-playground/validator.
- **Authentication**: golang-jwt/jwt/v5.
- **Configuration**: yaml.v3.

## Coding Conventions

### Frontend
- **TypeScript**: Always use TypeScript. Avoid `any` where possible.
- **Components**: Use Single File Components (SFC) with `<script setup>`.
- **Design System**: Always prefer reusing existing components from `frontend/src/components` over native HTML elements. For example, use `<Button>` instead of `<button>`, and extend existing components with new props/themes when new variants are needed.
- **I18n**: All user-facing text must be internationalized using `i18n.global.t()`. Do not hardcode strings in templates or scripts.
  - Add new keys to `frontend/src/plugins/i18n.ts`.
  - Structure keys hierarchically (e.g., `events.new.title`).
- **Styling**: Use scoped CSS (`<style scoped>`).
  - **CSS Variables**: Always use variables from `frontend/src/assets/base.css`. Never hardcode values for colors, font sizes, font weights, spacing, border radius, shadows, or any other design token — always use the corresponding CSS variable.
  - **Class Naming**: Use camelCase for CSS class names and always prefix them with the component name (e.g., `.myComponentHeader`).
- **Loading & Error Handling**: Do not add per-component loading or error state for API requests. Loading and errors are handled globally:
  - **Loading**: A counter-based Pinia store (`requestStore`) is incremented/decremented automatically by axios interceptors (`frontend/src/service/axiosInterceptors.ts`). A slim top progress bar in `App.vue` reflects this globally. Do not add local loading indicators except inside buttons (button spinners are fine and encouraged to prevent double-submits).
  - **Errors**: All non-422, non-401 errors are caught by the axios interceptor and stored in `requestStore.error`. Views place a `<RequestError />` component wherever they want the error to appear — this component reads from the store automatically. Do not add local error `<Alert>` components for API errors.
  - **Exceptions**: 422 (form validation) errors should be handled locally by the component. Client-side errors that don't involve an HTTP request (e.g. missing Google credential) should also be handled locally.

### Backend
- **Structure**: Follow standard Go project layout.
  - `cmd/`: Application entry points.
  - `internal/`: Private application and library code.
- **Error Handling**: Use custom error types where appropriate.

## Project Structure
- `frontend/src/components`: Reusable UI components.
- `frontend/src/views`: Page views mapped to routes.
- `frontend/src/plugins`: Plugin configurations (i18n, axios, etc.).
- `backend/cmd`: Main applications.
- `backend/internal`: Internal packages (api, app, users, events, apperrors).
