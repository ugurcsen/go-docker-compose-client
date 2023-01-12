package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
)

// Client is a docker compose client.
type Client struct {
	bin          string
	baseArgs     []string
	composePath  string
	projectName  string
	waitGroup    sync.WaitGroup
	ctx          context.Context
	dockerClient *client.Client
	UsePipes     bool
}

// NewClientWithContext create a new docker compose client with custom context.
func NewClientWithContext(ctx context.Context, composePath string) (*Client, error) {
	cmd := exec.Command("which", "docker")
	byt, err := cmd.Output()
	if err != nil {
		return nil, errors.New("docker not found")
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	c := &Client{
		bin:          string(byt[:len(byt)-1]),
		baseArgs:     []string{"compose"},
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
	cmd := exec.CommandContext(ctx, c.bin, c.baseArgs...)
	cmd.Dir = c.composePath
	return cmd
}

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

func (c *Client) Cp() error {
	panic("implement me")
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
func (c *Client) Exec(service string, commands ...string) error {
	panic("implement me")
}

// Images lists images used by the created containers.
func (c *Client) Images() error {
	panic("implement me")
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
	panic("implement me")
}

// LogsStream shows container logs as a stream.
func (c *Client) LogsStream(service string) (*Pipes, error) {
	panic("implement me")
}

// LogsAll shows all container logs.
func (c *Client) LogsAll() (*Pipes, error) {
	panic("implement me")
}

// LogsAllStream shows all container logs as a stream.
func (c *Client) LogsAllStream() (*Pipes, error) {
	panic("implement me")
}

// Pause pauses container.
func (c *Client) Pause(service string) (*Pipes, error) {
	return c.runCommand("pause", CompileNames(service))
}

// PauseAll pauses all containers.
func (c *Client) PauseAll() error {
	panic("implement me")
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
func (c *Client) Port(service string, innerPort uint16) (*Pipes, error) {
	return c.runCommand("port", CompileNames(service), fmt.Sprintf("%d", innerPort))
}

// Ps lists running containers.
func (c *Client) Ps() ([]types.Container, error) {
	return c.Containers(false)
}

// PsAll lists all containers.
func (c *Client) PsAll() ([]types.Container, error) {
	return c.Containers(true)
}

// Pull service images.
func (c *Client) Pull() error {
	panic("implement me")
}

// Push service images.
func (c *Client) Push() error {
	panic("implement me")
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
	panic("implement me")
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

func (c *Client) addWaitForCmd(cmd *exec.Cmd) {
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		cmd.Wait()
	}()
}

func (c *Client) createComposeFilterArgs() filters.Args {
	return filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.compose.project=" + c.projectName,
		},
	)
}
