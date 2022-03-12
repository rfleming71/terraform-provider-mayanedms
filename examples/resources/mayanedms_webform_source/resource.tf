
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
