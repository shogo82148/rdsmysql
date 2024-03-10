#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { IntegrationStack } from "../lib/integration-stack";

const app = new cdk.App();
new IntegrationStack(app, "IntegrationStack", {
  env: { account: "339712736426", region: "ap-northeast-1" },
});
