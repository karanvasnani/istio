// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/proto"

	authn "istio.io/api/authentication/v1alpha1"
	rbacproto "istio.io/api/rbac/v1alpha1"
	routing "istio.io/api/routing/v1alpha1"
	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/pkg/model"
	srmemory "istio.io/istio/pilot/pkg/serviceregistry/memory"
	mock_config "istio.io/istio/pilot/test/mock"
)

// getByMessageName finds a schema by message name if it is available
// In test setup, we do not have more than one descriptor with the same message type, so this
// function is ok for testing purpose.
func getByMessageName(descriptor model.ConfigDescriptor, name string) (model.ProtoSchema, bool) {
	for _, schema := range descriptor {
		if schema.MessageName == name {
			return schema, true
		}
	}
	return model.ProtoSchema{}, false
}

func TestConfigDescriptor(t *testing.T) {
	a := model.ProtoSchema{Type: "a", MessageName: "proxy.A"}
	descriptor := model.ConfigDescriptor{
		a,
		model.ProtoSchema{Type: "b", MessageName: "proxy.B"},
		model.ProtoSchema{Type: "c", MessageName: "proxy.C"},
	}
	want := []string{"a", "b", "c"}
	got := descriptor.Types()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("descriptor.Types() => got %+vwant %+v", spew.Sdump(got), spew.Sdump(want))
	}

	aType, aExists := descriptor.GetByType(a.Type)
	if !aExists || !reflect.DeepEqual(aType, a) {
		t.Errorf("descriptor.GetByType(a) => got %+v, want %+v", aType, a)
	}
	if _, exists := descriptor.GetByType("missing"); exists {
		t.Error("descriptor.GetByType(missing) => got true, want false")
	}

	aSchema, aSchemaExists := getByMessageName(descriptor, a.MessageName)
	if !aSchemaExists || !reflect.DeepEqual(aSchema, a) {
		t.Errorf("descriptor.GetByMessageName(a) => got %+v, want %+v", aType, a)
	}
	_, aSchemaNotExist := getByMessageName(descriptor, "blah")
	if aSchemaNotExist {
		t.Errorf("descriptor.GetByMessageName(blah) => got true, want false")
	}
}

func TestEventString(t *testing.T) {
	cases := []struct {
		in   model.Event
		want string
	}{
		{model.EventAdd, "add"},
		{model.EventUpdate, "update"},
		{model.EventDelete, "delete"},
	}
	for _, c := range cases {
		if got := c.in.String(); got != c.want {
			t.Errorf("Failed: got %q want %q", got, c.want)
		}
	}
}

func TestPortList(t *testing.T) {
	pl := model.PortList{
		{Name: "http", Port: 80, Protocol: model.ProtocolHTTP},
		{Name: "http-alt", Port: 8080, Protocol: model.ProtocolHTTP},
	}

	gotNames := pl.GetNames()
	wantNames := []string{"http", "http-alt"}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Errorf("GetNames() failed: got %v want %v", gotNames, wantNames)
	}

	cases := []struct {
		name  string
		port  *model.Port
		found bool
	}{
		{name: pl[0].Name, port: pl[0], found: true},
		{name: "foobar", found: false},
	}

	for _, c := range cases {
		gotPort, gotFound := pl.Get(c.name)
		if c.found != gotFound || !reflect.DeepEqual(gotPort, c.port) {
			t.Errorf("Get() failed: gotFound=%v wantFound=%v\ngot %+vwant %+v",
				gotFound, c.found, spew.Sdump(gotPort), spew.Sdump(c.port))
		}
	}
}

func TestServiceKey(t *testing.T) {
	svc := &model.Service{Hostname: "hostname"}

	// Verify Service.Key() delegates to ServiceKey()
	{
		want := "hostname|http|a=b,c=d"
		port := &model.Port{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}
		labels := model.Labels{"a": "b", "c": "d"}
		got := svc.Key(port, labels)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Service.Key() failed: got %v want %v", got, want)
		}
	}

	cases := []struct {
		port   model.PortList
		labels model.LabelsCollection
		want   string
	}{
		{
			port: model.PortList{
				{Name: "http", Port: 80, Protocol: model.ProtocolHTTP},
				{Name: "http-alt", Port: 8080, Protocol: model.ProtocolHTTP},
			},
			labels: model.LabelsCollection{{"a": "b", "c": "d"}},
			want:   "hostname|http,http-alt|a=b,c=d",
		},
		{
			port:   model.PortList{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}},
			labels: model.LabelsCollection{{"a": "b", "c": "d"}},
			want:   "hostname|http|a=b,c=d",
		},
		{
			port:   model.PortList{{Port: 80, Protocol: model.ProtocolHTTP}},
			labels: model.LabelsCollection{{"a": "b", "c": "d"}},
			want:   "hostname||a=b,c=d",
		},
		{
			port:   model.PortList{},
			labels: model.LabelsCollection{{"a": "b", "c": "d"}},
			want:   "hostname||a=b,c=d",
		},
		{
			port:   model.PortList{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}},
			labels: model.LabelsCollection{nil},
			want:   "hostname|http",
		},
		{
			port:   model.PortList{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}},
			labels: model.LabelsCollection{},
			want:   "hostname|http",
		},
		{
			port:   model.PortList{},
			labels: model.LabelsCollection{},
			want:   "hostname",
		},
	}

	for _, c := range cases {
		got := model.ServiceKey(svc.Hostname, c.port, c.labels)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Failed: got %q want %q", got, c.want)
		}
	}
}

func TestSubsetKey(t *testing.T) {
	hostname := model.Hostname("hostname")
	cases := []struct {
		hostname model.Hostname
		subset   string
		port     int
		want     string
	}{
		{
			hostname: "hostname",
			subset:   "subset",
			port:     80,
			want:     "outbound|80|subset|hostname",
		},
		{
			hostname: "hostname",
			subset:   "",
			port:     80,
			want:     "outbound|80||hostname",
		},
	}

	for _, c := range cases {
		got := model.BuildSubsetKey(model.TrafficDirectionOutbound, c.subset, hostname, c.port)
		if got != c.want {
			t.Errorf("Failed: got %q want %q", got, c.want)
		}

		// test parse subset key. ParseSubsetKey is the inverse of BuildSubsetKey
		_, s, h, p := model.ParseSubsetKey(got)
		if s != c.subset || h != c.hostname || p != c.port {
			t.Errorf("Failed: got %s,%s,%d want %s,%s,%d", s, h, p, c.subset, c.hostname, c.port)
		}
	}
}

func TestLabelsEquals(t *testing.T) {
	cases := []struct {
		a, b model.Labels
		want bool
	}{
		{
			a: nil,
			b: model.Labels{"a": "b"},
		},
		{
			a: model.Labels{"a": "b"},
			b: nil,
		},
		{
			a:    model.Labels{"a": "b"},
			b:    model.Labels{"a": "b"},
			want: true,
		},
	}
	for _, c := range cases {
		if got := c.a.Equals(c.b); got != c.want {
			t.Errorf("Failed: got eq=%v want=%v for %q ?= %q", got, c.want, c.a, c.b)
		}
	}
}

func TestConfigKey(t *testing.T) {
	config := mock_config.Make("ns", 2)
	want := "mock-config/ns/mock-config2"
	if key := config.ConfigMeta.Key(); key != want {
		t.Errorf("config.Key() => got %q, want %q", key, want)
	}
}

func TestResolveHostname(t *testing.T) {
	cases := []struct {
		meta model.ConfigMeta
		svc  *routing.IstioService
		want model.Hostname
	}{
		{
			meta: model.ConfigMeta{Namespace: "default", Domain: "cluster.local"},
			svc:  &routing.IstioService{Name: "hello"},
			want: "hello.default.svc.cluster.local",
		},
		{
			meta: model.ConfigMeta{Namespace: "foo", Domain: "foo"},
			svc: &routing.IstioService{Name: "hello",
				Namespace: "default", Domain: "svc.cluster.local"},
			want: "hello.default.svc.cluster.local",
		},
		{
			meta: model.ConfigMeta{},
			svc:  &routing.IstioService{Name: "hello"},
			want: "hello",
		},
		{
			meta: model.ConfigMeta{Namespace: "default"},
			svc:  &routing.IstioService{Name: "hello"},
			want: "hello.default",
		},
		{
			meta: model.ConfigMeta{Namespace: "default", Domain: "cluster.local"},
			svc:  &routing.IstioService{Service: "reviews.service.consul"},
			want: "reviews.service.consul",
		},
		{
			meta: model.ConfigMeta{Namespace: "foo", Domain: "foo"},
			svc: &routing.IstioService{Name: "hello", Service: "reviews.service.consul",
				Namespace: "default", Domain: "svc.cluster.local"},
			want: "reviews.service.consul",
		},
		{
			meta: model.ConfigMeta{Namespace: "default", Domain: "cluster.local"},
			svc:  &routing.IstioService{Service: "*cnn.com"},
			want: "*cnn.com",
		},
		{
			meta: model.ConfigMeta{Namespace: "foo", Domain: "foo"},
			svc: &routing.IstioService{Name: "hello", Service: "*cnn.com",
				Namespace: "default", Domain: "svc.cluster.local"},
			want: "*cnn.com",
		},
	}

	for _, test := range cases {
		if got := model.ResolveHostname(test.meta, test.svc); got != test.want {
			t.Errorf("ResolveHostname(%v, %v) => got %q, want %q", test.meta, test.svc, got, test.want)
		}
	}
}

func TestMatchSource(t *testing.T) {
	cases := []struct {
		meta      model.ConfigMeta
		svc       *routing.IstioService
		instances []*model.ServiceInstance
		want      bool
	}{
		{
			meta: model.ConfigMeta{Name: "test", Namespace: "default", Domain: "cluster.local"},
			want: true,
		},
		{
			meta: model.ConfigMeta{Name: "test", Namespace: "default", Domain: "cluster.local"},
			svc:  &routing.IstioService{Name: "hello"},
			want: false,
		},
		{
			meta:      model.ConfigMeta{Name: "test", Namespace: "default", Domain: "cluster.local"},
			svc:       &routing.IstioService{Name: "world"},
			instances: []*model.ServiceInstance{srmemory.MakeInstance(srmemory.HelloService, srmemory.GetPortHTTP(srmemory.HelloService), 0, "")},
			want:      false,
		},
		{
			meta:      model.ConfigMeta{Name: "test", Namespace: "default", Domain: "cluster.local"},
			svc:       &routing.IstioService{Name: "hello"},
			instances: []*model.ServiceInstance{srmemory.MakeInstance(srmemory.HelloService, srmemory.GetPortHTTP(srmemory.HelloService), 0, "")},
			want:      true,
		},
		{
			meta:      model.ConfigMeta{Name: "test", Namespace: "default", Domain: "cluster.local"},
			svc:       &routing.IstioService{Name: "hello", Labels: map[string]string{"version": "v0"}},
			instances: []*model.ServiceInstance{srmemory.MakeInstance(srmemory.HelloService, srmemory.GetPortHTTP(srmemory.HelloService), 0, "")},
			want:      true,
		},
	}

	for _, test := range cases {
		if got := model.MatchSource(test.meta, test.svc, test.instances); got != test.want {
			t.Errorf("MatchSource(%v) => got %v, want %v", test, got, test.want)
		}
	}
}

func TestSortRouteRules(t *testing.T) {
	rules := []model.Config{
		{
			ConfigMeta: model.ConfigMeta{Name: "d"},
			Spec:       &routing.RouteRule{Precedence: 2},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "b"},
			Spec:       &routing.RouteRule{Precedence: 3},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "c"},
			Spec:       &routing.RouteRule{Precedence: 2},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "a"},
		},
	}
	model.SortRouteRules(rules)
	if !(rules[0].Name == "a" && rules[1].Name == "b" && rules[2].Name == "c" && rules[3].Name == "d") {
		t.Errorf("SortRouteRules() => got %#v, want a, b, c, d", rules)
	}
}

type errorStore struct{}

func (errorStore) ConfigDescriptor() model.ConfigDescriptor {
	return model.IstioConfigTypes
}

func (errorStore) Get(typ, name, namespace string) (*model.Config, bool) {
	return nil, false
}

func (errorStore) List(typ, namespace string) ([]model.Config, error) {
	return nil, errors.New("fail")
}

func (errorStore) Create(config model.Config) (string, error) {
	return "", errors.New("fail more")
}

func (errorStore) Update(config model.Config) (string, error) {
	return "", errors.New("yes, fail again")
}

func (errorStore) Delete(typ, name, namespace string) error {
	return errors.New("just keep failing")
}

func TestRouteRules(t *testing.T) {
	instance := srmemory.MakeInstance(srmemory.HelloService, srmemory.GetPortHTTP(srmemory.HelloService), 0, "")
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:      model.RouteRule.Type,
			Name:      "example",
			Namespace: "default",
			Domain:    "cluster.local",
		},
		Spec: &routing.RouteRule{
			Match: &routing.MatchCondition{
				Source: &routing.IstioService{
					Name:   "hello",
					Labels: instance.Labels,
				},
			},
			Destination: &routing.IstioService{
				Name: "world",
			},
		},
	}

	if _, err := store.Create(config); err != nil {
		t.Error(err)
	}
	if out := store.RouteRules([]*model.ServiceInstance{instance}, srmemory.WorldService.Hostname.String()); len(out) != 1 ||
		!reflect.DeepEqual(config.Spec, out[0].Spec) {
		t.Errorf("RouteRules() => expected %#v but got %#v", config.Spec, out)
	}
	if out := store.RouteRules([]*model.ServiceInstance{instance}, srmemory.HelloService.Hostname.String()); len(out) != 0 {
		t.Error("RouteRules() => expected no match for destination-matched rules")
	}
	if out := store.RouteRules(nil, srmemory.WorldService.Hostname.String()); len(out) != 0 {
		t.Error("RouteRules() => expected no match for source-matched rules")
	}

	world := srmemory.MakeInstance(srmemory.WorldService, srmemory.GetPortHTTP(srmemory.WorldService), 0, "")
	if out := store.RouteRulesByDestination([]*model.ServiceInstance{world}); len(out) != 1 ||
		!reflect.DeepEqual(config.Spec, out[0].Spec) {
		t.Errorf("RouteRulesByDestination() => got %#v, want %#v", out, config.Spec)
	}
	if out := store.RouteRulesByDestination([]*model.ServiceInstance{instance}); len(out) != 0 {
		t.Error("RouteRulesByDestination() => expected no match")
	}

	// erroring out list
	if out := model.MakeIstioStore(errorStore{}).RouteRules([]*model.ServiceInstance{instance},
		srmemory.WorldService.Hostname.String()); len(out) != 0 {
		t.Errorf("RouteRules() => expected nil but got %v", out)
	}
	if out := model.MakeIstioStore(errorStore{}).RouteRulesByDestination([]*model.ServiceInstance{world}); len(out) != 0 {
		t.Errorf("RouteRulesByDestination() => expected nil but got %v", out)
	}
}

func TestEgressRules(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	rule := &routing.EgressRule{
		Destination: &routing.IstioService{
			Service: "*.foo.com",
		},
		Ports: []*routing.EgressRule_Port{{
			Port:     80,
			Protocol: "HTTP",
		}},
	}

	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:      model.EgressRule.Type,
			Name:      "example",
			Namespace: "default",
			Domain:    "cluster.local",
		},
		Spec: rule,
	}

	if _, err := store.Create(config); err != nil {
		t.Error(err)
	}

	got := store.EgressRules()
	if len(got) != 1 {
		t.Fatalf("EgressRules() => want 1 rule, got %d", len(got))
	}

	gotRule := got[0].Spec

	if !reflect.DeepEqual(gotRule, rule) {
		t.Errorf("EgressRules() => expected \n%#v, got \n%#v", rule, gotRule)
	}

	// erroring out list
	if out := model.MakeIstioStore(errorStore{}).EgressRules(); len(out) != 0 {
		t.Errorf("EgressRules() => expected nil but got %v", out)
	}
}

func TestDestinationPolicy(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	labels := map[string]string{"version": "v1"}
	instances := []*model.ServiceInstance{srmemory.MakeInstance(srmemory.HelloService, srmemory.GetPortHTTP(srmemory.HelloService), 0, "")}

	policy1 := &routing.DestinationPolicy{
		Source: &routing.IstioService{
			Name:   "hello",
			Labels: map[string]string{"version": "v0"},
		},
		Destination: &routing.IstioService{
			Name:   "world",
			Labels: labels,
		},
	}

	config1 := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:      model.DestinationPolicy.Type,
			Name:      "example",
			Namespace: "default",
			Domain:    "cluster.local",
		},
		Spec: policy1,
	}

	if _, err := store.Create(config1); err != nil {
		t.Error(err)
	}
	if out := store.Policy(instances, srmemory.WorldService.Hostname.String(), labels); out == nil ||
		!reflect.DeepEqual(policy1, out.Spec) {
		t.Errorf("Policy() => expected %#v but got %#v", policy1, out)
	}
	if out := store.Policy(instances, srmemory.HelloService.Hostname.String(), labels); out != nil {
		t.Error("Policy() => expected no match for destination-matched policy")
	}
	if out := store.Policy(instances, srmemory.WorldService.Hostname.String(), nil); out != nil {
		t.Error("Policy() => expected no match for labels-matched policy")
	}
	if out := store.Policy(nil, srmemory.WorldService.Hostname.String(), labels); out != nil {
		t.Error("Policy() => expected no match for source-matched policy")
	}

	// erroring out list
	if out := model.MakeIstioStore(errorStore{}).Policy(instances, srmemory.WorldService.Hostname.String(), labels); out != nil {
		t.Errorf("Policy() => expected nil but got %v", out)
	}
}

func TestAuthenticationPolicyConfig(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))

	authNPolicies := map[string]*authn.Policy{
		model.DefaultAuthenticationPolicyName: {},
		"hello": {
			Targets: []*authn.TargetSelector{{
				Name: "hello",
			}},
			Peers: []*authn.PeerAuthenticationMethod{{
				Params: &authn.PeerAuthenticationMethod_Mtls{},
			}},
		},
		"world": {
			Targets: []*authn.TargetSelector{{
				Name: "world",
				Ports: []*authn.PortSelector{
					{
						Port: &authn.PortSelector_Number{
							Number: 80,
						},
					},
				},
			}},
			Origins: []*authn.OriginAuthenticationMethod{
				{
					Jwt: &authn.Jwt{
						Issuer:  "abc.xzy",
						JwksUri: "https://secure.isio.io",
					},
				},
			},
			PrincipalBinding: authn.PrincipalBinding_USE_ORIGIN,
		},
	}
	for key, value := range authNPolicies {
		config := model.Config{
			ConfigMeta: model.ConfigMeta{
				Type:      model.AuthenticationPolicy.Type,
				Name:      key,
				Group:     "authentication",
				Version:   "v1alpha2",
				Namespace: "default",
				Domain:    "cluster.local",
			},
			Spec: value,
		}
		if _, err := store.Create(config); err != nil {
			t.Error(err)
		}
	}

	cases := []struct {
		hostname model.Hostname
		port     int
		expected string
	}{
		{
			hostname: "hello.default.svc.cluster.local",
			port:     80,
			expected: "hello",
		},
		{
			hostname: "world.default.svc.cluster.local",
			port:     80,
			expected: "world",
		},
		{
			hostname: "world.default.svc.cluster.local",
			port:     8080,
			expected: "default",
		},
		{
			hostname: "world.another-galaxy.svc.cluster.local",
			port:     8080,
			expected: "",
		},
	}

	for _, testCase := range cases {
		port := &model.Port{Port: testCase.port}
		expected := authNPolicies[testCase.expected]
		out := store.AuthenticationPolicyByDestination(testCase.hostname, port)
		if out == nil {
			if expected != nil {
				t.Errorf("AutheticationPolicy(%s:%d) => expected %#v but got nil",
					testCase.hostname, testCase.port, expected)
			}
		} else {
			policy := out.Spec.(*authn.Policy)
			if !reflect.DeepEqual(expected, policy) {
				t.Errorf("AutheticationPolicy(%s:%d) => expected %#v but got %#v",
					testCase.hostname, testCase.port, expected, out)
			}
		}
	}
}

func TestAuthenticationPolicyConfigWithGlobal(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))

	globalPolicy := authn.Policy{
		Peers: []*authn.PeerAuthenticationMethod{{
			Params: &authn.PeerAuthenticationMethod_Mtls{},
		}},
	}
	namespacePolicy := authn.Policy{}
	helloPolicy := authn.Policy{
		Targets: []*authn.TargetSelector{{
			Name: "hello",
		}},
		Peers: []*authn.PeerAuthenticationMethod{{
			Params: &authn.PeerAuthenticationMethod_Mtls{},
		}},
	}

	authNPolicies := []struct {
		name      string
		namespace string
		policy    *authn.Policy
	}{
		{
			name:   model.DefaultAuthenticationPolicyName,
			policy: &globalPolicy,
		},
		{
			name:      model.DefaultAuthenticationPolicyName,
			namespace: "default",
			policy:    &namespacePolicy,
		},
		{
			name:      "hello-policy",
			namespace: "default",
			policy:    &helloPolicy,
		},
	}
	for _, in := range authNPolicies {
		config := model.Config{
			ConfigMeta: model.ConfigMeta{
				Name:    in.name,
				Group:   "authentication",
				Version: "v1alpha2",
				Domain:  "cluster.local",
			},
			Spec: in.policy,
		}
		if in.namespace == "" {
			// Cluster-scoped policy
			config.ConfigMeta.Type = model.AuthenticationMeshPolicy.Type
		} else {
			config.ConfigMeta.Type = model.AuthenticationPolicy.Type
			config.ConfigMeta.Namespace = in.namespace
		}
		if _, err := store.Create(config); err != nil {
			t.Error(err)
		}
	}

	cases := []struct {
		hostname model.Hostname
		port     int
		expected *authn.Policy
	}{
		{
			hostname: "hello.default.svc.cluster.local",
			port:     80,
			expected: &helloPolicy,
		},
		{
			hostname: "world.default.svc.cluster.local",
			port:     80,
			expected: &namespacePolicy,
		},
		{
			hostname: "world.default.svc.cluster.local",
			port:     8080,
			expected: &namespacePolicy,
		},
		{
			hostname: "hello.another-galaxy.svc.cluster.local",
			port:     8080,
			expected: &globalPolicy,
		},
		{
			hostname: "world.another-galaxy.svc.cluster.local",
			port:     9090,
			expected: &globalPolicy,
		},
	}

	for _, testCase := range cases {
		port := &model.Port{Port: testCase.port}
		out := store.AuthenticationPolicyByDestination(testCase.hostname, port)

		if out == nil {
			// With global authentication policy, it's guarantee AuthenticationPolicyByDestination always
			// return non `nill` config.
			t.Errorf("AuthenticationPolicy(%s:%d) => cannot be nil", testCase.hostname, testCase.port)
		} else {
			policy := out.Spec.(*authn.Policy)
			if !reflect.DeepEqual(testCase.expected, policy) {
				t.Errorf("AuthenticationPolicy(%s:%d) => expected:\n%s\nbut got:\n%s\n(from %s/%s)",
					testCase.hostname, testCase.port, testCase.expected.String(), policy.String(), out.Name, out.Namespace)
			}
		}
	}
}

func TestResolveShortnameToFQDN(t *testing.T) {
	tests := []struct {
		name string
		meta model.ConfigMeta
		out  model.Hostname
	}{
		{
			"*", model.ConfigMeta{}, "*",
		},
		{
			"*", model.ConfigMeta{Namespace: "default", Domain: "cluster.local"}, "*",
		},
		{
			"foo", model.ConfigMeta{Namespace: "default", Domain: "cluster.local"}, "foo.default.svc.cluster.local",
		},
		{
			"foo.bar", model.ConfigMeta{Namespace: "default", Domain: "cluster.local"}, "foo.bar",
		},
		{
			"foo", model.ConfigMeta{Domain: "cluster.local"}, "foo.svc.cluster.local",
		},
		{
			"foo", model.ConfigMeta{Namespace: "default"}, "foo.default",
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", idx, tt.out), func(t *testing.T) {
			if actual := model.ResolveShortnameToFQDN(tt.name, tt.meta); actual != tt.out {
				t.Fatalf("model.ResolveShortnameToFQDN(%q, %v) = %q wanted %q", tt.name, tt.meta, actual, tt.out)
			}
		})
	}
}

func TestMostSpecificHostMatch(t *testing.T) {
	tests := []struct {
		in     []model.Hostname
		needle model.Hostname
		want   model.Hostname
	}{
		{[]model.Hostname{}, "*", ""},
		{[]model.Hostname{"*.com", "*.foo.com"}, "bar.foo.com", "*.foo.com"},
		{[]model.Hostname{"*.com", "*.foo.com"}, "foo.com", "*.com"},
		{[]model.Hostname{"*.com", "foo.com"}, "*.foo.com", "*.com"},

		{[]model.Hostname{"*.foo.com", "foo.com"}, "foo.com", "foo.com"},
		{[]model.Hostname{"*.foo.com", "foo.com"}, "*.foo.com", "*.foo.com"},

		// this passes because we sort alphabetically
		{[]model.Hostname{"bar.com", "foo.com"}, "*.com", "bar.com"},

		{[]model.Hostname{"bar.com", "*.foo.com"}, "*foo.com", "*.foo.com"},
		{[]model.Hostname{"foo.com", "*.foo.com"}, "*foo.com", "foo.com"},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", idx, tt.needle), func(t *testing.T) {
			actual, found := model.MostSpecificHostMatch(tt.needle, tt.in)
			if tt.want != "" && !found {
				t.Fatalf("model.MostSpecificHostMatch(%q, %v) = %v, %t; want: %v", tt.needle, tt.in, actual, found, tt.want)
			} else if actual != tt.want {
				t.Fatalf("model.MostSpecificHostMatch(%q, %v) = %v, %t; want: %v", tt.needle, tt.in, actual, found, tt.want)
			}
		})
	}
}

func TestServiceRoles(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	addRbacConfigToStore(model.ServiceRole.Type, "role1", "istio-system", store, t)
	addRbacConfigToStore(model.ServiceRole.Type, "role2", "default", store, t)
	addRbacConfigToStore(model.ServiceRole.Type, "role3", "istio-system", store, t)
	tests := []struct {
		namespace  string
		expectName map[string]bool
	}{
		{namespace: "wrong", expectName: nil},
		{namespace: "default", expectName: map[string]bool{"role2": true}},
		{namespace: "istio-system", expectName: map[string]bool{"role1": true, "role3": true}},
	}

	for _, tt := range tests {
		config := store.ServiceRoles(tt.namespace)
		if tt.expectName != nil {
			for _, cfg := range config {
				if !tt.expectName[cfg.Name] {
					t.Errorf("model.ServiceRoles: expecting %v, but got %v", tt.expectName, config)
				}
			}
		} else if len(config) != 0 {
			t.Errorf("model.ServiceRoles: expecting nil, but got %v", config)
		}
	}
}

func TestServiceRoleBindings(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	addRbacConfigToStore(model.ServiceRoleBinding.Type, "binding1", "istio-system", store, t)
	addRbacConfigToStore(model.ServiceRoleBinding.Type, "binding2", "default", store, t)
	addRbacConfigToStore(model.ServiceRoleBinding.Type, "binding3", "istio-system", store, t)
	tests := []struct {
		namespace  string
		expectName map[string]bool
	}{
		{namespace: "wrong", expectName: nil},
		{namespace: "default", expectName: map[string]bool{"binding2": true}},
		{namespace: "istio-system", expectName: map[string]bool{"binding1": true, "binding3": true}},
	}

	for _, tt := range tests {
		config := store.ServiceRoleBindings(tt.namespace)
		if tt.expectName != nil {
			for _, cfg := range config {
				if !tt.expectName[cfg.Name] {
					t.Errorf("model.ServiceRoleBinding: expecting %v, but got %v", tt.expectName, config)
				}
			}
		} else if len(config) != 0 {
			t.Errorf("model.ServiceRoleBinding: expecting nil, but got %v", config)
		}
	}
}

func TestRbacConfig(t *testing.T) {
	store := model.MakeIstioStore(memory.Make(model.IstioConfigTypes))
	addRbacConfigToStore(model.RbacConfig.Type, model.DefaultRbacConfigName, "", store, t)
	rbacConfig := store.RbacConfig()
	if rbacConfig.Name != model.DefaultRbacConfigName {
		t.Errorf("model.RbacConfig: expecting %s, but got %s", model.DefaultRbacConfigName, rbacConfig.Name)
	}
}

func addRbacConfigToStore(configType, name, namespace string, store model.IstioConfigStore, t *testing.T) {
	var value proto.Message
	switch configType {
	case model.ServiceRole.Type:
		value = &rbacproto.ServiceRole{Rules: []*rbacproto.AccessRule{
			{Services: []string{"service0"}, Methods: []string{"GET"}}}}
	case model.ServiceRoleBinding.Type:
		value = &rbacproto.ServiceRoleBinding{
			Subjects: []*rbacproto.Subject{{User: "User0"}},
			RoleRef:  &rbacproto.RoleRef{Kind: "ServiceRole", Name: "ServiceRole001"}}
	default:
		value = &rbacproto.RbacConfig{Mode: rbacproto.RbacConfig_ON}
	}
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:      configType,
			Name:      name,
			Namespace: namespace,
		},
		Spec: value, // Not used in test, added to pass validation.
	}
	if _, err := store.Create(config); err != nil {
		t.Error(err)
	}
}
