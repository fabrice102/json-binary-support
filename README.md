# Test of support of encoding of arbitrary byte strings in JSON

The goal of this repo is to test and understand how JSON libraries handle arbitrary non-UTF-8 strings.
While the JSON [ECMA standard](https://www.ecma-international.org/publications-and-standards/standards/ecma-404/) allows to encode arbitrary characters by escaping them as `\u0000`, `\u0001`, ..., `\u00ff`, many libraries do not support this.

The [RFC 7159](https://datatracker.ietf.org/doc/html/rfc7159#section-8.2) specifies that:
> However, the ABNF in this specification allows member names and string values to contain bit sequences that cannot encode Unicode characters; for example, "\uDEAD" (a single unpaired UTF-16  surrogate).  [..] The behavior of software that receives JSON texts containing such values is unpredictable; for example, implementations might return different values for the length of a string value or even suffer fatal runtime exceptions.

## In Go

Both the [standard JSON package](https://pkg.go.dev/encoding/json) and the Algorand version of the [codec package](https://pkg.go.dev/github.com/ugorji/go/codec) do not encode nor decode non-UTF-8 strings as described above.

During encoding, invalid UTF-8 characters are replaced by `\ufffd`.

During decoding, characters such as `\u0080` are replaced by the pair of two UTF-8 characters corresponding to the unicode character U+0080.

This is consistent with the fact that strings in Golang are read-only slices of bytes (instead of a slices of unicode characters). See https://blog.golang.org/strings

It looks like in Go, the only way to encode arbitrary string in JSON using the above-mentioned encoding is to create a custom type with custom serialization/de-serialization methods.

This is possible both in the standard JSON package and in the codec package by implementing [the Marshaler interface](https://pkg.go.dev/encoding/json#Marshaler).

## In Python

Using byte arrays raise an exception when serializing.

However, if the byte array is encoded using the codec `raw-unicode-escape`, then the output JSON is exactly what is expected (where non-UTF-8 characters are replaced by `\u0000`, ..., `\u00ff`).