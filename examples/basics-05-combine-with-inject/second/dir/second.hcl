module "second" {
  source = context.SOURCE_DIR
  var-b  = injectvar.var-b
}
