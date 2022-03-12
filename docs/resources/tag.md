---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mayanedms_tag Resource - terraform-provider-mayanedms"
subcategory: ""
description: |-
  
---

# mayanedms_tag (Resource)



## Example Usage

```terraform
resource "mayanedms_tag" "incoming" {
  label = "Incoming"
  color = "#ff0000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **color** (String) The RGB color values for the tag.
- **label** (String) Short text used as the tag name.

### Optional

- **id** (String) The ID of this resource.

