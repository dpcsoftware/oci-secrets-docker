package main

import (
	"context"
	"fmt"
	"encoding/base64"
	"strings"

	docker_secrets "github.com/docker/go-plugins-helpers/secrets"
	oci_auth "github.com/oracle/oci-go-sdk/v65/common/auth"
	oci_secrets "github.com/oracle/oci-go-sdk/v65/secrets"
)

type OCIDriver struct {
	client oci_secrets.SecretsClient
}

func (d *OCIDriver) Init() error {
	provider, err := oci_auth.InstancePrincipalConfigurationProvider()
	if (err != nil) {
		fmt.Println("Failed to get OCI instance principal configuration: ", err.Error())
		return err
	}

	client, err := oci_secrets.NewSecretsClientWithConfigurationProvider(provider)
	if (err != nil) {
		fmt.Println("Failed to create OCI secrets client: ", err.Error())
		return err
	}
	d.client = client

	return nil
}

func (d OCIDriver) Get(req docker_secrets.Request) docker_secrets.Response {
	var secretBundle oci_secrets.SecretBundle

	ctx := context.Background()

	ocid := req.SecretLabels["ocid"]
	if (ocid != "") {
		bundleRequest := oci_secrets.GetSecretBundleRequest{SecretId: &ocid}
		bundleResponse, err := d.client.GetSecretBundle(ctx, bundleRequest)
		if (err != nil) {
			fmt.Println("Failed to get OCI secret value: ", err.Error())
			return docker_secrets.Response{nil, err.Error(), true}
		}
		secretBundle = bundleResponse.SecretBundle
	} else {
		vaultOcid := req.SecretLabels["vault_ocid"]
		if (vaultOcid != "") {
			secretName := req.SecretLabels["name"]
			if (secretName == "") {
				secretName = req.SecretName

				// When using the secret name, strip the stack namespace
				namespace := req.SecretLabels["com.docker.stack.namespace"]
				if (namespace != "") {
					secretName = strings.TrimPrefix(secretName, namespace)
					secretName = strings.TrimLeft(secretName, "_-")
				}
			}

			bundleRequest := oci_secrets.GetSecretBundleByNameRequest{
				SecretName: &secretName,
				VaultId: &vaultOcid,
			}
			bundleResponse, err := d.client.GetSecretBundleByName(ctx, bundleRequest)
			if (err != nil) {
				fmt.Println("Failed to get OCI secret value: ", err.Error())
				return docker_secrets.Response{nil, err.Error(), true}
			}
			secretBundle = bundleResponse.SecretBundle
		} else {
			fmt.Println("No ocid or vault_ocid label found")
			return docker_secrets.Response{nil, "No ocid or vault_ocid label found", true}
		}
	}

	secretContent := secretBundle.SecretBundleContent.(oci_secrets.Base64SecretBundleContentDetails)

	value, err := base64.StdEncoding.DecodeString(*secretContent.Content)
	if (err != nil) {
		fmt.Println("Failed to decode secret value")
		return docker_secrets.Response{nil, "Failed to decode secret value", true}
	}

	return docker_secrets.Response{
		Value: value,
		DoNotReuse: true,
	}
}
