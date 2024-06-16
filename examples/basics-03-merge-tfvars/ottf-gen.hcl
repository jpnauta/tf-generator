generate {
  content = merge-tfvars([
    load("first.tfvars"),
    load("second.tfvars"),
    load("third.tfvars"),
  ])

  output = "tf-generator.tfvars"
}
