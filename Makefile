.PHONY: postgres adminer migrate

postgres:
	sudo docker run --rm -ti --network host -e POSTGRES_PASSWORD=secret postgres

adminer:
	sudo docker run --rm -ti --network host adminer

migrate:
	migrate -source file://migrations \
			-database postgres://postgres:secret@localhost/postgres?sslmode=disable up

migrate-down:
	migrate -source file://migrations \
			-database postgres://postgres:secret12@localhost/postgres?sslmode=disable down