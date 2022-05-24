package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/timflannagan/converter/internal/convert"
	"github.com/timflannagan/converter/internal/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	// TODO: avoid hardcoding the install namespace
	// Default to operators or openshift-operators unless the CSV's suggested-namespace
	// annotation is present?
	installNamespace = "rukpak-system"
	defaultOutputDir = "plain"
)

var (
	// TODO: avoid hardcoding the target namespaces
	targetNamespaces = []string{}
)

func main() {
	var (
		outputDir string
	)
	cmd := &cobra.Command{
		Use:  "convert",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputDir := args[0]
			info, err := os.Stat(inputDir)
			if err != nil {
				return fmt.Errorf("failed to stat the %s input directory: %s", inputDir, err)
			}
			if !info.IsDir() {
				return fmt.Errorf("failed to stat the %s input directory: input is not a directory", inputDir)
			}

			// TODO: Refactor the convert library to work with both static and runtime use cases.
			// TODO: Generate a plain manifest dockerfile?
			// TODO: Apply the generated plain bundle on a kind cluster
			// TODO: We need some way of projecting the installation namespace as the current
			// plain provisioner's BI controller assumes that resource exists in the Bundle
			// source type.
			// TODO: How to correctly generate the targetNamespaces without the usage of the OG
			// resource?
			objects, err := unstructured.FromDir(inputDir)
			if err != nil {
				return err
			}
			plain, err := convert.Convert(objects, installNamespace, targetNamespaces)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				return err
			}
			if err := ensurePlainManifests(plain.Objects, outputDir); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir, "Configures the directory that will contain the outputted set of decomposed bundle manifests")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ensurePlainManifests(objs []client.Object, outputDir string) error {
	// TODO: Return an aggregate error
	// TODO: Wrap this with an io.Writer?
	for _, obj := range objs {
		data, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		filename := generateManifestFileName(outputDir, obj)
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create the %s file: %w", filename, err)
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func generateManifestFileName(dir string, obj client.Object) string {
	return filepath.Join(dir, fmt.Sprintf("%s-%s.yaml", obj.GetName(), strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)))
}
