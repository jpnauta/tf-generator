## Summary

This example shows how to merge multiple `.tfvars` files into a single `.tfvars` file. This is useful when you want
to share `.tfvars` between terraform projects and/or create a hierarchy of `.tfvars`.

This example imports three source files: [first.tfvars](./first.tfvars), [second.tfvars](./second.tfvars), and
[third.tfvars](./third.tfvars). The `merge-tfvars(...)` function scans the `.tfvars` in each file, and combines them
into one. If any duplicate keys are found, the first value is used. Looking at [tf-generator.tfvars](./tf-generator.tfvars),
we can see that `first.tfvars` has the highest priority, followed by `second.tfvars`, and finally `third.tfvars`.

For consistency, the `merge-tfvars(...)` function also sorts keys alphabetically.
