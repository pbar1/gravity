package gravity

import (
	"strings"

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
