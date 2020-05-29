package gogen

type Gen interface {
	Generate(cfg *Config) (File, error)
}

func Generate(file Gen, opts ...Option) error {

	cfg := DefaultConfig

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if err := generate(file, &cfg); err != nil {
		return err
	}

	return nil
}

func isFiles(doc File) ([]File, bool) {
	if docs, ok := doc.(Files); ok {
		return docs, true
	}
	return nil, false
}

func generate(gen Gen, cfg *Config) error {
	file, err := gen.Generate(cfg)
	if err != nil {
		return err
	}
	if files, ok := isFiles(file); ok {
		for _, file = range files {
			if err := Write(file); err != nil {
				return err
			}
		}
	} else {
		if err := Write(file); err != nil {
			return err
		}
	}

	return nil
}

type Option interface {
	apply(cfg *Config)
}

type optionFunc func(*Config)

func (o optionFunc) apply(s *Config) {
	o(s)
}

func SetDescription(desc string) optionFunc {
	return func(cfg *Config) {
		cfg.Description = desc
	}
}
func SkipChildren(skip bool) optionFunc {
	return func(cfg *Config) {
		cfg.SkipChildren = skip
	}
}
