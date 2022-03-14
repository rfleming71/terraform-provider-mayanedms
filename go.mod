module github.com/rfleming71/terraform-provider-mayan-edms

go 1.16

require (
	github.com/hashicorp/terraform-plugin-docs v0.5.1
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.11.0
	github.com/rfleming71/terraform-provider-mayan-edms/client v0.0.0-00010101000000-000000000000
)

require google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect

replace github.com/rfleming71/terraform-provider-mayan-edms/client => ./internal/client
