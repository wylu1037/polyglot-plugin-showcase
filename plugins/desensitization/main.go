package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/plugins/desensitization/impl"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/desensitization"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: desensitization.Handshake,
		Plugins: map[string]plugin.Plugin{
			"desensitizer": &desensitization.DesensitzerGRPCPlugin{Impl: &impl.DesensitzerImpl{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
