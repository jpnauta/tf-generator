generate {
  content = merge-tfvars([
    load("locals.tfvars"),
  ])
  output = "tf-generator.tfvars"
}
