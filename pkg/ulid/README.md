# ulid

`ulid` is a thin wrapper around [`oklog/ulid`](https://github.com/oklog/ulid) for generating prefixed, time-sortable unique identifiers.

A ULID (Universally Unique Lexicographically Sortable Identifier) is a 128-bit value that is URL-safe, case-insensitive, and lexicographically sortable by creation time — making it a drop-in upgrade from UUIDs wherever ordering or readability matters.

## Quick start

```go
import "github.com/apexsudo/common/pkg/ulid"

id := ulid.New("user")
// → "user_01JCMK8NVDBBZJCF7R2TWKHQJH"
```

## API

### `New(prefix) string`

Generates a new ULID and returns it in the form `prefix_ULID`.

| Parameter | Type | Description |
|-----------|------|-------------|
| `prefix` | `string` | Short label identifying the entity type (e.g. `"user"`, `"order"`) |

**Example**

```go
userID  := ulid.New("user")   // "user_01JCMK8NVDBBZJCF7R2TWKHQJH"
orderID := ulid.New("order")  // "order_01JCMK8R3QKBM5XNPVTDHG2WYS"
```
