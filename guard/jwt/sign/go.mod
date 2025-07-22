module github.com/aionbit/ixo-plugins/guard/jwt/sign

go 1.24.4

require (
	github.com/aionbit/ixo-plugins/plugin v0.0.0-20250721025127-33acfc63bf41
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/aionbit/ixo-plugins/guard v0.0.0-00010101000000-000000000000
)

require github.com/mitchellh/mapstructure v1.5.0 // indirect

replace github.com/aionbit/ixo-plugins/guard => ../../
