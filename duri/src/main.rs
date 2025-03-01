use std::env;
use std::fs;
use std::io::Read;
use std::path::Path;
use base64::{Engine as _, engine::general_purpose::STANDARD};
use image;
use clipboard::{ClipboardContext, ClipboardProvider};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Get command line arguments
    let args: Vec<String> = env::args().collect();
    
    // Get the program name for error messages
    let program_name = env::args().next().unwrap_or_else(|| "program".to_string());
    
    // Check if an argument was provided
    if args.len() < 2 {
        eprintln!("Usage: {} <image_file_path>", program_name);
        std::process::exit(1);
    }
    
    let file_path = &args[1];
    let path = Path::new(file_path);
    
    // Check if file exists
    if !path.exists() {
        eprintln!("Error: File '{}' does not exist", file_path);
        std::process::exit(1);
    }
    
    // Read the file
    let mut file = fs::File::open(path)?;
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer)?;
    
    // Check file extension for ICNS specifically
    let is_icns = match path.extension() {
        Some(ext) => ext.to_string_lossy().to_lowercase() == "icns",
        None => false,
    };
    
    let mime_type = if is_icns {
        "image/icns"
    } else {
        // For other formats, use the image crate
        let format = match image::guess_format(&buffer) {
            Ok(format) => format,
            Err(_) => {
                if is_icns {
                    // If we already know it's ICNS, continue
                    image::ImageFormat::Png // Placeholder, won't be used
                } else {
                    eprintln!("Error: File '{}' is not a valid image", file_path);
                    std::process::exit(1);
                }
            }
        };
        
        // Convert format to MIME type
        match format {
            image::ImageFormat::Jpeg => "image/jpeg",
            image::ImageFormat::Png => "image/png",
            image::ImageFormat::WebP => "image/webp",
            image::ImageFormat::Ico => "image/x-icon",
            _ => "application/octet-stream", // Generic for other formats
        }
    };
    
    // Convert to base64
    let base64_string = STANDARD.encode(&buffer);
    
    // Create the complete data URI
    let data_uri = format!("data:{};base64,{}", mime_type, base64_string);
    
    // Copy to clipboard
    let mut ctx: ClipboardContext = ClipboardProvider::new()?;
    ctx.set_contents(data_uri.clone())?;
    
    println!("Successfully converted image to data URI and copied to clipboard!");
    println!("Image size: {} bytes, Data URI size: {} bytes", buffer.len(), data_uri.len());
    println!("MIME type detected: {}", mime_type);
    
    Ok(())
}