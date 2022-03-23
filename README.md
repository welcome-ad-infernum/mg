# mg

## Supported platforms

 * Windows binary
 * macOS binary
 * Linux binary and systemd script
 * Docker container
 * Helm chart for Kubernetes

## How to use

### Desktop binaries

Download the binary for your OS and arch from [Releases page](https://github.com/welcome-ad-infernum/mg/releases). Run the following commands from command line.

```
cd <downloads_folder>
./mg (for linux)
.\mg.exe (for windows)
```

### Linux systemd script

* Open the terminal on your Linux PC or server.
* Install it with:

    `curl -s https://raw.githubusercontent.com/welcome-ad-infernum/mg/main/examples/linux/install.sh | sudo sh -`
* Uninstall it with:

    `curl -s https://raw.githubusercontent.com/welcome-ad-infernum/mg/main/examples/linux/uninstall.sh | sudo sh -`

### Docker image

 * You can run it using command

   `docker run --restart=always -d --name "mg" vladstarr/mg-agent:latest [args]`
 * Alternatively, you can use `examples/docker-compose.yaml` file from our repo and configure for your needs. 

### Helm chart

 * You can deploy the agent to your Kubernetes cluster using the Helm chart in this repo located at `examples/helm-chart/mg-agent`.
 The command for deployment is 
 ```
helm upgrade mg-agent examples/helm-chart/mg-agent \
--namespace mg-agent \
--create-namespace \
--install
 ```
 * You can always customize the values.yaml of the chart for your needs. `agent` section values are treated as arguments for the binary.

### Available flags for binaries:

 * `-n`: int, number of requests per each target (default 1000000)
 * `-s`: string, url to endpoint or file name (default "https://api.itemstolist.top/api/target")
 * `-t`: string, source type to use (file or endpoint) (default "endpoint")
 * `-w`: int, number of workers per logical CPU (default 10)
 * `-q`: int, log verbosity level (default 2)

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

### For non-technical users:

 1. Go to [releases page](https://github.com/welcome-ad-infernum/mg/releases).
 2. Find the latest release.
 3. Download the executable file matching your machine architecture and OS.
 4. Open your terminal (`Terminal` for MacOS, `CMD` or `PowerShell` for Windows)
 6. Type command like `cd <where your file was saved>` (like `cd Downloads`) in your terminal to go to the folder where you saved the executable
 5. Run `./mg` or `.\mg.exe` depending on your platform with flags described above
