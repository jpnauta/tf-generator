generate {
  content = combine([
    load("hello.txt"),
    load("world.txt"),
  ])
  exclude-header = true
  output = "hello-world.txt"
}
