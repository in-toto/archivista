'use strict';
const path = require('path');
const assert = require('yeoman-assert');
const helpers = require('yeoman-test');

describe('generator-react-fc:app', () => {
  beforeAll(() => {
    return helpers
      .run(path.join(__dirname, '../generators/app'))
      .withPrompts({ name: 'MyComponent' });
  });

  it('creates files', () => {
    assert.file([
      'MyComponent/MyComponent.tsx',
      'MyComponent/__mocks__/MyComponent.mock.ts',
      'MyComponent/MyComponent.test.tsx',
    ]);
  });

  it('creates the component file with the correct content', () => {
    assert.fileContent('MyComponent/MyComponent.tsx', /import React from 'react';/);
    assert.fileContent('MyComponent/MyComponent.tsx', /export type MyComponentProps/);
    assert.fileContent('MyComponent/MyComponent.tsx', /const MyComponent/);
    assert.fileContent('MyComponent/MyComponent.tsx', /export default MyComponent;/);
  });

  it('creates the mock file with the correct content', () => {
    assert.fileContent(
      'MyComponent/__mocks__/MyComponent.mock.ts',
      /export default {.*}/s
    );
  });

  it('creates the test file with the correct content', () => {
    assert.fileContent('MyComponent/MyComponent.test.tsx', /import React from 'react';/);
    assert.fileContent('MyComponent/MyComponent.test.tsx', /import { render/);
    assert.fileContent('MyComponent/MyComponent.test.tsx', /describe\('MyComponent'/);
  });
});
