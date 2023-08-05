# Middlewares

In this small project, we inject some middlewares to our customized mux.

### Notes :
- Injecting idea is taken from ArdanLabs/Service project.
- httptreemux, uber's zap, and uuid are external dependencies.
- httptreemux's ContextMux is wrapped with custom app struct. In this way, context is added to each request.
- zap is useful for structured logging. Logged JSON has key-value nature. This might be useful while investigating logs.
- Each request has (kinda) unique TraceID which is embedded in request context.  


### References : 
- ArdanLabs/Service project. [REPO](https://github.com/ardanlabs/service/tree/service3)