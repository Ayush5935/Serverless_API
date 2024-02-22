import boto3
import csv

# Regions to search
regions = ["us-west-1", "us-west-2", "us-east-1", "us-east-2"]

# Function to get instances matching the specified naming pattern
def get_instances(region):
    ec2_client = boto3.client("ec2", region_name=region)
    response = ec2_client.describe_instances(
        Filters=[
            {"Name": "tag:Name", "Values": ["CS-PE*", "CS-R*"]}
        ]
    )
    instances = response["Reservations"]
    return instances

# Function to extract required information
def extract_info(instances):
    instance_data = []
    ec2 = boto3.resource('ec2')
    for reservation in instances:
        for instance in reservation["Instances"]:
            instance_id = instance["InstanceId"]
            instance_name = ''
            for tag in instance["Tags"]:
                if tag["Key"] == "Name":
                    instance_name = tag["Value"]
                    break
            subnet_id = instance.get("SubnetId", "")
            subnet_name = ""
            vpc_id = instance.get("VpcId", "")
            vpc_name = ""
            
            # Get subnet name
            if subnet_id:
                try:
                    subnet = ec2.Subnet(subnet_id)
                    subnet_name = subnet.tags[0]['Value'] if subnet.tags else ""
                except Exception as e:
                    print(f"Error retrieving subnet details for subnet {subnet_id}: {e}")
            
            # Get VPC name
            if vpc_id:
                try:
                    vpc = ec2.Vpc(vpc_id)
                    vpc_name = vpc.tags[0]['Value'] if vpc.tags else ""
                except Exception as e:
                    print(f"Error retrieving VPC details for VPC {vpc_id}: {e}")
            
            instance_data.append([instance_id, instance_name, subnet_id, subnet_name, vpc_id, vpc_name])
    return instance_data

# Function to save data to CSV
def save_to_csv(data):
    with open("aws_instances.csv", "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["Instance ID", "Instance Name", "Subnet ID", "Subnet Name", "VPC ID", "VPC Name"])
        writer.writerows(data)

# Main function
def main():
    all_instance_data = []
    for region in regions:
        instances = get_instances(region)
        instance_data = extract_info(instances)
        all_instance_data.extend(instance_data)
    save_to_csv(all_instance_data)

if __name__ == "__main__":
    main()
