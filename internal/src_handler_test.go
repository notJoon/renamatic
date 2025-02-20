package internal

import "testing"

func TestGnoFileProcessor_ShouldProcess(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "basic .gno file",
			path:     "test.gno",
			expected: true,
		},
		{
			name:     "extended .gnoA file",
			path:     "test.gnoA",
			expected: true,
		},
		{
			name:     "extended .gnoXXXXXX file",
			path:     "test.gnoXXXXXX",
			expected: true,
		},
		{
			name:     "non .gno file",
			path:     "test.go",
			expected: false,
		},
		{
			name:     "file without extension",
			path:     "testfile",
			expected: false,
		},
		{
			name:     ".gno file with path",
			path:     "/path/to/test.gno",
			expected: true,
		},
		{
			name:     ".gnoA file with path",
			path:     "/path/to/test.gnoA",
			expected: true,
		},
	}

	processor := newGnoFileProcessor(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ShouldProcess(tt.path)
			if result != tt.expected {
				t.Errorf("ShouldProcess(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
