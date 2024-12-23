use anyhow::Error;

use super::{
    config::{Defaults, Layout, Pane, Window},
    hooks::{
        look_path, resolve_dir, run_command, run_post_window_hooks, run_pre_window_hooks, send_keys,
    },
};

pub fn check_deps(deps: &Option<Vec<String>>) -> Result<(), Error> {
    if let Some(deps) = deps {
        for dep in deps.iter() {
            if look_path(dep).is_none() {
                return Err(Error::msg(format!("Dependency not found: {}", dep)));
            }
        }
    }

    Ok(())
}

pub fn create_session(sess_name: &str, config: &crate::Config) -> Result<(), Error> {
    let mut cmd = format!("tmux new-session -d -s {} -n placeholder", sess_name);
    run_command(&cmd)?;

    cmd = "tmux set-option -g base-index 1".to_string();
    run_command(&cmd)?;

    cmd = "tmux set-option -g pane-base-index 1".to_string();
    run_command(&cmd)?;

    config.windows.iter().enumerate().for_each(|(i, window)| {
        match create_window(sess_name, i, window, &config.defaults) {
            Ok(_) => (),
            Err(e) => eprintln!("Error creating window: {}", e),
        }
    });

    Ok(())
}

pub fn create_window(
    sess_name: &str,
    i: usize,
    window: &Window,
    defaults: &Option<Defaults>,
) -> Result<(), Error> {
    let mut cmd = String::new();
    let mut window_name = window.name.clone();

    if window_name.is_empty() {
        window_name = format!("window-{}", i + 1);
    }

    run_pre_window_hooks(window)?;

    if i == 0 {
        cmd = format!("tmux rename-window -t {}:1 {}", sess_name, window_name);
        run_command(&cmd)?;
    } else {
        cmd = format!("tmux new-window -t {sess_name}:{} -n {window_name}", i + 1);
        run_command(&cmd)?;
    }

    let mut resolved_dir = String::new();
    // set working dir
    if let Some(defaults) = defaults {
        if let Some(defaults_directory) = &defaults.directory {
            if let Some(window_directory) = &window.directory {
                resolved_dir = resolve_dir(defaults_directory, &window_directory.to_string());
                send_keys(sess_name, i + 1, 1, &format!("cd {}", resolved_dir))?;
            }
        }
    }

    // handle git branch
    if let Some(git_branch) = &window.git_branch {
        send_keys(sess_name, i + 1, 1, &format!("git checkout {}", git_branch))?;
    }

    // create panes
    // createPanes(sessionName, index+1, window, dir)
    create_panes(sess_name, i + 1, window, resolved_dir)?;

    // run post window hooks
    run_post_window_hooks(window)?;

    Ok(())
}

pub fn create_panes(
    sess_name: &str,
    window_idx: usize,
    window: &Window,
    default_dir: String,
) -> Result<(), Error> {
    // check if the window.panes (Option<Vec<Pane>>) is empty
    if let Some(panes) = &window.panes {
        let mut split_type = "-h";
        for (i, pane) in panes.iter().enumerate() {
            if i > 0 {
                // if the layout is present in the window and the layout has "vertical" set
                // split_type to "-v"
                if let Some(layout) = &window.layout {
                    if layout.to_string().contains("vertical") {
                        split_type = "-v";
                    }
                }
                let cmd = format!("tmux split-window {split_type} -t {sess_name}:{window_idx}");
                run_command(&cmd)?;
            }

            if let Some(pane_dir) = &pane.directory {
                let resolved_dir = resolve_dir(&default_dir, pane_dir);
                send_keys(
                    sess_name,
                    window_idx,
                    i + 1,
                    &format!("cd {}", resolved_dir),
                )?;
            }

            send_keys(sess_name, window_idx, i + 1, &pane.initial_command)?;
        }
    }

    if let Some(window_layout) = &window.layout {
        let cmd = format!(
            "tmux select-layout -t {sess_name}:{window_idx} {layout}",
            sess_name = sess_name,
            window_idx = window_idx,
            layout = window_layout
        );
        run_command(&cmd)?;
    }

    Ok(())
}
