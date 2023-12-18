package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the RedirectStore client is properly configured.
	// It is also possible to use the REDIRECT_STORE_* environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "redirect-store" {
  endpoint = "http://127.0.0.1:8030"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"redirect-store": providerserver.NewProtocol6WithError(New("test")()),
	}
)
