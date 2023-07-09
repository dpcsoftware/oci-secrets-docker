# Oracle Cloud Secrets for Docker

This is a plugin for Docker engine to retrieve Oracle Cloud vault secrets
and mount them inside Docker containers.

## Installation

Install it from Docker Hub picking the approriate version and platform.

For example, to install version 0.1 in a node with amd64 architecture, use:

```sh
docker install --alias oci-secrets dpcsoftware/oci-secrets:0.1-amd64
```

Check all tags available in Docker Hub: https://hub.docker.com/r/dpcsoftware/oci-secrets/tags

## Configuration

This first version only supports instance principal authentication.
Thus, it can only be used inside an Oracle Cloud compute instance.

To enable instance principal authentication, please check Oracle Cloud documentation:
https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm

## Usage

To configure an Oracle Cloud secret in Docker Swarm, you can specify the secret ID (OCID)
or the vault ID (OCID) and the secret name. These parameters are configured through labels.

### Get a secret by ID

```yaml
secrets:
    my_secret:
        driver: oci-secrets
        labels:
            ocid: ocid1.vaultsecret.oc1.sa-vinhedo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Get a secret by vault ID and name

```yaml
secrets:
    my_secret:
        driver: oci-secrets
        labels:
            vault_ocid: ocid1.vault.oc1.sa-vinhedo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

In this case, it will look for a secret with name `my_secret` in the specified vault. Docker stack namespace will
be stripped from secret name in the lookup process.

If, for any reason, you would like to create a secret with a name that must be different from the Oracle secret name, you can use the `name` label. In this case, the value of the label will be used
in the lookup process.

```yaml
secrets:
    my_secret:
        driver: oci-secrets
        labels:
            name: name_in_vault
            vault_ocid: ocid1.vault.oc1.sa-vinhedo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Using vault ID can simplify multiple secrets specification when combined with
Docker compose extensions (https://docs.docker.com/compose/compose-file/11-extension/).

```yaml
x-oci-vault: &oci-main-vault
    driver: oci-secrets
    labels:
        vault_ocid: ocid1.vault.oc1.sa-vinhedo-1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

secrets:
    my_first_secret_name: *oci-main-vault
    my_second_secret_name: *oci-main-vault
```

## Development

This plugin was written in Go using the Docker plugin Go helpers (https://github.com/docker/go-plugins-helpers).

## License

Copyright 2023 Daniel Pereira Coelho

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.