package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_NSGRulePriorityRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid priority in devops range",
			Content: `
resource "azurerm_network_security_rule" "valid" {
  name     = "AI-TCP-webServer"
  priority = 1050
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid priority in platform reserved low range",
			Content: `
resource "azurerm_network_security_rule" "reserved_low" {
  name     = "AI-TCP-webServer"
  priority = 500
}`,
			Expected: helper.Issues{
				{
					Rule:    NewNSGRulePriorityRule(),
					Message: "NSG rule priority 500 is outside the allowed Landing Zone DevOps range (1000-3999). Priorities 100-999 and 4000-4096 are reserved for the Platform Team.",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 17},
					},
				},
			},
		},
	}

	rule := NewNSGRulePriorityRule()
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

func Test_NSGRuleNamingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid nsg rule name",
			Content: `
resource "azurerm_network_security_rule" "valid" {
  name     = "AI-TCP-InternetToWebServer"
  priority = 1000
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid nsg rule name format",
			Content: `
resource "azurerm_network_security_rule" "invalid" {
  name     = "allow_inbound_tcp_web"
  priority = 1000
}`,
			Expected: helper.Issues{
				{
					Rule:    NewNSGRuleNamingRule(),
					Message: "NSG rule name 'allow_inbound_tcp_web' does not follow the required naming format '<Access><Direction>-<Protocol>-<RuleName>' (e.g., 'AI-TCP-InternetToWebServer' or 'DO-ANY-DLP').",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 37},
					},
				},
			},
		},
	}

	rule := NewNSGRuleNamingRule()
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