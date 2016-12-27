package nsnitro

type LBVServer struct {
	Name           string `json:"name"`
	Type           string `json:"servicetype"`
	IP             string `json:"ipv46"`
	Port           int    `json:"port"`
	Mode           string `json:"m,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	LBMethod       string `json:"lbmethod,omitempty"`
	CurState       string `json:"curstate,omitempty"`
	EffectiveState string `json:"effectivestate,omitempty"`
	Health         string `json:"health,omitempty"`
	TotalServices  string `json:"totalservices,omitempty"`
	ActiveServices string `json:"activeservices,omitempty"`
}

func NewLBVServer(name, serviceType, ip string, port int) LBVServer {
	return LBVServer{Name: name, Type: serviceType, IP: ip, Port: port}
}

func (c *Client) GetLBVServer(name string) (LBVServer, error) {
	lbvservers := []LBVServer{}
	err := c.fetch(nsrequest{
		Type: "lbvserver", Name: name,
	}, &lbvservers)

	if err != nil {
		return LBVServer{}, err
	}

	return lbvservers[0], nil
}

func (c *Client) GetLBVServers() ([]LBVServer, error) {
	lbvservers := []LBVServer{}
	err := c.fetch(nsrequest{Type: "lbvserver"}, &lbvservers)
	return lbvservers, err
}

func (c *Client) AddLBVServer(lbvserver LBVServer) error {
	return c.do("POST",
		nsrequest{
			Type: "lbvserver",
		},
		&nsresource{
			LBVServer: &lbvserver,
		})
}

func (c *Client) RemoveLBVServer(name string) error {
	return c.do("DELETE",
		nsrequest{
			Type: "lbvserver",
			Name: name,
		}, nil)
}
