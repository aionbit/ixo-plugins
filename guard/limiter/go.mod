module github.com/aionbit/ixo-plugins/guard/limiter

go 1.24.4

require github.com/aionbit/ixo-plugins/plugin v0.0.0-20250721025127-33acfc63bf41

require github.com/mitchellh/mapstructure v1.5.0 // indirect

replace github.com/aionbit/ixo-plugins/guard => ../../

require (
	github.com/bluele/gcache v0.0.2
	golang.org/x/time v0.12.0
)
