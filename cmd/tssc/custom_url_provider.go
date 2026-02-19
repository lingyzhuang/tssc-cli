package main

import (
	"context"
	"fmt"

	"github.com/redhat-appstudio/helmet/api"
)

// CustomURLProvider implements integrations.URLProvider by building
// GitHub App URLs from the cluster's OpenShift ingress domain.
type CustomURLProvider struct{}

// GetCallbackURL is set to target Developer Hub
func (CustomURLProvider) GetCallbackURL(ctx context.Context, ic api.IntegrationContext) (string, error) {
	ingressDomain, err := ic.GetOpenShiftIngressDomain(ctx)
	if err != nil {
		return "", fmt.Errorf("ingress domain unavailable (non-OpenShift cluster); "+
			"provide --callback-url explicitly: %w", err)
	}
	namespace, err := ic.GetProductNamespace("Developer Hub")
	if err != nil {
		return "", fmt.Errorf("product unavailable: %w", err)
	}
	return fmt.Sprintf("https://backstage-developer-hub-%s.%s/api/auth/github/handler/frame", namespace, ingressDomain), nil
}

// GetHomepageURL is set to target Developer Hub
func (CustomURLProvider) GetHomepageURL(ctx context.Context, ic api.IntegrationContext) (string, error) {
	ingressDomain, err := ic.GetOpenShiftIngressDomain(ctx)
	if err != nil {
		return "", fmt.Errorf("ingress domain unavailable (non-OpenShift cluster); "+
			"provide --homepage-url explicitly: %w", err)
	}
	namespace, err := ic.GetProductNamespace("Developer Hub")
	if err != nil {
		return "", fmt.Errorf("product unavailable: %w", err)
	}
	return fmt.Sprintf("https://backstage-developer-hub-%s.%s", namespace, ingressDomain), nil
}

// GetWebhookURL is set to target Tekton Pipelines as Code
func (CustomURLProvider) GetWebhookURL(ctx context.Context, ic api.IntegrationContext) (string, error) {
	ingressDomain, err := ic.GetOpenShiftIngressDomain(ctx)
	if err != nil {
		return "", fmt.Errorf("ingress domain unavailable (non-OpenShift cluster); "+
			"provide --webhook-url explicitly: %w", err)
	}
	return fmt.Sprintf(
		"https://pipelines-as-code-controller-openshift-pipelines.%s",
		ingressDomain,
	), nil
}
