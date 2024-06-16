## Summary

Welcome to the first example of `tf-generator`! This example shows the basics of how an [tf-generator.hcl](./tf-generator.hcl)
file can be used to generate output files.

In every project/example, the [tf-generator.hcl](./tf-generator.hcl) file determines how to generate files in the project.
In its most basic form, the `tf-generator.hcl` file simply copies the contents of one or more files, and outputs them
to a new file. In this example, the `load(...)` function loads [hello.txt](./hello.txt) and [world.txt](./world.txt),
then `combine(...)` joins the contents of both files into [hello-world.txt](./hello-world.txt).

The `exclude-header` header flag is an optional flag that disables the `#DO NOT EDIT!` message at the top of the
generated file. This is useful in some scenarios, but is otherwise unnecessary.

There are two ways to run the `tf-generator` command: generate mode, and check mode.
- In generate mode, the outputs specified in the `tf-generator.hcl` file are created or updated. This is useful for
  updating files during local development.
- In check mode, the command checks if the file contents are up-to-date, and does not perform any updates.
  This is useful for CI/CD pipelines, where the outputs should not be updated.

The unit tests in this project validate every example in check mode, so it is safe to assume that any output file like
[hello-world.txt](./hello-world.txt) is up-to-date, and demonstrates how the source files relate to the output file.
