package dag

import (
	"errors"
	"fmt"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

var (
	ErrEdgeAlreadyExists      = errors.New("edge already exists")
	ErrRuleNameAlreadyDefined = errors.New("rule name already exists")
	ErrRuleNotFound           = errors.New("rule not found")
)

type Dag struct {
	Rules map[string]*rule.Rule
}

func NewDag() Dag {
	return Dag{
		Rules: make(map[string]*rule.Rule),
	}
}

func (dag Dag) AddDependency(parentId, childId string) error {
	parent, ok := dag.Rules[parentId]
	if !ok {
		return fmt.Errorf("%s: %w", parentId, ErrRuleNotFound)
	}

	child, ok := dag.Rules[childId]
	if !ok {
		return fmt.Errorf("%s: %w", childId, ErrRuleNotFound)
	}

	for _, child := range parent.Children {
		if child == child {
			return fmt.Errorf("%s->%s: %w", parentId, childId, ErrEdgeAlreadyExists)
		}
	}

	parent.Children = append(parent.Children, child)
	child.Parents = append(child.Parents, parent)

	return nil
}

func (dag *Dag) AddRule(rule *rule.Rule) error {
	_, ok := dag.Rules[rule.Name]
	if ok {
		return fmt.Errorf("rule %s: %w", rule.Name, ErrRuleNameAlreadyDefined)
	}

	dag.Rules[rule.Name] = rule

	return nil
}

func (dag Dag) GetSources() map[string]*rule.Rule {
	sources := map[string]*rule.Rule{}

	for _, rule := range dag.Rules {
		if len(rule.Parents) == 0 {
			sources[rule.Name] = rule
		}
	}

	return sources
}
