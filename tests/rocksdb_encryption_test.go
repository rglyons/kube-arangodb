package tests

import (
	"context"
	"crypto/rand"
	"strings"
	"testing"

	"github.com/dchest/uniuri"

	api "github.com/arangodb/k8s-operator/pkg/apis/arangodb/v1alpha"
	"github.com/arangodb/k8s-operator/pkg/client"
	"github.com/arangodb/k8s-operator/pkg/util/k8sutil"
)

// TestRocksDBEncryptionSingle tests the creating of a single server deployment
// with RocksDB & Encryption.
func TestRocksDBEncryptionSingle(t *testing.T) {
	longOrSkip(t)
	image := getEnterpriseImageOrSkip(t)
	c := client.MustNewInCluster()
	kubecli := mustNewKubeClient(t)
	ns := getNamespace(t)

	// Prepare deployment config
	depl := newDeployment("test-rocksdb-enc-sng-" + uniuri.NewLen(4))
	depl.Spec.Mode = api.DeploymentModeSingle
	depl.Spec.Image = image
	depl.Spec.StorageEngine = api.StorageEngineRocksDB
	depl.Spec.RocksDB.Encryption.KeySecretName = strings.ToLower(uniuri.New())

	// Create encryption key secret
	key := make([]byte, 32)
	rand.Read(key)
	if err := k8sutil.CreateEncryptionKeySecret(kubecli, depl.Spec.RocksDB.Encryption.KeySecretName, ns, key); err != nil {
		t.Fatalf("Create encryption key secret failed: %v", err)
	}

	// Create deployment
	_, err := c.DatabaseV1alpha().ArangoDeployments(ns).Create(depl)
	if err != nil {
		t.Fatalf("Create deployment failed: %v", err)
	}

	// Wait for deployment to be ready
	apiObject, err := waitUntilDeployment(c, depl.GetName(), ns, deploymentHasState(api.DeploymentStateRunning))
	if err != nil {
		t.Errorf("Deployment not running in time: %v", err)
	}

	// Create database client
	ctx := context.Background()
	client := mustNewArangodDatabaseClient(ctx, kubecli, apiObject, t)

	// Wait for single server available
	if err := waitUntilVersionUp(client); err != nil {
		t.Fatalf("Single server not running returning version in time: %v", err)
	}

	// Cleanup
	removeDeployment(c, depl.GetName(), ns)
	removeSecret(kubecli, depl.Spec.RocksDB.Encryption.KeySecretName, ns)
}
