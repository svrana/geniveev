package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// This is probably a poor name if we've only got this structure
// to allow other rendering engines.
type GeneratorConfig struct {
	Code string
}

type Filename string

type TemplateConfig struct {
	TemplateConfigMap map[Filename]GeneratorConfig
	TemplateValues    map[string]*string
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

var cfgFile string = ".geniveev.toml"
var config *Config = NewConfig()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geniveev",
	Short: "An easy to use, language-agnostic code generation tool, powered by Go templates and an intuitive CLI.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func start(cmd string, args []string, templateConfig *TemplateConfig) {
	//fmt.Printf("got cmd: %s, args: %s, templateConfig: %s", cmd, args, templateConfig)
	// k := cfg[cmd]
	// if userDef, ok := k.(map[string]interface{}); ok {
	// 	//fmt.Printf("got userDef: %s\n", userDef)
	// 	for k := range userDef {
	// 		fmt.Printf("got key of userDef: %s\n", k)
	// 	}
	// }
	// for filename, generatorConfig := range templateConfig.TemplateConfigMap {
	//
	// 	//fmt.Printf("got filename: %s, code: %s\n", filename, generatorConfig.Code)
	// }

	// go through each key of service_stubs
	// substitute all parameterized values with those from templateValues
	// create directories as approrpriate
	// create go template, and write it to the file
}

// initConfig reads in config file and ENV variables if set.
func init() {
	// key splitting in cobra is borked with toml and keys containing periods, so just handle
	// the file reading ourselves here.
	b, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .geniveev.toml: %s", err)
		os.Exit(1)
	}

	cfg := make(map[string]interface{})

	if err := toml.Unmarshal(b, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal toml: %s\n", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`{(\w+)}`)

	for name, v := range cfg {
		config.Generator[name] = NewTemplateConfigEmpty()

		//fmt.Printf("key: %s\n", k)
		newCmd := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Generate %s", name),
			Run: func(_ *cobra.Command, args []string) {
				start(name, args, config.Generator[name])
			},
		}

		arguments := make(map[string]bool)
		var valueMap map[string]interface{}
		var ok bool

		if valueMap, ok = v.(map[string]interface{}); !ok {
			fmt.Fprintf(os.Stderr, "failed to parse generator config for %s", name)
			os.Exit(1)
		}

		for k, v := range valueMap {
			allMatches := re.FindAll([]byte(k), -1)
			for _, match := range allMatches {
				// where's the capture group?
				arguments[strings.Trim(string(match), "{}")] = true
			}
			if vMap, ok := v.(map[string]interface{}); ok {
				for key, value := range vMap {
					switch key {
					case "code":
						strValue, ok := value.(string)
						if !ok {
							fmt.Fprintf(os.Stderr, "failed to parse code value for %s", key)
							os.Exit(1)
						}
						//fmt.Printf("k: %s, key: %s, %s: %s\n", k, key, key, strValue)
						config.Generator[name].TemplateConfigMap[Filename(k)] = GeneratorConfig{Code: strValue}
					default:
						fmt.Fprintf(os.Stderr, "unknown key in %s: %s", k, key)
						os.Exit(1)

					}
				}
			}
		}

		// add each place holder as an argument stringvar flag to the new command above
		for t := range arguments {
			//fmt.Printf("adding required argument %s to newCmd\n", t)
			s := newCmd.Flags().String(t, "", fmt.Sprintf("specify %s value", t))
			if err := newCmd.MarkFlagRequired(t); err != nil {
				fmt.Fprintf(os.Stderr, "failed to require template value: %s", err)
			}
			config.Generator[name].TemplateValues[t] = s
		}
		//fmt.Println("adding newCmd")
		rootCmd.AddCommand(newCmd)

	}
}
