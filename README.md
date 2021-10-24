# 🛳️  XERXES
 XERXES is a cli wrapper written around docker and kong api gateway to make containers accessible over the gateway with declarative configuration;

## Why XERXES ?

The main goal of this tool is to provide an easy way to configure and run micro web backends without having to configure complex tools; 

Xerxes uses existing docker platform to run containers and also configure their discovery using the kong gateway this enables faster simulation and development of complex architectures locally.  
 
Since the platform is independent of underlying container runtime, it can be replaced with [Kubernetes](https://kubernetes.io) or  [Nomad](https://www.nomadproject.io); 

## Configuration

Xerxes primarily uses three config files (JSON) to run all of which should be located in `$HOME/.orchestrator/configuration` directory and should be accessible by the user running the cli

- **config.json :**
This JSON file contains the configuration to the resources needed by the cli to run such as the KONG Admin url , Datastore Path etc. A sample config.json  
```json
{
  "config": {
    "registry": {
      "AWS": {
        "AWS_ACCESS_KEY": "XXXXX",
        "AWS_SECRET_ACCESS_KEY": "XXXX",
        "AWS_REGION": "ap-south-1"
      }
    },
    "kong": {
      "host": "http://ip-of-kong-host",
      "admin": "8001"
    }
    "bitConf": {
      "dbpath": "/home/ubuntu/xerxesDB/data",
      "max_write_size": 1
    }
  }
}
```
- **host.json :** 
   This config file contains the list of available nodes (AWS EC2 Instances) which have docker client installed and is accessible over TCP://2375 , Our EC2 Images repo already contains an image pre-configured with docker . A sample host.json file 
```json
   {
  "nodes": {
    "node01": {
      "host": "tcp://ip-of-node:2375 //host over which docker server is accessible"
      "ip": "ip-of-node  // ip to be routed to from gateway",
      "version": "1.40 // docker server API version"
    }
  }
```
  
- **service.json :**
   This file contains all the configs of services and how they are configured at gateway 
   The current service.json file in production is 
```json
   {
  "services": {
    "test": {
      "image": "mysampleimage:latest //docker image to be used on the node",
      "container_port": 5000, 
      "base_port": 3000,
      "max_port": 4500, 
      "kong_conf": { 
        "service": { 
          "name": "test-service", 
          "route": "/hello-on-kong", 
          "target_path": "/" 
        },
        "upstream": { 
          "name": "sample.v1.test",
          "hashon": "none" 
        }
      },
      "health": {
        "endpoint" : "/health" 
      }
    }
  }
}
```

## Deploying a service  
  
### Step `1`  
The first step is to add the configuration for the service into the **service.json** file  
In this tutorial we'll be creating a new service named `testobj` , the config object for this service is  
```json
"testobj": {  
      "image": "testimg:v0.0.1",  
      "container_port": 5000,  
      "base_port": 3000,  
      "max_port": 4500,  
      "kong_conf": {   
        "service": {   
          "name": "test-service",  
          "route": "/test",   
          "target_path": "/"  
        },  
        "upstream": {   
          "name": "test.v1.service",  
          "hashon": "none"   
        }  
      },  
      "health": {  
        "endpoint" : "/health"  
      }  
}  
```  
the above given config does the following  
  
- cli can interact with this service using the given key `testobj`  
- port forwarding will work in the following way in  docker `0.0.0.0:[3000-4500]->5000/tcp`  
- all the requests containing path `/test` will be forwarded to `/` of this service  
- And the rest is the kong config  
  this config can now be appended to the file `service.json` file  
  
### step `2`  
  
Now we are ready to use the cli , the service can be accessed using the `testobj` keyword.  
  
- check if the config is loaded correctly  
  `./build/xerxesv3 service def[intition] testobj`  
  this should output the config we've just added  
- list all running services , `./build/xerxesv3 service ls testobj`  
- scale service to 1 ,  
  `./build/xerxesv2 service scale --number 1 testobj`  
  
- alternatively you can also specify on which node the service should run `./build/xerxesv3 service scale --number 1 --node m01 testobj`  
- `./build/xerxesv3 service ls testobj` now should output the running service  
  
```green
+----------------------+--------------+-----------+---------------+------+  
|          ID          |  CONTAINERID | HOST NODE |      IP       | PORT |  
+----------------------+--------------+-----------+---------------+------+  
| buqd6pobdd79giq9stf0 |efe0a250aadce0| node01    | 111.31.11.251 | 4249 |  
+----------------------+--------------+-----------+---------------+------+  
```  
- Now the command `./build/xerxesv3 flake logs buqd6pobdd79giq9stf0` can be used to view logs of this service  
- to scale more the same command can be used with changed args  `./build/xerxesv3 service scale --number 3 --node m01 testobj `  
- to shutdown we can scale to 0 or run `./build/xerxesv3 flake remove buqd6pobdd79giq9stf0`

## Configuration

Xerxes also comes bundled with a [scheduler](https://github.com/rahul0tripathi/xerxes/tree/master/scheduler) which is responsible for the health checks and deleting a container;

along with this the scheduler also exposes [GRPC](https://github.com/rahul0tripathi/xerxes/blob/master/proto/containerManager.proto) endpoint to query the container metadata, logs and the running services.

## Hyperlinks

 - [KONG](https://docs.konghq.com/2.2.x/)

created by [@rahul0tripathi](https://github.com/rahul0tripathi)
