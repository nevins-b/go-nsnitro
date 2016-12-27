package nsnitro

type ServerServiceGroupBinding struct {
	Name             string `json:"name"`
	ServiceGroupName string `json:"servicegroupname"`
	MonitorName      string `json:"monitor_name"`
	State            string `json:"svrstate,omitempty"`
	Port             int    `json:"port,omitempty"`
}

func (c *Client) GetServerServiceGroupBinding(serverName string) ([]ServerServiceGroupBinding, error) {
	bindings := []ServerServiceGroupBinding{}
	err := c.fetch(nsrequest{
		Type: "server_servicegroup_binding", Name: serverName,
	}, &bindings)
	return bindings, err
}
