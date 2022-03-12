resource "mayanedms_watchfolder_source" "folder_1" {
  label                  = "Folder 1"
  uncompress             = "ask"
  folder_path            = "/watch_folder_1/"
  include_subdirectories = false
  document_type_id       = mayanedms_document_type.scan.id
  interval               = 600
}
