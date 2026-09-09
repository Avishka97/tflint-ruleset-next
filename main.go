package main

import (
	"github.com/Avishka97/tflint-ruleset-next/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "next",
			Version: "0.1.3",
			Rules: []tflint.Rule{
				rules.NewResourceNamingRule(),
				rules.NewNSGRulePriorityRule(),
				rules.NewNSGRuleNamingRule(),
			},
		},
	})
}