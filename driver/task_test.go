package driver

import (
	"reflect"
	"testing"

	"github.com/hashicorp/consul-terraform-sync/client"
	"github.com/hashicorp/consul-terraform-sync/config"
	mocks "github.com/hashicorp/consul-terraform-sync/mocks/client"
	"github.com/hashicorp/consul-terraform-sync/templates/hcltmpl"
	"github.com/hashicorp/consul-terraform-sync/templates/tftmpl"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		clientType  string
		expectError bool
		expect      client.Client
	}{
		{
			"happy path with development client",
			developmentClient,
			false,
			&client.Printer{},
		},
		{
			"happy path with mock client",
			testClient,
			false,
			&mocks.Client{},
		},
		{
			"error when creating Terraform CLI client",
			"",
			true,
			&client.TerraformCLI{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := newClient(&clientConfig{
				clientType: tc.clientType,
			})
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, reflect.TypeOf(tc.expect), reflect.TypeOf(actual))
			}
		})
	}
}

func TestTask_ProviderNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		task     Task
		expected []string
	}{
		{
			"no provider",
			Task{},
			[]string{},
		},
		{
			"happy path",
			Task{
				providers: NewTerraformProviderBlocks(
					hcltmpl.NewNamedBlocksTest([]map[string]interface{}{
						{"local": map[string]interface{}{
							"configs": "stuff",
						}},
						{"null": map[string]interface{}{}},
					})),
			},
			[]string{"local", "null"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.task.ProviderNames()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTask_ServiceNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		task     Task
		expected []string
	}{
		{
			"no services",
			Task{},
			[]string{},
		},
		{
			"happy path",
			Task{
				services: []Service{
					Service{Name: "web"},
					Service{Name: "api"},
				},
			},
			[]string{"web", "api"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.task.ServiceNames()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTask_configureCondition(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		task     Task
		expected tftmpl.Condition
	}{
		{
			"services",
			Task{
				condition: &config.ServicesConditionConfig{
					config.ServicesMonitorConfig{
						Regexp: config.String("test"),
					},
				},
			},
			&tftmpl.ServicesCondition{
				tftmpl.ServicesMonitor{
					Regexp: "test",
				},
				true,
			},
		},
		{
			"catalog-services",
			Task{
				condition: &config.CatalogServicesConditionConfig{
					config.CatalogServicesMonitorConfig{
						Regexp:     config.String("test"),
						Datacenter: config.String("dc2"),
						Namespace:  config.String("ns2"),
						NodeMeta: map[string]string{
							"key1": "value1",
						},
						SourceIncludesVar: config.Bool(true),
					},
				},
			},
			&tftmpl.CatalogServicesCondition{
				tftmpl.CatalogServicesMonitor{
					Regexp:     "test",
					Datacenter: "dc2",
					Namespace:  "ns2",
					NodeMeta: map[string]string{
						"key1": "value1",
					},
				},
				true,
			},
		},
		{
			"consul-kv",
			Task{
				condition: &config.ConsulKVConditionConfig{
					config.ConsulKVMonitorConfig{
						Path:              config.String("key-path"),
						Recurse:           config.Bool(true),
						Datacenter:        config.String("dc2"),
						Namespace:         config.String("ns2"),
						SourceIncludesVar: config.Bool(true),
					},
				},
			},
			&tftmpl.ConsulKVCondition{
				tftmpl.ConsulKVMonitor{
					Path:       "key-path",
					Recurse:    true,
					Datacenter: "dc2",
					Namespace:  "ns2",
				},
				true,
			},
		},
		{
			"schedule",
			Task{
				condition: &config.ScheduleConditionConfig{
					Cron: config.String("10 2 * * 3"),
				},
			},
			&tftmpl.ServicesCondition{
				tftmpl.ServicesMonitor{},
				true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.task.configureCondition()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
