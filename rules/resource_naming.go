package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ResourceNamingRule struct {
	tflint.DefaultRule
}

func NewResourceNamingRule() *ResourceNamingRule {
	return &ResourceNamingRule{}
}

func (r *ResourceNamingRule) Name() string {
	return "company_resource_naming"
}

func (r *ResourceNamingRule) Enabled() bool {
	return true
}

func (r *ResourceNamingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *ResourceNamingRule) Check(runner tflint.Runner) error {
	// Query top-level resource blocks using GetModuleContent
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
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

		resourceType := block.Labels[0] // e.g., "azurerm_resource_group"
		localName := block.Labels[1]    // e.g., "my_rg" or "my-rg"

		// Check 1: Must use underscores instead of hyphens
		if strings.Contains(localName, "-") {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Resource local name '%s' must use underscores ('_') instead of hyphens ('-').", localName),
				block.DefRange,
			)
		}

		// Check 2: Prevent resource type suffix redundant inclusions
		cleanType := strings.TrimPrefix(resourceType, "azurerm_")
		typeParts := strings.Split(cleanType, "_")

		for _, part := range typeParts {
			if len(part) > 1 && (strings.HasSuffix(localName, "_"+part) || localName == part) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Resource local name '%s' contains redundant type suffix '_%s'. Omit type names from the local name.", localName, part),
					block.DefRange,
				)
				break
			}
		}
	}

	return nil
}