package clevercloud

// cleverCloudConfig defines the configuration for the Clever Cloud plugin connection.
type cleverCloudConfig struct {
	Token          *string `hcl:"token"`
	OrganizationID *string `hcl:"organization_id"`
	APIEndpoint    *string `hcl:"api_endpoint"`
}

// ConfigInstance returns a new instance of the connection config
func ConfigInstance() interface{} {
	return &cleverCloudConfig{}
}
