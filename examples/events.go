package main

import (
	"context"
	"fmt"
	"github.com/ugurcsen/go-docker-compose-client/client"
	"os"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	cli, _ := client.NewClientWithContext(ctx, os.Getenv("project_path"))
	pipes, err := cli.Events()
	if err != nil {
		panic(err)
	}
	bfr := make([]byte, 1024)
	i := 0
	for {
		i, err = pipes.Stdout.Read(bfr)
		if err != nil {
			break
		}
		if i > 0 {
			fmt.Println(string(bfr[:i]))
		}
	}
	cli.Wait()
}
