# client

## Functions

### func [CompileNames](/client.go#L57)

`func CompileNames(str string) string`

CompileNames regenerate strings to fit compose name standards.

## Types

### type [Client](/client.go#L17)

`type Client struct { ... }`

Client is a docker compose client.

#### func [NewClient](/client.go#L52)

`func NewClient(composePath string) (*Client, error)`

NewClient create a new docker compose client.

#### func [NewClientWithContext](/client.go#L28)

`func NewClientWithContext(ctx context.Context, composePath string) (*Client, error)`

NewClientWithContext create a new docker compose client with custom context.

#### func (*Client) [Build](/client.go#L156)

`func (c *Client) Build(service string) (*bytes.Buffer, error)`

Build or rebuild services.

#### func (*Client) [BuildAll](/client.go#L161)

`func (c *Client) BuildAll() (*bytes.Buffer, error)`

BuildAll Build or rebuild services.

#### func (*Client) [Containers](/client.go#L371)

`func (c *Client) Containers(all bool) ([]types.Container, error)`

Containers gets compose containers

#### func (*Client) [Convert](/client.go#L166)

`func (c *Client) Convert() error`

Convert converts the compose file to platform's canonical format.

#### func (*Client) [Cp](/client.go#L170)

`func (c *Client) Cp() error`

#### func (*Client) [Create](/client.go#L175)

`func (c *Client) Create(service string) error`

Create creates container for a service.

#### func (*Client) [CreateAll](/client.go#L180)

`func (c *Client) CreateAll() error`

CreateAll creates containers for a service.

#### func (*Client) [Down](/client.go#L185)

`func (c *Client) Down() (*bytes.Buffer, error)`

Down stops containers and removes containers, networks, volumes, and images created by up.

#### func (*Client) [Events](/client.go#L198)

`func (c *Client) Events() error`

Events displays real time events from containers.

#### func (*Client) [Exec](/client.go#L203)

`func (c *Client) Exec(service string, commands ...string) error`

Exec execute a command in a running container.

#### func (*Client) [Images](/client.go#L208)

`func (c *Client) Images() error`

Images lists images used by the created containers.

#### func (*Client) [Kill](/client.go#L227)

`func (c *Client) Kill(service string) error`

Kill stops running container without removing them.

#### func (*Client) [KillAll](/client.go#L232)

`func (c *Client) KillAll() error`

KillAll stops running containers without removing them.

#### func (*Client) [Logs](/client.go#L237)

`func (c *Client) Logs(service string) error`

Logs shows container logs.

#### func (*Client) [LogsAll](/client.go#L247)

`func (c *Client) LogsAll() error`

LogsAll shows all container logs.

#### func (*Client) [LogsAllStream](/client.go#L252)

`func (c *Client) LogsAllStream() error`

LogsAllStream shows all container logs as a stream.

#### func (*Client) [LogsStream](/client.go#L242)

`func (c *Client) LogsStream(service string) error`

LogsStream shows container logs as a stream.

#### func (*Client) [Networks](/client.go#L385)

`func (c *Client) Networks() ([]types.NetworkResource, error)`

Networks gets compose networks

#### func (*Client) [Pause](/client.go#L257)

`func (c *Client) Pause(service string) error`

Pause pauses container.

#### func (*Client) [PauseAll](/client.go#L268)

`func (c *Client) PauseAll() error`

PauseAll pauses all containers.

#### func (*Client) [Port](/client.go#L289)

`func (c *Client) Port(service string, innerPort uint16) error`

Port displays public facing port of the container.

#### func (*Client) [Ps](/client.go#L294)

`func (c *Client) Ps() ([]types.Container, error)`

Ps lists running containers.

#### func (*Client) [PsAll](/client.go#L299)

`func (c *Client) PsAll() ([]types.Container, error)`

PsAll lists all containers.

#### func (*Client) [Pull](/client.go#L304)

`func (c *Client) Pull() error`

Pull service images.

#### func (*Client) [Push](/client.go#L309)

`func (c *Client) Push() error`

Push service images.

#### func (*Client) [Restart](/client.go#L330)

`func (c *Client) Restart(service string) (*bytes.Buffer, error)`

Restart restart service container.

#### func (*Client) [RestartAll](/client.go#L335)

`func (c *Client) RestartAll() (*bytes.Buffer, error)`

RestartAll restart service containers.

#### func (*Client) [Rm](/client.go#L356)

`func (c *Client) Rm(service string) (*bytes.Buffer, error)`

Rm removes stopped service containers.

#### func (*Client) [RmAll](/client.go#L361)

`func (c *Client) RmAll() (*bytes.Buffer, error)`

RmAll removes all service containers.

#### func (*Client) [Run](/client.go#L366)

`func (c *Client) Run(service string, commands ...string) (*bytes.Buffer, error)`

Run a one-off command on a service.

#### func (*Client) [Start](/client.go#L119)

`func (c *Client) Start(service string) (*bytes.Buffer, error)`

Start service container.

#### func (*Client) [StartAll](/client.go#L124)

`func (c *Client) StartAll() (*bytes.Buffer, error)`

StartAll start all service containers.

#### func (*Client) [Stop](/client.go#L129)

`func (c *Client) Stop(service string) (*bytes.Buffer, error)`

Stop service container.

#### func (*Client) [StopAll](/client.go#L134)

`func (c *Client) StopAll() (*bytes.Buffer, error)`

StopAll stop all service containers.

#### func (*Client) [Unpause](/client.go#L273)

`func (c *Client) Unpause(service string) error`

Unpause unpauses container.

#### func (*Client) [UnpauseAll](/client.go#L284)

`func (c *Client) UnpauseAll() error`

UnpauseAll unpauses all containers.

#### func (*Client) [Up](/client.go#L81)

`func (c *Client) Up() (*bytes.Buffer, error)`

Up create and start containers.

#### func (*Client) [Volumes](/client.go#L394)

`func (c *Client) Volumes() ([]*types.Volume, error)`

Volumes gets compose volumes

#### func (*Client) [Wait](/client.go#L76)

`func (c *Client) Wait()`

Wait until all working processes completed.