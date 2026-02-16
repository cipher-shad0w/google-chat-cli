# gogchat

A lightweight CLI for sending messages and managing Google Chat spaces.

## Installation

### Homebrew (macOS & Linux)

```bash
brew tap cipher-shad0w/tap
brew install gogchat
```

### Download binary

Download the latest release for your platform from [GitHub Releases](https://github.com/cipher-shad0w/gogchat/releases).

### Build from source

Requires Go 1.25+:

```bash
go install github.com/cipher-shad0w/gogchat/cmd/gogchat@latest
```

> **Note:** When building from source, you must supply your own Google OAuth2 credentials.
> See [CLI.md](CLI.md) for details.

## Quick start

```bash
# Authenticate with your Google account
gogchat auth login

# List your spaces
gogchat spaces list

# Send a message
gogchat messages send --space SPACE_ID --text "Hello from the CLI!"
```

## Documentation

See [CLI.md](CLI.md) for the full command reference and configuration guide.

## License

MIT — see [LICENSE](LICENSE) for details.
