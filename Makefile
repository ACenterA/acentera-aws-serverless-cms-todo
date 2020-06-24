#
# == usage ==
#
#   - Launch local development
#     make dev
#

help:
	@printf '=======================================\n'
	@printf '[>] To launch local development \n'
	@printf '\n'
	@printf '[CMD] make dev\n'
	@printf '\n'
	@printf '=======================================\n'
	@printf '[>] If you are running acentera-core as a outside app\n'
	@printf '    You should only start the dev-plugin containers\n'
	@printf '\n'
	@printf '[CMD] make dev-plugin\n'
	@printf '\n'
	@printf '=======================================\n'

proxy: 
	@printf '[>] Provisioning the docker containers\n'                                        
	docker-compose -f docker-compose-core-db.yml up -d proxy-core
	#--force-recreate
dev: 
	@printf '[>] Provisioning the docker containers\n'                                        
	docker-compose -f docker-compose.yml -f docker-compose-core.yml -f docker-compose-core-db.yml up -d --build
	#--force-recreate

dev-plugin:
	@printf '[>] Provisioning the docker containers\n'                                        
	docker-compose -f docker-compose.yml up -d --build
	#--force-recreate

dev-plugin-with-db:
	@printf '[>] Provisioning the docker containers\n'                                        
	docker-compose -f docker-compose.yml -f docker-compose-core-db.yml up -d --build
	#--force-recreate

#	@(docker-compose up -d serverless-cms; docker-compose restart serverless-cms; sleep 5; docker-compose restart serverless-cms-graphql)
restart_graphql:
	@printf '[>] Resart graphql\n'                                        
	@(docker-compose exec -T serverless-cms make regen_graphql && docker-compose restart serverless-cms-graphql)
	@docker-compose logs serverless-cms-graphql | tail -n 30
