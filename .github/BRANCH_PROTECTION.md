# Branch protection and GitHub setup

Recommended settings to protect commits and prevent regression.

## Branch protection rules (Settings → Branches)

Add a rule for **main** (or your default branch):

1. **Require a pull request before merging**
   - Require at least 1 approval (optional but recommended).
   - Dismiss stale reviews when new commits are pushed (optional).

2. **Require status checks to pass before merging**
   - Require the **CI** workflow to pass:
     - Add status check: `test` (the CI job name).
   - This blocks merge when tests fail and protects against regression.

3. **Do not allow bypassing the above** (no “Allow specified actors to bypass” unless needed).

4. **Restrict who can push to matching branches** (optional): allow only maintainers to push to `main`; everyone else uses PRs.

5. **Require signed commits** (optional): enable “Require signed commits” if you want commit signing.

## Workflows in this repo

| Workflow | Trigger | Purpose |
|----------|--------|---------|
| **CI** | Push/PR to `main` | Build, test, vet. Must pass before merge. |
| **OpenAPI spec check** | Daily (cron) + manual | Fetches live Loops spec; if it differs from `openapi.json`, fails and opens a GitHub issue to notify maintainers. |

After enabling branch protection with the CI status check, every PR must have a green CI run before it can be merged.
