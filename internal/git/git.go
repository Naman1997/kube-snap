package git

import (
	"os/exec"
)

func CloneRepo(repoUrl string) error {
	cmd := exec.Command("git", "clone", repoUrl, ".")
	return cmd.Run()
}

func PullOrigin() error {
	cmd := exec.Command("git", "pull")
	return cmd.Run()
}

func SetupAuthor(email string, name string) {
	exec.Command("git", "config", "--global", "user.email", email).Run()
	exec.Command("git", "config", "--global", "user.name", name).Run()
}

func SwitchBranch(branch string) error {
	cmd := exec.Command("git", "switch", branch)
	return cmd.Run()
}

func AddAll() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

func Status() (string, error) {
	o, err := exec.Command("git", "status", "--porcelain").Output()
	return string(o), err
}

func CommitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

func Push() error {
	cmd := exec.Command("git", "push")
	return cmd.Run()
}
