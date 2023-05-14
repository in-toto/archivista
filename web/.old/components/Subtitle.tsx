import * as React from "react";
import { Theme, Typography } from "@mui/material";
import { SxProps } from "@mui/system";

export interface Props {
  sx?: SxProps<Theme>;
}

export function Subtitle({ children, sx }: React.PropsWithChildren<Props>) {
  return (
    <Typography variant="h6" sx={sx ?? {}} gutterBottom>
      {children}
    </Typography>
  );
}
