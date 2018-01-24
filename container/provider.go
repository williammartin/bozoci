package container

import (
	"os/exec"
	"path/filepath"
)

type Provider struct {
	ContainersDir string
}

func (p *Provider) Provide(containerID, rootfs, command string) error {
	cmd := exec.Command("runc", "run", "-d", containerID)
	cmd.Dir = filepath.Join(p.ContainersDir, containerID)
	return cmd.Run()
}
