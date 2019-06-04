# aws-serverless-cms-todo-app

[Jour our slacck community @ slack.acentera.com](https://slack.acentera.com/)

---

Vue.JS sample CMS Project & Todo App

This repository is based on the [Vue Element Admin](https://github.com/PanJiaChen/vue-element-admin) Vue.JS admin.

If this app has been deployed fully using the [serverless repository link](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:356769441913:applications~acentera-prod-serverless-cms-todo-app), once the cloudformatio nstack has been created, you should access the web application from the *WebsiteUrl cloudformation outputs*.


Once you navigate to the Cloudfront URL you should have a login page that looks like the following:

  ![01 - User Login](https://github.com/ACenterA/acentera-aws-serverless-cms-todo/raw/master/docs/images/01_login_user.png)

With a valid cognito user (from the ACenterA Core), you should be able to log-in using an MFA confirmation code.

  ![02 - User MFA](https://github.com/ACenterA/acentera-aws-serverless-cms-todo/raw/master/docs/images/02_login_user_mfa.png)

Once logged-in you should be able to navigate to the "Administration" section, only if your user is part of the "Admin" cognito group

  ![04 - Create Projects](https://github.com/ACenterA/acentera-aws-serverless-cms-todo/raw/master/docs/images/04_admin_create_project.png)
  
As admin you can list projects 

  ![06 - Manage Projects](https://github.com/ACenterA/acentera-aws-serverless-cms-todo/raw/master/docs/images/05_manage_projects.png)


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

# Launching using AWS Serverless Repo

Once you have deployed the serverless repo through the Serverless Repository, you should have a Cloudfront distributin created.

You may check in the stack output for the WebsiteUrl to use. 

Note: The Cloudfront should pass an 'X-Site' header that represent the SiteKey.

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


# License
The software is licensed under the MIT license.

# Attribution
By using this UI, we would like that any application which incorporates it shall prominently display the message “Made with ACenterA” in a legible manner in the footer of the admin console. This message must open a link to acentera.com when clicked or touched.
