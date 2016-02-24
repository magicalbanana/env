package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type config struct {
	FooStr string `env:"FOOSTR"`
	Bool   bool   `env:"BOOL"`
	Port   int    `env:"PORT"`
	DBUrl  string `env:"DBURL" envDefault:"postgres://localhost:5432/mbp"`
}

func TestParse(t *testing.T) {
	tests := []struct {
		desc      string
		envVar    map[string]string
		assertion func(m string)
	}{
		{
			desc:      "passes a non pointer",
			envVar:    nil,
			assertion: func(m string) { assert.Error(t, Parse(config{}), m) },
		},
		{
			desc:   "passes a non struct",
			envVar: nil,
			assertion: func(m string) {
				c := config{}
				assert.Error(t, Parse(&c.Port), m)
			},
		},
		{
			desc:   "unsupported type",
			envVar: map[string]string{"RUNE": "RUNE"},
			assertion: func(m string) {
				c := struct {
					Rune rune `env:"RUNE"`
				}{}
				assert.Error(t, Parse(&c), m)
			},
		},
		{
			desc: "incorrect type is set to an int type",
			envVar: map[string]string{
				"PORT": "manbearpig",
			},
			assertion: func(m string) { assert.Error(t, Parse(&config{}), m) },
		},
		{
			desc: "incorrect type is set to a bool type",
			envVar: map[string]string{
				"BOOL": "manbearpig",
			},
			assertion: func(m string) { assert.Error(t, Parse(&config{}), m) },
		},
		{
			desc: "parses ENV vars into the config struct successfully",
			envVar: map[string]string{
				"FOOSTR": "FOO",
				"BOOL":   "True",
				"PORT":   "3000",
			},
			assertion: func(m string) {
				err := Parse(&config{})
				assert.NoError(t, err, m)
			},
		},
		{
			desc: "sets default value when ENV var is not set",
			assertion: func(m string) {
				c := config{}
				Parse(&c)
				assert.Equal(
					t,
					c.DBUrl,
					"postgres://localhost:5432/mbp",
					m,
				)
			},
		},
	}

	for _, test := range tests {
		if test.envVar != nil {
			for k, v := range test.envVar {
				os.Setenv(k, v)
			}
		}

		defer func() {
			for k, _ := range test.envVar {
				os.Setenv(k, "")

			}
		}()

		test.assertion(test.desc)
	}
}
