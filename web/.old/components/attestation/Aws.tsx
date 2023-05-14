import * as React from "react";
import { AwsAttestor } from "../../lib/Witness";
import { Typography } from "@mui/material";

export function Aws(props: { data: AwsAttestor }) {
  const attestation = props.data;
  return (
    <>
      <Typography variant="body2">
        <b>Devpay Product Codes:</b> {attestation.devpayProductCodes}
      </Typography>
      <Typography variant="body2">
        <b>Marketplace Product Codes:</b> {attestation.marketplaceProductCodes}
      </Typography>
      <Typography variant="body2">
        <b>Availability Zone:</b> {attestation.availabilityZone}
      </Typography>
      <Typography variant="body2">
        <b>Private IP:</b> {attestation.privateIp}
      </Typography>
      <Typography variant="body2">
        <b>Version:</b> {attestation.version}
      </Typography>
      <Typography variant="body2">
        <b>Region:</b> {attestation.region}
      </Typography>
      <Typography variant="body2">
        <b>InstanceId:</b> {attestation.instanceId}
      </Typography>
      <Typography variant="body2">
        <b>Billing Products:</b> {attestation.billingProducts}
      </Typography>
      <Typography variant="body2">
        <b>Instance Type:</b> {attestation.instanceType}
      </Typography>
      <Typography variant="body2">
        <b>Account Id:</b> {attestation.accountId}
      </Typography>
      <Typography variant="body2">
        <b>Pending Time:</b> {attestation.pendingTime}
      </Typography>
      <Typography variant="body2">
        <b>Image Id:</b> {attestation.imageId}
      </Typography>
      <Typography variant="body2">
        <b>Kernal Id:</b> {attestation.kernalId}
      </Typography>
      <Typography variant="body2">
        <b>Ramdisk Id:</b> {attestation.ramdiskId}
      </Typography>
      <Typography variant="body2">
        <b>Architecture:</b> {attestation.architecture}
      </Typography>
      <Typography variant="body2">
        <b>Raw Iid:</b> {attestation.rawiid}
      </Typography>
      <Typography variant="body2">
        <b>Raw Signature:</b> {attestation.rawsig}
      </Typography>
      <Typography variant="body2">
        <b>Public Key:</b> {attestation.publickey}
      </Typography>
    </>
  );
}
