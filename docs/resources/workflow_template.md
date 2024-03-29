---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mayanedms_workflow_template Resource - terraform-provider-mayanedms"
subcategory: ""
description: |-
  
---

# mayanedms_workflow_template (Resource)



## Example Usage

```terraform
resource "mayanedms_workflow_template" "auto_processing" {
  label         = "Automated Processing"
  internal_name = "automated_processing"
  document_types = [
    mayanedms_document_type.scan.id,
    mayanedms_document_type.pdf.id,
    mayanedms_document_type.email.id,
  ]
}

resource "mayanedms_workflow_template_state" "auto_processing_new" {
  label             = "New"
  completion        = 0
  initial           = true
  workflow_template = mayanedms_workflow_template.auto_processing.id
}

resource "mayanedms_workflow_template_state" "auto_processing_ocr_finished" {
  label             = "OCR Finished"
  completion        = 25
  initial           = false
  workflow_template = mayanedms_workflow_template.auto_processing.id
}

resource "mayanedms_workflow_template_transition" "auto_processing_new" {
  label             = "OCR Trigger"
  origin_state      = mayanedms_workflow_template_state.auto_processing_new.id
  destination_state = mayanedms_workflow_template_state.auto_processing_ocr_finished.id
  workflow_template = mayanedms_workflow_template.auto_processing.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `internal_name` (String) This value will be used by other apps to reference this workflow. Can only contain letters, numbers, and underscores.
- `label` (String) Short text to describe the workflow

### Optional

- `document_types` (Set of Number)

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# import an existing workflow
terraform import "mayanedms_workflow_template.auto_processing" "8"

# import an existing state of the workflow
terraform import "mayanedms_workflow_template_state.auto_processing_new" "8-3"

# import an existing transition of the workflow
terraform import "mayanedms_workflow_template_transition.auto_processing_new" "8-12"
```
