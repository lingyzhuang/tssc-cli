package config

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/redhat-appstudio/tssc-cli/pkg/chartfs"

	"gopkg.in/yaml.v3"
)

// Settings represents a map of configuration settings.
type Settings map[string]interface{}

// ProductSpec represents a map of product name and specification.
type Products []Product

// Spec contains all configuration sections.
type Spec struct {
	// Namespace installer's namespace, where the installer's resources will be
	// deployed. Note, Helm charts deployed by the installer are likely to use a
	// different namespace.
	Namespace string `yaml:"namespace"`
	// Settings contains the configuration for the installer settings.
	Settings Settings `yaml:"settings"`
	// Products contains the configuration for the installer products.
	Products Products `yaml:"products"`
}

// Config root configuration structure.
type Config struct {
	cfs       *chartfs.ChartFS // embedded filesystem
	Installer Spec             `yaml:"tssc"` // root configuration for the installer
	root      yaml.Node        // yaml data representation
}

var (
	// ErrInvalidConfig indicates the configuration content is invalid.
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrEmptyConfig indicates the configuration file is empty.
	ErrEmptyConfig = errors.New("empty configuration")
	// ErrUnmarshalConfig indicates the configuration file structure is invalid.
	ErrUnmarshalConfig = errors.New("failed to unmarshal configuration")
)

// DefaultRelativeConfigPath default relative path to YAML configuration file.
var DefaultRelativeConfigPath = fmt.Sprintf("installer/%s", Filename)

// GetProduct returns a product by name, or an error if the product is not found.
func (c *Config) GetProduct(name string) (*Product, error) {
	for i := range c.Installer.Products {
		if c.Installer.Products[i].Name == name {
			return &c.Installer.Products[i], nil
		}
	}
	return nil, fmt.Errorf("product '%s' not found", name)
}

// GetEnabledProducts returns a map of enabled products.
func (c *Config) GetEnabledProducts() Products {
	enabled := Products{}
	for _, product := range c.Installer.Products {
		if product.Enabled {
			enabled = append(enabled, product)
		}
	}
	return enabled
}

// Validate validates the configuration, checking for missing fields.
func (c *Config) Validate() error {
	root := c.Installer
	// The installer itself must have a namespace.
	if root.Namespace == "" {
		return fmt.Errorf("%w: missing namespace", ErrInvalidConfig)
	}

	// The installer must have a settings section.
	if root.Settings == nil {
		return fmt.Errorf("%w: missing settings", ErrInvalidConfig)
	}

	// Validating the products, making sure every product entry is valid.
	for _, product := range root.Products {
		if err := product.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// DecodeNode returns a struct converted from *yaml.Node
func (c *Config) DecodeNode() error {
	if len(c.root.Content) == 0 {
		return fmt.Errorf("invalid configuration: content is empty")
	}
	if len(c.root.Content[0].Content) == 0 {
		return fmt.Errorf("invalid configuration: missing content")
	}
	root := c.root.Content[0].Content[1]
	if err := root.Decode(&c.Installer); err != nil {
		return err
	}
	return nil
}

// MarshalYAML marshals the Config into a YAML byte array.
func (c *Config) MarshalYAML() ([]byte, error) {
	var buf bytes.Buffer
	if len(c.root.Content) == 0 {
		return nil, fmt.Errorf("invalid configuration format: content is nil or empty")
	}
	buf.WriteString("---\n")
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	defer encoder.Close()
	if err := encoder.Encode(c.root.Content[0]); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalYAML Un-marshals the YAML payload into the Config struct, checking the
// validity of the configuration.
func (c *Config) UnmarshalYAML(payload []byte) error {
	if len(payload) == 0 {
		return ErrEmptyConfig
	}
	if err := yaml.Unmarshal(payload, &c.root); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalConfig, err)
	}
	if err := c.DecodeNode(); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalConfig, err)
	}
	return c.Validate()
}

// String returns this configuration as string, indented with two spaces.
func (c *Config) String() string {
	data, err := c.MarshalYAML()
	if err != nil {
		panic(err)
	}
	return string(data)
}

// NewConfigFromFile returns a new Config instance based on the informed file.
func NewConfigFromFile(cfs *chartfs.ChartFS, configPath string) (*Config, error) {
	c := &Config{cfs: cfs}
	var err error
	payload, err := c.cfs.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err = c.UnmarshalYAML(payload); err != nil {
		return nil, err
	}
	return c, nil
}

// NewConfigFromBytes instantiates a new Config from the bytes payload informed.
func NewConfigFromBytes(payload []byte) (*Config, error) {
	c := &Config{}
	if err := c.UnmarshalYAML(payload); err != nil {
		return nil, err
	}
	return c, nil
}

// NewConfigDefault returns a new Config instance with default values, i.e. the
// configuration payload is loading embedded data.
func NewConfigDefault() (*Config, error) {
	cfs, err := chartfs.NewChartFSForCWD()
	if err != nil {
		return nil, err
	}
	return NewConfigFromFile(cfs, DefaultRelativeConfigPath)
}
