package main

import (
	"fmt"
	"io/ioutil"

	"github.com/openhello/superlicense/pkg/license"
	"github.com/openhello/superlicense/pkg/log"

	"github.com/meilihao/goutil/file"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	appName = "superlicense"

	inGenerate  string
	outGenerate string

	inParse string

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: "superlicense",
		Run:   nil,
	}

	parseCmd = &cobra.Command{
		Use:   "parse",
		Short: "parse license",
		Run:   parese,
	}

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate license",
		Run:   generate,
	}
)

func main() {
	generateCmd.PersistentFlags().StringVar(&inGenerate, "i", `license.yaml`, "raw license yaml file")
	generateCmd.PersistentFlags().StringVar(&outGenerate, "o", `license.dat`, "license file")

	parseCmd.PersistentFlags().StringVar(&inParse, "i", `license.dat`, "license file")

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(generateCmd)

	rootCmd.Execute()
}

func generate(cmd *cobra.Command, args []string) {
	ld := &license.LicenseDisplay{}
	raw := file.GetFileValue(inGenerate)

	err := yaml.Unmarshal([]byte(raw), ld)
	if err != nil {
		log.Glog.Fatal(err.Error())
	}
	// ld.License.Mcode, _ = license.DefaultMcoder.Generate()
	ld.License.Mcode, _ = (&license.AdvanceMocoder{}).Generate()

	tmp, _ := yaml.Marshal(ld)
	fmt.Printf("generate input:\n%s\n---\n", string(tmp))
	ld, err = license.Generate(ld.License)
	if err != nil {
		log.Glog.Fatal(err.Error())
	}

	data, err := yaml.Marshal(ld)
	if err != nil {
		log.Glog.Fatal(err.Error())
	}

	if err = ioutil.WriteFile(outGenerate, data, 0660); err != nil {
		log.Glog.Fatal(err.Error())
	}
}

func parese(cmd *cobra.Command, args []string) {
	raw := file.GetFileValue(inParse)
	ld := &license.LicenseDisplay{}

	fmt.Printf("parese in:\n%s\n---\n", raw)

	err := yaml.Unmarshal([]byte(raw), ld)
	if err != nil {
		log.Glog.Warn(err.Error())

		ld.Raw = raw // only license data
	}
	fmt.Printf("paresing raw license:\n%s\n---\n", ld.Raw)

	l, err := license.ParseLicenseWithRaw(ld.Raw)
	if err != nil {
		log.Glog.Fatal(err.Error())
	}

	data, err := yaml.Marshal(l)
	if err != nil {
		log.Glog.Fatal(err.Error())
	}

	fmt.Println(string(data))
}
