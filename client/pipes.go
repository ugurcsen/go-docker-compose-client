package client

import (
	"bytes"
	"io"
)

type Pipes struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

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

func pipeBytes(pipe io.ReadCloser) []byte {
	bfr := bytes.Buffer{}
	bfr.ReadFrom(pipe)
	return bfr.Bytes()
}

func (p *Pipes) StdoutString() string {
	return string(p.StdoutBytes())
}

func (p *Pipes) StderrString() string {
	return string(p.StderrBytes())
}

func (p *Pipes) StdoutBytes() []byte {
	return pipeBytes(p.Stdout)
}

func (p *Pipes) StderrBytes() []byte {
	return pipeBytes(p.Stderr)
}

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
