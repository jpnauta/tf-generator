## Summary

This example shows how to subtract the keys of one `.tfvars` from another `.tfvars` file. This is useful for importing
`.tfvars` from a shared file and removing the keys that already exist in the current project.

This example imports two source files: [input.tfvars](./input.tfvars) and 
[removed-keys.tfvars](./removed-keys.tfvars). The `remove-tfvar-keys(...)` loads the `.tfvars` from both files, then
remove any keys in `removed-keys.tfvars` from `input.tfvars`. The result is written to 
[tf-generator.tfvars](./tf-generator.tfvars).

For consistency, the `remove-tfvar-keys(...)` function also sorts keys alphabetically.
