resource "mayanedms_index_template" "bills" {
  label = "Bills"
  slug  = "bills"
  document_types = [
    mayanedms_document_type.scan.id,
    mayanedms_document_type.email.id,
    mayanedms_document_type.pdf.id,
  ]
}

resource "mayanedms_index_template_node" "bills_node_1" {
  expression     = "{% for tag in document.tags.all %}{%if tag.label == \"Bill\" %}{{ document.metadata_value_of.company }}{% endif %}{% endfor %}"
  enabled        = true
  link_documents = false
  index_id       = mayanedms_index_template.bills.id
  parent_id      = mayanedms_index_template.bills.root_node_id
}

resource "mayanedms_index_template_node" "bills_node_2" {
  expression     = "{%if document.metadata_value_of.Year != \"\" %}{{ document.metadata_value_of.Year }}\r\n{% else %}Unknown{% endif %}"
  enabled        = true
  link_documents = false
  index_id       = mayanedms_index_template.bills.id
  parent_id      = mayanedms_index_template_node.bills_node_1.node_id
}