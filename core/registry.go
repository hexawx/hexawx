package core

type RemotePlugin struct {
	Name               string   `json:"name"`
	DisplayName        string   `json:"display_name"`
	Version            string   `json:"version"`
	SupportedPlatforms []string `json:"supported_platforms"`
	Type               string   `json:"type"`
	BinaryURL          string   `json:"binary_url"`
	Description        string   `json:"description"`
}
