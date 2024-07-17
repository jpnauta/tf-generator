final: prev:
{
  tf-generator = prev.buildGoModule {
    name = "tf-generator";

    src = ./.;

    vendorHash = "sha256-TIDg83EHqgzL9QOOsgGPm6f50qvnFbZPHwVo6pBaPQQ=";

    meta = {
      description = "Simple generator for Terraform/OpenTofu";
      homepage = "https://github.com/jpnauta/tf-generator";
      license = prev.lib.licenses.mit;
      maintainers = with prev.lib.maintainers; [ jpnauta ];
    };
  };
}