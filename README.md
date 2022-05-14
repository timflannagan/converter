# Overview

The converter CLI is responsible for converting a directory that contains an individual OLM registry+v1 bundle into a decomposed set of plain Kubernetes manifests.

## Usage

Pass an input directory that contains a registry+v1 bundle directory structure:

```bash
go run main.go ./bundle
```

Where that input `./bundle` directory may resemble the following structure:

```tree
bundle
    ├── bundle.Dockerfile
    ├── manifests
    │   ├── acme.cert-manager.io_challenges.yaml
    │   ├── acme.cert-manager.io_orders.yaml
    │   ├── cert-manager.clusterserviceversion.yaml
    │   ├── cert-manager-edit_rbac.authorization.k8s.io_v1_clusterrole.yaml
    │   ├── cert-manager.io_certificaterequests.yaml
    │   ├── cert-manager.io_certificates.yaml
    │   ├── cert-manager.io_clusterissuers.yaml
    │   ├── cert-manager.io_issuers.yaml
    │   ├── cert-manager_v1_service.yaml
    │   ├── cert-manager-view_rbac.authorization.k8s.io_v1_clusterrole.yaml
    │   ├── cert-manager-webhook_v1_configmap.yaml
    │   └── cert-manager-webhook_v1_service.yaml
    ├── metadata
    │   └── annotations.yaml
    └── tests
        └── scorecard
            └── config.yaml
```

By default, the CLI will output the decomposed manifests into the root `plain` directory. You can configure the output directory path through the `--output-dir` CLI flag:

```bash
go run main.go --output-dir ./manifests bundle
```
