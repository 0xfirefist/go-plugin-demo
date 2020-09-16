package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/kalradev/go-plugin-demo/grpc/shared"
)

type Greeter struct{}

func (Greeter) Greet() string {
	return "Hello!"
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"greeter": &shared.GreeterPlugin{Impl: &Greeter{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
