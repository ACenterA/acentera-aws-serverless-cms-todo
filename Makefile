#
# == usage ==
#
#   - Launch local development
#     make dev
#

help:
	@printf '=======================================\n'
	@printf '[>] To launch local development \n'
	@printf '[CMD] make dev\n'
	@printf '=======================================\n'

dev: 
	@printf '[>] Provisioning the docker containers\n'                                        
	docker-compose -f docker-compose.yml -f docker-compose-core.yml up -d --build --force-recreate

restart_graphql:
	@printf '[>] Resart graphql\n'                                        
	@(docker-compose up -d serverless-cms; docker-compose restart serverless-cms; sleep 2; docker-compose restart serverless-cms-graphql)
	@docker-compose logs serverless-cms-graphql | tail -n 30
