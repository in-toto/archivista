import { render, screen } from "@testing-library/react";

import React from "react";
import ResponsiveTable from "./ResponsiveSourcesTable";

describe("ResponsiveSourcesTable", () => {
  const rows = [
    {
      nodeId: "1",
      nodeType: "Type 1",
      nodeGroup: "Group 1",
      registeredAt: "2022-02-26T10:00:00.000Z",
      registeredBy: "User 1",
      selectors: { selectors: [{ key: "Key 1", value: "Value 1" }] },
    },
    {
      nodeId: "2",
      nodeType: "Type 2",
      nodeGroup: "Group 2",
      registeredAt: "2022-02-27T10:00:00.000Z",
      registeredBy: "User 2",
      selectors: {
        selectors: [{ key: "Key 2", value: "Value 2" }],
      },
    },
  ];

  it("renders table headers and rows correctly", () => {
    render(<ResponsiveTable rows={rows} deactiveNode={jest.fn()} isLoading={false} />);

    // Test headers
    expect(screen.getByText("Source Type")).toBeInTheDocument();
    expect(screen.getByText("Source Group")).toBeInTheDocument();
    expect(screen.getByText("Registered Date")).toBeInTheDocument();
    expect(screen.getByText("User")).toBeInTheDocument();
    expect(screen.getByText("Selector")).toBeInTheDocument();

    // Test rows
    expect(screen.getByText("Type 1")).toBeInTheDocument();
    expect(screen.getByText("Group 1")).toBeInTheDocument();
    expect(screen.getByText("2/26/2022, 10:00:00 AM")).toBeInTheDocument();
    expect(screen.getByText("User 1")).toBeInTheDocument();
    expect(screen.getByText("Key 1:Value 1")).toBeInTheDocument();
    expect(screen.getByText("Type 2")).toBeInTheDocument();
    expect(screen.getByText("Group 2")).toBeInTheDocument();
    expect(screen.getByText("2/27/2022, 10:00:00 AM")).toBeInTheDocument();
    expect(screen.getByText("User 2")).toBeInTheDocument();
    expect(screen.getByText("Key 2:Value 2")).toBeInTheDocument();
  });

  it("renders a skeleton loader when loading", () => {
    render(<ResponsiveTable rows={[]} deactiveNode={jest.fn()} isLoading={true} />);

    expect(screen.getByTestId("responsive-table-skeleton")).toBeInTheDocument();
  });
});
