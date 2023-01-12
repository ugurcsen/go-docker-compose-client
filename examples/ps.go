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
	containers, err := cli.PsAll()
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Println(container.Labels["com.docker.compose.service"], container.State)
	}
}
