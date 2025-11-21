#!/usr/bin/env python3
"""
Script to add renovate-config preset to a GitHub repository's renovate.json

Usage:
    python add_renovate_preset.py <owner> <repo>

Example:
    python add_renovate_preset.py myorg myrepo

This will:
1. Clone the repository
2. Find and update renovate.json to include the preset
3. Create a PR with the changes
"""

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Optional, Tuple


PRESET_REF = "github>lewtec/renovate-config:base"
BRANCH_NAME = "chore/add-renovate-config-preset"
COMMIT_MESSAGE = "chore: add renovate-config preset"
PR_TITLE = "Add renovate-config preset"
PR_BODY = """This PR adds the renovate-config preset to the repository's Renovate configuration.

The preset includes best practices and recommended configurations for dependency updates.

Changes:
- Added `extends` configuration to include `github>lewtec/renovate-config:base`
"""


def run_command(cmd: list[str], cwd: Optional[str] = None, check: bool = True) -> subprocess.CompletedProcess:
    """Run a shell command and return the result."""
    print(f"Running: {' '.join(cmd)}")
    return subprocess.run(cmd, cwd=cwd, check=check, capture_output=True, text=True)


def find_renovate_config(repo_path: Path) -> Optional[Path]:
    """Find renovate configuration file in the repository."""
    possible_locations = [
        "renovate.json",
        ".github/renovate.json",
        ".gitlab/renovate.json",
        ".renovaterc.json",
        ".renovaterc",
    ]

    for location in possible_locations:
        config_path = repo_path / location
        if config_path.exists():
            print(f"Found renovate config at: {config_path}")
            return config_path

    return None


def detect_indentation(content: str) -> str:
    """Detect the indentation used in the JSON file."""
    # Try to find indentation from the first indented line
    match = re.search(r'\n(\s+)', content)
    if match:
        indent = match.group(1)
        # Check if it's spaces or tabs
        if '\t' in indent:
            return '\t'
        else:
            return ' ' * len(indent)
    # Default to 2 spaces
    return '  '


def add_preset_to_config(config_path: Path, preset_ref: str) -> bool:
    """
    Add the preset to renovate config while preserving formatting.
    Returns True if changes were made, False otherwise.
    """
    with open(config_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Detect indentation
    indent = detect_indentation(content)

    # Parse JSON
    try:
        config = json.loads(content)
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON: {e}")
        return False

    # Check if extends already exists
    if 'extends' not in config:
        config['extends'] = []

    # Ensure extends is a list
    if not isinstance(config['extends'], list):
        config['extends'] = [config['extends']]

    # Check if preset already exists
    if preset_ref in config['extends']:
        print(f"Preset {preset_ref} already exists in extends")
        return False

    # Add preset at the beginning of extends array
    config['extends'].insert(0, preset_ref)

    # Write back with same indentation
    indent_str = indent if indent != '\t' else '\t'
    indent_level = 1 if indent == '\t' else len(indent)

    with open(config_path, 'w', encoding='utf-8') as f:
        json.dump(config, f, indent=indent_level if indent != '\t' else '\t', ensure_ascii=False)
        f.write('\n')  # Add trailing newline

    print(f"Added {preset_ref} to extends array")
    return True


def get_default_branch(repo_path: Path) -> str:
    """Get the default branch name."""
    result = run_command(['git', 'symbolic-ref', 'refs/remotes/origin/HEAD'], cwd=str(repo_path))
    branch = result.stdout.strip().replace('refs/remotes/origin/', '')
    return branch if branch else 'main'


def check_gh_cli() -> Tuple[bool, str]:
    """
    Check if gh CLI is installed and authenticated.
    Returns (is_available, error_message)
    """
    # Check if gh is in PATH
    try:
        result = run_command(['which', 'gh'], check=False)
        if result.returncode != 0:
            return False, "gh CLI not found in PATH. Please install it: https://cli.github.com/"
    except FileNotFoundError:
        return False, "gh CLI not found. Please install it: https://cli.github.com/"

    # Check if gh is authenticated
    try:
        result = run_command(['gh', 'auth', 'status'], check=False)
        if result.returncode != 0:
            return False, "gh CLI is not authenticated. Please run: gh auth login"
    except subprocess.CalledProcessError:
        return False, "gh CLI is not authenticated. Please run: gh auth login"

    return True, ""


def create_pr_with_gh(repo_path: Path, owner: str, repo: str, base_branch: str) -> bool:
    """Create a pull request using gh CLI."""
    # Verify gh CLI is available and authenticated
    is_available, error_msg = check_gh_cli()
    if not is_available:
        print(f"Error: {error_msg}")
        print(f"Branch '{BRANCH_NAME}' has been pushed to the repository.")
        print(f"You can create a PR manually at: https://github.com/{owner}/{repo}/compare/{base_branch}...{BRANCH_NAME}")
        return False

    try:
        # Create PR
        result = run_command([
            'gh', 'pr', 'create',
            '--title', PR_TITLE,
            '--body', PR_BODY,
            '--base', base_branch,
            '--head', BRANCH_NAME,
        ], cwd=str(repo_path))

        print(f"\n{result.stdout}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error creating PR with gh CLI: {e}")
        print(f"stderr: {e.stderr}")
        print(f"\nBranch '{BRANCH_NAME}' has been pushed to the repository.")
        print(f"You can create a PR manually at: https://github.com/{owner}/{repo}/compare/{base_branch}...{BRANCH_NAME}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description='Add renovate-config preset to a GitHub repository'
    )
    parser.add_argument('owner', help='GitHub repository owner')
    parser.add_argument('repo', help='GitHub repository name')
    parser.add_argument('--preset', default=PRESET_REF, help=f'Preset reference (default: {PRESET_REF})')
    parser.add_argument('--no-pr', action='store_true', help='Do not create a PR, just push the branch')

    args = parser.parse_args()

    # Create temporary directory
    with tempfile.TemporaryDirectory() as tmpdir:
        repo_path = Path(tmpdir) / args.repo

        # Clone repository
        print(f"\nCloning {args.owner}/{args.repo}...")
        try:
            run_command([
                'git', 'clone',
                f'https://github.com/{args.owner}/{args.repo}.git',
                str(repo_path)
            ])
        except subprocess.CalledProcessError as e:
            print(f"Error cloning repository: {e}")
            print(f"stderr: {e.stderr}")
            return 1

        # Find renovate config
        config_path = find_renovate_config(repo_path)
        if not config_path:
            print("\nError: No renovate configuration file found in the repository")
            print("Looked for: renovate.json, .github/renovate.json, .gitlab/renovate.json, .renovaterc.json, .renovaterc")
            return 1

        # Get default branch
        default_branch = get_default_branch(repo_path)
        print(f"Default branch: {default_branch}")

        # Create new branch
        print(f"\nCreating branch {BRANCH_NAME}...")
        run_command(['git', 'checkout', '-b', BRANCH_NAME], cwd=str(repo_path))

        # Add preset to config
        print("\nUpdating renovate configuration...")
        changes_made = add_preset_to_config(config_path, args.preset)

        if not changes_made:
            print("\nNo changes needed. Preset already configured or no changes made.")
            return 0

        # Show diff
        print("\nChanges made:")
        diff_result = run_command(['git', 'diff', str(config_path.name)], cwd=str(repo_path), check=False)
        print(diff_result.stdout)

        # Commit changes
        print("\nCommitting changes...")
        run_command(['git', 'add', str(config_path)], cwd=str(repo_path))
        run_command(['git', 'commit', '-m', COMMIT_MESSAGE], cwd=str(repo_path))

        # Push branch
        print(f"\nPushing branch {BRANCH_NAME}...")
        try:
            run_command(['git', 'push', '-u', 'origin', BRANCH_NAME], cwd=str(repo_path))
        except subprocess.CalledProcessError as e:
            print(f"Error pushing branch: {e}")
            print(f"stderr: {e.stderr}")
            return 1

        # Create PR
        if not args.no_pr:
            print("\nCreating pull request...")
            create_pr_with_gh(repo_path, args.owner, args.repo, default_branch)
        else:
            print(f"\nBranch {BRANCH_NAME} pushed successfully.")
            print(f"Create a PR at: https://github.com/{args.owner}/{args.repo}/compare/{default_branch}...{BRANCH_NAME}")

        print("\n✓ Done!")
        return 0


if __name__ == '__main__':
    sys.exit(main())
