generate {
  content = combine-with-inject(
    [
      load("first/dir/first.hcl"),
      load("second/dir/second.hcl"),
    ],
    load("locals.tfvars")
  )

  output = "tf-generator.tf"
}
