package handlebars

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input-all-fields.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	assert.Equal(t, "install", action.Name)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0]
	assert.Equal(t, "Render Template", step.Description, "Invalid description")

	assert.Equal(t, "stuff.yaml.hbs", step.Template, "Invalid source")
	assert.Equal(t, "other-stuff.yaml", step.Destination, "Invalid destination")
	assert.Equal(t, "mydata.json", step.Data, "Invalid data")
}

func TestAction_SetDefaults(t *testing.T) {
	a := Action{
		Name: "install",
		Steps: []Step{
			{
				&Instruction{
					Template: "foo.yaml",
				},
			},
		},
	}

	a.SetDefaults()

	step := a.Steps[0]
	assert.Equal(t, step.Template, step.Destination, "invalid destination")
}

func TestAction_getTempPath(t *testing.T) {
	s := Step{
		&Instruction{
			Template: "foo.yaml",
		},
	}

	got := s.getTempPath()
	want := "/porter/mixins/handlebars/output/foo.yaml"
	assert.Equal(t, want, got)
}
