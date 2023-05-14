const Generator = require('yeoman-generator');
const chalk = require('chalk');
const yosay = require('yosay');

module.exports = class extends Generator {
  prompting() {
    this.log(yosay(`Welcome to the ${chalk.green('React Component Generator')}! 🚀`));

    const prompts = [
      {
        type: 'input',
        name: 'componentName',
        message: `Enter a ${chalk.cyan('PascalCase')} name for your component:`,
        validate: function (input) {
          if (!/^[A-Z][a-zA-Z]+(?:[A-Z][a-zA-Z]+)*$/.test(input)) {
            return "Component name must be in PascalCase";
          }
          return true;
        },
      },
      {
        type: 'confirm',
        name: 'storybook',
        message: `Do you want to generate a ${chalk.green('Storybook')} story for this component? ${chalk.yellow('(y/n)')}`,
        default: true,
      },
    ];

    return this.prompt(prompts).then((answers) => {
      this.answers = answers;
    });
  }

  writing() {
    const componentName = this.answers.componentName;
    const componentFolder = componentName.replace(/([A-Z])/g, '-$1').toLowerCase().slice(1);
  
    // Component file
    this.fs.copyTpl(
      this.templatePath('MyComponent.tsx'),
      this.destinationPath(`${componentFolder}/${componentName}.tsx`),
      {
        componentName,
      }
    );
  
    // Storybook story file
    if (this.answers.storybook) {
      this.fs.copyTpl(
        this.templatePath('MyComponent.stories.tsx'),
        this.destinationPath(`${componentFolder}/${componentName}.stories.tsx`),
        {
          componentName, 
        }
      );
    }
  
    // Jest test file
    this.fs.copyTpl(
      this.templatePath('MyComponent.test.tsx'),
      this.destinationPath(`${componentFolder}/${componentName}.test.tsx`),
      {
        componentName, 
      }
    );
  
    // Mock file
    this.fs.copyTpl(
      this.templatePath('__mocks__/MyComponent.mock.tsx'),
      this.destinationPath(`${componentFolder}/__mocks__/${componentName}.mock.tsx`),
      {
        componentName, 
      }
    );
  }
};
