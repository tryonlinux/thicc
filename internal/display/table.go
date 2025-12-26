package display

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/tryonlinux/thicc/internal/models"
)

const asciiArt = `
 ████████╗██╗  ██╗██╗ ██████╗ ██████╗
 ╚══██╔══╝██║  ██║██║██╔════╝██╔════╝
    ██║   ███████║██║██║     ██║
    ██║   ██╔══██║██║██║     ██║
    ██║   ██║  ██║██║╚██████╗╚██████╗
    ╚═╝   ╚═╝  ╚═╝╚═╝ ╚═════╝ ╚═════╝
    Weight Tracker
`

// RenderWeightsTable creates a formatted table of weights with a line graph
func RenderWeightsTable(weights []models.Weight, settings *models.Settings, limit int) string {
	if len(weights) == 0 {
		return TitleStyle.Render(asciiArt) + "\n\nNo weights tracked. Add one with: thicc add <weight> [date]"
	}

	// Start with ASCII art
	var output strings.Builder
	output.WriteString(TitleStyle.Render(asciiArt))
	output.WriteString("\n")

	// Calculate stats
	var totalWeight, minWeight, maxWeight float64
	minWeight = math.MaxFloat64
	maxWeight = -math.MaxFloat64

	for _, w := range weights {
		totalWeight += w.Weight
		if w.Weight < minWeight {
			minWeight = w.Weight
		}
		if w.Weight > maxWeight {
			maxWeight = w.Weight
		}
	}

	avgWeight := totalWeight / float64(len(weights))
	latestWeight := weights[0].Weight
	latestBMI := weights[0].BMI
	startWeight := weights[len(weights)-1].Weight

	// Calculate change from start
	delta := latestWeight - startWeight
	var deltaStr string
	if delta < 0 {
		// Lost weight
		deltaStr = fmt.Sprintf("Lost %s", FormatWeight(math.Abs(delta), settings.WeightUnit))
	} else if delta > 0 {
		// Gained weight
		deltaStr = fmt.Sprintf("Gained %s", FormatWeight(delta, settings.WeightUnit))
	} else {
		deltaStr = "No change"
	}

	// Build stats header (goes with ASCII art header)
	var header strings.Builder
	header.WriteString(HeaderStyle.Render(fmt.Sprintf("Latest: %s | BMI: %s | Avg: %s | %s",
		FormatWeight(latestWeight, settings.WeightUnit),
		FormatBMI(latestBMI),
		FormatWeight(avgWeight, settings.WeightUnit),
		deltaStr)))
	header.WriteString("\n")
	header.WriteString(InfoStyle.Render(fmt.Sprintf("Min: %s | Max: %s | Entries: %d",
		FormatWeight(minWeight, settings.WeightUnit),
		FormatWeight(maxWeight, settings.WeightUnit),
		len(weights))))
	header.WriteString("\n\n\n")

	// Calculate goal difference
	goalDiff := latestWeight - settings.GoalWeight
	var goalDiffStr string
	if goalDiff > 0 {
		// Current weight is above goal - need to lose
		goalDiffStr = fmt.Sprintf("%.1f %s to lose", goalDiff, settings.WeightUnit)
	} else if goalDiff < 0 {
		// Current weight is below goal - need to gain
		goalDiffStr = fmt.Sprintf("%.1f %s to gain", math.Abs(goalDiff), settings.WeightUnit)
	} else {
		goalDiffStr = "at goal!"
	}

	// Build goal weight section (goes with table/graph below)
	goalHeader := fmt.Sprintf("Goal Weight: %s (%s)",
		FormatWeight(settings.GoalWeight, settings.WeightUnit),
		goalDiffStr)
	centeredGoalStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		Align(lipgloss.Center).
		Width(GoalHeaderWidth)
	header.WriteString(centeredGoalStyle.Render(goalHeader))
	header.WriteString("\n\n")

	// Truncate to TableMaxRows for table display
	displayWeights := weights
	if len(weights) > TableMaxRows {
		displayWeights = weights[:TableMaxRows]
	}

	// Create table and graph side by side
	weightTable := createWeightTable(displayWeights, settings)
	weightGraph := createLineGraph(weights, settings)

	// Combine table and graph
	combined := lipgloss.JoinHorizontal(lipgloss.Top, weightTable, "  ", weightGraph)

	return output.String() + header.String() + combined
}

// createWeightTable creates the weight table
func createWeightTable(weights []models.Weight, settings *models.Settings) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(TableBorderStyle).
		Headers("ID", "Date", "Weight", "BMI")

	for _, w := range weights {
		t.Row(
			fmt.Sprintf("%d", w.ID),
			FormatDate(w.Date),
			FormatWeight(w.Weight, settings.WeightUnit),
			FormatBMI(w.BMI),
		)
	}

	return t.Render()
}

// weightRange holds the min and max weight values for graph scaling
type weightRange struct {
	min float64
	max float64
}

// calculateWeightRange determines the min and max weights including goal weight and padding
func calculateWeightRange(weights []models.Weight, goalWeight float64) weightRange {
	minWeight := math.MaxFloat64
	maxWeight := -math.MaxFloat64

	for _, w := range weights {
		if w.Weight < minWeight {
			minWeight = w.Weight
		}
		if w.Weight > maxWeight {
			maxWeight = w.Weight
		}
	}

	// Include goal weight in range calculation
	if goalWeight < minWeight {
		minWeight = goalWeight
	}
	if goalWeight > maxWeight {
		maxWeight = goalWeight
	}

	// Add some padding to the range
	padding := (maxWeight - minWeight) * 0.1
	if padding == 0 {
		padding = 1
	}
	minWeight -= padding
	maxWeight += padding

	return weightRange{min: minWeight, max: maxWeight}
}

// createGraphGrid initializes an empty graph grid
func createGraphGrid(width, height int) [][]rune {
	graph := make([][]rune, height)
	for i := range graph {
		graph[i] = make([]rune, width)
		for j := range graph[i] {
			graph[i][j] = ' '
		}
	}
	return graph
}

// reverseWeights returns a reversed copy of the weights slice (oldest to newest)
func reverseWeights(weights []models.Weight) []models.Weight {
	reversed := make([]models.Weight, len(weights))
	copy(reversed, weights)
	for i := 0; i < len(reversed)/2; i++ {
		reversed[i], reversed[len(reversed)-1-i] = reversed[len(reversed)-1-i], reversed[i]
	}
	return reversed
}

// normalizeToGraphY converts a weight value to a Y coordinate on the graph
func normalizeToGraphY(weight float64, wr weightRange, height int) int {
	normalized := (weight - wr.min) / (wr.max - wr.min)
	y := height - 1 - int(normalized*float64(height-1))

	// Ensure y is within bounds
	if y < 0 {
		y = 0
	}
	if y >= height {
		y = height - 1
	}
	return y
}

// plotDataPoints plots weight data points and connects them with lines
func plotDataPoints(graph [][]rune, weights []models.Weight, wr weightRange, width, height int) {
	// Sample weights if we have more than width
	step := 1
	if len(weights) > width {
		step = len(weights) / width
	}

	prevX, prevY := -1, -1
	for i := 0; i < len(weights); i += step {
		w := weights[i]
		x := (i / step) % width
		y := normalizeToGraphY(w.Weight, wr, height)

		// Draw line from previous point
		if prevX >= 0 {
			drawLine(graph, prevX, prevY, x, y)
		}

		// Mark the point with smallest dot
		if x < width && y < height && y >= 0 {
			graph[y][x] = '·'
		}

		prevX, prevY = x, y
	}
}

// drawGoalLine draws a horizontal line representing the goal weight
func drawGoalLine(graph [][]rune, goalWeight float64, wr weightRange, width, height int) int {
	goalY := normalizeToGraphY(goalWeight, wr, height)
	if goalY >= 0 && goalY < height {
		for x := 0; x < width; x++ {
			// Don't overwrite weight data points
			if graph[goalY][x] != '·' {
				graph[goalY][x] = '─'
			}
		}
	}
	return goalY
}

// renderGraphWithLabels renders the graph grid with axis labels and styling
func renderGraphWithLabels(graph [][]rune, weights []models.Weight, settings *models.Settings, wr weightRange, goalY int) string {
	width := len(graph[0])
	height := len(graph)

	var graphOutput strings.Builder
	graphStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 1)

	var graphLines strings.Builder

	// Add max weight label
	graphLines.WriteString(fmt.Sprintf("%.1f %s ┤\n", wr.max, settings.WeightUnit))

	// Add graph lines with goal weight label
	for i := 0; i < height; i++ {
		if i == goalY {
			// Add goal weight label on the goal line
			goalLabel := fmt.Sprintf("Goal: %.1f", settings.GoalWeight)
			graphLines.WriteString(goalLabel)
			if len(goalLabel) < GoalLabelMinWidth {
				graphLines.WriteString(strings.Repeat(" ", GoalLabelMinWidth-len(goalLabel)))
			}
			graphLines.WriteString("┤")
		} else if i == 0 || i == height-1 {
			graphLines.WriteString("        │")
		} else {
			graphLines.WriteString("        │")
		}
		graphLines.WriteString(string(graph[i]))
		graphLines.WriteString("\n")
	}

	// Add min weight label and x-axis
	graphLines.WriteString(fmt.Sprintf("%.1f %s ┤", wr.min, settings.WeightUnit))
	graphLines.WriteString(strings.Repeat("─", width))
	graphLines.WriteString("\n")

	// Add x-axis labels (date range)
	if len(weights) > 0 {
		oldestDate := weights[0].Date
		newestDate := weights[len(weights)-1].Date
		graphLines.WriteString(fmt.Sprintf("        %s%s%s\n",
			oldestDate,
			strings.Repeat(" ", width-len(oldestDate)-len(newestDate)),
			newestDate))
	}

	graphOutput.WriteString(graphStyle.Render(graphLines.String()))
	return graphOutput.String()
}

// createLineGraph creates a simple ASCII line graph
func createLineGraph(weights []models.Weight, settings *models.Settings) string {
	if len(weights) == 0 {
		return ""
	}

	// Graph dimensions
	width := GraphWidth
	height := GraphHeight

	// Calculate weight range for scaling
	wr := calculateWeightRange(weights, settings.GoalWeight)

	// Create empty graph grid
	graph := createGraphGrid(width, height)

	// Reverse weights to show oldest to newest (left to right)
	reversedWeights := reverseWeights(weights)

	// Plot weight data points and connect them
	plotDataPoints(graph, reversedWeights, wr, width, height)

	// Draw horizontal goal weight line
	goalY := drawGoalLine(graph, settings.GoalWeight, wr, width, height)

	// Render graph with labels and styling
	return renderGraphWithLabels(graph, reversedWeights, settings, wr, goalY)
}

// drawLine draws a line between two points using Bresenham's algorithm
func drawLine(graph [][]rune, x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		// Draw very light connecting line (don't overwrite data points)
		if (x0 != x1 || y0 != y1) && x0 >= 0 && x0 < len(graph[0]) && y0 >= 0 && y0 < len(graph) {
			if graph[y0][x0] != '·' {
				graph[y0][x0] = '∙'
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
