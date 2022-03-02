# mg

## Supported platforms

 * Windows
 * macOS
 * Linux
 * Docker container
 * Helm chart for Kubernetes

## How to use

### Desktop binaries

Download the binary for your OS and arch from Releases page. Run it from command line.

I.e. `./mg -s https://api.itemstolist.top/api/target -t endpoint`

### Docker image

 * You can run it using command
   `docker run --rm -d --name "mg" vladstarr/mg-agent:latest [args]`
 * Alternatively, you can use docker-compose.yaml file from this repo.
   Run it using `docker-compose up -d`

### Helm chart

 * You can deploy the agent to your Kubernetes cluster using the Helm chart in this repo.
 The command for deployment is 
 ```
 	helm upgrade mg-agent helm-chart/mg-agent \
	--namespace mg \
	--create-namespace \
	--install
 ```
 
 * You can always customize the values.yaml of the chart for your needs. `agent` section options are treated as argument for the go binary.

### Available flags for binaries:

 * `-n`: int, number of requests per each target (default 1000000)
 * `-s`: string, url to endpoint or file name (default "ukraine.txt")
 * `-t`: string, source type to use (file or endpoint) (default "file")
 * `-w`: int, number of workers per logical CPU (default 10)

### Target model for endpoint to consume:

```json
{
	"id": 0,
	"url": "https://example.ru",
	"method": "GET",
	"data": null,
	"headers": [["<header-key>", "<header-value>"]],
	"proxy_url": null
}
```
