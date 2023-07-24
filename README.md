# RTMP Live

## What is this?

This repository provides a comprehensive guide and code samples for creating a small streaming platform based on the **RTMP** (Real-Time Messaging Protocol). The platform enables live streaming capabilities and leverages NGINX RTMP for receiving video streams. Additionally, the repository includes functionality to play the recorded videos via HTTP directly from another pool of servers.

Additionally, a service discovery process is included to report the active streams to an API. The API, integrated with Redis, returns the server and manifest path required for playback.

```mermaid
graph LR
    A[Edge] -- Which server should I request video? --> B[API]
    B -- Get server --> C[Redis]
    B -- Response with Origin Server 1 --> A
    A -- Request content --> D[Origin Server 1]
    E[Origin Server 2]
```

## What's the stack behind it?

This small live streaming platform relies on the following projects:

* [`NGINX-RTMP`](https://github.com/arut/nginx-rtmp-module) - the widely, battle-tested and probably the most famous RTMP server
* [`NGINX`](https://www.nginx.com/) - the most used werb server in the world
* [`Lua`](https://www.lua.org/) - a simple yet very powerful programing language 🇧🇷
* [`Go`](https://go.dev/) - a good language to build HTTP APIs, workers, daemons and every kind of distribued system service

## Edge - CDN

The Edge server, often referred to as "the frontend server" is an essential component of the Content Delivery Network (CDN). It plays a crucial role in the media streaming platform, facilitating a seamless viewing experience for users.

It is the server delivered by the platform you are using to watch the video, this is the server your media player will use to play the video.

The Edge server serves as the intermediary between the end-users and the video content they wish to watch. When you access a video on the platform, your media player interacts with the Edge server, which efficiently delivers the video content to your device. This playable URL comes through an HTTP API and it is out of the scope of this educational project.

Our Edge component here is responsible for retrieving from an HTTP API which Origin server the content will come from, and stick to it through the playback.

```mermaid
graph LR

style Edge fill:#FFC107, stroke:#000, stroke-width:2px, r:10px
style API fill:#4CAF50, stroke:#000, stroke-width:2px, r:10px
style Origin fill:#2196F3, stroke:#000, stroke-width:2px, r:10px

Edge("Edge") -- Which servers holds the content? --> API("HTTP API")
API -- Returns JSON with data --> Edge
Edge -- Proxies Request --> Origin("Origin Server")
```

A typical response from the HTTP API looks like this:

```json
{
    "name": "golive",
    "manifest": "index.m3u8",
    "host": "127.0.0.1"
}
```

We use these values in the [proxy_pass directive](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/) to proxy the request to the correct origin server.


```nginx
location ~ "/(?<stream>[^/]*)/index.m3u8$" {
    set $target "";

    access_by_lua_block {
        ngx.var.target = router.fetch_streams(ngx.var.stream)
    }

    proxy_pass http://$target;
}
```

## Origin

The Origin is the component responsible for receiving the video (ingest), storing and serving the original video content to the Edge servers and users.

Keys characteristics of the Origin service are:

* Ingest: it receives the video feed from an encoder, such as [Elemental](https://aws.amazon.com/elemental-server/) or [OBS Studio](https://obsproject.com/), serving as the entry point for content upload.
* Packager: the Origin service packages the video for user consumption, fragments it into segments, and generates [HLS](https://developer.apple.com/streaming/) manifests.
* Storage: in addition to packaging, the Origin service stores all the video content.
* Delivery: as the backbone of content distribution, it acts as an upstream to the Edge servers, efficiently delivering content when requested.

```mermaid
graph TD

style Encoder fill:#B2DFDB, stroke:#000, stroke-width:2px, r:10px
style RTMPLB fill:#FFCC80, stroke:#000, stroke-width:2px, r:10px
style OriginA fill:#BBDEFB, stroke:#000, stroke-width:2px, r:10px
style OriginB fill:#BBDEFB, stroke:#000, stroke-width:2px, r:10px

Encoder("Encoder (e.g., OBS)") --> RTMPLB("RTMP Load Balancer")
RTMPLB --> OriginA("Origin A")
RTMPLB --> OriginB("Origin B")
```

### Discovery

The Discovery service is responsible for tracking and identifying which server holds a specific streaming content. This becomes especially important when multiple encoders are feeding the Origin Service with different content, and the platform needs to determine the appropriate server(s) to deliver the content when requested by users.

It is deployed aside with the Origin service. The reason is because we continuously need to know whether the video feeding is up and running.

```mermaid
sequenceDiagram
    participant DS as Discovery Service
    participant FS as Filesystem
    participant API as HTTP API

    loop Watch filesystem events
        DS->>FS: Check if manifests are being created/updated
    end

    Note right of DS: Accesses the filesystem to verify if the streaming is working
    DS->>API: Report Host (IP), manifest path, stream name (e.g golive)
    Note right of DS: Sends relevant information to the HTTP API
```

## API

The HTTP API used by the Discovery Service performs two critical functions. Firstly, it enables real-time updates to Redis keys, ensuring the tracking of changes in streaming manifests. By utilizing TTL for Redis keys, the API automatically removes keys when the encoder goes offline or when the live streaming session ends. As a result, the platform stops offering the corresponding live content.
