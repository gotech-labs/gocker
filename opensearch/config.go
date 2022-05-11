package opensearch

type config struct {
	plugins []string
}

type configOption func(*config)

func WithPlugins(plugins ...string) configOption {
	return func(opt *config) {
		opt.plugins = plugins
	}
}
