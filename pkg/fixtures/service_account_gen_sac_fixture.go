// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package fixtures

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac/testconsts"
	"github.com/stackrox/rox/pkg/uuid"
)

// *storage.ServiceAccount represents a generic type that we use in the function below.

// GetSACTestStorageServiceAccountSet returns a set of mock *storage.ServiceAccount that can be used
// for scoped access control sets.
// It will include:
// 9 *storage.ServiceAccount scoped to Cluster1, 3 to each Namespace A / B / C.
// 9 *storage.ServiceAccount scoped to Cluster2, 3 to each Namespace A / B / C.
// 9 *storage.ServiceAccount scoped to Cluster3, 3 to each Namespace A / B / C.
func GetSACTestStorageServiceAccountSet(scopedStorageServiceAccountCreator func(id string, clusterID string, namespace string) *storage.ServiceAccount) []*storage.ServiceAccount {
	clusters := []string{testconsts.Cluster1, testconsts.Cluster2, testconsts.Cluster3}
	namespaces := []string{testconsts.NamespaceA, testconsts.NamespaceB, testconsts.NamespaceC}
	const numberOfAccounts = 3
	storageServiceAccounts := make([]*storage.ServiceAccount, 0, len(clusters)*len(namespaces)*numberOfAccounts)
	for _, cluster := range clusters {
		for _, namespace := range namespaces {
			for i := 0; i < numberOfAccounts; i++ {
				storageServiceAccounts = append(storageServiceAccounts, scopedStorageServiceAccountCreator(uuid.NewV4().String(), cluster, namespace))
			}
		}
	}
	return storageServiceAccounts
}