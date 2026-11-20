package tlsrule

import (
	"context"
	"testing"

	cato_models "github.com/catonetworks/cato-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func nullAttrs(ctx context.Context, t *testing.T, attrTypes map[string]attr.Type) map[string]attr.Value {
	t.Helper()
	values := make(map[string]attr.Value, len(attrTypes))
	for name, attrType := range attrTypes {
		value, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), nil))
		require.NoError(t, err)
		values[name] = value
	}
	return values
}

func tlsRulePlan(ctx context.Context, t *testing.T, category types.String) InspectionRule {
	t.Helper()

	applicationAttrs := nullAttrs(ctx, t, TLSApplicationAttrTypes)
	applicationAttrs["tls_inspect_category"] = category
	application, diags := types.ObjectValue(TLSApplicationAttrTypes, applicationAttrs)
	require.False(t, diags.HasError(), "%v", diags)

	ruleAttrs := nullAttrs(ctx, t, TLSInspectionRuleRuleAttrTypes)
	ruleAttrs["name"] = types.StringValue("inspect-category")
	ruleAttrs["application"] = application
	rule, diags := types.ObjectValue(TLSInspectionRuleRuleAttrTypes, ruleAttrs)
	require.False(t, diags.HasError(), "%v", diags)

	return InspectionRule{
		At:   types.ObjectNull(map[string]attr.Type{"position": types.StringType, "ref": types.StringType}),
		Rule: rule,
		ID:   types.StringNull(),
	}
}

func TestHydrateTLSRuleAPISendsTLSInspectCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	for _, category := range cato_models.AllTLSInspectSystemCategory {
		t.Run(category.String(), func(t *testing.T) {
			t.Parallel()
			input, diags := hydrateTLSRuleAPI(ctx, tlsRulePlan(ctx, t, types.StringValue(category.String())))
			require.Empty(t, diags)

			want := []cato_models.TLSInspectSystemCategory{category}
			require.NotNil(t, input.create.Rule.Application)
			require.Equal(t, want, input.create.Rule.Application.TLSInspectCategory)
			require.NotNil(t, input.update.Rule.Application)
			require.Equal(t, want, input.update.Rule.Application.TLSInspectCategory)
		})
	}
}

func TestHydrateTLSRuleAPIWithoutTLSInspectCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	input, diags := hydrateTLSRuleAPI(ctx, tlsRulePlan(ctx, t, types.StringNull()))
	require.Empty(t, diags)

	require.NotNil(t, input.create.Rule.Application)
	require.Empty(t, input.create.Rule.Application.TLSInspectCategory)
	require.NotNil(t, input.update.Rule.Application)
	require.NotNil(t, input.update.Rule.Application.TLSInspectCategory)
	require.Empty(t, input.update.Rule.Application.TLSInspectCategory)
}
