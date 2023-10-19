package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/rancher-sandbox/azure-janitor/action"
)

func main() {
	action.Log("running azure janitor")

	input, err := action.NewInput()
	if err != nil {
		action.LogErrorAndExit("failed to get input: %s", err.Error())
	}
	if err := input.Validate(); err != nil {
		action.LogErrorAndExit("failed input validation: %s", err.Error())
	}

	credential, err := azidentity.NewClientSecretCredential(input.AzureTenantID, input.AzureClientID, input.AzureClientSecret, nil)
	if err != nil {
		action.LogErrorAndExit("failed getting credentials: %s", err.Error())
	}

	client, err := armresources.NewResourceGroupsClient(input.AzureSubscriptionID, credential, nil)
	if err != nil {
		action.LogErrorAndExit("failed getting arm client: %s", err.Error())
	}

	a, err := action.New(client)
	if err != nil {
		action.LogErrorAndExit("failed to create action: %s", err.Error())
	}

	ctx := context.Background()
	if err := a.Cleanup(ctx, input.ResourceGroupPattern, input.Commit); err != nil {
		action.LogErrorAndExit("failed to cleanup resource groups: %s", err.Error())
	}
}
