package main

import (
	"os"
	"path/filepath"
	"syscall"
)

func setupOverlay(id string, lower string) string {
	base := filepath.Join(stateDir, id)
	upper := filepath.Join(base, "upper")
	work := filepath.Join(base, "work")
	merged := filepath.Join(base, "merged")

	os.MkdirAll(upper, 0755)
	os.MkdirAll(work, 0755)
	os.MkdirAll(merged, 0755)

	opts := "lowerdir=" + lower + ",upperdir=" + upper + ",workdir=" + work
	must(syscall.Mount("overlay", merged, "overlay", 0, opts))

	return merged
}

func cleanupOverlay(id string) {
	merged := filepath.Join(stateDir, id, "merged")
	syscall.Unmount(merged, 0)
}
