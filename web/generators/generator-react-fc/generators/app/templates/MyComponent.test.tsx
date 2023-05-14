import <%= componentName %> from './<%= componentName %>';
import React from 'react';
import { render } from '@testing-library/react';

describe('<%= componentName %>', () => {
  it('should render without error', () => {
    render(<<%= componentName %> />);
  });
});
