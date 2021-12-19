# Metrics interceptor
In this example, we will try to create fiber server with metrics interceptor enabled.

Metrics interceptor will collect bellow metrics with prometheus data format.
- RPC elapsed
- RPC error count
- RPC response code count

Users need to start a prometheus client locally export the data.
[rk-prom](https://github.com/rookie-ninja/rk-prom) would be a good option start prometheus client easily.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Quick start](#quick-start)
  - [Code](#code)
- [Options](#options)
  - [Override namespace and subsystem](#override-namespace-and-subsystem)
  - [Override Registerer](#override-registerer)
  - [Context Usage](#context-usage)
- [Example](#example)
  - [Start server](#start-server)
  - [Output](#output)
  - [Code](#code-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Quick start
Get rk-fiber package from the remote repository.

### Code
```go
import     "github.com/rookie-ninja/rk-fiber/interceptor/metrics/prom"
```
```go
    // ********************************************
    // ********** Enable interceptors *************
    // ********************************************
	interceptors := []fiber.Handler{
        rkfibermetrics.Interceptor(),
    }
```

## Options
In order to define prometheus style metrics, we need to define <namespace> and <subsystem>.
- namespace: rkentry.GlobalAppCtx().AppName ("rk" will be used by default.)
- subsystem: entryName (Provided as interceptor option. "fiber" will be used by default.)

| Name | Description | Default Values |
| ---- | ---- | ---- |
| rkfibermetrics.WithEntryNameAndType(entryName, entryType string) | Provide entry name and type if there are multiple extension interceptors needs to be used. | rk, fiber |
| rkfibermetrics.WithRegisterer(registerer prometheus.Registerer) | Provide prometheus registerer. | prometheus.DefaultRegisterer |

![arch](img/arch.png)

### Override namespace and subsystem
```go
func main() {
    // Override app name which would replace namespace value in prometheus.
    rkentry.GlobalAppCtx.GetAppInfoEntry().AppName = "newApp"

    // ********************************************
    // ********** Enable interceptors *************
    // ********************************************
	interceptors := []fiber.Handler{
        rkfibermetrics.Interceptor(
            // Add metrics interceptor with entry name and entry type.
            // subsystem would be replaced with newEntry.
            rkfibermetrics.Interceptor(rkfibermetrics.WithEntryNameAndType("newEntry", "fiber")),
        ),
    }

    // 1: Create fiber server
    server := startGreeterServer(opts...)
    ...
}
```

### Override Registerer
```go
	interceptors := []fiber.Handler{
        rkfibermetrics.Interceptor(
            // Provide new prometheus registerer.
            // Default value is prometheus.DefaultRegisterer
            rkfibermetrics.WithRegisterer(prometheus.NewRegistry()),
        ),
    }
```

### Context Usage
| Name | Functionality |
| ------ | ------ |
| rkfiberctx.GetLogger(*fiber.Ctx) | Get logger generated by log interceptor. If there are X-Request-Id or X-Trace-Id as headers in incoming and outgoing metadata, then loggers will has requestId and traceId attached by default. |
| rkfiberctx.GetEvent(*fiber.Ctx) | Get event generated by log interceptor. Event would be printed as soon as RPC finished. |
| rkfiberctx.GetIncomingHeaders(*fiber.Ctx) | Get incoming header. |
| rkfiberctx.AddHeaderToClient(ctx, "k", "v") | Add k/v to headers which would be sent to client. This is append operation. |
| rkfiberctx.SetHeaderToClient(ctx, "k", "v") | Set k/v to headers which would be sent to client. |

## Example
### Start server
```shell script
$ go run greeter-server.go
```

### Output
- Server: localhost:1608/metrics
```shell script
$ curl localhost:1608/metrics
...
# HELP rk_greeter_elapsedNano Summary for name:elapsedNano and labels:[entryName entryType realm region az domain instance appVersion appName restMethod restPath type resCode]
# TYPE rk_greeter_elapsedNano summary
rk_greeter_elapsedNano{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber",quantile="0.5"} 204503
rk_greeter_elapsedNano{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber",quantile="0.9"} 204503
rk_greeter_elapsedNano{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber",quantile="0.99"} 204503
rk_greeter_elapsedNano{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber",quantile="0.999"} 204503
rk_greeter_elapsedNano_sum{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber"} 204503
rk_greeter_elapsedNano_count{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber"} 1
# HELP rk_greeter_resCode counter for name:resCode and labels:[entryName entryType realm region az domain instance appVersion appName restMethod restPath type resCode]
# TYPE rk_greeter_resCode counter
rk_greeter_resCode{appName="rk",appVersion="",az="*",domain="*",entryName="greeter",entryType="fiber",instance="lark.local",realm="*",region="*",resCode="200",restMethod="GET",restPath="/rk/v1/greeter",type="fiber"} 1
```

### Code
- [greeter-server.go](greeter-server.go)