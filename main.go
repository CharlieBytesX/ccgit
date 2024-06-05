// Package main demonstrates how promptkit/selection is used.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

type CommitType struct {
	Type    string
	Meaning string
}

func (t CommitType) String() string {
	return fmt.Sprintf("%s: %s", t.Type, t.Meaning)
}

type Commit struct {
	Type                       string
	Scope                      string
	BreakingChanges            bool
	BreakingChangesDescription string
	Title                      string
	Description                string
}

func (c Commit) String() string {
	myString := c.Type
	if c.Scope != "" {
		myString += fmt.Sprintf("(%s)", c.Scope)
	}
	myString += ": " + c.Title
	if c.Description != "" {
		myString += "\n\n" + c.Description
	}

	if c.BreakingChanges {
		myString += "\n\nBREAKING CHANGE: " + c.BreakingChangesDescription
	}

	return myString

}

func main() {
	conventionalCommits := []CommitType{
		{"feat", "A new feature"},
		{"fix", "A bug fix"},
		{"docs", "Documentation only changes"},
		{"style", "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)"},
		{"refactor", "A code change that neither fixes a bug nor adds a feature"},
		{"perf", "A code change that improves performance"},
		{"test", "Adding missing tests or correcting existing tests"},
		{"build", "Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
		{"ci", "Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)"},
		{"chore", "Other changes that don't modify src or test files"},
		{"revert", "Reverts a previous commit"},
	}
	commitTypesOptions := []string{}
	for _, option := range conventionalCommits {
		commitTypesOptions = append(commitTypesOptions, option.String())
	}

	commit := Commit{}

	commitTypePrompt := selection.New("Select commit type:", commitTypesOptions)

	selectedCommitTypeFull, err := commitTypePrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}
	upToIndex := strings.Index(selectedCommitTypeFull, ":")
	selectedCommitType := selectedCommitTypeFull[:upToIndex]

	commit.Type = selectedCommitType

	titlePrompt := textinput.New("Title:")
	titlePrompt.Placeholder = "Insert commit title"
	titlePrompt.Validate = func(s string) error {
		length := len(s)
		if length > 50 {
			return errors.New("Bigger than 50")
		}
		if length > 0 && s[length-1] == '.' {
			return errors.New("Can't end with .")
		}
		if length < 1 {
			return errors.New("Can't be	0 length")
		}
		return nil
	}

	commitTitle, err := titlePrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// do something with the result
	commit.Title = commitTitle
	// do something with the final choice
	_ = selectedCommitType

	commitBodyPrompt := textinput.New("(optional) Body:")
	commitBodyPrompt.Placeholder = "Insert commit Body"
	commitBodyPrompt.Validate = func(s string) error {
		return nil
	}

	commitBody, err := commitBodyPrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	commit.Description = commitBody

	scopePrompt := textinput.New("(optional) Scope:")
	scopePrompt.Placeholder = "Insert scope"
	scopePrompt.Validate = func(s string) error {
		return nil
	}

	scope, err := scopePrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	commit.Scope = scope

	breakingChangesPrompt := confirmation.New("The commit has breaking changes?", confirmation.No)
	breakingChangesPrompt.Template = confirmation.TemplateYN
	commitHasBreakingChanges, err := breakingChangesPrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	commit.BreakingChanges = commitHasBreakingChanges

	if commitHasBreakingChanges {
		breakingChangesDescriptionPrompt := textinput.New("Breaking Changes Description:")
		bkChangesDescription, err := breakingChangesDescriptionPrompt.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		commit.BreakingChangesDescription = bkChangesDescription
	}

	fmt.Println("")
	cmd := exec.Command("git", "commit", "-m", commit.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing git commit: %v\nOutput: %s\n", err, output)
		os.Exit(1)
	}

	// Imprime la salida del comando git commit
	fmt.Printf("Output: %s\n", output)

}
