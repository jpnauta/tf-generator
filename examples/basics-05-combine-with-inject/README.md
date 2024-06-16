## Summary

This example demonstrates perhaps the most complex function in `tf-generator`: `combine-with-inject`. This function acts 
similar to `combine`, but it also allows for the injection of additional keys and values into the result. This is
useful for injecting hierarchical settings into a shared `.tf` file, or to inject the path to a shared file.

## Summary

In this example, [tf-generator.hcl](./tf-generator.hcl) imports three files: 
[first/dir/first.hcl](first/dir/first.hcl), [second/dir/second.hcl](second/dir/second.hcl), and
[locals.tfvars](./locals.tfvars). These files are passed into `combine-with-inject(...)` to generate the
contents of [tf-generator.tf](./tf-generator.tf).

The syntax of these imported `.hcl` files should look very familiar to ordinary
`.tf` contents. By convention, injected `.tf` files use the `.hcl` extension to distinguish them from ordinary `.tf`
files. You will see two new types of available variable references in the injected `.hcl` files: `context.*`
and `injectvar.*`.

- `context.*` variables provide context related to code generation. Most notably, the `context.SOURCE_DIR` variable
  provides a way to reference files relative to the imported file.
- `injectvar.*` variables provide a way to inject `.tfvars` into the module. When this is done, the
  variable is directly substituted into the generated file as a local variable.

In the resulting output file [tf-generator.tf](./tf-generator.tf), you will notice two locals defined: `INJECTED-var-a` and
`INJECTED-var-b`. The values for these locals are directly sourced from [locals.tfvars](./locals.tfvars). The original
`injectvar.*` references ar also substituted to reference the new injected locals.

The two references to `context.SOURCE_DIR` have also been substituted in [tf-generator.tf](./tf-generator.tf). The reference
from [first/dir/first.hcl](first/dir/first.hcl) has been substituted as `first/dir/`, while the reference from
[second/dir/second.hcl](second/dir/second.hcl) has been substituted as `second/dir/`. This shows that
`context.SOURCE_DIR` contains the relative path between the project directory and the imported `.hcl` file.

## Limitations

It is assumed that any injected `.tfvars` are static in nature, and will never be overwritten via the CLI. If a 
command like `terraform apply -var 'local-tfvar=foo'` is used, the value of `local-tfvar` will **not** be
overwritten.
