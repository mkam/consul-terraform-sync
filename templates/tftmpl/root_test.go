package tftmpl

import (
	"testing"

	"github.com/hashicorp/consul-terraform-sync/config"
	"github.com/hashicorp/consul-terraform-sync/templates/hcltmpl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendRootTerraformBlock_backend(t *testing.T) {
	consulBackend, err := config.DefaultTerraformBackend(&config.ConsulConfig{
		Address: config.String("consul.example.com"),
		TLS: &config.TLSConfig{
			Enabled: config.Bool(true),
			CACert:  config.String("ca_cert"),
			Cert:    config.String("cert"),
			Key:     config.String("key"),
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		rawBackend map[string]interface{}
		expected   string
	}{
		{
			"nil",
			nil,
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
}
`,
		}, {
			"empty",
			map[string]interface{}{"empty": map[string]interface{}{}},
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
  backend "empty" {
  }
}
`,
		}, {
			"invalid structure",
			map[string]interface{}{"invalid": "unexpected type"},
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
}
`,
		}, {
			"local",
			map[string]interface{}{"local": map[string]interface{}{
				"path": "relative/path/to/terraform.tfstate",
			}},
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
  backend "local" {
    path = "relative/path/to/terraform.tfstate"
  }
}
`,
		}, {
			"consul",
			consulBackend,
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
  backend "consul" {
    address   = "consul.example.com"
    ca_file   = "ca_cert"
    cert_file = "cert"
    gzip      = true
    key_file  = "key"
    path      = "consul-terraform-sync/terraform"
    scheme    = "https"
  }
}
`,
		}, {
			"postgres",
			map[string]interface{}{"pg": map[string]interface{}{
				"conn_str": "postgres://user:pass@db.example.com/terraform_backend",
			}},
			`terraform {
  required_version = ">= 0.13.0, < 0.16"
  backend "pg" {
    conn_str = "postgres://user:pass@db.example.com/terraform_backend"
  }
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hclFile := hclwrite.NewEmptyFile()
			body := hclFile.Body()

			var backend *hcltmpl.NamedBlock
			if tc.rawBackend != nil {
				b := hcltmpl.NewNamedBlock(tc.rawBackend)
				backend = &b
			}
			appendRootTerraformBlock(body, backend, nil)

			content := hclFile.Bytes()
			content = hclwrite.Format(content)
			assert.Equal(t, tc.expected, string(content))
		})
	}
}

func TestAppendRootProviderBlocks(t *testing.T) {
	testCases := []struct {
		name       string
		rawBackend map[string]interface{}
		expected   string
	}{
		{
			"nil",
			nil,
			`provider "" {
}
`,
		}, {
			"empty",
			map[string]interface{}{"empty": map[string]interface{}{}},
			`provider "empty" {
}
`,
		}, {
			"internal alias leak",
			map[string]interface{}{"foo": map[string]interface{}{
				"alias": "bar",
			}},
			`provider "foo" {
}
`,
		}, {
			"internal auto_commit leak",
			map[string]interface{}{"foo": map[string]interface{}{
				"auto_commit": "true",
			}},
			`provider "foo" {
}
`,
		}, {
			"invalid structure",
			map[string]interface{}{"invalid": "unexpected type"},
			`provider "" {
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hclFile := hclwrite.NewEmptyFile()
			body := hclFile.Body()

			backend := []hcltmpl.NamedBlock{hcltmpl.NewNamedBlock(tc.rawBackend)}
			appendRootProviderBlocks(body, backend)

			content := hclFile.Bytes()
			content = hclwrite.Format(content)
			assert.Equal(t, tc.expected, string(content))
		})
	}
}

func TestService_hcatQuery(t *testing.T) {
	testCases := []struct {
		name     string
		service  Service
		expected string
	}{
		{
			"empty",
			Service{},
			`""`,
		}, {
			"base",
			Service{Name: "app"},
			`"app"`,
		}, {
			"datacenter",
			Service{
				Name:       "app",
				Datacenter: "dc1",
			},
			`"app" "dc=dc1"`,
		}, {
			"namespace",
			Service{
				Name:      "app",
				Namespace: "namespace",
			},
			`"app" "ns=namespace"`,
		}, {
			"tag",
			Service{
				Name: "app",
				Tag:  "my-tag",
			},
			`"app" "\"my-tag\" in Service.Tags"`,
		}, {
			"all",
			Service{
				Name:       "app",
				Datacenter: "dc1",
				Namespace:  "namespace",
				Tag:        "my-tag",
			},
			`"app" "dc=dc1" "ns=namespace" "\"my-tag\" in Service.Tags"`,
		},
	}
	for _, tc := range testCases {
		actual := tc.service.hcatQuery()
		assert.Equal(t, tc.expected, actual)
	}
}
