# ENV - Lightweight Environment Variables Parser

## Author's POV
Parses environmental variables to a struct using the struct field's `env` tag.

## Usage

```go

func main() {
  config := struct {
      // note the tag we are using
      // env determines what ENV var the *Parse func will be reading
      // envDefault sets the default value if ENV var is blank or not set
      Foo string `env:"FOO" envDefault:"BooHoo"`
  }

  os.Setenv("FOO", "Hoho")
  env.Parse(&config)
}

```

## Supported Data Types
When creating the struct with it's fields the only supported types for now
are:

1. `string`
2. `int`
3. `bool`
