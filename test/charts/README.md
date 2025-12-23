# Test Helm Charts

This directory contains test Helm charts designed to validate the functionality of the TSSC-CLI installer components, particularly the resolver, engine, chartfs, and installer modules.

## Overview

The test charts are organized by functionality and test different aspects of the installer architecture:

- **Dependency Resolution**: Linear and complex dependency chains
- **Weight-Based Ordering**: Charts with different installation weights
- **Product Namespaces**: Charts that define and use product namespaces
- **Integration Requirements**: Charts with CEL-based integration requirements
- **Hooks**: Charts with pre-deploy and post-deploy hooks
- **Template Functions**: Charts that test custom template functions

## Chart Descriptions

### 1. Simple Test Chart

**Chart**: `test-simple`

**Purpose**: Basic chart with no dependencies for testing fundamental chart loading and rendering.

**Features**:
- No dependencies
- Simple ConfigMap resource
- Basic values
- Helm test pod

**Use Cases**:
- Test ChartFS can load charts
- Test basic template rendering
- Test Helm install/verify workflow

---

### 2. Dependency Chain Charts

**Charts**: `test-dep-a`, `test-dep-b`, `test-dep-c`

**Purpose**: Test linear dependency resolution (A → B → C).

**Dependency Flow**:
```
test-dep-a (base, no dependencies)
    ↓
test-dep-b (depends on test-dep-a)
    ↓
test-dep-c (depends on test-dep-b)
```

**Annotations**:
- `test-dep-a`: None
- `test-dep-b`: `tssc.redhat-appstudio.github.com/depends-on: test-dep-a`
- `test-dep-c`: `tssc.redhat-appstudio.github.com/depends-on: test-dep-b`

**Use Cases**:
- Test topology resolution with linear dependencies
- Test correct installation order
- Verify transitive dependency handling

---

### 3. Diamond Dependency Pattern Charts

**Charts**: `test-diamond-base`, `test-diamond-left`, `test-diamond-right`, `test-diamond-top`

**Purpose**: Test complex dependency resolution with diamond pattern.

**Dependency Flow**:
```
        test-diamond-base
              /    \
    test-diamond-left  test-diamond-right
              \    /
         test-diamond-top
```

**Annotations**:
- `test-diamond-base`: None
- `test-diamond-left`: `depends-on: test-diamond-base`
- `test-diamond-right`: `depends-on: test-diamond-base`
- `test-diamond-top`: `depends-on: test-diamond-left,test-diamond-right`

**Use Cases**:
- Test diamond dependency resolution
- Test prevention of circular dependencies
- Test multiple parallel dependencies
- Verify base chart is installed only once

---

### 4. Weight-Based Ordering Charts

**Charts**: `test-weight-low`, `test-weight-medium`, `test-weight-high`

**Purpose**: Test weight-based ordering of charts.

**Weights**:
- `test-weight-low`: weight `10` (installed first)
- `test-weight-medium`: weight `50` (installed second)
- `test-weight-high`: weight `90` (installed last)

**Annotations**:
- All charts have `tssc.redhat-appstudio.github.com/weight` annotation

**Use Cases**:
- Test weight-based ordering
- Test weight overrides default ordering
- Verify charts are installed in correct order

**Expected Order**:
1. `test-weight-low` (weight 10)
2. `test-weight-medium` (weight 50)
3. `test-weight-high` (weight 90)

---

### 5. Product Namespace Charts

**Charts**: `test-product`, `test-use-product-namespace`

**Purpose**: Test product namespace assignment and sharing.

**Annotations**:
- `test-product`: `product-name: "Test Product"` (defines a product)
- `test-use-product-namespace`: `use-product-namespace: "Test Product"` (uses the product's namespace)

**Dependency**:
- `test-use-product-namespace` depends on `test-product`

**Use Cases**:
- Test product namespace assignment
- Test `use-product-namespace` functionality
- Verify namespace sharing between charts

**Expected Behavior**:
- `test-product` creates namespace `test-product-ns`
- `test-use-product-namespace` resources are created in `test-product-ns`

---

### 6. Integration Requirements Charts

**Charts**: `test-integration-provider`, `test-integration-consumer-simple`, `test-integration-consumer-complex`

**Purpose**: Test integration requirements and CEL expression evaluation.

**Integration Provider**:
- `test-integration-provider`: Provides `test-github` and `test-quay` integrations
- Annotation: `integrations-provided: "test-github,test-quay"`

**Simple Consumer**:
- `test-integration-consumer-simple`: Requires `test-github`
- Annotation: `integrations-required: "test-github"`

**Complex Consumer**:
- `test-integration-consumer-complex`: Requires `(test-github OR test-gitlab) AND test-quay`
- Annotation: `integrations-required: "(test-github || test-gitlab) && test-quay"`

**Use Cases**:
- Test CEL expression validation
- Test integration requirement checking
- Test complex expressions with logical operators
- Verify integration dependency resolution

---

### 7. Hooks Test Chart

**Chart**: `test-hooks`

**Purpose**: Test pre-deploy and post-deploy hook execution.

**Hooks**:
- `hooks/pre-deploy.sh`: Runs before chart deployment
- `hooks/post-deploy.sh`: Runs after chart deployment

**Use Cases**:
- Test pre-deploy hook execution
- Test post-deploy hook execution
- Test environment variable passing to hooks
- Verify hook scripts can access chart metadata

**Environment Variables Available to Hooks**:
- `CHART_NAME`: Name of the chart
- `NAMESPACE`: Target namespace
- `RELEASE_NAME`: Helm release name

---

### 8. Template Functions Test Chart

**Chart**: `test-template-functions`

**Purpose**: Test custom template functions provided by the engine.

**Custom Functions Tested**:
- `toYaml`: Convert Go data structures to YAML
- `toJson`: Convert Go data structures to JSON
- `fromYaml`: Parse YAML strings into Go structures
- `fromJson`: Parse JSON strings into Go structures

**Template Features**:
- Helper templates (`_helpers.tpl`)
- Complex nested data structures
- JSON and YAML data transformations

**Use Cases**:
- Test engine custom functions
- Test complex value transformations
- Verify template rendering with custom functions
- Test Sprig functions integration

---

## Chart Annotations Reference

The following annotations are used by the TSSC-CLI resolver (defined in `pkg/resolver/annotations.go`):

| Annotation | Description | Example |
|------------|-------------|---------|
| `tssc.redhat-appstudio.github.com/product-name` | Product this chart belongs to | `"Test Product"` |
| `tssc.redhat-appstudio.github.com/depends-on` | Comma-separated list of chart dependencies | `"chart-a,chart-b"` |
| `tssc.redhat-appstudio.github.com/weight` | Integer weight for ordering (higher = later) | `"50"` |
| `tssc.redhat-appstudio.github.com/use-product-namespace` | Use another product's namespace | `"Product Name"` |
| `tssc.redhat-appstudio.github.com/integrations-provided` | Integrations this chart provides | `"github,quay"` |
| `tssc.redhat-appstudio.github.com/integrations-required` | CEL expression for required integrations | `"(github \|\| gitlab) && quay"` |

## Testing Workflows

### Unit Testing

Test charts can be used in Go unit tests:

```go
package resolver_test

import (
    "testing"
    "github.com/redhat-appstudio/tssc-cli/pkg/chartfs"
    "github.com/redhat-appstudio/tssc-cli/pkg/resolver"
    o "github.com/onsi/gomega"
)

func TestDependencyChain(t *testing.T) {
    g := o.NewWithT(t)

    // Load test charts
    cfs, err := chartfs.NewChartFS("../../test/charts")
    g.Expect(err).To(o.Succeed())

    charts, err := cfs.GetAllCharts()
    g.Expect(err).To(o.Succeed())

    // Create collection and topology
    collection := resolver.NewCollection(charts)
    topology := resolver.NewTopology(collection, "default-ns")

    // Add test-dep-c (should pull in test-dep-b and test-dep-a)
    err = topology.Add("test-dep-c")
    g.Expect(err).To(o.Succeed())

    // Verify order: test-dep-a, test-dep-b, test-dep-c
    deps := topology.GetAll()
    g.Expect(deps).To(o.HaveLen(3))
    g.Expect(deps[0].Chart().Name()).To(o.Equal("test-dep-a"))
    g.Expect(deps[1].Chart().Name()).To(o.Equal("test-dep-b"))
    g.Expect(deps[2].Chart().Name()).To(o.Equal("test-dep-c"))
}
```

### Integration Testing

Test charts can be used with the full installer:

```go
package installer_test

import (
    "testing"
    "github.com/redhat-appstudio/tssc-cli/pkg/chartfs"
    "github.com/redhat-appstudio/tssc-cli/pkg/installer"
    "github.com/redhat-appstudio/tssc-cli/pkg/k8s"
)

func TestHooksExecution(t *testing.T) {
    cfs, _ := chartfs.NewChartFS("../../test/charts")
    chart, _ := cfs.GetChart("test-hooks")

    kube := k8s.NewKube(clientset, "default")
    inst := installer.New(kube, chart, namespace, releaseName)

    // This should execute pre-deploy hook, install chart, then post-deploy hook
    err := inst.Install()
    // Assert hooks were executed...
}
```

### Manual Testing

You can also test charts manually with Helm:

```bash
# Lint a chart
helm lint test/charts/test-simple

# Template a chart
helm template test-release test/charts/test-simple

# Install a chart
helm install test-release test/charts/test-simple

# Run Helm tests
helm test test-release

# Uninstall
helm uninstall test-release
```

## Test Scenarios

### Scenario 1: Linear Dependency Resolution

**Charts**: `test-dep-a`, `test-dep-b`, `test-dep-c`

**Expected Behavior**:
1. When requesting `test-dep-c`, resolver should include all three charts
2. Installation order: `test-dep-a` → `test-dep-b` → `test-dep-c`
3. Each chart's ConfigMap should reference its dependency

### Scenario 2: Diamond Dependency Resolution

**Charts**: `test-diamond-*`

**Expected Behavior**:
1. When requesting `test-diamond-top`, resolver includes all four charts
2. `test-diamond-base` is installed only once
3. `test-diamond-left` and `test-diamond-right` can be installed in any order (both depend on base)
4. `test-diamond-top` is installed last

### Scenario 3: Weight-Based Ordering

**Charts**: `test-weight-*`

**Expected Behavior**:
1. Charts are ordered by weight regardless of name
2. Order: low (10) → medium (50) → high (90)

### Scenario 4: Integration Requirements

**Charts**: `test-integration-*`

**Expected Behavior**:
1. Provider chart makes integrations available
2. Simple consumer validates `test-github` exists
3. Complex consumer evaluates CEL: `(test-github || test-gitlab) && test-quay`
4. Should pass if `test-github` AND `test-quay` are available

### Scenario 5: Product Namespace Sharing

**Charts**: `test-product`, `test-use-product-namespace`

**Expected Behavior**:
1. `test-product` creates namespace `test-product-ns`
2. `test-use-product-namespace` resources go into `test-product-ns`
3. Both charts share the same namespace

## Directory Structure

```
test/charts/
├── README.md                              # This file
├── test-simple/                           # Simple test chart
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── configmap.yaml
│       ├── NOTES.txt
│       └── tests/
│           └── test.yaml
├── test-dep-a/                            # Dependency chain: A
├── test-dep-b/                            # Dependency chain: B (depends on A)
├── test-dep-c/                            # Dependency chain: C (depends on B)
├── test-diamond-base/                     # Diamond pattern: base
├── test-diamond-left/                     # Diamond pattern: left branch
├── test-diamond-right/                    # Diamond pattern: right branch
├── test-diamond-top/                      # Diamond pattern: top (depends on left + right)
├── test-weight-low/                       # Weight 10
├── test-weight-medium/                    # Weight 50
├── test-weight-high/                      # Weight 90
├── test-product/                          # Defines a product
├── test-use-product-namespace/            # Uses another product's namespace
├── test-integration-provider/             # Provides integrations
├── test-integration-consumer-simple/      # Simple integration requirement
├── test-integration-consumer-complex/     # Complex CEL integration requirement
├── test-hooks/                            # Pre/post-deploy hooks
│   ├── Chart.yaml
│   ├── values.yaml
│   ├── templates/
│   └── hooks/
│       ├── pre-deploy.sh
│       └── post-deploy.sh
└── test-template-functions/               # Custom template functions
    ├── Chart.yaml
    ├── values.yaml
    └── templates/
        ├── configmap.yaml
        ├── _helpers.tpl
        └── NOTES.txt
```

## Contributing

When adding new test charts:

1. **Follow the naming convention**: `test-<feature>`
2. **Keep charts simple**: Focus on testing one specific feature
3. **Add documentation**: Update this README with the new chart's purpose and use cases
4. **Include tests**: Add Helm test pods where appropriate
5. **Use annotations**: Leverage TSSC-CLI annotations for resolver testing

## Related Documentation

- **Resolver Documentation**: `pkg/resolver/` - Dependency topology resolution
- **ChartFS Documentation**: `pkg/chartfs/` - Chart filesystem abstraction
- **Engine Documentation**: `pkg/engine/` - Template rendering engine
- **Installer Documentation**: `pkg/installer/` - Chart installation logic
- **Annotations Reference**: `pkg/resolver/annotations.go` - Chart annotation constants
