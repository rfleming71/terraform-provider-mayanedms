# import an existing workflow
terraform import "mayanedms_workflow_template.auto_processing" "8"

# import an existing state of the workflow
terraform import "mayanedms_workflow_template_state.auto_processing_new" "8-3"

# import an existing transition of the workflow
terraform import "mayanedms_workflow_template_transition.auto_processing_new" "8-12"