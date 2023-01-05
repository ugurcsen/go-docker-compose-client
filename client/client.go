package client

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
)

// Client is a docker compose client
type Client struct {
	bin         string
	baseArgs    []string
	composePath string
	projectName string
	waitGroup   sync.WaitGroup
}

// NewClient create a new docker compose client
func NewClient(composePath string) *Client {
	cmd := exec.Command("which", "docker")
	byt, err := cmd.Output()
	if err != nil {
		panic("docker not found")
	}

	c := &Client{
		bin:         string(byt[:len(byt)-1]),
		baseArgs:    []string{"", "compose"},
		composePath: composePath,
		projectName: compileNames(filepath.Base(composePath)),
	}
	return c
}

func compileNames(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(str, "")
}

// initCmd initialize cmd for docker compose
func (c *Client) initCmd() *exec.Cmd {
	cmd := &exec.Cmd{}
	cmd.Path = c.bin
	cmd.Args = c.baseArgs
	cmd.Dir = c.composePath
	return cmd
}

// addCmdArgs add args to cmd
func addCmdArgs(cmd *exec.Cmd, args ...string) {
	cmd.Args = append(cmd.Args, args...)
}

func (c *Client) Wait() {
	c.waitGroup.Wait()
}

// Up create and start containers.
func (c *Client) Up() (*bytes.Buffer, error) {
	cmd := c.initCmd()
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	addCmdArgs(cmd, "up", "-d")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.addWaitForCmd(cmd)
	return stderr, nil
}

func (c *Client) baseStartStop(start bool, service *string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := c.initCmd()
	if start {
		if service == nil {
			addCmdArgs(cmd, "start")
		} else {
			addCmdArgs(cmd, "start", compileNames(*service))
		}
	} else {
		if service == nil {
			addCmdArgs(cmd, "stop")
		} else {
			addCmdArgs(cmd, "stop", compileNames(*service))
		}
	}
	cmd.Stderr = stdout
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

// Build or rebuild services
func (c *Client) Build() error {
	panic("implement me")
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
	cmd := c.initCmd()
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	addCmdArgs(cmd, "down")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	c.waitGroup.Add(1)
	c.addWaitForCmd(cmd)
	return stderr, nil
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

// Kill stops running container without removing them.
func (c *Client) Kill(service string) error {
	cmd := c.initCmd()
	addCmdArgs(cmd, "kill")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// KillAll stops running containers without removing them.
func (c *Client) KillAll() error {
	panic("implement me")
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
	cmd := c.initCmd()
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
	cmd := c.initCmd()
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

// ps lists containers.
func (c *Client) ps(all bool) ([]types.Container, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	// List all container-engines
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: all,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "label",
				Value: "com.docker.compose.project=" + c.projectName,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Ps lists running containers.
func (c *Client) Ps() ([]types.Container, error) {
	return c.ps(false)
}

// PsAll lists all containers.
func (c *Client) PsAll() ([]types.Container, error) {
	return c.ps(true)
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
func (c *Client) Restart(service string) error {
	cmd := c.initCmd()
	addCmdArgs(cmd, "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// RestartAll restart service containers.
func (c *Client) RestartAll() error {
	panic("implement me")
}

// Rm removes stopped service containers.
func (c *Client) Rm(service string) error {
	panic("implement me")
}

// RmAll removes all service containers.
func (c *Client) RmAll() error {
	panic("implement me")
}

// Run a one-off command on a service.
func (c *Client) Run(service string, commands ...string) error {
	panic("implement me")
}

func (c *Client) addWaitForCmd(cmd *exec.Cmd) {
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		cmd.Wait()
	}()
}
