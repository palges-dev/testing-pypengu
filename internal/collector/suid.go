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
	"syscall"
)

var gtfoBins = map[string]string{
	"vim":     "https://gtfobins.github.io/gtfobins/vim/",
	"nano":    "https://gtfobins.github.io/gtfobins/nano/",
	"find":    "https://gtfobins.github.io/gtfobins/find/",
	"nmap":    "https://gtfobins.github.io/gtfobins/nmap/",
	"python":  "https://gtfobins.github.io/gtfobins/python/",
	"perl":    "https://gtfobins.github.io/gtfobins/perl/",
	"ruby":    "https://gtfobins.github.io/gtfobins/ruby/",
	"bash":    "https://gtfobins.github.io/gtfobins/bash/",
	"sh":      "https://gtfobins.github.io/gtfobins/sh/",
	"less":    "https://gtfobins.github.io/gtfobins/less/",
	"more":    "https://gtfobins.github.io/gtfobins/more/",
	"awk":     "https://gtfobins.github.io/gtfobins/awk/",
	"gcc":     "https://gtfobins.github.io/gtfobins/gcc/",
	"tar":     "https://gtfobins.github.io/gtfobins/tar/",
	"zip":     "https://gtfobins.github.io/gtfobins/zip/",
	"unzip":   "https://gtfobins.github.io/gtfobins/unzip/",
	"git":     "https://gtfobins.github.io/gtfobins/git/",
	"docker":  "https://gtfobins.github.io/gtfobins/docker/",
}

func CollectSUID(verbose bool) []models.SuidBinary {
	binaries := []models.SuidBinary{}
	
	searchPaths := []string{
		"/bin", "/sbin", "/usr/bin", "/usr/sbin",
		"/usr/local/bin", "/usr/local/sbin",
		"/opt",
	}
	
	for _, searchPath := range searchPaths {
		filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			
			if info.IsDir() {
				return nil
			}
			
			if info.Mode()&os.ModeSetuid != 0 || info.Mode()&os.ModeSetgid != 0 {
				stat, ok := info.Sys().(*syscall.Stat_t)
				if !ok {
					return nil
				}
				
				owner := "root"
				if stat.Uid != 0 {
					owner = "other"
				}
				
				basename := filepath.Base(path)
				gtfoUrl, isGTFO := gtfoBins[basename]
				
				binary := models.SuidBinary{
					Path:    path,
					Owner:   owner,
					GTFOBin: isGTFO,
					GTFOUrl: gtfoUrl,
				}
				
				binaries = append(binaries, binary)
			}
			
			return nil
		})
	}
	
	if verbose {
		println("    Found", len(binaries), "SUID binaries")
	}
	
	return binaries
}
