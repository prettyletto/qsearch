# QSearch

QSearch is a Linux-first keyboard launcher for search providers.

Type fast, pick or refine a suggestion, press Enter, and QSearch opens the final search in your browser.

```text
qs g
qs yt
qs ytmusic
qs r golang bubbletea
```

The binary name is:

```bash
qs
```

## Install

From source while developing:

```bash
go build -o qs ./cmd/qs/main.go
```

Install with Go:

```bash
go install github.com/prettyletto/qsearch/cmd/qs@latest
```

Make sure your Go bin directory is in `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Config

QSearch creates a provider config on first run:

```text
~/.config/qsearch/providers.toml
```

If `XDG_CONFIG_HOME` is set, it uses:

```text
$XDG_CONFIG_HOME/qsearch/providers.toml
```

Create or restore the default provider config:

```bash
qs init
```

Overwrite it with defaults:

```bash
qs init --force
```

Custom providers are search-only for now. They use `{{query}}` as the escaped query placeholder:

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

Built-in providers currently include Google, YouTube, and YouTube Music with suggestions.

## Hyprland

QSearch is designed to be launched inside a small floating terminal.

The recommended shape is:

```text
keybind -> terminal with qsearch class/app-id -> qs g -> floating centered window
```

Hyprland binds use this form:

```ini
bind = MODS, key, exec, command
```

### Kitty

```ini
bind = SUPER, S, exec, kitty --class qsearch -e qs g

windowrulev2 = float, class:^(qsearch)$
windowrulev2 = center, class:^(qsearch)$
windowrulev2 = size 760 420, class:^(qsearch)$
windowrulev2 = stayfocused, class:^(qsearch)$
```

### Foot

```ini
bind = SUPER, S, exec, foot --app-id qsearch qs g

windowrulev2 = float, app-id:^(qsearch)$
windowrulev2 = center, app-id:^(qsearch)$
windowrulev2 = size 760 420, app-id:^(qsearch)$
windowrulev2 = stayfocused, app-id:^(qsearch)$
```

### Alacritty

```ini
bind = SUPER, S, exec, alacritty --class qsearch -e qs g

windowrulev2 = float, class:^(qsearch)$
windowrulev2 = center, class:^(qsearch)$
windowrulev2 = size 760 420, class:^(qsearch)$
windowrulev2 = stayfocused, class:^(qsearch)$
```

### Ghostty

Ghostty GTK builds support `--class`, which sets the Wayland app-id/class. The class must be a valid GTK application ID, so use a dotted name:

```ini
bind = SUPER, S, exec, ghostty --class=com.prettyletto.qsearch -e qs g

windowrulev2 = float, class:^(com.prettyletto.qsearch)$
windowrulev2 = center, class:^(com.prettyletto.qsearch)$
windowrulev2 = size 760 420, class:^(com.prettyletto.qsearch)$
windowrulev2 = stayfocused, class:^(com.prettyletto.qsearch)$
```

If your Ghostty build does not support `--class`, use your terminal’s title/class support or switch this binding to Kitty, Foot, Alacritty, or WezTerm.

### WezTerm

```ini
bind = SUPER, S, exec, wezterm start --class qsearch -- qs g

windowrulev2 = float, class:^(qsearch)$
windowrulev2 = center, class:^(qsearch)$
windowrulev2 = size 760 420, class:^(qsearch)$
windowrulev2 = stayfocused, class:^(qsearch)$
```

## Omarchy

Omarchy is configured through user dotfiles in `~/.config`. Prefer editing your user Hyprland config rather than files under `~/.local/share/omarchy`.

The best Omarchy integration is to use Omarchy's TUI launcher:

```ini
bindd = SUPER SHIFT, SLASH, QSearch, exec, omarchy-launch-tui qs g
```

`omarchy-launch-tui` opens QSearch in Omarchy's configured terminal and sets the app id/class from the command name:

```text
qs -> org.omarchy.qs
```

Add this window rule so QSearch uses Omarchy's existing floating TUI behavior, like `btop`, `impala`, and `wiremix`:

```ini
windowrule = tag +floating-window, match:class org.omarchy.qs
```

Omarchy already defines the shared `floating-window` behavior:

```ini
windowrule = float on, match:tag floating-window
windowrule = center on, match:tag floating-window
windowrule = size 875 600, match:tag floating-window
```

If you are testing a local repo build, this also works because Omarchy uses `basename` for the app id:

```ini
bindd = SUPER SHIFT, SLASH, QSearch, exec, omarchy-launch-tui /home/you/path/to/qsearch/qs g
```

The class is still:

```text
org.omarchy.qs
```

Recommended workflow:

1. Open the Omarchy menu with `Super + Alt + Space`.
2. Go to Setup / Configs / Hyprland.
3. Add the `omarchy-launch-tui` bind and the `org.omarchy.qs` window rule.
4. Reload Hyprland or let Omarchy restart the relevant process after saving.

If you do not want to use Omarchy's wrapper, use one of the plain Hyprland terminal snippets above instead. That path requires terminal-specific class/app-id flags such as `kitty --class`, `foot --app-id`, or `alacritty --class`.

If you want a tiny wrapper script, put this somewhere in `PATH`, for example `~/.local/bin/qsearch-prompt`:

```bash
#!/usr/bin/env bash
set -euo pipefail

exec omarchy-launch-tui qs g
```

Then bind the script:

```ini
bindd = SUPER SHIFT, SLASH, QSearch, exec, qsearch-prompt
```

Make it executable:

```bash
chmod +x ~/.local/bin/qsearch-prompt
```

## Theme

QSearch is meant to feel native to your terminal theme.

Neutral UI colors use the terminal palette where possible, with small light-theme contrast fixes. Provider tag colors and icons come from built-in provider metadata or `providers.toml`.

For the best visual result, use a Nerd Font-compatible terminal font. If an icon looks wrong, edit or remove the provider `icon` in:

```text
~/.config/qsearch/providers.toml
```

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build -o qs ./cmd/qs/main.go
```

Run:

```bash
./qs g
```

## References

- Hyprland bind syntax: https://wiki.hypr.land/Configuring/Binds/
- Hyprland window rules: https://wiki.hypr.land/Configuring/Window-Rules/
- Omarchy dotfiles and config ownership: https://learn.omacom.io/books/2/pages/65
- Omarchy navigation and menu shortcuts: https://learn.omacom.io/2/the-omarchy-manual
- Ghostty class/app-id behavior: https://man.archlinux.org/man/ghostty.1.en
