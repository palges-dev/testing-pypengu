package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"pypengu/internal/collector"
	"pypengu/internal/graph"
	"pypengu/internal/models"
	"time"
)

const VERSION = "1.0.0"

func main() {
	outputFile := flag.String("o", "pypengu-output.json", "Output JSON file")
	verbose := flag.Bool("v", false, "Verbose output")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("PyPengu v%s\n", VERSION)
		os.Exit(0)
	}

	fmt.Println("PyPengu: Linux PrivEsc Data Collector")
	fmt.Printf("Version: %s\n", VERSION)
	fmt.Println("Collecting data...")
	fmt.Println()

	metadata := collectMetadata(*verbose)
	
	if *verbose {
		fmt.Println("[*] Collecting users...")
	}
	users := collector.CollectUsers(*verbose)

	if *verbose {
		fmt.Println("[*] Collecting groups...")
	}
	groups := collector.CollectGroups(*verbose)

	if *verbose {
		fmt.Println("[*] Collecting sudo configuration...")
	}
	sudoEntries := collector.CollectSudo(*verbose)

	if *verbose {
		fmt.Println("[*] Scanning for SUID binaries...")
	}
	suidBinaries := collector.CollectSUID(*verbose)

	if *verbose {
		fmt.Println("[*] Checking Docker access...")
	}
	dockerInfo := collector.CollectDocker(*verbose)

	if *verbose {
		fmt.Println("[*] Scanning services...")
	}
	services := collector.CollectServices(*verbose)

	if *verbose {
		fmt.Println("[*] Scanning cron jobs...")
	}
	cronJobs := collector.CollectCron(*verbose)

	if *verbose {
		fmt.Println("[*] Checking kernel version...")
	}
	kernelInfo := collector.CollectKernel(*verbose)

	if *verbose {
		fmt.Println("[*] Building graph...")
	}
	graphData := graph.BuildGraph(
		metadata,
		users,
		groups,
		sudoEntries,
		suidBinaries,
		dockerInfo,
		services,
		cronJobs,
		kernelInfo,
	)

	jsonData, err := json.MarshalIndent(graphData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n[+] Data collected successfully\n")
	fmt.Printf("[+] Output: %s\n", *outputFile)
	fmt.Printf("[+] Nodes: %d | Edges: %d | Paths to root: %d\n",
		graphData.Stats.TotalNodes,
		graphData.Stats.TotalEdges,
		graphData.Stats.PathsToRoot,
	)
}

func collectMetadata(verbose bool) models.Metadata {
	hostname, _ := os.Hostname()
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("LOGNAME")
	}
	uid := os.Getuid()

	return models.Metadata{
		Hostname:    hostname,
		OS:          collector.GetOSInfo(),
		Kernel:      collector.GetKernelVersion(),
		Arch:        collector.GetArch(),
		CollectedAt: time.Now().UTC().Format(time.RFC3339),
		Collector:   "pypengu",
		Version:     VERSION,
		CollectedAs: currentUser,
		UID:         fmt.Sprintf("%d", uid),
	}
}
