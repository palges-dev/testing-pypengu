// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"os"
	"path/filepath"
	"pypengu/internal/models"
	"strings"
	"syscall"
)

func CollectServices(verbose bool) []models.ServiceInfo {
	services := []models.ServiceInfo{}
	
	systemdPaths := []string{
		"/etc/systemd/system",
		"/lib/systemd/system",
		"/usr/lib/systemd/system",
	}
	
	currentUID := os.Getuid()
	
	for _, systemdPath := range systemdPaths {
		filepath.Walk(systemdPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			
			if info.IsDir() {
				return nil
			}
			
			if !strings.HasSuffix(path, ".service") {
				return nil
			}
			
			stat, ok := info.Sys().(*syscall.Stat_t)
			if !ok {
				return nil
			}
			
			writable := false
			writableBy := ""
			
			if stat.Uid == uint32(currentUID) && info.Mode().Perm()&0200 != 0 {
				writable = true
				writableBy = "owner"
			} else if info.Mode().Perm()&0020 != 0 {
				writable = true
				writableBy = "group"
			} else if info.Mode().Perm()&0002 != 0 {
				writable = true
				writableBy = "world"
			}
			
			if writable {
				service := models.ServiceInfo{
					Name:       filepath.Base(path),
					Path:       path,
					RunAs:      "root",
					State:      "unknown",
					Writable:   true,
					WritableBy: writableBy,
				}
				services = append(services, service)
			}
			
			return nil
		})
	}
	
	if verbose {
		println("    Found", len(services), "writable services")
	}
	
	return services
}
