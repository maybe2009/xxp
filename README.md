# xxp
XXP is a network communication protocol written in go.
- Protocol raw data is in binary format
- For fixed size type(e.g. int32 int64 etc.), the correspongding protocol data is just the binray presentation of this type
- For variable length type(e.g. string array slice etc.), the corresponging protocol data is in "LV"(length-value) format
