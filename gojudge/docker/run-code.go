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

	var dockerImage string
	var command string

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
		return false, "", fmt.Errorf("unsupported language: %s", language)
	}

	// Ensure the image is pulled
	err = pullImageIfNotPresent(cli, dockerImage)
	if err != nil {
		return false, "", err
	}

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
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
	if err != nil {
		return false, "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return false, "", err
	}

	timeout := 10 * time.Second
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return false, "", err
		}
	case <-statusCh:
	case <-time.After(timeout):
	}

	// Fetch the container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return false, "", err
	}
	defer out.Close()

	result := ""
	buf := make([]byte, 1024)

	for {
		n, err := out.Read(buf)
		if err != nil && err != io.EOF {
			return false, "", err
		}
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}

	// Stop and remove the container
	stopTimeout := 10 * time.Second
	timeoutSecs := int(stopTimeout.Seconds())

	if err := cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{Timeout: &timeoutSecs}); err != nil {
		return false, "", err
	}
	if err := cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true}); err != nil {
		return false, "", err
	}

	log.Printf("Code executed successfully. Result: %s", result)


	trimmedResult := strings.TrimSpace(result)
	trimmedExpectedOutput := strings.TrimSpace(expectedOutput)

	cleanResult := cleanString(trimmedResult)
	cleanExpected := cleanString(trimmedExpectedOutput)
	log.Printf("Trimmed Result: %q", trimmedResult)
	log.Printf("Trimmed Expected: %q", trimmedExpectedOutput)
	log.Printf("Clean Result: %q", cleanResult)
	log.Printf("Clean Expected: %q", cleanExpected)

	
	matches := cleanResult == cleanExpected
	log.Printf("Code executed successfully. Result: %s Success: %t", cleanResult, matches)
	return matches, result, nil
}

func cleanString(input string) string {
	var output []rune
	for _, r := range input {
		if unicode.IsPrint(r) {
			output = append(output, r)
		}
	}
	return string(output)
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
