/*
Copyright 2016 The Kubernetes Authors.

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

package node

import (
	"fmt"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//clientset "k8s.io/client-go/kubernetes"
	fakeexternal "k8s.io/client-go/kubernetes/fake"
	//v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

func TestGetPreferredAddress(t *testing.T) {
	testcases := map[string]struct {
		Labels      map[string]string
		Addresses   []v1.NodeAddress
		Preferences []v1.NodeAddressType

		ExpectErr     string
		ExpectAddress string
	}{
		"no addresses": {
			ExpectErr: "no preferred addresses found; known addresses: []",
		},
		"missing address": {
			Addresses: []v1.NodeAddress{
				{Type: v1.NodeInternalIP, Address: "1.2.3.4"},
			},
			Preferences: []v1.NodeAddressType{v1.NodeHostName},
			ExpectErr:   "no preferred addresses found; known addresses: [{InternalIP 1.2.3.4}]",
		},
		"found address": {
			Addresses: []v1.NodeAddress{
				{Type: v1.NodeInternalIP, Address: "1.2.3.4"},
				{Type: v1.NodeExternalIP, Address: "1.2.3.5"},
				{Type: v1.NodeExternalIP, Address: "1.2.3.7"},
			},
			Preferences:   []v1.NodeAddressType{v1.NodeHostName, v1.NodeExternalIP},
			ExpectAddress: "1.2.3.5",
		},
		"found hostname address": {
			Labels: map[string]string{kubeletapis.LabelHostname: "label-hostname"},
			Addresses: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.2.3.5"},
				{Type: v1.NodeHostName, Address: "status-hostname"},
			},
			Preferences:   []v1.NodeAddressType{v1.NodeHostName, v1.NodeExternalIP},
			ExpectAddress: "status-hostname",
		},
		"found label address": {
			Labels: map[string]string{kubeletapis.LabelHostname: "label-hostname"},
			Addresses: []v1.NodeAddress{
				{Type: v1.NodeExternalIP, Address: "1.2.3.5"},
			},
			Preferences:   []v1.NodeAddressType{v1.NodeHostName, v1.NodeExternalIP},
			ExpectAddress: "label-hostname",
		},
	}

	for k, tc := range testcases {
		node := &v1.Node{
			ObjectMeta: metav1.ObjectMeta{Labels: tc.Labels},
			Status:     v1.NodeStatus{Addresses: tc.Addresses},
		}
		address, err := GetPreferredNodeAddress(node, tc.Preferences)
		errString := ""
		if err != nil {
			errString = err.Error()
		}
		if errString != tc.ExpectErr {
			t.Errorf("%s: expected err=%q, got %q", k, tc.ExpectErr, errString)
		}
		if address != tc.ExpectAddress {
			t.Errorf("%s: expected address=%q, got %q", k, tc.ExpectAddress, address)
		}
	}
}

func TestPatchNodeStatus(t *testing.T) {

	fmt.Println("PHIL start test")
	testcases := map[string]struct {
		oldNode *v1.Node
		newNode *v1.Node

		ExpectErr     string
		ExpectAddress string
	}{
		"one address": {
			oldNode: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "testKubeletHostname"},
				Spec:       v1.NodeSpec{},
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{Type: v1.NodeInternalIP, Address: "1.2.3.1"},
						{Type: v1.NodeInternalIP, Address: "1.2.3.5"},
						{Type: v1.NodeInternalIP, Address: "1.2.3.6"},
						{Type: v1.NodeHostName, Address: "testKubeletHostname"},
					},
				},
			},
			newNode: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{Name: "testKubeletHostname"},
				Spec:       v1.NodeSpec{},
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{Type: v1.NodeInternalIP, Address: "1.2.3.1"},
						{Type: v1.NodeInternalIP, Address: "1.2.3.5"},
						{Type: v1.NodeInternalIP, Address: "1.2.3.6"},
						{Type: v1.NodeInternalIP, Address: "1.2.3.7"},
						{Type: v1.NodeHostName, Address: "testKubeletHostname"},
					},
				},
			},

			ExpectErr: "PHIL no preferred addresses found; known addresses: []",
		},
	}

	for k, tc := range testcases {
		fakeClientset := fakeexternal.NewSimpleClientset()
		nodeName := types.NodeName("PHILnode")
		fmt.Printf("PHIL - nodeName %s\n\n", string(nodeName))
		fmt.Printf("PHIL - oldNode %v\n\n", tc.oldNode)
		fmt.Printf("PHIL - node %v\n\n", tc.newNode)
		patchedNode, patch, err := PatchNodeStatus(fakeClientset.CoreV1(), nodeName, tc.oldNode, tc.newNode)
		fmt.Printf("PHIL - patch %v\n\n", string(patch))
		fmt.Printf("PHIL - patchedNode, %v\n\n", patchedNode)
		errString := ""
		if err != nil {
			errString = err.Error()
		}
		if errString != tc.ExpectErr {
			t.Errorf("%s: expected err=%q, got %q", k, tc.ExpectErr, errString)
		}
		if patchedNode != nil {
			t.Errorf("hoho")
		}
		if patch != nil {
			t.Errorf("haha")
		}
		t.Errorf("PHIL-----")
		//if address != tc.ExpectAddress {
		//	t.Errorf("%s: expected address=%q, got %q", k, tc.ExpectAddress, address)
		//}
	}
	t.Errorf("PHIL  -----")
}

// func GetHostname(hostnameOverride string) string {
// func GetNodeHostIP(node *v1.Node) (net.IP, error) {
// func InternalGetNodeHostIP(node *api.Node) (net.IP, error) {
// func PatchNodeStatus(c v1core.CoreV1Interface, nodeName types.NodeName, oldNode *v1.Node, newNode *v1.Node) (*v1.Node, []byte, error) {
// func SetNodeCondition(c clientset.Interface, node types.NodeName, condition v1.NodeCondition) error {
// func PatchNodeCIDR(c clientset.Interface, node types.NodeName, cidr string) error {
// func PatchNodeStatus(c v1core.CoreV1Interface, nodeName types.NodeName, oldNode *v1.Node, newNode *v1.Node) (*v1.Node, []byte, error) {
// func preparePatchBytesforNodeStatus(nodeName types.NodeName, oldNode *v1.Node, newNode *v1.Node) ([]byte, error) {
