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
			Name: "invalid resource name with hyphens and suffix",
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
			},
		},
		{
			Name: "valid resource name",
			Content: `resource "azurerm_resource_group" "application" {
  name     = "rg-test"
  location = "East US"
}`,
			Expected: helper.Issues{},
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