.PHONY: run
run:
	@echo "build and run"
	docker-compose up --build -d

.PHONY: stop
stop:
	@echo "stop"
	docker-compose down

.PHONY: install-client
install-client:
	@echo "install client"
	${MAKE} -C grpc-upload install
