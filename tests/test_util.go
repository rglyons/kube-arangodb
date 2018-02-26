//
// DISCLAIMER
//
// Copyright 2018 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	driver "github.com/arangodb/go-driver"
	api "github.com/arangodb/k8s-operator/pkg/apis/arangodb/v1alpha"
	"github.com/arangodb/k8s-operator/pkg/generated/clientset/versioned"
	"github.com/arangodb/k8s-operator/pkg/util/arangod"
	"github.com/arangodb/k8s-operator/pkg/util/k8sutil"
	"github.com/arangodb/k8s-operator/pkg/util/retry"
)

const (
	deploymentReadyTimeout = time.Minute * 2
)

var (
	maskAny = errors.WithStack
)

// mustNewKubeClient creates a kubernetes client
// failing the test on errors.
func mustNewKubeClient(t *testing.T) kubernetes.Interface {
	c, err := k8sutil.NewKubeClient()
	if err != nil {
		t.Fatalf("Failed to create kube cli: %v", err)
	}
	return c
}

// mustNewArangodDatabaseClient creates a new database client,
// failing the test on errors.
func mustNewArangodDatabaseClient(kubecli kubernetes.Interface, apiObject *api.ArangoDeployment, t *testing.T) driver.Client {
	c, err := arangod.CreateArangodDatabaseClient(kubecli, apiObject)
	if err != nil {
		t.Fatalf("Failed to create arango database client: %v", err)
	}
	return c
}

// getNamespace returns the kubernetes namespace in which to run tests.
func getNamespace(t *testing.T) string {
	ns := os.Getenv("TEST_NAMESPACE")
	if ns == "" {
		t.Fatal("Missing environment variable TEST_NAMESPACE")
	}
	return ns
}

// newDeployment creates a basic ArangoDeployment with configured
// type & name.
func newDeployment(name string) *api.ArangoDeployment {
	return &api.ArangoDeployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.SchemeGroupVersion.String(),
			Kind:       api.ArangoDeploymentResourceKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(name),
		},
	}
}

// waitUntilDeployment waits until a deployment with given name in given namespace
// reached a state where the given predicate returns true.
func waitUntilDeployment(cli versioned.Interface, deploymentName, ns string, predicate func(*api.ArangoDeployment) error) (*api.ArangoDeployment, error) {
	var result *api.ArangoDeployment
	op := func() error {
		obj, err := cli.DatabaseV1alpha().ArangoDeployments(ns).Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			result = nil
			return maskAny(err)
		}
		result = obj
		if err := predicate(obj); err != nil {
			return maskAny(err)
		}
		return nil
	}
	if err := retry.Retry(op, deploymentReadyTimeout); err != nil {
		return nil, maskAny(err)
	}
	return result, nil
}

// waitUntilClusterHealth waits until an arango cluster
// reached a state where the given predicate returns nil.
func waitUntilClusterHealth(cli driver.Client, predicate func(driver.ClusterHealth) error) error {
	ctx := context.Background()
	op := func() error {
		cluster, err := cli.Cluster(ctx)
		if err != nil {
			return maskAny(err)
		}
		h, err := cluster.Health(ctx)
		if err != nil {
			return maskAny(err)
		}
		if err := predicate(h); err != nil {
			return maskAny(err)
		}
		return nil
	}
	if err := retry.Retry(op, deploymentReadyTimeout); err != nil {
		return maskAny(err)
	}
	return nil
}

// clusterHealthEqualsSpec returns nil when the given health matches
// with the given deployment spec.
func clusterHealthEqualsSpec(h driver.ClusterHealth, spec api.DeploymentSpec) error {
	agents := 0
	goodDBServers := 0
	goodCoordinators := 0
	for _, s := range h.Health {
		if s.Role == driver.ServerRoleAgent {
			agents++
		} else if s.Status == driver.ServerStatusGood {
			switch s.Role {
			case driver.ServerRoleDBServer:
				goodDBServers++
			case driver.ServerRoleCoordinator:
				goodCoordinators++
			}
		}
	}
	if spec.Agents.Count == agents &&
		spec.DBServers.Count == goodDBServers &&
		spec.Coordinators.Count == goodCoordinators {
		return nil
	}
	return fmt.Errorf("Expected %d,%d,%d got %d,%d,%d",
		spec.Agents.Count, spec.DBServers.Count, spec.Coordinators.Count,
		agents, goodDBServers, goodCoordinators,
	)
}

// updateDeployment updates a deployment
func updateDeployment(cli versioned.Interface, deploymentName, ns string, update func(*api.DeploymentSpec)) (*api.ArangoDeployment, error) {
	for {
		// Get current version
		current, err := cli.Database().ArangoDeployments(ns).Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			return nil, maskAny(err)
		}
		update(&current.Spec)
		current, err = cli.Database().ArangoDeployments(ns).Update(current)
		if k8sutil.IsConflict(err) {
			// Retry
		} else if err != nil {
			return nil, maskAny(err)
		}
		return current, nil
	}
}

// removeDeployment removes a deployment
func removeDeployment(cli versioned.Interface, deploymentName, ns string) error {
	if err := cli.Database().ArangoDeployments(ns).Delete(deploymentName, nil); err != nil && k8sutil.IsNotFound(err) {
		return maskAny(err)
	}
	return nil
}