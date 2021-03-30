//go:generate packr2
package handlebars

import (
	"get.porter.sh/porter/pkg/context"
)

const (
	defaultClientVersion string = "1.4.0"
	tempDestination      string = "/porter/mixins/handlebars/output"
	helperScript         string = "/porter/mixins/handlebars/handlebars-helpers.js"
)

type Mixin struct {
	*context.Context
}

// New mixin client.
func New() (*Mixin, error) {
	return &Mixin{
		Context: context.New(),
	}, nil
}
