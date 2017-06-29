# About
Use Kubernetes injected pod information as environment variables in a scratch container (i.e. exec environment without a shell interpreting envvars). 
Also, gracefully shuts down the http process when the pod gets a SIGTERM.  

Includes several http endpoints, e.g. /debug as well as /healthz and /ready, to test Kubernetes behaviour (liveness and readiness probes).  
Can also be used to demo zero-downtime rolling-updates (use ready_probe_deploy.yaml in v1 folder, then `kubectl rolling-update probes --image=embano1/probes:2.0`).  

# Build (binary, docker image) with:  
`make all` (make sure to tweak Makefile variables)  

# Or directly run (from my container image)
Obiously, have a working kubernetes cluster (e.g. minikube), then:  
`kubectl create -f examples/`  

Tested with Kubernetes v1.6.x. You might run into trouble with older API servers due to the API deployment spec change (extensions to apps/v1beta1). The yaml files have a commented line to easily switch. 