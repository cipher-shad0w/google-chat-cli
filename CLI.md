# gogchat CLI Reference

## Installation

```
go install github.com/cipher-shad0w/gogchat/cmd/gogchat@latest
```

---

## Quick Start

```bash
# Authenticate with Google (opens browser)
gogchat auth login

# List spaces you belong to
gogchat spaces list

# Send a message to a space
gogchat messages send spaces/AAAABBBBcccc --text "Hello from gogchat!"
```

---

## Top-Level Help

```
$ gogchat -h
gogchat - A CLI for the Google Chat API

Usage:
  gogchat <command> <subcommand> [flags]

Available Commands:
  auth            Authenticate with Google
  spaces          Manage Google Chat spaces
  messages        Send and manage messages
  members         Manage space memberships
  reactions       Manage message reactions
  attachments     Get attachment metadata
  emoji           Manage custom emojis
  media           Upload and download media
  events          List and inspect space events
  readstate       Manage read state for spaces and threads
  notifications   Manage space notification settings

Global Flags:
  -j, --json        Output in JSON format
      --admin        Use admin access (for admin-only operations)
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat <command> -h" for more information about a command.
Use "gogchat <command> <subcommand> -h" for more information about a subcommand.
```

---

## auth

Authenticate with Google and manage stored credentials.

```
$ gogchat auth -h
Manage authentication for gogchat.

Usage:
  gogchat auth <subcommand> [flags]

Available Subcommands:
  login       Authenticate with Google (OAuth2 browser flow)
  logout      Clear stored authentication tokens
  status      Show current authentication status

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat auth <subcommand> -h" for more information about a subcommand.
```

### auth login

Authenticate with Google using the OAuth2 browser flow. Opens a browser window to complete the authorization. Tokens are stored locally for subsequent commands.

gogchat ships with built-in OAuth2 credentials, so you can simply run the command — no setup required.

```
$ gogchat auth login -h
Authenticate with Google using the OAuth2 browser flow.

Opens your default browser to complete the Google authorization flow.
After authorization, tokens are stored at ~/.config/gogchat/token.json.

Usage:
  gogchat auth login [flags]

Flags:
      --client-id       string   Override the built-in OAuth2 client ID
      --client-secret   string   Override the built-in OAuth2 client secret

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Authenticate (opens browser, click your Google account)
  $ gogchat auth login
  Opening browser for authentication...
  ✓ Successfully logged in!
    Token saved to: /home/user/.config/gogchat/token.json
```

**Advanced: Custom OAuth2 Credentials**

If you want to use your own Google Cloud OAuth2 credentials instead of the
built-in ones, you can override them via:

- **CLI flags**: `gogchat auth login --client-id "YOUR_ID" --client-secret "YOUR_SECRET"`
- **Config file** (`~/.config/gogchat/config.yaml`):
  ```yaml
  client_id: "your-client-id.apps.googleusercontent.com"
  client_secret: "your-client-secret"
  ```
- **Environment variables**: `GOGCHAT_CLIENT_ID` and `GOGCHAT_CLIENT_SECRET`

### auth logout

Clear stored authentication tokens from the local credential store.

```
$ gogchat auth logout -h
Clear stored authentication tokens.

Removes the locally stored OAuth2 tokens. You will need to run
"gogchat auth login" again to re-authenticate.

Usage:
  gogchat auth logout [flags]

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  $ gogchat auth logout
  Logged out. Stored tokens removed.
```

### auth status

Display the current authentication status, including the authenticated user and token expiry.

```
$ gogchat auth status -h
Show current authentication status.

Displays the authenticated user, token validity, and scopes granted.

Usage:
  gogchat auth status [flags]

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  $ gogchat auth status
  Authenticated as: user@example.com
  Token valid until: 2026-02-16T18:30:00Z
  Scopes: chat.spaces, chat.messages, chat.memberships

  $ gogchat auth status --json
  {
    "user": "user@example.com",
    "valid": true,
    "expires": "2026-02-16T18:30:00Z",
    "scopes": ["chat.spaces", "chat.messages", "chat.memberships"]
  }
```

---

## spaces

Manage Google Chat spaces: create, list, search, update, and delete spaces.

```
$ gogchat spaces -h
Manage Google Chat spaces.

Usage:
  gogchat spaces <subcommand> [flags]

Available Subcommands:
  list              List spaces the authenticated user is a member of
  get               Get details of a space
  create            Create a named space
  update            Update a space
  delete            Delete a space (cascading delete)
  search            Search for spaces (admin only)
  setup             Create a space and add members in one step
  find-dm           Find a direct message space with another user
  complete-import   Complete the import process for a space

Global Flags:
  -j, --json        Output in JSON format
      --admin        Use admin access (for admin-only operations)
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat spaces <subcommand> -h" for more information about a subcommand.
```

### spaces list

List spaces the authenticated user is a member of.

```
$ gogchat spaces list -h
List spaces the authenticated user is a member of.

Returns a paginated list of spaces. Use --all to automatically
paginate through all results.

Usage:
  gogchat spaces list [flags]

Flags:
      --filter       string   Filter spaces (e.g. "spaceType = \"SPACE\"")
      --page-size    int      Number of results per page (default 100, max 1000)
      --page-token   string   Page token for pagination
      --all                   Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List all spaces
  $ gogchat spaces list
  NAME                  DISPLAY NAME         TYPE          MEMBERS
  spaces/AAAABBBBcccc   Engineering Team     SPACE         42
  spaces/DDDDEEEEffff   Project Alpha        SPACE         12
  spaces/GGGGHHHHiiii   alice@example.com    DIRECT_MESSAGE 2

  # List only named spaces
  $ gogchat spaces list --filter 'spaceType = "SPACE"'

  # List all spaces as JSON, paginate automatically
  $ gogchat spaces list --all --json

  # List with custom page size
  $ gogchat spaces list --page-size 50
```

### spaces get

Get details of a specific space.

```
$ gogchat spaces get -h
Get details of a space.

Accepts a space ID or full resource name. If a bare ID is provided,
it is automatically expanded to "spaces/<ID>".

Usage:
  gogchat spaces get <space> [flags]

Arguments:
  space   Space ID or resource name (e.g. "spaces/AAAABBBBcccc" or "AAAABBBBcccc")

Flags:
      --admin   Use admin access to retrieve the space

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get space by full resource name
  $ gogchat spaces get spaces/AAAABBBBcccc
  Name:          spaces/AAAABBBBcccc
  Display Name:  Engineering Team
  Type:          SPACE
  Threaded:      true
  Members:       42
  Created:       2025-06-15T10:30:00Z

  # Get space by ID only
  $ gogchat spaces get AAAABBBBcccc

  # Get space as JSON
  $ gogchat spaces get spaces/AAAABBBBcccc --json

  # Get space using admin privileges
  $ gogchat spaces get spaces/AAAABBBBcccc --admin
```

### spaces create

Create a new named space.

```
$ gogchat spaces create -h
Create a named space.

Creates a new Google Chat space with the specified display name. The
authenticated user is automatically added as a member.

Usage:
  gogchat spaces create [flags]

Flags:
      --display-name   string   Display name for the space (required)
      --space-type     string   Type of space: SPACE or GROUP_CHAT (default "SPACE")
      --description    string   Description of the space
      --request-id     string   Unique request ID for idempotency

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Create a new space
  $ gogchat spaces create --display-name "Backend Team"
  Created space: spaces/AAAAMMMNNNoo
  Display Name:  Backend Team
  Type:          SPACE

  # Create a group chat
  $ gogchat spaces create --display-name "Quick Sync" --space-type GROUP_CHAT

  # Create with description and idempotency key
  $ gogchat spaces create \
      --display-name "Release Planning" \
      --description "Coordinate release milestones and blockers" \
      --request-id "create-release-planning-001"
```

### spaces update

Update an existing space.

```
$ gogchat spaces update -h
Update a space's properties.

Update one or more fields on a space. Use --update-mask to specify
which fields to update. If --update-mask is omitted, it is inferred
from the flags provided.

Usage:
  gogchat spaces update <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --display-name    string   New display name
      --description     string   New description
      --history-state   string   History state: HISTORY_ON or HISTORY_OFF
      --update-mask     string   Comma-separated list of fields to update
      --admin                    Use admin access to update the space

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Update the display name
  $ gogchat spaces update spaces/AAAABBBBcccc --display-name "Backend Team v2"

  # Update description and enable history
  $ gogchat spaces update spaces/AAAABBBBcccc \
      --description "Updated team channel" \
      --history-state HISTORY_ON

  # Update with explicit update mask
  $ gogchat spaces update spaces/AAAABBBBcccc \
      --display-name "New Name" \
      --update-mask "displayName"

  # Update as admin
  $ gogchat spaces update spaces/AAAABBBBcccc \
      --display-name "Renamed by Admin" \
      --admin
```

### spaces delete

Delete a space. This is a cascading delete that removes all messages, members, and other data in the space.

```
$ gogchat spaces delete -h
Delete a space (cascading).

Permanently deletes a space and all its contents including messages,
memberships, reactions, and attachments. This action cannot be undone.

Usage:
  gogchat spaces delete <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --admin    Use admin access to delete the space
      --force    Skip confirmation prompt

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Delete a space (with confirmation prompt)
  $ gogchat spaces delete spaces/AAAABBBBcccc
  Are you sure you want to delete space "Engineering Team" (spaces/AAAABBBBcccc)?
  This will permanently delete all messages, members, and data. [y/N]: y
  Deleted space spaces/AAAABBBBcccc.

  # Delete without confirmation
  $ gogchat spaces delete spaces/AAAABBBBcccc --force

  # Delete as admin
  $ gogchat spaces delete spaces/AAAABBBBcccc --admin --force
```

### spaces search

Search for spaces across the organization. Requires admin access.

```
$ gogchat spaces search -h
Search for spaces (admin only).

Search across all spaces in the organization using a query string.
This command requires admin access and automatically sets the --admin flag.

Usage:
  gogchat spaces search [flags]

Flags:
      --query        string   Search query (required). Supports operators such as
                              "customer=\"customers/my_customer\"" combined with
                              filters like "spaceType = \"SPACE\""
      --page-size    int      Number of results per page (default 100, max 1000)
      --page-token   string   Page token for pagination
      --order-by     string   Sort order (e.g. "displayName", "createTime desc")
      --admin                 Use admin access (automatically enabled)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Search for spaces by display name
  $ gogchat spaces search --query 'customer="customers/my_customer"'
  NAME                  DISPLAY NAME         TYPE     MEMBERS
  spaces/AAAABBBBcccc   Engineering Team     SPACE    42
  spaces/XXXXYYYYzzzz   Sales Team           SPACE    28

  # Search with ordering
  $ gogchat spaces search \
      --query 'customer="customers/my_customer"' \
      --order-by "createTime desc" \
      --page-size 10

  # Search and output as JSON
  $ gogchat spaces search \
      --query 'customer="customers/my_customer" AND spaceType = "SPACE"' \
      --json
```

### spaces setup

Create a new space and add members in a single step.

```
$ gogchat spaces setup -h
Create a space and add members in one step.

Sets up a new space with the specified display name and immediately
adds the given members. This is equivalent to calling "spaces create"
followed by "members add" for each member.

Usage:
  gogchat spaces setup [flags]

Flags:
      --display-name   string   Display name for the space (required)
      --members        string   Comma-separated list of user resource names to add
      --space-type     string   Type of space: SPACE or GROUP_CHAT (default "SPACE")

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Create a space with three members
  $ gogchat spaces setup \
      --display-name "Project Phoenix" \
      --members "users/123456789,users/987654321,users/111222333"
  Created space: spaces/AAAANNNNoooo
  Display Name:  Project Phoenix
  Added 3 members.

  # Create a group chat with members
  $ gogchat spaces setup \
      --display-name "Quick Standup" \
      --space-type GROUP_CHAT \
      --members "users/123456789,users/987654321"
```

### spaces find-dm

Find the direct message space between the authenticated user and another user.

```
$ gogchat spaces find-dm -h
Find a direct message space with a user.

Returns the DM space between the authenticated user and the
specified user, if one exists.

Usage:
  gogchat spaces find-dm [flags]

Flags:
      --user   string   User resource name (required, e.g. "users/123456789")

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Find DM with a user
  $ gogchat spaces find-dm --user users/123456789
  Name:          spaces/GGGGHHHHiiii
  Display Name:  alice@example.com
  Type:          DIRECT_MESSAGE

  # Find DM and output as JSON
  $ gogchat spaces find-dm --user users/123456789 --json
```

### spaces complete-import

Complete the import process for a space that was created via the import mode.

```
$ gogchat spaces complete-import -h
Complete the import process for a space.

Finalizes the import of a space that was created in import mode.
After completion, the space becomes a regular space and is visible
to all members.

Usage:
  gogchat spaces complete-import <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Complete import for a space
  $ gogchat spaces complete-import spaces/AAAABBBBcccc
  Import completed for space spaces/AAAABBBBcccc.
```

---

## messages

Send, list, update, and delete messages in Google Chat spaces.

```
$ gogchat messages -h
Send and manage messages in Google Chat spaces.

Usage:
  gogchat messages <subcommand> [flags]

Available Subcommands:
  list      List messages in a space
  get       Get details of a message
  send      Send a message to a space
  update    Update a message
  delete    Delete a message
  replace   Full replacement update (PUT) of a message

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat messages <subcommand> -h" for more information about a subcommand.
```

### messages list

List messages in a space.

```
$ gogchat messages list -h
List messages in a space.

Returns a paginated list of messages in the specified space. Messages
are returned in reverse chronological order by default. Use --all to
automatically paginate through all results.

Usage:
  gogchat messages list <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --page-size      int      Number of results per page (default 25, max 1000)
      --page-token     string   Page token for pagination
      --filter         string   Filter messages (e.g. "createTime > \"2025-01-01T00:00:00Z\"")
      --order-by       string   Sort order (e.g. "createTime desc")
      --show-deleted              Include deleted messages in the list
      --all                       Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List recent messages in a space
  $ gogchat messages list spaces/AAAABBBBcccc
  MESSAGE                                          SENDER              TEXT                     CREATED
  spaces/AAAABBBBcccc/messages/123456.789012        users/111222333     Hey team, standup time!  2026-02-16T09:00:00Z
  spaces/AAAABBBBcccc/messages/123456.789013        users/444555666     On my way!               2026-02-16T09:01:00Z

  # List messages with a filter
  $ gogchat messages list spaces/AAAABBBBcccc \
      --filter 'createTime > "2026-02-01T00:00:00Z"'

  # List all messages as JSON
  $ gogchat messages list spaces/AAAABBBBcccc --all --json

  # Include deleted messages
  $ gogchat messages list spaces/AAAABBBBcccc --show-deleted

  # Custom page size and order
  $ gogchat messages list spaces/AAAABBBBcccc --page-size 100 --order-by "createTime asc"
```

### messages get

Get details of a specific message.

```
$ gogchat messages get -h
Get details of a message.

Retrieves the full message resource including text, sender, thread
information, annotations, and attachment metadata.

Usage:
  gogchat messages get <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get a message
  $ gogchat messages get spaces/AAAABBBBcccc/messages/123456.789012
  Name:      spaces/AAAABBBBcccc/messages/123456.789012
  Sender:    users/111222333
  Created:   2026-02-16T09:00:00Z
  Text:      Hey team, standup time!
  Thread:    spaces/AAAABBBBcccc/threads/abcDEF123

  # Get as JSON
  $ gogchat messages get spaces/AAAABBBBcccc/messages/123456.789012 --json
```

### messages send

Send a message to a space.

```
$ gogchat messages send -h
Send a message to a space.

Sends a text message to the specified space. Supports threading by
providing a --thread-key or using --reply-option to control reply behavior.

Usage:
  gogchat messages send <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --text           string   Message text content (required). Supports Google Chat
                                formatting (e.g. *bold*, _italic_, `code`,
                                ```code block```, ~strikethrough~)
      --thread-key     string   Thread key for creating or replying in a named thread
      --request-id     string   Unique request ID for idempotency
      --message-id     string   Custom message ID (must start with "client-")
      --reply-option   string   Reply behavior:
                                  REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD - reply to thread
                                    or create new if thread not found
                                  REPLY_MESSAGE_OR_FAIL - reply to thread or fail

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Send a simple message
  $ gogchat messages send spaces/AAAABBBBcccc --text "Hello, team!"
  Sent: spaces/AAAABBBBcccc/messages/678901.234567

  # Send a threaded message
  $ gogchat messages send spaces/AAAABBBBcccc \
      --text "This is a reply" \
      --thread-key "standup-2026-02-16" \
      --reply-option REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD

  # Send with formatting
  $ gogchat messages send spaces/AAAABBBBcccc \
      --text "*Build passed* for \`main\` branch. See details at https://ci.example.com/123"

  # Send with custom message ID
  $ gogchat messages send spaces/AAAABBBBcccc \
      --text "Deployment complete" \
      --message-id "client-deploy-20260216-001"

  # Send quietly (only output the message name)
  $ gogchat messages send spaces/AAAABBBBcccc --text "Silent ping" --quiet
  spaces/AAAABBBBcccc/messages/678901.234568
```

### messages update

Update an existing message.

```
$ gogchat messages update -h
Update a message.

Updates one or more fields of an existing message. Use --update-mask to
specify which fields to update. If --update-mask is omitted, it is
inferred from the flags provided.

Usage:
  gogchat messages update <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
      --text            string   New message text content
      --update-mask     string   Comma-separated list of fields to update (e.g. "text")
      --allow-missing              Create the message if it does not exist

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Update message text
  $ gogchat messages update spaces/AAAABBBBcccc/messages/123456.789012 \
      --text "Updated: Hello, team! (edited)"

  # Update with explicit update mask
  $ gogchat messages update spaces/AAAABBBBcccc/messages/123456.789012 \
      --text "Corrected text" \
      --update-mask "text"

  # Update or create if missing
  $ gogchat messages update spaces/AAAABBBBcccc/messages/123456.789012 \
      --text "This message may or may not exist" \
      --allow-missing
```

### messages delete

Delete a message.

```
$ gogchat messages delete -h
Delete a message.

Permanently deletes a message. By default, prompts for confirmation.

Usage:
  gogchat messages delete <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
      --force           Skip confirmation prompt
      --force-threads   Also delete all threaded replies to this message

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Delete with confirmation
  $ gogchat messages delete spaces/AAAABBBBcccc/messages/123456.789012
  Are you sure you want to delete message spaces/AAAABBBBcccc/messages/123456.789012? [y/N]: y
  Deleted message spaces/AAAABBBBcccc/messages/123456.789012.

  # Delete without confirmation
  $ gogchat messages delete spaces/AAAABBBBcccc/messages/123456.789012 --force

  # Delete message and all threaded replies
  $ gogchat messages delete spaces/AAAABBBBcccc/messages/123456.789012 \
      --force --force-threads
```

### messages replace

Perform a full replacement update (PUT) of a message, replacing the entire message resource.

```
$ gogchat messages replace -h
Full replacement update (PUT) of a message.

Replaces the entire message resource. Fields not specified in the
request are reset to their default values. This differs from
"messages update" (PATCH) which only modifies the specified fields.

Usage:
  gogchat messages replace <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
      --text            string   Message text content
      --update-mask     string   Comma-separated list of fields to update
      --allow-missing              Create the message if it does not exist

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Replace entire message content
  $ gogchat messages replace spaces/AAAABBBBcccc/messages/123456.789012 \
      --text "Completely replaced message content"

  # Replace with allow-missing (upsert)
  $ gogchat messages replace spaces/AAAABBBBcccc/messages/123456.789012 \
      --text "Replacement text" \
      --allow-missing
```

---

## members

Manage memberships in Google Chat spaces.

```
$ gogchat members -h
Manage space memberships.

Usage:
  gogchat members <subcommand> [flags]

Available Subcommands:
  list      List members of a space
  get       Get member details
  add       Add a member to a space
  update    Update a membership (e.g. change role)
  remove    Remove a member from a space

Global Flags:
  -j, --json        Output in JSON format
      --admin        Use admin access (for admin-only operations)
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat members <subcommand> -h" for more information about a subcommand.
```

### members list

List members of a space.

```
$ gogchat members list -h
List members of a space.

Returns a paginated list of memberships in the specified space.
Use --all to automatically paginate through all results.

Usage:
  gogchat members list <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --page-size      int      Number of results per page (default 100, max 1000)
      --page-token     string   Page token for pagination
      --filter         string   Filter memberships (e.g. "member.type = \"HUMAN\"")
      --show-invited              Include invited members who have not yet joined
      --show-groups               Include Google Groups in the results
      --admin                     Use admin access to list members
      --all                       Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List all members of a space
  $ gogchat members list spaces/AAAABBBBcccc
  MEMBER                                        USER              ROLE           STATE
  spaces/AAAABBBBcccc/members/111222333          users/111222333   ROLE_MANAGER   JOINED
  spaces/AAAABBBBcccc/members/444555666          users/444555666   ROLE_MEMBER    JOINED
  spaces/AAAABBBBcccc/members/777888999          users/777888999   ROLE_MEMBER    JOINED

  # List only human members
  $ gogchat members list spaces/AAAABBBBcccc --filter 'member.type = "HUMAN"'

  # Include invited members
  $ gogchat members list spaces/AAAABBBBcccc --show-invited

  # List all members as admin
  $ gogchat members list spaces/AAAABBBBcccc --admin --all --json
```

### members get

Get details of a specific membership.

```
$ gogchat members get -h
Get member details.

Retrieves details of a specific membership in a space.

Usage:
  gogchat members get <member> [flags]

Arguments:
  member   Membership resource name (e.g. "spaces/AAAABBBBcccc/members/111222333")

Flags:
      --admin   Use admin access to retrieve the membership

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get membership details
  $ gogchat members get spaces/AAAABBBBcccc/members/111222333
  Name:    spaces/AAAABBBBcccc/members/111222333
  User:    users/111222333
  Role:    ROLE_MANAGER
  State:   JOINED
  Created: 2025-06-15T10:30:00Z

  # Get as JSON
  $ gogchat members get spaces/AAAABBBBcccc/members/111222333 --json

  # Get using admin access
  $ gogchat members get spaces/AAAABBBBcccc/members/111222333 --admin
```

### members add

Add a member to a space.

```
$ gogchat members add -h
Add a member to a space.

Adds a user to the specified space with the given role. The user
receives a notification and the space appears in their space list.

Usage:
  gogchat members add <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --user    string   User resource name to add (e.g. "users/123456789")
      --role    string   Member role: ROLE_MEMBER or ROLE_MANAGER (default "ROLE_MEMBER")
      --admin              Use admin access to add the member

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Add a member with default role
  $ gogchat members add spaces/AAAABBBBcccc --user users/123456789
  Added users/123456789 to spaces/AAAABBBBcccc as ROLE_MEMBER.

  # Add a member as manager
  $ gogchat members add spaces/AAAABBBBcccc \
      --user users/123456789 \
      --role ROLE_MANAGER

  # Add member as admin
  $ gogchat members add spaces/AAAABBBBcccc \
      --user users/987654321 \
      --admin
```

### members update

Update a membership (e.g., change role).

```
$ gogchat members update -h
Update a membership (e.g. change role).

Update properties of an existing membership. Commonly used to promote
or demote members between ROLE_MEMBER and ROLE_MANAGER.

Usage:
  gogchat members update <member> [flags]

Arguments:
  member   Membership resource name (e.g. "spaces/AAAABBBBcccc/members/111222333")

Flags:
      --role          string   New role: ROLE_MEMBER or ROLE_MANAGER
      --update-mask   string   Comma-separated list of fields to update (e.g. "role")
      --admin                  Use admin access to update the membership

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Promote a member to manager
  $ gogchat members update spaces/AAAABBBBcccc/members/444555666 \
      --role ROLE_MANAGER
  Updated spaces/AAAABBBBcccc/members/444555666: role -> ROLE_MANAGER.

  # Demote a manager to member
  $ gogchat members update spaces/AAAABBBBcccc/members/444555666 \
      --role ROLE_MEMBER

  # Update as admin with explicit mask
  $ gogchat members update spaces/AAAABBBBcccc/members/444555666 \
      --role ROLE_MANAGER \
      --update-mask "role" \
      --admin
```

### members remove

Remove a member from a space.

```
$ gogchat members remove -h
Remove a member from a space.

Removes the specified membership. The user will no longer see the
space in their space list and will lose access to space contents.

Usage:
  gogchat members remove <member> [flags]

Arguments:
  member   Membership resource name (e.g. "spaces/AAAABBBBcccc/members/111222333")

Flags:
      --admin   Use admin access to remove the member
      --force   Skip confirmation prompt

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Remove a member (with confirmation)
  $ gogchat members remove spaces/AAAABBBBcccc/members/444555666
  Are you sure you want to remove spaces/AAAABBBBcccc/members/444555666? [y/N]: y
  Removed spaces/AAAABBBBcccc/members/444555666.

  # Remove without confirmation
  $ gogchat members remove spaces/AAAABBBBcccc/members/444555666 --force

  # Remove as admin
  $ gogchat members remove spaces/AAAABBBBcccc/members/444555666 --admin --force
```

---

## reactions

Manage reactions on messages.

```
$ gogchat reactions -h
Manage message reactions.

Usage:
  gogchat reactions <subcommand> [flags]

Available Subcommands:
  list      List reactions on a message
  add       Add a reaction to a message
  remove    Remove a reaction

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat reactions <subcommand> -h" for more information about a subcommand.
```

### reactions list

List reactions on a message.

```
$ gogchat reactions list -h
List reactions on a message.

Returns a paginated list of reactions on the specified message.
Use --all to automatically paginate through all results.

Usage:
  gogchat reactions list <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
      --page-size    int      Number of results per page (default 25, max 200)
      --page-token   string   Page token for pagination
      --filter       string   Filter reactions (e.g. "emoji.unicode = \"👍\"" or
                              "user.name = \"users/123456789\"")
      --all                   Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List all reactions on a message
  $ gogchat reactions list spaces/AAAABBBBcccc/messages/123456.789012
  REACTION                                                              USER              EMOJI
  spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR111           users/111222333    👍
  spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR222           users/444555666    🎉
  spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR333           users/777888999    👍

  # Filter for a specific emoji
  $ gogchat reactions list spaces/AAAABBBBcccc/messages/123456.789012 \
      --filter 'emoji.unicode = "👍"'

  # List as JSON
  $ gogchat reactions list spaces/AAAABBBBcccc/messages/123456.789012 --all --json
```

### reactions add

Add a reaction to a message.

```
$ gogchat reactions add -h
Add a reaction to a message.

Adds a Unicode emoji or custom emoji reaction to the specified message.

Usage:
  gogchat reactions add <message> [flags]

Arguments:
  message   Message resource name (e.g. "spaces/AAAABBBBcccc/messages/123456.789012")

Flags:
      --emoji   string   Emoji to react with. Either a Unicode emoji character
                         (e.g. "👍", "🎉", "❤️") or a custom emoji UID

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Add a thumbs-up reaction
  $ gogchat reactions add spaces/AAAABBBBcccc/messages/123456.789012 --emoji "👍"
  Added reaction 👍 to spaces/AAAABBBBcccc/messages/123456.789012.

  # Add a custom emoji reaction
  $ gogchat reactions add spaces/AAAABBBBcccc/messages/123456.789012 \
      --emoji "custom-emoji-uid-12345"

  # Add a party popper
  $ gogchat reactions add spaces/AAAABBBBcccc/messages/123456.789012 --emoji "🎉"
```

### reactions remove

Remove a reaction.

```
$ gogchat reactions remove -h
Remove a reaction.

Removes the specified reaction from a message.

Usage:
  gogchat reactions remove <reaction> [flags]

Arguments:
  reaction   Reaction resource name
             (e.g. "spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR111")

Flags:
      --force   Skip confirmation prompt

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Remove a reaction (with confirmation)
  $ gogchat reactions remove spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR111
  Are you sure you want to remove this reaction? [y/N]: y
  Removed reaction.

  # Remove without confirmation
  $ gogchat reactions remove \
      spaces/AAAABBBBcccc/messages/123456.789012/reactions/RRR111 --force
```

---

## attachments

Get metadata for message attachments.

```
$ gogchat attachments -h
Get attachment metadata.

Usage:
  gogchat attachments <subcommand> [flags]

Available Subcommands:
  get       Get attachment metadata

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat attachments <subcommand> -h" for more information about a subcommand.
```

### attachments get

Get metadata for a specific attachment.

```
$ gogchat attachments get -h
Get attachment metadata.

Retrieves metadata for an attachment including its name, content type,
size, and download URI. Use "gogchat media download" to download the
actual file content.

Usage:
  gogchat attachments get <attachment> [flags]

Arguments:
  attachment   Attachment resource name
               (e.g. "spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get attachment metadata
  $ gogchat attachments get \
      spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001
  Name:          spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001
  Content Type:  application/pdf
  Filename:      quarterly-report.pdf
  Size:          2,456,789 bytes
  Source:         UPLOADED_CONTENT
  Download URI:  https://chat.googleapis.com/v1/media/...

  # Get as JSON
  $ gogchat attachments get \
      spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001 --json
```

---

## emoji

Manage custom emojis in Google Chat.

```
$ gogchat emoji -h
Manage custom emojis.

Usage:
  gogchat emoji <subcommand> [flags]

Available Subcommands:
  list      List custom emojis
  get       Get custom emoji details
  create    Create a custom emoji
  delete    Delete a custom emoji

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat emoji <subcommand> -h" for more information about a subcommand.
```

### emoji list

List custom emojis.

```
$ gogchat emoji list -h
List custom emojis.

Returns a paginated list of custom emojis available to the authenticated
user. Use --all to automatically paginate through all results.

Usage:
  gogchat emoji list [flags]

Flags:
      --page-size    int      Number of results per page (default 25, max 200)
      --page-token   string   Page token for pagination
      --filter       string   Filter custom emojis (e.g. "creator.name = \"users/123456789\"")
      --all                   Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List all custom emojis
  $ gogchat emoji list
  NAME                    SHORT NAME       CREATOR
  customEmojis/AAA111     :party-parrot:   users/111222333
  customEmojis/BBB222     :ship-it:        users/444555666
  customEmojis/CCC333     :lgtm:           users/777888999

  # List all custom emojis as JSON
  $ gogchat emoji list --all --json

  # Filter emojis by creator
  $ gogchat emoji list --filter 'creator.name = "users/111222333"'
```

### emoji get

Get details of a specific custom emoji.

```
$ gogchat emoji get -h
Get custom emoji details.

Retrieves details of a custom emoji including its name, short name,
image URI, and creator.

Usage:
  gogchat emoji get <emoji> [flags]

Arguments:
  emoji   Custom emoji resource name (e.g. "customEmojis/AAA111")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get custom emoji details
  $ gogchat emoji get customEmojis/AAA111
  Name:        customEmojis/AAA111
  Short Name:  :party-parrot:
  Creator:     users/111222333
  Created:     2025-09-10T14:20:00Z

  # Get as JSON
  $ gogchat emoji get customEmojis/AAA111 --json
```

### emoji create

Create a custom emoji.

```
$ gogchat emoji create -h
Create a custom emoji.

Creates a new custom emoji by uploading an image file. Supported
formats are PNG, GIF, and JPEG. Maximum file size is 256 KB.
Image dimensions should be 128x128 pixels.

Usage:
  gogchat emoji create [flags]

Flags:
      --name          string   Short name for the emoji, without colons (required)
      --image-file    string   Path to the image file (required). Supported formats:
                               PNG, GIF, JPEG. Max size: 256 KB

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Create a custom emoji from a PNG file
  $ gogchat emoji create --name "ship-it" --image-file ./ship-it.png
  Created custom emoji: customEmojis/DDD444
  Short Name:  :ship-it:

  # Create from a GIF
  $ gogchat emoji create --name "party-parrot" --image-file ~/emojis/parrot.gif
```

### emoji delete

Delete a custom emoji.

```
$ gogchat emoji delete -h
Delete a custom emoji.

Permanently deletes a custom emoji. Messages that reference this
emoji will show a placeholder instead.

Usage:
  gogchat emoji delete <emoji> [flags]

Arguments:
  emoji   Custom emoji resource name (e.g. "customEmojis/AAA111")

Flags:
      --force   Skip confirmation prompt

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Delete a custom emoji (with confirmation)
  $ gogchat emoji delete customEmojis/AAA111
  Are you sure you want to delete custom emoji ":party-parrot:" (customEmojis/AAA111)? [y/N]: y
  Deleted custom emoji customEmojis/AAA111.

  # Delete without confirmation
  $ gogchat emoji delete customEmojis/AAA111 --force
```

---

## media

Upload and download media files in Google Chat.

```
$ gogchat media -h
Upload and download media.

Usage:
  gogchat media <subcommand> [flags]

Available Subcommands:
  upload      Upload an attachment to a space
  download    Download media

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat media <subcommand> -h" for more information about a subcommand.
```

### media upload

Upload an attachment to a space.

```
$ gogchat media upload -h
Upload an attachment to a space.

Uploads a file as an attachment to the specified space. The uploaded
file can then be referenced when sending messages. Maximum file size
is 200 MB.

Usage:
  gogchat media upload <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --file   string   Path to the file to upload (required)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Upload a file
  $ gogchat media upload spaces/AAAABBBBcccc --file ./report.pdf
  Uploaded: spaces/AAAABBBBcccc/attachments/ATT002
  Filename:     report.pdf
  Content Type: application/pdf
  Size:         1,234,567 bytes

  # Upload and get JSON output
  $ gogchat media upload spaces/AAAABBBBcccc --file ./screenshot.png --json
```

### media download

Download media from Google Chat.

```
$ gogchat media download -h
Download media.

Downloads a media resource (attachment) from Google Chat. If no
output path is specified, the file is saved to the current directory
using the original filename.

Usage:
  gogchat media download <resource> [flags]

Arguments:
  resource   Media resource name or URI

Flags:
  -o, --output   string   Output file path. If not specified, uses the
                           original filename in the current directory

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Download to current directory
  $ gogchat media download spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001
  Downloaded: quarterly-report.pdf (2,456,789 bytes)

  # Download to a specific path
  $ gogchat media download \
      spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001 \
      -o ~/Downloads/report.pdf

  # Download quietly
  $ gogchat media download \
      spaces/AAAABBBBcccc/messages/123456.789012/attachments/ATT001 \
      -o ./file.pdf --quiet
```

---

## events

List and inspect space events.

```
$ gogchat events -h
List and inspect space events.

Space events represent changes to a space, such as new members,
removed members, new messages, and reactions.

Usage:
  gogchat events <subcommand> [flags]

Available Subcommands:
  list      List space events
  get       Get space event details

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat events <subcommand> -h" for more information about a subcommand.
```

### events list

List space events.

```
$ gogchat events list -h
List space events.

Returns a paginated list of events in the specified space. A filter
specifying at least one event type is required.

Usage:
  gogchat events list <space> [flags]

Arguments:
  space   Space resource name (e.g. "spaces/AAAABBBBcccc")

Flags:
      --filter       string   Filter events (required). Must specify at least one
                              event_type. Examples:
                                'event_types:"google.workspace.chat.message.v1.created"'
                                'event_types:"google.workspace.chat.membership.v1.created"
                                 AND event_types:"google.workspace.chat.membership.v1.deleted"'
      --page-size    int      Number of results per page (default 100, max 1000)
      --page-token   string   Page token for pagination
      --all                   Automatically paginate through all results

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # List new message events
  $ gogchat events list spaces/AAAABBBBcccc \
      --filter 'event_types:"google.workspace.chat.message.v1.created"'
  EVENT                                         TYPE                                             TIME
  spaces/AAAABBBBcccc/spaceEvents/EVT001        google.workspace.chat.message.v1.created          2026-02-16T09:00:00Z
  spaces/AAAABBBBcccc/spaceEvents/EVT002        google.workspace.chat.message.v1.created          2026-02-16T09:05:00Z

  # List membership changes
  $ gogchat events list spaces/AAAABBBBcccc \
      --filter 'event_types:"google.workspace.chat.membership.v1.created" AND event_types:"google.workspace.chat.membership.v1.deleted"'

  # List all events as JSON
  $ gogchat events list spaces/AAAABBBBcccc \
      --filter 'event_types:"google.workspace.chat.message.v1.created"' \
      --all --json
```

### events get

Get details of a specific space event.

```
$ gogchat events get -h
Get space event details.

Retrieves details of a specific space event, including its type,
time, and the affected resource.

Usage:
  gogchat events get <event> [flags]

Arguments:
  event   Space event resource name (e.g. "spaces/AAAABBBBcccc/spaceEvents/EVT001")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get event details
  $ gogchat events get spaces/AAAABBBBcccc/spaceEvents/EVT001
  Name:       spaces/AAAABBBBcccc/spaceEvents/EVT001
  Type:       google.workspace.chat.message.v1.created
  Time:       2026-02-16T09:00:00Z
  Resource:   spaces/AAAABBBBcccc/messages/123456.789012

  # Get as JSON
  $ gogchat events get spaces/AAAABBBBcccc/spaceEvents/EVT001 --json
```

---

## readstate

Manage read state for spaces and threads.

```
$ gogchat readstate -h
Manage read state for spaces and threads.

Read state tracks which messages the user has read in a space or thread.
Updating read state marks messages as read up to the specified time.

Usage:
  gogchat readstate <subcommand> [flags]

Available Subcommands:
  get-space       Get the space read state for a user
  update-space    Update the space read state for a user
  get-thread      Get the thread read state for a user

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat readstate <subcommand> -h" for more information about a subcommand.
```

### readstate get-space

Get the space read state for a user.

```
$ gogchat readstate get-space -h
Get space read state.

Retrieves the read state for the authenticated user in the specified
space, including the timestamp of the last read message.

Usage:
  gogchat readstate get-space <user_space> [flags]

Arguments:
  user_space   Space read state resource name
               (e.g. "users/me/spaces/AAAABBBBcccc/spaceReadState"
                or "users/123456789/spaces/AAAABBBBcccc/spaceReadState")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get space read state for the current user
  $ gogchat readstate get-space users/me/spaces/AAAABBBBcccc/spaceReadState
  Name:            users/me/spaces/AAAABBBBcccc/spaceReadState
  Last Read Time:  2026-02-16T08:45:00Z

  # Get as JSON
  $ gogchat readstate get-space users/me/spaces/AAAABBBBcccc/spaceReadState --json
```

### readstate update-space

Update the space read state for a user.

```
$ gogchat readstate update-space -h
Update space read state.

Updates the read state for the authenticated user in the specified
space. This marks all messages up to the given time as read.

Usage:
  gogchat readstate update-space <user_space> [flags]

Arguments:
  user_space   Space read state resource name
               (e.g. "users/me/spaces/AAAABBBBcccc/spaceReadState")

Flags:
      --last-read-time   string   Timestamp to mark as last read (RFC 3339 format,
                                  e.g. "2026-02-16T09:00:00Z")
      --update-mask      string   Comma-separated list of fields to update

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Mark space as read up to a specific time
  $ gogchat readstate update-space users/me/spaces/AAAABBBBcccc/spaceReadState \
      --last-read-time "2026-02-16T09:00:00Z"
  Updated read state for users/me/spaces/AAAABBBBcccc/spaceReadState.
  Last Read Time: 2026-02-16T09:00:00Z

  # Update with explicit mask
  $ gogchat readstate update-space users/me/spaces/AAAABBBBcccc/spaceReadState \
      --last-read-time "2026-02-16T12:00:00Z" \
      --update-mask "lastReadTime"
```

### readstate get-thread

Get the thread read state for a user.

```
$ gogchat readstate get-thread -h
Get thread read state.

Retrieves the read state for the authenticated user in the specified
thread, including the timestamp of the last read message in that thread.

Usage:
  gogchat readstate get-thread <user_thread> [flags]

Arguments:
  user_thread   Thread read state resource name
                (e.g. "users/me/spaces/AAAABBBBcccc/threads/abcDEF123/threadReadState")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get thread read state
  $ gogchat readstate get-thread \
      users/me/spaces/AAAABBBBcccc/threads/abcDEF123/threadReadState
  Name:            users/me/spaces/AAAABBBBcccc/threads/abcDEF123/threadReadState
  Last Read Time:  2026-02-16T09:15:00Z

  # Get as JSON
  $ gogchat readstate get-thread \
      users/me/spaces/AAAABBBBcccc/threads/abcDEF123/threadReadState --json
```

---

## notifications

Manage notification settings for spaces.

```
$ gogchat notifications -h
Manage space notification settings.

Control how and when you receive notifications for activity in a space.

Usage:
  gogchat notifications <subcommand> [flags]

Available Subcommands:
  get       Get notification settings for a space
  update    Update notification settings for a space

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Use "gogchat notifications <subcommand> -h" for more information about a subcommand.
```

### notifications get

Get notification settings for a space.

```
$ gogchat notifications get -h
Get notification setting for a space.

Retrieves the notification setting for the authenticated user in the
specified space.

Usage:
  gogchat notifications get <user_space> [flags]

Arguments:
  user_space   Space notification setting resource name
               (e.g. "users/me/spaces/AAAABBBBcccc/spaceNotificationSetting")

Flags:
  (none)

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Get notification settings
  $ gogchat notifications get users/me/spaces/AAAABBBBcccc/spaceNotificationSetting
  Name:                  users/me/spaces/AAAABBBBcccc/spaceNotificationSetting
  Notification Setting:  ALL_NEW_MESSAGES
  Mute Setting:          UNMUTED

  # Get as JSON
  $ gogchat notifications get \
      users/me/spaces/AAAABBBBcccc/spaceNotificationSetting --json
```

### notifications update

Update notification settings for a space.

```
$ gogchat notifications update -h
Update notification setting for a space.

Updates the notification setting for the authenticated user in the
specified space. Control notification level and mute state.

Usage:
  gogchat notifications update <user_space> [flags]

Arguments:
  user_space   Space notification setting resource name
               (e.g. "users/me/spaces/AAAABBBBcccc/spaceNotificationSetting")

Flags:
      --notification-setting   string   Notification level:
                                          ALL_NEW_MESSAGES - notify on all messages
                                          MENTIONS_AND_THREADS_FOLLOWING - notify on
                                            @mentions and followed threads only
                                          OFF - no notifications
      --mute-setting           string   Mute state: MUTED or UNMUTED
      --update-mask            string   Comma-separated list of fields to update

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet        Suppress non-essential output
  -v, --verbose      Enable verbose/debug output
      --config       Path to config file (default: ~/.config/gogchat/config.yaml)
  -h, --help         Show help for a command

Examples:
  # Set notifications to mentions only
  $ gogchat notifications update \
      users/me/spaces/AAAABBBBcccc/spaceNotificationSetting \
      --notification-setting MENTIONS_AND_THREADS_FOLLOWING
  Updated notification setting for users/me/spaces/AAAABBBBcccc/spaceNotificationSetting.
  Notification Setting: MENTIONS_AND_THREADS_FOLLOWING

  # Mute a space
  $ gogchat notifications update \
      users/me/spaces/AAAABBBBcccc/spaceNotificationSetting \
      --mute-setting MUTED

  # Turn off notifications and mute
  $ gogchat notifications update \
      users/me/spaces/AAAABBBBcccc/spaceNotificationSetting \
      --notification-setting OFF \
      --mute-setting MUTED \
      --update-mask "notificationSetting,muteSetting"
```

---

## Configuration

### Config File

`gogchat` reads configuration from `~/.config/gogchat/config.yaml` by default. You can override this with the `--config` flag.

```yaml
# ~/.config/gogchat/config.yaml

# Default output format (json or text)
output: text

# Default page size for list operations
page_size: 100

# OAuth2 client configuration (for custom OAuth apps)
client_id: "your-client-id.apps.googleusercontent.com"
client_secret: "your-client-secret"

# Token storage path (default: ~/.config/gogchat/credentials.json)
credentials_path: "~/.config/gogchat/credentials.json"
```

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `GOGCHAT_CONFIG` | Path to config file | `~/.config/gogchat/config.yaml` |
| `GOGCHAT_OUTPUT` | Default output format (`json` or `text`) | `text` |
| `GOGCHAT_CLIENT_ID` | OAuth2 client ID | (built-in) |
| `GOGCHAT_CLIENT_SECRET` | OAuth2 client secret | (built-in) |
| `GOGCHAT_CREDENTIALS` | Path to stored credentials | `~/.config/gogchat/credentials.json` |
| `NO_COLOR` | Disable colored output when set | (unset) |

Environment variables take precedence over config file values. Command-line flags take precedence over both.

---

## Global Flags Reference

| Flag | Short | Description |
|---|---|---|
| `--json` | `-j` | Output in JSON format. All commands support JSON output for scripting and automation. |
| `--admin` | | Use admin access (Workspace admin privileges). Required for some operations like `spaces search`. Automatically set where required. |
| `--quiet` | `-q` | Suppress non-essential output. Only print resource names or critical errors. Useful in scripts. |
| `--verbose` | `-v` | Enable verbose/debug output. Prints HTTP request and response details for troubleshooting. |
| `--config` | | Path to config file. Overrides the default path of `~/.config/gogchat/config.yaml`. |
| `--help` | `-h` | Show help for any command or subcommand. |

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | General error (e.g. invalid arguments, API error) |
| `2` | Authentication error (not logged in or token expired) |
| `3` | Permission denied (insufficient scopes or not a space member) |
| `4` | Resource not found |
| `5` | Rate limited (Google API quota exceeded) |

---

## See Also

- [Google Chat API Reference](https://developers.google.com/workspace/chat/api/reference/rest)
- [Google Chat API Guides](https://developers.google.com/workspace/chat/overview)
