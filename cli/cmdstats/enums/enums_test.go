package enums

import "testing"

func TestStatsOutputFormatType_Valid(t *testing.T) {
	formats := AllStatsOutputFormats()
	if len(formats) != 4 {
		t.Fatalf("expected 4 output formats, got %d", len(formats))
	}

	for _, f := range formats {
		if !f.IsValid() {
			t.Errorf("expected %s to be valid", f)
		}

		if f.String() != string(f) {
			t.Errorf("expected String() == %s, got %s", string(f), f.String())
		}
	}

	if StatsOutputFormatType("invalid").IsValid() {
		t.Errorf("expected invalid format to be invalid")
	}
}

func TestStatsOutputFormatType_Booleans(t *testing.T) {
	if !StatsOutputFormatTable.IsTable() || StatsOutputFormatJSON.IsTable() {
		t.Errorf("unexpected IsTable behavior")
	}

	if !StatsOutputFormatJSON.IsJSON() || StatsOutputFormatCSV.IsJSON() {
		t.Errorf("unexpected IsJSON behavior")
	}

	if !StatsOutputFormatTUI.IsTUI() || StatsOutputFormatTable.IsTUI() {
		t.Errorf("unexpected IsTUI behavior")
	}

	if !StatsOutputFormatCSV.IsCSV() || StatsOutputFormatTUI.IsCSV() {
		t.Errorf("unexpected IsCSV behavior")
	}

	if !StatsOutputFormatTable.HasFormatting() || StatsOutputFormatType("").HasFormatting() {
		t.Errorf("unexpected HasFormatting behavior")
	}
}

func TestStatsOutputFormatType_Parse(t *testing.T) {
	parsed, err := ParseStatsOutputFormat("json")
	if err != nil || parsed != StatsOutputFormatJSON {
		t.Errorf("expected json, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatsOutputFormatType(" TABLE ")
	if errAlias != nil || parsedAlias != StatsOutputFormatTable {
		t.Errorf("expected table, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatsOutputFormat("unknown_fmt")
	if invalidErr == nil {
		t.Errorf("expected error for unknown output format")
	}

	if len(AllStatsOutputFormatTypes()) != 4 {
		t.Errorf("expected AllStatsOutputFormatTypes count to match")
	}
}

func TestStatsMetricType_Valid(t *testing.T) {
	metrics := AllStatsMetrics()
	if len(metrics) != 10 {
		t.Fatalf("expected 10 metrics, got %d", len(metrics))
	}

	for _, m := range metrics {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("expected String() == %s, got %s", string(m), m.String())
		}
	}

	if StatsMetricType("not_a_metric").IsValid() {
		t.Errorf("expected invalid metric to be invalid")
	}
}

func TestStatsMetricType_Booleans(t *testing.T) {
	if !StatsMetricTotalRuns.IsTotalRuns() || StatsMetricSuccessCount.IsTotalRuns() {
		t.Errorf("unexpected IsTotalRuns behavior")
	}

	if !StatsMetricSuccessCount.IsSuccessCount() || StatsMetricFailCount.IsSuccessCount() {
		t.Errorf("unexpected IsSuccessCount behavior")
	}

	if !StatsMetricFailCount.IsFailCount() || StatsMetricFailRate.IsFailCount() {
		t.Errorf("unexpected IsFailCount behavior")
	}

	if !StatsMetricFailRate.IsFailRate() || StatsMetricTotalRuns.IsFailRate() {
		t.Errorf("unexpected IsFailRate behavior")
	}

	if !StatsMetricAvgDurationMs.IsAvgDurationMs() || StatsMetricMinDurationMs.IsAvgDurationMs() {
		t.Errorf("unexpected IsAvgDurationMs behavior")
	}

	if !StatsMetricMinDurationMs.IsMinDurationMs() || StatsMetricMaxDurationMs.IsMinDurationMs() {
		t.Errorf("unexpected IsMinDurationMs behavior")
	}

	if !StatsMetricMaxDurationMs.IsMaxDurationMs() || StatsMetricAvgDurationMs.IsMaxDurationMs() {
		t.Errorf("unexpected IsMaxDurationMs behavior")
	}

	if !StatsMetricTotalCommands.IsTotalCommands() || StatsMetricUniqueCommands.IsTotalCommands() {
		t.Errorf("unexpected IsTotalCommands behavior")
	}

	if !StatsMetricUniqueCommands.IsUniqueCommands() || StatsMetricTotalCommands.IsUniqueCommands() {
		t.Errorf("unexpected IsUniqueCommands behavior")
	}

	if !StatsMetricLastUsed.IsLastUsed() || StatsMetricTotalRuns.IsLastUsed() {
		t.Errorf("unexpected IsLastUsed behavior")
	}

	if !StatsMetricAvgDurationMs.HasDuration() || !StatsMetricMinDurationMs.HasDuration() ||
		!StatsMetricMaxDurationMs.HasDuration() || StatsMetricTotalRuns.HasDuration() {
		t.Errorf("unexpected HasDuration behavior")
	}

	if !StatsMetricTotalRuns.HasCount() || !StatsMetricSuccessCount.HasCount() ||
		!StatsMetricFailCount.HasCount() || !StatsMetricTotalCommands.HasCount() ||
		!StatsMetricUniqueCommands.HasCount() || StatsMetricFailRate.HasCount() {
		t.Errorf("unexpected HasCount behavior")
	}
}

func TestStatsMetricType_Parse(t *testing.T) {
	parsed, err := ParseStatsMetric("total_runs")
	if err != nil || parsed != StatsMetricTotalRuns {
		t.Errorf("expected total_runs, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatsMetricType(" FAIL_RATE ")
	if errAlias != nil || parsedAlias != StatsMetricFailRate {
		t.Errorf("expected fail_rate, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatsMetric("bogus_metric")
	if invalidErr == nil {
		t.Errorf("expected error for bogus metric")
	}

	if len(AllStatsMetricTypes()) != 10 {
		t.Errorf("expected AllStatsMetricTypes count to match")
	}
}

func TestStatsTimeRangeType_Valid(t *testing.T) {
	ranges := AllStatsTimeRanges()
	if len(ranges) != 5 {
		t.Fatalf("expected 5 time ranges, got %d", len(ranges))
	}

	for _, r := range ranges {
		if !r.IsValid() {
			t.Errorf("expected %s to be valid", r)
		}

		if r.String() != string(r) {
			t.Errorf("expected String() == %s, got %s", string(r), r.String())
		}
	}

	if StatsTimeRangeType("999d").IsValid() {
		t.Errorf("expected invalid time range to be invalid")
	}
}

func TestStatsTimeRangeType_Booleans(t *testing.T) {
	if !StatsTimeRangeAll.IsAll() || StatsTimeRangeToday.IsAll() {
		t.Errorf("unexpected IsAll behavior")
	}

	if !StatsTimeRangeToday.IsToday() || StatsTimeRange7D.IsToday() {
		t.Errorf("unexpected IsToday behavior")
	}

	if !StatsTimeRange7D.Is7D() || StatsTimeRange30D.Is7D() {
		t.Errorf("unexpected Is7D behavior")
	}

	if !StatsTimeRange30D.Is30D() || StatsTimeRange90D.Is30D() {
		t.Errorf("unexpected Is30D behavior")
	}

	if !StatsTimeRange90D.Is90D() || StatsTimeRangeAll.Is90D() {
		t.Errorf("unexpected Is90D behavior")
	}

	if StatsTimeRangeAll.HasBoundedTime() || !StatsTimeRange7D.HasBoundedTime() ||
		!StatsTimeRangeToday.HasBoundedTime() {
		t.Errorf("unexpected HasBoundedTime behavior")
	}
}

func TestStatsTimeRangeType_Parse(t *testing.T) {
	parsed, err := ParseStatsTimeRange("7d")
	if err != nil || parsed != StatsTimeRange7D {
		t.Errorf("expected 7d, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatsTimeRangeType(" TODAY ")
	if errAlias != nil || parsedAlias != StatsTimeRangeToday {
		t.Errorf("expected today, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatsTimeRange("1year")
	if invalidErr == nil {
		t.Errorf("expected error for invalid range")
	}

	if len(AllStatsTimeRangeTypes()) != 5 {
		t.Errorf("expected AllStatsTimeRangeTypes count to match")
	}
}

func TestStatsCategoryType_Valid(t *testing.T) {
	cats := AllStatsCategories()
	if len(cats) != 6 {
		t.Fatalf("expected 6 categories, got %d", len(cats))
	}

	for _, c := range cats {
		if !c.IsValid() {
			t.Errorf("expected %s to be valid", c)
		}

		if c.String() != string(c) {
			t.Errorf("expected String() == %s, got %s", string(c), c.String())
		}
	}

	if StatsCategoryType("database").IsValid() {
		t.Errorf("expected invalid category to be invalid")
	}
}

func TestStatsCategoryType_Booleans(t *testing.T) {
	if !StatsCategoryAll.IsAll() || StatsCategoryRepository.IsAll() {
		t.Errorf("unexpected IsAll behavior")
	}

	if !StatsCategoryRepository.IsRepository() || StatsCategoryRemote.IsRepository() {
		t.Errorf("unexpected IsRepository behavior")
	}

	if !StatsCategoryRemote.IsRemote() || StatsCategorySSH.IsRemote() {
		t.Errorf("unexpected IsRemote behavior")
	}

	if !StatsCategorySSH.IsSSH() || StatsCategoryAutomation.IsSSH() {
		t.Errorf("unexpected IsSSH behavior")
	}

	if !StatsCategoryAutomation.IsAutomation() || StatsCategoryMaintenance.IsAutomation() {
		t.Errorf("unexpected IsAutomation behavior")
	}

	if !StatsCategoryMaintenance.IsMaintenance() || StatsCategoryAll.IsMaintenance() {
		t.Errorf("unexpected IsMaintenance behavior")
	}

	if StatsCategoryAll.HasCategoryFilter() || !StatsCategorySSH.HasCategoryFilter() {
		t.Errorf("unexpected HasCategoryFilter behavior")
	}
}

func TestStatsCategoryType_Parse(t *testing.T) {
	parsed, err := ParseStatsCategory("ssh")
	if err != nil || parsed != StatsCategorySSH {
		t.Errorf("expected ssh, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatsCategoryType(" MAINTENANCE ")
	if errAlias != nil || parsedAlias != StatsCategoryMaintenance {
		t.Errorf("expected maintenance, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatsCategory("gaming")
	if invalidErr == nil {
		t.Errorf("expected error for invalid category")
	}

	if len(AllStatsCategoryTypes()) != 6 {
		t.Errorf("expected AllStatsCategoryTypes count to match")
	}
}

func TestStatsChartModeType_Valid(t *testing.T) {
	modes := AllStatsChartModes()
	if len(modes) != 4 {
		t.Fatalf("expected 4 chart modes, got %d", len(modes))
	}

	for _, m := range modes {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("expected String() == %s, got %s", string(m), m.String())
		}
	}

	if StatsChartModeType("pie").IsValid() {
		t.Errorf("expected invalid chart mode to be invalid")
	}
}

func TestStatsChartModeType_Booleans(t *testing.T) {
	if !StatsChartModeNone.IsNone() || StatsChartModeBar.IsNone() {
		t.Errorf("unexpected IsNone behavior")
	}

	if !StatsChartModeBar.IsBar() || StatsChartModeSparkline.IsBar() {
		t.Errorf("unexpected IsBar behavior")
	}

	if !StatsChartModeSparkline.IsSparkline() || StatsChartModeHistogram.IsSparkline() {
		t.Errorf("unexpected IsSparkline behavior")
	}

	if !StatsChartModeHistogram.IsHistogram() || StatsChartModeNone.IsHistogram() {
		t.Errorf("unexpected IsHistogram behavior")
	}

	if StatsChartModeNone.HasVisualChart() || !StatsChartModeBar.HasVisualChart() ||
		!StatsChartModeSparkline.HasVisualChart() || !StatsChartModeHistogram.HasVisualChart() {
		t.Errorf("unexpected HasVisualChart behavior")
	}
}

func TestStatsChartModeType_Parse(t *testing.T) {
	parsed, err := ParseStatsChartMode("bar")
	if err != nil || parsed != StatsChartModeBar {
		t.Errorf("expected bar, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseStatsChartModeType(" SPARKLINE ")
	if errAlias != nil || parsedAlias != StatsChartModeSparkline {
		t.Errorf("expected sparkline, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseStatsChartMode("scatter")
	if invalidErr == nil {
		t.Errorf("expected error for invalid chart mode")
	}

	if len(AllStatsChartModeTypes()) != 4 {
		t.Errorf("expected AllStatsChartModeTypes count to match")
	}
}
