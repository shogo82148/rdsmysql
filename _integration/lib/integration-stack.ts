import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as rds from "aws-cdk-lib/aws-rds";

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

    // Relational Database Service
    const rdsSG = new ec2.SecurityGroup(this, "RDSSG", {
      vpc,
    });
    rdsSG.addIngressRule(bastionSG, ec2.Port.tcp(3306));
    const cluster = new rds.DatabaseCluster(this, "Database", {
      engine: rds.DatabaseClusterEngine.auroraMysql({
        version: rds.AuroraMysqlEngineVersion.VER_3_05_2,
      }),
      writer: rds.ClusterInstance.serverlessV2("writer"),
      vpc,
      securityGroups: [rdsSG],
    });
    cluster.secret?.grantRead(bastion);
  }
}
