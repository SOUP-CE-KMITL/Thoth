package thoth

import (
	"os/exec"
	"strconv"
)

/**
* scale-out replicas via cli
**/
func ScaleOutViaCli(replicas int, namespace, name string) (string, error) {
	var err error
	var cmd []byte
	cmd, err = exec.Command("kubectl", "scale", "--replicas="+strconv.Itoa(replicas), "rc", name, "--namespace", namespace).Output()
	return string(cmd), err
}
