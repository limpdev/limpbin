import os
from pathlib import Path
import html2text
from bs4 import BeautifulSoup
import re
from urllib.parse import urljoin, urlparse

class DocConverter:
    def __init__(self, root_dir):
        self.root_dir = Path(root_dir)
        self.converter = html2text.HTML2Text()
        self.converter.body_width = 0  # Disable text wrapping
        self.converter.ignore_images = False
        self.converter.ignore_links = False
        
    def fix_relative_links(self, content, file_path):
        """
        Adjust relative links to maintain structure when converting to .md
        """
        soup = BeautifulSoup(content, 'html.parser')
        
        # Fix links
        for a in soup.find_all('a', href=True):
            href = a['href']
            if not urlparse(href).netloc:  # If it's a relative link
                # Convert .html extensions to .md
                if href.endswith('.html'):
                    href = href[:-5] + '.md'
                a['href'] = href
                
        # Fix image sources
        for img in soup.find_all('img', src=True):
            src = img['src']
            if not urlparse(src).netloc:
                # Make image paths relative to new location
                img['src'] = os.path.relpath(
                    os.path.join(os.path.dirname(file_path), src),
                    os.path.dirname(file_path)
                )
                
        return str(soup)
    
    def convert_file(self, html_path):
        """
        Convert a single HTML file to Markdown
        """
        with open(html_path, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # Fix relative links before conversion
        content = self.fix_relative_links(content, html_path)
        
        # Convert to markdown
        markdown = self.converter.handle(content)
        
        # Create corresponding markdown file path
        md_path = html_path.with_suffix('.md')
        
        # Ensure the directory exists
        md_path.parent.mkdir(parents=True, exist_ok=True)
        
        # Write the markdown file
        with open(md_path, 'w', encoding='utf-8') as f:
            f.write(markdown)
            
        return md_path
    
    def convert_directory(self):
        """
        Convert all HTML files in the directory structure
        """
        converted_files = []
        
        for html_file in self.root_dir.rglob('*.html'):
            try:
                md_file = self.convert_file(html_file)
                converted_files.append((html_file, md_file))
                print(f"Converted: {html_file} -> {md_file}")
            except Exception as e:
                print(f"Error converting {html_file}: {str(e)}")
                
        return converted_files

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Convert HTML documentation to Markdown')
    parser.add_argument('root_dir', help='Root directory of the HTML documentation')
    parser.add_argument('--delete-html', action='store_true', 
                      help='Delete original HTML files after conversion')
    
    args = parser.parse_args()
    
    converter = DocConverter(args.root_dir)
    converted_files = converter.convert_directory()
    
    print(f"\nConverted {len(converted_files)} files")
    
    if args.delete_html:
        for html_file, _ in converted_files:
            html_file.unlink()
            print(f"Deleted original file: {html_file}")

if __name__ == '__main__':
    main()