package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type NSGRulePriorityRule struct {
	tflint.DefaultRule
}

func NewNSGRulePriorityRule() *NSGRulePriorityRule {
	return &NSGRulePriorityRule{}
}

func (r *NSGRulePriorityRule) Name() string {
	return "azurerm_nsg_rule_priority"
}

func (r *NSGRulePriorityRule) Enabled() bool {
	return true
}

func (r *NSGRulePriorityRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NSGRulePriorityRule) Check(runner tflint.Runner) error {
	// Inspect both standalone security rules and inline security rules inside NSG blocks
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "priority"},
					},
					Blocks: []hclext.BlockSchema{
						{
							Type: "security_rule",
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "priority"},
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

		// Check standalone azurerm_network_security_rule resources
		if resourceType == "azurerm_network_security_rule" {
			checkPriorityAttribute(runner, r, block.Body.Attributes["priority"])
		}

		// Check inline security_rule blocks within azurerm_network_security_group
		if resourceType == "azurerm_network_security_group" {
			for _, inlineRule := range block.Body.Blocks {
				if inlineRule.Type == "security_rule" {
					checkPriorityAttribute(runner, r, inlineRule.Body.Attributes["priority"])
				}
			}
		}
	}

	return nil
}

func checkPriorityAttribute(runner tflint.Runner, r *NSGRulePriorityRule, attr *hclext.Attribute) {
	if attr == nil {
		return
	}

	var priority int
	err := runner.EvaluateExpr(attr.Expr, &priority, nil)
	if err != nil {
		return // Skip dynamic/unresolvable expressions
	}

	if priority < 1000 || priority > 3999 {
		runner.EmitIssue(
			r,
			fmt.Sprintf("NSG rule priority %d is outside the allowed Landing Zone DevOps range (1000-3999). Priorities 100-999 and 4000-4096 are reserved for the Platform Team.", priority),
			attr.Range,
		)
	}
}