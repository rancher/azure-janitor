package action

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func (a *action) cleanupResourceGroup(ctx context.Context, rg *armresources.ResourceGroup, commit bool) error {
	Log("Cleaning up resource group %s", *rg.Name)

	if !commit {
		LogDebug("Dry run (i.e. commit=false), skipping deletion of resource group %s", *rg.Name)
		return nil
	}

	pollerResp, err := a.client.BeginDelete(ctx, *rg.Name, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deletion of resource group: %w", err)
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete resource group: %w", err)
	}

	Log("Succesfully deleted resource group %s", *rg.Name)

	return nil
}
