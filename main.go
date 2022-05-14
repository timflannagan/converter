package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/timflannagan/converter/internal/convert"
	"github.com/timflannagan/converter/internal/unstructured"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

const (
	installNamespace = "rukpak-system"
	// TODO: flag to control output directory
	outputDir = "./plain"
)

var (
	targetNamespaces = []string{}
)

func main() {
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
			objects, err := unstructured.FromDir(inputDir)
			if err != nil {
				return err
			}
			if len(objects) == 0 {
				return fmt.Errorf("failed to read bundle manifests")
			}

			var (
				reg convert.RegistryV1
			)
			for _, obj := range objects {
				obj := obj

				switch obj.GetObjectKind().GroupVersionKind().Kind {
				case "ClusterServiceVersion":
					csv := v1alpha1.ClusterServiceVersion{}
					if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &csv); err != nil {
						return err
					}
					csv.SetNamespace(installNamespace)
					reg.CSV = csv
				case "CustomResourceDefinition":
					crd := apiextensionsv1.CustomResourceDefinition{}
					if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &crd); err != nil {
						return err
					}
					reg.CRDs = append(reg.CRDs, crd)
				default:
					reg.Others = append(reg.Others)
				}
			}
			objs, err := convert.Convert(reg, installNamespace, targetNamespaces)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create the %s output directory: %w", outputDir, err)
			}
			for _, obj := range objs.Objects {
				data, err := yaml.Marshal(obj)
				if err != nil {
					return err
				}
				filename := filepath.Join(outputDir, fmt.Sprintf("%s-%s.yaml", obj.GetName(), strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)))
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
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
