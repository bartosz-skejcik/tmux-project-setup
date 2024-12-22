use std::env;
use std::path::Path;
use std::process::Command;

use anyhow::Error;

use super::config::Window;

pub fn run_command(command: &str) -> Result<(), Error> {
    let output = std::process::Command::new("sh")
        .arg("-c")
        .arg(command)
        .output()?;

    if output.status.success() {
        Ok(())
    } else {
        Err(Error::msg(format!(
            "Command failed with exit code {}: {}",
            output.status.code().unwrap(),
            String::from_utf8_lossy(&output.stderr)
        )))
    }
}

pub fn run_pre_session_hooks(config: &crate::Config) -> Result<(), Error> {
    match &config.defaults {
        Some(defaults) => {
            if let Some(pre_cmd) = &defaults.pre_command {
                run_command(pre_cmd)?;
            }

            Ok(())
        }
        None => Ok(()),
    }
}

pub fn run_post_session_hooks(config: &crate::Config) -> Result<(), Error> {
    match &config.defaults {
        Some(defaults) => {
            if let Some(post_cmd) = &defaults.post_command {
                run_command(post_cmd)?;
            }

            Ok(())
        }
        None => Ok(()),
    }
}

pub fn run_pre_window_hooks(window: &Window) -> Result<(), Error> {
    if let Some(pre_cmd) = &window.pre_command {
        run_command(pre_cmd)?;
    }

    Ok(())
}

pub fn run_post_window_hooks(window: &Window) -> Result<(), Error> {
    if let Some(post_cmd) = &window.post_command {
        run_command(post_cmd)?;
    }

    Ok(())
}

pub fn look_path(exec: &str) -> Option<String> {
    if Path::new(exec).is_absolute() {
        // if the path is absolute, check if it exists
        if Path::new(exec).exists() {
            return Some(exec.to_string());
        } else {
            return None;
        }
    }

    // get the PATH env variable
    if let Ok(path) = env::var("PATH") {
        for dir in path.split(":") {
            let full_path = Path::new(dir).join(exec);
            if full_path.exists() {
                return Some(full_path.to_string_lossy().into_owned());
            }
        }
    }

    None
}

pub fn resolve_dir(parent: &str, dir: &str) -> String {
    if Path::new(dir).is_absolute() {
        return dir.to_string();
    }

    let parent_path = Path::new(parent);
    let full_path = parent_path.join(dir);

    full_path.to_string_lossy().into_owned()
}

pub fn send_keys(
    sess_name: &str,
    window_idx: usize,
    pane_idx: usize,
    command: &str,
) -> Result<(), Error> {
    let cmd = format!(
        "tmux send-keys -t {sess_name}:{window_idx}.{pane_idx} {command} C-m",
        sess_name = sess_name,
        window_idx = window_idx + 1,
        pane_idx = pane_idx + 1,
        command = command
    );

    run_command(&cmd)
}
