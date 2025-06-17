#!/usr/bin/env python3
import sys
import re
import os
from bs4 import BeautifulSoup

def read_svg_file(svg_file_path):
    """Read SVG content from the provided file path."""
    try:
        with open(svg_file_path, 'r') as file:
            return file.read().strip()
    except FileNotFoundError:
        print(f"Error: SVG file '{svg_file_path}' not found.")
        sys.exit(1)
    except Exception as e:
        print(f"Error reading SVG file: {e}")
        sys.exit(1)

def add_icon_to_html(html_file_path, svg_content):
    """Add SVG icon to HTML file."""
    try:
        # Create the new icon wrapper
        new_icon = '''<div class="icon-wrapper" onclick="copySvgCode(this)">{}</div'''.format(svg_content)

        # Read the HTML file
        with open(html_file_path, 'r') as file:
            content = file.read()

        # Find the position to insert the new icon
        # We'll insert it before the last </div> that appears before </section>
        pattern = r'(<\/div>[^\n]*\n[^\n]*<\/section>)'
        replacement = new_icon + '\n                \\1'
        updated_content = re.sub(pattern, replacement, content, count=1)
        
        # Format the HTML using BeautifulSoup
        soup = BeautifulSoup(updated_content, 'html.parser')
        formatted_html = soup.prettify()
        
        # Write the formatted content back to the file
        with open(html_file_path, 'w') as file:
            file.write(formatted_html)

        print(f"New icon added successfully to {html_file_path}")
        return True
    except Exception as e:
        print(f"Error updating HTML file: {e}")
        return False

def main():
    if len(sys.argv) != 3:
        print("Usage: ./add_icon.py <html_file> <svg_file>")
        sys.exit(1)

    html_file = sys.argv[1]
    svg_file = sys.argv[2]

    # Check if files exist
    if not os.path.isfile(html_file):
        print(f"Error: HTML file '{html_file}' not found.")
        sys.exit(1)

    # Read SVG content from file
    svg_content = read_svg_file(svg_file)

    # Add the icon to the HTML file
    add_icon_to_html(html_file, svg_content)

if __name__ == "__main__":
    main()