package mcptools

import (
	"context"
	"fmt"
	"strings"

	"github.com/redhat-appstudio/tssc-cli/pkg/constants"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type IntegrationTools struct {
	integrationCmd *cobra.Command // integration subcommand
}

// IntegrationListTool list integrations tool.
const (
	// IntegrationListTool
	IntegrationListTool = constants.AppName + "_integration_list"
	// IntegrationScaffoldTool
	IntegrationScaffoldTool = constants.AppName + "_integration_scaffold"
	// IntegrationStatusTool
	IntegrationStatusTool = constants.AppName + "_integration_status"
	// MissingIntegrations
	MissingIntegrations = "integration"
)

func (i *IntegrationTools) listHandler(
	ctx context.Context,
	ctr mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	var output strings.Builder
	output.WriteString(fmt.Sprintf(`
# Integration Commands

The detailed description of each '%s integration' command is found below.
`,
		constants.AppName,
	))

	for _, subCmd := range i.integrationCmd.Commands() {
		var flagsInfo strings.Builder
		subCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			required := ""
			if _, value := f.Annotations[cobra.BashCompOneRequiredFlag]; value {
				if len(f.Annotations[cobra.BashCompOneRequiredFlag]) > 0 &&
					f.Annotations[cobra.BashCompOneRequiredFlag][0] == "true" {
					required = " (REQUIRED)"
				}
			}

			flagsInfo.WriteString(fmt.Sprintf(
				"  - \"--%s\" %s%s: %s.\n",
				f.Name,
				f.Value.Type(),
				required,
				f.Usage,
			))
		})
		output.WriteString(fmt.Sprintf(`
## '$ %s integration %s'

%s
%s

### Flags

%s
`,
			constants.AppName,
			subCmd.Name(),
			subCmd.Short,
			subCmd.Long,
			flagsInfo.String(),
		))
	}
	return mcp.NewToolResultText(output.String()), nil
}

func (i *IntegrationTools) scaffoldHandler(
	ctx context.Context,
	ctr mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	var output strings.Builder
	// Validate integrations
	if integrations, ok := ctr.GetArguments()[MissingIntegrations].([]any); ok {
		for _, integration := range integrations {
			for _, subCmd := range i.integrationCmd.Commands() {
				if integration == subCmd.Name() {
					output.WriteString(fmt.Sprintf(`
## Integration: %s is missing, please create the integration with following command:

%s
`,
						subCmd.Name(),
						subCmd.Example,
					))
				}
			}
		}
	}

	return mcp.NewToolResultText(output.String()), nil
}

func (i *IntegrationTools) Init(s *server.MCPServer) {
	s.AddTools([]server.ServerTool{{
		Tool: mcp.NewTool(
			IntegrationListTool,
			mcp.WithDescription(`
List the TSSC integrations available for the user. Certain integrations are
required for certain features, make sure to configure the integrations
accordingly.`),
		),
		Handler: i.listHandler,
	},
		{
			Tool: mcp.NewTool(
				IntegrationScaffoldTool,
				mcp.WithDescription(`
Scaffold the configuration required for a specific TSSC integration. The
scaffolded configuration can be used as a reference to create the integration
using the 'tssc integration <name> ...' command.`),
				mcp.WithArray(
					MissingIntegrations,
					mcp.Description(`
The missing integrations for deployment.`,
					),
				),
			),
			Handler: i.scaffoldHandler,
		}}...)
}

func NewIntegrationTools(integrationCmd *cobra.Command) *IntegrationTools {
	return &IntegrationTools{
		integrationCmd: integrationCmd,
	}
}
