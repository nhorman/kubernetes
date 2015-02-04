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
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// REST implements the RESTStorage interface
type REST struct {
	registry Registry
}

// NewREST returns a new REST.
func NewREST(registry Registry) *REST {
	return &REST{
		registry: registry,
	}
}

func (rs *REST) Create(ctx api.Context, obj runtime.Object) (<-chan apiserver.RESTResult, error) {
	req := obj.(*api.ExternalIPRequest)
	if !api.ValidNamespace(ctx, &req.ObjectMeta) {
		return nil, errors.NewConflict("externaliprequest", req.Namespace, fmt.Errorf("ExternalIPRequest.Namespace does not match the provided context"))
	}
	api.FillObjectMetaSystemFields(ctx, &req.ObjectMeta)
	if len(req.Name) == 0 {
		// TODO properly handle auto-generated names.
		// See https://github.com/GoogleCloudPlatform/kubernetes/issues/148 170 & 1135
		req.Name = string(req.UID)
	}

	return apiserver.MakeAsync(func() (runtime.Object, error) {
		if err := rs.registry.CreateExternalIPRequest(ctx, req); err != nil {
			return nil, err
		}
		return rs.registry.GetExternalIPRequest(ctx, req.Name)
	}), nil
}

func (rs *REST) Delete(ctx api.Context, id string) (<-chan apiserver.RESTResult, error) {
	return apiserver.MakeAsync(func() (runtime.Object, error) {
		_, found := api.NamespaceFrom(ctx)
		if !found {
			return &api.Status{Status: api.StatusFailure}, nil
		}

		return &api.Status{Status: api.StatusSuccess}, rs.registry.DeleteExternalIPRequest(ctx, id)
	}), nil
}

func (rs *REST) Get(ctx api.Context, id string) (runtime.Object, error) {
	req, err := rs.registry.GetExternalIPRequest(ctx, id)
	if err != nil {
		return req, err
	}
	if req == nil {
		return req, nil
	}

	return req, err
}

func ExternalIPRequestToSelectableFields(req *api.ExternalIPRequest) labels.Set {

	// TODO we are populating both Status and DesiredState because selectors are not aware of API versions
	// see https://github.com/GoogleCloudPlatform/kubernetes/pull/2503

	return labels.Set{
		"name":                req.Name,
	}
}

// filterFunc returns a predicate based on label & field selectors
func (rs *REST) filterFunc(label, field labels.Selector) func(*api.ExternalIPRequest) bool {
	return func(req *api.ExternalIPRequest) bool {
		fields := ExternalIPRequestToSelectableFields(req)
		return label.Matches(labels.Set(req.Labels)) && field.Matches(fields)
	}
}

func (rs *REST) List(ctx api.Context, label, field labels.Selector) (runtime.Object, error) {
	reqs, err := rs.registry.ListExternalIpRequestsPredicate(ctx, rs.filterFunc(label, field))
	return reqs, err
}

// Watch begins watching for new, changed, or deleted pods.
func (rs *REST) Watch(ctx api.Context, label, field labels.Selector, resourceVersion string) (watch.Interface, error) {
	// TODO: Add pod status to watch command
	return rs.registry.WatchExternalIPRequests(ctx, label, field, resourceVersion)
}

func (*REST) New() runtime.Object {
	return &api.ExternalIPRequest{}
}

func (*REST) NewList() runtime.Object {
	return &api.ExternalIPRequestList{}
}

func (rs *REST) Update(ctx api.Context, obj runtime.Object) (<-chan apiserver.RESTResult, error) {
	req := obj.(*api.ExternalIPRequest)
	if !api.ValidNamespace(ctx, &req.ObjectMeta) {
		return nil, errors.NewConflict("externaliprequest", req.Namespace, fmt.Errorf("ExternalIPRequest.Namespace does not match the provided context"))
	}
	return apiserver.MakeAsync(func() (runtime.Object, error) {
		if err := rs.registry.UpdateExternalIPRequest(ctx, req); err != nil {
			return nil, err
		}
		return rs.registry.GetExternalIPRequest(ctx, req.Name)
	}), nil
}
