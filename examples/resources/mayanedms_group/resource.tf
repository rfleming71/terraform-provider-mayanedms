
resource "mayanedms_group" "editors" {
  name = "Editors"
}

resource "mayanedms_group" "test-group" {
  name  = "Test Group"
  users = [2, 3, 5, 10]
}
