package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/kalradev/go-plugin-demo/grpc/protoc"
	"google.golang.org/grpc"
)

// Handshake are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"greeter": &GreeterPlugin{},
}

type Greeter interface {
	Greet() string
}

// GRPCClient is an implementation of Greeter that talks over RPC.
type GRPCClient struct{ client protoc.GreeterServiceClient }

func (m *GRPCClient) Greet() string {
	res, _ := m.client.Greet(context.Background(), &protoc.Empty{})
	return res.Message
}

type GRPCServer struct{ Impl Greeter }

func (m *GRPCServer) Greet(ctx context.Context, req *protoc.Empty) (*protoc.GreetResponse, error) {
	v := m.Impl.Greet()
	return &protoc.GreetResponse{Message: v}, nil
}

type GreeterPlugin struct {
	plugin.Plugin
	Impl Greeter
}

func (p *GreeterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	protoc.RegisterGreeterServiceServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GreeterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: protoc.NewGreeterServiceClient(c)}, nil
}
