resource "mayanedms_metadata_type" "company" {
  label = "Company Name"
  name  = "company"
}

resource "mayanedms_metadata_type" "year" {
  label   = "Year"
  name    = "year"
  default = "2022"
}
