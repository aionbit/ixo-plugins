module github.com/aionbit/ixo-plugins/guard/signature/verify

go 1.24.4

require (
	github.com/aionbit/ixo-plugins/plugin v0.0.0-20250721025127-33acfc63bf41
	github.com/aionbit/ixo-plugins/guard v0.0.0-00010101000000-000000000000
)

require github.com/mitchellh/mapstructure v1.5.0 // indirect

replace github.com/aionbit/ixo-plugins/guard => ../../
