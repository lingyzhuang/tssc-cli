package subcmd

import (
	"log/slog"

	"github.com/redhat-appstudio/tssc-cli/pkg/config"
	"github.com/redhat-appstudio/tssc-cli/pkg/integration"
	"github.com/redhat-appstudio/tssc-cli/pkg/k8s"

	"github.com/spf13/cobra"
)

// IntegrationArtifactory is the sub-command for the "integration artifactory",
// responsible for creating and updating the Artifactory integration secret.
type IntegrationArtifactory struct {
	cmd         *cobra.Command           // cobra command
	logger      *slog.Logger             // application logger
	cfg         *config.Config           // installer configuration
	kube        *k8s.Kube                // kubernetes client
	integration *integration.Integration // integration instance

	apiToken         string // web API token
	dockerconfigjson string // credentials to push/pull from the registry
}

var _ Interface = &IntegrationArtifactory{}

const artifactoryIntegrationLongDesc = `
Manages the artifactory integration with TSSC, by storing the required
credentials required by the TSSC services to interact with artifactory.

The credentials are stored in a Kubernetes Secret in the configured namespace
for RHDH.
`
const artifactoryIntegrationExample = `
  $ tssc integration artifactory \
	  --dockerconfigjson '{ "auths": { "artifactory.io/my-repository": { "auth": "REDACTED" } } }' \
	  --token "REDACTED" \
	  --url 'https://artifactory.io'
`

// Cmd exposes the cobra instance.
func (a *IntegrationArtifactory) Cmd() *cobra.Command {
	return a.cmd
}

// Complete is a no-op in this case.
func (a *IntegrationArtifactory) Complete(args []string) error {
	var err error
	a.cfg, err = bootstrapConfig(a.cmd.Context(), a.kube)
	return err
}

// Validate checks if the required configuration is set.
func (a *IntegrationArtifactory) Validate() error {
	return a.integration.Validate()
}

// Run creates or updates the Artifactory integration secret.
func (a *IntegrationArtifactory) Run() error {
	return a.integration.Create(a.cmd.Context(), a.cfg)
}

// NewIntegrationArtifactory creates the sub-command for the "integration artifactory"
// responsible to manage the TSSC integrations with a Artifactory image registry.
func NewIntegrationArtifactory(
	logger *slog.Logger,
	kube *k8s.Kube,
	i *integration.Integration,
) *IntegrationArtifactory {
	a := &IntegrationArtifactory{
		cmd: &cobra.Command{
			Use:          "artifactory [flags]",
			Short:        "Integrates a Artifactory instance into TSSC",
			Long:         artifactoryIntegrationLongDesc,
			Example:      artifactoryIntegrationExample,
			SilenceUsage: true,
		},

		logger:      logger,
		kube:        kube,
		integration: i,
	}
	i.PersistentFlags(a.cmd)
	return a
}
