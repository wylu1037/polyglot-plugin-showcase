package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/plugins/desensitization/adapter"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

func main() {
	// Use the common plugin interface with type-specific interface name
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			"desensitization": &common.PluginGRPCPlugin{
				Impl: adapter.NewDesensitizationAdapter(),
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
