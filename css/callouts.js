/**
 * Converts GitHub-style markdown callouts to HTML callouts
 * Looks for the pattern: > [!TYPE] followed by indented content
 */
function convertMarkdownCalloutsToHtml(markdownText) {
  // Define callout types and their icons
  const calloutTypes = {
    NOTE: '<i class="fas fa-info-circle"></i>',
    TIP: '<i class="fas fa-lightbulb"></i>',
    IMPORTANT: '<i class="fas fa-exclamation-circle"></i>',
    WARNING: '<i class="fas fa-exclamation-triangle"></i>',
    CAUTION: '<i class="fas fa-fire"></i>',
  };

  // Regex to match GitHub-style callouts
  // Matches: > [!TYPE] followed by lines starting with >
  const calloutRegex = />\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*\n((?:>\s*.*(?:\n|$))*)/gm;

  return markdownText.replace(calloutRegex, function (match, type, content) {
    // Normalize the type to handle case variations
    const normalizedType = type.toUpperCase();

    // Make sure we have a valid type, or default to NOTE
    const calloutType = Object.keys(calloutTypes).includes(normalizedType) ? normalizedType : "NOTE";

    // Process the content by removing the leading >
    const processedContent = content
      .split("\n")
      .map((line) => line.replace(/^\s*>\s?/, ""))
      .filter((line) => line !== "")
      .join("\n");

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
