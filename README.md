# QSearch

QSearch is a Linux-first search launcher for your terminal.

It is not trying to be a browser, a search-results page, or a terminal clone of Google or YouTube. The goal is simpler:

```text
type fast -> get suggestions -> choose or keep typing -> press enter -> open the browser
```

The binary is named:

```sh
qs
```

Typical use:

```sh
qs g
qs yt
qs ytmusic
qs r golang bubbletea
```

When you pass a query, QSearch opens the search immediately. When you do not pass a query, it opens the TUI.

```sh
qs g golang channels
# opens Google immediately

qs g
# opens the TUI with Google selected
```

## Features

- Fast CLI-first search launcher
- TUI prompt for interactive searches
- Google, YouTube, and YouTube Music providers
- Suggestions/autocomplete for built-in providers
- Custom search-only providers from `providers.toml`
- Browser opening through `xdg-open`
- Provider switching inside the TUI
- Nerd Font-friendly provider icons and keycaps

## Requirements

You need:

- Linux
- Go
- `xdg-open`
- a terminal with a Nerd Font configured

For the best experience, use QSearch from a small floating terminal window. It works fine in a normal terminal too.

## Install

From this repo:

```sh
make install
```

That installs `qs` into:

```sh
$(go env GOPATH)/bin
```

Make sure that directory is in your shell `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

For a permanent shell setup, put that line in your shell config, for example `~/.zshrc`, `~/.bashrc`, or your fish equivalent.

You can also build a local binary without installing:

```sh
make build
./qs g
```

Or install directly with Go:

```sh
go install github.com/prettyletto/qsearch/cmd/qs@latest
```

## Usage

Show help:

```sh
qs help
```

Open the TUI with Google:

```sh
qs g
```

Search immediately:

```sh
qs g linux clipboard manager
qs yt bubble tea tui
qs ytmusic radiohead
```

Built-in providers:

```text
google, g
youtube, y, yt
ytmusic, ym, music
```

Custom providers are loaded from your config file and can add names like `r` for Reddit or `c` for ChatGPT.

## TUI Keys

Inside the TUI:

```text
tab      switch provider
up/down  select suggestion
enter    open selected suggestion or typed query
esc      exit
```

The TUI expects a Nerd Font. If symbols look wrong, fix the terminal font first.

## Config

QSearch creates this file on first run:

```text
~/.config/qsearch/providers.toml
```

If `XDG_CONFIG_HOME` is set, it uses:

```text
$XDG_CONFIG_HOME/qsearch/providers.toml
```

Create the default config:

```sh
qs init
```

Overwrite it with the defaults:

```sh
qs init --force
```

The default custom providers look like this:

```toml
[[providers]]
name = "reddit"
aliases = ["r"]
url = "https://www.reddit.com/search/?q={{query}}"
icon = ""
tag_bg = "#FF4500"
icon_color = "#FFFFFF"
text_color = "#FFFFFF"

[[providers]]
name = "chatgpt"
aliases = ["c", "gpt"]
url = "https://chatgpt.com/?q={{query}}"
icon = "⌬"
tag_bg = "#10A37F"
icon_color = "#FFFFFF"
text_color = "#FFFFFF"
```

Custom providers are search-only for now. Use `{{query}}` where the escaped search text should go.

## Omarchy

QSearch works nicely with Omarchy because Omarchy already has a TUI launcher:

```sh
omarchy-launch-tui
```

The basic keybind in `~/.config/hypr/bindings.conf` is:

```ini
bindd = SUPER SHIFT, SLASH, QSearch, exec, omarchy-launch-tui qs g
```

However, there is one important detail: Hyprland/Omarchy may not see your shell `PATH`. If you installed QSearch with Go and the keybind does nothing, use the full path to `qs`.

Find it with:

```sh
go env GOPATH
```

Usually the binary will be here:

```text
/home/your-user/go/bin/qs
```

Then bind that full path:

```ini
bindd = SUPER SHIFT, SLASH, QSearch, exec, omarchy-launch-tui /home/your-user/go/bin/qs g
```

The Omarchy launcher sets the app id/class from the binary name:

```text
qs -> org.omarchy.qs
```

To make it floating and centered with your own size, add direct window rules in `~/.config/hypr/bindings.conf`:

```ini
windowrule = float on, match:class org.omarchy.qs
windowrule = center on, match:class org.omarchy.qs
windowrule = size 900 420, match:class org.omarchy.qs
```

Do not tag QSearch as `floating-window` if you want a custom size. Omarchy's default floating-window tag applies its own shared size:

```ini
windowrule = size 875 600, match:tag floating-window
```

After changing Hyprland config:

```sh
omarchy-restart-hyprctl
hyprctl configerrors
```

Then close and reopen QSearch. Window rules apply most reliably when the window is newly created.

## Plain Hyprland

If you are not using Omarchy, launch QSearch in a terminal with a stable class or app id, then write rules for that class.

Kitty:

```ini
bind = SUPER SHIFT, SLASH, exec, kitty --class qsearch -e qs g

windowrule = float on, match:class qsearch
windowrule = center on, match:class qsearch
windowrule = size 900 420, match:class qsearch
```

Foot:

```ini
bind = SUPER SHIFT, SLASH, exec, foot --app-id qsearch qs g

windowrule = float on, match:class qsearch
windowrule = center on, match:class qsearch
windowrule = size 900 420, match:class qsearch
```

Ghostty:

```ini
bind = SUPER SHIFT, SLASH, exec, ghostty --class=com.prettyletto.qsearch -e qs g

windowrule = float on, match:class com.prettyletto.qsearch
windowrule = center on, match:class com.prettyletto.qsearch
windowrule = size 900 420, match:class com.prettyletto.qsearch
```

If your compositor uses older `windowrulev2` syntax, translate the same idea to your local Hyprland version.

## Development

Common commands:

```sh
make help
make test
make build
make run
make install
```

Build output:

```text
./qs
```

Run from source:

```sh
go run ./cmd/qs g
```

Run the installed binary:

```sh
qs g
```

Clean local build output:

```sh
make clean
```

## Project Shape

The project is intentionally small:

```text
cmd/qs/main.go                    app entrypoint
internal/dispatch                 provider argument dispatch
internal/app/search               search flow coordination
internal/domain/provider          provider interfaces
internal/providers                built-in providers
internal/tui/search               Bubble Tea TUI
internal/infra/browser            xdg-open integration
internal/config                   providers.toml loading
```

QSearch should stay focused: choose a provider, get suggestions, open the final browser URL.
