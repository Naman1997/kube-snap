package ansible

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"kubesnap.io/kubesnap/internal/utilities"
)

func TriggerPlaybook(playbook string, level int, eventYaml string) {
	verbosity := "-" + (strings.Repeat("v", level))
	path := "event.yaml"
	os.Remove(path)
	utilities.CreateFile("event", eventYaml)
	cmd := exec.Command("ansible-playbook", playbook, verbosity, "--extra-vars", "@"+path)
	var errb bytes.Buffer
	cmd.Stderr = &errb
	response, err := cmd.Output()
	utilities.CheckIfError(err, "Auto-remidiation step failed: "+errb.String())
	utilities.CreateTimedLog(string(response))
}
