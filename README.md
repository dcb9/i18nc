# i18nc
Generate type safe Go code from english locale content

## Usage

### Create i18nc generate.go

```
# file i18nc/generate.go

package i18nc

//go:generate go run -mod=mod github.com/dcb9/i18nc/cmd/i18nc $GOPACKAGE i18nc_generated.go ../locales/en.json
```

### generate Go code

```bash
$ go generate ./i18nc
```

### Init i18nc and use

```
package main

import (
	"your_module/i18nc"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func main() {
	i18nc.Localizer = initLocalizer()

	fmt.Println(i18nc.YourMessageID())
}

func initLocalizer() *i18n.Localizer {
	// FIXME: do your own init procedure
	return nil
}
```

## Acknowledgements

- [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n)
- [gobeam/stringy](https://github.com/gobeam/stringy)

## License

Released under the MIT License - see `LICENSE` for details.
