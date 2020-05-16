provider "demo" {
  url = "http://localhost:9200"
}

resource "demo_person" "somebody" {
  name  = "dds"
  age   = 100
  hobby = "demos"
}
