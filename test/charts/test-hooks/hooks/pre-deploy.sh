#!/bin/bash
set -e

echo "================================"
echo "Pre-Deploy Hook for test-hooks"
echo "================================"
echo "Chart: ${CHART_NAME:-test-hooks}"
echo "Namespace: ${NAMESPACE:-default}"
echo "Release: ${RELEASE_NAME:-test}"
echo ""
echo "Running pre-deployment checks..."
echo "✓ Pre-deploy hook completed successfully"
echo "================================"
