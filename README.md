
# Web Scraping Server

A suite of endpoints to start the scraping process to different website


## Clone the project

To clone this project run

```bash
  git clone git@github.com:asmejia1993/web-scraping-server.git
```

Move to directory

```bash
  cd web-scraping-server
```

Run the server

```bash
  docker-compose -d up
```
## API Reference

#### Get all items

```http
  GET /api/v1/web-scraping
```

#### Get item

```http
  GET /api/v1/web-scraping/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### Create the franchises

```http
  POST /api/v1/web-scraping
```

### Architecture

![Alt text](/docs/architecture.PNG "Architecture Diagram")