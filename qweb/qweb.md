# qweb - markdown -> html, with golang

### USAGE: qweb [doc.md] [style.css] -> then, you let it fuckin rip baby

> pretty much as simple as that, nothing special about it, besides the fact that it does the job right, with no clever shit or goofy ass command-line syntax.

### -> READ HERE FOR CODE HIGHLIGHTING <-

`qweb` will automatically add the link for code highlighting + a script, both courtesy of **prismJS**
> Here is what is does...
```html
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Markdown Documentation</title>
<link href="prism.css" rel="stylesheet" />
<style>%s</style>
</head>
<body>
<script src="prism.js"></script>
%s
</body>
</html>
```
> Note the `link` & `script` tags in the head and body, respectively.
#### 📌 CSS and JS must be in the same directory as the highlighted HTML file...
> Unfortunately, due to restrictions in the scopability of web browsers (security, duh), Chrome cannot simply reference locally stored CSS/JS files even if they are on `PATH`
>> PERHAPS I can figure a way to import the inline javscript and css so that each HTML has it's own sort of *engine* - the major concern here would be the unbelievable bloat this would add to the HTML document.
