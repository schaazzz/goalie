package main

import (
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/schaazzz/goalie/shared"
)

type EchoServer struct {
	logger hclog.Logger
}

func (this *EchoServer) HandleCommand(cmd []string) error {
	this.logger.Debug("!! ECHO Server:", cmd[0])
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ECHO-SERVER",
	MagicCookieValue: "echo-server",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	echoServer := &EchoServer{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"echo-server": &shared.ServicePlugin{Impl: echoServer},
	}

	logger.Debug("message from plugin", "baz", "loo")

	go func() {
		for {
			time.Sleep(10 * time.Second)
			logger.Debug("running...")
		}
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
