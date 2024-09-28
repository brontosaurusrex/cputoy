package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// Global variable for terminal width
var width int = 30
var widthold int = 30

// Struct to hold CPU statistics
type CPUStat struct {
	User, Nice, System, Idle, IOWait, IRQ, SoftIRQ, Steal int64
}

// Function to parse a line from /proc/stat
func parseCPUStat(statLine string) (CPUStat, error) {
	fields := strings.Fields(statLine)
	if len(fields) < 8 {
		return CPUStat{}, fmt.Errorf("malformed /proc/stat data")
	}

	user, _ := strconv.ParseInt(fields[1], 10, 64)
	nice, _ := strconv.ParseInt(fields[2], 10, 64)
	system, _ := strconv.ParseInt(fields[3], 10, 64)
	idle, _ := strconv.ParseInt(fields[4], 10, 64)
	iowait, _ := strconv.ParseInt(fields[5], 10, 64)
	irq, _ := strconv.ParseInt(fields[6], 10, 64)
	softirq, _ := strconv.ParseInt(fields[7], 10, 64)
	steal := int64(0)
	if len(fields) >= 9 {
		steal, _ = strconv.ParseInt(fields[8], 10, 64)
	}

	return CPUStat{user, nice, system, idle, iowait, irq, softirq, steal}, nil
}

// Function to get CPU statistics from /proc/stat
func getCPUStats() ([]CPUStat, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var stats []CPUStat

	// Parse each line starting with "cpu" (for all CPUs)
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu") && !strings.HasPrefix(line, "cpu ") {
			stat, err := parseCPUStat(line)
			if err != nil {
				return nil, err
			}
			stats = append(stats, stat)
		}
	}

	return stats, nil
}

// Function to calculate the CPU usage percentage between two CPUStat snapshots
func calculateUsage(stat1, stat2 CPUStat) float64 {
	idle1 := stat1.Idle + stat1.IOWait
	idle2 := stat2.Idle + stat2.IOWait

	total1 := stat1.User + stat1.Nice + stat1.System + stat1.Idle +
		stat1.IOWait + stat1.IRQ + stat1.SoftIRQ + stat1.Steal
	total2 := stat2.User + stat2.Nice + stat2.System + stat2.Idle +
		stat2.IOWait + stat2.IRQ + stat2.SoftIRQ + stat2.Steal

	totalDiff := total2 - total1
	idleDiff := idle2 - idle1

	if totalDiff == 0 {
		return 0.0
	}

	return float64(totalDiff-idleDiff) / float64(totalDiff) * 100.0
}

// Function to generate a bar with full blocks for CPU usage visualization
func generateBar(usage float64) string {
	totalBlocks := width // Total width of the terminal
	//fmt.Printf("%d", totalBlocks) // debug
	// Calculate how many full blocks should be filled
	fullBlocks := int(usage / 100.0 * float64(totalBlocks))

	// Create the full bar using full blocks and empty space
	filledBar := strings.Repeat("|", fullBlocks)

	// Ensure that the remaining empty space count is non-negative
	remainingBlocks := totalBlocks - fullBlocks
	if remainingBlocks < 0 {
		remainingBlocks = 0
	}

	emptyBar := strings.Repeat(" ", remainingBlocks)

	return filledBar + emptyBar
}

// Function to clear and refresh the terminal output
func clearTerminal() {
	fmt.Print("\033[H\033[2J") // Move to top-left and clear the screen
}

// Function to move cursor to the top
func cursorTerminal() {
	fmt.Print("\033[H") // Move to top-left
}

func main() {

	// Set up channel to listen for Ctrl+C (SIGINT)
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Clear terminal
	clearTerminal()

	go func() {
		// Start the main loop
		for {
			// Get initial CPU stats
			initialStats, err := getCPUStats()
			if err != nil {
				fmt.Println("Error reading CPU stats:", err)
				time.Sleep(1200 * time.Millisecond) // Wait a bit before retrying
				continue
			}

			// term width
			widthold = width
			ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
			if err != nil {
				fmt.Println("Error getting terminal size:", err)
				return
			}
			width = int(ws.Col) // replacing the global variable
			//fmt.Printf("Terminal width: %d\n", width)
			//if width is different than previous width, clear the term
			if width != widthold {
				clearTerminal()
			}

			// Wait for N milliseconds
			time.Sleep(500 * time.Millisecond)

			// Get final CPU stats
			finalStats, err := getCPUStats()
			if err != nil {
				fmt.Println("Error reading CPU stats:", err)
				time.Sleep(1 * time.Second) // Wait a bit before retrying
				continue
			}

			// Calculate and print CPU usage for each core, renumbering CPUs from 1 to LAST
			for i := range initialStats {
				usage := calculateUsage(initialStats[i], finalStats[i])
				bar := generateBar(usage) // No need to pass terminalWidth
				fmt.Printf("%s\n", bar)   // Display CPU number
			}

			// Move cursor to the top-left
			cursorTerminal()
		}
	}()

	// Block the main thread until Ctrl+C is received
	<-sigChannel
	fmt.Println("Exiting...")
}
