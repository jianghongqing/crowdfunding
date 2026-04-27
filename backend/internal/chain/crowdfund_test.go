package chain

import (
	"math/big"
	"testing"
	"time"
)

func TestDeriveStatus(t *testing.T) {
	now := uint64(time.Now().Unix())

	tests := []struct {
		name      string
		pledged   *big.Int
		goal      *big.Int
		deadline  *big.Int
		withdrawn bool
		want      string
	}{
		{
			name:      "withdrawn campaign",
			pledged:   big.NewInt(2),
			goal:      big.NewInt(1),
			deadline:  big.NewInt(int64(now + 100)),
			withdrawn: true,
			want:      "succeeded_withdrawn",
		},
		{
			name:      "goal reached",
			pledged:   big.NewInt(10),
			goal:      big.NewInt(5),
			deadline:  big.NewInt(int64(now + 100)),
			withdrawn: false,
			want:      "goal_reached_pending_withdraw",
		},
		{
			name:      "expired campaign",
			pledged:   big.NewInt(1),
			goal:      big.NewInt(5),
			deadline:  big.NewInt(int64(now - 100)),
			withdrawn: false,
			want:      "failed_refundable",
		},
		{
			name:      "active campaign",
			pledged:   big.NewInt(1),
			goal:      big.NewInt(5),
			deadline:  big.NewInt(int64(now + 100)),
			withdrawn: false,
			want:      "active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveStatus(tt.pledged, tt.goal, tt.deadline, tt.withdrawn)
			if got != tt.want {
				t.Fatalf("DeriveStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}
