package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type NSGRuleNamingRule struct {
	tflint.DefaultRule
}

func NewNSGRuleNamingRule() *NSGRuleNamingRule {
	return &NSGRuleNamingRule{}
}

func (r *NSGRuleNamingRule) Name() string {
	return "azurerm_nsg_rule_naming"
}

func (r *NSGRuleNamingRule) Enabled() bool {
	return true
}

func (r *NSGRuleNamingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Format: <Access>-<Direction>-<Protocol>-<RuleName(camelCase)>
// Example: AI-TCP-InternetToWebServer, DO-ANY-DLP
var nsgRuleNameRegex = regexp.MustCompile(`^(A|D)(I|O)-(TCP|UDP|ICMP|ANY)-([a-z][a-zA-Z0-9]*|[A-Z][a-zA-Z0-9]*)$`)

func (r *NSGRuleNamingRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "name"},
					},
					Blocks: []hclext.BlockSchema{
						{
							Type: "security_rule",
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "name"},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0]

		if resourceType == "azurerm_network_security_rule" {
			checkRuleNameAttribute(runner, r, block.Body.Attributes["name"])
		}

		if resourceType == "azurerm_network_security_group" {
			for _, inlineRule := range block.Body.Blocks {
				if inlineRule.Type == "security_rule" {
					checkRuleNameAttribute(runner, r, inlineRule.Body.Attributes["name"])
				}
			}
		}
	}

	return nil
}

func checkRuleNameAttribute(runner tflint.Runner, r *NSGRuleNamingRule, attr *hclext.Attribute) {
	if attr == nil {
		return
	}

	var ruleName string
	err := runner.EvaluateExpr(attr.Expr, &ruleName, nil)
	if err != nil {
		return
	}

	if !nsgRuleNameRegex.MatchString(ruleName) {
		runner.EmitIssue(
			r,
			fmt.Sprintf("NSG rule name '%s' does not follow the required naming format '<Access><Direction>-<Protocol>-<RuleName>' (e.g., 'AI-TCP-InternetToWebServer' or 'DO-ANY-DLP').", ruleName),
			attr.Range,
		)
	}
}