package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

func PrintSettings(settings Settings, format Format, w io.Writer) (err error) {
	var encoder interface {
		Encode(v interface{}) error
	}
	switch format {
	case "yaml":
		encoder = yaml.NewEncoder(w)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
	return encoder.Encode(settings)
}

func PrintSection(section interface{ IsEmpty() bool }, format Format, w io.Writer) (err error) {
	if section.IsEmpty() {
		return nil
	}
	var encoder interface {
		Encode(v interface{}) error
	}
	switch format {
	case "yaml":
		encoder = yaml.NewEncoder(w)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
	return encoder.Encode(section)
}
