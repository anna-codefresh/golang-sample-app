package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
    		panic(err)
    }

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    	if err != nil {
    		panic(err)
    	}

    for _, container := range containers {
    		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
    	}
}
