package common

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	// MinSupportedProtocolVersion is the minimum protocol version we support.
	// Plugins with versions below this will be rejected.
	MinSupportedProtocolVersion = 1

	// CurrentProtocolVersion is the current protocol version.
	// New plugins should use this version.
	CurrentProtocolVersion = 1

	// MaxSupportedProtocolVersion is the maximum protocol version we support.
	// This allows forward compatibility with newer plugins.
	MaxSupportedProtocolVersion = 1

	// This is a randomly generated 64-character hex string
	// to prevent unauthorized processes from being mistakenly identified as plugins.
	MagicCookieValue = "8f3e9a2d7c1b5e4f6a8d9c2b1e5f7a3d4c6b8e1f9a2d5c7b3e8f1a4d6c9b2e5f"
)

// Version history:
// v1 (current): Initial release with GetMetadata and Execute methods

// IsProtocolVersionSupported checks if a given protocol version is supported.
// This allows the host to reject incompatible plugins early during handshake.
func IsProtocolVersionSupported(version int) bool {
	return version >= MinSupportedProtocolVersion && version <= MaxSupportedProtocolVersion
}

// Handshake is a common handshake that is shared by all plugins.
// Reference: Terraform's plugin system uses similar handshake mechanism.
//
// IMPORTANT: The MagicCookie values should NEVER be changed after release,
// as it would break compatibility with all existing plugins.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  CurrentProtocolVersion,
	MagicCookieKey:   "PLUGIN_INTERFACE",
	MagicCookieValue: MagicCookieValue,
}

// PluginMap is the map of plugins we can dispense
var PluginMap = map[string]plugin.Plugin{
	"plugin": &PluginGRPCPlugin{},
}

// PluginGRPCPlugin is the implementation of plugin.GRPCPlugin
type PluginGRPCPlugin struct {
	plugin.Plugin                 // Embedded plugin.Plugin to satisfy the interface
	Impl          PluginInterface // Impl Injection
}

func (p *PluginGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *PluginGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (any, error) {
	return &GRPCClient{client: NewPluginClient(c)}, nil
}

// GRPCClient is an implementation of PluginInterface that talks over RPC
type GRPCClient struct {
	client PluginClient
}

func (m *GRPCClient) GetMetadata() (*MetadataResponse, error) {
	resp, err := m.client.GetMetadata(context.Background(), &MetadataRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *GRPCClient) Execute(method string, params map[string]string) (*ExecuteResponse, error) {
	resp, err := m.client.Execute(context.Background(), &ExecuteRequest{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GRPCServer is the gRPC server that GRPCClient talks to
type GRPCServer struct {
	Impl PluginInterface
	UnimplementedPluginServer
}

func (m *GRPCServer) GetMetadata(ctx context.Context, req *MetadataRequest) (*MetadataResponse, error) {
	return m.Impl.GetMetadata()
}

func (m *GRPCServer) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	return m.Impl.Execute(req.Method, req.Params)
}
