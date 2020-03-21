# dflmon

## Installation

### Docker

See [DockerHub](https://hub.docker.com/r/duffleman/dflmon). The docker file as no requirements, so deploy as you need with the following environment variables:

- `DFLMON_CACHET_URL` - URL for the cachet status page
- `DFLMON_CACHET_TOKEN` - API token for the cachet status page
- `DFLMON_CONFIG` - file path or raw JSON for the jobs to handle

### Config

Here is an example config. "Interval" is always in seconds. You can save the config as minified JSON as insert it as raw JSON via an environment variable called `DFLMON_CONFIG`, or save it to a file and pass the path to find that file as an env var called `DFLMON_CONFIG`.

```json
{
	"jobs": [
		{
			"name": "radarr",
			"component_name": "Radarr",
			"type": "https",
			"host": "radarr.int.dfl.mn",
			"interval": 5
		},
		{
			"name": "vpn_tunnel",
			"component_name": "VPN Bridge",
			"type": "icmp",
			"host": "192.168.254.254",
			"interval": 5
		}
	]
}
```

So far only these two forms are supported:

- icmp
- https
