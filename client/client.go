package client

import (
	"bytes"
	"context"
	"errors"
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
func (c *Client) initCmd(ctx context.Context, buffer *bytes.Buffer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.bin, c.baseArgs...)
	cmd.Dir = c.composePath
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	return cmd
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
func (c *Client) Up() (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	addCmdArgs(cmd, "up", "-d")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.addWaitForCmd(cmd)
	return stdout, nil
}

func (c *Client) baseStartStop(start bool, service *string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	if start {
		if service == nil {
			addCmdArgs(cmd, "start")
		} else {
			addCmdArgs(cmd, "start", CompileNames(*service))
		}
	} else {
		if service == nil {
			addCmdArgs(cmd, "stop")
		} else {
			addCmdArgs(cmd, "stop", CompileNames(*service))
		}
	}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	c.addWaitForCmd(cmd)
	return stdout, nil
}

// Start service container.
func (c *Client) Start(service string) (*bytes.Buffer, error) {
	return c.baseStartStop(true, &service)
}

// StartAll start all service containers.
func (c *Client) StartAll() (*bytes.Buffer, error) {
	return c.baseStartStop(true, nil)
}

// Stop service container.
func (c *Client) Stop(service string) (*bytes.Buffer, error) {
	return c.baseStartStop(false, &service)
}

// StopAll stop all service containers.
func (c *Client) StopAll() (*bytes.Buffer, error) {
	return c.baseStartStop(false, nil)
}

func (c *Client) baseBuild(service *string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	if service == nil {
		addCmdArgs(cmd, "build")
	} else {
		addCmdArgs(cmd, "build", CompileNames(*service))
	}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	c.addWaitForCmd(cmd)
	return stdout, nil
}

// Build or rebuild services.
func (c *Client) Build(service string) (*bytes.Buffer, error) {
	return c.baseBuild(&service)
}

// BuildAll Build or rebuild services.
func (c *Client) BuildAll() (*bytes.Buffer, error) {
	return c.baseBuild(nil)
}

// Convert converts the compose file to platform's canonical format.
func (c *Client) Convert() error {
	panic("implement me")
}

func (c *Client) Cp() error {
	panic("implement me")
}

// Create creates container for a service.
func (c *Client) Create(service string) error {
	panic("implement me")
}

// CreateAll creates containers for a service.
func (c *Client) CreateAll() error {
	panic("implement me")
}

// Down stops containers and removes containers, networks, volumes, and images created by up.
func (c *Client) Down() (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	addCmdArgs(cmd, "down")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.addWaitForCmd(cmd)
	return stdout, nil
}

// Events displays real time events from containers.
func (c *Client) Events() error {
	panic("implement me")
}

// Exec execute a command in a running container.
func (c *Client) Exec(service string, commands ...string) error {
	panic("implement me")
}

// Images lists images used by the created containers.
func (c *Client) Images() error {
	panic("implement me")
}

func (c *Client) baseKill(service *string) error {
	cmd := c.initCmd(c.ctx, nil)
	if service == nil {
		addCmdArgs(cmd, "kill")
	} else {
		addCmdArgs(cmd, "kill", CompileNames(*service))
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// Kill stops running container without removing them.
func (c *Client) Kill(service string) error {
	return c.baseKill(&service)
}

// KillAll stops running containers without removing them.
func (c *Client) KillAll() error {
	return c.baseKill(nil)
}

// Logs shows container logs.
func (c *Client) Logs(service string) error {
	panic("implement me")
}

// LogsStream shows container logs as a stream.
func (c *Client) LogsStream(service string) error {
	panic("implement me")
}

// LogsAll shows all container logs.
func (c *Client) LogsAll() error {
	panic("implement me")
}

// LogsAllStream shows all container logs as a stream.
func (c *Client) LogsAllStream() error {
	panic("implement me")
}

// Pause pauses container.
func (c *Client) Pause(service string) error {
	cmd := c.initCmd(c.ctx, nil)
	addCmdArgs(cmd, "pause")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// PauseAll pauses all containers.
func (c *Client) PauseAll() error {
	panic("implement me")
}

// Unpause unpauses container.
func (c *Client) Unpause(service string) error {
	cmd := c.initCmd(c.ctx, nil)
	addCmdArgs(cmd, "unpause")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// UnpauseAll unpauses all containers.
func (c *Client) UnpauseAll() error {
	panic("implement me")
}

// Port displays public facing port of the container.
func (c *Client) Port(service string, innerPort uint16) error {
	panic("implement me")
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

func (c *Client) baseRestart(service *string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	if service == nil {
		addCmdArgs(cmd, "restart")
	} else {
		addCmdArgs(cmd, "restart", *service)
	}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.addWaitForCmd(cmd)
	return stdout, nil
}

// Restart restart service container.
func (c *Client) Restart(service string) (*bytes.Buffer, error) {
	return c.baseRestart(&service)
}

// RestartAll restart service containers.
func (c *Client) RestartAll() (*bytes.Buffer, error) {
	return c.baseRestart(nil)
}

func (c *Client) baseRm(service *string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd(c.ctx, stdout)
	if service == nil {
		addCmdArgs(cmd, "rm")
	} else {
		addCmdArgs(cmd, "rm", *service)
	}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.addWaitForCmd(cmd)
	return stdout, nil
}

// Rm removes stopped service containers.
func (c *Client) Rm(service string) (*bytes.Buffer, error) {
	return c.baseRm(&service)
}

// RmAll removes all service containers.
func (c *Client) RmAll() (*bytes.Buffer, error) {
	return c.baseRm(nil)
}

// Run a one-off command on a service.
func (c *Client) Run(service string, commands ...string) (*bytes.Buffer, error) {
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
