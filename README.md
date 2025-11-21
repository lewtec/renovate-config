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

Use the `renovate-config` CLI tool to automatically add this preset to any GitHub repository.

### Installation

#### Using UV (recommended)

```bash
# Install UV if you don't have it
curl -LsSf https://astral.sh/uv/install.sh | sh

# Install the CLI tool
uv tool install git+https://github.com/lewtec/renovate-config.git

# Or for development
uv pip install -e .
```

#### Using pip

```bash
pip install git+https://github.com/lewtec/renovate-config.git
```

### Prerequisites

- Python 3.8+
- Git
- GitHub CLI (`gh`) - optional, for creating PRs automatically

### Usage

```bash
renovate-config <owner> <repo>
```

**Example:**
```bash
renovate-config myorg myrepo
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
renovate-config myorg myrepo --preset "github>lewtec/renovate-config:custom"
```

### Development

```bash
# Clone the repository
git clone https://github.com/lewtec/renovate-config.git
cd renovate-config

# Install with UV
uv pip install -e ".[dev]"

# Run tests
uv run pytest

# Or run directly
python -m renovate_config_cli.cli <owner> <repo>
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