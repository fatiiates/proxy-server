# Proxy Server
The proxy exists between the principal services and the clients. It is a methodology generally used for various operations such as protection from online monitoring and load balancing.

# How to work?

The reverse proxy receives requests from the client and directs them to the services running in the background. When the redirect is complete, the response is returned to the client.

The reverse proxy requires a 'config.yaml' in the following format

    proxy:
    method: string
    listen:
        address: string 
        port: int
    services:
        - name: string
        domain: string
        hosts:
            - address: string
            port: int
            - address: string
            port: int

You can find a sample config.yaml [here](proxy/example-config.yaml).

## Testing

After downloading the repository, first build it with docker-compose and get the images

    docker-compose build

Configure the example-config.yaml file according to you. It will also meet your wishes in its default configurations.

If you are going to add a new service instance, the addresses must be in ascending order and should be suitable for the subnet you will give. For example:

    hosts:
        - address: "10.0.0.3"
          port: 5000
        - address: "10.0.0.4"
          port: 5000
        - address: "10.0.0.5"
          port: 5000
        - address: "10.0.0.6"
          port: 5000
        - address: "10.0.0.7"
          port: 5000
        - address: "10.0.0.8"
          port: 5000

Afterwards, you can raise the shuttles directly.

    docker-compose up -d

You must scale the summation-service to the number of services you defined in the example-config.yaml file. If the default file is 2 services, you should scale it by 2.

    docker-compose up --scale service=2 -d

[NOTE] If you make a request to localhost:5000 after starting the containers, it will not work because the host will be localhost. For example, you need to add a record like the one below to the /etc/hosts directory.

    127.0.0.1	proxy-test.fatiiates.com

Finally, the proxy is now ready to receive requests. You can send a request to the server with the following format.

    curl -X POST proxy-test.fatiiates.com:5000/sum\?n1=4\&n2=16 \
    -H 'Content-Type: application/json' \
    -H 'Accept: application/json'

![image](https://user-images.githubusercontent.com/51250249/164884714-7b44642c-c735-4a45-87ec-a32b5c855c94.png)

# REFERENCES

- Main ref: https://www.youtube.com/watch?v=nDZO5m2ExBE 

- https://dev.to/b0r/implement-reverse-proxy-in-gogolang-2cp4

- https://gobyexample.com/atomic-counters  

- https://avinetworks.com/glossary/round-robin-load-balancing/

# FUTURE WORKS

- Caching mechanism is definitely a must, but I think it should be a centrally managed mechanism, not in-memory caching.