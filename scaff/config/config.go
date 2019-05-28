package config

type TemplateKey string

func (x TemplateKey) AsRaw() string {
	return string(x)
}

type Description string

type ResolvedConfig map[TemplateKey]ParsedValue
type UnresolvedConfig map[TemplateKey]TemplateValue

type HookConfig struct {
	Command string
	Args    []string
}

type TemplateValue struct {
	Default     string
	Description Description
	VerifyHook  HookConfig `yaml:"verify_hook"`
}

type ParsedValue struct {
	Source      TemplateValue
	ParsedValue string
}

type FileConfig struct {
	FileDelims    map[string]string `yaml:"lang_delims"`
	LanguageRules LanguageRules     `yaml:"lang_rules"`
}

type GoRules struct {
	SourcePackage string `yaml:"pkg"`
	ReplaceRule   string `yalm:"replace_with_id"`
}

type LanguageRules struct {
	Go *GoRules `yaml:"go"`
}

type ScaffConfig struct {
	Context    UnresolvedConfig
	Dynamics   map[string]string
	FileConfig FileConfig `yaml:"file_config"`
}

func (r ResolvedConfig) AsRaw() map[string]string {
	raw := make(map[string]string)
	for k, v := range r {
		raw[k.AsRaw()] = v.ParsedValue
	}

	return raw
}
