"""
Tests for bulk update functionality in renovate-config CLI
"""

import argparse
import json
import subprocess
from unittest.mock import MagicMock, patch

import pytest
from renovate_config_cli.cli import get_repositories, main


@pytest.fixture
def mock_run_command():
    with patch('renovate_config_cli.cli.run_command') as mock:
        yield mock


@pytest.fixture
def mock_process_repository():
    with patch('renovate_config_cli.cli.process_repository') as mock:
        yield mock


def test_get_repositories_excludes_forks_by_default(mock_run_command):
    """Test that get_repositories excludes forks when include_forks is False"""
    mock_run_command.return_value = MagicMock(
        stdout=json.dumps([
            {"name": "repo1", "isFork": False},
            {"name": "repo2", "isFork": True},
            {"name": "repo3", "isFork": False},
        ]),
        returncode=0
    )

    repos = get_repositories("test-owner", include_forks=False)

    assert repos == ["repo1", "repo3"]
    mock_run_command.assert_called_once()
    args = mock_run_command.call_args[0][0]
    assert args[:4] == ['gh', 'repo', 'list', 'test-owner']
    assert '--no-archived' in args


def test_get_repositories_includes_forks_when_requested(mock_run_command):
    """Test that get_repositories includes forks when include_forks is True"""
    mock_run_command.return_value = MagicMock(
        stdout=json.dumps([
            {"name": "repo1", "isFork": False},
            {"name": "repo2", "isFork": True},
        ]),
        returncode=0
    )

    repos = get_repositories("test-owner", include_forks=True)

    assert repos == ["repo1", "repo2"]


def test_main_bulk_mode(mock_run_command, mock_process_repository):
    """Test that main iterates over repositories in bulk mode"""
    # Mock check_gh_cli
    with patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")):
        # Mock get_repositories output
        mock_run_command.return_value = MagicMock(
            stdout=json.dumps([
                {"name": "repo1", "isFork": False},
                {"name": "repo2", "isFork": False},
            ]),
            returncode=0
        )
        mock_process_repository.return_value = True

        # Run main with only owner arg
        with patch('argparse.ArgumentParser.parse_args',
                  return_value=argparse.Namespace(
                      owner='test-owner',
                      repo=None,
                      preset='default-preset',
                      no_pr=False,
                      include_forks=False
                  )):
            main()

        # Check process_repository called for each repo
        assert mock_process_repository.call_count == 2
        # Use ANY for the args parameter since we can't easily match the Namespace object exactly
        from unittest.mock import ANY
        mock_process_repository.assert_any_call('test-owner', 'repo1', ANY, True, "")
        mock_process_repository.assert_any_call('test-owner', 'repo2', ANY, True, "")


def test_main_single_repo_mode(mock_run_command, mock_process_repository):
    """Test that main processes single repo when repo arg is provided"""
    with patch('renovate_config_cli.cli.check_gh_cli', return_value=(True, "")):
        mock_process_repository.return_value = True

        with patch('argparse.ArgumentParser.parse_args',
                  return_value=argparse.Namespace(
                      owner='test-owner',
                      repo='test-repo',
                      preset='default-preset',
                      no_pr=False,
                      include_forks=False
                  )):
            main()

        assert mock_process_repository.call_count == 1
        from unittest.mock import ANY
        mock_process_repository.assert_called_with('test-owner', 'test-repo', ANY, True, "")
