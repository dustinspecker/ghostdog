package rule

import (
	"sort"
	"testing"
)

func TestSortByName(t *testing.T) {
	rules := []*Rule{
		&Rule{Name: "test"},
		&Rule{Name: "build"},
		&Rule{Name: "publish"},
	}

	sort.Sort(ByName(rules))

	if rules[0].Name != "build" {
		t.Errorf("expected build to be first but got %s", rules[0].Name)
	}
	if rules[1].Name != "publish" {
		t.Errorf("expected publish to be first but got %s", rules[1].Name)
	}
	if rules[2].Name != "test" {
		t.Errorf("expected test to be first but got %s", rules[2].Name)
	}
}
