resource "mayanedms_index_template" "bills" {
  label = "Bills"
  slug  = "bills"
  document_types = [
    mayanedms_document_type.scan.id,
    mayanedms_document_type.email.id,
    mayanedms_document_type.pdf.id,
  ]
}
