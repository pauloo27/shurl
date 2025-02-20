# SHURL - HELM CHART

Deploy Shurl with a helm chart.

# Usage

First, add the helm repo to your helm client:

> `helm repo add dbcafe https://code.db.cafe/api/packages/pauloo27/helm`

> `helm repo update`

Then, generate the default values.yml: 

> `helm show values dbcafe/shurl > my-values-file.yaml`

Finally, install the chart:

> `helm install shurl dbcafe/shurl -f my-values-file.yaml -n <namespace>`

## Deploy a new version of the chart

First, pack the chart:

> `make pack`

Then, push the chart to the repo:

> `REGISTRY_PASSWORD=<password> make push ./shurl-<version>.tgz`
