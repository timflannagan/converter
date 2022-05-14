package unstructured

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func FromReader(reader io.Reader) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 1)

	unst := &unstructured.Unstructured{}
	err := decoder.Decode(unst)
	if err != nil {
		return nil, err
	}

	return unst, nil
}

func FromString(str string) (*unstructured.Unstructured, error) {
	return FromReader(strings.NewReader(str))
}

func FromBytes(b []byte) (*unstructured.Unstructured, error) {
	return FromReader(bytes.NewReader(b))
}

func FromFile(filepath string) (*unstructured.Unstructured, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return FromReader(file)
}

func FromDir(dirpath string) ([]*unstructured.Unstructured, error) {
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	unsts := make([]*unstructured.Unstructured, 0, 0)
	for _, file := range files {
		unst, err := FromFile(path.Join(dirpath, file.Name()))
		if err != nil {
			return nil, err
		}

		unsts = append(unsts, unst)
	}

	return unsts, nil
}
