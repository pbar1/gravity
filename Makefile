consul-start:
	docker-compose -f test/docker-compose.yml up -d

consul-stop:
	docker-compose -f test/docker-compose.yml down

# needs consul on local PATH...
consul-config: consul-start
	bash ./test/consul-config.sh
