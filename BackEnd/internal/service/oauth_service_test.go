package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomState(t *testing.T) {
	state1, err := GenerateRandomState()
	assert.NoError(t, err)
	assert.NotEmpty(t, state1)

	state2, err := GenerateRandomState()
	assert.NoError(t, err)
	assert.NotEmpty(t, state2)

	// State phải ngẫu nhiên và không trùng nhau
	assert.NotEqual(t, state1, state2)
	assert.GreaterOrEqual(t, len(state1), 32)
}

func TestValidateState(t *testing.T) {
	tests := []struct {
		name        string
		cookieState string
		queryState  string
		expected    bool
	}{
		{
			name:        "Valid matching states",
			cookieState: "random-state-1234567890",
			queryState:  "random-state-1234567890",
			expected:    true,
		},
		{
			name:        "Mismatched states",
			cookieState: "state-abc",
			queryState:  "state-xyz",
			expected:    false,
		},
		{
			name:        "Empty cookie state",
			cookieState: "",
			queryState:  "state-xyz",
			expected:    false,
		},
		{
			name:        "Empty query state",
			cookieState: "state-abc",
			queryState:  "",
			expected:    false,
		},
		{
			name:        "Both empty states",
			cookieState: "",
			queryState:  "",
			expected:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ValidateState(tc.cookieState, tc.queryState)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
