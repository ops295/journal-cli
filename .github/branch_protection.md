Branch protection guidance

This repository aims to keep `main` protected with an open-source friendly policy.

Recommendations (manual or via API):

- Require pull request reviews before merging (1+ approvals).
- Require status checks to pass (CI/build) before merge.
- Require code owner reviews for files covered by `CODEOWNERS` (e.g., critical paths).
- Enforce for administrators if the project wants strict protection.

Example GitHub CLI command to set a minimal protection for `main` (requires `gh` + repo admin):

```bash
gh api --method PUT /repos/:owner/:repo/branches/main/protection -f required_status_checks.contexts='[]' -f enforce_admins=true -f required_pull_request_reviews.dismiss_stale_reviews=true -f required_pull_request_reviews.require_code_owner_reviews=true
```

There is also a GitHub Actions workflow in `.github/workflows/enforce-branch-protection.yml` that can be used (requires a PAT with repo admin rights set as `BG_PROTECT_TOKEN` secret) to set protection automatically.
