package desensitization

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DESENSITIZER_PLUGIN",
	MagicCookieValue: "desensitize",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"desensitizer": &DesensitzerGRPCPlugin{},
}

// DesensitzerGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type DesensitzerGRPCPlugin struct {
	plugin.Plugin
	Impl Desensitizer // Impl Injection
}

func (p *DesensitzerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterDesensitizerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *DesensitzerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: NewDesensitizerClient(c)}, nil
}

// GRPCClient is an implementation of Desensitizer that talks over RPC.
type GRPCClient struct {
	client DesensitizerClient
}

func (m *GRPCClient) DesensitizeName(name string) (string, error) {
	resp, err := m.client.DesensitizeName(context.Background(), &DesensitizeRequest{
		Data: name,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func (m *GRPCClient) DesensitizeTelNo(telNo string) (string, error) {
	resp, err := m.client.DesensitizeTelNo(context.Background(), &DesensitizeRequest{
		Data: telNo,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func (m *GRPCClient) DesensitizeIDNumber(idNumber string) (string, error) {
	resp, err := m.client.DesensitizeIDNumber(context.Background(), &DesensitizeRequest{
		Data: idNumber,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func (m *GRPCClient) DesensitizeEmail(email string) (string, error) {
	resp, err := m.client.DesensitizeEmail(context.Background(), &DesensitizeRequest{
		Data: email,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func (m *GRPCClient) DesensitizeBankCard(cardNumber string) (string, error) {
	resp, err := m.client.DesensitizeBankCard(context.Background(), &DesensitizeRequest{
		Data: cardNumber,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func (m *GRPCClient) DesensitizeAddress(address string) (string, error) {
	resp, err := m.client.DesensitizeAddress(context.Background(), &DesensitizeRequest{
		Data: address,
	})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Desensitizer
	UnimplementedDesensitizerServer
}

func (m *GRPCServer) DesensitizeName(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeName(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}

func (m *GRPCServer) DesensitizeTelNo(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeTelNo(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}

func (m *GRPCServer) DesensitizeIDNumber(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeIDNumber(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}

func (m *GRPCServer) DesensitizeEmail(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeEmail(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}

func (m *GRPCServer) DesensitizeBankCard(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeBankCard(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}

func (m *GRPCServer) DesensitizeAddress(ctx context.Context, req *DesensitizeRequest) (*DesensitizeResponse, error) {
	result, err := m.Impl.DesensitizeAddress(req.Data)
	if err != nil {
		return nil, err
	}
	return &DesensitizeResponse{Result: result}, nil
}
