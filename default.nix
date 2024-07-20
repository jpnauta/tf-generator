final: prev:
{
  tf-generator = prev.buildGoModule {
    name = "tf-generator";

    src = ./.;

    vendorHash = "sha256-InRf1tpAVo2x6IFUmqERfRcT9P0qAd3lEvZ1u5Nfbn4=";

    meta = {
      description = "Simple generator for Terraform/OpenTofu";
      homepage = "https://github.com/jpnauta/tf-generator";
      license = prev.lib.licenses.mit;
      maintainers = with prev.lib.maintainers; [ jpnauta ];
    };
  };
}