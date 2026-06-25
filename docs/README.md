# deploy-terraform — docs

**Terraform deploy.** Generic IaC driver — `terraform init/apply/destroy` in a config directory.

## Install

```bash
togo install togo-framework/deploy-terraform
```

Registers on the [`deploy`](https://github.com/togo-framework/deploy) base; select it with **deploy.provider in togo.yaml (or DEPLOY_PROVIDER)**, then use **`togo deploy`**.

## Interface

`Deployer` — `Provision`/`Deploy`/`Destroy`/`Status` over a `Spec{App,Dir,BuildCmd,Host,User,Image,Region,Domain}` built from your `togo.yaml`.

## Configuration

| Env var | Description |
|---|---|
| `TF_DIR` | Path to the Terraform config directory (default `<project>/infra`). |

## Usage & notes

Runs Terraform in `TF_DIR` (default `<project>/infra`); `Deploy`→apply, `Destroy`→destroy. Pass values from `togo.yaml` as `-var`.

## Example

```bash
togo deploy --provider terraform --dry-run   # preview the plan
togo deploy --provider terraform
```

## Links

- [Terraform](https://developer.hashicorp.com/terraform/docs)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/deploy-terraform)
