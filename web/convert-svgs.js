const fs = require('fs');
const path = require('path');
const glob = require('glob');

const srcDir = 'src/svgs';

// Get a list of all .svg files in the srcDir directory
glob(`${srcDir}/*.svg`, (err, files) => {
  if (err) throw err;

  files.forEach((filePath) => {
    const fileName = path.basename(filePath, '.svg');
    const inlineFilePath = path.join(srcDir, `${fileName}.inline.svg`);

    // If a matching .inline.svg file already exists, skip this file
    if (fs.existsSync(inlineFilePath)) {
      console.log(`Skipping ${filePath} - matching .inline.svg file already exists`);
      return;
    }

    // Read the contents of the .svg file
    const svg = fs.readFileSync(filePath, { encoding: 'utf8' });

    // Convert the svg to an inline svg string
    const inlineSvg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${svg.match(/ viewBox="([^"]*)"/)[1]}">${svg.match(/<svg[^>]*>([\s\S]*)<\/svg>/)[1]}</svg>`;

    // Write the inline svg to a new file with the same name but with .inline.svg appended
    fs.writeFileSync(inlineFilePath, inlineSvg);

    console.log(`Converted ${filePath} to ${inlineFilePath}`);
  });
});
