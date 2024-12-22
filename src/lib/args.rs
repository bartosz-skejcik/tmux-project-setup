use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "tmux-setup")]
#[command(about = "A tool to set up tmux sessions", long_about = None)]
pub struct TmuxSetup {
    /// Use an existing template
    #[arg(short, long, value_name = "template-name")]
    pub template: Option<String>,

    #[command(subcommand)]
    pub command: Option<Commands>, // Keep this as an Option to allow running without args
}

#[derive(Subcommand)]
pub enum Commands {
    /// Run the setup wizard
    Wizard {
        /// Create a new template with the given name
        #[arg(long, value_name = "template-name")]
        create_template: Option<String>,
    },
}
