import markdown2
import os
import sys

def md_to_html(markdown_file, output_file=None):
    """
    Convert a markdown file to HTML using the markdown2 library and a template.
    
    Args:
        markdown_file (str): Path to the markdown file
        output_file (str, optional): Path to the output HTML file. If None, 
                                     will use the same name as the markdown file but with .html extension.
    """
    # Read the HTML template
    with open("index.html", "r") as f:
        html_template = f.read()
    
    # Read the markdown file
    try:
        with open(markdown_file, "r", encoding="utf-8") as f:
            markdown_content = f.read()
    except Exception as e:
        print(f"Error: Could not read markdown file: {e}")
        return False
    
    # Convert markdown to HTML
    try:
        markdowner = markdown2.Markdown(extras=[
            "fenced-code-blocks", 
            "tables", 
            "header-ids", 
            "footnotes",
            "code-friendly",
            "link-shortrefs",
            "html-classes",
            "highlightjs-lang"
        ])
        html_content = markdowner.convert(markdown_content)
    except Exception as e:
        print(f"Error: Could not convert markdown to html: {e}")
        return False
    
    # Replace the %s placeholders in the template
    # First %s is for the title (use the filename without extension)
    title = os.path.basename(markdown_file).split('.')[0]
    complete_html = html_template % (title, html_content)
    
    # Determine output file name if not provided
    if output_file is None:
        output_file = os.path.splitext(markdown_file)[0] + ".html"
    
    # Write the HTML to file
    try:
        with open(output_file, "w", encoding="utf-8") as f:
            f.write(complete_html)
        print(f"Successfully converted {markdown_file} to {output_file}")
        return True
    except Exception as e:
        print(f"Error: Could not write to output file: {e}")
        return False

if __name__ == "__main__":
    # Check if a markdown file was provided
    if len(sys.argv) < 2:
        print("Usage: python pymo.py <markdown_file> [output_html_file]")
        sys.exit(1)
    
    markdown_file = sys.argv[1]
    
    # Check if output file was provided
    output_file = sys.argv[2] if len(sys.argv) > 2 else None
    
    # Convert markdown to HTML
    success = md_to_html(markdown_file, output_file)
    
    # Exit with appropriate code
    sys.exit(0 if success else 1)