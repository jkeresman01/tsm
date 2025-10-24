<div align="center">

  <h1>tsm</h1>
  <h6> TMUX session manager </h4>

[![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
![Status](https://img.shields.io/badge/status-work--in--progress-yellow?style=for-the-badge)


</div>

[![asciicast](https://asciinema.org/a/HuFKZYETW9Xy9AADF5uu32l54.svg)](https://asciinema.org/a/HuFKZYETW9Xy9AADF5uu32l54)

## Project structure ##

```
tsm/
├── main.go                  # Application entry point
├── config/                  # Configuration management
├── logger_factory/          # Logging utilities
├── modes/                   # Mode implementations
├── styles/                  # UI styling
├── tmux/                    # Tmux integration
├── utils/                   # Utility functions
└── view/                    # views
    └── model/               # view models
```
## Architecture

### Mode System

TSM uses a strategy pattern for different operational modes:

See: [Strategy Design Pattern](https://www.amazon.com/Design-Patterns-Elements-Reusable-Object-Oriented/dp/0201633612/ref=sr_1_1?sr=8-1)


```go
type ModeStrategy interface {
    Update(msg tea.Msg) (ModeStrategy, tea.Cmd)
    View() string
    ModeName() string
    Reset()
    GetCurrentSession() string
}
```

## Configuration

On first run, TSM will create a default configuration file at `~/.config/tsm/config.json`.

### Defaults

```json
{
  "search_paths": [
    "~/projects",
    "~/code",
    "~/work",
    "~/dev"
  ],
  "max_depth": 3,
  "theme": "dark"
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `search_paths` | array | `["~/projects", "~/code", "~/work", "~/dev"]` | Directories to scan for projects |
| `max_depth` | number | `3` | How deep to scan subdirectories |
| `theme` | string | `"dark"` | UI theme: `"dark"` or `"light"` |


### Excluded Directories

TSM automatically excludes common non-project directories:
* `node_modules`
* `vendor`
* `build`
* `dist`
* `target`
* `.git`
* `.cache`
* `__pycache__`






