generate {
  content = load("does-not-exist.tfvars")
  output = "tf-generator.tfvars"
}