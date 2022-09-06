module github.com/rfleming71/terraform-provider-mayan-edms

go 1.16

require (
	github.com/hashicorp/terraform-plugin-docs v0.13.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.0
	github.com/rfleming71/terraform-provider-mayan-edms/client v0.0.0-00010101000000-000000000000
)

require (
	github.com/hashicorp/hcl/v2 v2.11.1 // indirect
	github.com/hashicorp/terraform-plugin-go v0.8.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.3.0
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
)

replace github.com/rfleming71/terraform-provider-mayan-edms/client => ./internal/client
