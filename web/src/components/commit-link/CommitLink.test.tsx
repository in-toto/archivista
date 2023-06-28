import '@testing-library/jest-dom'; // Import testing-library/jest-dom

import CommitLink from './CommitLink';
import React from 'react';
import { render } from '@testing-library/react';

describe('CommitLink', () => {
  test('renders the component with the provided commit', () => {
    const commit = 'abcd1234efgh5678';

    const { getByText } = render(<CommitLink commit={commit} />);

    const commitText = getByText('Commit: efgh567');

    expect(commitText).toBeInTheDocument();
  });
});
