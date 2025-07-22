module github.com/aionbit/ixo-plugins/net/proxy

go 1.24.4

require (
	github.com/aionbit/ixo-plugins/plugin v0.0.0-20250715122039-2f9dc1e563a1
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/mitchellh/mapstructure v1.5.0 // indirect

replace github.com/aionbit/ixo-plugins/plugin => ../../plugin