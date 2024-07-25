package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
)

func LogInfo(message string) {
	hclog.Default().Info(message)
}
func LogError(message string) {
	hclog.Default().Error(message)
}
func LogWarn(message string) {
	hclog.Default().Warn(message)

}
func ExecCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
func ExitError(err error) {
	LogError(err.Error())
	if err != nil {
		os.Exit(1)
	}
}
func WriteFile(filePath string, content string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	return nil
}
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	return data, err
}

func InitNodeName(releaseName string) string {
	return fmt.Sprintf("%s-vault-0", releaseName)
}
