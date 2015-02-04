/*
Copyright 2014 Google Inc. All rights reserved.

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

package externaliprequest 

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// Registry is an interface implemented by things that know how to store Pod objects.
type Registry interface {
	// ListExternalIPRequests obtains a list of external ip requests having labels which match selector.
	ListExternalIpRequests(ctx api.Context, selector labels.Selector) (*api.ExternalIPRequestList, error)
	ListExternalIpRequestsPredicate(ctx api.Context, filter func(*api.ExternalIPRequest) bool) (*api.ExternalIPRequestList, error)
	// Watch for new/changed/deleted external ip requests 
	WatchExternalIPRequests(ctx api.Context, label, field labels.Selector, resourceVersion string) (watch.Interface, error)
	// Get a specific external ip request 
	GetExternalIPRequest(ctx api.Context, request string) (*api.ExternalIPRequest, error)
	// Create a external ip request based on a specification.
	CreateExternalIPRequest(ctx api.Context, request *api.ExternalIPRequest) error
	// Update an existing request 
	UpdateExternalIPRequest(ctx api.Context, request *api.ExternalIPRequest) error
	// Delete an existing request 
	DeleteExternalIPRequest(ctx api.Context, request string) error
}
