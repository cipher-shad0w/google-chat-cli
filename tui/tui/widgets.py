"""Custom widgets for the gchat TUI."""

from textual import events, work
from textual.app import ComposeResult
from textual.message import Message
from textual.widgets import Label, ListItem, ListView, RichLog, Static, TextArea


class SpaceItem(ListItem):
    """A list item representing a chat space."""

    def __init__(self, space_name: str, display_name: str) -> None:
        """Initialize with space metadata.

        Args:
            space_name: The space resource name (e.g., "spaces/AAAA")
            display_name: Human-readable name to show in the list
        """
        # Create label with display_name, fallback to space_name
        label_text = display_name if display_name else space_name.split("/")[-1]
        super().__init__(Label(label_text))
        self.space_name = space_name
        self.display_name = display_name


class GroupsPanel(Static):
    """Left panel displaying the list of chat groups."""

    class SpaceSelected(Message):
        """Posted when a space is selected."""

        def __init__(self, space_name: str, display_name: str) -> None:
            self.space_name = space_name
            self.display_name = display_name
            super().__init__()

    def compose(self) -> ComposeResult:
        """Create the groups list view."""
        yield ListView(id="groups-list")

    def on_mount(self) -> None:
        """Set the border title and load spaces."""
        self.border_title = "Groups"
        self.load_spaces()

    @work(thread=True)
    def load_spaces(self) -> None:
        """Load spaces from the CLI in a background thread."""
        from tui.cli import list_spaces  # Import here to avoid circular imports

        spaces = list_spaces()

        # Update UI from the worker thread
        self.app.call_from_thread(self._populate_spaces, spaces)

    def _populate_spaces(self, spaces: list[dict]) -> None:
        """Populate the ListView with space items."""
        groups_list = self.query_one("#groups-list", ListView)
        groups_list.clear()

        if not spaces:
            # Show a message if no spaces found
            groups_list.append(ListItem(Label("No spaces found")))
            return

        for space in spaces:
            space_name = space.get("name", "")
            display_name = space.get("displayName", "")
            groups_list.append(SpaceItem(space_name, display_name))

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle space selection."""
        if isinstance(event.item, SpaceItem):
            self.post_message(
                self.SpaceSelected(event.item.space_name, event.item.display_name)
            )


class ChatPanel(Static):
    """Main panel displaying chat messages."""

    def compose(self) -> ComposeResult:
        """Create the chat log view."""
        yield RichLog(id="chat-log", highlight=True, markup=True)

    def on_mount(self) -> None:
        """Set the border title when mounted."""
        self.border_title = "Chat"


class MessageInput(TextArea):
    """A TextArea that sends on Enter and inserts newline on Shift+Enter."""

    class Submitted(Message):
        """Posted when the user presses Enter to send a message."""

        def __init__(self, value: str) -> None:
            self.value = value
            super().__init__()

    def _on_key(self, event: events.Key) -> None:
        if event.key == "enter":
            event.prevent_default()
            event.stop()
            text = self.text.strip()
            if text:
                self.post_message(self.Submitted(text))
                self.clear()
            return
        super()._on_key(event)


class InputPanel(Static):
    """Bottom panel for composing messages."""

    def compose(self) -> ComposeResult:
        """Create the message input area."""
        yield MessageInput(id="message-input")

    def on_mount(self) -> None:
        """Set the border title when mounted."""
        self.border_title = "Message"
