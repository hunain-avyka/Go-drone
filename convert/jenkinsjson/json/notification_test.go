package json

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	harness "github.com/hunain-avyka/go-spec/dist/go"
)

func TestConvertNotification(t *testing.T) { // Convert bytes to string

	var tests []runner
	tests = append(tests, prepare(t, "/Notification/notification_snippet", &harness.Step{
		Id:   "notifyEndpointsc1ca8f",
		Name: "Notification",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/webhook",
			With: map[string]interface{}{
				"urls":         []string{"https://webhook.site/9ffae84b-a338-43ef-9283-319d70574bf4"},
				"method":       "POST",
				"token-value":  "<+input>",
				"content-type": "application/json",
				"template": `{
    "status": "SUCCESSFUL",
    "notes": "Build metrics for analysis"
}`,
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertNotification(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertNotification() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
