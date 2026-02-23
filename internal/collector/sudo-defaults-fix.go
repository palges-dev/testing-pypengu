// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"bufio"
	"os"
	"os/exec"
	"pypengu/internal/models"
	"strings"
)

var sudoGTFOBins = map[string]string{
	"yelp":    "https://gtfobins.github.io/gtfobins/yelp/",
	"pkexec":  "https://gtfobins.github.io/gtfobins/pkexec/",
	"time":    "https://gtfobins.github.io/gtfobins/time/",
	"mtr":     "https://gtfobins.github.io/gtfobins/mtr/",
	"rlogin":  "https://gtfobins.github.io/gtfobins/rlogin/",
	"finger":  "https://gtfobins.github.io/gtfobins/finger/",
	"cancel":  "https://gtfobins.github.io/gtfobins/cancel/",
	"whois":   "https://gtfobins.github.io/gtfobins/whois/",
	"vim":     "https://gtfobins.github.io/gtfobins/vim/",
	"nano":    "https://gtfobins.github.io/gtfobins/nano/",
	"find":    "https://gtfobins.github.io/gtfobins/find/",
	"bash":    "https://gtfobins.github.io/gtfobins/bash/",
	"sh":      "https://gtfobins.github.io/gtfobins/sh/",
	"python":  "https://gtfobins.github.io/gtfobins/python/",
	"python3": "https://gtfobins.github.io/gtfobins/python/",
	"perl":    "https://gtfobins.github.io/gtfobins/perl/",
	"ruby":    "https://gtfobins.github.io/gtfobins/ruby/",
	"awk":     "https://gtfobins.github.io/gtfobins/awk/",
	"git":     "https://gtfobins.github.io/gtfobins/git/",
	"less":    "https://gtfobins.github.io/gtfobins/less/",
	"more":    "https://gtfobins.github.io/gtfobins/more/",
	"dmf":     "https://gtfobins.github.io/",
}

func CollectSudo(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	entries = append(entries, runSudoLCurrent(verbose)...)
	entries = append(entries, parseSudoersFile(verbose)...)
	entries = append(entries, parseSudoersD(verbose)...)
	
	if verbose {
		println("    Found", len(entries), "sudo entries")
	}
	
	return entries
}

func runSudoLCurrent(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	cmd := exec.Command("sudo", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if verbose {
			println("    [-] Could not run sudo -l")
		}
		return entries
	}
	
	lines := strings.Split(string(output), "\n")
	inCommandSection := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "Matching Defaults") || strings.HasPrefix(line, "Defaults entries") {
			continue
		}
		
		if strings.Contains(line, "may run the following commands") {
			inCommandSection = true
			continue
		}
		
		if !inCommandSection || line == "" {
			continue
		}
		
		if strings.HasPrefix(line, "User ") || strings.HasPrefix(line, "Runas") || strings.HasPrefix(line, "Host") {
			continue
		}
		
		if !strings.Contains(line, "(root)") && !strings.Contains(line, "(ALL)") {
			continue
		}
		
		noPasswd := strings.Contains(line, "NOPASSWD:")
		
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		
		var command string
		var runAs string = "root"
		
		for i, part := range parts {
			if part == "NOPASSWD:" && i+1 < len(parts) {
				command = parts[i+1]
				break
			}
			if strings.HasPrefix(part, "/") {
				command = part
				break
			}
		}
		
		if command == "" {
			continue
		}
		
		basename := command
		if idx := strings.LastIndex(command, "/"); idx != -1 {
			basename = command[idx+1:]
		}
		
		gtfoUrl, isGTFO := sudoGTFOBins[basename]
		
		entry := models.SudoEntry{
			User:     "current",
			Command:  command,
			NoPasswd: noPasswd,
			RunAs:    runAs,
			Entry:    line,
		}
		
		if isGTFO {
			entry.GTFOBin = true
			entry.GTFOUrl = gtfoUrl
		}
		
		entries = append(entries, entry)
		
		if verbose {
			println("      [+]", command, "NOPASSWD:", noPasswd, "GTFOBin:", isGTFO)
		}
	}
	
	return entries
}

func parseSudoersFile(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	file, err := os.Open("/etc/sudoers")
	if err != nil {
		return entries
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "Defaults") {
			continue
		}
		
		entry := parseSudoLine(line)
		if entry.User != "" && entry.User != "Defaults" {
			entries = append(entries, entry)
		}
	}
	
	return entries
}

func parseSudoersD(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	files, err := os.ReadDir("/etc/sudoers.d")
	if err != nil {
		return entries
	}
	
	for _, f := range files {
		if f.IsDir() || strings.HasPrefix(f.Name(), ".") {
			continue
		}
		
		path := "/etc/sudoers.d/" + f.Name()
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "Defaults") {
				continue
			}
			
			entry := parseSudoLine(line)
			if entry.User != "" && entry.User != "Defaults" {
				entries = append(entries, entry)
			}
		}
		file.Close()
	}
	
	return entries
}

func parseSudoLine(line string) models.SudoEntry {
	entry := models.SudoEntry{Entry: line}
	
	if strings.HasPrefix(line, "Defaults") {
		return entry
	}
	
	noPasswd := strings.Contains(line, "NOPASSWD")
	entry.NoPasswd = noPasswd
	
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return entry
	}
	
	entry.User = parts[0]
	
	if strings.Contains(line, "ALL=(ALL)") || strings.Contains(line, "ALL=(ALL:ALL)") {
		entry.RunAs = "ALL"
		if len(parts) > 2 {
			entry.Command = strings.Join(parts[2:], " ")
		} else {
			entry.Command = "ALL"
		}
	}
	
	return entry
}
