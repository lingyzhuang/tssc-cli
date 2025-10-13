package subcmd

import (
	"log/slog"

	"github.com/redhat-appstudio/tssc-cli/pkg/config"
	"github.com/redhat-appstudio/tssc-cli/pkg/integration"
	"github.com/redhat-appstudio/tssc-cli/pkg/k8s"

	"github.com/spf13/cobra"
)

// IntegrationTrustification is the sub-command for the "integration trustification",
// responsible for creating and updating the Trustification integration secret.
type IntegrationTrustification struct {
	cmd         *cobra.Command           // cobra command
	logger      *slog.Logger             // application logger
	cfg         *config.Config           // installer configuration
	kube        *k8s.Kube                // kubernetes client
	integration *integration.Integration // integration instance
}

var _ Interface = &IntegrationTrustification{}

const trustificationIntegrationLongDesc = `
Manages the Trustification integration with TSSC, by storing the required
credentials required by the TSSC services to interact with Trustification.

The credentials are stored in a Kubernetes Secret in the configured namespace
for RHDH.
`
const trustificationIntegrationExample = `
  $ tssc integration trustification \
	  --bombastic-api-url "https://sbom.trustification.dev" \
	  --oidc-issuer-url "https://sso.trustification.dev/realms/chicken" \
	  --oidc-client-id "REDACTED" \
	  --oidc-client-secret "REDACTED"
`

// Cmd exposes the cobra instance.
func (t *IntegrationTrustification) Cmd() *cobra.Command {
	return t.cmd
}

// Complete is a no-op in this case.
func (t *IntegrationTrustification) Complete(args []string) error {
	var err error
	t.cfg, err = bootstrapConfig(t.cmd.Context(), t.kube)
	return err
}

// Validate checks if the required configuration is set.
func (t *IntegrationTrustification) Validate() error {
	return t.integration.Validate()
}

// Run creates or updates the Trustification integration secret.
func (t *IntegrationTrustification) Run() error {
	return t.integration.Create(t.cmd.Context(), t.cfg)
}

// NewIntegrationTrustification creates the sub-command for the "integration
// trustification" responsible to manage the TSSC integrations with the
// Trustification service.
func NewIntegrationTrustification(
	logger *slog.Logger,
	kube *k8s.Kube,
	i *integration.Integration,
) *IntegrationTrustification {
	t := &IntegrationTrustification{
		cmd: &cobra.Command{
			Use:          "trustification [flags]",
			Short:        "Integrates a Trustification instance into TSSC",
			Long:         trustificationIntegrationLongDesc,
			Example:      trustificationIntegrationExample,
			SilenceUsage: true,
		},

		logger:      logger,
		kube:        kube,
		integration: i,
	}
	i.PersistentFlags(t.cmd)
	return t
}
