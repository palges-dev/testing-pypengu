// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"pypengu/internal/models"
	"strings"
)

var kernelCVEs = map[string][]models.CVEInfo{
	"5.15.0": {
		{
			ID:          "CVE-2023-0386",
			Description: "OverlayFS privilege escalation",
			Reference:   "https://nvd.nist.gov/vuln/detail/CVE-2023-0386",
		},
	},
	"5.4.0": {
		{
			ID:          "CVE-2022-0847",
			Description: "Dirty Pipe - overwrite data in arbitrary read-only files",
			Reference:   "https://nvd.nist.gov/vuln/detail/CVE-2022-0847",
		},
	},
	"4.15.0": {
		{
			ID:          "CVE-2021-3493",
			Description: "OverlayFS privilege escalation",
			Reference:   "https://nvd.nist.gov/vuln/detail/CVE-2021-3493",
		},
	},
}

func CollectKernel(verbose bool) models.KernelInfo {
	version := GetKernelVersion()
	
	info := models.KernelInfo{
		Version: version,
		CVEs:    []models.CVEInfo{},
	}
	
	for kernelVer, cves := range kernelCVEs {
		if strings.Contains(version, kernelVer) {
			info.CVEs = append(info.CVEs, cves...)
		}
	}
	
	if verbose {
		println("    Kernel:", version)
		if len(info.CVEs) > 0 {
			println("    Found", len(info.CVEs), "potential CVEs")
		}
	}
	
	return info
}
