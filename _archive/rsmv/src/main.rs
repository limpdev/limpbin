use clap::Parser;
use std::fs;
use std::io;
use std::path::{Path, PathBuf};
use std::time::Instant;

/// High-performance file movement utility
#[derive(Parser)]
#[clap(author, version, about)]
struct Args {
    /// Source file or directory
    source: String,

    /// Destination file or directory
    destination: String,

    /// Delete source files after copying (default is to preserve source files)
    #[clap(short, long)]
    delete: bool,
}

fn main() -> io::Result<()> {
    let args = Args::parse();
    let source = &args.source;
    let destination = &args.destination;
    let delete_source = args.delete;

    let start = Instant::now();
    let (files_processed, bytes_processed) = process_files(source, destination, delete_source)?;
    let duration = start.elapsed();

    // Convert bytes to appropriate unit
    let size_str = if bytes_processed >= 1_000_000_000 {
        format!("{:.2} GB", bytes_processed as f64 / 1_000_000_000.0)
    } else if bytes_processed >= 1_000_000 {
        format!("{:.2} MB", bytes_processed as f64 / 1_000_000.0)
    } else if bytes_processed >= 1_000 {
        format!("{:.2} KB", bytes_processed as f64 / 1_000.0)
    } else {
        format!("{} bytes", bytes_processed)
    };

    // Calculate speed
    let speed_bytes_per_sec = if duration.as_secs() > 0 || duration.subsec_nanos() > 0 {
        bytes_processed as f64 / duration.as_secs_f64()
    } else {
        f64::INFINITY
    };

    let speed_str = if speed_bytes_per_sec >= 1_000_000_000.0 {
        format!("{:.2} GB/s", speed_bytes_per_sec / 1_000_000_000.0)
    } else if speed_bytes_per_sec >= 1_000_000.0 {
        format!("{:.2} MB/s", speed_bytes_per_sec / 1_000_000.0)
    } else if speed_bytes_per_sec >= 1_000.0 {
        format!("{:.2} KB/s", speed_bytes_per_sec / 1_000.0)
    } else {
        format!("{:.2} bytes/s", speed_bytes_per_sec)
    };

    let operation = if delete_source { "Moved" } else { "Copied" };
    println!(
        "{} {} files ({}) in {:.2?} ({}).",
        operation, files_processed, size_str, duration, speed_str
    );

    Ok(())
}

fn process_files(source: &str, destination: &str, delete_source: bool) -> io::Result<(u64, u64)> {
    let source_path = Path::new(source);
    let destination_path = Path::new(destination);

    if !source_path.exists() {
        return Err(io::Error::new(
            io::ErrorKind::NotFound,
            format!("Source '{}' does not exist", source),
        ));
    }

    // Handle file-to-file operation
    if source_path.is_file() {
        let dest_path = if destination_path.is_dir() {
            // If destination is a directory, create a path with the source filename
            PathBuf::from(destination).join(source_path.file_name().unwrap())
        } else {
            // Otherwise use the destination as-is
            PathBuf::from(destination)
        };

        let bytes = source_path.metadata()?.len();
        
        if delete_source {
            // Try to use rename for speed (works on same filesystem)
            match fs::rename(source_path, &dest_path) {
                Ok(_) => return Ok((1, bytes)),
                Err(e) => {
                    // Check if error is likely due to cross-device operation
                    // (Workaround for not using unstable CrossesDevices error kind)
                    let err_string = e.to_string().to_lowercase();
                    let is_cross_device = err_string.contains("cross-device") || 
                                         err_string.contains("different drive") ||
                                         err_string.contains("not same device");
                                         
                    if is_cross_device {
                        // Fall back to copy + remove if across devices
                        fs::copy(source_path, &dest_path)?;
                        fs::remove_file(source_path)?;
                        return Ok((1, bytes));
                    } else {
                        return Err(e);
                    }
                }
            }
        } else {
            // Just copy without deleting
            fs::copy(source_path, &dest_path)?;
            return Ok((1, bytes));
        }
    }

    // Handle directory-to-directory operation
    if source_path.is_dir() {
        // Create destination directory if it doesn't exist
        if !destination_path.exists() {
            fs::create_dir_all(destination_path)?;
        } else if !destination_path.is_dir() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                format!("Destination '{}' is not a directory", destination),
            ));
        }

        let mut files_processed = 0;
        let mut bytes_processed = 0;

        // Process all entries in the source directory
        let entries = fs::read_dir(source_path)?;
        for entry in entries {
            let entry = entry?;
            let file_type = entry.file_type()?;
            let path = entry.path();
            
            // Process files and directories differently
            if file_type.is_file() {
                let file_name = path.file_name().unwrap();
                let dest_file = PathBuf::from(destination).join(file_name);
                
                let file_size = path.metadata()?.len();
                
                if delete_source {
                    // Try rename first (faster when on same filesystem)
                    match fs::rename(&path, &dest_file) {
                        Ok(_) => {
                            files_processed += 1;
                            bytes_processed += file_size;
                        },
                        Err(e) => {
                            // Check if error is likely due to cross-device operation
                            let err_string = e.to_string().to_lowercase();
                            let is_cross_device = err_string.contains("cross-device") || 
                                                 err_string.contains("different drive") ||
                                                 err_string.contains("not same device");
                                                 
                            if is_cross_device {
                                // Use copy + delete for cross-device moves
                                fs::copy(&path, &dest_file)?;
                                fs::remove_file(&path)?;
                                files_processed += 1;
                                bytes_processed += file_size;
                            } else {
                                return Err(e);
                            }
                        }
                    }
                } else {
                    // Just copy without deleting
                    fs::copy(&path, &dest_file)?;
                    files_processed += 1;
                    bytes_processed += file_size;
                }
            } else if file_type.is_dir() {
                // Recursively process subdirectories
                let dir_name = path.file_name().unwrap();
                let dest_dir = PathBuf::from(destination).join(dir_name);
                
                // Create the destination subdirectory
                if !dest_dir.exists() {
                    fs::create_dir_all(&dest_dir)?;
                }
                
                // Recursively process files from the subdirectory
                let (sub_files, sub_bytes) = process_files(
                    path.to_str().unwrap(),
                    dest_dir.to_str().unwrap(),
                    delete_source,
                )?;
                
                files_processed += sub_files;
                bytes_processed += sub_bytes;
                
                // Remove the now-empty source directory if deletion is enabled
                if delete_source && path.read_dir()?.next().is_none() {
                    fs::remove_dir(path)?;
                }
            }
        }
        
        // Remove the source directory if it's now empty and deletion is enabled
        if delete_source && source_path.read_dir()?.next().is_none() {
            fs::remove_dir(source_path)?;
        }
        
        return Ok((files_processed, bytes_processed));
    }

    Err(io::Error::new(
        io::ErrorKind::InvalidInput,
        format!("Unsupported operation for '{}' to '{}'", source, destination),
    ))
}