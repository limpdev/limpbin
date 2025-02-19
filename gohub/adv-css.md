# Advanced `CSS`
> assorted features/functions using more advanced `css`, `html`, and `javascript`

# Gunmetal Background for Pages

```css
.background {
    background-image: linear-gradient #272a2c, #161819;
    background-color: initial;
```

## Button -> `copy code`
Make sure the `pre` blocks within the HTML are wrapped by a `div`, perhaps called "pre-container":
```html
<div class="code-container">
    <pre>
        <code class="language-javascript">
         // This is some example javascript code
         function myFunction() {
             console.log("some shit for the garbage collector");
         }
        </code>
    </pre>
    <button class="copy-button">Copy Code</button>
</div>
```
> the `button` can go right at the end, before closing the pre-container.

#### **now we need to make the button**
```css
.copy-button {
     position: absolute; /* Positioned relative to .code-container */
     top: 10px;
     right: 10px;
     background-color: #4CAF50;
     color: white;
     border: none;
     padding: 8px 12px;
     border-radius: 4px;
     cursor: pointer;
     font-size: 14px;
     opacity: 0.7; /* semi-transparent */
 }

 .copy-button:hover {
     opacity: 1; /* Fully opaque on hover */
 }

.copy-button:focus {
     outline: none; /* Remove focus outline */
 }

.code-container {
    position: relative; /* Needed for absolute positioning of the button */
    margin-bottom: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
    padding: 10px;
    background-color: #f9f9f9;
}

pre {
    margin: 0;
    overflow-x: auto; /* Allows horizontal scrolling if code is too wide */
    padding: 10px;
}
```
### 🧠Clipboard Logic, in Javascript
```js
document.addEventListener('DOMContentLoaded', function() {

    document.querySelectorAll('.code-container').forEach(container => {
        const copyButton = container.querySelector('.copy-button');
        const codeBlock = container.querySelector('pre code') || container.querySelector('pre');

        copyButton.addEventListener('click', () => {
            const codeToCopy = codeBlock.textContent;

            navigator.clipboard.writeText(codeToCopy)
              .then(() => {
                copyButton.textContent = 'Copied!';
                setTimeout(() => {
                    copyButton.textContent = 'Copy Code';
                }, 2000);
               })
              .catch(err => {
                console.error('Failed to copy', err);
                copyButton.textContent = 'Copy Failed';
                setTimeout(() => {
                    copyButton.textContent = 'Copy Code';
                }, 2000);
               });
        });
    });
});
```

---
### *blurred theme for zed editor*
---

> `catppuccino espresso w/ blur`
```json
{
  "name": "Catppuccin Espresso (Blur)",
  "appearance": "dark",
  "style": {
    "accents": [
      "#c6aaf666",
      "#bac1f766",
      "#8cc7e766",
      "#add8a866",
      "#e6d4b066",
      "#ecb19766",
      "#e696a966"
    ],
    "background.appearance": "blurred",
    "border": "#00000000",
    "border.variant": "#00000000",
    "border.focused": "#b7bdf8",
    "border.selected": "#c6a0f6",
    "border.transparent": "#a6da95",
    "border.disabled": "#6e738d",
    "elevated_surface.background": "#0a0a0a",
    "surface.background": "#000000d0",
    "background": "#00000887",
    "element.background": "#181926",
    "element.hover": "#494d644d",
    "element.active": "#00000000",
    "element.selected": "#363a4f4d",
    "element.disabled": "#6e738d",
    "drop_target.background": "#f4dbd618",
    "ghost_element.background": "#18192659",
    "ghost_element.hover": "#f4dbd608",
    "ghost_element.active": "#f4dbd612",
    "ghost_element.selected": "#f4dbd612",
    "ghost_element.disabled": "#6e738d",
    "text": "#cad3f5",
    "text.muted": "#b8c0e0",
    "text.placeholder": "#5b6078",
    "text.disabled": "#494d64",
    "text.accent": "#c6a0f6",
    "icon": "#cad3f5",
    "icon.muted": "#8087a2",
    "icon.disabled": "#6e738d",
    "icon.placeholder": "#5b6078",
    "icon.accent": "#c6a0f6",
    "status_bar.background": "#000000d7",
    "title_bar.background": "#00000887",
    "title_bar.inactive_background": "#181926d9",
    "toolbar.background": "#00000000",
    "tab_bar.background": "#00000000",
    "tab.inactive_background": "#00000000",
    "tab.active_background": "#f4dbd612",
    "search.match_background": "#8bd5ca33",
    "panel.background": "#00000000",
    "panel.focused_border": "00000000",
    "panel.indent_guide": "#363a4f99",
    "panel.indent_guide_active": "#5b6078",
    "panel.indent_guide_hover": "#c6a0f6",
    "pane.focused_border": "#cad3f5",
    "pane_group.border": "#363a4f",
    "scrollbar.thumb.background": "#f4dbd612",
    "scrollbar.thumb.hover_background": "#6e738d",
    "scrollbar.thumb.border": "#c6a0f6",
    "scrollbar.track.background": "#00000000",
    "scrollbar.track.border": "#00000000",
    "editor.foreground": "#cad3f5",
    "editor.background": "#00000000",
    "editor.gutter.background": "#00000000",
    "editor.subheader.background": "#1e2030",
    "editor.active_line.background": "#00000000",
    "editor.highlighted_line.background": "#f4dbd612",
    "editor.line_number": "#ffffff20",
    "editor.active_line_number": "#f4dbd690",
    "editor.invisible": "#939ab766",
    "editor.wrap_guide": "#5b6078",
    "editor.active_wrap_guide": "#5b6078",
    "editor.document_highlight.bracket_background": "#f4dbd640",
    "editor.document_highlight.read_background": "#a5adcb29",
    "editor.document_highlight.write_background": "#a5adcb29",
    "editor.indent_guide": "#363a4f99",
    "editor.indent_guide_active": "#5b6078",
    "terminal.background": "#00000000",
    "terminal.ansi.background": "#24273a",
    "terminal.foreground": "#cad3f5",
    "terminal.dim_foreground": "#8087a2",
    "terminal.bright_foreground": "#cad3f5",
    "terminal.ansi.black": "#494d64",
    "terminal.ansi.red": "#ed8796",
    "terminal.ansi.green": "#a6da95",
    "terminal.ansi.yellow": "#eed49f",
    "terminal.ansi.blue": "#8aadf4",
    "terminal.ansi.magenta": "#f5bde6",
    "terminal.ansi.cyan": "#8bd5ca",
    "terminal.ansi.white": "#b8c0e0",
    "terminal.ansi.bright_black": "#5b6078",
    "terminal.ansi.bright_red": "#ed8796",
    "terminal.ansi.bright_green": "#a6da95",
    "terminal.ansi.bright_yellow": "#eed49f",
    "terminal.ansi.bright_blue": "#8aadf4",
    "terminal.ansi.bright_magenta": "#f5bde6",
    "terminal.ansi.bright_cyan": "#8bd5ca",
    "terminal.ansi.bright_white": "#a5adcb",
    "terminal.ansi.dim_black": "#494d64",
    "terminal.ansi.dim_red": "#ed8796",
    "terminal.ansi.dim_green": "#a6da95",
    "terminal.ansi.dim_yellow": "#eed49f",
    "terminal.ansi.dim_blue": "#8aadf4",
    "terminal.ansi.dim_magenta": "#f5bde6",
    "terminal.ansi.dim_cyan": "#8bd5ca",
    "terminal.ansi.dim_white": "#b8c0e0",
    "link_text.hover": "#91d7e3",
    "conflict": "#eed49f",
    "conflict.border": "#eed49f",
    "conflict.background": "#1e2030",
    "created": "#a6da95",
    "created.border": "#a6da95",
    "created.background": "#1e2030",
    "deleted": "#ed8796",
    "deleted.border": "#ed8796",
    "deleted.background": "#1e2030",
    "hidden": "#6e738d",
    "hidden.border": "#6e738d",
    "hidden.background": "#1e2030",
    "hint": "#5b6078",
    "hint.border": "#5b6078",
    "hint.background": "#1a1a1ac0",
    "ignored": "#6e738d",
    "ignored.border": "#6e738d",
    "ignored.background": "#1e2030",
    "modified": "#eed49f",
    "modified.border": "#eed49f",
    "modified.background": "#1e2030",
    "predictive": "#6e738d",
    "predictive.border": "#b7bdf8",
    "predictive.background": "#1e2030",
    "renamed": "#7dc4e4",
    "renamed.border": "#7dc4e4",
    "renamed.background": "#1e2030",
    "info": "#8bd5ca",
    "info.border": "#8bd5ca",
    "info.background": "#7ab3b3",
    "warning": "#eed49f",
    "warning.border": "#eed49f",
    "warning.background": "#d3a168",
    "error": "#ed8796",
    "error.border": "#ed8796",
    "error.background": "#ed8796",
    "success": "#a6da95",
    "success.border": "#a6da95",
    "success.background": "#8cbe6c",
    "unreachable": "#ed8796",
    "unreachable.border": "#ed8796",
    "unreachable.background": "#ed87961f",
    "players": [
      {
        "cursor": "#f4dbd6",
        "selection": "#5b607880",
        "background": "#f4dbd6"
      },
      {
        "cursor": "#c6aaf6",
        "selection": "#c6aaf633",
        "background": "#c6aaf6"
      },
      {
        "cursor": "#bac1f7",
        "selection": "#bac1f733",
        "background": "#bac1f7"
      },
      {
        "cursor": "#8cc7e7",
        "selection": "#8cc7e733",
        "background": "#8cc7e7"
      },
      {
        "cursor": "#add8a8",
        "selection": "#add8a833",
        "background": "#add8a8"
      },
      {
        "cursor": "#e6d4b0",
        "selection": "#e6d4b033",
        "background": "#e6d4b0"
      },
      {
        "cursor": "#ecb197",
        "selection": "#ecb19733",
        "background": "#ecb197"
      },
      {
        "cursor": "#e696a9",
        "selection": "#e696a933",
        "background": "#e696a9"
      }
    ],
    "syntax": {
      "variable": {
        "color": "#cad3f5",
        "font_style": null,
        "font_weight": null
      },
      "variable.builtin": {
        "color": "#ed8796",
        "font_style": null,
        "font_weight": null
      },
      "variable.parameter": {
        "color": "#ee99a0",
        "font_style": null,
        "font_weight": null
      },
      "variable.member": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "variable.special": {
        "color": "#f5bde6",
        "font_style": "italic",
        "font_weight": null
      },
      "constant": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "constant.builtin": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "constant.macro": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "module": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "label": {
        "color": "#7dc4e4",
        "font_style": null,
        "font_weight": null
      },
      "string": {
        "color": "#a6da95",
        "font_style": null,
        "font_weight": null
      },
      "string.documentation": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": null
      },
      "string.regexp": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "string.escape": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "string.special": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "string.special.path": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "string.special.symbol": {
        "color": "#f0c6c6",
        "font_style": null,
        "font_weight": null
      },
      "string.special.url": {
        "color": "#f4dbd6",
        "font_style": "italic",
        "font_weight": null
      },
      "character": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": null
      },
      "character.special": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "boolean": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "number": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "number.float": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "type": {
        "color": "#eed49f",
        "font_style": null,
        "font_weight": null
      },
      "type.builtin": {
        "color": "#c6a0f6",
        "font_style": "italic",
        "font_weight": null
      },
      "type.definition": {
        "color": "#eed49f",
        "font_style": null,
        "font_weight": null
      },
      "type.interface": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "type.super": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "attribute": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "property": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "function": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "function.builtin": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "function.call": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "function.macro": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": null
      },
      "function.method": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "function.method.call": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "constructor": {
        "color": "#f0c6c6",
        "font_style": null,
        "font_weight": null
      },
      "operator": {
        "color": "#91d7e3",
        "font_style": null,
        "font_weight": null
      },
      "keyword": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.modifier": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.type": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.coroutine": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.function": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.operator": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.import": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.repeat": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.return": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.debug": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.exception": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.conditional": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.conditional.ternary": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.directive": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.directive.define": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "keyword.export": {
        "color": "#91d7e3",
        "font_style": null,
        "font_weight": null
      },
      "punctuation": {
        "color": "#939ab7",
        "font_style": null,
        "font_weight": null
      },
      "punctuation.delimiter": {
        "color": "#939ab7",
        "font_style": null,
        "font_weight": null
      },
      "punctuation.bracket": {
        "color": "#939ab7",
        "font_style": null,
        "font_weight": null
      },
      "punctuation.special": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "punctuation.special.symbol": {
        "color": "#f0c6c6",
        "font_style": null,
        "font_weight": null
      },
      "punctuation.list_marker": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": null
      },
      "comment": {
        "color": "#939ab7",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.doc": {
        "color": "#939ab7",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.documentation": {
        "color": "#939ab7",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.error": {
        "color": "#ed8796",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.warning": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.hint": {
        "color": "#8aadf4",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.todo": {
        "color": "#f0c6c6",
        "font_style": "italic",
        "font_weight": null
      },
      "comment.note": {
        "color": "#f4dbd6",
        "font_style": "italic",
        "font_weight": null
      },
      "diff.plus": {
        "color": "#a6da95",
        "font_style": null,
        "font_weight": null
      },
      "diff.minus": {
        "color": "#ed8796",
        "font_style": null,
        "font_weight": null
      },
      "tag": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "tag.attribute": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "tag.delimiter": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": null
      },
      "parameter": {
        "color": "#ee99a0",
        "font_style": null,
        "font_weight": null
      },
      "field": {
        "color": "#b7bdf8",
        "font_style": null,
        "font_weight": null
      },
      "namespace": {
        "color": "#eed49f",
        "font_style": "italic",
        "font_weight": null
      },
      "float": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "symbol": {
        "color": "#f5bde6",
        "font_style": null,
        "font_weight": null
      },
      "string.regex": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "text": {
        "color": "#cad3f5",
        "font_style": null,
        "font_weight": null
      },
      "emphasis.strong": {
        "color": "#ee99a0",
        "font_style": null,
        "font_weight": 700
      },
      "emphasis": {
        "color": "#ee99a0",
        "font_style": "italic",
        "font_weight": null
      },
      "embedded": {
        "color": "#ee99a0",
        "font_style": null,
        "font_weight": null
      },
      "text.literal": {
        "color": "#a6da95",
        "font_style": null,
        "font_weight": null
      },
      "concept": {
        "color": "#7dc4e4",
        "font_style": null,
        "font_weight": null
      },
      "enum": {
        "color": "#8bd5ca",
        "font_style": null,
        "font_weight": 700
      },
      "function.decorator": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "type.class.definition": {
        "color": "#eed49f",
        "font_style": null,
        "font_weight": 700
      },
      "hint": {
        "color": "#5b6078",
        "font_style": "italic",
        "font_weight": null
      },
      "link_text": {
        "color": "#8aadf4",
        "font_style": null,
        "font_weight": null
      },
      "link_uri": {
        "color": "#f4dbd6",
        "font_style": "italic",
        "font_weight": null
      },
      "parent": {
        "color": "#f5a97f",
        "font_style": null,
        "font_weight": null
      },
      "predictive": {
        "color": "#6e738d",
        "font_style": null,
        "font_weight": null
      },
      "predoc": {
        "color": "#ed8796",
        "font_style": null,
        "font_weight": null
      },
      "primary": {
        "color": "#ee99a0",
        "font_style": null,
        "font_weight": null
      },
      "tag.doctype": {
        "color": "#c6a0f6",
        "font_style": null,
        "font_weight": null
      },
      "string.doc": {
        "color": "#8bd5ca",
        "font_style": "italic",
        "font_weight": null
      },
      "title": {
        "color": "#cad3f5",
        "font_style": null,
        "font_weight": 800
      },
      "variant": {
        "color": "#ed8796",
        "font_style": null,
        "font_weight": null
      }
    }
  }
}
]
}

```
