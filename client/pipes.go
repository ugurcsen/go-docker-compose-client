package client

import (
	"bytes"
	"io"
)

// Pipes is a struct that holds the pipes of a command.
type Pipes struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

// Close closes the pipes.
func (p *Pipes) Close() error {
	if p.Stdin != nil {
		err := p.Stdin.Close()
		if err != nil {
			return err
		}
	}

	if p.Stdout != nil {
		err := p.Stdout.Close()
		if err != nil {
			return err
		}
	}

	if p.Stderr != nil {
		err := p.Stderr.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// pipeBytes reads all bytes from pipe.
func pipeBytes(pipe io.ReadCloser) []byte {
	bfr := bytes.Buffer{}
	bfr.ReadFrom(pipe)
	return bfr.Bytes()
}

// StdoutString reads all bytes from stdout pipe and convert string.
func (p *Pipes) StdoutString() string {
	return string(p.StdoutBytes())
}

// StderrString reads all bytes from stderr pipe and convert string.
func (p *Pipes) StderrString() string {
	return string(p.StderrBytes())
}

// StdoutBytes reads all bytes from stdout pipe.
func (p *Pipes) StdoutBytes() []byte {
	return pipeBytes(p.Stdout)
}

// StderrBytes reads all bytes from stderr pipe.
func (p *Pipes) StderrBytes() []byte {
	return pipeBytes(p.Stderr)
}

// String returns a string representation of the pipes.
func (p *Pipes) String() string {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		c1 <- p.StdoutString()
	}()

	go func() {
		c2 <- p.StderrString()
	}()

	return <-c1 + <-c2
}
