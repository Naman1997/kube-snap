package ansible

import (
	"os/exec"
	"strings"
)

func TriggerPlaybook(playbook string, level int) ([]byte, error) {
	verbosity := "-" + (strings.Repeat("v", level))
	cmd := exec.Command("ansible-playbook", playbook, verbosity)
	return cmd.Output()
}
