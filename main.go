package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nevins-b/go-nsnitro/nsnitro"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("nsnitro", "A NetScaler 10+ Nitro API cli")
	ns_server   = app.Flag("server", "URL of the NetScaler server").Envar("NSNITRO_SERVER").Required().URL()
	ns_username = app.Flag("username", "NetScaler Nitro API user name").Envar("NSNITRO_USERNAME").String()
	ns_password = app.Flag("password", "NetScaler Nitro API password").Envar("NSNITRO_PASSWORD").String()

	// let's resemble the stock cli
	// https://support.citrix.com/servlet/KbServlet/download/23190-102-666049/NS-CommandReference-Guide.pdf
	// limitations: case-sensitive & no long-short options

	show                       = app.Command("show", "")
	show_responder             = show.Command("responder", "")
	show_responder_action      = show_responder.Command("action", "Print one or all responder actions")
	show_responder_action_name = show_responder_action.Arg("name", "Name of a responder action").String()
	show_responder_policy      = show_responder.Command("policy", "Print one or all responder policies")
	show_responder_policy_name = show_responder_policy.Arg("name", "Name of a responder policy").String()
	show_rewrite               = show.Command("rewrite", "")
	show_rewrite_action        = show_rewrite.Command("action", "Print one or all rewrite actions")
	show_rewrite_action_name   = show_rewrite_action.Arg("name", "Name of a rewrite action").String()
	show_rewrite_policy        = show_rewrite.Command("policy", "Print one or all rewrite policies")
	show_rewrite_policy_name   = show_rewrite_policy.Arg("name", "Name of a rewrite policy").String()
	show_server                = show.Command("server", "Print one or all servers")
	show_server_name           = show_server.Arg("name", "Name of a server").String()
	show_lb                    = show.Command("lb", "")
	show_lb_vserver            = show_lb.Command("vserver", "Print one or all lb vservers")
	show_lb_vserver_name       = show_lb_vserver.Arg("name", "Name of an lb vserver").String()
	show_lb_monitor            = show_lb.Command("monitor", "Print one or all lb monitors")
	show_lb_monitor_name       = show_lb_monitor.Arg("name", "Name of an lb monitor").String()
	show_servicegroup          = show.Command("servicegroup", "Print one or all servicegroups")
	show_servicegroup_name     = show_servicegroup.Arg("name", "Name of a servicegroup").String()
	show_ssl                   = show.Command("ssl", "")
	show_ssl_vserver           = show_ssl.Command("vserver", "Print one or all ssl vservers")
	show_ssl_vserver_name      = show_ssl_vserver.Arg("name", "Name of an ssl vserver").String()
	show_version               = show.Command("version", "Print the NetScalar version")

	add                     = app.Command("add", "")
	add_lb                  = add.Command("lb", "")
	add_lb_monitor          = add_lb.Command("monitor", "Add an lb monitor")
	add_lb_monitor_name     = add_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	add_lb_monitor_type     = add_lb_monitor.Arg("type", "Type of an lb monitor").Required().String()
	add_lb_monitor_send     = add_lb_monitor.Flag("send", "String to send to a service").String()
	add_lb_monitor_recv     = add_lb_monitor.Flag("recv", "String that expected from a service").String()
	add_lb_monitor_port     = add_lb_monitor.Flag("destport", "The port the probe is sent to").Int()
	add_lb_monitor_interval = add_lb_monitor.Flag("interval", "Frequency of the probe sent to a service").Int()
	add_lb_vserver          = add_lb.Command("vserver", "Add an lb vserver")
	add_lb_vserver_name     = add_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	add_lb_vserver_type     = add_lb_vserver.Arg("servicetype", "Type of an lb vserver").Required().String()
	add_lb_vserver_ipv4     = add_lb_vserver.Arg("ip", "IP address of an lb vserver").Required().IP()
	add_lb_vserver_port     = add_lb_vserver.Arg("port", "Port of an lb vserver").Required().Int()
	add_server              = add.Command("server", "Add a server")
	add_server_name         = add_server.Arg("name", "Name of a server").Required().String()
	add_server_ipv4         = add_server.Arg("ip", "IP address of a server").Required().IP()
	add_servicegroup        = add.Command("servicegroup", "Add a servicegroup")
	add_servicegroup_name   = add_servicegroup.Arg("name", "Name of a servicegroup").Required().String()
	add_servicegroup_type   = add_servicegroup.Arg("servicetype", "Type of servicegroup").Required().String()

	rm                   = app.Command("rm", "")
	rm_lb                = rm.Command("lb", "")
	rm_lb_monitor        = rm_lb.Command("monitor", "Remove an lb monitor")
	rm_lb_monitor_name   = rm_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	rm_lb_monitor_type   = rm_lb_monitor.Arg("type", "Type of an lb monitor").Required().String()
	rm_lb_vserver        = rm_lb.Command("vserver", "Remove an lb vserver")
	rm_lb_vserver_name   = rm_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	rm_server            = rm.Command("server", "Remove a server")
	rm_server_name       = rm_server.Arg("name", "Name of a server").Required().String()
	rm_servicegroup      = rm.Command("servicegroup", "Remove a servicegroup")
	rm_servicegroup_name = rm_servicegroup.Arg("name", "Name of a servicegroup").Required().String()

	bind                         = app.Command("bind", "")
	bind_lb                      = bind.Command("lb", "")
	bind_lb_monitor              = bind_lb.Command("monitor", "Bind an lb monitor to a service or servicegroup")
	bind_lb_monitor_name         = bind_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	bind_lb_monitor_servicegroup = bind_lb_monitor.Flag("servicegroup", "Name of a servicegroup").String()
	bind_lb_vserver              = bind_lb.Command("vserver", "Bind an lb vserver to a service or servicegroup")
	bind_lb_vserver_name         = bind_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	bind_lb_vserver_servicegroup = bind_lb_vserver.Flag("servicegroup", "Name of a servicegroup").String()
	bind_servicegroup            = bind.Command("servicegroup", "Bind a servicegroup to a service")
	bind_servicegroup_name       = bind_servicegroup.Arg("name", "Name of an servicegroup").Required().String()
	bind_servicegroup_server     = bind_servicegroup.Arg("server", "Name of a server").Required().String()
	bind_servicegroup_port       = bind_servicegroup.Arg("port", "Port of a server").Required().Int()

	unbind                         = app.Command("unbind", "")
	unbind_lb                      = unbind.Command("lb", "")
	unbind_lb_monitor              = unbind_lb.Command("monitor", "Unbind an lb monitor from a service or servicegroup")
	unbind_lb_monitor_name         = unbind_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	unbind_lb_monitor_servicegroup = unbind_lb_monitor.Flag("servicegroup", "Name of a servicegroup").String()
	unbind_lb_vserver              = unbind_lb.Command("vserver", "Unbind an lb vserver from a service or servicegroup")
	unbind_lb_vserver_name         = unbind_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	unbind_lb_vserver_servicegroup = unbind_lb_vserver.Flag("servicegroup", "Name of a servicegroup").String()
	unbind_servicegroup            = unbind.Command("servicegroup", "Unbind a servicegroup from a service")
	unbind_servicegroup_name       = unbind_servicegroup.Arg("name", "Name of an servicegroup").Required().String()
	unbind_servicegroup_server     = unbind_servicegroup.Arg("server", "Name of a server").Required().String()
	unbind_servicegroup_port       = unbind_servicegroup.Arg("port", "Port of a server").Required().Int()
)

func main() {
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		kingpin.Fatalf("%s, try %s help", err, os.Args[:1])
	}

	client, err := newClient(*ns_server, *ns_username, *ns_password)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	switch cmd {
	case "show responder action":
		doShowResponderAction(client)
	case "show responder policy":
		doShowResponderPolicy(client)
	case "show rewrite action":
		doShowRewriteAction(client)
	case "show rewrite policy":
		doShowRewritePolicy(client)
	case "show lb monitor":
		doShowLBMonitor(client)
	case "show lb vserver":
		doShowLBVServer(client)
	case "show ssl vserver":
		doShowSSLVServer(client)
	case "show server":
		doShowServer(client)
	case "show servicegroup":
		doShowServiceGroup(client)
	case "show version":
		doShowVersion(client)

	case "add lb monitor":
		doAddLBMonitor(client)
	case "rm lb monitor":
		doRemoveLBMonitor(client)

	case "add lb vserver":
		doAddLBVServer(client)
	case "rm lb vserver":
		doRemoveLBVServer(client)

	case "add server":
		doAddServer(client)
	case "rm server":
		doRemoveServer(client)

	case "add servicegroup":
		doAddServiceGroup(client)
	case "rm servicegroup":
		doRemoveServiceGroup(client)

	case "bind lb monitor":
		doBindLBMonitor(client)
	case "unbind lb monitor":
		doUnbindLBMonitor(client)

	case "bind lb vserver":
		doBindLBVServer(client)
	case "unbind lb vserver":
		doUnbindLBVServer(client)

	case "bind servicegroup":
		doBindServiceGroup(client)
	case "unbind servicegroup":
		doUnbindServiceGroup(client)
	}
}

func newClient(uri *url.URL, username, password string) (*nsnitro.Client, error) {
	if uri.User != nil {
		if username == "" {
			username = uri.User.Username()
		}
		if password == "" {
			password, _ = uri.User.Password()
		}

		uri.User = nil
	}

	config := &nsnitro.Config{
		URL:      uri.String(),
		User:     username,
		Password: password,
	}

	config.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client := nsnitro.NewClient(config)
	_, err := client.Version()
	return client, err
}
