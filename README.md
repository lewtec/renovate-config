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

Use the `lewtec-renovate-config` CLI tool to automatically add this preset to any GitHub repository.

### Installation

You can install the CLI using `go install`:

```bash
go install github.com/lewtec/renovate-config/cmd/lewtec-renovate-config@latest
```

### Prerequisites

- Go 1.23+
- Git
- GitHub CLI (`gh`)

### Usage

The CLI can be used to update a single repository or all repositories in an organization.

#### Single Repository

```bash
lewtec-renovate-config <owner> <repo>
```

**Example:**

```bash
lewtec-renovate-config myorg myrepo
```

This will:

1. Clone the target repository.
2. Find the `renovate.json` file.
3. Add the preset to the `extends` array.
4. Create a new branch (`chore/add-renovate-config-preset`).
5. Commit and push the changes.
6. Create a Pull Request (if `gh` CLI is available).

#### Bulk Update

To update all repositories for an owner (organization), omit the `repo` argument.

```bash
lewtec-renovate-config <owner>
```

### Options

- `--preset <ref>`: Custom preset reference (default: `github>lewtec/renovate-config:base`).
- `--no-pr`: Skip PR creation, only push the branch.
- `--include-forks`: Include forked repositories during bulk updates.

### Development

```bash
# Clone the repository
git clone https://github.com/lewtec/renovate-config.git
cd renovate-config

# Build the binary
mise run build

# Run tests
mise run test

# Run directly
go run ./cmd/lewtec-renovate-config <owner> <repo>
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
