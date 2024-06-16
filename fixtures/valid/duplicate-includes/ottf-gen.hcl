generate {
  content = combine-with-inject(
    [
      load("test.hcl"),
      load("test.hcl"),
    ],
    load("test.tfvars"),
  )
  output = "tf-generator.tf"
}
