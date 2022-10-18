package encrypt

import (
	"github.com/no-src/gofs/conf"
)

// Option the encryption option
type Option struct {
	Encrypt       bool
	EncryptPath   string
	EncryptSecret []byte

	Decrypt       bool
	DecryptPath   string
	DecryptSecret []byte
	DecryptOut    string
}

// NewOption create an encryption option
func NewOption(config conf.Config) Option {
	if !config.Encrypt && !config.Decrypt {
		return EmptyOption()
	}
	return Option{
		Encrypt:       config.Encrypt,
		EncryptPath:   config.EncryptPath,
		EncryptSecret: []byte(config.EncryptSecret),
		Decrypt:       config.Decrypt,
		DecryptPath:   config.DecryptPath,
		DecryptSecret: []byte(config.DecryptSecret),
		DecryptOut:    config.DecryptOut,
	}
}

// EmptyOption returns an empty encryption option
func EmptyOption() Option {
	return Option{}
}
