package tree_test

import (
	"fmt"
	"github.com/HunterPie/Longinus/core/signature"
	"github.com/HunterPie/Longinus/core/tree"
	"reflect"
	"testing"
)

func TestPatternTree_FindPattern(t *testing.T) {
	owners := []*signature.PatternOwner{
		signature.New("test_1", "01 ?? 03 04"),
		signature.New("test_2", "02 03 04"),
		signature.New("test_3", "01 02 03 ??"),
	}
	tests := []struct {
		Owners   []*signature.PatternOwner
		TestCase []uint8
		Expected []*signature.PatternOwner
	}{
		{owners, []uint8{0x01, 0x02, 0x03, 0x04}, []*signature.PatternOwner{owners[0], owners[2]}},
		{owners, []uint8{0x02, 0x03, 0x04}, []*signature.PatternOwner{owners[1]}},
		{owners, []uint8{0x01, 0x04, 0x05}, []*signature.PatternOwner{}},
	}

	for i, test := range tests {
		name := fmt.Sprintf("TestPatternTreeNode_Find %d", i)

		patternTree := tree.New(test.Owners...)

		t.Run(name, func(t *testing.T) {
			actual := patternTree.FindPattern(test.TestCase)

			if !reflect.DeepEqual(actual, test.Expected) {
				t.Errorf("got %v, want %v", actual, test.Expected)
			}
		})
	}
}
