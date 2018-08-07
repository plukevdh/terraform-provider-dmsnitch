# Terraform [Dead Man's Snitch](https://deadmanssnitch.com/) Provider

## Requirements

- Terraform 0.10.x or higher
- Go 1.8 or higher

## Building The Provider

Clone the provider source code

```sh
$ mkdir -p $GOPATH/src/github.com/plukevdh; cd $GOPATH/src/github.com/plukevdh
$ git clone https://github.com/plukevdh/terraform-provider-dmsnitch.git
```

Build the source into executable binary

```sh
$ cd $GOPATH/src/github.com/plukevdh/terraform-provider-dmsnitch
$ make build
```

## Usage

First create an API token in the Dead Man's Snitch dashboard.[^1]

![](http://img.plukevdh.me/0M2i1K2n2T1a/Image%2525202018-08-07%252520at%2525203.45.04%252520PM.png)

Copy this key and configure the DMS provider:

```hcl-terraform
provider "dmsnitch" {
  api_key = "${var.dms_key}"
}
```

Optionally, this can be configured using the envvar `DMS_TOKEN` to avoid storing the token in plaintext config.

Then you can create and manage your DMS snitches like so:

```hcl-terraform
resource "dmsnitch_snitch" "mysnitch" {
  name = "My Important Service"
  notes = "Description or other notes about this snitch."
  
  interval = "daily" 
  type = "basic"
  tags = ["one", "two"]
}
```

You then might use this resource to perform database backup event check-ins with another provider, such as AWS RDS events.

```hcl-terraform
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