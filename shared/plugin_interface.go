package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type ServiceIntf interface {
	HandleCommand([]string) error
}

type ServiceRPC struct {
	client *rpc.Client
}

func (this *ServiceRPC) HandleCommand(cmd []string) error {
	var resp error
	err := this.client.Call("Plugin.HandleCommand", cmd, &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type ServiceRPCServer struct {
	Impl ServiceIntf
}

func (this *ServiceRPCServer) HandleCommand(cmd []string, resp *error) error {
	*resp = this.Impl.HandleCommand(cmd)
	return nil
}

type ServicePlugin struct {
	Impl ServiceIntf
}

func (this *ServicePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ServiceRPCServer{Impl: this.Impl}, nil
}

func (this *ServicePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServiceRPC{client: c}, nil
}
