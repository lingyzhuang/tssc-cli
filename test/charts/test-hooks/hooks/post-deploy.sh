#!/bin/bash
set -e

echo "================================"
echo "Post-Deploy Hook for test-hooks"
echo "================================"
echo "Chart: ${CHART_NAME:-test-hooks}"
echo "Namespace: ${NAMESPACE:-default}"
echo "Release: ${RELEASE_NAME:-test}"
echo ""
echo "Running post-deployment tasks..."
echo "✓ Post-deploy hook completed successfully"
echo "================================"
