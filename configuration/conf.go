package configuration

type LogConfiguration struct {
	LogToFile            bool        `json:"logToFile"`
	MinimumSeverityLevel string      `json:"minimumSeverityLevel"`
	LogToConsole         bool        `json:"logToConsole"`
	FileJsonFormat       bool        `json:"fileJsonFormat"`
	ConsoleJsonFormat    bool        `json:"consoleJsonFormat"`
	CompressLogs         interface{} `json:"compressLogs"`
	DisableConsoleColor  bool        `json:"disableConsoleColor,omitempty"`
}
