package main

import (
	"github.com/tdannenmuller/steampipe-plugin-clevercloud/clevercloud"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: clevercloud.Plugin,
	})
}
