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

	// Map of Azure resource types to common redundant suffix abbreviations
	commonSuffixes := map[string][]string{
		"azurerm_resource_group": {"rg", "resource_group", "group"},
		"azurerm_key_vault":      {"kv", "vault", "key_vault"},
		"azurerm_storage_account":{"st", "sa", "storage", "storage_account"},
		"azurerm_virtual_network":{"vnet", "network"},
		"azurerm_subnet":         {"snet", "subnet"},
	}

	for _, block := range content.Blocks {
		if len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0] // e.g., "azurerm_resource_group"
		localName := block.Labels[1]    // e.g., "demo_rg"

		// Check 1: Must use underscores instead of hyphens
		if strings.Contains(localName, "-") {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Resource local name '%s' must use underscores ('_') instead of hyphens ('-').", localName),
				block.DefRange,
			)
		}

		// Check 2: Prevent resource type suffix redundant inclusions (both full words and common abbreviations)
		cleanType := strings.TrimPrefix(resourceType, "azurerm_")
		typeParts := strings.Split(cleanType, "_")

		// Collect all redundant suffixes to test (word parts + explicit abbreviations)
		forbiddenSuffixes := append([]string{}, typeParts...)
		if extra, exists := commonSuffixes[resourceType]; exists {
			forbiddenSuffixes = append(forbiddenSuffixes, extra...)
		}

		for _, suffix := range forbiddenSuffixes {
			if len(suffix) > 1 && (strings.HasSuffix(localName, "_"+suffix) || localName == suffix) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Resource local name '%s' contains redundant type suffix '_%s'. Omit type names or abbreviations from the local name.", localName, suffix),
					block.DefRange,
				)
				break
			}
		}
	}

	return nil
}