<h6 align="center">here is Mehran, start coding --</h6>

<p align="center">
<img src="./assets/jeddi.png" height="250px">
</p>

# rate limiter



## Introduction
This project is a rate limiter,
a critical component in a distributed system.
It's designed to control the rate of actions that can be performed by a user or a system.
The main goal is to prevent any part of the system from being overwhelmed by too many requests.

We use sliding-window algorithm to implement rate limiting. 
There is also a middleware that imposes a monthly quota on each user. 
This can be used as an entry for a data-processing pipeline. 

## Installation 

first you should have docker installed in your system.
you can follow instruction on this [link](https://docs.docker.com/engine/install/).

after that build and run the project image running this command
```bash
docker compose up -d .  
```


## Getting started
first you need to signup:
```bash
curl --request POST \
  --url http://localhost:8080/signup \
  --header 'Content-Type: application/json' \
  --data '{
    "month_size_limit": 100,
    "minute_rate_limit": 4
    }
```
you will see response like this:(save the ID for the next request)
```bash
{
	"ID": "d8fff3b4-388d-48f3-82c3-81924c0de5b4",
	"month_size_limit": 100,
	"minute_rate_limit": 4
}
```
after that you can send your data as follows
```bash
curl --request GET \
  --url http://localhost:8080/hello \
  --header 'Data-ID: id_kscdsdsdc' \
  --header 'Data-Size: 2' \
  --header 'User-ID: d8fff3b4-388d-48f3-82c3-81924c0de5b4'
```
passing the middlewares,
you will access to api which here is a "HI".
```bash
HI
```
