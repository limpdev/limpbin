const path = require("path");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
    // Entry point for our content script
    entry: {
        content: "./src/content.js",
    },
    // Output configuration
    output: {
        path: path.resolve(__dirname, "dist"),
        filename: "[name].bundle.js",
        clean: true, // Clean the dist folder before each build
    },
    // Plugins
    plugins: [
        new CopyPlugin({
            patterns: [
                { from: "public", to: "." }, // Copies files from public/ to dist/
            ],
        }),
    ],
    // Optional: configuration for resolving modules, etc.
    resolve: {
        extensions: [".js"],
    },
};