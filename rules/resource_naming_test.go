package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ResourceNamingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid resource name",
			Content: `resource "azurerm_resource_group" "application" {
  name     = "rg-test"
  location = "East US"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid camelCase name without redundant tokens",
			Content: `resource "azurerm_key_vault" "primaryApp" {
  name     = "kv-test"
  location = "East US"
}`,
			Expected: helper.Issues{},
		},

		// --- CASE 1: Hyphen Check ---
		{
			Name: "case 1 - invalid resource name with hyphens and suffix",
			Content: `resource "azurerm_resource_group" "my-app-rg" {
  name     = "rg-test"
  location = "East US"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'my-app-rg' must use underscores ('_') instead of hyphens ('-').",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 46},
					},
				},
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'my-app-rg' contains redundant type suffix or abbreviation 'rg'. Omit type names or abbreviations from the local name.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 46},
					},
				},
			},
		},

		// --- CASE 2: Casing Format Check ---
		{
			Name: "case 2 - invalid PascalCase resource name",
			Content: `resource "azurerm_key_vault" "PrimarySecrets" {
  name     = "kv-test"
  location = "East US"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'PrimarySecrets' must start with a lowercase letter.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 46},
					},
				},
			},
		},

		// --- CASE 3: Redundant Token Check ---
		{
			Name: "case 3 - redundant abbreviation prefix",
			Content: `resource "azurerm_key_vault" "kvSecrets" {
  name     = "kv-test"
  location = "East US"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'kvSecrets' contains redundant type suffix or abbreviation 'kv'. Omit type names or abbreviations from the local name.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 41},
					},
				},
			},
		},
		{
			Name: "case 3 - redundant type suffix attached",
			Content: `resource "azurerm_storage_account" "primarystorage" {
  name     = "sttest"
  location = "East US"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'primarystorage' contains redundant type suffix or abbreviation 'storage'. Omit type names or abbreviations from the local name.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 52},
					},
				},
			},
		},
		{
			Name: "case 3 - redundant abbreviation infix with underscore",
			Content: `resource "azurerm_subnet" "internal_snet_backend" {
  name     = "snet-test"
  location = "East US"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamingRule(),
					Message: "Resource local name 'internal_snet_backend' contains redundant type suffix or abbreviation 'snet'. Omit type names or abbreviations from the local name.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 50},
					},
				},
			},
		},
	}

	rule := NewResourceNamingRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}