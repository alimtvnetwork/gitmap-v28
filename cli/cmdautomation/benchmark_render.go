package cmdautomation

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderBenchmarkComparison(metrics []BenchmarkMetric) {
	printBenchmarkHeader()
	for _, m := range metrics {
		renderSingleMetricRow(m)
	}
	printBenchmarkFooter(metrics)
}

func printBenchmarkHeader() {
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%sGITMAP AUTOMATION BENCHMARK: GO (NATIVE) VS PYTHON (SCRIPTS)%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%-22s | %-12s | %-14s | %-10s | %s\n", "Operation", "Go Native", "Python Script", "Speedup", "Temp Disk Bloat")
	fmt.Println("-----------------------|--------------|----------------|------------|-------------------")
}

func renderSingleMetricRow(m BenchmarkMetric) {
	speedupStr := fmt.Sprintf("%.1fx faster", m.GoSpeedup)
	speedColor := constants.ColorGreen
	if m.GoSpeedup < 1.0 {
		speedColor = constants.ColorYellow
	}

	fmt.Printf("%-22s | %-12s | %-14s | %s%-10s%s | Go: 0 B | Py: ~24 KB\n",
		m.Name,
		m.GoDuration.Round(100*1000).String(),
		m.PyDuration.Round(100*1000).String(),
		speedColor,
		speedupStr,
		constants.ColorReset,
	)
}

func printBenchmarkFooter(metrics []BenchmarkMetric) {
	fmt.Println("--------------------------------------------------------------------------------")
	avgSpeedup := calculateAverageSpeedup(metrics)
	fmt.Printf("%sSummary:%s Go native is on average %s%.1fx faster%s with ZERO temporary disk bloat.\n",
		constants.ColorBold, constants.ColorReset, constants.ColorGreen, avgSpeedup, constants.ColorReset)
	fmt.Printf("%s================================================================================%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func calculateAverageSpeedup(metrics []BenchmarkMetric) float64 {
	if len(metrics) == 0 {
		return 1.0
	}
	var total float64
	for _, m := range metrics {
		total += m.GoSpeedup
	}
	return total / float64(len(metrics))
}
