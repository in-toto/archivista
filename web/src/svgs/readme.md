# Svgs

## Introduction

The `convert-svgs` script is a Node.js script that converts SVG files to inline SVGs for us, and we have it installed in our package.json scripts. We also have configured a gatsby plugin for loading svgs from the svgs folder. This can be used to make it easier to use SVGs in your project, as inline SVGs can be styled using CSS and are more accessible than using SVGs as images.

## Prerequisites

Before using the `convert-svgs` script, you'll need to have Node.js and npm installed on your system. You can download and install Node.js from the official website: https://nodejs.org/

## Usage

To use the `convert-svgs` script, follow these steps:

1. Drop an SVG file in the `src/svgs` folder of the project.
1. Open a terminal or command prompt and navigate to the root directory of the project.
1. Run the following command to convert your SVG files to inline SVGs:

```sh
npm run convert-svgs
```

This will convert any SVG files in the `src/svgs` folder that don't already have a matching `.inline.svg` file.
The inline SVGs will be saved with the same name as the original SVG file, but with `.inline.svg` appended to the end.

And that's it! Your SVG files should now be converted to inline SVGs, ready to be used in your project.
