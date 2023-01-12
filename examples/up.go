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
	pipes, err := cli.Up()
	if err != nil {
		panic(err)
	}
	fmt.Println(pipes.String())
	cli.Wait()

}
