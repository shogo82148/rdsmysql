import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as rds from 'aws-cdk-lib/aws-rds';

export class IntegrationStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // VPC
    const vpc = new ec2.Vpc(this, 'MyVpc', {
      maxAzs: 2,
      subnetConfiguration: [
        {
          name: 'public',
          subnetType: ec2.SubnetType.PUBLIC,
        },
      ],
    });

    const rdsInstance = new rds.DatabaseInstance(this, 'MyDatabaseInstance', {
      engine: rds.DatabaseInstanceEngine.mysql({
        version: rds.MysqlEngineVersion.VER_8_0_32
      }),
      allocatedStorage: 100,
      credentials: {
        username: 'admin',
        password: cdk.SecretValue.unsafePlainText('secret'),
      },
      deletionProtection: true,
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      vpc: vpc
    });
  }
}
