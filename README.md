# Terraform [Dead Man's Snitch](https://deadmanssnitch.com/) Provider

## Note on Compatability

This branch is meant to work against the latest version of Terraform. For previous versions, see below

- 0.0.x: Terraform 0.11 and earlier
- 0.1.x: Terraform 0.12 

## Requirements

- Terraform 0.12.x or higher
- Go 1.11 or higher
- Go module support

Please note [the following details](https://www.terraform.io/docs/extend/terraform-0.12-compatibility.html) if you have built this plugin prior to Terraform 0.12. 

## Setting up the provider

```sh
$ go get -u github.com/plukevdh/terraform-provider-dmsnitch
$ mkdir -p ~/.terraform.d/plugins
$ mv ${GOPATH}/bin/terraform-provider-dmsnitch ~/.terraform.d/plugins
```

## Usage

First create an API token in the Dead Man's Snitch dashboard.[^1]

![](http://img.plukevdh.me/0M2i1K2n2T1a/Image%2525202018-08-07%252520at%2525203.45.04%252520PM.png)

Copy this key and configure the DMS provider:

```hcl
provider "dmsnitch" {
  api_key = "${var.dms_key}"
}
```

Optionally, this can be configured using the envvar `DMS_TOKEN` to avoid storing the token in plaintext config.

Then you can create and manage your DMS snitches like so:

```hcl
resource "dmsnitch_snitch" "mysnitch" {
  name = "My Important Service"
  notes = "Description or other notes about this snitch."
  
  interval = "daily" 
  type = "basic"
  tags = ["one", "two"]
}
```

You then might use this resource to perform database backup event check-ins with another provider, such as AWS RDS events.

```hcl
resource "aws_sns_topic" "backup_event" {
  name = "db-backup-events"
}

resource "aws_db_event_subscription" "backup_event" {
  name      = "db-backup-events"
  sns_topic = "${aws_sns_topic.sbackup_event.arn}"

  source_type = "db-instance"
  source_ids  = ["${aws_db_instance.db.id}"]

  event_categories = [
    "backup",
  ]
}

resource "aws_sns_topic_subscription" "backup_event" {
  endpoint               = "${dmsnitch_snitch.mysnitch.url}"
  protocol               = "https"
  endpoint_auto_confirms = true
  topic_arn              = "${aws_sns_topic.backup_event.arn}"
}
```

You can also import existing snitches using their token found in the snitch's page URL:

`terraform import dmsnitch_snitch.mysnitch 5b025eecf3`

![](http://img.plukevdh.me/1X2N462b0J3a/%255B5a117e75fd66875d1a7c61c65ceaaae3%255D_Image%2525202018-08-07%252520at%2525204.27.59%252520PM.png)

                                     
### Fields

| Field | Required | Values | Defaults |
|---|---|---|---|
| `name` | yes |
| `notes`| no | | `Managed by Terraform` | 
| `interval` | yes | `15_minute`, `30_minute`, `hourly`, `daily`, `weekly`, `monthly` | `daily` |
| `type` | yes; `smart` is only valid for `weekly` or `monthly` intervals  | `basic`, `smart` | `basic` |
| `tags` | no | an array of values | 
 
 ### Attributes

| Attribute | Description |
|---|---|
| `token`, `id` | The unique snitch ID. |
| `url`| The snitch checkin URL (for performing the check-in ping). | 
| `status` | Health status for the snitch. |
  
For additional details about these fields and their purposes, see the [API documentation](https://deadmanssnitch.com/docs/api/v1). 

## Acknowledgements

This codebase is based heavily off of the [Bitbucket Provider](https://github.com/terraform-providers/terraform-provider-bitbucket) codebase.

[^1]: Use of the DMS API requires [a paid plan](https://deadmanssnitch.com/plans).
