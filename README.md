# aws-serverless-cms-todo-app

Vue.JS Project / Todo App

This repository is based on the [Vue Element Admin](https://github.com/PanJiaChen/vue-element-admin) Vue.JS admin.

# ENV File Informations

    $ cat .env

| Environment   | Value          | Info  |
| ------------- |---------------| -----:|
| PLUGIN_NAME   | serverless-cms | Plugin Name used in template registration process |
| PLUGIN_KEY    | a2f9b231-88b5-40c9-9f75-cc14f854a7bb     | Unique Site Key passed as HTTP Header |
| PLUGIN_DATA_KEY    | a2f9b231-88b5-40c9-9f75-cc14f854a7bb     | Data Key may be usefull to refer to an existing environment data |
| DEV_ACCOUNTID | 0      | AWS Account Id simulator (in dev) |
| ADMIN_USERNAME | admin@acentera.com | Admin Username provisionned |
| SITE_TITLE | My Serverless Admin - Local | Vue.JS Site Title |
| ENV_CONFIG | local | Specify local to use local graphql, anything else would use values from the DB |

# How to launch a local Development environment

Using Makefile

    $ make dev

Using Docker Compose only

    $ docker-compose -f docker-compose.yml -f docker-compose-core.yml up -d --build --force-recreate


> Note: The local environment doesnt use Cognito, while the production will using Cognito for user / groups.

------

> Note: The local environment usees an appsync nodejs proxy to call sam model. That doesnt yet supports GraphQL Subscriptions.
> There are way to actually uses a real GraphQL / Cognito Endpoint

# Endpoints

| Endpoint Name  | Value          | Info  |
| ------------- |---------------| -----:|
| Vue.JS App | http://127.0.0.1/ | Vue.JS App |
| DynamoDB   | http://127.0.0.1:8001/ | DynamoDB Tables |
| GraphQL    | http://127.0.0.1:4000/ | Apollo GraphQL Server |
| API Gateway Proxy | http://127.0.0.1:2000/ | API Gateway Proxy (acentera core or plugin routing) |

# Deploy on AWS

Using an AWS CodeBuild Pipeline (or simply using the sam package / publish ).

You have a pre-requisit of the [ACenterA Core](https://github.com/ACenterA/acentera-aws-core) serverless repo.

# How to Rename the Plugin

Create a Fork and rename

find src/     -type f -exec sed -ri 's/serverless-cms/my-new-plugin/g' {} \;
find config/  -type f -exec sed -ri 's/serverless-cms/my-new-plugin/g' {} \;
find build/   -type f -exec sed -ri 's/serverless-cms/my-new-plugin/g' {} \;


# Special Thanks

  * [Georgy Nemtsov](https://github.com/gnemtsov)
     For the Node.js script, which works like local AppSync: appsync-local.js

