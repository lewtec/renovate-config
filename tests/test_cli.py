"""
Tests for the renovate-config CLI tool
"""

import json
import tempfile
from pathlib import Path

from renovate_config_cli.cli import add_preset_to_config, detect_indentation


def test_add_preset_to_empty_extends():
    """Test adding preset to a config with empty extends array"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        config = {
            "$schema": "https://docs.renovatebot.com/renovate-schema.json",
            "extends": []
        }
        json.dump(config, f, indent=2)
        f.flush()
        temp_path = Path(f.name)

    try:
        result = add_preset_to_config(temp_path, "github>lewtec/renovate-config:base")
        assert result == True, "Should return True when changes are made"

        with open(temp_path, 'r') as f:
            updated_config = json.load(f)

        assert "github>lewtec/renovate-config:base" in updated_config["extends"]
        print("✓ test_add_preset_to_empty_extends passed")
    finally:
        temp_path.unlink()


def test_add_preset_to_existing_extends():
    """Test adding preset to a config with existing extends"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        config = {
            "$schema": "https://docs.renovatebot.com/renovate-schema.json",
            "extends": ["config:recommended"]
        }
        json.dump(config, f, indent=2)
        f.flush()
        temp_path = Path(f.name)

    try:
        result = add_preset_to_config(temp_path, "github>lewtec/renovate-config:base")
        assert result == True, "Should return True when changes are made"

        with open(temp_path, 'r') as f:
            updated_config = json.load(f)

        assert "github>lewtec/renovate-config:base" in updated_config["extends"]
        assert "config:recommended" in updated_config["extends"]
        # Should be inserted at the beginning
        assert updated_config["extends"][0] == "github>lewtec/renovate-config:base"
        print("✓ test_add_preset_to_existing_extends passed")
    finally:
        temp_path.unlink()


def test_preset_already_exists():
    """Test that script detects when preset already exists"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        config = {
            "$schema": "https://docs.renovatebot.com/renovate-schema.json",
            "extends": ["github>lewtec/renovate-config:base", "config:recommended"]
        }
        json.dump(config, f, indent=2)
        f.flush()
        temp_path = Path(f.name)

    try:
        result = add_preset_to_config(temp_path, "github>lewtec/renovate-config:base")
        assert result == False, "Should return False when preset already exists"
        print("✓ test_preset_already_exists passed")
    finally:
        temp_path.unlink()


def test_indentation_detection():
    """Test that indentation detection works"""
    # Test 2 spaces
    content_2_spaces = '{\n  "extends": []\n}'
    indent = detect_indentation(content_2_spaces)
    assert indent == '  ', f"Expected 2 spaces, got {repr(indent)}"

    # Test 4 spaces
    content_4_spaces = '{\n    "extends": []\n}'
    indent = detect_indentation(content_4_spaces)
    assert indent == '    ', f"Expected 4 spaces, got {repr(indent)}"

    # Test tabs
    content_tabs = '{\n\t"extends": []\n}'
    indent = detect_indentation(content_tabs)
    assert indent == '\t', f"Expected tab, got {repr(indent)}"

    print("✓ test_indentation_detection passed")


def test_preserves_formatting():
    """Test that formatting is preserved"""
    original_content = """{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "config:recommended"
  ],
  "automerge": true
}
"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        f.write(original_content)
        f.flush()
        temp_path = Path(f.name)

    try:
        add_preset_to_config(temp_path, "github>lewtec/renovate-config:base")

        with open(temp_path, 'r') as f:
            updated_content = f.read()

        # Check that file still has proper JSON structure
        updated_config = json.loads(updated_content)
        assert updated_config["extends"][0] == "github>lewtec/renovate-config:base"
        assert updated_config["automerge"] == True

        # Check indentation is preserved (2 spaces)
        assert '  "$schema"' in updated_content or '  "extends"' in updated_content

        print("✓ test_preserves_formatting passed")
    finally:
        temp_path.unlink()


if __name__ == '__main__':
    print("Running tests...\n")
    test_indentation_detection()
    test_add_preset_to_empty_extends()
    test_add_preset_to_existing_extends()
    test_preset_already_exists()
    test_preserves_formatting()
    print("\n✓ All tests passed!")
