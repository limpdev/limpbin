use std::fs;
use std::path::{Path, PathBuf};
use std::io;

fn move_file(source: &Path, destination: &Path) -> io::Result<()> {
    fs::rename(source, destination)
}

fn main() {
    let args: Vec<String> = std::env::args().collect();

    if args.len() != 3 {
        eprintln!("Usage: file_mover.exe <destination> <source>");
        return;
    }

    let destination_path = PathBuf::from(args[1].as_str());
    let source_path = PathBuf::from(args[2].as_str());

    match move_file(&source_path, &destination_path) {
        Ok(_) => println!("File moved successfully."),
        Err(e) => eprintln!(
            "Failed to move file from {} to {}: {}",
            source_path.display(),
            destination_path.display(),
            e
        ),
    }
}


