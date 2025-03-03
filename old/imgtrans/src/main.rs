use std::path::Path;

use clap::{Parser, ValueEnum};

use image::{GenericImageView, ImageError, ImageFormat, ImageOutputFormat};

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Cli {
    /// Path to the input image
    input: String,

    /// Operation to perform
    #[arg(value_enum)]
    operation: Operation,
}

#[derive(Clone, ValueEnum)]
enum Operation {
    FlipH,
    FlipV,
    RotateCw,
    RotateCcw,
}

fn main() -> Result<(), ImageError> {
    let cli = Cli::parse();

    let input_path = &cli.input;
    let operation = cli.operation;

    // Load the image
    let mut img = image::open(&input_path)?;

    // Perform the operation
    match operation {
        Operation::FlipH => {
            img = img.fliph();
        }
        Operation::FlipV => {
            img = img.flipv();
        }
        Operation::RotateCw => {
            img = img.rotate90();
        }
        Operation::RotateCcw => {
            img = img.rotate270();
        }
    }

    // Determine the image format
    let img_format = image::ImageFormat::from_path(&input_path)?;

    // Overwrite the original image
    let output_path = Path::new(&input_path);
    let mut output_file = std::fs::File::create(&output_path)?;

    // Save the image in the original format
    img.write_to(&mut output_file, img_format)?;

    println!("Image has been processed and saved to {}", input_path);

    Ok(())
}
