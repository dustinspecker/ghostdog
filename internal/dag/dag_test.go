package dag

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

func TestAddRuleAddRuleToListOfRules(t *testing.T) {
	dag := NewDag()

	ruleA := &rule.Rule{
		Name: "a",
	}
	if err := dag.AddRule(ruleA); err != nil {
		t.Fatalf("should have added ruleA: %s", err)
	}

	if dag.Rules["a"] != ruleA {
		t.Errorf("expected rule a to be ruleA, but was %v", dag.Rules["a"])
	}

	ruleB := &rule.Rule{
		Name: "b",
	}
	if err := dag.AddRule(ruleB); err != nil {
		t.Fatalf("should have added ruleB: %s", err)
	}

	if dag.Rules["b"] != ruleB {
		t.Errorf("expected rule b to be ruleB, but was %v", dag.Rules["b"])
	}
}

func TestAddRuleReturnsErrorIfIdAlreadyUsed(t *testing.T) {
	dag := NewDag()

	ruleA := &rule.Rule{}
	ruleB := &rule.Rule{}

	dag.AddRule(ruleA)

	err := dag.AddRule(ruleB)
	if err == nil {
		t.Fatalf("adding a rule with a used id should return an error")
	}

	if !errors.Is(err, ErrRuleNameAlreadyDefined) {
		t.Error("expected error to be an ErrRuleNameAlreadyDefined")
	}
}

func TestAddDependency(t *testing.T) {
	dag := NewDag()

	parent := rule.Rule{
		Name: "parent",
	}
	child := rule.Rule{
		Name: "child",
	}
	if err := dag.AddRule(&parent); err != nil {
		t.Fatalf("should have added parent: %s", err)
	}
	if err := dag.AddRule(&child); err != nil {
		t.Fatalf("should have added child: %s", err)
	}

	if err := dag.AddDependency("parent", "child"); err != nil {
		t.Fatalf("should have added edge between parent and child: %s", err)
	}

	if !reflect.DeepEqual(parent.Children, []*rule.Rule{&child}) {
		t.Errorf("parent should have child as child, but had %v", parent.Children)
	}

	if !reflect.DeepEqual(child.Parents, []*rule.Rule{&parent}) {
		t.Errorf("child should have parent as parent, but had %v", child.Parents)
	}
}

func TestAddDependencyReturnsErrorWhenparentIdNotFound(t *testing.T) {
	dag := NewDag()

	err := dag.AddDependency("parentId", "childId")
	if err == nil {
		t.Fatal("should return error when parentId not found")
	}

	if !errors.Is(err, ErrRuleNotFound) {
		t.Errorf("expected ErrRuleNotFound error, but got %s", err)
	}

	if !strings.HasPrefix(err.Error(), "parentId:") {
		t.Errorf("expected error message to start with parentId:, but got %s", err.Error())
	}
}

func TestAddDependencyReturnsErrorWhenchildIdNotFound(t *testing.T) {
	dag := NewDag()

	if err := dag.AddRule(&rule.Rule{Name: "parentId"}); err != nil {
		t.Fatalf("expected no error adding parentId rule, but got %s", err)
	}

	err := dag.AddDependency("parentId", "childId")
	if err == nil {
		t.Fatal("should return error when childId not found")
	}

	if !errors.Is(err, ErrRuleNotFound) {
		t.Errorf("expected ErrRuleNotFound error, but got %s", err)
	}

	if !strings.HasPrefix(err.Error(), "childId:") {
		t.Errorf("expected error message to start with childId:, but got %s", err.Error())
	}
}

func TestAddDependencyReturnsErrorWhenEdgeAlreadyExists(t *testing.T) {
	dag := NewDag()

	if err := dag.AddRule(&rule.Rule{Name: "parent"}); err != nil {
		t.Fatalf("expected no error adding parent rule, but got %s", err)
	}

	if err := dag.AddRule(&rule.Rule{Name: "child"}); err != nil {
		t.Fatalf("expected no error adding child rule, but got %s", err)
	}

	if err := dag.AddRule(&rule.Rule{Name: "another_child"}); err != nil {
		t.Fatalf("expected no error adding another_child rule, but got %s", err)
	}

	if err := dag.AddDependency("parent", "child"); err != nil {
		t.Fatal("expected no error adding edge between parent and child")
	}

	if err := dag.AddDependency("parent", "another_child"); err != nil {
		t.Fatal("expected no error adding edge between parent and another_child")
	}

	err := dag.AddDependency("parent", "child")
	if err == nil {
		t.Fatal("expected error when adding edge that already exists")
	}

	if !errors.Is(err, ErrEdgeAlreadyExists) {
		t.Errorf("expected ErrEdgeAlreadyExists error, but got %s", err)
	}
}

func TestGetSources(t *testing.T) {
	dag := NewDag()

	sourceA := &rule.Rule{
		Name: "a",
	}
	sourceB := &rule.Rule{
		Name: "b",
	}
	ruleC := &rule.Rule{
		Name: "c",
	}
	ruleD := &rule.Rule{
		Name: "d",
	}

	dag.AddRule(sourceA)
	dag.AddRule(sourceB)
	dag.AddRule(ruleC)
	dag.AddRule(ruleD)

	dag.AddDependency("a", "c")
	dag.AddDependency("b", "d")

	sources := dag.GetSources()

	expectedSources := []*rule.Rule{
		sourceA,
		sourceB,
	}

	if !reflect.DeepEqual(sources, expectedSources) {
		t.Errorf("expected sources to be a map containing only sourceA and sourceB, but got %v", sources)
	}
}
