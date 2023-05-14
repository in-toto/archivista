module.exports = {
  env: {
    browser: true,
  },
  root: true,
  overrides: [
    {
      files: ["*.ts", "*.tsx"],
      processor: "@graphql-eslint/graphql",
      parser: "@typescript-eslint/parser",
      plugins: ["react", "@typescript-eslint", "prettier"],
      extends: [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:@typescript-eslint/recommended-requiring-type-checking",
        "plugin:react-hooks/recommended",
        "prettier",
      ],
      rules: {
        "@typescript-eslint/no-unsafe-assignment": "off",
        quotes: [2, "single", { avoidEscape: true }],
        "prettier/prettier": ["error", { singleQuote: true }],
      },
      env: {
        es6: true,
      },
      parser: "@typescript-eslint/parser",
      parserOptions: {
        project: "./tsconfig.json",
        tsconfigRootDir: __dirname,
        ecmaFeatures: {
          jsx: true,
        },
        ecmaVersion: 2021,
        sourceType: "module",
      },
      settings: {
        react: {
          version: "detect",
        },
      },
    },
    {
      files: ["*.graphql"],
      parser: "@graphql-eslint/eslint-plugin",
      plugins: ["@graphql-eslint"],
      rules: {
        "@graphql-eslint/no-anonymous-operations": "error",
        "@graphql-eslint/naming-convention": [
          "error",
          {
            OperationDefinition: {
              style: "PascalCase",
              forbiddenPrefixes: ["Query", "Mutation", "Subscription", "Get"],
              forbiddenSuffixes: ["Query", "Mutation", "Subscription"],
            },
          },
        ],
      },
    },
    {
      files: ["*.test.ts", "*.test.tsx"],
      extends: ["plugin:jest/recommended"],
      plugins: ["jest"],
      env: {
        "jest/globals": true,
      },
      parserOptions: {
        project: "./testing.tsconfig.json",
        tsconfigRootDir: __dirname,
      },
    },
  ],
};
