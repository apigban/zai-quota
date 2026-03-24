package models

type QuotaResponse struct {
	Limits []Limit `json:"limits"`
	Level  string  `json:"level"`
}

type Limit struct {
	Type          string        `json:"type"`
	Unit          int           `json:"unit"`
	Number        int           `json:"number"`
	Usage         int           `json:"usage,omitempty"`
	CurrentValue  int           `json:"currentValue,omitempty"`
	Remaining     int           `json:"remaining,omitempty"`
	Percentage    int           `json:"percentage"`
	NextResetTime int64         `json:"nextResetTime"`
	UsageDetails  []UsageDetail `json:"usageDetails,omitempty"`
}

type UsageDetail struct {
	ModelCode string `json:"modelCode"`
	Usage     int    `json:"usage"`
}
