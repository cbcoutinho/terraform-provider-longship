package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// useStateForUnknownModifier implements the plan modifier.
type useStateForUnknownModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m useStateForUnknownModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForUnknownModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyMap implements the plan modification logic.
func (m useStateForUnknownModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {

	tflog.Warn(ctx, fmt.Sprintf("State Value: %+v", req.StateValue))
	tflog.Warn(ctx, fmt.Sprintf("Plan Value: %+v", req.PlanValue))

	// Check if the resource is being created.
	if req.State.Raw.IsNull() {
		tflog.Warn(ctx, "Resource is being created")
	}

	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		tflog.Warn(ctx, "State Value is Null")
		return
	}

	if req.PlanValue.IsNull() {
		tflog.Warn(ctx, "Plan Value is null")
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		tflog.Warn(ctx, "Plan Value is Known")
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		tflog.Warn(ctx, "Config Value is Unknown")
		return
	}

	tflog.Warn(ctx, fmt.Sprintf("Setting Plan Value to %+v", req.StateValue.String()))
	resp.PlanValue = req.StateValue
}

// UseStateForUnknown returns a plan modifier that copies a known prior state
// value into the planned value. Use this when it is known that an unconfigured
// value will remain the same after a resource update.
//
// To prevent Terraform errors, the framework automatically sets unconfigured
// and Computed attributes to an unknown value "(known after apply)" on update.
// Using this plan modifier will instead display the prior state value in the
// plan, unless a prior plan modifier adjusts the value.
func MyUseStateForUnknown() planmodifier.Map {
	return useStateForUnknownModifier{}
}
