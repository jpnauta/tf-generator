generate {
  content = remove-tfvar-keys(
    load("input.tfvars"),
    load("removed-keys.tfvars"),
  )

  output = "tf-generator.tfvars"
}
