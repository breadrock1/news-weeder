# News-Weeder Project

## Overview

There is wrapper service for `redis-stack-server` service to store and search similar news 
articles by content loaded by news-rss and news-telegram projects.

## Quick Start

1. Clone the repository:

    ```shell
    git clone http://<gitlab-domain-address>/data-lake/news-weeder.git
    cd news-weeder
    ```

2. Build docker image from sources:

    ```shell
   docker build -t news-weeder:latest .
    ```

3. Edit configs file `configs/production.toml` to launch docker compose services:

   The main parameter is `redis.address` - you have to pass redis server address

4. Start the application using Docker Compose:

    ```shell
    docker compose up -d news-weeder <other-needed-services>
    ```

5. The application should now be running. Check the logs with:

    ```shell
    docker compose logs -f
    ```
