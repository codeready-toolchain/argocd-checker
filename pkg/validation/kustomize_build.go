package validation

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
	kbuild "sigs.k8s.io/kustomize/kustomize/v5/commands/build"
	kfsys "sigs.k8s.io/kustomize/kyaml/filesys"
)

func lookupKustomizationFile(logger Logger, afs afero.Afero, basedir string) (string, bool) {
	for _, k := range []string{"kustomization.yaml", "kustomization.yml", "Kustomization"} {
		p := filepath.Join(basedir, k)
		if _, err := afs.Open(p); err == nil {
			logger.Debug("found Kustomization file", "path", p)
			return p, true
		}
	}
	return "", false
}

// verifies that `kustomize build` completes successfully
func checkBuild(logger Logger, fsys kfsys.FileSystem, path string) error {
	logger.Debug("👀 checking kustomize build", "path", path)
	buffy := new(bytes.Buffer)

	kcmd := kbuild.NewCmdBuild(fsys, &kbuild.Help{}, buffy)
	if err := kcmd.RunE(kcmd, []string{path}); err != nil {
		return err
	}

	exec.Command("kustomize", "build")
	return nil
}
