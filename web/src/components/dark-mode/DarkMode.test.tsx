import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/react';

import { DarkMode } from './DarkMode';
import React from 'react';
import { ThemeContext } from '../../shared/contexts/theme-context/ThemeContext';

describe('DarkMode', () => {
  // eslint-disable-next-line jest/expect-expect
  it('should render without crashing', () => {
    render(<DarkMode />);
  });

  it("should render the label 'Light' when isDarkMode is false", () => {
    const themeContext = { isDarkMode: false, toggleDarkMode: jest.fn() };
    render(
      <ThemeContext.Provider value={themeContext}>
        <DarkMode />
      </ThemeContext.Provider>
    );
    expect(screen.getByText('Light')).toBeInTheDocument();
  });

  it("should render the label 'Dark' when isDarkMode is true", () => {
    const themeContext = { isDarkMode: true, toggleDarkMode: jest.fn() };
    render(
      <ThemeContext.Provider value={themeContext}>
        <DarkMode />
      </ThemeContext.Provider>
    );
    expect(screen.getByText('Dark')).toBeInTheDocument();
  });

  // eslint-disable-next-line jest/no-commented-out-tests
  // it("should call toggleDarkMode when Switch is toggled", () => {
  //   const toggleDarkMode = jest.fn();
  //   render(
  //     <ThemeContext.Provider value={{ isDarkMode: false, toggleDarkMode }}>
  //       <DarkMode />
  //     </ThemeContext.Provider>
  //   );
  //   userEvent.click(screen.getByRole("checkbox"));
  //   expect(toggleDarkMode).toHaveBeenCalledTimes(1);
  // });
});
