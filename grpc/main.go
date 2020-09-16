package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/kalradev/go-plugin-demo/grpc/shared"
)

func main() {
	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command("./plugin/plugin"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("greeter")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	greeter := raw.(shared.Greeter)
	fmt.Println(greeter.Greet())
}
