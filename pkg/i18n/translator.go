package i18n

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)

const DefaultLanguageCode = "ru"

var (
	langs = map[string]*Language{}

	ErrKeyNotFound = errors.New("key not found")
	ErrCasting     = errors.New("casting error")
)

type TranslationBuilder struct {
	f    string
	args []interface{}
}

func NewBuilder(f string) *TranslationBuilder {
	return &TranslationBuilder{f: f, args: []interface{}{}}
}

func (t *TranslationBuilder) Arg(args ...interface{}) *TranslationBuilder {
	t.args = append(t.args, args...)

	return t
}

func (t *TranslationBuilder) Translate(lang string) string {
	return Translatef(t.f, lang, t.args...)
}

func LoadLangs(path string) error {
	files, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".toml") {
			var lang Language

			data, err := os.ReadFile(filepath.Join(path, file.Name()))

			if err != nil {
				return err
			}

			if _, err := toml.Decode(string(data), &lang); err != nil {
				return err
			}

			langIds := []string{strings.TrimSuffix(file.Name(), ".toml")}

			langIds = append(langIds, lang.Aliases...)

			for _, id := range langIds {
				langs[id] = &lang
			}
		}
	}

	return nil
}

func Translatef(key string, language string, args ...interface{}) string {
	lang, ok := langs[language]

	if !ok {
		lang = langs[DefaultLanguageCode]
	}

	if m, ok := lang.Translations[key]; ok {
		if s, ok := m.(string); ok {
			return fmt.Sprintf(s, args...)
		}
	}

	return key
}

func Get[T any](key string, language string) (val T, err error) {
	lang, ok := langs[language]

	if !ok {
		lang = langs[DefaultLanguageCode]
	}

	v, ok := lang.Translations[key]

	if !ok {
		return val, ErrKeyNotFound
	}

	if value, ok := v.(T); ok {
		return value, nil
	}

	return val, ErrCasting
}
