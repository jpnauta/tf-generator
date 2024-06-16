generate {
  content = combine-with-inject([
    load("test.hcl"),
  ], load("empty.tfvars"))
  output = "tf-generator.tf"
}
