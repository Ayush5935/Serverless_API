import csv
import boto3
from time import sleep
import webbrowser

# Function to obtain SSO access token
def get_sso_access_token():
    session = boto3.session.Session()
    region = 'us-west-2'  # Update with your region
    sso_oidc = session.client('sso-oidc', region_name=region)

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

    webbrowser.open(url, autoraise=True)
    print(f"Open the following URL in your browser to authenticate: {url}")
    print(f"Waiting for authentication...")

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

# Function to get VPC ID and Subnet ID for an instance in a specific region
def get_vpc_subnet_id(instance_id, account_id, access_token, region):
    session = boto3.session.Session()
    sso = session.client('sso', region_name=region)

    role_credentials = sso.get_role_credentials(
        roleName='DishWPaaSAdministrator',
        accountId=account_id,
        accessToken=access_token
    )['roleCredentials']

    ec2 = boto3.client(
        'ec2',
        region_name=region,
        aws_access_key_id=role_credentials['accessKeyId'],
        aws_secret_access_key=role_credentials['secretAccessKey'],
        aws_session_token=role_credentials['sessionToken']
    )

    response = ec2.describe_instances(InstanceIds=[instance_id])
    vpc_id = response['Reservations'][0]['Instances'][0]['VpcId']
    subnet_id = response['Reservations'][0]['Instances'][0]['SubnetId']
    return vpc_id, subnet_id

# Main function
def main():
    input_file = 'input.csv'
    output_file = 'output.csv'
    access_token = get_sso_access_token()

    with open(output_file, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Account ID', 'Instance ID', 'Region', 'VPC ID', 'Subnet ID'])

        with open(input_file, mode='r') as csvfile:
            reader = csv.reader(csvfile)
            next(reader)
            for row in reader:
                account_id, instance_id, region = row

                try:
                    vpc_id, subnet_id = get_vpc_subnet_id(instance_id, account_id, access_token, region)
                    writer.writerow([account_id, instance_id, region, vpc_id, subnet_id])
                except Exception as e:
                    print(f"Error processing instance {instance_id} in region {region}: {e}")

    print(f"Results saved to {output_file}")

if __name__ == "__main__":
    main()
