# 🚀 Tmux Project Setup in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/bartosz-skejcik/tmux-setup)](https://goreportcard.com/report/github.com/bartosz-skejcik/tmux-setup)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.com/bartosz-skejcik/tmux-setup.svg?branch=main)](https://travis-ci.com/bartosz-skejcik/tmux-setup)

This Go-based application automates the creation of simple & complex `tmux` sessions based on a YAML configuration file. It supports features like pane layouts, global defaults, dependency checks, and Git branch integration.

## ✨ Features

-   Automated creation of complex `tmux` sessions
-   YAML-based configuration
-   Pane layouts and global defaults
-   Dependency checks
-   Git branch integration
-   Interactive configuration wizard
-   Template support for reusable configurations

## 📚 Table of Contents

-   [Installation](#-installation)
    -   [Prerequisites](#prerequisites)
    -   [Install the Application System-Wide](#install-the-application-system-wide)
-   [Running the Application](#-running-the-application)
-   [Using the Configuration Wizard](#-using-the-configuration-wizard)
    -   [Why Use the Wizard?](#why-use-the-wizard)
    -   [Creating a Configuration File](#creating-a-configuration-file)
    -   [Creating a Template](#creating-a-template)
    -   [What Are Templates?](#what-are-templates)
    -   [Using a Template](#using-a-template)
    -   [When to Use Templates](#when-to-use-templates)
-   [Configuration Options](#-configuration-options)
    -   [Top-Level Properties](#top-level-properties)
    -   [Defaults Properties](#defaults-properties)
    -   [Windows Properties](#windows-properties)
    -   [Panes Properties](#panes-properties)
-   [Example Configuration Files](#-example-configuration-files)
    -   [Minimal Example](#minimal-example)
    -   [Advanced Example](#advanced-example)
-   [Contribution Guide](#-contribution-guide)

## 📦 Installation

### Prerequisites

-   `Go` version 1.23 or higher
-   `tmux` installed on your system

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

## 🚀 Running the Application

To start the app, navigate to the directory containing your `tmux.conf.yml` file and run:

```bash
  tmux-setup
```

This will:

-   Parse the `tmux.conf.yml` file.
-   Create a `tmux` session based on the configuration.
-   Attach you to the session if no arguments are provided.

## 🧙‍♂️ Using the Configuration Wizard

The application includes an interactive wizard to help you create a configuration file or template.

### Why Use the Wizard?

The wizard simplifies the process of creating `tmux` configurations by guiding you through each step interactively. This is especially useful for users who are not familiar with YAML syntax or the specific configuration options available.

### Creating a Configuration File

To start the wizard for creating a configuration file, run:

```bash
tmux-setup wizard
```

Follow the prompts to configure your session, windows, and panes. The wizard will save the configuration to `tmux.conf.yml` by default.

### Creating a Template

To create a new template using the wizard, run:

```bash
tmux-setup wizard --create-template <template-name>
```

This will guide you through the process of creating a template and save it to `~/.config/tmux-setup/templates/<template-name>.yml`.

### What Are Templates?

Templates are reusable configuration files that can be used to quickly set up `tmux` sessions with predefined settings. They are useful for standardizing setups across different projects or environments.

### Using a Template

To use an existing template, run:

```bash
tmux-setup --template <template-name>
```

This will load the specified template from `~/.config/tmux-setup/templates/` and create a `tmux` session based on it.

### When to Use Templates

Use templates when you have a common setup that you want to reuse across multiple projects or environments. Templates save time and ensure consistency by providing a predefined configuration that can be easily applied.

## ⚙️ Configuration Options

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
| `pre_command`     | No       | `""`          | Command to run before the session starts.                            |
| `post_command`    | No       | `""`          | Command to run after the session ends.                               |

### `windows` Properties

| Property       | Required | Default Value | Description                                                             |
| -------------- | -------- | ------------- | ----------------------------------------------------------------------- |
| `name`         | No       | `window-N`    | Name of the window.                                                     |
| `directory`    | No       | `""`          | Directory to switch to before running any commands in the window.       |
| `layout`       | No       | `""`          | Predefined layout for panes (`even-horizontal`, `even-vertical`, etc.). |
| `git_branch`   | No       | `""`          | Git branch to check out in the window's directory.                      |
| `panes`        | No       | `[]`          | List of panes to create in the window (see below).                      |
| `pre_command`  | No       | `""`          | Command to run before the window starts.                                |
| `post_command` | No       | `""`          | Command to run after the window ends.                                   |

### `panes` Properties

| Property           | Required | Default Value | Description                                               |
| ------------------ | -------- | ------------- | --------------------------------------------------------- |
| `directory`        | No       | `""`          | Directory to switch to before running the pane's command. |
| `initial_command`  | No       | `""`          | Command to run in the pane.                               |
| `refresh_interval` | No       | `0`           | Interval in seconds to refresh the pane's command.        |
| `pre_command`      | No       | `""`          | Command to run before the pane starts.                    |
| `post_command`     | No       | `""`          | Command to run after the pane ends.                       |

## 📄 Example Configuration Files

### Minimal Example

```yaml
session_name: my_session
dependencies:
    - nvim
    - pnpm
windows:
    - name: editor
      panes:
          - initial_command: nvim .
    - name: server
      panes:
          - initial_command: pnpm run dev
```

### Advanced Example

```yaml
session_name: dev
focus_window: 1
defaults:
    directory: ./src
    initial_command: echo "Welcome to tmux!"
    pre_command: ./scripts/pre_session.sh
    post_command: ./scripts/post_session.sh
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
      pre_command: ./scripts/pre_window.sh
      post_command: ./scripts/post_window.sh
    - name: logs
      panes:
          - directory: ./logs
            initial_command: tail -f app.log
          - initial_command: htop
            refresh_interval: 10
```

## 🤝 Contribution Guide

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
