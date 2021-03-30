package handlebars

import (
	"fmt"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the handlebars mixin in porter.yaml
// mixins:
// - handlebars:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
const dockerfileLines = `
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates xz-utils tree
RUN curl -o /tmp/node.tar.xz https://nodejs.org/dist/v14.16.0/node-v14.16.0-linux-x64.tar.xz
RUN tar -xf /tmp/node.tar.xz && \
    mv node-v14.16.0-linux-x64 /usr/local/lib/nodejs && \
    rm /tmp/node.tar.xz

ENV NODEJS_HOME=/usr/local/lib/nodejs
ENV PATH="${NODEJS_HOME}/bin:${PATH}"

RUN npm install -g hbs-cli@%s

RUN mkdir -p /porter/mixins/handlebars && \
    cd /porter/mixins/handlebars && \
    npm install handlebars-helpers
`

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	if input.Config.ClientVersion == "" {
		input.Config.ClientVersion = defaultClientVersion
	}

	fmt.Fprintf(m.Out, dockerfileLines, input.Config.ClientVersion)

	return nil
}
