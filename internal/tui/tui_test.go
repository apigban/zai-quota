package tui

import (
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zai-quota/internal/config"
	"zai-quota/internal/models"
)

func TestNewModel(t *testing.T) {
	cfg := &config.Config{
		APIKey:         "test-key",
		Endpoint:       "https://api.test.com",
		TimeoutSeconds: 5,
	}

	m := NewModel(cfg, false)

	assert.Equal(t, stateLoading, m.state)
	assert.NotNil(t, m.client)
	assert.Equal(t, cfg, m.cfg)
	assert.NotNil(t, m.processed)
	assert.NotNil(t, m.expanded)
	assert.NotNil(t, m.limits)
}

func TestNewModel_EmptyAPIKey(t *testing.T) {
	cfg := &config.Config{
		APIKey:         "",
		Endpoint:       "https://api.test.com",
		TimeoutSeconds: 5,
	}

	m := NewModel(cfg, false)

	assert.Equal(t, stateSetup, m.state, "empty API key should trigger setup state")
	assert.NotNil(t, m.setupInput, "setup input should be initialized")
}
func TestNewModel_WithAPIKey(t *testing.T) {
	cfg := &config.Config{
		APIKey:         "valid-api-key",
		Endpoint:       "https://api.test.com",
		TimeoutSeconds: 5,
	}

	m := NewModel(cfg, false)

	assert.Equal(t, stateLoading, m.state, "with API key should start in loading state")
}

func TestModel_Init(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Endpoint: "https://api.test.com",
	}
	m := NewModel(cfg, false)

	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func TestModel_Update_KeyPressQuit(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)

	tests := []struct {
		name string
		code rune
		mod  tea.KeyMod
	}{
		{"q", 'q', 0},
		{"Q", 'Q', 0},
		{"ctrl+c", 'c', tea.ModCtrl},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := tea.Key{Code: tt.code, Mod: tt.mod}
			msg := tea.KeyPressMsg(key)

			updatedModel, cmd := m.Update(msg)
			assert.NotNil(t, cmd, "expected quit command for key %s", tt.name)
			_ = updatedModel
		})
	}
}

func TestModel_Update_KeyPressRefresh(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoaded

	updatedModel, cmd := m.Update(tea.KeyPressMsg{})
	_ = updatedModel
	_ = cmd
}

func TestModel_Update_KeyPressExpand(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoaded
	m.processed = map[string]ProcessedLimitData{
		"TIME_LIMIT": {
			UsageDetails: []models.UsageDetail{
				{ModelCode: "test-model", Usage: 5},
			},
		},
	}

	key := tea.Key{Code: 'e'}
	msg := tea.KeyPressMsg(key)

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	assert.True(t, updated.expanded["TIME_LIMIT"])

	updatedModel2, _ := updated.Update(msg)
	updated2 := updatedModel2.(Model)

	assert.False(t, updated2.expanded["TIME_LIMIT"])
}

func TestModel_Update_RefreshWhileRefreshing(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateRefreshing

	updatedModel, cmd := m.Update(tea.KeyPressMsg{})
	assert.Nil(t, cmd)
	_ = updatedModel
}

func TestModel_Update_WindowSize(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	assert.Equal(t, 80, updated.width)
	assert.Equal(t, 24, updated.height)
}

func TestModel_Update_QuotaFetchedSuccess(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoading

	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{
				Type:          "TIME_LIMIT",
				Percentage:    25,
				Usage:         1000,
				CurrentValue:  250,
				Remaining:     750,
				NextResetTime: time.Now().Add(5 * time.Hour).Unix(),
			},
		},
		Level: "pro",
	}

	msg := quotaFetchedMsg{quota: quota}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	assert.Equal(t, stateLoaded, updated.state)
	assert.NotNil(t, updated.limits)
	assert.Equal(t, "pro", updated.level)
	assert.NotNil(t, updated.processed)
	assert.Nil(t, updated.err)
}

func TestModel_Update_QuotaFetchedError(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoading

	testErr := errors.New("network error")
	msg := quotaFetchedMsg{err: testErr}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	assert.Equal(t, stateError, updated.state)
	assert.Equal(t, testErr, updated.err)
}

func TestModel_View_Loading(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoading

	view := m.View()
	content := view.Content

	assert.Contains(t, content, "Z.AI Quota Monitor")
}

func TestModel_View_Loaded(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateLoaded
	m.level = "pro"
	m.processed = map[string]ProcessedLimitData{
		"TIME_LIMIT": {
			Label:        "[5-Hour Prompt Limit]",
			Percentage:   25,
			Used:         250,
			Total:        1000,
			Remaining:    750,
			ResetDisplay: "2h 30m",
		},
	}

	view := m.View()
	content := view.Content

	assert.Contains(t, content, "Z.AI Quota Monitor")
	assert.Contains(t, content, "[Pro]")
	assert.Contains(t, content, "5-Hour Prompt Limit")
	assert.Contains(t, content, "250 / 1000")
}

func TestModel_View_ErrorNoData(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateError
	m.err = errors.New("connection failed")
	m.processed = nil

	view := m.View()
	content := view.Content

	assert.Contains(t, content, "Error:")
	assert.Contains(t, content, "No quota data available")
}

func TestModel_View_ErrorWithData(t *testing.T) {
	cfg := &config.Config{APIKey: "test"}
	m := NewModel(cfg, false)
	m.state = stateError
	m.err = errors.New("connection failed")
	m.processed = map[string]ProcessedLimitData{
		"TIME_LIMIT": {
			Label:        "[5-Hour Prompt Limit]",
			Percentage:   50,
			Used:         500,
			Total:        1000,
			ResetDisplay: "1h",
		},
	}

	view := m.View()
	content := view.Content

	assert.Contains(t, content, "Error:")
	assert.Contains(t, content, "5-Hour Prompt Limit")
	assert.Contains(t, content, "[r] Refresh")
}

func TestProcessLimits(t *testing.T) {
	tests := []struct {
		name      string
		limits    []models.Limit
		wantType  string
		wantPct   int
		wantLevel string
	}{
		{
			name: "TIME_LIMIT 25 percent usage",
			limits: []models.Limit{
				{
					Type:          "TIME_LIMIT",
					Percentage:    25,
					Usage:         1000,
					CurrentValue:  250,
					Remaining:     750,
					NextResetTime: time.Now().Add(5 * time.Hour).Unix(),
				},
			},
			wantType:  "TIME_LIMIT",
			wantPct:   25,
			wantLevel: "safe",
		},
		{
			name: "TIME_LIMIT 85 percent usage",
			limits: []models.Limit{
				{
					Type:          "TIME_LIMIT",
					Percentage:    85,
					Usage:         1000,
					CurrentValue:  850,
					Remaining:     150,
					NextResetTime: time.Now().Add(1 * time.Hour).Unix(),
				},
			},
			wantType:  "TIME_LIMIT",
			wantPct:   85,
			wantLevel: "warning",
		},
		{
			name: "TIME_LIMIT 92 percent usage",
			limits: []models.Limit{
				{
					Type:          "TIME_LIMIT",
					Percentage:    92,
					Usage:         1000,
					CurrentValue:  920,
					Remaining:     80,
					NextResetTime: time.Now().Add(30 * time.Minute).Unix(),
				},
			},
			wantType:  "TIME_LIMIT",
			wantPct:   92,
			wantLevel: "critical",
		},
		{
			name: "TIME_LIMIT 97 percent usage",
			limits: []models.Limit{
				{
					Type:          "TIME_LIMIT",
					Percentage:    97,
					Usage:         1000,
					CurrentValue:  970,
					Remaining:     30,
					NextResetTime: time.Now().Add(10 * time.Minute).Unix(),
				},
			},
			wantType:  "TIME_LIMIT",
			wantPct:   97,
			wantLevel: "emergency",
		},
		{
			name: "TOKENS_LIMIT",
			limits: []models.Limit{
				{
					Type:          "TOKENS_LIMIT",
					Percentage:    33,
					NextResetTime: time.Now().Add(5 * 24 * time.Hour).Unix(),
				},
			},
			wantType:  "TOKENS_LIMIT",
			wantPct:   33,
			wantLevel: "safe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processLimits(tt.limits)
			require.NotNil(t, result)
			data, exists := result[tt.wantType]
			require.True(t, exists)
			assert.Equal(t, tt.wantPct, data.Percentage)
			assert.Equal(t, tt.wantLevel, data.WarningLevel)
		})
	}
}

func TestProcessLimits_Empty(t *testing.T) {
	result := processLimits([]models.Limit{})
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestGetWarningLevel(t *testing.T) {
	tests := []struct {
		pct  int
		want string
	}{
		{0, "safe"},
		{50, "safe"},
		{79, "safe"},
		{80, "warning"},
		{85, "warning"},
		{89, "warning"},
		{90, "critical"},
		{92, "critical"},
		{94, "critical"},
		{95, "emergency"},
		{100, "emergency"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			level := getWarningLevel(tt.pct)
			assert.Equal(t, tt.want, level)
		})
	}
}

func TestGetColorForPercentage(t *testing.T) {
	tests := []struct {
		pct  int
		name string
	}{
		{0, "safe"},
		{50, "safe"},
		{79, "safe"},
		{80, "warning"},
		{85, "warning"},
		{89, "warning"},
		{90, "critical"},
		{92, "critical"},
		{94, "critical"},
		{95, "emergency"},
		{100, "emergency"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := getColorForPercentage(tt.pct)
			assert.NotNil(t, color)
		})
	}
}

func TestRenderProgressBar(t *testing.T) {
	m := Model{}
	result := m.renderProgressBar(50, colorSafe)
	assert.Contains(t, result, "█")
	assert.Contains(t, result, "░")
}
