# Adding a new app

Each app maps to one Kubernetes namespace and lives in a single Go file under `charts/apps/`. The file defines a constructor that uses the `kubehelpers` builder API to declare the desired resources. cdk8s synthesizes those into a `dist/<name>.k8s.yaml` file, which Flux CD then reconciles against the cluster.

## Minimal example — a simple stateless app

`charts/apps/excalidraw.go` is the simplest possible app: one deployment, one ingress.

```go
package apps

import "git.mkz.me/mycroft/k8s-home/internal/kubehelpers"

func NewExcalidrawChart(builder *kubehelpers.Builder) *kubehelpers.Chart {
    name := "excalidraw"

    chart := builder.NewChart(name)
    chart.NewNamespace(name) // required — sets chart.Namespace used by all helpers

    labels := map[string]string{
        "app.kubernetes.io/name":      name,
        "app.kubernetes.io/component": "main",
    }

    chart.NewDeployment(&kubehelpers.Deployment{
        Name:            name,
        Labels:          labels,
        Image:           builder.RegisterContainerImage("excalidraw/excalidraw"),
        ImagePullPolicy: "Always", // omit to use the default (IfNotPresent)
    })

    chart.NewIngress(&kubehelpers.Ingress{
        Name:   name,
        Labels: labels,
        Port:   80,
        Ingresses: []string{
            "excalidraw.services.mkz.me",
        },
    })

    return chart
}
```

`NewIngress` creates both the `Service` and the `Ingress` resource. TLS via cert-manager and the Traefik ingress class are applied automatically.

## App with environment variables and secrets

`charts/apps/memos.go` shows how to pass plain env vars and pull a secret from Vault via External Secrets Operator.

```go
env := []kubehelpers.EnvEntry{
    // plain value
    {Name: "MEMOS_DRIVER", Value: kubehelpers.EnvValue{Value: "postgres"}},
    // value read from a Kubernetes Secret key
    {Name: "MEMOS_DSN", Value: kubehelpers.EnvValue{
        ValueFromSecret: kubehelpers.EnvValueFromSecret{
            Name: "postgres", // Secret name
            Key:  "dsn",      // key within the Secret
        },
    }},
}

// Provision the SecretStore (points to Vault) and the ExternalSecret that
// fetches the "postgres" secret from path secret/namespaces/<namespace>/postgres
chart.CreateSecretStore(namespace)
chart.CreateExternalSecret(namespace, "postgres")

chart.NewDeployment(&kubehelpers.Deployment{
    Name:   name,
    Image:  "neosmemo/memos",
    Labels: labels,
    Env:    env,
})
```

The Vault path convention is `secret/namespaces/<namespace>/<secret-name>`.

## App with a Redis sidecar and Prometheus monitoring

`charts/apps/useless.go` adds a Redis StatefulSet and a Prometheus `ServiceMonitor`.

```go
// Provision Redis; returns the StatefulSet name and the ClusterIP service name.
_, redisServiceName := chart.NewRedisStatefulsetWithOpts(namespace, kubehelpers.RedisOpts{
    AppPort: 6379, // optional, defaults to 6379
})

env := []kubehelpers.EnvEntry{
    {Name: "REDIS_HOST", Value: kubehelpers.EnvValue{
        Value: fmt.Sprintf("%s.%s", redisServiceName, namespace),
    }},
    {Name: "REDIS_PORT", Value: kubehelpers.EnvValue{Value: "6379"}},
}

chart.NewDeployment(&kubehelpers.Deployment{
    Name:            name,
    Labels:          labels,
    Env:             env,
    Image:           appImage,
    ImagePullPolicy: "Always",
})

// Creates a ServiceMonitor picked up by kube-prometheus-stack.
// ServicePortName defaults to "http".
chart.NewServiceMonitor(&kubehelpers.ServiceMonitor{
    Labels: labels,
})
```

## Registering the app

Add the constructor to the slice in `charts/homelab.go`:

```go
charts_apps.NewMyAppChart,
```

Comment it out to disable the app without deleting its definition.

## Adding the container image version

If the app uses a container image (not a Helm chart), add it to `versions.yaml`:

```yaml
images:
  myorg/myapp: 1.2.3
```

Then reference it in the constructor:

```go
image := builder.RegisterContainerImage("myorg/myapp")
```

`RegisterContainerImage` looks up the version from `versions.yaml` and returns the full `image:tag` string. The `check-versions` command uses this registry to detect outdated images.

## Generating the manifest

```sh
mise run generate
```

This compiles the Go code and synthesizes `dist/<name>.k8s.yaml`. Inspect the output to verify the generated resources look correct before committing.
