package analysis

import (
	"sort"
	"time"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
)

// ============================================================================
// Label Health Types (bv-100)
// Foundation for all label-centric analysis features
// ============================================================================

// LabelHealth represents the overall health assessment of a single label
// Health is a composite score based on velocity, freshness, flow, and criticality
type LabelHealth struct {
	Label       string             `json:"label"`            // The label name
	IssueCount  int                `json:"issue_count"`      // Total issues with this label
	OpenCount   int                `json:"open_count"`       // Open issues with this label
	ClosedCount int                `json:"closed_count"`     // Closed issues with this label
	Blocked     int                `json:"blocked_count"`    // Blocked issues with this label
	Health      int                `json:"health"`           // Composite health score 0-100
	HealthLevel string             `json:"health_level"`     // "healthy", "warning", "critical"
	Velocity    VelocityMetrics    `json:"velocity"`         // Work completion rate
	Freshness   FreshnessMetrics   `json:"freshness"`        // How recently updated
	Flow        FlowMetrics        `json:"flow"`             // Cross-label dependencies
	Criticality CriticalityMetrics `json:"criticality"`      // Graph-based importance
	Issues      []string           `json:"issues,omitempty"` // Issue IDs with this label
}

// VelocityMetrics tracks the rate of work completion for a label
type VelocityMetrics struct {
	ClosedLast7Days  int     `json:"closed_last_7_days"`  // Issues closed in past week
	ClosedLast30Days int     `json:"closed_last_30_days"` // Issues closed in past month
	AvgDaysToClose   float64 `json:"avg_days_to_close"`   // Average time from open to close
	TrendDirection   string  `json:"trend_direction"`     // "improving", "stable", "declining"
	TrendPercent     float64 `json:"trend_percent"`       // Percent change vs prior period
	VelocityScore    int     `json:"velocity_score"`      // Normalized 0-100 score
}

// FreshnessMetrics tracks how recently issues in a label have been updated
type FreshnessMetrics struct {
	MostRecentUpdate   time.Time `json:"most_recent_update"`    // When was any issue last updated
	OldestOpenIssue    time.Time `json:"oldest_open_issue"`     // Created date of oldest open issue
	AvgDaysSinceUpdate float64   `json:"avg_days_since_update"` // Average staleness
	StaleCount         int       `json:"stale_count"`           // Issues with no updates > threshold
	StaleThresholdDays int       `json:"stale_threshold_days"`  // What we consider stale (default 14)
	FreshnessScore     int       `json:"freshness_score"`       // Normalized 0-100 score (higher = fresher)
}

// FlowMetrics captures cross-label dependency relationships
type FlowMetrics struct {
	IncomingDeps      int      `json:"incoming_deps"`       // Other labels blocking this one
	OutgoingDeps      int      `json:"outgoing_deps"`       // Labels this one blocks
	IncomingLabels    []string `json:"incoming_labels"`     // Which labels block this one
	OutgoingLabels    []string `json:"outgoing_labels"`     // Which labels this one blocks
	BlockedByExternal int      `json:"blocked_by_external"` // Issues blocked by other labels
	BlockingExternal  int      `json:"blocking_external"`   // Issues blocking other labels
	FlowScore         int      `json:"flow_score"`          // 0-100, higher = better flow (less blocked)
}

// CriticalityMetrics measures the importance of a label in the dependency graph
type CriticalityMetrics struct {
	AvgPageRank       float64 `json:"avg_pagerank"`        // Average PageRank of issues in label
	AvgBetweenness    float64 `json:"avg_betweenness"`     // Average betweenness centrality
	MaxBetweenness    float64 `json:"max_betweenness"`     // Highest betweenness (bottleneck indicator)
	CriticalPathCount int     `json:"critical_path_count"` // Issues on critical path
	BottleneckCount   int     `json:"bottleneck_count"`    // Issues identified as bottlenecks
	CriticalityScore  int     `json:"criticality_score"`   // 0-100, higher = more critical
}

// LabelDependency represents a dependency relationship between two labels
type LabelDependency struct {
	FromLabel     string         `json:"from_label"`               // Blocking label
	ToLabel       string         `json:"to_label"`                 // Blocked label
	IssueCount    int            `json:"issue_count"`              // Number of cross-label dependencies
	IssueIDs      []string       `json:"issue_ids,omitempty"`      // Specific issue pairs
	BlockingPairs []BlockingPair `json:"blocking_pairs,omitempty"` // Individual blocking relationships
}

// BlockingPair represents a single issue blocking another across labels
type BlockingPair struct {
	BlockerID    string `json:"blocker_id"`    // Issue doing the blocking
	BlockedID    string `json:"blocked_id"`    // Issue being blocked
	BlockerLabel string `json:"blocker_label"` // Label of blocker
	BlockedLabel string `json:"blocked_label"` // Label of blocked
}

// CrossLabelFlow captures the complete flow of work between labels
type CrossLabelFlow struct {
	Labels              []string          `json:"labels"`                 // All labels in analysis
	FlowMatrix          [][]int           `json:"flow_matrix"`            // [from][to] dependency counts
	Dependencies        []LabelDependency `json:"dependencies"`           // Detailed dependency list
	CriticalPaths       []LabelPath       `json:"critical_paths"`         // Label-level critical paths
	BottleneckLabels    []string          `json:"bottleneck_labels"`      // Labels causing most blockage
	TotalCrossLabelDeps int               `json:"total_cross_label_deps"` // Total inter-label dependencies
}

// LabelPath represents a sequence of labels in a dependency chain
type LabelPath struct {
	Labels      []string `json:"labels"`       // Ordered sequence of labels
	Length      int      `json:"length"`       // Number of label transitions
	IssueCount  int      `json:"issue_count"`  // Total issues in this path
	TotalWeight float64  `json:"total_weight"` // Sum of dependency weights
}

// LabelSummary provides a quick overview for display
type LabelSummary struct {
	Label          string `json:"label"`
	IssueCount     int    `json:"issue_count"`
	OpenCount      int    `json:"open_count"`
	Health         int    `json:"health"`              // 0-100
	HealthLevel    string `json:"health_level"`        // "healthy", "warning", "critical"
	TopIssue       string `json:"top_issue,omitempty"` // Highest priority open issue
	NeedsAttention bool   `json:"needs_attention"`     // Flag for labels requiring action
}

// LabelAnalysisResult is the top-level result for label analysis
type LabelAnalysisResult struct {
	GeneratedAt     time.Time       `json:"generated_at"`
	TotalLabels     int             `json:"total_labels"`
	HealthyCount    int             `json:"healthy_count"`              // Labels with health >= 70
	WarningCount    int             `json:"warning_count"`              // Labels with health 40-69
	CriticalCount   int             `json:"critical_count"`             // Labels with health < 40
	Labels          []LabelHealth   `json:"labels"`                     // Detailed per-label health
	Summaries       []LabelSummary  `json:"summaries"`                  // Quick overview list
	CrossLabelFlow  *CrossLabelFlow `json:"cross_label_flow,omitempty"` // Inter-label analysis
	AttentionNeeded []string        `json:"attention_needed"`           // Labels requiring attention
}

// ComputeCrossLabelFlow analyzes blocking dependencies between labels and returns counts.
// It respects cfg.IncludeClosedInFlow: when false, closed issues are ignored.
func ComputeCrossLabelFlow(issues []model.Issue, cfg LabelHealthConfig) CrossLabelFlow {
	labels := ExtractLabels(issues)
	labelList := make([]string, len(labels.Labels))
	copy(labelList, labels.Labels)
	sort.Strings(labelList)

	n := len(labelList)
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	index := make(map[string]int, n)
	for i, l := range labelList {
		index[l] = i
	}

	// Build issue map for lookup
	issueMap := make(map[string]model.Issue, len(issues))
	for _, iss := range issues {
		issueMap[iss.ID] = iss
	}

	// Dependency aggregation
	type pairKey struct{ from, to string }
	depMap := make(map[pairKey]*LabelDependency)
	totalDeps := 0

	for _, blocked := range issues {
		if !cfg.IncludeClosedInFlow && blocked.Status == model.StatusClosed {
			continue
		}
		for _, dep := range blocked.Dependencies {
			if dep == nil || dep.Type != model.DepBlocks {
				continue
			}
			blocker, ok := issueMap[dep.DependsOnID]
			if !ok {
				continue
			}
			if !cfg.IncludeClosedInFlow && blocker.Status == model.StatusClosed {
				continue
			}
			// Cross-product of labels
			for _, from := range blocker.Labels {
				for _, to := range blocked.Labels {
					if from == "" || to == "" || from == to {
						continue // skip empty/self
					}
					iFrom, okFrom := index[from]
					iTo, okTo := index[to]
					if !okFrom || !okTo {
						continue
					}
					matrix[iFrom][iTo]++
					totalDeps++
					key := pairKey{from: from, to: to}
					entry, exists := depMap[key]
					if !exists {
						entry = &LabelDependency{
							FromLabel: key.from,
							ToLabel:   key.to,
							IssueIDs:  []string{},
						}
						depMap[key] = entry
					}
					entry.IssueCount++
					entry.IssueIDs = append(entry.IssueIDs, blocked.ID)
					entry.BlockingPairs = append(entry.BlockingPairs, BlockingPair{
						BlockerID:    blocker.ID,
						BlockedID:    blocked.ID,
						BlockerLabel: from,
						BlockedLabel: to,
					})
				}
			}
		}
	}

	// Build dependency list deterministically
	var deps []LabelDependency
	for _, d := range depMap {
		deps = append(deps, *d)
	}
	sort.Slice(deps, func(i, j int) bool {
		if deps[i].FromLabel != deps[j].FromLabel {
			return deps[i].FromLabel < deps[j].FromLabel
		}
		if deps[i].ToLabel != deps[j].ToLabel {
			return deps[i].ToLabel < deps[j].ToLabel
		}
		return deps[i].IssueCount > deps[j].IssueCount
	})

	// Bottleneck labels: highest outgoing deps
	outCounts := make(map[string]int, n)
	maxOut := 0
	for i, row := range matrix {
		sum := 0
		for _, v := range row {
			sum += v
		}
		outCounts[labelList[i]] = sum
		if sum > maxOut {
			maxOut = sum
		}
	}
	var bottlenecks []string
	for label, c := range outCounts {
		if c == maxOut && c > 0 {
			bottlenecks = append(bottlenecks, label)
		}
	}
	sort.Strings(bottlenecks)

	return CrossLabelFlow{
		Labels:              labelList,
		FlowMatrix:          matrix,
		Dependencies:        deps,
		BottleneckLabels:    bottlenecks,
		TotalCrossLabelDeps: totalDeps,
	}
}

// ComputeVelocityMetrics calculates simple velocity stats for a label.
// It looks at closed issues and recent closures to give a quick pulse.
func ComputeVelocityMetrics(issues []model.Issue, now time.Time) VelocityMetrics {
	const day = 24 * time.Hour
	var closed7, closed30 int
	var totalCloseDur time.Duration
	var closeSamples int

	// Rolling windows
	weekAgo := now.Add(-7 * day)
	monthAgo := now.Add(-30 * day)
	prevWeekStart := now.Add(-14 * day)

	var prevWeek, currentWeek int

	for _, iss := range issues {
		if iss.ClosedAt == nil {
			continue
		}
		closedAt := *iss.ClosedAt
		if closedAt.After(weekAgo) {
			closed7++
		}
		if closedAt.After(monthAgo) {
			closed30++
		}
		if closedAt.After(prevWeekStart) && closedAt.Before(weekAgo) {
			prevWeek++
		} else if closedAt.After(weekAgo) {
			currentWeek++
		}
		if !iss.CreatedAt.IsZero() {
			totalCloseDur += closedAt.Sub(iss.CreatedAt)
			closeSamples++
		}
	}

	avgDays := 0.0
	if closeSamples > 0 {
		avgDays = totalCloseDur.Hours() / 24.0 / float64(closeSamples)
	}

	trendPercent := 0.0
	trendDir := "stable"
	if prevWeek > 0 {
		trendPercent = (float64(currentWeek-prevWeek) / float64(prevWeek)) * 100
		switch {
		case trendPercent > 10:
			trendDir = "improving"
		case trendPercent < -10:
			trendDir = "declining"
		}
	} else if currentWeek > 0 {
		trendDir = "improving"
		trendPercent = 100
	}

	// Simple score: closed in last month scaled plus recency bonus
	velocityScore := 0
	if closed30 > 0 {
		velocityScore = int(minFloat(100, float64(closed30)*10))
	}
	// Bonus if trend improving
	if trendDir == "improving" && velocityScore < 100 {
		velocityScore = clampScore(velocityScore + 10)
	}

	return VelocityMetrics{
		ClosedLast7Days:  closed7,
		ClosedLast30Days: closed30,
		AvgDaysToClose:   avgDays,
		TrendDirection:   trendDir,
		TrendPercent:     trendPercent,
		VelocityScore:    velocityScore,
	}
}

// ComputeFreshnessMetrics calculates freshness and staleness for a label.
func ComputeFreshnessMetrics(issues []model.Issue, now time.Time, staleDays int) FreshnessMetrics {
	if staleDays <= 0 {
		staleDays = DefaultStaleThresholdDays
	}
	var mostRecent time.Time
	var oldestOpen time.Time
	var totalStaleness float64
	var count int
	staleCount := 0
	threshold := float64(staleDays)

	for _, iss := range issues {
		if iss.UpdatedAt.After(mostRecent) {
			mostRecent = iss.UpdatedAt
		}
		if iss.Status != model.StatusClosed {
			if oldestOpen.IsZero() || iss.CreatedAt.Before(oldestOpen) {
				oldestOpen = iss.CreatedAt
			}
		}
		if !iss.UpdatedAt.IsZero() {
			days := now.Sub(iss.UpdatedAt).Hours() / 24.0
			totalStaleness += days
			count++
			if days >= threshold {
				staleCount++
			}
		}
	}

	avgStaleness := 0.0
	if count > 0 {
		avgStaleness = totalStaleness / float64(count)
	}
	// Freshness score: 100 when avg=0, declines linearly to 0 at 2x threshold
	freshnessScore := int(maxFloat(0, 100-(avgStaleness/(threshold*2))*100))

	return FreshnessMetrics{
		MostRecentUpdate:   mostRecent,
		OldestOpenIssue:    oldestOpen,
		AvgDaysSinceUpdate: avgStaleness,
		StaleCount:         staleCount,
		StaleThresholdDays: staleDays,
		FreshnessScore:     clampScore(freshnessScore),
	}
}

// ComputeLabelHealthForLabel computes health for a single label.
// If stats is nil, it will compute graph stats once for the provided issues.
func ComputeLabelHealthForLabel(label string, issues []model.Issue, cfg LabelHealthConfig, now time.Time, stats *GraphStats) LabelHealth {
	health := NewLabelHealth(label)
	health.Issues = []string{}

	// Collect issues with this label
	var labeled []model.Issue
	for _, iss := range issues {
		for _, l := range iss.Labels {
			if l == label {
				labeled = append(labeled, iss)
				health.Issues = append(health.Issues, iss.ID)
				break
			}
		}
	}

	health.IssueCount = len(labeled)
	if health.IssueCount == 0 {
		health.Health = 0
		health.HealthLevel = HealthLevelCritical
		return health
	}

	// Status counts
	for _, iss := range labeled {
		switch iss.Status {
		case model.StatusClosed:
			health.ClosedCount++
		case model.StatusInProgress:
			health.OpenCount++
		case model.StatusBlocked:
			health.Blocked++
		default:
			health.OpenCount++
		}
	}

	velocity := ComputeVelocityMetrics(labeled, now)
	freshness := ComputeFreshnessMetrics(labeled, now, cfg.StaleThresholdDays)

	// Flow: count cross-label deps
	flow := FlowMetrics{}
	seenIn := make(map[string]struct{})
	seenOut := make(map[string]struct{})
	for _, iss := range labeled {
		for _, dep := range iss.Dependencies {
			if dep == nil || dep.Type != model.DepBlocks {
				continue
			}
			blockerLabels := GetLabelsForIssue(issues, dep.DependsOnID)
			targetLabels := iss.Labels
			// incoming: other label blocks this
			for _, bl := range blockerLabels {
				if bl != label {
					flow.IncomingDeps++
					seenIn[bl] = struct{}{}
				}
			}
			// outgoing: this label blocks others
			for _, tl := range targetLabels {
				if tl == label {
					continue
				}
				flow.OutgoingDeps++
				seenOut[tl] = struct{}{}
			}
		}
	}
	for l := range seenIn {
		flow.IncomingLabels = append(flow.IncomingLabels, l)
	}
	for l := range seenOut {
		flow.OutgoingLabels = append(flow.OutgoingLabels, l)
	}
	sort.Strings(flow.IncomingLabels)
	sort.Strings(flow.OutgoingLabels)
	flow.FlowScore = clampScore(100 - (flow.IncomingDeps * 5))

	// Criticality: derive from graph metrics (reuse precomputed stats when supplied)
	if stats == nil {
		analyzer := NewAnalyzer(issues)
		s := analyzer.Analyze()
		stats = &s
	}
	pr := stats.PageRank()
	bw := stats.Betweenness()
	maxPR := findMax(pr)
	maxBW := findMax(bw)

	var prSum, bwSum float64
	maxBwLabel := 0.0
	var critCount, bottleneckCount int
	for _, iss := range labeled {
		prSum += pr[iss.ID]
		bwVal := bw[iss.ID]
		bwSum += bwVal
		if bwVal > maxBwLabel {
			maxBwLabel = bwVal
		}
		if stats.GetCriticalPathScore(iss.ID) > 0 {
			critCount++
		}
		if bwVal > 0 {
			bottleneckCount++
		}
	}
	avgPR := 0.0
	avgBW := 0.0
	if health.IssueCount > 0 {
		avgPR = prSum / float64(health.IssueCount)
		avgBW = bwSum / float64(health.IssueCount)
	}
	critScore := 0
	if maxPR > 0 {
		critScore += int((avgPR / maxPR) * 50)
	}
	if maxBW > 0 {
		critScore += int((maxBwLabel / maxBW) * 50)
	}
	critScore = clampScore(critScore)

	health.Velocity = velocity
	health.Freshness = freshness
	health.Flow = flow
	health.Criticality = CriticalityMetrics{
		AvgPageRank:       avgPR,
		AvgBetweenness:    avgBW,
		MaxBetweenness:    maxBwLabel,
		CriticalPathCount: critCount,
		BottleneckCount:   bottleneckCount,
		CriticalityScore:  critScore,
	}

	health.Health = ComputeCompositeHealth(velocity.VelocityScore, freshness.FreshnessScore, flow.FlowScore, critScore, cfg)
	health.HealthLevel = HealthLevelFromScore(health.Health)
	return health
}

// ComputeAllLabelHealth computes health for all labels in the issue set.
func ComputeAllLabelHealth(issues []model.Issue, cfg LabelHealthConfig, now time.Time) LabelAnalysisResult {
	labels := ExtractLabels(issues)
	result := LabelAnalysisResult{
		GeneratedAt:     now,
		TotalLabels:     labels.LabelCount,
		Labels:          []LabelHealth{},
		Summaries:       []LabelSummary{},
		AttentionNeeded: []string{},
	}

	// Deterministic traversal
	sort.Strings(labels.Labels)

	// Precompute stats once for efficiency
	analyzer := NewAnalyzer(issues)
	fullStats := analyzer.Analyze()

	for _, label := range labels.Labels {
		health := ComputeLabelHealthForLabel(label, issues, cfg, now, &fullStats)
		result.Labels = append(result.Labels, health)
		summary := LabelSummary{
			Label:          label,
			IssueCount:     health.IssueCount,
			OpenCount:      health.OpenCount,
			Health:         health.Health,
			HealthLevel:    health.HealthLevel,
			NeedsAttention: NeedsAttention(health),
		}
		if len(health.Issues) > 0 {
			summary.TopIssue = health.Issues[0]
		}
		result.Summaries = append(result.Summaries, summary)
		switch health.HealthLevel {
		case HealthLevelHealthy:
			result.HealthyCount++
		case HealthLevelWarning:
			result.WarningCount++
			result.AttentionNeeded = append(result.AttentionNeeded, label)
		case HealthLevelCritical:
			result.CriticalCount++
			result.AttentionNeeded = append(result.AttentionNeeded, label)
		}
	}

	sort.Slice(result.Summaries, func(i, j int) bool {
		if result.Summaries[i].Health != result.Summaries[j].Health {
			return result.Summaries[i].Health > result.Summaries[j].Health
		}
		return result.Summaries[i].Label < result.Summaries[j].Label
	})

	return result
}

func clampScore(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ============================================================================
// Health Score Constants and Thresholds
// ============================================================================

// HealthLevel constants for categorizing label health
const (
	HealthLevelHealthy  = "healthy"  // Health >= 70
	HealthLevelWarning  = "warning"  // Health 40-69
	HealthLevelCritical = "critical" // Health < 40
)

// Default thresholds for health calculations
const (
	DefaultStaleThresholdDays = 14   // Days without update to consider stale
	HealthyThreshold          = 70   // Min health score for "healthy"
	WarningThreshold          = 40   // Min health score for "warning"
	VelocityWeight            = 0.25 // Weight for velocity in composite score
	FreshnessWeight           = 0.25 // Weight for freshness in composite score
	FlowWeight                = 0.25 // Weight for flow in composite score
	CriticalityWeight         = 0.25 // Weight for criticality in composite score
)

// ============================================================================
// Configuration Types
// ============================================================================

// LabelHealthConfig configures label health computation
type LabelHealthConfig struct {
	StaleThresholdDays  int     `json:"stale_threshold_days"`   // Days to consider issue stale
	VelocityWeight      float64 `json:"velocity_weight"`        // Weight for velocity component
	FreshnessWeight     float64 `json:"freshness_weight"`       // Weight for freshness component
	FlowWeight          float64 `json:"flow_weight"`            // Weight for flow component
	CriticalityWeight   float64 `json:"criticality_weight"`     // Weight for criticality component
	MinIssuesForHealth  int     `json:"min_issues_for_health"`  // Min issues to compute health
	IncludeClosedInFlow bool    `json:"include_closed_in_flow"` // Include closed issues in flow analysis
}

// DefaultLabelHealthConfig returns sensible defaults
func DefaultLabelHealthConfig() LabelHealthConfig {
	return LabelHealthConfig{
		StaleThresholdDays:  DefaultStaleThresholdDays,
		VelocityWeight:      VelocityWeight,
		FreshnessWeight:     FreshnessWeight,
		FlowWeight:          FlowWeight,
		CriticalityWeight:   CriticalityWeight,
		MinIssuesForHealth:  1,
		IncludeClosedInFlow: false,
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// HealthLevelFromScore returns the health level string for a score
func HealthLevelFromScore(score int) string {
	if score >= HealthyThreshold {
		return HealthLevelHealthy
	}
	if score >= WarningThreshold {
		return HealthLevelWarning
	}
	return HealthLevelCritical
}

// NeedsAttention returns true if a label needs attention based on health
func NeedsAttention(health LabelHealth) bool {
	return health.Health < HealthyThreshold
}

// ComputeCompositeHealth calculates the overall health score from components
func ComputeCompositeHealth(velocity, freshness, flow, criticality int, cfg LabelHealthConfig) int {
	weighted := float64(velocity)*cfg.VelocityWeight +
		float64(freshness)*cfg.FreshnessWeight +
		float64(flow)*cfg.FlowWeight +
		float64(criticality)*cfg.CriticalityWeight

	// Normalize to 0-100 and clamp
	return clampScore(int(weighted + 0.5))
}

// NewLabelHealth creates a new LabelHealth with default values
func NewLabelHealth(label string) LabelHealth {
	return LabelHealth{
		Label:       label,
		Health:      100, // Start healthy, reduce based on issues
		HealthLevel: HealthLevelHealthy,
		Velocity: VelocityMetrics{
			TrendDirection: "stable",
			VelocityScore:  100,
		},
		Freshness: FreshnessMetrics{
			StaleThresholdDays: DefaultStaleThresholdDays,
			FreshnessScore:     100,
		},
		Flow: FlowMetrics{
			FlowScore: 100,
		},
		Criticality: CriticalityMetrics{
			CriticalityScore: 50, // Neutral starting point
		},
	}
}

// ============================================================================
// Label Extraction (bv-101)
// Functions to extract and aggregate labels from issues
// ============================================================================

// LabelStats provides basic statistics about a label
type LabelStats struct {
	Label       string         `json:"label"`
	TotalCount  int            `json:"total_count"`  // Total issues with this label
	OpenCount   int            `json:"open_count"`   // Open issues
	ClosedCount int            `json:"closed_count"` // Closed issues
	InProgress  int            `json:"in_progress"`  // In-progress issues
	Blocked     int            `json:"blocked"`      // Issues blocked by dependencies
	ByPriority  map[int]int    `json:"by_priority"`  // Count by priority level
	ByType      map[string]int `json:"by_type"`      // Count by issue type
	IssueIDs    []string       `json:"issue_ids"`    // All issue IDs with this label
}

// LabelExtractionResult contains all extracted label data
type LabelExtractionResult struct {
	Labels         []string               `json:"labels"`          // Unique labels in sorted order
	LabelCount     int                    `json:"label_count"`     // Number of unique labels
	Stats          map[string]*LabelStats `json:"stats"`           // Per-label statistics
	IssueCount     int                    `json:"issue_count"`     // Total issues analyzed
	UnlabeledCount int                    `json:"unlabeled_count"` // Issues without labels
	TopLabels      []string               `json:"top_labels"`      // Labels sorted by issue count
}

// ExtractLabels extracts unique labels from a slice of issues with statistics
// Handles edge cases: nil issues, empty labels, duplicate labels
func ExtractLabels(issues []model.Issue) LabelExtractionResult {
	result := LabelExtractionResult{
		Stats:     make(map[string]*LabelStats),
		Labels:    []string{},
		TopLabels: []string{},
	}

	if len(issues) == 0 {
		return result
	}

	result.IssueCount = len(issues)
	labelSet := make(map[string]bool)

	for _, issue := range issues {
		// Track issues without labels
		if len(issue.Labels) == 0 {
			result.UnlabeledCount++
		}

		// Process each label on the issue
		for _, label := range issue.Labels {
			// Skip empty labels
			if label == "" {
				continue
			}

			// Track unique labels
			labelSet[label] = true

			// Initialize stats if needed
			stats, exists := result.Stats[label]
			if !exists {
				stats = &LabelStats{
					Label:      label,
					ByPriority: make(map[int]int),
					ByType:     make(map[string]int),
					IssueIDs:   []string{},
				}
				result.Stats[label] = stats
			}

			// Update counts
			stats.TotalCount++
			stats.IssueIDs = append(stats.IssueIDs, issue.ID)

			// Count by status
			switch issue.Status {
			case model.StatusOpen:
				stats.OpenCount++
			case model.StatusClosed:
				stats.ClosedCount++
			case model.StatusInProgress:
				stats.InProgress++
			}

			// Count by priority
			stats.ByPriority[issue.Priority]++

			// Count by type
			stats.ByType[string(issue.IssueType)]++
		}
	}

	// Build sorted label list
	for label := range labelSet {
		result.Labels = append(result.Labels, label)
	}
	sort.Strings(result.Labels)
	result.LabelCount = len(result.Labels)

	// Build top labels by issue count
	result.TopLabels = sortLabelsByCount(result.Stats)

	return result
}

// sortLabelsByCount returns labels sorted by total issue count (descending)
func sortLabelsByCount(stats map[string]*LabelStats) []string {
	type labelCount struct {
		label string
		count int
	}

	var lc []labelCount
	for label, s := range stats {
		lc = append(lc, labelCount{label: label, count: s.TotalCount})
	}

	sort.Slice(lc, func(i, j int) bool {
		if lc[i].count != lc[j].count {
			return lc[i].count > lc[j].count
		}
		return lc[i].label < lc[j].label // Alphabetical for ties
	})

	result := make([]string, len(lc))
	for i, l := range lc {
		result[i] = l.label
	}
	return result
}

// GetLabelIssues returns all issues that have a specific label
func GetLabelIssues(issues []model.Issue, label string) []model.Issue {
	var result []model.Issue
	for _, issue := range issues {
		for _, l := range issue.Labels {
			if l == label {
				result = append(result, issue)
				break
			}
		}
	}
	return result
}

// GetLabelsForIssue returns all labels for a specific issue ID
func GetLabelsForIssue(issues []model.Issue, issueID string) []string {
	for _, issue := range issues {
		if issue.ID == issueID {
			return issue.Labels
		}
	}
	return nil
}

// GetCommonLabels returns labels that appear in multiple provided label sets
func GetCommonLabels(labelSets ...[]string) []string {
	if len(labelSets) == 0 {
		return nil
	}

	// Count occurrences
	counts := make(map[string]int)
	for _, labels := range labelSets {
		seen := make(map[string]bool)
		for _, label := range labels {
			if !seen[label] {
				counts[label]++
				seen[label] = true
			}
		}
	}

	// Find labels in all sets
	var common []string
	for label, count := range counts {
		if count == len(labelSets) {
			common = append(common, label)
		}
	}
	sort.Strings(common)
	return common
}

// GetLabelCooccurrence builds a co-occurrence matrix showing which labels appear together
func GetLabelCooccurrence(issues []model.Issue) map[string]map[string]int {
	cooc := make(map[string]map[string]int)

	for _, issue := range issues {
		labels := issue.Labels
		// For each pair of labels on the same issue
		for i := 0; i < len(labels); i++ {
			for j := i + 1; j < len(labels); j++ {
				l1, l2 := labels[i], labels[j]
				// Ensure consistent ordering
				if l1 > l2 {
					l1, l2 = l2, l1
				}

				// Initialize maps if needed
				if cooc[l1] == nil {
					cooc[l1] = make(map[string]int)
				}
				if cooc[l2] == nil {
					cooc[l2] = make(map[string]int)
				}

				// Increment both directions
				cooc[l1][l2]++
				cooc[l2][l1]++
			}
		}
	}

	return cooc
}

// ComputeBlockedByLabel determines which issues are blocked, grouped by label
// Returns a map of label -> count of blocked issues with that label
func ComputeBlockedByLabel(issues []model.Issue, analyzer *Analyzer) map[string]int {
	blocked := make(map[string]int)

	for _, issue := range issues {
		if issue.Status == model.StatusClosed {
			continue
		}

		// Check if issue is blocked
		blockers := analyzer.GetOpenBlockers(issue.ID)
		if len(blockers) > 0 {
			// This issue is blocked - count for each of its labels
			for _, label := range issue.Labels {
				blocked[label]++
			}
		}
	}

	return blocked
}
