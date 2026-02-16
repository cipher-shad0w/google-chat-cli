"""Entry point for running the TUI as a module: python -m tui"""

from tui.app import ChatApp


def main() -> None:
    """Run the gchat TUI application."""
    app = ChatApp()
    app.run()


if __name__ == "__main__":
    main()
