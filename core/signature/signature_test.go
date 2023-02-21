package signature_test

import (
	"fmt"
	"github.com/HunterPie/Longinus/core/signature"
	"testing"
)

func flattenEdges(edge *signature.PatternEdge) []uint8 {
	bytes := make([]uint8, 0)

	for edge != nil {
		bytes = append(bytes, edge.Value)
	}

	return bytes
}

func assertEqual(edge *signature.PatternEdge, bytes []uint8) bool {
	if edge == nil {
		return len(bytes) == 0
	}

	return assertEqual(edge.Next, bytes[1:])
}

func TestPatternOwner_New(t *testing.T) {
	tests := []struct {
		Name      string
		Signature string
		Expected  []uint8
	}{
		{"test_1", "01 ?? 03 04", []uint8{0x01, 0x00, 0x03, 0x04}},
		{"test_2", "0A 48 ?? ??", []uint8{0x0A, 0x48, 0x00, 0x00}},
	}

	for i, test := range tests {
		name := fmt.Sprintf("TestPatternOwner_New %d", i)

		t.Run(name, func(t *testing.T) {
			actual := signature.New(test.Name, test.Signature)

			if !assertEqual(actual.Pattern, test.Expected) {
				t.Errorf("got %v, want %v", flattenEdges(actual.Pattern), test.Expected)
			}

			if actual.Name != test.Name {
				t.Errorf("got %s, want %s", actual.Name, test.Name)
			}
		})
	}
}
