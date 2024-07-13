generate {
  content = combine-with-inject([
    load("test.hcl")
  ], load("locals.tfvars"))
  output = "tf-generator.tf"
}
