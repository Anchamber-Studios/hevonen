package translation

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Languages map[string]Messages

type Messages map[string]string

//go:embed *.toml
var LocaleFS embed.FS

var bundle *i18n.Bundle

func SetupTranslations() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "active.en.toml")
}

func GetBundle() *i18n.Bundle {
	if bundle == nil {
		SetupTranslations()
	}
	return bundle
}
