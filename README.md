# Tmux Project Setup in Go

This Go-based application automates the creation of complex `tmux` sessions based on a YAML configuration file. It supports features like pane layouts, global defaults, dependency checks, and Git branch integration.

## Installation

### Prerequisites

- `Go` version 1.23 or higher
- `tmux` installed on your system

### Install the Application System-Wide

1. Clone the repository:

   ```bash
   git clone https://github.com/bartosz-skejcik/tmux-project-setup.git
   cd tmux-project-setup
   ```

2. Build the binary:

   ```bash
   go build -o tmux-setup
   ```

3. Install the binary:

   ```bash
   sudo mv tmux-setup /usr/local/bin/tmux-setup
   ```

4. Verify the installation:
   ```bash
   tmux-setup --help
   ```

## Running the Application

To start the app, navigate to the directory containing your `tmux.conf.yml` file and run:

```bash
  tmux-setup
```

This will:

- Parse the `tmux.conf.yml` file.
- Create a `tmux` session based on the configuration.
- Attach you to the session if no arguments are provided.

## Configuration Options

The configuration file must be named `tmux.conf.yml`. Below is a list of supported properties and their descriptions.

### Top-Level Properties

| Property       | Required | Default Value | Description                                                      |
| -------------- | -------- | ------------- | ---------------------------------------------------------------- |
| `session_name` | No       | `dev`         | Name of the `tmux` session to create.                            |
| `focus_window` | No       | `1`           | Index of the window to focus when attaching to the session.      |
| `defaults`     | No       | `{}`          | Global defaults applied to all windows and panes (see below).    |
| `dependencies` | No       | `[]`          | List of required system commands. Will abort if any are missing. |
| `windows`      | Yes      | `[]`          | List of windows to create in the session.                        |

### `defaults` Properties

| Property          | Required | Default Value | Description                                                          |
| ----------------- | -------- | ------------- | -------------------------------------------------------------------- |
| `directory`       | No       | `""`          | Default directory for all windows and panes unless overridden.       |
| `initial_command` | No       | `""`          | Default initial command for all windows and panes unless overridden. |

### `windows` Properties

| Property          | Required | Default Value | Description                                                             |
| ----------------- | -------- | ------------- | ----------------------------------------------------------------------- |
| `name`            | No       | `window-N`    | Name of the window.                                                     |
| `directory`       | No       | `""`          | Directory to switch to before running any commands in the window.       |
| `initial_command` | No       | `""`          | Command to run in the first pane of the window.                         |
| `layout`          | No       | `""`          | Predefined layout for panes (`even-horizontal`, `even-vertical`, etc.). |
| `git_branch`      | No       | `""`          | Git branch to check out in the window's directory.                      |
| `panes`           | No       | `[]`          | List of panes to create in the window (see below).                      |

### `panes` Properties

| Property          | Required | Default Value | Description                                               |
| ----------------- | -------- | ------------- | --------------------------------------------------------- |
| `directory`       | No       | `""`          | Directory to switch to before running the pane's command. |
| `initial_command` | No       | `""`          | Command to run in the pane.                               |

## Example Configuration Files

### Minimal Example

```yaml
session_name: my_session
windows:
  - name: editor
    initial_command: nvim
```

### Advanced Example

```yaml
session_name: dev
focus_window: 1
defaults:
  directory: ./src
  initial_command: echo \"Welcome to tmux!\"
dependencies:
  - git
  - nvim
windows:
  - name: code
    layout: even-horizontal
    panes:
      - initial_command: nvim
      - directory: ./tests
        initial_command: pytest
  - name: server
    git_branch: main
    initial_command: npm start
  - name: logs
    panes:
      - directory: ./logs
        initial_command: tail -f app.log
      - initial_command: htop
```

## Contribution Guide

We welcome contributions to improve this project. Here's how you can help:

1. **Fork the Repository:**

   - Clone your fork locally: `git clone https://github.com/your-username/tmux-project-setup.git`

2. **Create a Feature Branch:**

   - `git checkout -b feature/my-new-feature`

3. **Write Code:**

   - Follow Go conventions and document your code.

4. **Test Changes:**

   - Ensure your changes don't break existing functionality.

5. **Submit a Pull Request:**

   - Provide a clear description of your changes and why they are useful.

6. **Respond to Feedback:**
   - Be ready to address any comments or requested changes.

---

Feel free to open an issue for questions, bugs, or feature requests!
