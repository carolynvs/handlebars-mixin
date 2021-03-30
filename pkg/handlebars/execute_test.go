package handlebars

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Fake out executing a command
	// It's okay to use os.LookupEnv here because it's running in it's own process, and won't impact running tests in parallel.
	if _, mockCommand := os.LookupEnv(test.MockedCommandEnv); mockCommand {
		fmt.Println("MOCK_COMMAND", mockCommand)
		if expectedCmdEnv, doAssert := os.LookupEnv(test.ExpectedCommandEnv); doAssert {
			fmt.Println("EXPECTED_COMMAND", expectedCmdEnv)

			gotCmd := strings.Join(os.Args[1:len(os.Args)], " ")

			// There may be multiple expected commands, separated by a newline character
			wantCmds := strings.Split(expectedCmdEnv, "\n")

			commandNotFound := true
			for _, wantCmd := range wantCmds {
				if wantCmd == gotCmd {
					commandNotFound = false
				}
			}

			if commandNotFound {
				fmt.Printf("WANT COMMANDS : %q\n", wantCmds)
				fmt.Printf("GOT COMMAND : %q\n", gotCmd)
				os.Exit(127)
			}
		}
		os.Exit(0)
	}

	// Otherwise, run the tests
	os.Exit(m.Run())
}

func TestMixin_Execute(t *testing.T) {
	testcases := []struct {
		name        string // Test case name
		file        string // Path to th test input yaml
		wantCommand string // Full command that you expect to be called based on the input YAML
	}{
		{"action", "testdata/step-input.yaml",
			"hbs --data /porter/mixins/handlebars/template-data.json --extension yaml --helper /porter/mixins/handlebars/handlebars-helpers.js --output /porter/mixins/handlebars/output -- stuff.yaml"},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewTestMixin(t)
			m.TestContext.AddTestFile("testdata/stuff.yaml", "/porter/mixins/handlebars/output/stuff.yaml")

			m.Setenv(test.ExpectedCommandEnv, tc.wantCommand)
			mixinInputB, err := ioutil.ReadFile(tc.file)
			require.NoError(t, err)

			m.In = bytes.NewBuffer(mixinInputB)

			err = m.Execute()
			require.NoError(t, err, "execute failed")
		})
	}
}
