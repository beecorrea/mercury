# API Reference

## Shortlinks

### List Shortlinks

* **Endpoint:** `GET /api/links`
* **Response Status:** `200 OK`
* **Response Body:**
```json
[
  {
    "key": "github",
    "url": "https://github.com",
    "domain": "localhost:45800",
    "created_at": "2026-05-23T19:51:13Z"
  }
]
```

### Create Shortlink

* **Endpoint:** `POST /api/shorten`
* **Request Body:**
```json
{
  "key": "google",
  "url": "https://google.com",
  "domain": "localhost:45800"
}
```
* **Response Status:** `201 Created`
* **Response Body:**
```json
{
  "message": "shortlink created successfully"
}
```

### Delete Shortlink

* **Endpoint:** `DELETE /api/links/:key`
* **Response Status:** `200 OK`
* **Response Body:**
```json
{
  "message": "shortlink deleted successfully"
}
```
