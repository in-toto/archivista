import * as React from 'react';

import App from './src/App';

export type RootElementProps = {
  element: JSX.Element;
};

/**
 * # gatsby-browser.tsx
 * - The file you go to when you need to wrap the app with something global and dynamic
 * - see: https://www.gatsbyjs.com/docs/reference/config-files/gatsby-browser/
 *
 * This file contains the implementation of the `wrapPageElement` function,
 * which is used by Gatsby to wrap each page element with a custom component.
 *
 * The function takes in one parameters: `element`, which is the page element
 * to be wrapped.
 *
 * The function returns a JSX element that wraps the `element` parameter with the
 * custom component.
 *
 * Note that this function is executed in the browser, not during the build process.
 * @param {RootElementProps} { element }
 * @returns
 */
export const wrapRootElement = ({ element }: RootElementProps) => {
  return <App>{element}</App>;
};
