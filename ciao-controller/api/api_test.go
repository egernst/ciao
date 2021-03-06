// Copyright (c) 2016 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/01org/ciao/ciao-controller/types"
	"github.com/01org/ciao/payloads"
	"github.com/01org/ciao/service"
)

type test struct {
	method           string
	request          string
	requestBody      string
	media            string
	expectedStatus   int
	expectedResponse string
}

var tests = []test{
	{
		"GET",
		"/",
		"",
		"application/text",
		http.StatusOK,
		`[{"rel":"pools","href":"/pools","version":"x.ciao.pools.v1","minimum_version":"x.ciao.pools.v1"},{"rel":"external-ips","href":"/external-ips","version":"x.ciao.external-ips.v1","minimum_version":"x.ciao.external-ips.v1"},{"rel":"workloads","href":"/workloads","version":"x.ciao.workloads.v1","minimum_version":"x.ciao.workloads.v1"},{"rel":"tenants","href":"/tenants","version":"x.ciao.tenants.v1","minimum_version":"x.ciao.tenants.v1"}]`,
	},
	{
		"GET",
		"/pools",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusOK,
		`{"pools":[{"id":"ba58f471-0735-4773-9550-188e2d012941","name":"testpool","free":0,"total_ips":0,"links":[{"rel":"self","href":"/pools/ba58f471-0735-4773-9550-188e2d012941"}]}]}`,
	},
	{
		"GET",
		"/pools?name=testpool",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusOK,
		`{"pools":[{"id":"ba58f471-0735-4773-9550-188e2d012941","name":"testpool","free":0,"total_ips":0,"links":[{"rel":"self","href":"/pools/ba58f471-0735-4773-9550-188e2d012941"}]}]}`,
	},
	{
		"POST",
		"/pools",
		`{"name":"testpool"}`,
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"GET",
		"/pools/ba58f471-0735-4773-9550-188e2d012941",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusOK,
		`{"id":"ba58f471-0735-4773-9550-188e2d012941","name":"testpool","free":0,"total_ips":0,"links":[{"rel":"self","href":"/pools/ba58f471-0735-4773-9550-188e2d012941"}],"subnets":[],"ips":[]}`,
	},
	{
		"DELETE",
		"/pools/ba58f471-0735-4773-9550-188e2d012941",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"POST",
		"/pools/ba58f471-0735-4773-9550-188e2d012941",
		`{"subnet":"192.168.0.0/24"}`,
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"DELETE",
		"/pools/ba58f471-0735-4773-9550-188e2d012941/subnets/ba58f471-0735-4773-9550-188e2d012941",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"DELETE",
		"/pools/ba58f471-0735-4773-9550-188e2d012941/external-ips/ba58f471-0735-4773-9550-188e2d012941",
		"",
		fmt.Sprintf("application/%s", PoolsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"GET",
		"/external-ips",
		"",
		fmt.Sprintf("application/%s", ExternalIPsV1),
		http.StatusOK,
		`[{"mapping_id":"ba58f471-0735-4773-9550-188e2d012941","external_ip":"192.168.0.1","internal_ip":"172.16.0.1","instance_id":"","tenant_id":"8a497c68-a88a-4c1c-be56-12a4883208d3","pool_id":"f384ffd8-e7bd-40c2-8552-2efbe7e3ad6e","pool_name":"mypool","links":[{"rel":"self","href":"/external-ips/ba58f471-0735-4773-9550-188e2d012941"},{"rel":"pool","href":"/pools/f384ffd8-e7bd-40c2-8552-2efbe7e3ad6e"}]}]`,
	},
	{
		"POST",
		"/19df9b86-eda3-489d-b75f-d38710e210cb/external-ips",
		`{"pool_name":"apool","instance_id":"validinstanceID"}`,
		fmt.Sprintf("application/%s", ExternalIPsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"POST",
		"/workloads",
		`{"id":"","description":"testWorkload","fw_type":"legacy","vm_type":"qemu","image_name":"","config":"this will totally work!","defaults":[]}`,
		fmt.Sprintf("application/%s", WorkloadsV1),
		http.StatusCreated,
		`{"workload":{"id":"ba58f471-0735-4773-9550-188e2d012941","description":"testWorkload","fw_type":"legacy","vm_type":"qemu","image_name":"","config":"this will totally work!","defaults":[],"storage":null},"link":{"rel":"self","href":"/workloads/ba58f471-0735-4773-9550-188e2d012941"}}`,
	},
	{
		"DELETE",
		"/workloads/76f4fa99-e533-4cbd-ab36-f6c0f51292ed",
		"",
		fmt.Sprintf("application/%s", WorkloadsV1),
		http.StatusNoContent,
		"null",
	},
	{
		"GET",
		"/workloads/ba58f471-0735-4773-9550-188e2d012941",
		"",
		fmt.Sprintf("application/%s", WorkloadsV1),
		http.StatusOK,
		`{"id":"ba58f471-0735-4773-9550-188e2d012941","description":"testWorkload","fw_type":"legacy","vm_type":"qemu","image_name":"","config":"this will totally work!","defaults":null,"storage":null}`,
	},
	{
		"GET",
		"/tenants/093ae09b-f653-464e-9ae6-5ae28bd03a22/quotas",
		"",
		fmt.Sprintf("application/%s", TenantsV1),
		http.StatusOK,
		`{"quotas":[{"name":"test-quota-1","value":"10","usage":"3"},{"name":"test-quota-2","value":"unlimited","usage":"10"},{"name":"test-limit","value":"123"}]}`,
	},
}

type testCiaoService struct{}

func (ts testCiaoService) ListPools() ([]types.Pool, error) {
	self := types.Link{
		Rel:  "self",
		Href: "/pools/ba58f471-0735-4773-9550-188e2d012941",
	}

	resp := types.Pool{
		ID:       "ba58f471-0735-4773-9550-188e2d012941",
		Name:     "testpool",
		Free:     0,
		TotalIPs: 0,
		Subnets:  []types.ExternalSubnet{},
		IPs:      []types.ExternalIP{},
		Links:    []types.Link{self},
	}

	return []types.Pool{resp}, nil
}

func (ts testCiaoService) AddPool(name string, subnet *string, ips []string) (types.Pool, error) {
	return types.Pool{}, nil
}

func (ts testCiaoService) ShowPool(id string) (types.Pool, error) {
	fmt.Println("ShowPool")
	self := types.Link{
		Rel:  "self",
		Href: "/pools/ba58f471-0735-4773-9550-188e2d012941",
	}

	resp := types.Pool{
		ID:       "ba58f471-0735-4773-9550-188e2d012941",
		Name:     "testpool",
		Free:     0,
		TotalIPs: 0,
		Subnets:  []types.ExternalSubnet{},
		IPs:      []types.ExternalIP{},
		Links:    []types.Link{self},
	}

	return resp, nil
}

func (ts testCiaoService) DeletePool(id string) error {
	return nil
}

func (ts testCiaoService) AddAddress(poolID string, subnet *string, ips []string) error {
	return nil
}

func (ts testCiaoService) RemoveAddress(poolID string, subnet *string, extIP *string) error {
	return nil
}

func (ts testCiaoService) ListMappedAddresses(tenant *string) []types.MappedIP {
	var ref string

	m := types.MappedIP{
		ID:         "ba58f471-0735-4773-9550-188e2d012941",
		ExternalIP: "192.168.0.1",
		InternalIP: "172.16.0.1",
		TenantID:   "8a497c68-a88a-4c1c-be56-12a4883208d3",
		PoolID:     "f384ffd8-e7bd-40c2-8552-2efbe7e3ad6e",
		PoolName:   "mypool",
	}

	if tenant != nil {
		ref = fmt.Sprintf("%s/external-ips/%s", *tenant, m.ID)
	} else {
		ref = fmt.Sprintf("/external-ips/%s", m.ID)
	}

	link := types.Link{
		Rel:  "self",
		Href: ref,
	}

	m.Links = []types.Link{link}

	if tenant == nil {
		ref := fmt.Sprintf("/pools/%s", m.PoolID)

		link := types.Link{
			Rel:  "pool",
			Href: ref,
		}

		m.Links = append(m.Links, link)
	}

	return []types.MappedIP{m}
}

func (ts testCiaoService) MapAddress(tenantID string, name *string, instanceID string) error {
	return nil
}

func (ts testCiaoService) UnMapAddress(string) error {
	return nil
}

func (ts testCiaoService) CreateWorkload(req types.Workload) (types.Workload, error) {
	req.ID = "ba58f471-0735-4773-9550-188e2d012941"
	return req, nil
}

func (ts testCiaoService) DeleteWorkload(tenant string, workload string) error {
	return nil
}

func (ts testCiaoService) ShowWorkload(tenant string, ID string) (types.Workload, error) {
	return types.Workload{
		ID:          "ba58f471-0735-4773-9550-188e2d012941",
		TenantID:    tenant,
		Description: "testWorkload",
		FWType:      payloads.Legacy,
		VMType:      payloads.QEMU,
		Config:      "this will totally work!",
	}, nil
}

func (ts testCiaoService) ListQuotas(tenantID string) []types.QuotaDetails {
	return []types.QuotaDetails{
		{Name: "test-quota-1", Value: 10, Usage: 3},
		{Name: "test-quota-2", Value: -1, Usage: 10},
		{Name: "test-limit", Value: 123, Usage: 0},
	}
}

func (ts testCiaoService) UpdateQuotas(tenantID string, qds []types.QuotaDetails) error {
	return nil
}

func TestResponse(t *testing.T) {
	var ts testCiaoService

	mux := Routes(Config{"", ts}, nil)

	for i, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.request, bytes.NewBuffer([]byte(tt.requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(service.SetPrivilege(req.Context(), true))

		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", tt.media)

		mux.ServeHTTP(rr, req)

		status := rr.Code
		if status != tt.expectedStatus {
			t.Errorf("test %d: got %v, expected %v", i, status, tt.expectedStatus)
		}

		if rr.Body.String() != tt.expectedResponse {
			t.Errorf("test %d: %s: failed\ngot: %v\nexp: %v", i, tt.request, rr.Body.String(), tt.expectedResponse)
		}
	}
}

func TestRoutes(t *testing.T) {
	var ts testCiaoService
	config := Config{"", ts}

	r := Routes(config, nil)
	if r == nil {
		t.Fatalf("No routes returned")
	}
}
