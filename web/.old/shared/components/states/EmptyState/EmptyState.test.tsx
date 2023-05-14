import { fireEvent, render, screen } from "@testing-library/react";

import EmptyState from "./EmptyState";
import React from "react";

describe("EmptyState component", () => {
  it("should render with header, subheader, and button text", () => {
    const headerText = "No items found";
    const subheaderText = "Please add some items to continue";
    const addText = "Add item";
    const onAddButtonClick = jest.fn();
    render(<EmptyState headerText={headerText} subheaderText={subheaderText} addText={addText} onAddButtonClick={onAddButtonClick} />);

    expect(screen.getByText(headerText)).toBeInTheDocument();
    expect(screen.getByText(subheaderText)).toBeInTheDocument();
    expect(screen.getByText(addText)).toBeInTheDocument();
  });

  it("should call onAddButtonClick when add button is clicked", () => {
    const headerText = "No items found";
    const subheaderText = "Please add some items to continue";
    const addText = "Add item";
    const onAddButtonClick = jest.fn();
    render(<EmptyState headerText={headerText} subheaderText={subheaderText} addText={addText} onAddButtonClick={onAddButtonClick} />);

    fireEvent.click(screen.getByText(addText));
    expect(onAddButtonClick).toHaveBeenCalled();
  });
});
