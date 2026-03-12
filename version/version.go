package version

// Version 构建时通过 ldflags 注入，例如：-X smartwatch-server/version.Version=v1.0.0-abc1234
var Version = "unknown"
