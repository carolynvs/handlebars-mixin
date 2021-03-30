package handlebars

import (
	"encoding/json"

	"github.com/gobuffalo/packr/v2"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const dataFile = "/porter/mixins/handlebars/template-data.json"

func (m *Mixin) loadAction() (*Action, error) {
	var action Action
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute() error {
	action, err := m.loadAction()
	if err != nil {
		return err
	}

	action.SetDefaults()
	err = m.writeTemplateData()
	if err != nil {
		return err
	}
	err = m.writeHelpers()
	if err != nil {
		return err
	}

	_, err = builder.ExecuteSingleStepAction(m.Context, action)
	if err != nil {
		return err
	}

	for _, s := range action.Steps {
		err := m.FileSystem.Rename(s.getTempPath(), s.Destination)
		if err != nil {
			return errors.Wrapf(err, "error overwriting existing file %s with the rendered template", s.Template)
		}
	}

	return nil
}

func (m *Mixin) writeTemplateData() error {
	data := map[string]interface{}{
		"name": "carolyn was here",
	}

	dump, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "could not marshal template data as json")
	}

	err = m.FileSystem.WriteFile(dataFile, dump, 0600)
	return errors.Wrap(err, "error writing template data")
}

func (m *Mixin) writeHelpers() error {
	t := packr.New("helpers", "./helpers")

	b, err := t.Find("handlebars-helpers.js")
	if err != nil {
		return errors.Wrap(err, "error loading handlebars helper script template")
	}

	err = m.FileSystem.WriteFile(helperScript, b, 0600)
	return errors.Wrap(err, "error writing handlebars helper script")
}
