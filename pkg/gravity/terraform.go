package gravity

import (
	"strings"

	"github.com/runatlantis/atlantis/server/events/terraform"
)

// Init runs terraform init on the given directory
// TODO: this is ugly
func Init(path string) (*string, error) {
	// logger := logging.NewSimpleLogger("fooSource", false, logging.Debug)
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
// TODO: this is ugly
func Plan(path string) (*string, error) {
	// logger := logging.NewSimpleLogger("fooSource", false, logging.Debug)
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
// TODO: this is ugly
func Apply(path string) (*string, error) {
	// logger := logging.NewSimpleLogger("fooSource", false, logging.Debug)
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
