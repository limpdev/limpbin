/**
 * Converts GitHub-style markdown callouts to HTML callouts within existing HTML blockquotes
 */
function convertMarkdownCalloutsToHtml(htmlText) {
  // Define callout types and their icons
  const calloutTypes = {
    NOTE: '<i class="fas fa-info-circle"></i>',
    TIP: '<i class="fas fa-lightbulb"></i>',
    IMPORTANT: '<i class="fas fa-exclamation-circle"></i>',
    WARNING: '<i class="fas fa-exclamation-triangle"></i>',
    CAUTION: '<i class="fas fa-fire"></i>',
  };

  // Regex to match GitHub-style callouts inside blockquotes in HTML
  // Matches: <blockquote><p>[!TYPE] ... </p></blockquote>
  const calloutRegex = /<blockquote>\s*<p>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*([\s\S]*?)<\/p>\s*<\/blockquote>/gm;

  return htmlText.replace(calloutRegex, function (match, type, content) {
    // Normalize the type to handle case variations
    const normalizedType = type.toUpperCase();

    // Make sure we have a valid type, or default to NOTE
    const calloutType = Object.keys(calloutTypes).includes(normalizedType) ? normalizedType : "NOTE";

    // Process the content - trim whitespace
    const processedContent = content.trim();

    // Build the HTML replacement
    return `<div class="callout callout-${calloutType.toLowerCase()}">
  <div class="callout-header">
    <span class="callout-icon">${calloutTypes[calloutType]}</span>
    <span class="callout-title">${calloutType}</span>
  </div>
  <div class="callout-content">
    <p>${processedContent}</p>
  </div>
</div>`;
  });
}

document.addEventListener("DOMContentLoaded", function () {
  // Get the current HTML content of the body
  let bodyHTML = document.body.innerHTML;

  // Convert the markdown callouts in the HTML to proper callout divs
  bodyHTML = convertMarkdownCalloutsToHtml(bodyHTML);

  // Replace the body's HTML with the converted content
  document.body.innerHTML = bodyHTML;
});
