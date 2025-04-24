# Append SVG Content to HTML

- This iteration of apynd is allegedly optimzed for handling many files well. Should inspect as the current version works well and this one is much, much different...

```python
#!/usr/bin/env python3
import sys
import re
import os
import glob
import argparse

def read_svg_file(svg_file_path):
    """Read SVG content from the provided file path."""
    try:
        with open(svg_file_path, 'r') as file:
            return file.read().strip()
    except FileNotFoundError:
        print(f"Error: SVG file '{svg_file_path}' not found.")
        return None
    except Exception as e:
        print(f"Error reading SVG file {svg_file_path}: {e}")
        return None

def add_icon_to_html(html_file_path, svg_content):
    """Add SVG icon to HTML file."""
    try:
        # Create the new icon wrapper
        new_icon = '''                <div class="icon-wrapper" onclick="copySvgCode(this)">
                    {}
                </div>'''.format(svg_content)

        # Read the HTML file
        with open(html_file_path, 'r') as file:
            content = file.read()

        # Find the position to insert the new icon
        # We'll insert it before the last </div> that appears before </section>
        pattern = r'(</div>[^\n]*\n[^\n]*</section>)'
        replacement = new_icon + '\n                \\1'
        updated_content = re.sub(pattern, replacement, content, count=1)

        # Write the updated content back to the file
        with open(html_file_path, 'w') as file:
            file.write(updated_content)

        return True
    except Exception as e:
        print(f"Error updating HTML file: {e}")
        return False

def main():
    parser = argparse.ArgumentParser(description='Add SVG icons to an HTML file')
    parser.add_argument('html_file', help='Path to the HTML file to update')
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument('-f', '--files', nargs='+', help='One or more SVG files to add')
    group.add_argument('-d', '--directory', help='Directory containing SVG files to add')
    parser.add_argument('-e', '--extension', default='.svg', help='File extension to look for (default: .svg)')
    args = parser.parse_args()

    # Check if HTML file exists
    if not os.path.isfile(args.html_file):
        print(f"Error: HTML file '{args.html_file}' not found.")
        sys.exit(1)
    
    # Gather list of SVG files
    svg_files = []
    
    if args.files:
        svg_files = args.files
    elif args.directory:
        if not os.path.isdir(args.directory):
            print(f"Error: Directory '{args.directory}' not found.")
            sys.exit(1)
        # Get all SVG files in the directory
        svg_files = glob.glob(os.path.join(args.directory, f'*{args.extension}'))
    
    if not svg_files:
        print("No SVG files found to process.")
        sys.exit(1)
    
    # Process each SVG file
    success_count = 0
    for svg_file in svg_files:
        svg_content = read_svg_file(svg_file)
        if svg_content:
            if add_icon_to_html(args.html_file, svg_content):
                success_count += 1
                print(f"Added icon from {svg_file} to {args.html_file}")
    
    print(f"\nSummary: Added {success_count} out of {len(svg_files)} icons to {args.html_file}")

if __name__ == "__main__":
    main()
```

### How to use this enhanced script:

1. Save this script as `add_icons.py`
2. Make it executable: `chmod +x add_icons.py`
3. Use it in one of these ways:

#### Add multiple individual SVG files:
```bash
./add_icons.py index.html -f icon1.svg icon2.svg icon3.svg
```

#### Add all SVG files from a directory:
```bash
./add_icons.py index.html -d ./icons/
```

#### Add files with a specific extension from a directory:
```bash
./add_icons.py index.html -d ./icons/ -e .icon.svg
```

### Key improvements:

1. **Command-line argument parsing**: Uses `argparse` for better user interface and help messages
2. **Multiple file handling**: Processes individual files or all matching files in a directory
3. **Batch processing**: Adds all icons in a single run
4. **Custom extension support**: Can filter files by a specific extension
5. **Summary reporting**: Shows how many files were successfully processed
6. **Error handling**: Continues processing remaining files if one fails
7. **Proper feedback**: Shows progress and results of operations

This script should efficiently handle large batches of SVG icons while providing good feedback on the process.