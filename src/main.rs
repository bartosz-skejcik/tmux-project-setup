use std::env::current_dir;

use anyhow::Error;
use clap::Parser;
use lib::{
    args::{Commands, TmuxSetup},
    config::Config,
    hooks::{
        attach_session, get_template_path, run_post_session_hooks, run_pre_session_hooks,
        session_exists,
    },
    session::{check_deps, create_session},
};

mod lib;

fn main() -> Result<(), Error> {
    let args = TmuxSetup::parse();

    // Check for the top-level template option
    if let Some(template_name) = args.template {
        let filepath = get_template_path(&template_name)?;

        start(&filepath)?;

        return Ok(());
    }

    match args.command {
        Some(Commands::Wizard { create_template }) => {
            if let Some(template_name) = create_template {
                println!("Creating template: {}", template_name);
                // Add logic for creating a template here

                // run the setup wizard here
            } else {
                println!("Creating a new setup wizard and saving to pwd...");

                // run the setup wizard here
            }
        }
        None => {
            println!("No command provided, starting with default tmux.conf.yml file...");
            let filepath = "tmux.conf.yml".to_string();
            start(&filepath)?;
        }
    }

    Ok(())
}

fn start(file_path: &str) -> Result<(), Error> {
    let config = Config::load(file_path)?;

    if session_exists(&config.session_name) {
        attach_session(&config.session_name)?;
    }

    if let Some(defaults) = &config.defaults {
        check_deps(&defaults.dependencies)?;
    }

    let sess_name = match config.session_name.is_empty() {
        true => current_dir()?
            .file_name()
            .unwrap()
            .to_str()
            .unwrap()
            .to_string(),
        false => config.session_name.clone(),
    };

    if config.defaults.is_some() {
        run_pre_session_hooks(&config)?;
    }

    create_session(&sess_name, &config)?;

    if config.defaults.is_some() {
        run_post_session_hooks(&config)?;
    }

    println!("{:#?}", config);

    Ok(())
}
