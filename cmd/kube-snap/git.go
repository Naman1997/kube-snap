package main

import (
	"os/exec"
	"time"
)

func cloneRepo() error {
	cmd := exec.Command("git", "clone", getValueOf("repo-url"), ".")
	return cmd.Run()
}

func pullOrigin() error {
	cmd := exec.Command("git", "pull")
	return cmd.Run()
}

func setupAuthor() {
	exec.Command("git", "config", "--global", "user.email", "kubesnap@kubesnap.com").Run()
	exec.Command("git", "config", "--global", "user.name", "kube-snap").Run()
}

func switchBranch(branch string) error {
	cmd := exec.Command("git", "switch", branch)
	return cmd.Run()
}

func addAll() error {
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	return err
}

func status() (string, error) {
	o, err := exec.Command("git", "status", "--porcelain").Output()
	return string(o), err
}

func commitChanges() error {
	cmd := exec.Command("git", "commit", "-m", "Snapshot taken on "+time.Now().UTC().String())
	return cmd.Run()
}

func push() error {
	cmd := exec.Command("git", "push")
	return cmd.Run()
}
