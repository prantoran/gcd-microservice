#### Notes
* Tutorial link
- https://outcrawl.com/getting-started-microservices-go-grpc-kubernetes
* Compiling proto with grpc
- $ protoc --go_out=plugins=grpc:./ ./pb/*.proto


Building Docker images

Now that your services are ready, you need to containerize them. Create Dockerfiles, one for each service, inside your project's root directory. Dockerfiles can't be located in a subdirectory, because Docker's COPY command can't copy files from the parent directory, which is needed in this example.



Build both images. If you switched to Minikube's Docker daemon, they will become available inside the VM.

$ docker build -t local/gcd -f Dockerfile.gcd .
$ docker build -t local/cli -f Dockerfile.cli .






Deploying to Kubernetes cluster

For each service, you'll need to configure two Kubernetes objects—a deployment and a service.

Explained simply, deployments manage sets of pods to keep the "desired state" of the cluster. Pods are collections of containers. Because they are created and destroyed dynamically, services are needed to provide fixed addresses by which to access them. Which pods are targeted by which services, is determined by label selectors.

Configuration files below are fairly standard. The important parts to keep note of are the ports, names of Docker images built earlier, and labels. Also, imagePullPolicy is set to Never to ensure Kubernetes uses images built locally.

Create gcd.yaml file. It declares a deployment and a service in the same file. Another practice is to separate them into gcd-deployment.yaml and gcd-service.yaml.




Create api.yaml file. The service type is set to NodePort, which makes this service accessible outside of the cluster. For the GCD service, this is set to a default value of ClusterIP, which means a "cluster-internal" IP address.



To create these resources inside the cluster, run the following commands.

$ kubectl create -f cli.yaml
$ kubectl create -f gcd.yaml




Check if all pods are running. By specifying -w flag, you can watch for changes.

$ kubectl get pods -w
NAME                             READY     STATUS    RESTARTS   AGE
api-deployment-778049682-3vd0z   1/1       Running   0          3s
gcd-deployment-544390878-0zgc8   1/1       Running   0          2s
gcd-deployment-544390878-p78g0   1/1       Running   0          2s
gcd-deployment-544390878-r26nx   1/1       Running   0          2s



### kubectl
* deleting service
- kubectl delete service servicename
* deleting deployment
- kbuectl delete deployment deployname 




As set in the configuration files, API service runs on a single pod and the GCD service runs on 3.

Get the URL of the CLI service.
$ minikube service cli-service --url



Finally, try it out.
$ curl http://192.168.99.100:32602/gcd/294/462




