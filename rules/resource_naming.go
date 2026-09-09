package rules

import (
	"fmt"
	"regexp"
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

// Regex matching lower camelCase (starts with lowercase letter, alphanumeric)
var camelCaseRegex = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)

// Common provider-specific abbreviations/synonyms mapped to resource types.
var commonTypeTokens = map[string][]string{
	"azurerm_resource_group":  {"rg", "resourcegroup", "resgroup", "group"},
	"azurerm_key_vault":       {"kv", "keyvault", "vault"},
	"azurerm_storage_account": {"st", "sa", "storage", "storageaccount"},
	"azurerm_virtual_network": {"vnet", "vn", "virtualnetwork", "network"},
	"azurerm_subnet":          {"snet", "subnet", "sub"},
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

	for _, block := range content.Blocks {
		if len(block.Labels) < 2 {
			continue
		}

		resourceType := block.Labels[0] // e.g., "azurerm_resource_group"
		localName := block.Labels[1]    // e.g., "my-app-rg", "DemoRg"

		// --- CHECK 1: Hyphen Violation ---
		if strings.Contains(localName, "-") {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Resource local name '%s' must use underscores ('_') instead of hyphens ('-').", localName),
				block.DefRange,
			)
		}

		// --- CHECK 2: Casing Format Violation (Starts with uppercase/PascalCase) ---
		if !strings.Contains(localName, "-") && !strings.Contains(localName, "_") {
			if !camelCaseRegex.MatchString(localName) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Resource local name '%s' must start with a lowercase letter.", localName),
					block.DefRange,
				)
			}
		}

		// --- CHECK 3: Redundant Type Tokens / Abbreviations ---
		tokens := collectForbiddenTokens(resourceType)
		if matchedToken := findRedundantToken(localName, tokens); matchedToken != "" {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Resource local name '%s' contains redundant type suffix or abbreviation '%s'. Omit type names or abbreviations from the local name.", localName, matchedToken),
				block.DefRange,
			)
		}
	}

	return nil
}

// collectForbiddenTokens extracts type words from the resource type string and appends custom overrides.
func collectForbiddenTokens(resourceType string) []string {
	tokenMap := make(map[string]struct{})

	cleanType := strings.TrimPrefix(resourceType, "azurerm_")
	cleanType = strings.TrimPrefix(cleanType, "aws_")
	cleanType = strings.TrimPrefix(cleanType, "google_")

	for _, part := range strings.Split(cleanType, "_") {
		if len(part) > 1 {
			tokenMap[strings.ToLower(part)] = struct{}{}
		}
	}

	if custom, exists := commonTypeTokens[resourceType]; exists {
		for _, token := range custom {
			tokenMap[strings.ToLower(token)] = struct{}{}
		}
	}

	tokens := make([]string, 0, len(tokenMap))
	for t := range tokenMap {
		tokens = append(tokens, t)
	}

	return tokens
}

// findRedundantToken checks for token occurrences across prefixes, suffixes, segments, or attached words.
func findRedundantToken(name string, tokens []string) string {
	segments := splitNameSegments(name)

	for _, token := range tokens {
		// Exact segment match (e.g., "demo_rg", "demoRg" -> segment "rg")
		for _, segment := range segments {
			if segment == token {
				return token
			}
		}

		// Attached match at start or end (e.g., "demorg", "rgdemo")
		normalized := strings.ToLower(name)
		if len(token) >= 2 {
			if strings.HasPrefix(normalized, token) || strings.HasSuffix(normalized, token) {
				return token
			}
		}
	}

	return ""
}

// splitNameSegments breaks strings by delimiters and camelCase boundaries.
func splitNameSegments(s string) []string {
	rawParts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})

	var segments []string
	for _, part := range rawParts {
		var current strings.Builder
		for i, r := range part {
			if i > 0 && (r >= 'A' && r <= 'Z') {
				if current.Len() > 0 {
					segments = append(segments, strings.ToLower(current.String()))
					current.Reset()
				}
			}
			current.WriteRune(r)
		}
		if current.Len() > 0 {
			segments = append(segments, strings.ToLower(current.String()))
		}
	}

	return segments
}