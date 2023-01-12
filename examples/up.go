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
	stdout, err := cli.Up()
	if err != nil {
		panic(err)
	}
	cli.Wait()
	fmt.Println(stdout.String())
}
