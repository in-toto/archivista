/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
import { fireEvent, render } from "@testing-library/react";

import React from "react";
import Wizard from "./NewRegistrationWizard";

test("renders wizard with initial step", () => {
  const { getByText, getByLabelText } = render(<Wizard />);
  expect(getByText("Choose Node Type")).toBeInTheDocument();
  expect(getByLabelText("TPM")).toBeInTheDocument();
});

test("advances to next step on selecting node type", () => {
  const { getByLabelText, getByText } = render(<Wizard />);
  fireEvent.click(getByLabelText("Google Cloud"));
  fireEvent.click(getByText("Next"));
  expect(getByText("Choose Node Group")).toBeInTheDocument();
});

test("shows appropriate selector key based on node type", () => {
  const { getByLabelText, getByText } = render(<Wizard />);
  fireEvent.click(getByLabelText("GitHub"));
  fireEvent.click(getByText("Next"));
  expect(getByText("Selector Key: Repository")).toBeInTheDocument();
});

test("advances to next step on selecting node group", () => {
  const { getByLabelText, getByText } = render(<Wizard />);
  fireEvent.click(getByLabelText("prod"));
  fireEvent.click(getByText("Next"));
  expect(getByText("Choose Selector Key")).toBeInTheDocument();
});

test("advances to next step on entering selector value", () => {
  const { getByLabelText, getByText } = render(<Wizard />);
  fireEvent.click(getByLabelText("Google Cloud"));
  fireEvent.click(getByText("Next"));
  fireEvent.click(getByLabelText("prod"));
  fireEvent.click(getByText("Next"));
  fireEvent.change(getByLabelText("Selector Value"), { target: { value: "abcd1234" } });
  fireEvent.click(getByText("Next"));
  expect(getByText("All steps completed")).toBeInTheDocument();
});

test("resets wizard on clicking reset button", () => {
  const { getByLabelText, getByText } = render(<Wizard />);
  fireEvent.click(getByLabelText("Google Cloud"));
  fireEvent.click(getByText("Next"));
  fireEvent.click(getByLabelText("prod"));
  fireEvent.click(getByText("Next"));
  fireEvent.change(getByLabelText("Selector Value"), { target: { value: "abcd1234" } });
  fireEvent.click(getByText("Next"));
  fireEvent.click(getByText("Reset"));
  expect(getByText("Choose Node Type")).toBeInTheDocument();
});
