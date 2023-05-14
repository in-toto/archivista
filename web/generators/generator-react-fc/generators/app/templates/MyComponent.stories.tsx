import <%= componentName %> from './<%= componentName %>';
import React from 'react';
import { storiesOf } from '@storybook/react';

storiesOf('MyComponent', module).add('default', () => <<%= componentName %> />);
