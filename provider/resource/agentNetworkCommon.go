package resource

import (
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// AgentNetworkTokenLimit is a per-policy token cap, applied independently per source group and per user.
type AgentNetworkTokenLimit struct {
	Enabled       bool `pulumi:"enabled"`
	GroupCap      int  `pulumi:"groupCap"`
	UserCap       int  `pulumi:"userCap"`
	WindowSeconds int  `pulumi:"windowSeconds"`
}

// Annotate provides documentation for AgentNetworkTokenLimit fields.
func (t *AgentNetworkTokenLimit) Annotate(a infer.Annotator) {
	a.Describe(&t.Enabled, "Whether the token limit is enforced.")
	a.Describe(&t.GroupCap, "Tokens allowed per source group within the window (each group has its own bucket). 0 means uncapped.")
	a.Describe(&t.UserCap, "Tokens allowed per individual user within the window. 0 means uncapped.")
	a.Describe(&t.WindowSeconds, "Reset frequency in seconds. Minimum 60 when the limit is enabled.")
}

// AgentNetworkBudgetLimit is a per-policy USD spend cap, applied independently per source group and per user.
type AgentNetworkBudgetLimit struct {
	Enabled       bool    `pulumi:"enabled"`
	GroupCapUsd   float64 `pulumi:"groupCapUsd"`
	UserCapUsd    float64 `pulumi:"userCapUsd"`
	WindowSeconds int     `pulumi:"windowSeconds"`
}

// Annotate provides documentation for AgentNetworkBudgetLimit fields.
func (b *AgentNetworkBudgetLimit) Annotate(a infer.Annotator) {
	a.Describe(&b.Enabled, "Whether the budget limit is enforced.")
	a.Describe(&b.GroupCapUsd, "USD allowed per source group within the window (each group has its own bucket). 0 means uncapped.")
	a.Describe(&b.UserCapUsd, "USD allowed per individual user within the window. 0 means uncapped.")
	a.Describe(&b.WindowSeconds, "Reset frequency in seconds. Minimum 60 when the limit is enabled.")
}

// AgentNetworkLimits bundles the token and budget caps attached to a policy or budget rule.
type AgentNetworkLimits struct {
	TokenLimit  AgentNetworkTokenLimit  `pulumi:"tokenLimit"`
	BudgetLimit AgentNetworkBudgetLimit `pulumi:"budgetLimit"`
}

// Annotate provides documentation for AgentNetworkLimits fields.
func (l *AgentNetworkLimits) Annotate(a infer.Annotator) {
	a.Describe(&l.TokenLimit, "Token cap composed with any guardrail-level checks.")
	a.Describe(&l.BudgetLimit, "USD budget cap composed with any guardrail-level checks.")
}

func toAPIAgentNetworkLimits(limits AgentNetworkLimits) nbapi.AgentNetworkPolicyLimits {
	return nbapi.AgentNetworkPolicyLimits{
		TokenLimit: nbapi.AgentNetworkPolicyTokenLimit{
			Enabled:       limits.TokenLimit.Enabled,
			GroupCap:      int64(limits.TokenLimit.GroupCap),
			UserCap:       int64(limits.TokenLimit.UserCap),
			WindowSeconds: int64(limits.TokenLimit.WindowSeconds),
		},
		BudgetLimit: nbapi.AgentNetworkPolicyBudgetLimit{
			Enabled:       limits.BudgetLimit.Enabled,
			GroupCapUsd:   limits.BudgetLimit.GroupCapUsd,
			UserCapUsd:    limits.BudgetLimit.UserCapUsd,
			WindowSeconds: int64(limits.BudgetLimit.WindowSeconds),
		},
	}
}

func fromAPIAgentNetworkLimits(limits nbapi.AgentNetworkPolicyLimits) AgentNetworkLimits {
	return AgentNetworkLimits{
		TokenLimit: AgentNetworkTokenLimit{
			Enabled:       limits.TokenLimit.Enabled,
			GroupCap:      int(limits.TokenLimit.GroupCap),
			UserCap:       int(limits.TokenLimit.UserCap),
			WindowSeconds: int(limits.TokenLimit.WindowSeconds),
		},
		BudgetLimit: AgentNetworkBudgetLimit{
			Enabled:       limits.BudgetLimit.Enabled,
			GroupCapUsd:   limits.BudgetLimit.GroupCapUsd,
			UserCapUsd:    limits.BudgetLimit.UserCapUsd,
			WindowSeconds: int(limits.BudgetLimit.WindowSeconds),
		},
	}
}

// equalStringMapPtr compares two *map[string]string values by content, treating nil and an
// empty map as equal.
func equalStringMapPtr(mapA, mapB *map[string]string) bool {
	var valueA, valueB map[string]string

	if mapA != nil {
		valueA = *mapA
	}

	if mapB != nil {
		valueB = *mapB
	}

	if len(valueA) != len(valueB) {
		return false
	}

	for key, val := range valueA {
		if valueB[key] != val {
			return false
		}
	}

	return true
}

// equalAgentNetworkLimitsPtr compares two optional AgentNetworkLimits, treating a nil pointer
// as equivalent to the zero-value limits (matches how the NetBird API always echoes back a
// full limits struct even when the caller never configured one).
func equalAgentNetworkLimitsPtr(limitsA, limitsB *AgentNetworkLimits) bool {
	var zero AgentNetworkLimits

	valueA := zero
	if limitsA != nil {
		valueA = *limitsA
	}

	valueB := zero
	if limitsB != nil {
		valueB = *limitsB
	}

	return valueA == valueB
}
