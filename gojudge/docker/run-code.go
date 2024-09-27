package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func RunCodeInContainer(language, code, input, expectedOutput string) (bool, string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.43"))
	if err != nil {
		return false, "", err
	}

	dockerImage, command, err := getDockerImageAndCommand(language, code)
	if err != nil {
		return false, "", err
	}

	// Ensure the image is pulled
	if err = pullImageIfNotPresent(cli, dockerImage); err != nil {
		return false, "", err
	}

	// Create and start the container
	resp, err := createContainer(cli, dockerImage, command)
	if err != nil {
		return false, "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return false, "", err
	}

	// Wait for the container to finish
	if err := waitForContainer(cli, resp.ID); err != nil {
		return false, "", err
	}

	// Fetch and process container logs
	result, err := fetchContainerLogs(cli, resp.ID)
	if err != nil {
		return false, "", err
	}

	// Stop and remove the container
	if err := cleanupContainer(cli, resp.ID); err != nil {
		return false, "", err
	}

	log.Printf("Code executed successfully. Result: %s", result)

	// Compare results
	matches := compareOutputs(result, expectedOutput)
	return matches, result, nil
}
func getDockerImageAndCommand(language, code string) (string, string, error) {
	var dockerImage, command string

	switch language {
	case "python":
		dockerImage = "python"
		command = fmt.Sprintf("python -c \"%s\"", code)

	case "java":
		dockerImage = "openjdk:11"
		command = fmt.Sprintf("echo '%s' > Main.java && javac Main.java && java Main", code)

	case "javascript":
		dockerImage = "node"
		command = fmt.Sprintf("node -e \"%s\"", code)

	default:
		return "", "", fmt.Errorf("unsupported language: %s", language)
	}
	return dockerImage, command, nil
}

// createContainer creates a Docker container with the specified image and command.
func createContainer(cli *client.Client, dockerImage, command string) (container.CreateResponse, error) {
	return cli.ContainerCreate(context.Background(), &container.Config{
		Image: dockerImage,
		Cmd:   []string{"sh", "-c", command},
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode:    "none",
		CapDrop:        []string{"ALL"},
		ReadonlyRootfs: true,
		Resources: container.Resources{
			Memory:    256 * 1024 * 1024,
			CPUShares: 512,
		},
	}, nil, nil, "")
}

// waitForContainer waits for the specified container to finish running.
func waitForContainer(cli *client.Client, containerID string) error {
	timeout := 10 * time.Second
	statusCh, errCh := cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for container %s to finish", containerID)
	}

	return nil
}

// compareOutputs compares the cleaned result with the expected output.
func compareOutputs(result, expectedOutput string) bool {
	trimmedResult := strings.TrimSpace(result)
	trimmedExpectedOutput := strings.TrimSpace(expectedOutput)

	cleanResult := cleanString(trimmedResult)
	cleanExpected := cleanString(trimmedExpectedOutput)

	log.Printf("Trimmed Result: %q", trimmedResult)
	log.Printf("Trimmed Expected: %q", trimmedExpectedOutput)
	log.Printf("Clean Result: %q", cleanResult)
	log.Printf("Clean Expected: %q", cleanExpected)

	return cleanResult == cleanExpected
}

// cleanString cleans a string by removing non-printable characters.
func cleanString(input string) string {
	var output []rune
	for _, r := range input {
		if unicode.IsPrint(r) {
			output = append(output, r)
		}
	}
	return string(output)
}

// fetchContainerLogs retrieves and cleans logs from the specified container.
func fetchContainerLogs(cli *client.Client, containerID string) (string, error) {
	out, err := cli.ContainerLogs(context.Background(), containerID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	var result strings.Builder
	buf := make([]byte, 1024)

	for {
		n, err := out.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		result.Write(buf[:n])
	}

	cleanedResult := cleanDockerOutput(result.String())
	return cleanedResult, nil
}


// cleanupContainer stops and removes the specified container.
func cleanupContainer(cli *client.Client, containerID string) error {
	stopTimeout := 10 * time.Second
	timeoutSecs := int(stopTimeout.Seconds())

	if err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{Timeout: &timeoutSecs}); err != nil {
		return err
	}
	return cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
}

// pullImageIfNotPresent pulls the image if it's not available locally
func pullImageIfNotPresent(cli *client.Client, dockerImage string) error {
	ctx := context.Background()

	// Try to inspect the image
	_, _, err := cli.ImageInspectWithRaw(ctx, dockerImage)
	if err == nil {
		// Image already exists
		return nil
	}

	// Image not present, pull it
	fmt.Printf("Pulling image: %s\n", dockerImage)
	out, err := cli.ImagePull(ctx, dockerImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()

	// Read the output from the pull command (optional, just for logging)
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return fmt.Errorf("error reading image pull response: %w", err)
	}

	return nil
}


func cleanDockerOutput(output string) string {
	// Remove ANSI escape codes and non-printable characters
	output = strings.ReplaceAll(output, "\x1B[0m", "") // Reset
	output = strings.ReplaceAll(output, "\x1B[1m", "") // Bold
	output = strings.ReplaceAll(output, "\x1B[2m", "") // Faint
	output = strings.ReplaceAll(output, "\x1B[3m", "") // Italic
	output = strings.ReplaceAll(output, "\x1B[4m", "") // Underline
	output = strings.ReplaceAll(output, "\x1B[5m", "") // Blink
	output = strings.ReplaceAll(output, "\x1B[7m", "") // Inverse
	output = strings.ReplaceAll(output, "\x1B[8m", "") // Hidden
	output = strings.ReplaceAll(output, "\x1B[9m", "") // Strikethrough

	// Remove any other unwanted characters (like newlines and special characters)
	output = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, output)

	return strings.TrimSpace(output)
}