locals {
  a = 1
}

locals {
  a = 2
}

generate {
  content = {
    file-name = ""
    content   = local.a
  }
  output = "output.txt"
}
