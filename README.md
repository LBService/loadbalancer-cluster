# LoadBalancer-Cluster

## Context

With IaaS and PaaS development In DC, our data center typically has VM clusters managed by Openstack/VMWare, container clusters managed by Paas Plateform and bare-metal clusters. 

Public Cloud already has complete solution including loadbalancer services. But in private cloud, we have not yet a good (standard) solution, (HardWare F5 is too expensive and not scalable), lvs nginx are goold software, but setup them by hand needs lots of work in a cloud context. 

Building a management software for loadbalancer service to meet automatic and scalability functoins in order to relieve ops engineers and business application engineers from burden is a real necessity 

This project aimed to offer layer 4 tcp & udp and layer 7 http & https loadbalancing service managed by Kubernetes.

## Status

**In progress**

## Architect


## Getting start
