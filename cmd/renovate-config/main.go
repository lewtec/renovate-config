package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	preset       string
	noPr         bool
	includeForks bool
)

const (
	defaultPreset = "github>lewtec/renovate-config:base"
	branchName    = "chore/add-renovate-config-preset"
	commitMessage = "chore: add renovate-config preset"
	prTitle       = "Add renovate-config preset"
	prBody        = `This PR adds the renovate-config preset to the repository's Renovate configuration.

The preset includes best practices and recommended configurations for dependency updates.

Changes:
- Added ` + "`extends`" + ` configuration to include ` + "`github>lewtec/renovate-config:base`"
)

func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func findRenovateConfig(repoPath string) string {
	possibleLocations := []string{
		"renovate.json",
		".github/renovate.json",
		".gitlab/renovate.json",
		".renovaterc.json",
		".renovaterc",
	}

	for _, loc := range possibleLocations {
		path := filepath.Join(repoPath, loc)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Found renovate config at: %s\n", path)
			return path
		}
	}
	return ""
}

func detectIndentation(content string) string {
	re := regexp.MustCompile(`\n(\s+)`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		indent := match[1]
		if strings.Contains(indent, "\t") {
			return "\t"
		}
		return indent
	}
	return "  "
}

func addPresetToConfig(configPath, presetRef string) (bool, error) {
	contentBytes, err := os.ReadFile(configPath)
	if err != nil {
		return false, err
	}
	content := string(contentBytes)

	indent := detectIndentation(content)

	var config map[string]interface{}
	if err := json.Unmarshal(contentBytes, &config); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return false, err
	}

	var extends []interface{}
	if val, ok := config["extends"]; ok {
		if list, ok := val.([]interface{}); ok {
			extends = list
		} else if str, ok := val.(string); ok {
			extends = []interface{}{str}
		}
	} else {
		extends = []interface{}{}
	}

	// Check if preset already exists
	foundIndex := -1
	for i, v := range extends {
		if s, ok := v.(string); ok && s == presetRef {
			foundIndex = i
			break
		}
	}

	if foundIndex != -1 {
		if foundIndex == len(extends)-1 {
			fmt.Printf("Preset %s already exists at the end of extends\n", presetRef)
			return false, nil
		}
		// Remove it to append at end
		extends = append(extends[:foundIndex], extends[foundIndex+1:]...)
		fmt.Printf("Moving %s to the end of extends array\n", presetRef)
	}

	extends = append(extends, presetRef)
	config["extends"] = extends

	// We need to marshal with indentation
	// Go's standard library doesn't easily support preserving exact indentation of the whole file if it's mixed,
	// but we will use the detected indent.
	// Use encoder to prevent HTML escaping (e.g., > becoming \u003e)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	if err := encoder.Encode(config); err != nil {
		return false, err
	}

	outputBytes := buf.Bytes()

	if err := os.WriteFile(configPath, outputBytes, 0644); err != nil {
		return false, err
	}

	fmt.Printf("Added %s to extends array\n", presetRef)
	return true, nil
}

func getDefaultBranch(repoPath string) string {
	output, err := runCommand(repoPath, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "main"
	}
	branch := strings.TrimSpace(output)
	branch = strings.Replace(branch, "refs/remotes/origin/", "", 1)
	if branch == "" {
		return "main"
	}
	return branch
}

func checkGhCli() (bool, string) {
	_, err := runCommand("", "gh", "--version")
	if err != nil {
		return false, "gh CLI not found. Please install it: https://cli.github.com/"
	}

	_, err = runCommand("", "gh", "auth", "status")
	if err != nil {
		return false, "gh CLI is not authenticated. Please run: gh auth login"
	}
	return true, ""
}

func cloneRepository(owner, repo, destination string) bool {
	_, err := runCommand("", "gh", "repo", "clone", fmt.Sprintf("%s/%s", owner, repo), destination)
	if err != nil {
		fmt.Printf("Error cloning with gh: %v\nFallback to git clone...\n", err)
		_, err = runCommand("", "git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", owner, repo), destination)
		if err != nil {
			fmt.Printf("Error cloning with git: %v\n", err)
			return false
		}
	}
	return true
}

func createPrWithGh(repoPath, owner, repo, baseBranch string) bool {
	output, err := runCommand(repoPath, "gh", "pr", "create",
		"--repo", fmt.Sprintf("%s/%s", owner, repo),
		"--title", prTitle,
		"--body", prBody,
		"--base", baseBranch,
		"--head", branchName,
	)
	if err != nil {
		fmt.Printf("Error creating PR: %v\n%s\n", err, output)
		fmt.Printf("\nBranch '%s' has been pushed. Create PR manually.\n", branchName)
		return false
	}
	fmt.Printf("\n✓ Pull request created successfully!\n%s\n", output)
	return true
}

type ghRepo struct {
	Name   string `json:"name"`
	IsFork bool   `json:"isFork"`
}

func getRepositories(owner string, includeForks bool) ([]string, error) {
	fmt.Printf("Fetching repositories for %s...\n", owner)
	output, err := runCommand("", "gh", "repo", "list", owner, "--json", "name,isFork", "--limit", "1000", "--no-archived")
	if err != nil {
		return nil, err
	}

	var repos []ghRepo
	if err := json.Unmarshal([]byte(output), &repos); err != nil {
		return nil, err
	}

	var names []string
	for _, r := range repos {
		if !includeForks && r.IsFork {
			continue
		}
		names = append(names, r.Name)
	}
	return names, nil
}

func processRepository(owner, repo, presetRef string, noPr bool, ghAvailable bool, ghError string) int {
	fmt.Printf("\n%s\nProcessing repository: %s/%s\n%s\n\n", strings.Repeat("=", 50), owner, repo, strings.Repeat("=", 50))

	tmpDir, err := os.MkdirTemp("", "renovate-cli-*")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	repoPath := filepath.Join(tmpDir, repo)

	fmt.Printf("Cloning %s/%s...\n", owner, repo)
	if !cloneRepository(owner, repo, repoPath) {
		return 1
	}

	configPath := findRenovateConfig(repoPath)
	if configPath == "" {
		fmt.Printf("\nError: No renovate configuration file found in %s/%s\n", owner, repo)
		return 1
	}

	defaultBranch := getDefaultBranch(repoPath)
	fmt.Printf("Default branch: %s\n", defaultBranch)

	fmt.Printf("\nCreating branch %s...\n", branchName)
	if _, err := runCommand(repoPath, "git", "checkout", "-b", branchName); err != nil {
		fmt.Printf("Error creating branch: %v\n", err)
		return 1
	}

	fmt.Println("\nUpdating renovate configuration...")
	changesMade, err := addPresetToConfig(configPath, presetRef)
	if err != nil {
		fmt.Printf("Error updating config: %v\n", err)
		return 1
	}

	if !changesMade {
		fmt.Println("\nNo changes needed.")
		return 0
	}

	fmt.Println("\nChanges made:")
	out, _ := runCommand(repoPath, "git", "diff", filepath.Base(configPath))
	fmt.Println(out)

	fmt.Println("\nCommitting changes...")
	runCommand(repoPath, "git", "add", configPath)
	runCommand(repoPath, "git", "commit", "-m", commitMessage)

	// Delete remote branch if it exists to avoid conflicts
	fmt.Println("\nChecking if remote branch exists...")
	runCommand(repoPath, "git", "push", "origin", "--delete", branchName)

	fmt.Printf("\nPushing branch %s...\n", branchName)
	if output, err := runCommand(repoPath, "git", "push", "-u", "origin", branchName); err != nil {
		fmt.Printf("Error pushing branch: %v\n%s\n", err, output)
		return 1
	}

	if !noPr {
		if ghAvailable {
			fmt.Println("\nCreating pull request...")
			createPrWithGh(repoPath, owner, repo, defaultBranch)
		} else {
			fmt.Printf("\nBranch %s pushed. Cannot create PR: %s\n", branchName, ghError)
		}
	} else {
		fmt.Printf("\nBranch %s pushed successfully.\n", branchName)
	}

	fmt.Println("\n✓ Done with repository!")
	return 0
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "renovate-config [owner] [repo]",
		Short: "Add renovate-config preset to a GitHub repository",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			owner := args[0]

			ghAvailable, ghError := checkGhCli()
			if !ghAvailable {
				fmt.Printf("⚠ Warning: %s\n", ghError)
				if !noPr {
					fmt.Println("  PR creation will not be available. Use --no-pr to suppress this warning.")
					fmt.Println()
				}
			}

			if len(args) == 2 {
				repo := args[1]
				os.Exit(processRepository(owner, repo, preset, noPr, ghAvailable, ghError))
			} else {
				if !ghAvailable {
					fmt.Printf("Error: %s\n", ghError)
					os.Exit(1)
				}
				repos, err := getRepositories(owner, includeForks)
				if err != nil {
					fmt.Printf("Error fetching repos: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Found %d repositories: %s\n", len(repos), strings.Join(repos, ", "))

				failures := 0
				for _, repo := range repos {
					if res := processRepository(owner, repo, preset, noPr, ghAvailable, ghError); res != 0 {
						failures++
					}
				}
				if failures > 0 {
					fmt.Printf("\nFinished with %d failures out of %d repositories.\n", failures, len(repos))
					os.Exit(1)
				}
				fmt.Printf("\nSuccessfully processed all %d repositories.\n", len(repos))
			}
		},
	}

	rootCmd.Flags().StringVar(&preset, "preset", defaultPreset, "Preset reference")
	rootCmd.Flags().BoolVar(&noPr, "no-pr", false, "Do not create a PR, just push the branch")
	rootCmd.Flags().BoolVar(&includeForks, "include-forks", false, "Include forked repositories when processing all repos")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
