go test -coverprofile=cover.test.out ./... &&\
go test -covermode=count -coverprofile=cover_count.test.out ./... &&\
go tool cover -func=cover.test.out &&\
go tool cover -html=cover.test.out &&\
go tool cover -html=cover_count.test.out
