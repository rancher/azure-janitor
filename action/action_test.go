package action_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

	"github.com/rancher-sandbox/azure-janitor/action"
)

var _ = Describe("Test Azure Janitor", func() {

	Describe("Cleaning up a single resource group", func() {
		It("should delete the resource group", func() {
			fc := &fakeResourcesClient{
				existingGroups: []string{"e2e-1", "rg2"},
			}

			a, err := action.New(fc)
			Expect(err).NotTo(HaveOccurred())

			err = a.Cleanup(context.TODO(), "e2e*", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(fc.DeleteCalled()).To(BeTrue())
			rg := fc.ResourceGroups()
			Expect(rg).To(HaveLen(1))
			Expect(rg).Should(ContainElement("rg2"))
			Expect(rg).ShouldNot(ContainElement("e2e-1"))
		})
	})

	Describe("Cleaning up resource group with multiple matches", func() {
		It("should delete the matching resource groups", func() {
			fc := &fakeResourcesClient{
				existingGroups: []string{"e2e-1", "rg2", "e2e-2"},
			}

			a, err := action.New(fc)
			Expect(err).NotTo(HaveOccurred())

			err = a.Cleanup(context.TODO(), "e2e*", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(fc.DeleteCalled()).To(BeTrue())
			rg := fc.ResourceGroups()
			Expect(rg).To(HaveLen(1))
			Expect(rg).Should(ContainElement("rg2"))
			Expect(rg).ShouldNot(ContainElements("e2e-1", "e2e-2"))
		})
	})

	Describe("Cleaning up resource group with no matches", func() {
		It("should not delete the resource group", func() {
			fc := &fakeResourcesClient{
				existingGroups: []string{"rg1", "rg2"},
			}

			a, err := action.New(fc)
			Expect(err).NotTo(HaveOccurred())

			err = a.Cleanup(context.TODO(), "e2e*", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(fc.DeleteCalled()).To(BeFalse())
			rg := fc.ResourceGroups()
			Expect(rg).To(HaveLen(2))
			Expect(rg).Should(ContainElements("rg1", "rg2"))
		})
	})
})

type fakeResourcesClient struct {
	existingGroups []string
	deleteCalled   bool
}

func (f *fakeResourcesClient) ResourceGroups() []string {
	return f.existingGroups
}

func (f *fakeResourcesClient) DeleteCalled() bool {
	return f.deleteCalled
}

func (f *fakeResourcesClient) NewListPager(options *armresources.ResourceGroupsClientListOptions) *runtime.Pager[armresources.ResourceGroupsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[armresources.ResourceGroupsClientListResponse]{
		More: func(page armresources.ResourceGroupsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, rgclr *armresources.ResourceGroupsClientListResponse) (armresources.ResourceGroupsClientListResponse, error) {
			values := []*armresources.ResourceGroup{}
			for i, _ := range f.existingGroups {
				name := f.existingGroups[i]
				values = append(values, &armresources.ResourceGroup{
					Name: &name,
				})
			}

			return armresources.ResourceGroupsClientListResponse{
				ResourceGroupListResult: armresources.ResourceGroupListResult{
					Value: values,
				},
			}, nil
		},
	})
}

func (f *fakeResourcesClient) BeginDelete(ctx context.Context, resourceGroupName string, options *armresources.ResourceGroupsClientBeginDeleteOptions) (*runtime.Poller[armresources.ResourceGroupsClientDeleteResponse], error) {
	f.deleteCalled = true
	if len(f.existingGroups) == 0 {
		return nil, fmt.Errorf("cannot delete resource group %s as fake doesn't have it", resourceGroupName)
	}
	newGroups := []string{}
	for _, rg := range f.existingGroups {
		if rg != resourceGroupName {
			newGroups = append(newGroups, rg)
		}
	}
	f.existingGroups = newGroups

	resp := &http.Response{}
	resp.StatusCode = http.StatusNoContent
	resp.Request = &http.Request{
		Method: "DELETE",
	}
	resp.Body = io.NopCloser(strings.NewReader(""))

	pipeline := runtime.Pipeline{}

	return runtime.NewPoller[armresources.ResourceGroupsClientDeleteResponse](resp, pipeline, nil)
}
