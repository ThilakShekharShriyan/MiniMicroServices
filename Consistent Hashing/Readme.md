# Consistent Hashing

This project implements **Consistent Hashing** with virtual nodes, exposing an HTTP API to manage and query a distributed hash ring.

---

## What is Consistent Hashing?

**Consistent Hashing** is a technique used to distribute keys across a dynamic set of nodes (like servers or databases) in a way that minimizes re-distribution when nodes are added or removed.

Instead of placing keys directly on nodes, consistent hashing maps both keys and nodes to points on a **circular hash ring**. A key is assigned to the **next node clockwise** on the ring.

---

## Problem with Simple Modulo Hashing

Using modulo (e.g., `hash(key) % N`) for assigning keys to nodes has a critical flaw:

- When you **add or remove a node**, **most keys are remapped**.
- This leads to **huge data movement**, cache invalidation, and poor scalability.

**Consistent Hashing** fixes this by ensuring **only a small portion of keys need to be moved** when nodes change.

---

## Why Is It Required?

Consistent Hashing is essential in distributed systems to:

- Balance load across nodes
- Minimize re-distribution of data
- Handle dynamic scaling (adding/removing nodes)
- Enable fault tolerance and partitioning

---

## Where Is It Used?

Consistent Hashing is widely adopted in distributed systems:

| Technology | Use Case |
|-----------|----------|
| **Redis Cluster** | Distributes keys across shards |
| **Cassandra** | Maps partitions to nodes using consistent hashing |
| **Amazon DynamoDB** | Uses consistent hashing to partition and replicate data |

---

## Main Logic

### `hashring` 

- Implements a **consistent hash ring** with configurable **replicas** (virtual nodes) per physical node.
- Keys and nodes are both hashed using CRC32.
- Supports:
  - Adding nodes
  - Removing nodes
  - Looking up which node a key maps to
  - Listing all active nodes

### `api` 

- A **wrapper around the hashring logic**, exposed as REST APIs.
- Built using [Gorilla Mux](https://github.com/gorilla/mux).

---

## Commands to run after you run the server.


# Add Node
curl -X POST -H "Content-Type: application/json" -d '{"name":"NodeA"}' localhost:8080/nodes

# Remove Node
curl -X DELETE localhost:8080/nodes/NodeA

# Lookup Key
curl localhost:8080/lookup?key=myfile123

# List Nodes
curl localhost:8080/nodes




You can containerize the above by running

```
docker build -t consistent-hashing-api .

docker run -p 8080:8080 consistent-hashing-api

```




