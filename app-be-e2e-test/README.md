This is an end-to-end test for backend only.
This means there'll be no UI automation.
There'll only be HTTP requests to hit backend API.

Benefits:

- Testing backend functionality.
- Make sure an exceptional bug happens only once.
  1. An exception is reported by QA or end user.
  2. We debug the root cause, fix it, and push the hotfix to production.
  3. We create an end-to-end test that covers that exceptional/edge case.
  4. We'll never encounter the same bug again.
- Documenting API flow (what APIs need to be hit to achieve a user story).
- Manually hitting an API in dev/staging/prod for operation purposes.

The `libSingle` directory is filled with functions and structures that hit one API endpoint only.
Useful for local testing when developing one API.

The `libMultiple` directory is filled with functions and structures that hit multiple API endpoints.
Useful for end-to-end testing scenarios that simulate business flow of user stories.

## E2E Backend Run

- `cd e2e-backend`
- `go mod tidy`
- Create an `.env` file from `.env.sample`
- Open `executeSingle/main.go`
  - Update the import to the libSingle you want to execute
- Execute: `go run executeSingle/*.go`
  - In Windows: `go run .\executeSingle\`

## E2E Backend Init

- go mod init diskon-hunter/price-monitoring-e2e-test