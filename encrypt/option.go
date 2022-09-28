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
	DecryptPath   string
	DecryptSecret []byte
	DecryptSuffix string
	DecryptOut    string
}

// NewOption create an encryption option
func NewOption(config conf.Config) Option {
	return Option{
		Encrypt:       config.Encrypt,
		EncryptPath:   config.EncryptPath,
		EncryptSecret: []byte(config.EncryptSecret),
		EncryptSuffix: config.EncryptSuffix,
		Decrypt:       config.Decrypt,
		DecryptPath:   config.DecryptPath,
		DecryptSecret: []byte(config.DecryptSecret),
		DecryptSuffix: config.DecryptSuffix,
		DecryptOut:    config.DecryptOut,
	}
}

// EmptyOption returns an empty encryption option
func EmptyOption() Option {
	return Option{}
}
