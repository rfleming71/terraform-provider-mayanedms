---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mayanedms_webform_source Resource - terraform-provider-mayanedms"
subcategory: ""
description: |-
  
---

# mayanedms_webform_source (Resource)



## Example Usage

```terraform
resource "mayanedms_webform_source" "default" {
  label      = "Default"
  uncompress = "ask"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **label** (String)

### Optional

- **enabled** (Boolean) Defaults to `true`.
- **id** (String) The ID of this resource.
- **uncompress** (String) Defaults to `ask`.


