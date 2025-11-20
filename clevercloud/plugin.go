package clevercloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const pluginName = "clevercloud"

// Plugin defines the Steampipe plugin for Clever Cloud.
func Plugin(ctx context.Context) *plugin.Plugin {
	return &plugin.Plugin{
		Name: pluginName,
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"clevercloud_billing": tableCleverCloudBilling(ctx),
		},
	}
}
