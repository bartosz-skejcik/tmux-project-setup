use std::fmt;

use anyhow::Error;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize)]
pub struct Config {
    pub session_name: String,
    pub focus_window: Option<usize>,
    pub defaults: Option<Defaults>,
    pub windows: Vec<Window>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Defaults {
    pub directory: Option<String>,
    pub initial_command: Option<String>,
    pub dependencies: Option<Vec<String>>,
    pub pre_command: Option<String>,
    pub post_command: Option<String>,
    pub branch: Option<String>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Window {
    pub name: String,
    pub directory: Option<String>,
    pub initial_command: Option<String>,
    pub layout: Option<Layout>,
    pub git_branch: Option<String>,
    pub panes: Option<Vec<Pane>>,
    pub pre_command: Option<String>,
    pub post_command: Option<String>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Pane {
    pub directory: Option<String>,
    pub initial_command: String,
    pub pre_command: Option<String>,
    pub post_command: Option<String>,
}

#[derive(Debug, Deserialize, Serialize)]
pub enum Layout {
    EvenHorizontal,
    EvenVertical,
    MainHorizontal,
    MainVertical,
    Tiled,
}

impl fmt::Display for Layout {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Layout::EvenHorizontal => write!(f, "even-horizontal"),
            Layout::EvenVertical => write!(f, "even-vertical"),
            Layout::MainHorizontal => write!(f, "main-horizontal"),
            Layout::MainVertical => write!(f, "main-vertical"),
            Layout::Tiled => write!(f, "tiled"),
        }
    }
}

impl Config {
    pub fn new() -> Config {
        Config {
            session_name: String::from(""),
            focus_window: None,
            defaults: None,
            windows: vec![],
        }
    }

    pub fn write_to_file(filepath: &str, content: &str) -> Result<(), Error> {
        std::fs::write(filepath, content)?;

        Ok(())
    }

    pub fn read_from_file(filepath: &str) -> Result<String, Error> {
        let content = std::fs::read_to_string(filepath)?;

        Ok(content)
    }

    pub fn load(filepath: &str) -> Result<Self, Error> {
        let content =
            Self::read_from_file(filepath).expect("No config found in the given directory");

        let config: Config = serde_yaml::from_str(&content)?;

        Ok(config)
    }

    pub fn save(&self, filepath: &str) -> Result<(), Error> {
        let content = serde_yaml::to_string(self)?;

        Self::write_to_file(filepath, &content)?;

        Ok(())
    }
}
