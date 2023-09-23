Using the official protobuf example "addressbook" as the reference

https://github.com/protocolbuffers/protobuf/tree/main/examples/go

1. Use uber golang style to rewrite the code
    - https://github.com/uber-go/guide/blob/master/style.md
2. Add more commands based on CRUD
    - Create
    - Read
        - Get by user id
        - List with sorting
    - Update
        - Update by id
    - Delete
        - Delete by id

Operations:
$protoc -I=proto --go_out=go/addressbookpb --go_opt=paths=source_relative adressbook.proto