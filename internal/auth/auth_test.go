package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Define a structure for test cases
	type testCase struct {
		name        string
		headers     http.Header
		expectedKey string
		expectedErr error
	}

	// Define the test cases
	testCases := []testCase{
		// 1. Success Case
		{
			name: "Success with Valid ApiKey",
			headers: http.Header{
				"Authorization": []string{"ApiKey my-secret-key-123"},
			},
			expectedKey: "my-secret-key-123",
			expectedErr: nil,
		},

		// 2. Missing Header
		{
			name:        "Missing Authorization Header",
			headers:     http.Header{},
			expectedKey: "",
			expectedErr: ErrNoAuthHeaderIncluded, // Expecting the specific error variable
		},

		// 3. Malformed Header - Missing prefix
		{
			name: "Malformed Header - Missing ApiKey prefix",
			headers: http.Header{
				"Authorization": []string{"Bearer my-secret-key-123"},
			},
			expectedKey: "",
			// We check for the error string for unexported errors
			expectedErr: errors.New("malformed authorization header"),
		},

		// Bonus Malformed Test: Only one part
		{
			name: "Malformed Header - Only one part",
			headers: http.Header{
				"Authorization": []string{"ApiKey"},
			},
			expectedKey: "",
			expectedErr: errors.New("malformed authorization header"),
		},
	}

	// Iterate and run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualKey, actualErr := GetAPIKey(tc.headers)

			// Check for Key Mismatch
			if actualKey != tc.expectedKey {
				t.Errorf("expected key %q, but got %q", tc.expectedKey, actualKey)
			}

			// Check for Error Mismatch
			if tc.expectedErr == nil && actualErr != nil {
				t.Errorf("expected nil error, but got %v", actualErr)
			} else if tc.expectedErr != nil && actualErr == nil {
				t.Errorf("expected error %v, but got nil", tc.expectedErr)
			} else if tc.expectedErr != nil && actualErr != nil && actualErr.Error() != tc.expectedErr.Error() {
				// Handle specific error value check for ErrNoAuthHeaderIncluded
				if tc.expectedErr != ErrNoAuthHeaderIncluded || actualErr != ErrNoAuthHeaderIncluded {
					t.Errorf("expected error message %q, but got %q", tc.expectedErr.Error(), actualErr.Error())
				}
			}
		})
	}
}
