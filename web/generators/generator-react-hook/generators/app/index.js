const Generator = require("yeoman-generator");
const chalk = require("chalk");
const yosay = require("yosay");

module.exports = class extends Generator {
  prompting() {
    this.log(yosay(`Welcome to the ${chalk.green("React Hook Generator")}! 🎣`));

    const prompts = [
      {
        type: "input",
        name: "hookName",
        message: `Enter a ${chalk.cyan("camelCase")} name for your hook:`,
        validate: function(hookName) {
          const camelCaseRegex = /^[a-z]+[A-Za-z0-9]*$/;
          if (!camelCaseRegex.test(hookName)) {
            return "Hook name should be in camelCase";
          }

          return true;
        }
      },
      {
        type: "input",
        name: "propsName",
        message: "What is the name of the props type (optional)? If none, just press enter."
      },
      {
        type: "input",
        name: "desc",
        message: `Provide a ${chalk.cyan("description")} for your hook.`
      }
    ];

    return this.prompt(prompts).then(answers => {
      this.hookName = answers.hookName;
      this.desc = answers.desc;
      this.propsName = answers.propsName;
      this.answers = answers;
    });
  }

  writing() {
    const hookName = this.answers.hookName;
    const desc = this.desc;
    const propsName = this.propsName;
    const hookFileName = hookName + ".tsx";
    const testFileName = hookName + ".test.tsx";
    const hookFolderInKebab = hookName.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();

    this.fs.copyTpl(this.templatePath("useMyThing.tsx"), this.destinationPath(`${hookFolderInKebab}/${hookFileName}`), {
      desc,
      hookName,
      propsName
    });

    this.fs.copyTpl(this.templatePath("useMyThing.test.tsx"), this.destinationPath(`${hookFolderInKebab}/${testFileName}`), {
      desc,
      hookName,
      propsName
    });
  }
};
