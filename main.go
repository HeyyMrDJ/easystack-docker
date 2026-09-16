package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moby/moby/client"
)

// ContainerInfo is a simplified view of the fields we care about.
type ContainerInfo struct {
	ID     string   `json:"id"`
	Names  []string `json:"names"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	State  string   `json:"state"`
}

type listContainersArgs struct{}

type listContainersResult struct {
	Containers []ContainerInfo `json:"containers"`
}

func main() {
	cli, err := client.New(
		client.WithHost("tcp://127.0.0.1:2375"),
	)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	server := mcp.NewServer(&mcp.Implementation{Name: "easystack-docker", Version: "0.1.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_containers",
		Description: "List Docker containers on the configured host",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listContainersArgs) (*mcp.CallToolResult, listContainersResult, error) {
		infos, err := listContainers(ctx, cli)
		if err != nil {
			return nil, listContainersResult{}, err
		}
		return nil, listContainersResult{Containers: infos}, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

func listContainers(ctx context.Context, cli *client.Client) ([]ContainerInfo, error) {
	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	infos := make([]ContainerInfo, 0, len(containers.Items))
	for _, c := range containers.Items {
		infos = append(infos, ContainerInfo{
			ID:     c.ID,
			Names:  c.Names,
			Image:  c.Image,
			Status: c.Status,
			State:  string(c.State),
		})
	}
	return infos, nil
}
