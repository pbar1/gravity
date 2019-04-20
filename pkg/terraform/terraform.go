package terraform

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	tfconfig "github.com/hashicorp/terraform/config"
	tf "github.com/runatlantis/atlantis/server/events/terraform"
)

func commandHelper(path string, args []string) (*string, error) {
	tfClient, err := tf.NewClient(nil, "", "", "", "", nil)
	if err != nil {
		return nil, err
	}

	_, outCh := tfClient.RunCommandAsync(nil, path, args, nil, "")

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

// Init runs terraform init on the given directory
func Init(path string) (*string, error) {
	args := []string{"init", "-no-color", "-input=false", "-upgrade=true"}
	outStr, err := commandHelper(path, args)
	if err != nil {
		return nil, err
	}
	return outStr, nil
}

// WorkspaceList returns a string slice of the workspaces found in
// provided Terrafrom project directory
func WorkspaceList(path string) ([]string, error) {
	args := []string{"workspace", "list"}
	outStr, err := commandHelper(path, args)
	if err != nil {
		return nil, err
	}

	wsDirty := strings.Split(strings.Replace(*outStr, "*", "", 1), "\n")
	var workspaces []string
	for _, w := range wsDirty {
		workspaces = append(workspaces, strings.TrimSpace(w))
	}
	return workspaces, nil
}

// WorkspaceSelect selects the provided Terraform workspace name
func WorkspaceSelect(path, workspace string) (*string, error) {
	args := []string{"workspace", "select", workspace}
	outStr, err := commandHelper(path, args)
	if err != nil {
		return nil, err
	}
	return outStr, nil
}

// Plan runs a Terraform plan on the given directory, and returns the output
func Plan(path string) (*string, error) {
	args := []string{"plan", "-no-color", "-input=false"}
	outStr, err := commandHelper(path, args)
	if err != nil {
		return nil, err
	}
	return outStr, nil
}

// Apply runs a Terraform apply on the given directory, and returns the output
func Apply(path string) (*string, error) {
	args := []string{"apply", "-no-color", "-auto-approve", "-lock=false"}
	outStr, err := commandHelper(path, args)
	if err != nil {
		return nil, err
	}
	return outStr, nil
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

// FindWorkspaceVarFile attempts to find a Terraform tfvars file with a
// filename matching the provided workspace name, searching recursively
// from the provided directory. Will return an empty string if no
// matching file is found.
func FindWorkspaceVarFile(dir, workspace string) (*string, error) {
	var workspaceVarFile string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		f, err := os.Stat(path)
		if err != nil {
			return err
		}

		if !f.IsDir() && filepath.Base(f.Name()) == workspace+".tfvars" {
			workspaceVarFile = f.Name()
			return io.EOF // break out of walk on first match
		}

		return nil
	})

	if err == io.EOF {
		return &workspaceVarFile, nil
	}

	return nil, err
}
