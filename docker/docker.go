package docker

import "os/exec"

var DockerExec string

func DockerCmd(args ...string) (string, error) {

	cmd := exec.Command(DockerExec, args...)
	out, err := cmd.Output()

	if err != nil {
		println("ERROR: " + err.Error())
		return "", err
	}

	return string(out), nil
}
