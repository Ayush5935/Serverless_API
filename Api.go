import csv
import boto3
from botocore.exceptions import ClientError

# Function to obtain SSO access token
def get_sso_access_token():
    session = boto3.session.Session()
    region = 'us-west-2'  # Update with your region
    sso_oidc = session.client('sso-oidc', region_name=region)

    # Start device authorization
    start_url = 'https://d-92670ca28f.awsapps.com/start#/'  # Update with your start URL
    client_creds = sso_oidc.register_client(clientName='myapp', clientType='public')
    device_authorization = sso_oidc.start_device_authorization(
        clientId=client_creds['clientId'],
        clientSecret=client_creds['clientSecret'],
        startUrl=start_url
    )
    url = device_authorization['verificationUriComplete']
    device_code = device_authorization['deviceCode']
    expires_in = device_authorization['expiresIn']
    interval = device_authorization['interval']

    # Poll for token
    for _ in range(1, expires_in // interval + 1):
        sleep(interval)
        try:
            token = sso_oidc.create_token(
                grantType='urn:ietf:params:oauth:grant-type:device_code',
                deviceCode=device_code,
                clientId=client_creds['clientId'],
                clientSecret=client_creds['clientSecret']
            )
            return token['accessToken']
        except sso_oidc.exceptions.AuthorizationPendingException:
            pass
    raise Exception("Failed to obtain SSO access token")

# Function to get VPC ID and Subnet ID for an instance
def get_vpc_subnet_id(instance_id, access_token, region):
    session = boto3.session.Session()
    sso = session.client('sso', region_name=region)

    # Get account ID
    account_id = boto3.client('sts').get_caller_identity().get('Account')

    # Assume role
    role_credentials = sso.get_role_credentials(
        roleName='DishWPaaSAdministrator',
        accountId=account_id,
        accessToken=access_token
    )['roleCredentials']

    # Create EC2 client with assumed role credentials
    ec2 = boto3.client(
        'ec2',
        region_name=region,
        aws_access_key_id=role_credentials['accessKeyId'],
        aws_secret_access_key=role_credentials['secretAccessKey'],
        aws_session_token=role_credentials['sessionToken']
    )

    try:
        # Describe instance
        response = ec2.describe_instances(InstanceIds=[instance_id])

        # Extract VPC ID and Subnet ID
        vpc_id = response['Reservations'][0]['Instances'][0]['VpcId']
        subnet_id = response['Reservations'][0]['Instances'][0]['SubnetId']

        return vpc_id, subnet_id
    except ClientError as e:
        if e.response['Error']['Code'] == 'InvalidInstanceID.NotFound':
            print(f"Instance {instance_id} not found in region {region}")
        else:
            raise e

# Main function
def main():
    input_file = 'input.csv'  # Input CSV file containing Account ID and Instance ID columns
    output_file = 'output.csv'  # Output CSV file to store results

    # Get SSO access token
    access_token = get_sso_access_token()

    # List of US regions
    us_regions = ['us-east-1', 'us-east-2', 'us-west-1', 'us-west-2']

    # Open output CSV file for writing
    with open(output_file, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Account ID', 'Instance ID', 'VPC ID', 'Subnet ID'])  # Write header

        # Open input CSV file for reading
        with open(input_file, mode='r') as csvfile:
            reader = csv.reader(csvfile)

            # Skip header row
            next(reader)

            # Iterate over rows in input CSV file
            for row in reader:
                account_id, instance_id = row

                # Get VPC ID and Subnet ID using SSO for each region
                for region in us_regions:
                    try:
                        vpc_id, subnet_id = get_vpc_subnet_id(instance_id, access_token, region)
                        writer.writerow([account_id, instance_id, vpc_id, subnet_id])
                        break  # If successful, break out of the loop and proceed to the next instance
                    except Exception as e:
                        print(f"Error processing instance {instance_id} in region {region}: {e}")

    print(f"Results saved to {output_file}")

if __name__ == "__main__":
    main()
