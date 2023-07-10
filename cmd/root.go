package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/svrana/geniveev"
)

var cfgFile string = ".geniveev.toml"
var config *geniveev.Config = geniveev.NewConfig()

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

func TemplateFuncTitle(input string) string {
	return cases.Title(language.Make("en")).String(input)
}

func giniTemplateParse(name string, templateStr string, templateValues geniveev.TemplateValues) (string, error) {
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"Title": TemplateFuncTitle,
	}).Parse(templateStr)
	if err != nil {
		return "", err
	}
	out := &bytes.Buffer{}
	if err := tmpl.Execute(out, templateValues); err != nil {
		return "", err
	}
	return out.String(), nil
}

func constructFilename(templatedFilename geniveev.Filename, templateValues geniveev.TemplateValues) (string, error) {
	var filename = string(templatedFilename)
	return giniTemplateParse(filename, filename, templateValues)
}

func createPath(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to check %s for existence: %s", filename, err)

		}
	} else {
		return fmt.Errorf("%s already exists, will not overwrite", filename)
	}

	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create directory %s: %s", path.Dir(filename), err)
		os.Exit(1)
	}
	return nil
}

func generate(generatorConfig geniveev.GeneratorConfig, filename string, templateConfig *geniveev.TemplateConfig) error {
	code, err := giniTemplateParse(filename, generatorConfig.Code, templateConfig.TemplateValues)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
		return err
	}
	return nil
}

func start(templateConfig *geniveev.TemplateConfig) error {
	for templatedFilename, config := range templateConfig.TemplateConfigMap {
		filename, err := constructFilename(templatedFilename, templateConfig.TemplateValues)
		if err != nil {
			return err
		}
		if err := createPath(filename); err != nil {
			return err
		}
		if err := generate(config, filename, templateConfig); err != nil {
			return err
		}
	}
	return nil
}

// Initialize reads in config file and ENV variables if set.
func Initialize() {
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

	re := regexp.MustCompile(`{(.\w+)}`)

	for name, v := range cfg {
		config.Generator[name] = geniveev.NewTemplateConfigEmpty()

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
				arguments[strings.Trim(string(match), ".{}")] = true
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

						// FIXME: go through code and create parameters out of any variables used there
						config.Generator[name].TemplateConfigMap[geniveev.Filename(k)] = geniveev.GeneratorConfig{Code: strValue}
					default:
						fmt.Fprintf(os.Stderr, "unknown key in %s: %s", k, key)
						os.Exit(1)

					}
				}
			}
		}

		newCmd := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Generate %s", name),
			Run: func(_ *cobra.Command, _ []string) {
				if err := start(config.Generator[name]); err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					os.Exit(1)
				}
			},
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

		rootCmd.AddCommand(newCmd)
	}
}
