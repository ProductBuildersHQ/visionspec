package workflow

import "testing"

func TestReviewGateStepID(t *testing.T) {
	tests := []struct {
		name string
		gate ReviewGate
		want string
	}{
		{
			name: "defaults to kebab-cased action",
			gate: ReviewGate{After: "faq", Action: "prfaq_review"},
			want: "prfaq-review",
		},
		{
			name: "single-word action unchanged",
			gate: ReviewGate{After: "narrative-6p", Action: "decision_meeting", Required: true},
			want: "decision-meeting",
		},
		{
			name: "explicit id wins",
			gate: ReviewGate{ID: "board-review", After: "narrative-6p", Action: "decision_meeting"},
			want: "board-review",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gate.StepID(); got != tt.want {
				t.Errorf("StepID() = %q, want %q", got, tt.want)
			}
		})
	}
}
