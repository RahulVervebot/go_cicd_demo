# Contributing

## Workflow
1. Create a branch from `main`
   - `feat/<short-name>` or `fix/<short-name>`
2. Make small commits with clear messages
3. Open a PR into `main`
4. PR must pass CI and get at least 1 review

## Guidelines to avoid conflicts
- Keep PRs small and focused
- Don’t reformat unrelated files
- Prefer adding new files over editing the same file heavily
- Rebase frequently if your branch lives > 1 day

## Local checks before pushing
```bash
make ci
