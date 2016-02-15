## Requirement
 1. docker

## Installation Environment
 1. setup local kubenetes single node for experiment.
 	- in setup dir have command script that help you to setup kubernete environment. (/setup/kube_start.sh)
 2. making at least one replication controller from experiment dir with
 	` $kubectl create -f <.yaml, .json file> `
 	(eg. experiment/app/nginx-rc.yaml)
 3. enable docker Remote API
 	- you need to edit `/etc/init/docker.conf` file at `DOCKER_OPTS=` and paste `'-H tcp://0.0.0.0:4243 -H unix:///var/run/` behind it. 
 	
## API
 You can run ThothEyeAPI which stored inside `api` dir (pkg/api/start.sh). The API is running on port 8182, so you can access the API via `localhost:8182` 

read more routes details at 

`https://docs.google.com/document/d/11aDN-w_Ib1Bw0bLuvYmRxVyD-S44UJwT1EsmLDnV9wk/edit`
