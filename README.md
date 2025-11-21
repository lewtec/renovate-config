# renovate-config

Shareable Renovate configuration presets.

## Usage

To use this preset in your repository, add it to your `renovate.json`:

```json
{
  "extends": ["github>lewtec/renovate-config:base"]
}
```

## Automated Setup Script

Use the `add_renovate_preset.py` script to automatically add this preset to any GitHub repository:

### Prerequisites

- Python 3.7+
- Git
- GitHub CLI (`gh`) - optional, for creating PRs automatically

### Usage

```bash
python add_renovate_preset.py <owner> <repo>
```

**Example:**
```bash
python add_renovate_preset.py myorg myrepo
```

This will:
1. Clone the target repository
2. Find the `renovate.json` file
3. Add the preset to the `extends` array (preserving existing formatting)
4. Create a new branch
5. Commit the changes
6. Push the branch
7. Create a Pull Request (if `gh` CLI is available)

### Options

- `--preset` - Custom preset reference (default: `github>lewtec/renovate-config:base`)
- `--no-pr` - Skip PR creation, only push the branch

### Example with custom preset

```bash
python add_renovate_preset.py myorg myrepo --preset "github>lewtec/renovate-config:custom"
```

## Presets

### base

The base preset includes:
- Recommended configuration
- Dependency dashboard
- Vulnerability alerts
- Docker support with digest pinning
- GitHub Actions digest pinning
- Lock file maintenance
- Automerge for dev dependencies (minor updates)
- Custom package rules for Go, Terraform, and Cloudflare packages