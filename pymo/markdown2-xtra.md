usage: markdown2 [PATHS...]

A fast and complete Python implementation of Markdown, a
text-to-HTML conversion tool for web writers.

Supported extra syntax options (see -x|--extras option below and
see <https://github.com/trentm/python-markdown2/wiki/Extras> for details):

* admonitions: Enable parsing of RST admonitions.
* breaks: Control where hard breaks are inserted in the markdown.
  Options include:
  - on_newline: Replace single new line characters with <br> when True
  - on_backslash: Replace backslashes at the end of a line with <br>
* break-on-newline: Alias for the on_newline option in the breaks extra.
* code-friendly: Disable _ and __ for em and strong.
* cuddled-lists: Allow lists to be cuddled to the preceding paragraph.
* fenced-code-blocks: Allows a code block to not have to be indented
  by fencing it with '```' on a line before and after. Based on
  <http://github.github.com/github-flavored-markdown/> with support for
  syntax highlighting.
* footnotes: Support footnotes as in use on daringfireball.net and
  implemented in other Markdown processors (tho not in Markdown.pl v1.0.1).
* header-ids: Adds "id" attributes to headers. The id value is a slug of
  the header text.
* highlightjs-lang: Allows specifying the language which used for syntax
  highlighting when using fenced-code-blocks and highlightjs.
* html-classes: Takes a dict mapping html tag names (lowercase) to a
  string to use for a "class" tag attribute. Currently only supports "img",
  "table", "thead", "pre", "code", "ul" and "ol" tags. Add an issue if you require
  this for other tags.
* link-patterns: Auto-link given regex patterns in text (e.g. bug number
  references, revision number references).
* link-shortrefs: allow shortcut reference links, not followed by `[]` or
  a link label.
* markdown-in-html: Allow the use of `markdown="1"` in a block HTML tag to
  have markdown processing be done on its contents. Similar to
  <http://michelf.com/projects/php-markdown/extra/#markdown-attr> but with
  some limitations.
* metadata: Extract metadata from a leading '---'-fenced block.
  See <https://github.com/trentm/python-markdown2/issues/77> for details.
* middle-word-em: Allows or disallows emphasis syntax in the middle of words,
  defaulting to allow. Disabling this means that `this_text_here` will not be
  converted to `this<em>text</em>here`.
* nofollow: Add `rel="nofollow"` to add `<a>` tags with an href. See
  <http://en.wikipedia.org/wiki/Nofollow>.
* numbering: Support of generic counters.  Non standard extension to
  allow sequential numbering of figures, tables, equations, exhibits etc.
* pyshell: Treats unindented Python interactive shell sessions as <code>
  blocks.
* smarty-pants: Replaces ' and " with curly quotation marks or curly
  apostrophes.  Replaces --, ---, ..., and . . . with en dashes, em dashes,
  and ellipses.
* spoiler: A special kind of blockquote commonly hidden behind a
  click on SO. Syntax per <http://meta.stackexchange.com/a/72878>.
* strike: text inside of double tilde is ~~strikethrough~~
* tag-friendly: Requires atx style headers to have a space between the # and
  the header text. Useful for applications that require twitter style tags to
  pass through the parser.
* tables: Tables using the same format as GFM
  <https://help.github.com/articles/github-flavored-markdown#tables> and
  PHP-Markdown Extra <https://michelf.ca/projects/php-markdown/extra/#table>.
* toc: The returned HTML string gets a new "toc_html" attribute which is
  a Table of Contents for the document. (experimental)
* use-file-vars: Look for an Emacs-style markdown-extras file variable to turn
  on Extras.
* wiki-tables: Google Code Wiki-style tables. See
  <http://code.google.com/p/support/wiki/WikiSyntax#Tables>.
* wavedrom: Support for generating Wavedrom digital timing diagrams
* xml: Passes one-liner processing instructions and namespaced XML tags.

positional arguments:
  paths                 optional list of files to convert.If none are given,
                        stdin will be used

options:
  -h, --help            show this help message and exit
  --version             show program's version number and exit
  -v, --verbose         more verbose output
  --encoding ENCODING   specify encoding of text content
  --html4tags           use HTML 4 style for empty element tags
  -s MODE, --safe MODE  sanitize literal HTML: 'escape' escapes HTML meta
                        chars, 'replace' replaces with an [HTML_REMOVED] note
  -x EXTRAS, --extras EXTRAS
                        Turn on specific extra features (not part of the core
                        Markdown spec). See above.
  --use-file-vars USE_FILE_VARS
                        Look for and use Emacs-style 'markdown-extras' file
                        var to turn on extras. See
                        <https://github.com/trentm/python-
                        markdown2/wiki/Extras>
  --link-patterns-file LINK_PATTERNS_FILE
                        path to a link pattern file
  --self-test           run internal self-tests (some doctests)
  --compare             run against Markdown.pl as well (for testing)
  --output OUTPUT       output to a file instead of stdout
