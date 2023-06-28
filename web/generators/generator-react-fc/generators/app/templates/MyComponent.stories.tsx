import <%= componentName %> from './<%= componentName %>';
import React from 'react';
import { storiesOf } from '@storybook/react';

storiesOf('<%= componentName %>', module).add('default', () => <<%= componentName %> />);
