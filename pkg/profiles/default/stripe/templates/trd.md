# API Specification: {{.Title}}

> **Contract-First**: This specification IS the source of truth. Implementation follows this spec.

## Overview

Brief description of what this API does and who it's for.

**Base URL**: `https://api.example.com/v1`

**Authentication**: Bearer token / API key

---

## Resources

### Resource: `{{resource_name}}`

A `{{resource_name}}` represents...

#### The {{resource_name}} object

```json
{
  "id": "res_1234567890",
  "object": "{{resource_name}}",
  "created": 1234567890,
  "livemode": true,

  // Core attributes
  "name": "string",
  "description": "string | null",
  "metadata": {},

  // Status
  "status": "active | pending | archived"
}
```

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | Unique identifier for the object. |
| `object` | string | String representing the object's type. Always `{{resource_name}}`. |
| `created` | timestamp | Time at which the object was created. Measured in seconds since Unix epoch. |
| `livemode` | boolean | Has the value `true` if the object exists in live mode or `false` if in test mode. |
| `name` | string | The name of the resource. |
| `description` | string | Optional description. |
| `metadata` | hash | Set of key-value pairs for storing additional information. |
| `status` | enum | Current status: `active`, `pending`, or `archived`. |

---

## Endpoints

### Create a {{resource_name}}

Creates a new {{resource_name}} object.

```
POST /v1/{{resource_name}}s
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | The name of the resource. Maximum 255 characters. |
| `description` | string | No | An optional description. Maximum 1000 characters. |
| `metadata` | hash | No | Set of key-value pairs. Keys can be up to 40 characters; values up to 500 characters. |

#### Request

```bash
curl https://api.example.com/v1/{{resource_name}}s \
  -u sk_test_xxx: \
  -d name="Example Resource" \
  -d "metadata[order_id]=12345"
```

#### Response

Returns the created `{{resource_name}}` object if successful. Returns an error if parameters are invalid.

```json
{
  "id": "res_1234567890",
  "object": "{{resource_name}}",
  "created": 1234567890,
  "livemode": false,
  "name": "Example Resource",
  "description": null,
  "metadata": {
    "order_id": "12345"
  },
  "status": "active"
}
```

---

### Retrieve a {{resource_name}}

Retrieves the details of an existing {{resource_name}}.

```
GET /v1/{{resource_name}}s/:id
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | The identifier of the {{resource_name}} to retrieve. |

#### Request

```bash
curl https://api.example.com/v1/{{resource_name}}s/res_1234567890 \
  -u sk_test_xxx:
```

#### Response

Returns a `{{resource_name}}` object if a valid identifier was provided. Returns an error otherwise.

---

### Update a {{resource_name}}

Updates the specified {{resource_name}} by setting the values of the parameters passed.

```
POST /v1/{{resource_name}}s/:id
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | No | The name of the resource. |
| `description` | string | No | The description of the resource. |
| `metadata` | hash | No | Set of key-value pairs. |

#### Request

```bash
curl https://api.example.com/v1/{{resource_name}}s/res_1234567890 \
  -u sk_test_xxx: \
  -d name="Updated Name"
```

#### Response

Returns the updated `{{resource_name}}` object if successful.

---

### Delete a {{resource_name}}

Permanently deletes a {{resource_name}}. This cannot be undone.

```
DELETE /v1/{{resource_name}}s/:id
```

#### Request

```bash
curl https://api.example.com/v1/{{resource_name}}s/res_1234567890 \
  -u sk_test_xxx: \
  -X DELETE
```

#### Response

Returns an object with `deleted: true` if successful.

```json
{
  "id": "res_1234567890",
  "object": "{{resource_name}}",
  "deleted": true
}
```

---

### List all {{resource_name}}s

Returns a list of {{resource_name}}s. The {{resource_name}}s are returned sorted by creation date, with the most recent first.

```
GET /v1/{{resource_name}}s
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `limit` | integer | No | Number of objects to return. Default 10, maximum 100. |
| `starting_after` | string | No | Cursor for pagination. Object ID to start after. |
| `ending_before` | string | No | Cursor for pagination. Object ID to end before. |
| `created` | hash | No | Filter by creation date. Supports `gt`, `gte`, `lt`, `lte`. |
| `status` | string | No | Filter by status. |

#### Request

```bash
curl https://api.example.com/v1/{{resource_name}}s?limit=3 \
  -u sk_test_xxx:
```

#### Response

Returns a paginated list object containing the matching {{resource_name}}s.

```json
{
  "object": "list",
  "url": "/v1/{{resource_name}}s",
  "has_more": true,
  "data": [
    {
      "id": "res_1234567890",
      "object": "{{resource_name}}",
      ...
    }
  ]
}
```

---

## Errors

This API uses conventional HTTP response codes. Codes in the 2xx range indicate success. Codes in the 4xx range indicate an error with the information provided. Codes in the 5xx range indicate a server error.

### Error Response Format

```json
{
  "error": {
    "type": "invalid_request_error",
    "code": "resource_missing",
    "message": "No such {{resource_name}}: res_invalid",
    "param": "id",
    "doc_url": "https://docs.example.com/errors#resource_missing"
  }
}
```

### Error Types

| Type | Description |
|------|-------------|
| `invalid_request_error` | Invalid parameters were supplied. |
| `authentication_error` | Invalid or missing API key. |
| `rate_limit_error` | Too many requests hit the API too quickly. |
| `api_error` | Something went wrong on our end. |

### Error Codes

| Code | HTTP Status | Description | Remediation |
|------|-------------|-------------|-------------|
| `resource_missing` | 404 | The requested resource doesn't exist. | Verify the ID is correct and the resource hasn't been deleted. |
| `parameter_invalid` | 400 | A parameter is invalid. | Check the `param` field and correct the value. |
| `parameter_missing` | 400 | A required parameter is missing. | Add the required parameter specified in `param`. |
| `idempotency_key_in_use` | 409 | The idempotency key is already in use with different parameters. | Use a new idempotency key or retry with the same parameters. |
| `rate_limit_exceeded` | 429 | Too many requests. | Back off and retry with exponential backoff. |

---

## Pagination

All list endpoints support cursor-based pagination using `starting_after` and `ending_before` parameters.

```bash
# First page
curl https://api.example.com/v1/{{resource_name}}s?limit=10 \
  -u sk_test_xxx:

# Next page (use last object's ID)
curl https://api.example.com/v1/{{resource_name}}s?limit=10&starting_after=res_xyz \
  -u sk_test_xxx:
```

---

## Idempotency

The API supports idempotency for safely retrying requests. Pass a unique `Idempotency-Key` header.

```bash
curl https://api.example.com/v1/{{resource_name}}s \
  -u sk_test_xxx: \
  -H "Idempotency-Key: unique-key-12345" \
  -d name="Example"
```

Keys expire after 24 hours.

---

## Versioning

The API version is specified in the URL path (`/v1/`). Breaking changes will be released under a new version.

Current version: `2024-01-01`

You can pin to a specific version using the `API-Version` header.

---

## Webhooks

Events are sent to configured webhook endpoints.

### Event Types

| Event | Description |
|-------|-------------|
| `{{resource_name}}.created` | A new {{resource_name}} was created. |
| `{{resource_name}}.updated` | A {{resource_name}} was updated. |
| `{{resource_name}}.deleted` | A {{resource_name}} was deleted. |

### Webhook Payload

```json
{
  "id": "evt_1234567890",
  "object": "event",
  "type": "{{resource_name}}.created",
  "created": 1234567890,
  "data": {
    "object": {
      "id": "res_1234567890",
      ...
    }
  }
}
```

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| All endpoints | 100 requests per second |
| List endpoints | 25 requests per second |

Rate limit headers are included in all responses:

- `X-RateLimit-Limit`: Maximum requests per second
- `X-RateLimit-Remaining`: Remaining requests in window
- `X-RateLimit-Reset`: Unix timestamp when limit resets

---

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 2024-01-01 | 2024-01-01 | Initial release |
