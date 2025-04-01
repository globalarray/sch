package i18n

type Language struct {
	Name         string                 `toml:"Name"`
	Aliases      []string               `toml:"Aliases"`
	Translations map[string]interface{} `toml:"Translations"`
}
