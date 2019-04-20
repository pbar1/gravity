package gravity

import (
	"os"
	"path/filepath"
	"strings"

	tfconfig "github.com/hashicorp/terraform/config"

	"github.com/runatlantis/atlantis/server/events/terraform"
)

// TODO: these are all fairly ugly and use internal loggers

// Init runs terraform init on the given directory
func Init(path string) (*string, error) {
	tf, err := terraform.NewClient(nil, "", "", "", "", nil)
	if err != nil {
		return nil, err
	}

	// TODO: lock should not be false
	args := []string{"init", "-no-color", "-lock=false"}
	_, outCh := tf.RunCommandAsync(nil, path, args, nil, "")

	var outStrB strings.Builder
	for l := range outCh {
		if l.Err != nil {
			return nil, l.Err
		}
		outStrB.WriteString(l.Line)
	}
	outStr := outStrB.String()

	return &outStr, nil
}

// Plan runs a Terraform plan on the given directory, and returns the output
func Plan(path string) (*string, error) {
	tf, err := terraform.NewClient(nil, "", "", "", "", nil)
	if err != nil {
		return nil, err
	}

	args := []string{"plan", "-no-color", "-lock=false"}
	_, outCh := tf.RunCommandAsync(nil, path, args, nil, "")

	var outStrB strings.Builder
	for l := range outCh {
		if l.Err != nil {
			return nil, l.Err
		}
		outStrB.WriteString(l.Line)
	}
	outStr := outStrB.String()

	return &outStr, nil
}

// Apply runs a Terraform apply on the given directory, and returns the output
func Apply(path string) (*string, error) {
	tf, err := terraform.NewClient(nil, "", "", "", "", nil)
	if err != nil {
		return nil, err
	}

	args := []string{"apply", "-no-color", "-auto-approve", "-lock=false"}
	_, outCh := tf.RunCommandAsync(nil, path, args, nil, "")

	var outStrB strings.Builder
	for l := range outCh {
		if l.Err != nil {
			return nil, l.Err
		}
		outStrB.WriteString(l.Line)
	}
	outStr := outStrB.String()

	return &outStr, nil
}

// FindStatefulDirs takes directory root and returns the paths within
// that contain a Terraform "backend" definition
func FindStatefulDirs(dir string) ([]string, error) {
	var tfStatefulDirs []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		f, err := os.Stat(path)
		if err != nil {
			return err
		}

		if f.IsDir() {
			tfConfig, err := tfconfig.LoadDir(path)
			if err != nil {
				if strings.HasPrefix(err.Error(), "No Terraform configuration files found in directory") {
					return nil
				}
				return err
			}

			if tfConfig.Terraform.Backend != nil {
				tfStatefulDirs = append(tfStatefulDirs, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tfStatefulDirs, nil
}
