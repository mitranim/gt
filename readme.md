## Overview

Short for "**G**o **T**ypes". Important data types missing from the Go standard library.

* `Date`: civil date without time.
* `NullDate`: civil date where zero value is empty/null.
* `NullTime`: time where zero value is empty/null.
* `Interval`: ISO 8601 duration, corresponds to Postgres `interval`.
* `NullInterval`: interval where zero value is empty/null.
* `Uuid`: simple implementation of UUID version 4.
* `NullUuid`: UUID where zero value is empty/null.
* `NullInt`: int where zero value is empty/null.
* `NullUint`: uint where zero value is empty/null.
* `NullFloat`: float where zero value is empty/null.
* `NullUrl`: actually usable variant of `url.URL`, used by value rather than pointer, and where zero value is empty/null.
* `Ter`: nullable boolean (ternary), more usable and efficient than either `*bool` or `sql.NullBool`.

API docs: https://pkg.go.dev/github.com/mitranim/gt

Important features:

All types implement all relevant encoding/decoding interfaces for text, JSON, and SQL. Types can be seamlessly used for database fields, JSON fields, and so on, without the need for manual conversions.

All "nullable" `gt` types are type aliases of "normal" types, where the zero value of the normal type is "null". This lets you eliminate some invalid states at the type system level.

For example, for nullable DB enums, `gt.NullString` is a better choice than `*string` or `sql.NullString`, because it can't represent the often-invalid state of non-null `""`. Eliminating invalid states eliminates bugs. Similarly, `gt.NullTime` is a better choice than `time.Time`, `*time.Time` or `sql.NullTime`. It avoids the hassle of dealing with pointers or manual JSON conversions, while preventing your from accidentally inserting `0001-01-01`.

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
