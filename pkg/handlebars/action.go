package handlebars

import (
	"path"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MarshalYAML converts the action back to a YAML representation
// install:
//   handlebars:
//     ...
func (a Action) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{a.Name: a.Steps}, nil
}

// MakeSteps builds a slice of Step for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - handlebars: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

func (a Action) SetDefaults() {
	for _, s := range a.Steps {
		s.SetDefaults()
	}
}

type Step struct {
	*Instruction `yaml:"handlebars"`
}

// Actions is a set of actions, and the steps, passed from Porter.
type Actions []Action

// UnmarshalYAML takes chunks of a porter.yaml file associated with this mixin
// and populates it on the current action set.
// install:
//   handlebars:
//     ...
//   handlebars:
//     ...
// upgrade:
//   handlebars:
//     ...
func (a *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, Action{})
	if err != nil {
		return err
	}

	for actionName, action := range results {
		for _, result := range action {
			s := result.(*[]Step)
			*a = append(*a, Action{
				Name:  actionName,
				Steps: *s,
			})
		}
	}
	return nil
}

var _ builder.HasOrderedArguments = Instruction{}
var _ builder.ExecutableStep = Instruction{}

type Instruction struct {
	Template       string `yaml:"template"`
	Destination    string `yaml:"destination,omitempty"`
	Description    string `yaml:"description"`
	SuppressOutput bool   `yaml:"suppress-output,omitempty"`
}

// hbs -s --helper /myhelpers.js --data /data.json -- /template.yaml
func (s Instruction) GetCommand() string {
	return "hbs"
}

func (s Instruction) GetArguments() []string {
	return nil
}

func (s Instruction) GetSuffixArguments() []string {
	return []string{
		"--",
		s.Template,
	}
}

func (s Instruction) GetFlags() builder.Flags {
	ext := strings.TrimPrefix(path.Ext(s.Template), ".")
	return builder.Flags{
		builder.NewFlag("helper", helperScript),
		builder.NewFlag("data", dataFile),
		builder.NewFlag("output", tempDestination),
		builder.NewFlag("extension", ext),
	}
}

func (s Instruction) SuppressesOutput() bool {
	return s.SuppressOutput
}

func (s *Instruction) SetDefaults() {
	if s.Destination == "" {
		s.Destination = s.Template
	}
}

func (s Instruction) getTempPath() string {
	return path.Join(tempDestination, path.Base(s.Template))
}
