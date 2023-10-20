package action

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type AzureJanitorAction interface {
	Cleanup(ctx context.Context, resourceGroupPattern string, commit bool) error
}

func New(client *armresources.ResourceGroupsClient) (AzureJanitorAction, error) {
	return &action{
		client: client,
	}, nil
}

type action struct {
	client *armresources.ResourceGroupsClient
}

func (a *action) Cleanup(ctx context.Context, resourceGroupPattern string, commit bool) error {
	Log("cleaning up resource groups with pattern %s (commit=%t)", resourceGroupPattern, commit)

	LogDebug("Getting list of Azure resource groups")
	pager := a.client.NewListPager(nil)

	resourceGroups := make([]*armresources.ResourceGroup, 0)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("getting resource groups page: %w", err)
		}

		resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
	}

	for _, resGroup := range resourceGroups {
		match, err := regexp.MatchString(resourceGroupPattern, *resGroup.Name)
		if err != nil {
			return fmt.Errorf("failed matching %s using pattern %s: %w", *resGroup.Name, resourceGroupPattern, err)
		}
		if !match {
			LogDebug("Resource group %s doesn't match pattern %s, skipping cleanup", *resGroup.Name, resourceGroupPattern)
			continue
		}

		if err := a.cleanupResourceGroup(ctx, resGroup, commit); err != nil {
			return fmt.Errorf("cleaning up resource group %s: %w", *resGroup.Name, err)
		}
	}

	return nil
}
