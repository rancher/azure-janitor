# Azure Janitor

A GitHub Action to cleanup Azure Resource Groups.

> By default the action will not perform the delete (i.e. it will be a dry-run). You need to explicitly set commit to `true`.

## Inputs

| Name            | Required | Description                                                           |
| --------------- | -------- | ------------------------------------------------------------------_-- |
| resource-groups | Y        | A pattern for the resource groups to delete (e.g. `e2e*`)  |
| subscription-id | Y        | The azure subscription id                                             |
| client-id       | Y        | The azure client id to use                                            |
| client-secret   | Y        | The azure client secret to use                                        |
| tenant-id       | Y        | The azure tenant id to use                                            |
| commit          | N        | Whether to perform the delete. Defaults to `false` which is a dry run |

## Example Usage

```yaml
jobs:
  cleanup:
    runs-on: ubuntu-latest
    name: Cleanup resource groups
    steps:
      - name: Cleanup
        uses: ghcr.io/rancher-sandbox/azure-janitor@abcdef1234
        with:
            resource-groups: e2e*
            subscription-id: ${{ secrets.AZURE_SUBSCRIPTION_ID}}
            client-id: ${{ secrets.AZURE_CLIENT_ID}}
            client-secret: ${{ secrets.AZURE_CLIENT_SECRET}}
            tenant-id: ${{ secrets.AZURE_TENANT_ID}}
```
