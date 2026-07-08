package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func AdvertiseRefs(service string, owner string, name string) ([]byte, error) {
	subcommand, valid := serviceSubcommand(service)
	if !valid {
		return nil, fmt.Errorf(UnknownService, service)
	}

	command := exec.Command("git", subcommand, "--stateless-rpc", "--advertise-refs", RepoPath(owner, name))
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if runError := command.Run(); runError != nil {
		return nil, fmt.Errorf(ServiceFailed, runError, stderr.String())
	}

	var buffer bytes.Buffer
	buffer.WriteString(packetLine(ServicePrefix + service + "\n"))
	buffer.WriteString(FlushPacket)
	buffer.Write(stdout.Bytes())
	return buffer.Bytes(), nil
}

func UploadPack(owner string, name string, input io.Reader, output io.Writer) error {
	command := exec.Command("git", "upload-pack", "--stateless-rpc", RepoPath(owner, name))
	return pipeService(command, input, output)
}

func ReceivePack(owner string, name string, token string, input io.Reader, output io.Writer) error {
	command := exec.Command("git",
		"-c", "credential.helper=",
		"-c", CredentialHelper,
		"-c", "core.hooksPath="+HooksDir(),
		"receive-pack", "--stateless-rpc", RepoPath(owner, name),
	)
	command.Env = append(os.Environ(), GitTokenEnv+"="+token)
	return pipeService(command, input, output)
}

func pipeService(command *exec.Cmd, input io.Reader, output io.Writer) error {
	command.Stdin = input
	command.Stdout = output
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if runError := command.Run(); runError != nil {
		return fmt.Errorf(ServiceFailed, runError, stderr.String())
	}
	return nil
}

func serviceSubcommand(service string) (string, bool) {
	switch service {
	case ServiceUploadPack:
		return "upload-pack", true
	case ServiceReceivePack:
		return "receive-pack", true
	default:
		return "", false
	}
}

func packetLine(payload string) string {
	return fmt.Sprintf("%04x%s", len(payload)+4, payload)
}
