//go:build test_all || test_unit

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"testing"

	dspav1alpha1 "github.com/opendatahub-io/data-science-pipelines-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func TestDeployVisualizationServer(t *testing.T) {
	testNamespace := "testnamespace"
	testDSPAName := "testdspa"
	expectedVisualizationServerName := "ds-pipeline-visualizationserver-testdspa"

	// Construct DSPASpec with deployed Visualization Server
	dspa := &dspav1alpha1.DataSciencePipelinesApplication{
		Spec: dspav1alpha1.DSPASpec{
			VisualizationServer: &dspav1alpha1.VisualizationServer{
				Deploy: true,
			},
			Database: &dspav1alpha1.Database{
				DisableHealthCheck: false,
				MariaDB: &dspav1alpha1.MariaDB{
					Deploy: true,
				},
			},
			ObjectStorage: &dspav1alpha1.ObjectStorage{
				DisableHealthCheck: false,
				Minio: &dspav1alpha1.Minio{
					Deploy: false,
					Image:  "someimage",
				},
			},
		},
	}

	// Enrich DSPA with name+namespace
	dspa.Namespace = testNamespace
	dspa.Name = testDSPAName

	// Create Context, Fake Controller and Params
	ctx, params, reconciler := CreateNewTestObjects()
	err := params.ExtractParams(ctx, dspa, reconciler.Client, reconciler.Log)
	assert.Nil(t, err)

	// Ensure Visualization Server Deployment doesn't yet exist
	deployment := &appsv1.Deployment{}
	created, err := reconciler.IsResourceCreated(ctx, deployment, expectedVisualizationServerName, testNamespace)
	assert.False(t, created)
	assert.Nil(t, err)

	// Run test reconciliation
	err = reconciler.ReconcileVisualizationServer(dspa, params)
	assert.Nil(t, err)

	// Ensure Visualization Server Deployment now exists
	deployment = &appsv1.Deployment{}
	created, err = reconciler.IsResourceCreated(ctx, deployment, expectedVisualizationServerName, testNamespace)
	assert.True(t, created)
	assert.Nil(t, err)

}

func TestDontDeployVisualizationServer(t *testing.T) {
	testNamespace := "testnamespace"
	testDSPAName := "testdspa"
	expectedVisualizationServerName := "ds-pipeline-visualization-server-testdspa"

	// Construct DSPASpec with non-deployed Visualization Server
	dspa := &dspav1alpha1.DataSciencePipelinesApplication{
		Spec: dspav1alpha1.DSPASpec{
			VisualizationServer: &dspav1alpha1.VisualizationServer{
				Deploy: false,
			},
		},
	}

	// Enrich DSPA with name+namespace
	dspa.Name = testDSPAName
	dspa.Namespace = testNamespace

	// Create Context, Fake Controller and Params
	ctx, params, reconciler := CreateNewTestObjects()

	// Ensure Visualization Server Deployment doesn't yet exist
	deployment := &appsv1.Deployment{}
	created, err := reconciler.IsResourceCreated(ctx, deployment, expectedVisualizationServerName, testNamespace)
	assert.False(t, created)
	assert.Nil(t, err)

	// Run test reconciliation
	err = reconciler.ReconcileVisualizationServer(dspa, params)
	assert.Nil(t, err)

	// Ensure Visualization Server Deployment still doesn't exist
	deployment = &appsv1.Deployment{}
	created, err = reconciler.IsResourceCreated(ctx, deployment, expectedVisualizationServerName, testNamespace)
	assert.False(t, created)
	assert.Nil(t, err)
}
