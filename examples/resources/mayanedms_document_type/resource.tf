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
