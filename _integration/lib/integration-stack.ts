import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";

export class IntegrationStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // VPC
    const vpc = new ec2.Vpc(this, "MyVpc", {});

    // Bastion EC2 Instance
    const bastionSG = new ec2.SecurityGroup(this, "BastionSG", {
      vpc,
    });
    const bastion = new ec2.BastionHostLinux(this, "Bastion", {
      vpc,
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T4G,
        ec2.InstanceSize.MICRO
      ),
      machineImage: new ec2.AmazonLinux2023ImageSsmParameter({
        cpuType: ec2.AmazonLinuxCpuType.ARM_64,
      }),
      securityGroup: bastionSG,
    });
  }
}
