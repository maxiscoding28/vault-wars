package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"vault-wars/util"
)

func isDoormatInstalled() error {
	if _, err := exec.LookPath("doormat"); err != nil {
		return errors.New("doormat is not installed")
	}
	util.LogInfo("Doormat installed")
	return nil
}

// func getAwsAccountNumber() (string, error) {
// 	fmt.Print("Please enter your AWS account number: ")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()
// 	awsAccountNumber := scanner.Text()

// 	if err := scanner.Err(); err != nil {
// 		return "", fmt.Errorf("error reading input: %v", err)
// 	}

// 	return awsAccountNumber, nil
// }

// func getAwsRegion() (string, error) {
// 	fmt.Print("Please enter your desired AWS region: ")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()
// 	awsRegion := scanner.Text()

// 	if err := scanner.Err(); err != nil {
// 		return "", fmt.Errorf("error reading input: %v", err)
// 	}

// 	return awsRegion, nil
// }

func GetDoormatCredentials() error {
	if err := isDoormatInstalled(); err != nil {
		return err
	}

	// if awsAccountNumber, err := getAwsAccountNumber(); err != nil {
	// 	return err
	// }

	// if awsRegion, err := getAwsRegion(); err != nil {
	// 	return err
	// }

	cmd := exec.Command("doormat", "aws", "export", "--account", "580037024501")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing doormat command: %v", err)
	}

	envVars := strings.Split(out.String(), " && ")
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "export ") {
			envVar = strings.TrimPrefix(envVar, "export ")
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if err := os.Setenv(key, value); err != nil {
					return fmt.Errorf("error setting environment variable %s: %v", key, err)
				}
			}
		}
	}

	return nil
}
