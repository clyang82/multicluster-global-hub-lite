package globalhubcontroller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	policyv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type placementBindingController struct {
	client dynamic.Interface
	gvr    schema.GroupVersionResource
}

func NewPlacementBindingController(dynamicClient dynamic.Interface) IController {
	return &placementBindingController{
		client: dynamicClient,
		gvr:    policyv1.SchemeGroupVersion.WithResource("placementbindings"),
	}
}

func (c *placementBindingController) GetName() string {
	return "placementbinding-controller"
}

func (c *placementBindingController) GetGVR() schema.GroupVersionResource {
	return c.gvr
}

func (c *placementBindingController) CreateInstanceFunc() func() client.Object {
	return func() client.Object {
		return &policyv1.PlacementBinding{}
	}
}

func (c *placementBindingController) ReconcileFunc() func(stopCh <-chan struct{}, obj interface{}) error {
	return func(stopCh <-chan struct{}, obj interface{}) error {
		unObj, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return fmt.Errorf("cann't convert obj(%+v) to *unstructured.Unstructured", obj)
		}

		labels := unObj.GetLabels()
		_, ok = labels[GlobalHubPolicyNamespaceLabel]
		if !ok {
			if labels == nil {
				labels = map[string]string{}
			}
			labels[GlobalHubPolicyNamespaceLabel] = unObj.GetNamespace()
			unObj.SetLabels(labels)
			if _, err := c.client.Resource(c.gvr).Namespace(unObj.GetNamespace()).Update(context.TODO(), unObj, metav1.UpdateOptions{}); err != nil {
				return err
			}
			return nil
		}
		return nil
	}
}
