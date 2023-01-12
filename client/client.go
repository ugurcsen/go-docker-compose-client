package client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sync"
)

// Client is a docker compose client.
type Client struct {
	// bin is the docker compose binary path.
	bin string
	// composePath is the path to the docker compose project path.
	composePath string
	// projectName is the name of the project.
	projectName string
	// waitGroup is used to wait for all running processes.
	waitGroup sync.WaitGroup
	// ctx is the context for the client.
	ctx context.Context
	// dockerClient is the docker client.
	dockerClient *client.Client
	// UsePipes use pipes to communicate with docker compose.
	UsePipes bool
}

// NewClientWithContext create a new docker compose client with custom context.
func NewClientWithContext(ctx context.Context, composePath string) (*Client, error) {
	err := exec.Command("docker").Run()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path.Join(composePath, "docker-compose.yml"))
	if err != nil {
		return nil, fmt.Errorf("docker-compose.yml not found in %s", composePath)
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	c := &Client{
		bin:          "docker",
		composePath:  composePath,
		projectName:  CompileNames(filepath.Base(composePath)),
		dockerClient: dockerClient,
		ctx:          ctx,
		UsePipes:     true,
	}
	return c, nil
}

// NewClient create a new docker compose client.
func NewClient(composePath string) (*Client, error) {
	return NewClientWithContext(context.Background(), composePath)
}

// CompileNames regenerate strings to fit compose name standards.
func CompileNames(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(str, "")
}

// initCmd initialize cmd for docker compose.
func (c *Client) initCmd(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.bin, "compose")
	cmd.Dir = c.composePath
	return cmd
}

// runCommand runs a docker compose command and create pipes.
func (c *Client) runCommand(args ...string) (*Pipes, error) {
	cmd := c.initCmd(c.ctx)
	addCmdArgs(cmd, args...)
	var err error
	if c.UsePipes {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}
		pipes := &Pipes{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		}

		err = cmd.Start()
		if err != nil {
			return nil, err
		}
		c.addWaitForCmd(cmd)

		return pipes, nil
	}

	err = cmd.Run()
	return nil, err
}

// addCmdArgs add args to cmd.
func addCmdArgs(cmd *exec.Cmd, args ...string) {
	cmd.Args = append(cmd.Args, args...)
}

// Wait until all working processes completed.
func (c *Client) Wait() {
	c.waitGroup.Wait()
}

// Up create and start containers.
func (c *Client) Up() (*Pipes, error) {
	return c.runCommand("up", "-d")
}

// Start service container.
func (c *Client) Start(service string) (*Pipes, error) {
	return c.runCommand("start", CompileNames(service))
}

// StartAll start all service containers.
func (c *Client) StartAll() (*Pipes, error) {
	return c.runCommand("start")
}

// Stop service container.
func (c *Client) Stop(service string) (*Pipes, error) {
	return c.runCommand("stop", CompileNames(service))
}

// StopAll stop all service containers.
func (c *Client) StopAll() (*Pipes, error) {
	return c.runCommand("stop")
}

// Build or rebuild services.
func (c *Client) Build(service string) (*Pipes, error) {
	return c.runCommand("build", CompileNames(service))
}

// BuildAll Build or rebuild services.
func (c *Client) BuildAll() (*Pipes, error) {
	return c.runCommand("build")
}

// Convert converts the compose file to platform's canonical format.
func (c *Client) Convert() (*Pipes, error) {
	return c.runCommand("convert")
}

// Create creates container for a service.
func (c *Client) Create(service string) (*Pipes, error) {
	return c.runCommand("create", CompileNames(service))
}

// CreateAll creates containers for a service.
func (c *Client) CreateAll() (*Pipes, error) {
	return c.runCommand("create")
}

// Down stops containers and removes containers, networks, volumes, and images created by up.
func (c *Client) Down() (*Pipes, error) {
	return c.runCommand("down")
}

// Events displays real time events from containers.
func (c *Client) Events() (*Pipes, error) {
	return c.runCommand("events")
}

// Exec execute a command in a running container.
func (c *Client) Exec(service string, commands ...string) (*Pipes, error) {
	return c.runCommand(append([]string{"exec", CompileNames(service)}, commands...)...)
}

// Images lists images used by the created containers.
func (c *Client) Images() ([]types.ImageSummary, error) {
	images, err := c.dockerClient.ImageList(c.ctx, types.ImageListOptions{
		Filters: c.createComposeFilterArgs(),
	})
	if err != nil {
		return nil, err
	}

	return images, nil
}

// Kill stops running container without removing them.
func (c *Client) Kill(service string) (*Pipes, error) {
	return c.runCommand("kill", CompileNames(service))
}

// KillAll stops running containers without removing them.
func (c *Client) KillAll() (*Pipes, error) {
	return c.runCommand("kill")
}

// Logs shows container logs.
func (c *Client) Logs(service string) (*Pipes, error) {
	return c.runCommand("logs", CompileNames(service))
}

// LogsStream shows container logs as a stream.
func (c *Client) LogsStream(service string) (*Pipes, error) {
	return c.runCommand("logs", "-f", CompileNames(service))
}

// LogsAll shows all container logs.
func (c *Client) LogsAll() (*Pipes, error) {
	return c.runCommand("logs")
}

// LogsAllStream shows all container logs as a stream.
func (c *Client) LogsAllStream() (*Pipes, error) {
	return c.runCommand("logs", "-f")
}

// Pause pauses container.
func (c *Client) Pause(service string) (*Pipes, error) {
	return c.runCommand("pause", CompileNames(service))
}

// PauseAll pauses all containers.
func (c *Client) PauseAll() (*Pipes, error) {
	return c.runCommand("pause")
}

// Unpause unpauses container.
func (c *Client) Unpause(service string) (*Pipes, error) {
	return c.runCommand("unpause", CompileNames(service))
}

// UnpauseAll unpauses all containers.
func (c *Client) UnpauseAll() (*Pipes, error) {
	return c.runCommand("unpause")
}

// Port displays public facing port of the container.
func (c *Client) Port(service string) ([]types.Port, error) {
	containers, err := c.Containers(false)
	ports := make([]types.Port, 0)
	if err != nil {
		return nil, err
	}
	for i, container := range containers {
		if container.Labels["com.docker.compose.service"] == service {
			ports = append(ports, containers[i].Ports...)
		}
	}

	return ports, nil
}

// Ps lists running containers.
func (c *Client) Ps() ([]types.Container, error) {
	return c.Containers(false)
}

// PsAll lists all containers.
func (c *Client) PsAll() ([]types.Container, error) {
	return c.Containers(true)
}

// Top lists processes running inside a container.
func (c *Client) Top(service string) (*Pipes, error) {
	return c.runCommand("top", CompileNames(service))
}

// TopAll lists processes running inside all containers.
func (c *Client) TopAll() (*Pipes, error) {
	return c.runCommand("top")
}

// Restart restart service container.
func (c *Client) Restart(service string) (*Pipes, error) {
	return c.runCommand("restart", CompileNames(service))
}

// RestartAll restart service containers.
func (c *Client) RestartAll() (*Pipes, error) {
	return c.runCommand("restart")
}

// Rm removes stopped service containers.
func (c *Client) Rm(service string) (*Pipes, error) {
	return c.runCommand("rm", CompileNames(service))
}

// RmAll removes all service containers.
func (c *Client) RmAll() (*Pipes, error) {
	return c.runCommand("rm")
}

// Run a one-off command on a service.
func (c *Client) Run(service string, commands ...string) (*Pipes, error) {
	return c.runCommand(append([]string{"run", "--rm", service}, commands...)...)
}

// Containers gets compose containers
func (c *Client) Containers(all bool) ([]types.Container, error) {
	// List all containers
	containers, err := c.dockerClient.ContainerList(c.ctx, types.ContainerListOptions{
		All:     all,
		Filters: c.createComposeFilterArgs(),
	})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Networks gets compose networks
func (c *Client) Networks() ([]types.NetworkResource, error) {
	networks, err := c.dockerClient.NetworkList(c.ctx, types.NetworkListOptions{Filters: c.createComposeFilterArgs()})
	if err != nil {
		return nil, err
	}
	return networks, nil
}

// Volumes gets compose volumes
func (c *Client) Volumes() ([]*types.Volume, error) {
	volumeListOkBody, err := c.dockerClient.VolumeList(c.ctx, c.createComposeFilterArgs())
	if err != nil {
		return nil, err
	}
	return volumeListOkBody.Volumes, nil
}

// addWaitForCmd increase the wait for command counter
func (c *Client) addWaitForCmd(cmd *exec.Cmd) {
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		err := cmd.Wait()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}

// createComposeFilterArgs creates compose filter args
func (c *Client) createComposeFilterArgs() filters.Args {
	return filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.compose.project=" + c.projectName,
		},
	)
}
