package collector

import (
	"bufio"
	"os"
	"os/exec"
	"pypengu/internal/models"
	"strings"
)

func CollectSudo(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	entries = append(entries, parseSudoersFile(verbose)...)
	entries = append(entries, parseSudoersD(verbose)...)
	entries = append(entries, runSudoL(verbose)...)
	
	if verbose {
		println("    Found", len(entries), "sudo entries")
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
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		entry := parseSudoLine(line)
		if entry.User != "" {
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
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			
			entry := parseSudoLine(line)
			if entry.User != "" {
				entries = append(entries, entry)
			}
		}
		file.Close()
	}
	
	return entries
}

func parseSudoLine(line string) models.SudoEntry {
	entry := models.SudoEntry{Entry: line}
	
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

func runSudoL(verbose bool) []models.SudoEntry {
	entries := []models.SudoEntry{}
	
	cmd := exec.Command("sudo", "-l")
	output, err := cmd.Output()
	if err != nil {
		return entries
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "may run") {
			continue
		}
		
		if strings.Contains(line, "NOPASSWD") {
			entry := models.SudoEntry{
				User:     "current",
				NoPasswd: true,
				Entry:    line,
			}
			entries = append(entries, entry)
		}
	}
	
	return entries
}
