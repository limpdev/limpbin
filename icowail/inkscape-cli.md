Usage:
  org.inkscape.Inkscape [OPTION…] file1 [file2 [fileN]]

Process (or open) one or more files.

Help Options:
  -?, --help                                 Show help options
  --help-all                                 Show all help options
  --help-gapplication                        Show GApplication options
  --help-gtk                                 Show GTK+ Options

Application Options:
  -V, --version                              Print Inkscape version
  --debug-info                               Print debugging information
  --system-data-directory                    Print system data directory
  --user-data-directory                      Print user data directory
  --app-id-tag                               Create a unique instance of Inkscape with the application ID 'org.inkscape.Inkscape.TAG'
  -p, --pipe                                 Read input file from standard input (stdin)
  --pdf-page=PAGE                            PDF page number to import
  --pdf-poppler                              Use poppler when importing via commandline
  --convert-dpi-method=METHOD                Method used to convert pre-0.92 document dpi, if needed: [none|scale-viewbox|scale-document]
  --no-convert-text-baseline-spacing         Do not fix pre-0.92 document's text baseline spacing on opening
  -o, --export-filename=FILENAME             Output file name (defaults to input filename; file type is guessed from extension if present; use '-' to write to stdout)
  --export-overwrite                         Overwrite input file (otherwise add '_out' suffix if type doesn't change)
  --export-type=TYPE[,TYPE]*                 File type(s) to export: [svg,png,ps,eps,pdf,emf,wmf,xaml]
  --export-extension=EXTENSION-ID            Extension ID to use for exporting
  -C, --export-area-page                     Area to export is page
  -D, --export-area-drawing                  Area to export is whole drawing (ignoring page size)
  -a, --export-area=x0:y0:x1:y1              Area to export in SVG user units
  --export-area-snap                         Snap the bitmap export area outwards to the nearest integer values
  -d, --export-dpi=DPI                       Resolution for bitmaps and rasterized filters; default is 96
  -w, --export-width=WIDTH                   Bitmap width in pixels (overrides --export-dpi)
  -h, --export-height=HEIGHT                 Bitmap height in pixels (overrides --export-dpi)
  --export-margin=MARGIN                     Margin around export area: units of page size for SVG, mm for PS/PDF
  -i, --export-id=OBJECT-ID[;OBJECT-ID]*     ID(s) of object(s) to export
  -j, --export-id-only                       Hide all objects except object with ID selected by export-id
  -l, --export-plain-svg                     Remove Inkscape-specific SVG attributes/properties
  --export-ps-level=LEVEL                    Postscript level (2 or 3); default is 3
  --export-pdf-version=VERSION               PDF version (1.4 or 1.5); default is 1.5
  -T, --export-text-to-path                  Convert text to paths (PS/EPS/PDF/SVG)
  --export-latex                             Export text separately to LaTeX file (PS/EPS/PDF)
  --export-ignore-filters                    Render objects without filters instead of rasterizing (PS/EPS/PDF)
  -t, --export-use-hints                     Use stored filename and DPI hints when exporting object selected by --export-id
  -b, --export-background=COLOR              Background color for exported bitmaps (any SVG color string)
  -y, --export-background-opacity=VALUE      Background opacity for exported bitmaps (0.0 to 1.0, or 1 to 255)
  --export-png-color-mode=COLOR-MODE         Color mode (bit depth and color type) for exported bitmaps (Gray_1/Gray_2/Gray_4/Gray_8/Gray_16/RGB_8/RGB_16/GrayAlpha_8/GrayAlpha_16/RGBA_8/RGBA_16)
  --export-png-use-dithering=false|true      Force dithering or disables it
  -I, --query-id=OBJECT-ID[,OBJECT-ID]*      ID(s) of object(s) to be queried
  -S, --query-all                            Print bounding boxes of all objects
  -X, --query-x                              X coordinate of drawing or object (if specified by --query-id)
  -Y, --query-y                              Y coordinate of drawing or object (if specified by --query-id)
  -W, --query-width                          Width of drawing or object (if specified by --query-id)
  -H, --query-height                         Height of drawing or object (if specified by --query-id)
  --vacuum-defs                              Remove unused definitions from the <defs> section(s) of document
  --select=OBJECT-ID[,OBJECT-ID]*            Select objects: comma-separated list of IDs
  --actions=ACTION(:ARG)[;ACTION(:ARG)]*     List of actions (with optional arguments) to execute
  --action-list                              List all available actions
  -g, --with-gui                             With graphical user interface (required by some actions)
  --batch-process                            Close GUI after executing all actions
  --shell                                    Start Inkscape in interactive shell mode
  --display=DISPLAY                          X display to use


Examples:
  Export input SVG (in.svg) to PDF (out.pdf) format:
	inkscape --export-filename=out.pdf in.svg
  Export input files (in1.svg, in2.svg) to PNG format keeping original name (in1.png, in2.png):
	inkscape --export-type=png in1.svg in2.svg
  See 'man inkscape' and http://wiki.inkscape.org/wiki/index.php/Using_the_Command_Line for more details.
