.PHONY: run
run:
	SECRET_KEY="gbS6TmleOxMtEsOqNCnO9anxNlxJlfDB8e+7KxAjSDg=" go run main.go

.PHONY: start_db
start_db:
	docker-compose up -d

.PHONY: stop_db
stop_db:
	docker-compose down

.PHONY: fill_test
fill_test:
	PGPASSWORD=password psql -U username -d bookstore -h localhost < testdata/initdata.sql

.PHONY: erase_test
erase_test:
	PGPASSWORD=password psql -U username -d bookstore -h localhost < testdata/erase_initdata.sql

.PHONY: print_test
print_test:
	PGPASSWORD=password psql -U username -d bookstore -h localhost < testdata/print_all.sql

.PHONY: generate_uml
generate_uml:
	mkdir -p doc/img/ && java -jar doc/plantuml.jar -o ../img/ doc/uml/*.puml
