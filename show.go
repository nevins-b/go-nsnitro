package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/fatih/color"
	"github.com/nevins-b/go-nsnitro/nsnitro"
	"github.com/rodaine/table"
)

func doShowResponderAction(client *nsnitro.Client) {
	var t table.Table
	if *show_responder_action_name == "" {
		t = newTable("Name", "Type")
		actions, err := client.GetResponderActions()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, action := range actions {
			t.AddRow(action.Name, action.Type)
		}
	} else {
		t = newPanel("ResponderAction")
		action, err := client.GetResponderAction(*show_responder_action_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", action.Name)
		t.AddRow("Type", action.Type)
		t.AddRow("Target", action.Target)
	}

	t.Print()
}

func doShowResponderPolicy(client *nsnitro.Client) {
	var t table.Table
	if *show_responder_policy_name == "" {
		t = newTable("Name", "Action")
		policies, err := client.GetResponderPolicies()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, policy := range policies {
			t.AddRow(policy.Name, policy.Action)
		}
	} else {
		t = newPanel("ResponderPolicy")
		policy, err := client.GetResponderPolicy(*show_responder_policy_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", policy.Name)
		t.AddRow("Action", policy.Action)
		t.AddRow("Rule", policy.Rule)
	}

	t.Print()
}

func doShowRewriteAction(client *nsnitro.Client) {
	var t table.Table
	if *show_rewrite_action_name == "" {
		t = newTable("Name", "Type")
		actions, err := client.GetRewriteActions()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, action := range actions {
			t.AddRow(action.Name, action.Type)
		}
	} else {
		t = newPanel("RewriteAction")
		action, err := client.GetRewriteAction(*show_rewrite_action_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", action.Name)
		t.AddRow("Type", action.Type)
		t.AddRow("Target", action.Target)
		t.AddRow("Expression", action.Expression)
	}

	t.Print()
}

func doShowRewritePolicy(client *nsnitro.Client) {
	var t table.Table
	if *show_rewrite_policy_name == "" {
		t = newTable("Name", "Action")
		policies, err := client.GetRewritePolicies()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, policy := range policies {
			t.AddRow(policy.Name, policy.Action)
		}
	} else {
		t = newPanel("RewritePolicy")
		policy, err := client.GetRewritePolicy(*show_rewrite_policy_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", policy.Name)
		t.AddRow("Action", policy.Action)
		t.AddRow("Rule", policy.Rule)
	}

	t.Print()
}

func doShowServer(client *nsnitro.Client) {
	var t table.Table
	if *show_server_name == "" {
		t = newTable("Name", "IP")
		servers, err := client.GetServers()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, server := range servers {
			t.AddRow(server.Name, server.IP)
		}
	} else {
		t = newPanel("Server")
		server, err := client.GetServer(*show_server_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", server.Name)
		t.AddRow("IP", server.IP)
		t.AddRow("State", server.State)
	}

	t.Print()
}

func doShowServiceGroup(client *nsnitro.Client) {
	var t table.Table
	if *show_servicegroup_name == "" {
		t = newTable("Name", "Type")
		servicegroups, err := client.GetServiceGroups()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, servicegroup := range servicegroups {
			t.AddRow(servicegroup.Name, servicegroup.Type)
		}
	} else {
		t = newPanel("ServiceGroup")
		servicegroup, err := client.GetServiceGroup(*show_servicegroup_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", servicegroup.Name)
		t.AddRow("Type", servicegroup.Type)
		t.AddRow("State", servicegroup.State)

		serverBindings, err := client.GetServiceGroupServerBindings(*show_servicegroup_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		for i, binding := range serverBindings {
			t.AddRow(fmt.Sprintf("Server.%d", i), binding.ServerName)
		}

		lbMonitorBindings, err := client.GetServiceGroupLBMonitorBindings(*show_servicegroup_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		for i, binding := range lbMonitorBindings {
			t.AddRow(fmt.Sprintf("LBMonitor.%d", i), binding.MonitorName)
		}
	}

	t.Print()
}

func doShowLBMonitor(client *nsnitro.Client) {
	var t table.Table
	if *show_lb_monitor_name == "" {
		t = newTable("Name", "Type")
		lbmonitors, err := client.GetLBMonitors()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, lbmonitor := range lbmonitors {
			t.AddRow(lbmonitor.Name, lbmonitor.Type)
		}
	} else {
		t = newPanel("LBMonitor")
		lbmonitor, err := client.GetLBMonitor(*show_lb_monitor_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", lbmonitor.Name)
		t.AddRow("Type", lbmonitor.Type)
		t.AddRow("State", lbmonitor.State)
		t.AddRow("Port", lbmonitor.Port)
		t.AddRow("Send", lbmonitor.Send)
		t.AddRow("Recv", lbmonitor.Recv)
		t.AddRow("Interval", lbmonitor.Interval)
	}

	t.Print()
}

func doShowLBVServer(client *nsnitro.Client) {
	var t table.Table
	if *show_lb_vserver_name == "" {
		t = newTable("Name", "Type", "IP", "Port")
		lbvservers, err := client.GetLBVServers()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, lbvserver := range lbvservers {
			t.AddRow(lbvserver.Name, lbvserver.Type, lbvserver.IP, lbvserver.Port)
		}
	} else {
		t = newPanel("LBVServer")
		lbvserver, err := client.GetLBVServer(*show_lb_vserver_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", lbvserver.Name)
		t.AddRow("Type", lbvserver.Type)
		t.AddRow("IP", lbvserver.IP)
		t.AddRow("Port", lbvserver.Port)
		t.AddRow("Mode", lbvserver.Mode)
		t.AddRow("Weight", lbvserver.Weight)
		t.AddRow("LB Method", lbvserver.LBMethod)

		serviceGroupBindings, err := client.GetLBVServerServiceGroupBindings(*show_lb_vserver_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		for i, binding := range serviceGroupBindings {
			t.AddRow(fmt.Sprintf("ServiceGroup.%d", i), binding.ServiceGroupName)
		}
	}

	t.Print()
}

func doShowSSLVServer(client *nsnitro.Client) {
	var t table.Table
	if *show_ssl_vserver_name == "" {
		t = newTable("Name")
		sslvservers, err := client.GetSSLVServers()
		if err != nil {
			kingpin.Fatalf(err.Error())
		}

		for _, sslvserver := range sslvservers {
			t.AddRow(sslvserver.Name)
		}
	} else {
		t = newPanel("SSLVServer")
		sslvserver, err := client.GetSSLVServer(*show_ssl_vserver_name)
		if err != nil {
			kingpin.Fatalf(err.Error())
		}
		t.AddRow("Name", sslvserver.Name)
		t.AddRow("Client Auth", sslvserver.ClientAuth)
		t.AddRow("Cipher Redirect", sslvserver.CipherRedirect)
		t.AddRow("DH", sslvserver.DH)
		t.AddRow("SSL2", sslvserver.SSL2)
		t.AddRow("SSL3", sslvserver.SSL3)
		t.AddRow("Session Reuse", sslvserver.SessionReuse)
		t.AddRow("Session Timeout", sslvserver.SessionTimeout)
		t.AddRow("SSL Redirect", sslvserver.SSLRedirect)
		t.AddRow("TLS1", sslvserver.TLS1)
		t.AddRow("TLS11", sslvserver.TLS11)
		t.AddRow("TLS22", sslvserver.TLS12)
	}

	t.Print()
}

func doShowVersion(client *nsnitro.Client) {
	version, err := client.Version()
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	t := newTable("Version")
	t.AddRow(version)
	t.Print()
}

func newTable(columns ...string) table.Table {
	labels := make([]interface{}, len(columns))
	for index, value := range columns {
		labels[index] = value
	}

	tbl := table.New(labels...)
	tbl.WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc())
	if len(columns) > 1 {
		tbl.WithFirstColumnFormatter(color.New(color.FgYellow).SprintfFunc())
	}
	return tbl
}

func newPanel(resourceType string) table.Table {
	tbl := table.New(resourceType, "")
	tbl.WithHeaderFormatter(func(string, ...interface{}) string {
		return ""
	})

	tbl.WithFirstColumnFormatter(color.New(color.FgYellow).SprintfFunc())
	return tbl
}
