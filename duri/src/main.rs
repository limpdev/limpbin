use base64::{engine::general_purpose::STANDARD, Engine};
use clap::Parser;
use clipboard::{ClipboardContext, ClipboardProvider};
use mime_guess::from_path;
use std::fs::File;
use std::io::Read;
use std::path::PathBuf;

#[derive(Parser)]
#[command(name = "image-to-datauri")]
#[command(about = "Converts an image file to a DataURI and copies it to clipboard")]
struct Cli {
    /// Path to the image file
    #[arg(value_name = "FILE")]
    path: PathBuf,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Cli::parse();
    
    // Read the file
    let mut file = File::open(&args.path)?;
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer)?;
    
    // Guess MIME type from file extension
    let mime_type = from_path(&args.path)
        .first_or_octet_stream()
        .to_string();
    
    // Encode to base64
    let base64_string = STANDARD.encode(&buffer);
    
    // Create DataURI
    let data_uri = format!("data:{};base64,{}", mime_type, base64_string);
    
    // Copy to clipboard
    let mut ctx: ClipboardContext = ClipboardProvider::new()?;
    ctx.set_contents(data_uri.clone())?;
    
    println!("Successfully copied DataURI to clipboard!");
    println!("Length: {} characters", data_uri.len());
    
    Ok(())
}