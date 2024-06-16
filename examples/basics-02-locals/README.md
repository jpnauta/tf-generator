## Summary

Similar to `.tf` files, `ott-gen` has a concept of `locals {}` blocks. Any variables defined can be referenced 
in any `generate {}` block within the same file. This is especially useful if multiple `generate {}` blocks
are defined.

This example loads the files [a.tf](a.tf), [b.tf](b.tf) and [c.tf](c.tf) into their own respective variables.
These files are combined into two pairs into the variables `a-plus-b` and `a-plus-c`, then exported as the files
[a-plus-b.tf](./a-plus-b.tf) and [a-plus-c.tf](./a-plus-c.tf).
