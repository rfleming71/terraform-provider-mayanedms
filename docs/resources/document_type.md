---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mayanedms_document_type Resource - terraform-provider-mayanedms"
subcategory: ""
description: |-
  
---

# mayanedms_document_type (Resource)



## Example Usage

```terraform
resource "mayanedms_document_type" "email" {
  label = "Email"
}

resource "mayanedms_document_type" "pdf" {
  label                      = "Pdf"
  delete_time_period         = 30
  delete_time_unit           = "days"
  trash_time_period          = 48
  trash_time_unit            = "hours"
  filename_generator_backend = "original"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String)

### Optional

- `delete_time_period` (Number) Amount of time after which documents of this type will be moved to the trash Defaults to `30`.
- `delete_time_unit` (String) Unit of delete_time_period. (minutes, hours, days) Defaults to `days`.
- `filename_generator_backend` (String) The class responsible for producing the actual filename used to store the uploaded documents Defaults to `uuid`.
- `filename_generator_backend_arguments` (String) The arguments for the filename generator backend as a YAML dictionary. Defaults to ``.
- `trash_time_period` (Number) Amount of time after which documents of this type in the trash will be deleted.
- `trash_time_unit` (String) Unit of trash_time_period. (minutes, hours, days)

### Read-Only

- `id` (String) The ID of this resource.


