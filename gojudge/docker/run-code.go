package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func RunCodeInContainer(language, code, input string) (string, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.43"))

	if err != nil {
		return "", err
	}

	var dockerImage string
	var command string

	switch language {
	case "python":
		dockerImage = "python"
		command = fmt.Sprintf("python -c \"%s\"", code) // Direct execution of Python code

	case "java":
		dockerImage = "openjdk:11"
		command = fmt.Sprintf("echo '%s' > Main.java && javac Main.java && java Main", code)

	case "javascript":
		dockerImage = "node"
		command = fmt.Sprintf("node -e \"%s\"", code) // Execute JavaScript with node

	default:
		return "", err

	}

	// Ensure the image is pulled
	err = pullImageIfNotPresent(cli, dockerImage)
	if err != nil {
		return "", err
	}

	// Create the container
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: dockerImage,
		Cmd:   []string{"sh", "-c", command},
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode:    "none",          // Disable networking
		CapDrop:        []string{"ALL"}, // Drop all capabilities
		ReadonlyRootfs: true,            // Enforce read-only filesystem
		Resources: container.Resources{
			Memory:    256 * 1024 * 1024, // 256 MB memory limit
			CPUShares: 512,               // CPU limit
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", nil
	}

	timeout := 10 * time.Second // Set execution timeout
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	case <-time.After(timeout):
	}

	// Fetch the container logs (this will contain the output of the code execution)
	out, err := cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	defer out.Close()

	// read the logs

	result := ""

	buf := make([]byte, 1024)

	for {
		n, err := out.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}

	// Stop the container after execution

	stopTimeout := 10 * time.Second           // Use time.Duration for the timeout
	timeoutSecs := int(stopTimeout.Seconds()) // Convert to int representing seconds

	if err := cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{Timeout: &timeoutSecs}); err != nil {
		return "", err
	}
	// Remove the container after execution
	if err := cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true}); err != nil {
		return "", err
	}

	return result, nil
}

func languageExt(language string) string {
	switch language {
	case "python":
		return "py"
	case "java":
		return "java"
	case "javascript":
		return "js"
	default:
		return ""
	}
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
