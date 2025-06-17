To use :not() in CSS to exclude a `<pre>` tag when targeting a combination of hljs, code, and pre, you need to understand how the elements are nested and how CSS specificity works.

You want to style code blocks that use the .hljs class, typically added by highlight.js, but only apply styles when it’s not a `<pre>` tag.

## Highlight.js

Highlight.js often wraps code like this:

`<pre><code class="hljs">...</code></pre>`

`<pre class="hljs">...</pre>`

`<code class="hljs">...</code>`


##### Example CSS Rule

If you’re trying to apply a style to .hljs that is not inside a `<pre>`, here’s how to write it:

```css
.hljs:not(pre) {
  background-color: yellow;
}
```

But this only matches elements with the .hljs class that are not `<pre>` elements. If you’re targeting code.hljs but want to exclude cases where it’s inside a `<pre>`, you need a different approach:

> Exclude code.hljs inside pre

```css
code.hljs:not(pre code.hljs) {
  /* This doesn't work the way you might hope */
}
```

> Instead, use the direct descendant selector to limit styling to only code.hljs not nested inside pre:

```css
/* Style code.hljs NOT inside pre */
:not(pre) > code.hljs {
  color: red;
}
```

Any `code.hljs` whose direct parent is not a `pre`.

#### Summary of Working Selectors

|             Use Case              |        Selector         |
| :-------------------------------: | :---------------------: |
| Target .hljs that’s not a `<pre>` |    `.hljs:not(pre)`     |
| Target code.hljs not in a `<pre>` | `:not(pre) > code.hljs` |
|       Target only pre.hljs        |       `pre.hljs`        |
| Target only code.hljs inside pre  |    `pre > code.hljs`    |

```css
/* Style code snippets that are not inside a `<pre>` */
:not(pre) > code.hljs {
  background: #f0f0f0;
  padding: 0.2em;
}

/* Style pre.hljs normally */
pre.hljs {
  background: #1e1e1e;
  color: white;
  padding: 1em;
}
```

---

## Additional  JavaScript

> _add a little flair why don't ya..._

```js
// Add ripple effect for every mouse click, anywhere on the page using an SVG
document.addEventListener('click', function (e) {
	// Create a container for the ripple effect
	const rippleContainer = document.createElement('div');
	rippleContainer.style.position = 'fixed';
	rippleContainer.style.left = e.clientX - 48 + 'px'; // Center the ripple at click position
	rippleContainer.style.top = e.clientY - 48 + 'px';
	rippleContainer.style.pointerEvents = 'none'; // Don't interfere with further clicks
	rippleContainer.style.zIndex = '9999';

	// Create SVG element
	const svgNS = 'http://www.w3.org/2000/svg';
	const svg = document.createElementNS(svgNS, 'svg');
	svg.setAttribute('width', '96');
	svg.setAttribute('height', '96');
	svg.setAttribute('viewBox', '0 0 24 24');

	// Create circle element
	const circle = document.createElementNS(svgNS, 'circle');
	circle.setAttribute('cx', '12');
	circle.setAttribute('cy', '12');
	circle.setAttribute('r', '0');
	circle.setAttribute('fill', 'rgba(168, 168, 168, 0.5)');

	// Create animate elements
	const animateRadius = document.createElementNS(svgNS, 'animate');
	animateRadius.setAttribute('attributeName', 'r');
	animateRadius.setAttribute('calcMode', 'spline');
	animateRadius.setAttribute('dur', '0.4s');
	animateRadius.setAttribute('keySplines', '.52,.6,.25,.99');
	animateRadius.setAttribute('values', '0;11');
	animateRadius.setAttribute('fill', 'freeze');

	const animateOpacity = document.createElementNS(svgNS, 'animate');
	animateOpacity.setAttribute('attributeName', 'opacity');
	animateOpacity.setAttribute('calcMode', 'spline');
	animateOpacity.setAttribute('dur', '0.4s');
	animateOpacity.setAttribute('keySplines', '.52,.6,.25,.99');
	animateOpacity.setAttribute('values', '1;0');
	animateOpacity.setAttribute('fill', 'freeze');

	// Assemble the SVG
	circle.appendChild(animateRadius);
	circle.appendChild(animateOpacity);
	svg.appendChild(circle);
	rippleContainer.appendChild(svg);

	// Add to document
	document.body.appendChild(rippleContainer);

	// Remove after animation completes
	setTimeout(() => {
		document.body.removeChild(rippleContainer);
	}, 500); // Match the duration of the animation
});

```