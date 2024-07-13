locals {
  a        = load("a.tf")
  b        = load("b.tf")
  c        = load("c.tf")
  a-plus-b = combine([local.a, local.b])
  a-plus-c = combine([local.a, local.c])
}

generate {
  content = local.a-plus-b
  output  = "a-plus-b.tf"
}

generate {
  content = local.a-plus-c
  output  = "a-plus-c.tf"
}
