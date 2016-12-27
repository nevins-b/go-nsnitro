package nsnitro

type ServiceGroup struct {
	Name           string `json:"servicegroupname"`
	Type           string `json:"servicetype,omitempty"`
	State          string `json:"state,omitempty"`
	Server         string `json:"servername,omitempty"`
	Port           int    `json:"port,omitempty"`
	Graceful       string `json:"graceful,omitempty"`
	Delay          int    `json:"delay,omitempty"`
	EffectiveState string `json:"servicegroupeffectivestate,omitempty"`
}

func NewServiceGroup(name, serviceType string) ServiceGroup {
	return ServiceGroup{Name: name, Type: serviceType}
}

func (c *Client) GetServiceGroup(name string) (ServiceGroup, error) {
	servicegroups := []ServiceGroup{}
	err := c.fetch(nsrequest{
		Type: "servicegroup", Name: name,
	}, &servicegroups)
	if err != nil {
		return ServiceGroup{}, err
	}
	return servicegroups[0], nil
}

func (c *Client) GetServiceGroups() ([]ServiceGroup, error) {
	servicegroups := []ServiceGroup{}
	err := c.fetch(nsrequest{Type: "servicegroup"}, &servicegroups)
	return servicegroups, err
}

func (c *Client) RemoveServiceGroup(name string) error {
	return c.do("DELETE",
		nsrequest{
			Type: "servicegroup",
			Name: name,
		}, nil)
}

func (c *Client) AddServiceGroup(serviceGroup ServiceGroup) error {
	return c.do("POST",
		nsrequest{
			Type: "servicegroup",
		},
		&nsresource{
			ServiceGroup: &serviceGroup,
		})
}

func (c *Client) ServiceGroupPerformAction(serviceGroup ServiceGroup, params *map[string]string) error {
	return c.do("POST",
		nsrequest{
			Type: "servicegroup",
		},
		&nsresource{
			Params:       params,
			ServiceGroup: &serviceGroup,
		})
}
