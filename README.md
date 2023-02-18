# vault-secrets-operator
Kubernetes operator for the custom resources `VaultSecrets`.

## Requirements
* Go >= 1.19.3
* Helm >= 3.10.2
* Hashicorp Vault >= 1.12.3
* Kubernetes >= 1.24.7

## Installation
### With Helm
```console
$ helm dependency build
$ helm install vault-secrets-operator -n vault-secrets --create-namespace
```

#### Values

| Name           | Description              | Default value            |
|----------------|--------------------------|--------------------------|
| image.registry | Name of the docker registry | `erizzardi.mine.io`   |
| image.name     | Name of the docker image | `vault-secrets-operator` |
| image.pullPolicy | Operator image pull policy | `IfNotPresent` |
| vault.enabled | Deploy in the cluster an instance of Vault in dev mode | `true` |
| vault.dev.devRootToken | Default root token for Vault | `myroot` |
| operator.extraEnvs | Extra environment variables to pass to the operator container | `{}` |
| operator.resources{} | Resources limits to impose on to the operator container | `{}` |

The helm chart depends on Vault's, thus read [here](https://github.com/hashicorp/vault-helm) for additional Vault configurations.

### With Docker
```console
$ make image
```
This command runs the unit tests and builds the image. Use the environment variables `REGISTRY_NAME`, `APPLICATION_NAME` and `VERSION` to manipulate the tag of the resulting image. The unit tests need access to a Kubernetes cluster and a Vault instance; configure your environment with the variables `VAULT_URL`, `VAULT_TOKEN`, `KUBECONFIG` and `NAMESPACE`, as such:

```console
$ VAULT_URL=http://localhost:8200 VAULT_TOKEN=myroot KUBECONFIG=/path/to/kube/config NAMESPACE=vault-secrets make image
```
**N.B.** the namespace `NAMESPACE` needs to aready exist.

Once the image has been built, the operator can be started with
```console
$ docker run -e VAULT_URL=${VAULT_URL} -e VAULT_TOKEN=${VAULT_TOKEN} -e KUBECONFIG=${KUBECONFIG} --name vault-secrets ${REGISTRY_NAME}/${APPLICATION_NAME}:${VERSION}
```

## Usage
This operator defines a custom resource: `VaultSecrets`. An example of a manifest is as follows:
```yaml
apiVersion: erizzardi.mine.io/v1alpha1
kind: VaultSecret
metadata:
  name: vaultsecret-test
spec:
  mountPath: operator-engine
  secretPath: secretpath-test
  data:
    - name: secret1
      value: foo
    - name: secret2
      value: bar
    - name: secret3
      value: '{"foo":"bar"}'
```
The section `spec.data` contains the body of the secrets, that is going to be written at the path identified by `mountPath` and `secretPath` (**N.B.** at the moment the operator supports **only** KVv2 secret engines). The above resource's data will be written into Vault under the name `operator-engine:secretpath-test` in the form
```json
{
    "secret1": "foo",
	"secret2": "bar",
	"foo": "bar"
}
```

When the operator detects a new `VaultSecret`, it reads the definition of the object and writes a secret into Vault, according to the object's manifest. If a `VaultSecret` is patched/updated, the operator checks whether there are differences with the current configuration and the latest applied (read from the `erizzardi.mine.io/last-applied-configuration` annotation), and if so it writes into Vault a new version of the secret. If an object is deleted, the operator deletes from Vault the version of the secret associated with the deleted resource. 

## Configuration

| Flag | Env variable | Description | Default value |
|------|--------------|-------------|---------------|
| `--vault-url` | `VAULT_URL` | Complete Hashicorp Vault URL | `http://localhost:8200' |
| `--vault-token` | `VAULT_TOKEN` | Hashicorp Vault auth token | |
| `--resync-period` | `RESYNC_PERIOD` | Operator's Informer resync period, in seconds (read [here](https://groups.google.com/g/kubernetes-sig-api-machinery/c/PbSCXdLDno0) for more info) | 60 |
| `--loop-period` | `LOOP_PERIOD` | Operator's main loop period, in seconds | 1 |
| `--log-level` | `LOG_LEVEL` | Logging level | `INFO` |
| `--local-testing` | `LOCAL_TESTING` | Toggle to launch the operator in local-testing mode. It requires a kubernetes configuration, specified with `--kubeconfig` | false |
| `--kubeconfig` | `KUBECONFIG` | Location of the kubeconfig file, ignored if --local-testing is not set | `.kube/config` |

## TODO
* Improve test coverage
* Edit the operator such that it handles regular k8s Secrets
* Add support for KVv1 secret engines