package output

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Table represents a table output
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

// Render renders the table to stdout
func (t *Table) Render() {
	table := tablewriter.NewWriter(os.Stdout)

	// Enable borders and separators for visible table lines
	table.SetBorder(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(false)

	// Set header first (before any other configuration that might affect it)
	// For tablewriter, we need to set headers without color first for width calculation
	// Then we can apply colors, but tablewriter doesn't handle ANSI codes well in headers
	// So we'll use plain headers for now to avoid width calculation issues
	table.SetHeader(t.headers)

	// Ensure header is displayed
	if len(t.headers) > 0 {
		table.SetHeaderLine(true)
	}

	// Find STATUS and READY column indices
	statusColIndex := -1
	readyColIndex := -1
	for i, header := range t.headers {
		headerUpper := strings.ToUpper(header)
		if headerUpper == "STATUS" {
			statusColIndex = i
		}
		if headerUpper == "READY" {
			readyColIndex = i
		}
	}

	// Add rows with colorized status
	for _, row := range t.rows {
		coloredRow := make([]string, len(row))
		for i, cell := range row {
			// Colorize STATUS column if it exists
			switch {
			case i == statusColIndex && statusColIndex >= 0:
				status := strings.TrimSpace(cell)
				if isStatusField(status) {
					coloredRow[i] = ColorizeStatus(status)
				} else {
					coloredRow[i] = cell
				}
			case i == readyColIndex && readyColIndex >= 0:
				// Colorize READY column (e.g., "1/1", "2/3")
				ready := strings.TrimSpace(cell)
				if strings.Contains(ready, "/") {
					parts := strings.Split(ready, "/")
					if len(parts) == 2 {
						if parts[0] == parts[1] {
							// All ready - green
							coloredRow[i] = StatusReady.Sprint(ready)
						} else {
							// Partially ready - yellow
							coloredRow[i] = StatusPending.Sprint(ready)
						}
					} else {
						coloredRow[i] = ready
					}
				} else {
					coloredRow[i] = ready
				}
			default:
				coloredRow[i] = cell
			}
		}
		table.Append(coloredRow)
	}

	table.Render()
}

// isStatusField checks if a string looks like a status field
func isStatusField(s string) bool {
	statuses := []string{
		"Running", "Pending", "Failed", "Succeeded", "Error",
		"Ready", "NotReady", "Active", "Inactive",
		"ContainerCreating", "PodInitializing", "CrashLoopBackOff",
		"ImagePullBackOff", "ErrImagePull", "Completed",
	}
	for _, status := range statuses {
		if s == status {
			return true
		}
	}
	return false
}
