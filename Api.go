import csv
import boto3
import webbrowser
from time import sleep

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

    # Open browser for user authentication
    webbrowser.open(url, autoraise=True)

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

# Function to assume IAM role and create EC2 client
def assume_role_and_create_ec2_client(access_token, region):
    sso = boto3.client('sso', region_name=region)

    # Assume role
    role_credentials = sso.get_role_credentials(
        roleName='DishWPaaSAdministrator',
        accessToken=access_token
    )['roleCredentials']

    # Create EC2 client with assumed role credentials
    session = boto3.Session(
        region_name=region,
        aws_access_key_id=role_credentials['accessKeyId'],
        aws_secret_access_key=role_credentials['secretAccessKey'],
        aws_session_token=role_credentials['sessionToken']
    )
    return session.client('ec2')

# Main function
def main():
    # Read CSV file containing Instance ID column
    input_file = 'instances.csv'  # Update with your input file path
    output_file = 'instance_details.csv'  # Update with your output file path

    access_token = get_sso_access_token()

    with open(input_file, 'r') as csvfile:
        csvreader = csv.reader(csvfile)
        next(csvreader)  # Skip header row
        instances = [row[0] for row in csvreader]

    # Assume IAM role and create EC2 client
    ec2_client = assume_role_and_create_ec2_client(access_token, 'us-west-2')  # Update with your region

    # List to store instance details
    instance_details = []

    # Iterate over each instance ID
    for instance_id in instances:
        # Get network info for instance
        try:
            response = ec2_client.describe_instances(InstanceIds=[instance_id])
            instance = response['Reservations'][0]['Instances'][0]
            vpc_id = instance['VpcId']
            subnet_id = instance['SubnetId']
            instance_details.append([instance_id, vpc_id, subnet_id])
        except Exception as e:
            print(f"Error occurred while retrieving network info for instance {instance_id}: {e}")

    # Write instance details to output CSV file
    with open(output_file, 'w', newline='') as csvfile:
        csvwriter = csv.writer(csvfile)
        csvwriter.writerow(['Instance ID', 'VPC ID', 'Subnet ID'])
        csvwriter.writerows(instance_details)

    print(f"Instance details written to {output_file}")

if __name__ == "__main__":
    main()
