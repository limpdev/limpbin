#!/usr/bin/env python3
# --- START OF FILE mdit.py ---

import argparse
import sys
from pathlib import Path
from markdown_it import MarkdownIt
from markdown_it.token import Token
from markdown_it.utils import OptionsDict
from typing import Optional # Removed List, Sequence as they weren't strictly needed for type hints here

# --- Markdown-it Setup ---
# Using gfm-like preset and enabling specific rules/plugins
# Note: 'replacements' and 'smartquotes' are often part of the 'typographer' setting.
# Enabling them individually like this is fine, but using `md.enable('typographer')`
# might be simpler if you want all typographic features.
# Also, 'custom_replacements' is not a built-in rule; we add it via the plugin below.
md = MarkdownIt("gfm-like").enable(
    [
        "table",
        "strikethrough",
        "linkify",
        "replacements", # Handles things like (c) -> ©, (tm) -> ™, etc.
        "smartquotes",  # Handles smart quotes and dashes
        "html_inline",
        "html_block",
        "list",
    ]
)
# Note: gfm-like already includes most block and inline rules. Enabling them explicitly
#       like in the original example isn't strictly necessary unless you start from 'zero'.
#       The list above enables optional features on top of gfm-like.


# --- Custom Replacement Rule ---
def customReplacementsRule(state):
    """
    Core rule to perform custom text replacements (e.g., -> to →).
    Operates on the token stream after block/inline parsing.
    """
    replacements = {
        "->": "→",
        "=>": "⇒",
        "<-": "←",
        "<=": "⇐",
        "<->": "↔",
        "<=>": "⇔",
    }
    for i, token in enumerate(state.tokens):
        # Text content usually resides within 'inline' tokens.
        if token.type == "inline" and token.children:
            for child in token.children:
                # Only modify plain 'text' tokens.
                if child.type == "text":
                    content = child.content
                    for key, value in replacements.items():
                        content = content.replace(key, value)
                    child.content = content
        # This automatically avoids replacements inside code blocks ('fence')
        # or inline code ('code_inline') because their content is stored
        # directly in the token's 'content' attribute or handled differently,
        # not as 'text' children of an 'inline' token.


# --- Custom Replacement Plugin ---
def custom_replacements_plugin(md_instance: MarkdownIt, options: Optional[OptionsDict] = None):
    """
    Plugin function to add the custom replacement rule to the core ruler.
    """
    # Add our rule. 'custom_replacer' is a unique name for the rule.
    # We add it after 'inline' processing and potentially after standard 'replacements'.
    md_instance.core.ruler.push("custom_replacer", customReplacementsRule)


# --- Apply the Custom Plugin ---
# This is crucial: the plugin needs to be explicitly used by the md instance.
md.use(custom_replacements_plugin)


# --- CORE PROCESSING FUNCTION ---
def process_markdown_file(input_path: Path, output_path: Path):
    """
    Reads a Markdown file, converts it to HTML using the configured
    MarkdownIt instance, and writes the HTML to the output file.
    """
    masterHTML = Path("./master.html")
    print(f"Processing: {input_path} -> {output_path}")
    try:
        markdown_text = input_path.read_text(encoding="utf-8")
        html_text = md.render(markdown_text)
        # Wrap the output HTML with the template and inject the content to "{{CONTENT}}"
        with open(masterHTML, "r", encoding="utf-8") as f:
            template = f.read()
            html_text = template.replace("{{CONTENT}}", html_text)
        # Ensure parent directory exists for the output file
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(html_text, encoding="utf-8")
    except Exception as e:
        print(f"Error processing file {input_path}: {e}", file=sys.stderr)


# --- Main Execution Block ---
if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Convert Markdown file(s) to HTML using markdown-it-py with custom replacements.",
        epilog="Example: mit my_document.md OR python mit.py ./docs/"
    )
    parser.add_argument(
        "input_path",
        type=str,
        help="Path to a Markdown file (.md, .markdown) or a directory containing Markdown files.",
    )
    parser.add_argument(
        "-o", "--output",
        type=str,
        help="Optional: Specify an output file or directory.",
        default=None
    )
    parser.add_argument(
        "-r", "--recursive",
        action="store_true",
        help="If input_path is a directory, process Markdown files recursively."
    )

    args = parser.parse_args()

    input_path = Path(args.input_path)
    output_spec = args.output

    if not input_path.exists():
        print(f"Error: Input path '{input_path}' does not exist.", file=sys.stderr)
        sys.exit(1)

    if input_path.is_file():
        # Handle single file input
        if input_path.suffix.lower() not in [".md", ".markdown"]:
            print(f"Warning: Input file '{input_path}' does not have a typical Markdown extension (.md, .markdown). Processing anyway.", file=sys.stderr)

        if output_spec:
            output_path = Path(output_spec)
            # If output is specified as a directory, put the file inside it with original name + .html
            if output_path.suffix == "" or output_path.is_dir():
                 output_path = output_path / input_path.with_suffix(".html").name
        else:
            # Default: Output next to input file
            output_path = input_path.with_suffix(".html")

        process_markdown_file(input_path, output_path)

    elif input_path.is_dir():
        # Handle directory input
        output_base_dir = Path(output_spec) if output_spec else input_path
        if output_base_dir.exists() and not output_base_dir.is_dir():
             print(f"Error: Specified output '{output_base_dir}' exists but is not a directory.", file=sys.stderr)
             sys.exit(1)

        # Determine search pattern (recursive or not)
        search_pattern = "**/*" if args.recursive else "*"
        markdown_files_found = False

        # Iterate through potential markdown files
        for md_file in input_path.glob(f"{search_pattern}.md"):
            markdown_files_found = True
            relative_path = md_file.relative_to(input_path)
            output_path = (output_base_dir / relative_path).with_suffix(".html")
            process_markdown_file(md_file, output_path)

        for md_file in input_path.glob(f"{search_pattern}.markdown"):
             markdown_files_found = True
             # Avoid processing if it was already handled by the .md glob (e.g., file.md.markdown)
             # Though this case is unlikely, checking suffix prevents double processing.
             if md_file.suffix.lower() == ".markdown":
                 relative_path = md_file.relative_to(input_path)
                 output_path = (output_base_dir / relative_path).with_suffix(".html")
                 process_markdown_file(md_file, output_path)

        if not markdown_files_found:
             print(f"No Markdown files (.md, .markdown) found in '{input_path}' {'recursively' if args.recursive else ''}.", file=sys.stderr)


    else:
        print(f"Error: Input path '{input_path}' is neither a file nor a directory.", file=sys.stderr)
        sys.exit(1)

    print("Done.")
# --- END OF FILE mdit.py ---
