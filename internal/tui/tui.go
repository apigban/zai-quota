package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"image/color"

	"zai-quota/internal/api"
	"zai-quota/internal/config"
	"zai-quota/internal/models"
	"zai-quota/internal/processor"
)

func Run(cfg *config.Config, debug bool) error {
	m := NewModel(cfg, debug)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

type state int

const (
	stateLoading state = iota
	stateLoaded
	stateRefreshing
	stateError
	stateSetup
	stateSetupMergeChoice
	stateSetupClosing
)

type Model struct {
	state              state
	cfg                *config.Config
	client             *api.Client
	limits             []models.Limit
	level              string
	processed          map[string]ProcessedLimitData
	expanded           map[string]bool
	err                error
	width              int
	height             int
	debug              bool
	debugLog           []string
	setupInput         textinput.Model
	existingConfig     bool
	saveError          error
	setupClosingReason string
}

type ProcessedLimitData struct {
	Type          string
	Label         string
	Percentage    int
	Used          int
	Total         int
	Remaining     int
	ResetDisplay  string
	ResetDateTime string
	WarningLevel  string
	UsageDetails  []models.UsageDetail
}

type quotaFetchedMsg struct {
	quota     *models.QuotaResponse
	err       error
	debugInfo []string
}

func NewModel(cfg *config.Config, debug bool) Model {
	client := api.NewClient(cfg.APIKey, cfg.Endpoint, cfg.TimeoutSeconds)

	ti := textinput.New()
	ti.Placeholder = "Enter your API key"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()

	existingConfig, _ := config.ConfigFileExists()

	initialState := stateLoading
	if cfg.APIKey == "" {
		initialState = stateSetup
	}

	return Model{
		state:          initialState,
		cfg:            cfg,
		client:         client,
		limits:         []models.Limit{},
		level:          "",
		processed:      make(map[string]ProcessedLimitData),
		expanded:       make(map[string]bool),
		debug:          debug,
		debugLog:       []string{},
		setupInput:     ti,
		existingConfig: existingConfig,
	}
}

func (m Model) Init() tea.Cmd {
	if m.state == stateSetup {
		return nil
	}
	return fetchQuotaCmd(m.client, m.cfg.Endpoint, m.debug)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.state == stateSetup {
			return m.handleSetupKeyPress(msg)
		}
		if m.state == stateSetupMergeChoice {
			return m.handleSetupMergeChoiceKeyPress(msg)
		}
		if m.state == stateSetupClosing {
			return m.handleSetupClosingKeyPress(msg)
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case quotaFetchedMsg:
		return m.handleQuotaFetched(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "Q", "ctrl+c":
		return m, tea.Quit

	case "r", "R":
		if m.state != stateRefreshing {
			m.state = stateRefreshing
			return m, fetchQuotaCmd(m.client, m.cfg.Endpoint, m.debug)
		}
	case "e", "E":
		if m.expanded == nil {
			m.expanded = make(map[string]bool)
		}
		if data, exists := m.processed["TOKENS_LIMIT"]; exists && len(data.UsageDetails) > 0 {
			m.expanded["TOKENS_LIMIT"] = !m.expanded["TOKENS_LIMIT"]
		}
		if data, exists := m.processed["TIME_LIMIT"]; exists && len(data.UsageDetails) > 0 {
			m.expanded["TIME_LIMIT"] = !m.expanded["TIME_LIMIT"]
		}
	}

	return m, nil
}

func (m Model) handleSetupKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.setupInput, cmd = m.setupInput.Update(msg)

	switch msg.String() {
	case "enter":
		if m.setupInput.Value() == "" {
			return m, cmd
		}
		if m.existingConfig {
			m.state = stateSetupMergeChoice
			return m, nil
		}
		return m.saveConfigAndProceed()

	case "esc", "ctrl+c":
		m.state = stateSetupClosing
		m.setupClosingReason = "API Key is required to use this tool. Closing..."
		return m, nil
	}

	return m, cmd
}

func (m Model) handleSetupMergeChoiceKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a", "A":
		return m.mergeAndSave()
	case "o", "O":
		return m.saveConfigAndProceed()
	case "esc", "ctrl+c":
		m.state = stateSetupClosing
		m.setupClosingReason = "API Key is required to use this tool. Closing..."
		return m, nil
	}
	return m, nil
}

func (m Model) handleSetupClosingKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "Q", "enter", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) mergeAndSave() (tea.Model, tea.Cmd) {
	existingCfg, err := config.LoadConfigFromFile()
	if err != nil {
		m.saveError = err
		m.state = stateSetupClosing
		m.setupClosingReason = fmt.Sprintf("Could not read existing config: %s. Please restart the TUI to retry.", err)
		return m, nil
	}

	existingCfg.APIKey = m.setupInput.Value()
	m.cfg.APIKey = existingCfg.APIKey
	m.cfg.Endpoint = existingCfg.Endpoint
	m.cfg.TimeoutSeconds = existingCfg.TimeoutSeconds

	if err := config.SaveConfig(m.cfg); err != nil {
		m.saveError = err
		m.state = stateSetupClosing
		m.setupClosingReason = fmt.Sprintf("Could not save configuration: %s. Please restart the TUI to retry.", err)
		return m, nil
	}

	m.client = api.NewClient(m.cfg.APIKey, m.cfg.Endpoint, m.cfg.TimeoutSeconds)
	m.state = stateLoading
	return m, fetchQuotaCmd(m.client, m.cfg.Endpoint, m.debug)
}

func (m Model) saveConfigAndProceed() (tea.Model, tea.Cmd) {
	m.cfg.APIKey = m.setupInput.Value()

	if err := config.SaveConfig(m.cfg); err != nil {
		m.saveError = err
		m.state = stateSetupClosing
		m.setupClosingReason = fmt.Sprintf("Could not save configuration: %s. Please restart the TUI to retry.", err)
		return m, nil
	}

	m.client = api.NewClient(m.cfg.APIKey, m.cfg.Endpoint, m.cfg.TimeoutSeconds)
	m.state = stateLoading
	return m, fetchQuotaCmd(m.client, m.cfg.Endpoint, m.debug)
}

func (m Model) handleQuotaFetched(msg quotaFetchedMsg) (tea.Model, tea.Cmd) {
	if len(msg.debugInfo) > 0 {
		m.debugLog = append(m.debugLog, msg.debugInfo...)
	}

	if msg.err != nil {
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	m.limits = msg.quota.Limits
	m.level = msg.quota.Level
	m.processed = processLimits(msg.quota.Limits)
	m.err = nil
	m.state = stateLoaded

	return m, nil
}

func (m Model) View() tea.View {
	var b strings.Builder

	switch m.state {
	case stateSetup:
		b.WriteString(m.renderSetupView())
	case stateSetupMergeChoice:
		b.WriteString(m.renderMergeChoiceView())
	case stateSetupClosing:
		b.WriteString(m.renderSetupClosingView())
	default:
		b.WriteString(m.renderTitle())
		b.WriteString("\n")

		if m.state == stateError && m.processed == nil {
			b.WriteString(m.renderErrorOverlay())
			b.WriteString("\n")
			b.WriteString(emptyStateStyle.Render("No quota data available"))
			b.WriteString("\n")
		} else {
			if m.state == stateError {
				b.WriteString(m.renderErrorOverlay())
				b.WriteString("\n")
			}

			if len(m.processed) > 0 {
				b.WriteString(m.renderQuotaDisplay())
				b.WriteString("\n")
			}
		}

		b.WriteString(m.renderHelp())

		if m.debug && len(m.debugLog) > 0 {
			b.WriteString("\n")
			b.WriteString(m.renderDebugLog())
		}
	}

	return tea.NewView(b.String())
}

func (m Model) renderSetupView() string {
	var b strings.Builder

	b.WriteString(setupTitleStyle.Render("Welcome to Z.AI Quota Monitor!"))
	b.WriteString("\n\n")
	b.WriteString(setupPromptStyle.Render("To get started, please provide your API Key."))
	b.WriteString("\n")
	b.WriteString(setupPromptStyle.Render("You can find this in your z.ai dashboard."))
	b.WriteString("\n\n")
	b.WriteString(setupPromptStyle.Render("API Key: "))
	b.WriteString(m.setupInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("(Enter to confirm, Esc to quit)"))

	return b.String()
}

func (m Model) renderMergeChoiceView() string {
	var b strings.Builder

	b.WriteString(setupTitleStyle.Render("Existing Config Found"))
	b.WriteString("\n\n")
	b.WriteString(setupPromptStyle.Render("A configuration file already exists."))
	b.WriteString("\n")
	b.WriteString(setupPromptStyle.Render("How would you like to proceed?"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("[A] Append key to existing config  [O] Overwrite all  [Esc] Cancel"))

	return b.String()
}

func (m Model) renderSetupClosingView() string {
	var b strings.Builder

	b.WriteString(setupTitleStyle.Render("Setup"))
	b.WriteString("\n\n")

	if m.saveError != nil {
		b.WriteString(errorTextStyle.Render("Error: "))
		b.WriteString(m.saveError.Error())
		b.WriteString("\n\n")
	}

	b.WriteString(errorTextStyle.Render(m.setupClosingReason))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("[q] Quit"))

	return b.String()
}

func (m Model) renderTitle() string {
	var status string
	switch m.state {
	case stateLoading:
		status = lipgloss.NewStyle().Foreground(colorPurple).Render("Loading...")
	case stateRefreshing:
		status = lipgloss.NewStyle().Foreground(colorPurple).Render("Refreshing...")
	case stateLoaded, stateError:
		status = time.Now().Format("15:04")
	}

	title := titleStyle.Render("Z.AI Quota Monitor")

	if m.level != "" {
		level := strings.Title(strings.ToLower(m.level))
		levelText := levelStyle.Render(fmt.Sprintf(" [%s]", level))
		title = lipgloss.JoinHorizontal(lipgloss.Top, title, levelText, statusStyle.Render(status))
	} else {
		title = lipgloss.JoinHorizontal(lipgloss.Top, title, statusStyle.Render(status))
	}

	return title
}

func processLimits(limits []models.Limit) map[string]ProcessedLimitData {
	result := make(map[string]ProcessedLimitData)

	for _, limit := range limits {
		resetTime := processor.ConvertTimestamp(limit.NextResetTime)
		data := ProcessedLimitData{
			Type:          limit.Type,
			Percentage:    limit.Percentage,
			ResetDisplay:  processor.FormatTimeUntil(resetTime),
			ResetDateTime: processor.FormatResetDateTime(resetTime),
			UsageDetails:  limit.UsageDetails,
			Total:         limit.Usage,
			Used:          limit.CurrentValue,
			Remaining:     limit.Remaining,
		}

		switch limit.Type {
		case "TOKENS_LIMIT":
			data.Label = "[5-Hour Prompt Limit]"
		case "TIME_LIMIT":
			data.Label = "[Tool Quota]"
		}

		data.WarningLevel = getWarningLevel(data.Percentage)
		result[limit.Type] = data
	}

	return result
}

func getWarningLevel(percentage int) string {
	switch {
	case percentage >= 95:
		return "emergency"
	case percentage >= 90:
		return "critical"
	case percentage >= 80:
		return "warning"
	default:
		return "safe"
	}
}

func (m Model) renderQuotaDisplay() string {
	var b strings.Builder

	if data, exists := m.processed["TOKENS_LIMIT"]; exists {
		b.WriteString(m.renderLimit(data))
		b.WriteString("\n")
	}

	if data, exists := m.processed["TIME_LIMIT"]; exists {
		b.WriteString(m.renderLimit(data))
		b.WriteString("\n")
	}

	return b.String()
}

func fetchQuotaCmd(client *api.Client, endpoint string, debug bool) tea.Cmd {
	return func() tea.Msg {
		var debugInfo []string
		if debug {
			debugInfo = append(debugInfo, fmt.Sprintf("Fetching quota from %s", endpoint))
		}

		quota, err := client.FetchQuota(context.Background())
		if debug && err != nil {
			debugInfo = append(debugInfo, fmt.Sprintf("Fetch error: %v", err))
		}

		return quotaFetchedMsg{
			quota:     quota,
			err:       err,
			debugInfo: debugInfo,
		}
	}
}

func (m Model) renderLimit(data ProcessedLimitData) string {
	var b strings.Builder

	b.WriteString(quotaLabelStyle.Render(data.Label))
	b.WriteString("\n")

	progressColor := getColorForPercentage(data.Percentage)
	progressStyle := lipgloss.NewStyle().Foreground(progressColor)
	progressBar := m.renderProgressBar(data.Percentage, progressColor)
	b.WriteString(progressStyle.Render(progressBar))
	b.WriteString("\n")

	if data.Total > 0 {
		b.WriteString(quotaValueStyle.Render(fmt.Sprintf("%d / %d", data.Used, data.Total)))
		b.WriteString("\n")
	} else {
		b.WriteString(quotaValueStyle.Render(fmt.Sprintf("%d%% used", data.Percentage)))
		b.WriteString("\n")
	}

	if data.Remaining > 0 {
		b.WriteString(quotaValueStyle.Render(fmt.Sprintf("Remaining: %d", data.Remaining)))
		b.WriteString("\n")
	}

	b.WriteString(quotaValueStyle.Render(fmt.Sprintf("Reset: %s (%s)", data.ResetDisplay, data.ResetDateTime)))
	b.WriteString("\n")

	if len(data.UsageDetails) > 0 {
		if m.expanded[data.Type] {
			b.WriteString(quotaValueStyle.Render("Usage breakdown:"))
			b.WriteString("\n")
			for _, detail := range data.UsageDetails {
				b.WriteString(quotaValueStyle.Render(fmt.Sprintf("  └─ %s: %d", detail.ModelCode, detail.Usage)))
				b.WriteString("\n")
			}
			b.WriteString(helpStyle.Render("[e] Collapse"))
		} else {
			b.WriteString(helpStyle.Render("[e] Expand details"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderProgressBar(percentage int, colorVal color.Color) string {
	width := 40
	filled := int(float64(width) * float64(percentage) / 100)

	barStyle := lipgloss.NewStyle().Foreground(colorVal)
	emptyStyle := lipgloss.NewStyle().Foreground(colorGray)

	return barStyle.Render(strings.Repeat("█", filled)) + emptyStyle.Render(strings.Repeat("░", width-filled))
}

func (m Model) renderHelp() string {
	hasDetails := false
	if data, exists := m.processed["TOKENS_LIMIT"]; exists && len(data.UsageDetails) > 0 {
		hasDetails = true
	}
	if data, exists := m.processed["TIME_LIMIT"]; exists && len(data.UsageDetails) > 0 {
		hasDetails = true
	}

	if hasDetails {
		return helpStyle.Render("[r] Refresh  [e] Expand  [q] Quit")
	}
	return helpStyle.Render("[r] Refresh  [q] Quit")
}

func (m Model) renderErrorOverlay() string {
	if m.err == nil {
		return ""
	}

	return errorBoxStyle.Render(
		errorTextStyle.Render("Error: ") + m.err.Error(),
	)
}

func (m Model) renderDebugLog() string {
	var b strings.Builder
	b.WriteString(helpStyle.Render("Debug Log:"))
	b.WriteString("\n")
	for _, log := range m.debugLog {
		b.WriteString(helpStyle.Render(log))
		b.WriteString("\n")
	}
	return b.String()
}
