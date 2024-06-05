package git

import (
	"fmt"
	"os/exec"
)

func IsInsideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		fmt.Println("No Git repository found.")
		return false
	}

	return true
}
