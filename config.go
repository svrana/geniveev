package geniveev

import "fmt"

// This is probably a poor name if we've only got this structure
// to allow other rendering engines.
type GeneratorConfig struct {
	Code string
}

type Filename string

type TemplateValues map[string]*string

type TemplateConfig struct {
	TemplateConfigMap map[Filename]GeneratorConfig
	TemplateValues    TemplateValues
}

func (tc *TemplateConfig) String() string {
	return fmt.Sprintf("configMap: %+v, values: %+v", tc.TemplateConfigMap, tc.TemplateValues)
}

func NewTemplateConfigEmpty() *TemplateConfig {
	return &TemplateConfig{
		TemplateConfigMap: make(map[Filename]GeneratorConfig),
		TemplateValues:    make(map[string]*string),
	}
}

type Config struct {
	Generator map[string]*TemplateConfig
}

func NewConfig() *Config {
	return &Config{
		Generator: make(map[string]*TemplateConfig),
	}
}
