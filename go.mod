module github.com/scottmckendry/beam

go 1.24.4

require (
	github.com/Oudwins/tailwind-merge-go v0.2.1
	github.com/a-h/templ v0.3.906
	github.com/dustin/go-humanize v1.0.1
	github.com/go-chi/chi/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/securecookie v1.1.2
	github.com/joho/godotenv v1.5.1
	github.com/starfederation/datastar-go v1.0.1
	github.com/tursodatabase/go-libsql v0.0.0-20250609073118-9c24e0e7fa97
	github.com/yuin/goldmark v1.7.12
	golang.org/x/oauth2 v0.30.0
)

// pin the release canditate version of datastar (prevents falling back to the beta versions)
replace github.com/starfederation/datastar => github.com/starfederation/datastar v1.0.0-RC.1

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/CAFxX/httpcompression v0.0.9 // indirect
	github.com/a-h/parse v0.0.0-20250122154542-74294addb73e // indirect
	github.com/air-verse/air v1.62.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/axzilla/templui v0.82.0 // indirect
	github.com/bep/godartsass/v2 v2.5.0 // indirect
	github.com/bep/golibsass v1.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/cli/browser v1.3.0 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gohugoio/hugo v0.147.6 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/libsql/sqlite-antlr4-parser v0.0.0-20240327125255-dbf53b6cbf06 // indirect
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tdewolff/minify v2.3.6+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/parse/v2 v2.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

tool (
	github.com/a-h/templ/cmd/templ
	github.com/air-verse/air
	github.com/axzilla/templui/cmd/templui
	github.com/tdewolff/minify/cmd/minify
)
