package main

import (
	"context"
	"fmt"
	"encoding/base64"

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
	ocid := req.SecretLabels["ocid"]

	if (ocid == "") {
		fmt.Println("No ocid label found")
		return docker_secrets.Response{nil, "No ocid label found", true}
	}

	ctx := context.Background()

	bundleRequest := oci_secrets.GetSecretBundleRequest{}
	bundleRequest.SecretId = &ocid
	bundleResponse, err := d.client.GetSecretBundle(ctx, bundleRequest)
	
	if (err != nil) {
		fmt.Println("Failed to get OCI secret value: ", err.Error())
		return docker_secrets.Response{nil, err.Error(), true}
	}

	secretContent := bundleResponse.SecretBundle.SecretBundleContent.(oci_secrets.Base64SecretBundleContentDetails)

	value, err := base64.StdEncoding.DecodeString(*secretContent.Content)

	resp := docker_secrets.Response{}
	resp.DoNotReuse = true
	resp.Value = value

	return resp
}
