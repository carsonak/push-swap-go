package pushswap

import (
	"testing"
)

func TestParseNumberSlice(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		allowDups bool
		want      []float64
		wantErr   bool
	}{
		// --- valid input ---
		{
			name:      "empty input",
			input:     []string{},
			allowDups: false,
			want:      []float64{},
		},
		{
			name:      "single integer",
			input:     []string{"42"},
			allowDups: false,
			want:      []float64{42},
		},
		{
			name:      "multiple integers",
			input:     []string{"1", "2", "3"},
			allowDups: false,
			want:      []float64{1, 2, 3},
		},
		{
			name:      "negative numbers",
			input:     []string{"-5", "-1", "0", "3"},
			allowDups: false,
			want:      []float64{-5, -1, 0, 3},
		},
		{
			name:      "float values",
			input:     []string{"1.5", "2.7", "-0.3"},
			allowDups: false,
			want:      []float64{1.5, 2.7, -0.3},
		},
		{
			name:      "duplicates allowed",
			input:     []string{"3", "1", "3"},
			allowDups: true,
			want:      []float64{3, 1, 3},
		},
		{
			name:      "preserves input order",
			input:     []string{"5", "2", "8", "1"},
			allowDups: false,
			want:      []float64{5, 2, 8, 1},
		},
		// --- error cases ---
		{
			name:      "non-numeric string",
			input:     []string{"abc"},
			allowDups: false,
			wantErr:   true,
		},
		{
			name:      "mixed valid and invalid",
			input:     []string{"1", "two", "3"},
			allowDups: false,
			wantErr:   true,
		},
		{
			name:      "empty string element",
			input:     []string{""},
			allowDups: false,
			wantErr:   true,
		},
		{
			name:      "duplicate rejected when allowDups=false",
			input:     []string{"1", "2", "1"},
			allowDups: false,
			wantErr:   true,
		},
		{
			name:      "duplicate floats rejected",
			input:     []string{"1.5", "1.5"},
			allowDups: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumberSlice(tt.input, tt.allowDups)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseNumberSlice() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("ParseNumberSlice() len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ParseNumberSlice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
