locals {
  a = 1
  b = local.a
  c = local.b
}

locals {
  d = local.b + local.c
  e = local.c + local.a
}

generate {
  content = {
    file-name = ""
    content   = "${local.d}+${local.e}"
  }
  exclude-header = true
  output         = "output.txt"
}
