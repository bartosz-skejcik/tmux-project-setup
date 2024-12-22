use std::{default, env::current_dir};

use anyhow::Error;
use clap::Parser;
use lib::{
    args::{Commands, TmuxSetup},
    config::Config,
    session::check_deps,
};

mod lib;

fn main() -> Result<(), Error> {
    let args = TmuxSetup::parse();

    // Check for the top-level template option
    if let Some(template_name) = args.template {
        println!("Using template: {}", template_name);
        // Add your logic for using an existing template here
    } else {
        match args.command {
            Some(Commands::Wizard { create_template }) => {
                if let Some(template_name) = create_template {
                    println!("Creating template: {}", template_name);
                    // Add your logic for creating a template here

                    // run the setup wizard here
                } else {
                    println!("Creating a new setup wizard and saving to pwd...");

                    // run the setup wizard here
                }
            }
            None => {
                let current_directory = current_dir()?;
                let file_path = current_directory.join("tmux.conf.yaml");
                let file_path = file_path.to_str().unwrap();
                let config = Config::load(file_path)?;

                if let Some(defaults) = &config.defaults {
                    // will throw for `asdf` 😎
                    check_deps(&defaults.dependencies)?;
                }

                println!("{:#?}", config);
            }
        }
    }

    Ok(())
}
