## Requirements

- Terraform 0.13.x or higher
- Go 1.11 or higher
- Go module support

## Setting up the provider

```sh
$ go get -u github.com/repaygithub/terraform-provider-dmsnitch
$ mkdir -p ~/.terraform.d/plugins
$ mv ${GOPATH}/bin/terraform-provider-dmsnitch ~/.terraform.d/plugins
```

## Usage

First create an API token in the Dead Man's Snitch dashboard. **Note**: Use of the DMS API requires [a paid plan](https://deadmanssnitch.com/plans).

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
  alert_email = ["example@example.com", "support@example.com"]
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

### Fields

| Field | Required | Values | Defaults |
|---|---|---|---|
| `name` | yes |
| `notes`| no | | `Managed by Terraform` | 
| `interval` | yes | `15_minute`, `30_minute`, `hourly`, `daily`, `weekly`, `monthly` | `daily` |
| `type` | yes; `smart` is only valid for `weekly` or `monthly` intervals  | `basic`, `smart` | `basic` |
| `alert_email` | no | an array of values |  
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
