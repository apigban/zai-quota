package exporter

import (
	"net/http"
)

const landingPageHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Z.ai Quota Exporter</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        a { color: #0066cc; }
        code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>Z.ai Quota Exporter</h1>
    <p>Prometheus exporter for Z.ai API quota metrics.</p>
    <p><a href="/metrics">Metrics</a> | <a href="/health">Health</a></p>
    <h2>Available Metrics</h2>
    <ul>
        <li><code>zai_quota_prompt_usage_ratio</code> - Prompt usage as ratio (0-1)</li>
        <li><code>zai_quota_tool_calls_used</code> - Tool calls used</li>
        <li><code>zai_quota_tool_calls_limit</code> - Tool call limit</li>
        <li><code>zai_quota_tool_calls_by_tool</code> - Per-tool breakdown</li>
        <li><code>zai_quota_up</code> - Exporter health status</li>
    </ul>
</body>
</html>
`

func NewServerHandler(metricsHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", metricsHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", landingHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(landingPageHTML))
}
