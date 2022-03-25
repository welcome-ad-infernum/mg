# Kubernetes agent deployment

## Prerequisites

* [helm](https://helm.sh) needs to be installed on your PC
* [kubectl](https://kubernetes.io/docs/tasks/tools/) needs to be installed on your PC

## Deploy the agent using Helm chart

The command for deployment is:

```
$ helm upgrade mg-agent examples/helm-chart/mg-agent \
--namespace mg-agent \
--create-namespace \
--install
```

## Verify that agent is running correctly

You can view agent's logs with:

```
$ kubectl -n mg-agent get pod
$ kubectl -n mg-agent logs -f <pod_name>
```

## Customization

You can always customize the `values.yaml` of the chart for your needs. 

`agent` section values are treated as arguments for the agent binary. See `templates/deployment.yaml` for details.