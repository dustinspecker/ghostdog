package rule

type Rule struct {
	Name     string
	Commands []string
	Sources  []string
	Outputs  []string
	Parents  []*Rule
	Children []*Rule
}
