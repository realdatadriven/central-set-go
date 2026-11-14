# Documentation site (Hugo + Plume)

This folder contains a Hugo site built using the Plume theme.

Theme is added as a Git submodule under `docs/themes/plume`.

Quickstart:

- Ensure `hugo` and Dart Sass are installed (https://gohugo.io/getting-started/quick-start/ and https://gohugo.io/functions/css/sass/#dart-sass).
- From repo root, update submodules: `git submodule update --init --recursive`.
- Start local server from `docs`: `hugo server --source docs --themesDir themes -D --disableFastRender`.

The site content is taken from the theme's exampleSite (https://github.com/Dobefu/hugo-theme-plume/tree/main/exampleSite).

Customize `docs/hugo.toml` or the files in `docs/content` as needed.


`$HOME/go/bin/hugo serve --source docs --themesDir themes -D --disableFastRender`
