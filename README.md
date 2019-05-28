
# Docker Registry Manager [![Go Report Card](https://goreportcard.com/badge/github.com/snagles/docker-registry-manager)](https://goreportcard.com/report/github.com/snagles/docker-registry-manager) [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many).

| Service   |  Master  | Develop  |   
|---|---|---|
| Status   | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master)  | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=develop)   |
| Coverage  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/develop/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  |

![Example](https://github.com/snagles/resources/blob/master/docker-registry-manager-updated.gif)

## Current Features
 1. Support for docker distribution registry v2 (https and http)
 2. Viewable image/tags stages, commands, and sizes.
 3. Bulk deletes of tags
 4. Registry activity logs
 5. Comparison of registry images to public Dockerhub images

## Planned Features
 1. Authentication for users with admin/read only rights using TLS
 2. Global search
 3. List image shared layers
 4. Event timeline

## Quickstart
 The below steps assume you have a docker registry currently running (with delete mode enabled (https://docs.docker.com/registry/configuration/). To add a registry to manage, add via the interface... or via the registries.yml file

### Docker-Compose (Recommended)
 Install compose (https://docs.docker.com/compose/install/), and then run the below commands

```bash
rm -fr ./docker-registrees
git clone https://github.com/Jean-Baptiste-Lasselle/docker-registrees.git && cd docker-registrees 

mkdir -p ./better

if [ -f ./better/registries.yml ]; then  rm -f ./better/registries.yml;fi;

echo "registries:" >> ./better/registries.yml
echo "  localRegistry:" >> ./better/registries.yml
echo "    displayname: registry.example.pour.vinse.com:5000" >> ./better/registries.yml
echo "    url: http://localhost # Example https://localhost, http://remotehost.com" >> ./better/registries.yml
echo "    port: 5000  # Example: 443, 8080, 5000" >> ./better/registries.yml
echo "    username: exampleUser" >> ./better/registries.yml
echo "    password: examplePassword" >> ./better/registries.yml
echo "    refresh-rate: "5m" # Example: 60s, 5m, 1h" >> ./better/registries.yml
echo "    skip-tls-validation: true # REQUIRED for self signed certificates" >> ./better/registries.yml
echo "    dockerhub-integration: true # Optional - compares to dockerhub to determine if image up to date" >> ./better/registries.yml
docker-compose up -d
firefox localhost:8081
```

#### Environment Options:
 - MANAGER_PORT: Port to run on inside the docker container
 - MANAGER_REGISTRIES: Registries.yml file location inside the docker container
 - MANAGER_LOG_LEVEL: Log level for logs (fatal, panic, error, warn, info, debug)
 - MANAGER_ENABLE_HTTPS: true/false for using HTTPS. When using HTTPS the below options must be set
 - MANAGER_KEY: key file location inside the docker container
 - MANAGER_CERTIFICATE: Certificate location inside the docker container

### Go
 ```bash
    git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
    vim registries.yml # add your registry
    cd app && go build . && ./app --port 8080 --log-level warn --registries "../registries.yml"
    firefox localhost:8080
 ```

#### CLI Options
  - port, p: Port to run on
  - registries, r: Registrys.yml file location
  - log-level, l: Log level for logs (fatal, panic, error, warn, info, debug)
  - enable-https, e: true/false for using HTTPS. When using HTTPS the below options must be set
  - tls-key, k: key file location inside the docker container
  - tls-certificate, cert: Certificate location inside the docker container

### Dockerfile
 ```bash
    vim registries.yml # add your registry
    docker run --detach --name docker-registry-manager -p 8080:8080 -e MANAGER_PORT=8080 -e MANAGER_REGISTRIES=/app/registries.yml -e MANAGER_LOG_LEVEL=warn docker-registry-manager
    firefox localhost:8080
 ```

#### Environment Options:
- MANAGER_PORT: Port to run on inside the docker container
- MANAGER_REGISTRIES: Registries.yml file location inside the docker container
- MANAGER_LOG_LEVEL: Log level for logs (fatal, panic, error, warn, info, debug)
- MANAGER_ENABLE_HTTPS: true/false for using HTTPS. When using HTTPS the below options must be set
- MANAGER_KEY: key file location inside the docker container
- MANAGER_CERTIFICATE: Certificate location inside the docker container

### Registries.yml Example
```yml
registries:
  localRegistry:
    displayname: registry.example.com:5000
    url: http://localhost # Example https://localhost, http://remotehost.com
    port: 5000  # Example: 443, 8080, 5000
    username: exampleUser
    password: examplePassword
    refresh-rate: "5m" # Example: 60s, 5m, 1h
    skip-tls-validation: true # REQUIRED for self signed certificates
    dockerhub-integration: true # Optional - compares to dockerhub to determine if image up to date
```

# ANNEXE 1

## Pushing an image to registry, from remote machine

Le test ci-dessous, démontre que docker interdit une connexion HTTP , et oblige une
connexion `SSL/TLS` `HTTPS`, pour faire un `docker push` (cf. la doc. officielle Docker,  [`Run an externally-accessible registry`](https://docs.docker.com/registry/deploying/#run-an-externally-accessible-registry) ) :

```bash
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker run -d -p 5000:5000 --name awx-registry -e REGISTRY_DELETE_ENABLED=true registry:2
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker ps -a
CONTAINER ID        IMAGE                             COMMAND                  CREATED             STATUS              PORTS                                            NAMES
3e06d249d86a        registry:2                        "/entrypoint.sh /etc…"   4 seconds ago       Up 3 seconds        0.0.0.0:5000->5000/tcp                           awx-registry
a745392b8ff1        snagles/docker-registry-manager   "/app/app"               7 minutes ago       Up 7 minutes        0.0.0.0:5001->5000/tcp, 0.0.0.0:8081->8080/tcp   docker-registry-manager
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker ps -a
CONTAINER ID        IMAGE                             COMMAND                  CREATED             STATUS              PORTS                                            NAMES
3e06d249d86a        registry:2                        "/entrypoint.sh /etc…"   23 minutes ago      Up 23 minutes       0.0.0.0:5000->5000/tcp                           awx-registry
a745392b8ff1        snagles/docker-registry-manager   "/app/app"               30 minutes ago      Up 30 minutes       0.0.0.0:5001->5000/tcp, 0.0.0.0:8081->8080/tcp   docker-registry-manager
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker pull centos:8
Error response from daemon: manifest for centos:8 not found
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker pull centos:6
6: Pulling from library/centos
ff50d722b382: Pull complete
Digest: sha256:dec8f471302de43f4cfcf82f56d99a5227b5ea1aa6d02fa56344986e1f4610e7
Status: Downloaded newer image for centos:6
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker images
REPOSITORY                        TAG                 IMAGE ID            CREATED             SIZE
centos                            6                   d0957ffdf8a2        2 months ago        194MB
registry                          2                   f32a97de94e1        2 months ago        25.8MB
hello-world                       latest              fce289e99eb9        4 months ago        1.84kB
snagles/docker-registry-manager   latest              9e75f9744e04        7 months ago        23MB
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker tag d0957ffdf8a2 192.168.1.22:5000/pegasus/centos:6
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push 192.168.1.22:5000/pegasus/centos:6
The push refers to repository [192.168.1.22:5000/pegasus/centos]
Get https://192.168.1.22:5000/v2/: http: server gave HTTP response to HTTPS client
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker images
REPOSITORY                         TAG                 IMAGE ID            CREATED             SIZE
192.168.1.22:5000/pegasus/centos   6                   d0957ffdf8a2        2 months ago        194MB
centos                             6                   d0957ffdf8a2        2 months ago        194MB
registry                           2                   f32a97de94e1        2 months ago        25.8MB
hello-world                        latest              fce289e99eb9        4 months ago        1.84kB
snagles/docker-registry-manager    latest              9e75f9744e04        7 months ago        23MB
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push d0957ffdf8a2
The push refers to repository [docker.io/library/d0957ffdf8a2]
An image does not exist locally with the tag: d0957ffdf8a2
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker pull 192.168.1.22:5000/pegasus/centos:6
Error response from daemon: Get https://192.168.1.22:5000/v2/: http: server gave HTTP response to HTTPS client
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push http://192.168.1.22:5000/pegasus/centos:6
invalid reference format
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push 192.168.1.22:5000/pegasus/centos:6
The push refers to repository [192.168.1.22:5000/pegasus/centos]
Get https://192.168.1.22:5000/v2/: http: server gave HTTP response to HTTPS client
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push localhost:5000/pegasus/centos:6
The push refers to repository [localhost:5000/pegasus/centos]
An image does not exist locally with the tag: localhost:5000/pegasus/centos
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker tag d0957ffdf8a2 localhost:5000/pegasus/centos:6
jbl@poste-devops-typique:~/docker-registrees/registrees$ docker push localhost:5000/pegasus/centos:6
The push refers to repository [localhost:5000/pegasus/centos]
af6bf1987c2e: Pushed
6: digest: sha256:9aae95c8043f4e401178d68006756dc68982ae6d0693b71a714754227ce0abc6 size: 529
jbl@poste-devops-typique:~/docker-registrees/registrees$

```

* The configuration file that the private registry docker relies on, is (le fichier de configuration du registry docker privé est) :

```bash
docker exec -it awx_dock_registry sh -c "cat /etc/docker/registry/config.yml"
```
* The configuration that worked for me, which was required, even if using the `REGISTRY_DELETE_ENABLED=true|false` envrionment variable in the `docker-compse.yml`, for the docker registry manager app to be able to delete `oci` images by digests  (Et la configuration qui a fonctionné, pour que le `docker_registry_manager` puisse supprimer des `digest`, c-a-d. des versions d'images docker) :

```YAML
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
  delete:                                                                                                                                                                             
    enabled: true
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
```

* A more refined configuration, infering a perisstence of the images in an external  database (Une configuration plus étoffée, pour persister les données de registry; les images docker, dans une base de données dans un conteneur distinct du conteneur du `private docker registry` ) :

```YAML
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
  delete:                                                                                                                                                                             
    enabled: true
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
```
