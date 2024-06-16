generate {
  content = merge-tfvars([
    load("test/locals.tfvars"),
  ])
  output = "tf-generator.tfvars"
}
