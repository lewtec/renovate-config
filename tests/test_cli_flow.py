"""
Tests for the CLI flow logic (argument parsing, repository iteration, etc.)
"""
import argparse
import sys
from unittest.mock import MagicMock, patch

import pytest
from renovate_config_cli import cli

def test_get_repositories_success():
    """Test get_repositories parses gh output correctly and excludes forks by default."""
    mock_run = MagicMock()
    # Default behavior: exclude forks
    mock_run.stdout = '[{"name": "repo1", "isFork": false}, {"name": "repo2", "isFork": true}, {"name": "repo3", "isFork": false}]'
    mock_run.returncode = 0

    with patch('renovate_config_cli.cli.run_command', return_value=mock_run) as mock_run_cmd:
        repos = cli.get_repositories("owner")

        assert repos == ["repo1", "repo3"]
        mock_run_cmd.assert_called_with(['gh', 'repo', 'list', 'owner', '--json', 'name,isFork', '--limit', '1000', '--no-archived'])


def test_get_repositories_include_forks():
    """Test get_repositories includes forks when requested."""
    mock_run = MagicMock()
    mock_run.stdout = '[{"name": "repo1", "isFork": false}, {"name": "repo2", "isFork": true}]'
    mock_run.returncode = 0

    with patch('renovate_config_cli.cli.run_command', return_value=mock_run) as mock_run_cmd:
        repos = cli.get_repositories("owner", include_forks=True)

        assert repos == ["repo1", "repo2"]
        mock_run_cmd.assert_called_with(['gh', 'repo', 'list', 'owner', '--json', 'name,isFork', '--limit', '1000', '--no-archived'])

def test_get_repositories_failure():
    """Test get_repositories handles failure gracefully."""
    with patch('renovate_config_cli.cli.run_command', side_effect=cli.subprocess.CalledProcessError(1, 'cmd')):
        repos = cli.get_repositories("owner")
        assert repos == []

def test_get_repositories_json_error():
    """Test get_repositories handles invalid JSON."""
    mock_run = MagicMock()
    mock_run.stdout = 'invalid json'

    with patch('renovate_config_cli.cli.run_command', return_value=mock_run):
        repos = cli.get_repositories("owner")
        assert repos == []

def test_process_repository_success():
    """Test process_repository flow on success."""
    with patch('renovate_config_cli.cli.clone_repository', return_value=True), \
         patch('renovate_config_cli.cli.find_renovate_config', return_value=cli.Path("renovate.json")), \
         patch('renovate_config_cli.cli.get_default_branch', return_value="main"), \
         patch('renovate_config_cli.cli.run_command'), \
         patch('renovate_config_cli.cli.add_preset_to_config', return_value=True), \
         patch('renovate_config_cli.cli.create_pr_with_gh', return_value=True):

        result = cli.process_repository("owner", "repo", "preset", False, True, "")
        assert result == 0

def test_process_repository_failure_clone():
    """Test process_repository fails if clone fails."""
    with patch('renovate_config_cli.cli.clone_repository', return_value=False):
        result = cli.process_repository("owner", "repo", "preset", False, True, "")
        assert result == 1

def test_process_repository_no_config():
    """Test process_repository fails if no config found."""
    with patch('renovate_config_cli.cli.clone_repository', return_value=True), \
         patch('renovate_config_cli.cli.find_renovate_config', return_value=None):
        result = cli.process_repository("owner", "repo", "preset", False, True, "")
        assert result == 1

def test_main_single_repo():
    """Test main with specific repo argument."""
    with patch('argparse.ArgumentParser.parse_args',
               return_value=argparse.Namespace(owner="owner", repo="repo", preset="preset", no_pr=False, include_forks=False)), \
         patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")), \
         patch('renovate_config_cli.cli.process_repository', return_value=0) as mock_process:

        cli.main()

        mock_process.assert_called_once_with("owner", "repo", "preset", False, True, "")

def test_main_all_repos():
    """Test main without repo argument (process all)."""
    with patch('argparse.ArgumentParser.parse_args',
               return_value=argparse.Namespace(owner="owner", repo=None, preset="preset", no_pr=False, include_forks=False)), \
         patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")), \
         patch('renovate_config_cli.cli.get_repositories', return_value=["repo1", "repo2"]) as mock_get_repos, \
         patch('renovate_config_cli.cli.process_repository', return_value=0) as mock_process:

        cli.main()

        mock_get_repos.assert_called_once_with("owner", False)
        assert mock_process.call_count == 2
        mock_process.assert_any_call("owner", "repo1", "preset", False, True, "")
        mock_process.assert_any_call("owner", "repo2", "preset", False, True, "")

def test_main_all_repos_include_forks():
    """Test main passes include_forks flag."""
    with patch('argparse.ArgumentParser.parse_args',
               return_value=argparse.Namespace(owner="owner", repo=None, preset="preset", no_pr=False, include_forks=True)), \
         patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")), \
         patch('renovate_config_cli.cli.get_repositories', return_value=["repo1"]) as mock_get_repos, \
         patch('renovate_config_cli.cli.process_repository', return_value=0) as mock_process:

        cli.main()

        mock_get_repos.assert_called_once_with("owner", True)

def test_main_all_repos_mixed_success():
    """Test main reporting failure count."""
    with patch('argparse.ArgumentParser.parse_args',
               return_value=argparse.Namespace(owner="owner", repo=None, preset="preset", no_pr=False, include_forks=False)), \
         patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")), \
         patch('renovate_config_cli.cli.get_repositories', return_value=["repo1", "repo2"]), \
         patch('renovate_config_cli.cli.process_repository', side_effect=[0, 1]):

        ret = cli.main()
        assert ret == 1  # Should return 1 because one repo failed

def test_main_all_repos_gh_unavailable():
    """Test main fails gracefully if gh is not available and repo is not specified."""
    with patch('argparse.ArgumentParser.parse_args',
               return_value=argparse.Namespace(owner="owner", repo=None, preset="preset", no_pr=False, include_forks=False)), \
         patch('renovate_config_cli.cli.check_gh_cli', return_value=(False, "error message")):

        ret = cli.main()
        assert ret == 1
