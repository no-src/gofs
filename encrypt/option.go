package encrypt

import (
	"github.com/no-src/gofs/conf"
)

// Option the encryption option
type Option struct {
	Encrypt       bool
	EncryptPath   string
	EncryptSecret []byte
	EncryptSuffix string

	Decrypt       bool
	DecryptSecret []byte
	DecryptPath   string
	DecryptOut    string
	DecryptSuffix string
}

// NewOption create an encryption option
func NewOption(config conf.Config) Option {
	return Option{
		Encrypt:       config.Encrypt,
		EncryptPath:   config.EncryptPath,
		EncryptSecret: []byte(config.EncryptSecret),
		EncryptSuffix: config.EncryptSuffix,
		Decrypt:       config.Decrypt,
		DecryptSecret: []byte(config.DecryptSecret),
		DecryptPath:   config.DecryptPath,
		DecryptOut:    config.DecryptOut,
		DecryptSuffix: config.DecryptSuffix,
	}
}

// EmptyOption returns an empty encryption option
func EmptyOption() Option {
	return Option{}
}
