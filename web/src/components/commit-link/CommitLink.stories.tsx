import CommitLink from './CommitLink';
import React from 'react';
import { storiesOf } from '@storybook/react';

const commit = 'asdfbu892jlkj23';

storiesOf('MyComponent', module).add('default', () => <CommitLink commit={commit} />);
