const { ApolloServer, PubSub } = require("apollo-server");
var { Source } = require("graphql");
const velocity = require("velocityjs");
const yaml = require("js-yaml");
const axios = require("axios");
const chalk = require("chalk");
const fs = require("fs");
const moment = require("moment");

var AWS = require('aws-sdk');
AWS.config.region = 'us-eastt-1'

AWS.config.update({
   accessKeyId: 'AKID', secretAccessKey: 'SECRET',
   region: 'us-east-1',
   endpoint: 'http://serverless-cms:3001',
   sslEnabled: false
});

/*
AWS.config.lambda = { endpoint: 'http://serverless-cms:3001',
                      sslEnabled: false
}
*/

const {GraphQLScalarType} = require('graphql');

const { visit } = require("graphql/language/visitor");
const { parse } = require("graphql/language");

require("dotenv").config();

const options = {
    quiet: false
};
process.argv.forEach(option => {
    if (option === "-q" || option === "--quiet") options.quiet = true;
});

const pubsub = new PubSub();

const CF_SCHEMA = yaml.Schema.create([
    new yaml.Type("!Ref", {
        kind: "scalar",
        construct: function(data) {
            return {
                Ref: data
            };
        }
    }),
    new yaml.Type("!Equals", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::Equals": data
            };
        }
    }),
    new yaml.Type("!Not", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::Not": data
            };
        }
    }),
    new yaml.Type("!Sub", {
        kind: "scalar",
        construct: function(data) {
            return {
                "Fn::Sub": data
            };
        }
    }),
    new yaml.Type("!If", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::If": data
            };
        }
    }),
    new yaml.Type("!Join", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::Join": data
            };
        }
    }),
    new yaml.Type("!Select", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::Select": data
            };
        }
    }),
    new yaml.Type("!FindInMap", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::FindInMap": data
            };
        }
    }),
    new yaml.Type("!GetAtt", {
        kind: "sequence",
        construct: function(data) {
            return {
                "Fn::GetAtt": data
            };
        }
    }),
    new yaml.Type("!GetAZs", {
        kind: "scalar",
        construct: function(data) {
            return {
                "Fn::GetAZs": data
            };
        }
    }),
    new yaml.Type("!Base64", {
        kind: "mapping",
        construct: function(data) {
            return {
                "Fn::Base64": data
            };
        }
    })
]);

//Read GraphQL schema
var typeDefs = fs.readFileSync("schema.graphql", "utf8");
typeDefs = 'scalar AWSDateTime\n' + typeDefs;
const schemaAST = parse(typeDefs);

//Fill in subscriptions from the Schema
let subscriptions = {
    Query: {},
    Mutation: {}
};
let eventLists = {};
const visitor = {
    enter(node, key, parent, path, ancestors) {
        if (node.kind === "ObjectTypeDefinition" && node.name.value === "Subscription") {
            node.fields.forEach(node => {
                if (node.kind === "FieldDefinition") {
                    const subscriptionName = node.name.value;
                    node.directives.forEach(node => {
                        if (node.name.value === "aws_subscribe") {
                            node.arguments.forEach(node => {
                                let operationType;
                                switch (node.name.value) {
                                    case "queries":
                                        operationType = "Query";
                                        break;
                                    case "mutations":
                                        operationType = "Mutation";
                                        break;
                                }
                                if (operationType) {
                                    node.value.values.forEach(node => {
                                        if (
                                            subscriptions[operationType][node.value] === undefined
                                        ) {
                                            subscriptions[operationType][node.value] = [
                                                subscriptionName
                                            ];
                                        } else {
                                            subscriptions[operationType][node.value].push(
                                                subscriptionName
                                            );
                                        }

                                        const eventName = operationType + "_" + node.value;
                                        if (eventLists[subscriptionName] === undefined) {
                                            eventLists[subscriptionName] = [eventName];
                                        } else {
                                            eventLists[subscriptionName].push(eventName);
                                        }
                                    });
                                }
                            });
                        }
                    });
                }
            });
        }
    }
};

visit(schemaAST, visitor);

//Read CloudFront template
const cfTemplate = yaml.load(fs.readFileSync("./template.yml", "utf8"), {
    schema: CF_SCHEMA
});

//Fill in datasources from CF template
let dataSources = {};
Object.keys(cfTemplate.Resources).forEach(name => {
    const R = cfTemplate.Resources[name];
    /*
	console.error(name);
    console.error(R);
    console.error(R.RequestMappingTemplate);
	*/
    if (R.Type === "AWS::AppSync::DataSource") {
        const {
            Properties: { LambdaConfig }
        } = R;

        console.error('LAMBDA CONFIG IS FOR ' + name);
        console.error(LambdaConfig);
        if (LambdaConfig !== undefined) {
            const {
                LambdaFunctionArn: {
                    "Fn::GetAtt": functionName
                }
            } = LambdaConfig;

	    console.error('func name is ' + functionName);
            if (functionName !== undefined) {
                dataSources[name] = functionName;
            }
        }
    }
});

//Fill in resolvers from CF template
let resolvers = {
    Subscription: {},
    Query: {},
    Mutation: {}
};
Object.keys(cfTemplate.Resources).forEach(name => {
    const R = cfTemplate.Resources[name];
    if (R.Type === "AWS::AppSync::Resolver") {
        console.error('GOTT RESOLVER for ' + name)
        const {
            Properties: {
                TypeName: typeName,
                FieldName: fieldName,
                DataSourceName: {
                    "Fn::GetAtt": { badnotvaliddataSourceName }
                },
            }
        } = R;
/*
 { Type: 'AWS::AppSync::Resolver',
serverless-cms-graphql_1  |   Properties: 
serverless-cms-graphql_1  |    { ApiId: { 'Fn::ImportValue': [Object] },
serverless-cms-graphql_1  |      TypeName: 'Mutation',
serverless-cms-graphql_1  |      FieldName: 'updateProject',
serverless-cms-graphql_1  |      DataSourceName: { 'Fn::GetAtt': 'AppSyncModelProjectDataSource.Name' },
serverless-cms-graphql_1  |      RequestMappingTemplate: '{ "version" : "2017-02-28", "operation": "Invoke", "payload": { "resolve": "mutation.updateProject", "context": $utils.toJson($context) } }',
serverless-cms-graphql_1  |      ResponseMappingTemplate: '$util.toJson($context.result)' } }
serverless-cms-graphql_1  | undefined
*/
        var dataSourceName = R.Properties.DataSourceName['Fn::GetAtt'].split('.')[0]

        console.error(dataSourceName)
        let {
            Properties: { RequestMappingTemplate: requestMappingTemplate }
        } = R;
        console.error(requestMappingTemplate)
        requestMappingTemplate = requestMappingTemplate.replace(
            /(\$utils\.toJson\()([^()]+)(\))/g,
            "$2"
        );

        if (typeName === "Subscription") {
            resolvers["Subscription"][fieldName] = {
                subscribe: () => {
                    console.log(
                        `Resolver`,
                        chalk.black.bgBlue(fieldName),
                        `executed at ${d.toLocaleTimeString()}`
                    );
                    console.log(`Returned asyncIterator(${eventLists[fieldName]})`);
                    return pubsub.asyncIterator(eventLists[fieldName]);
                }
            };
        } else {
            //local lambda endpoint for the resolver
            let lambdaEndpoint;
            console.error('tetst');
            console.error(dataSources)
            console.error(dataSourceName)
            if (dataSources[dataSourceName] !== undefined) {
                lambdaEndpoint = `http://${process.env.LOCAL_LAMBDA_HOST}:${
                    process.env.LOCAL_LAMBDA_PORT
                }/2015-03-31/functions/${dataSources[dataSourceName]}/invocations`;
            } else {
                console.log(
                    chalk.black.bgYellow("WARNING"),
                    `Lambda endpoint is not defined for the ${fieldName} resolver`
                );
            }

            //resolver function
            console.error(resolvers[typeName][fieldName])
            resolvers[typeName][fieldName] = async (root, args, context) => {
                const d = new Date();
                console.log(
                    `Resolver`,
                    chalk.black.bgBlue(fieldName),
                    `executed at ${d.toLocaleTimeString()}`
                );

                console.log("Rendering velocity template...");
                let template = velocity.render(requestMappingTemplate, {
                    context: {
                        arguments: JSON.stringify(args),
                        request: {
                            headers: JSON.stringify(context.request.headers)
                        },
                        identity: JSON.stringify(context.identity)
                    }
                });
                // TODO arguments = with arguments : enclosed by cquotes... and others
		template = template.replace('{arguments={},', '{"arguments":{},')
		template = template.replace('{arguments={"', '{"arguments":{"')
		template = template.replace('request={headers=', '"request":{"headers":')
		template = template.replace('arguments={}, request={headers=', '"arguments":{}, "request":{"headers":')
		template = template.replace('identity={"source', '"identity":{"source')
		console.error(template);
                template = JSON.parse(template);
                const payload = JSON.parse(JSON.stringify(template.payload));
		console.error(payload);

                if (!options.quiet) {
                    const headers = template.payload.headers;
                    if (headers !== undefined && headers !== null) {
                        Object.keys(headers).map(function(key, index) {
                            headers[key] =
                                key === "x-amz-security-token"
                                    ? headers[key].substring(0, 15) + "...[truncated]"
                                    : headers[key];
                        });
                    }
                    console.log("Resulting template:", JSON.stringify(template));
                }

                if (lambdaEndpoint === undefined) {
                    console.log("No endpoint defined, nothing to do..");
                    return;
                }

                console.log("Invoking lambda function with payload...");
                console.log(lambdaEndpoint);
                console.log(payload);
		var lambda = new AWS.Lambda({maxRetries: 1});
		var params = {
		  FunctionName: 'ModelLambda',
		  Payload: JSON.stringify(payload)
		};

		const lambdaResult = await lambda.invoke(params).promise();

		// new Lambda({ region: 'us-east-1', endpoint: 'http://docker.for.mac.localhost:3001' })
		/*
                const response = await axios.post(lambdaEndpoint, payload, {
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json"
                    }
                });
		*/

		console.error('received of ')
		console.error(lambdaResult)
		console.error(lambdaResult.Payload)
                const rr  = JSON.parse(lambdaResult.Payload);
                const data = rr
                console.error(data);
                if (data.errorMessage !== undefined) {
                    console.log("Lambda response:", chalk.black.bgRed("ERROR"));
                    console.log("Error: ", JSON.stringify(data));
                    console.log("");
                    throw new Error(data.errorMessage);
                } else {
                    console.log("Lambda response:", chalk.black.bgGreen("DATA"));
                    if (!options.quiet) {
                        console.log("Data: ", JSON.stringify(data));
                    }

                    if (subscriptions[typeName][fieldName] !== undefined) {
                        const eventName = typeName + "_" + fieldName;
                        console.log(
                            `Publishing event "${eventName}" to subscriptions: ${JSON.stringify(
                                subscriptions[typeName][fieldName]
                            )}...`
                        );
                        subscriptions[typeName][fieldName].forEach(subscriptionName =>
                            pubsub.publish(eventName, { [subscriptionName]: data })
                        );
                    }

                    console.log("");
                    return data;
                }
            };
        }
    }
});

/*
const Date = new GraphQLScalarType({
  name: 'Date',
  serialize(value) {
    return value;
  },
});
*/

const AWSDateTime = new GraphQLScalarType({
    name: 'AWSDateTime',
    description: 'AWS DateTime...Date custom scalar type',
    parseValue(value) {
      console.log('teest ' + value);
      return moment(value).toDate(); // value from the client
    },
    serialize(value) {
      console.log('111 PLEASE!!!')
      console.log(value)
      if (!moment(value).isValid()) {
        console.log('PLEASE!!! err')
        return;
      } else{
         console.error('validD?');
      }
      return moment(value);
    },
    parseLiteral(ast) {
      console.log('zzz 111 PLEASE!!!')
      console.log(ast);
      if (ast.value && moment(ast.value).isValid()) {
        console.log('OK TEST');
        return moment(ast.value).format();
      }
      console.log('NULLL NO');
      return null;
      /*
      if (ast.kind === Kind.INT) {
        return parseInt(ast.value, 10); // ast value is always in string format
      }
      return null;
      */
    },
})

// MAYBE ? 
const DateTimeScalar = {
  __parseValue(value) {
    return moment(value).toDate(); // value from the client
  },
  __serialize(value) {
    if (!moment(value).isValid()) {
      console.log('ZKLKZ PLEASE!!!')
      return;
    }
    return moment(value).format(); // value sent to the client
  },
  __parseLiteral(ast) {
    if (ast.value && moment(ast.value).isValid()) {
      return moment(ast.value).format();
    }
    return null;
  }
}

/*
const SubscriptionNoop = new GraphQLScalarType({
    name: 'Subscription',
    description: 'Subscription',
    parseValue(value) {
      return value; //new Date(value); // value from the client
    },
    serialize(value) {
      return value // value.getTime(); // value sent to the client
    },
    parseLiteral(ast) {
      return null;
    },
})
*/

resolvers['AWSDateTime'] = AWSDateTime // DateTimeScalar // = {
// resolvers['Subscription'] = SubscriptionNoop
delete resolvers['Subscription']
console.error(resolvers)
/*
DDtypeName][fieldName] = async (root, args, context) => {
              const d = new Date();
                console.log(
                    `Resolver`,
*/
//creating and starting Apollo-server
const server = new ApolloServer({
    typeDefs,
    resolvers,
    context: async ({ req, connection }) => {
        if (connection) {
            return {};
        } else {
            return {
                request: { headers: req.headers },
                identity: { sourceIp: "127.0.0.1" }
            };
        }
    },
    tracing: true
});

server.listen().then(({ url, subscriptionsUrl }) => {
    console.log(chalk.bold(`Local AppSync ready at ${url}`));
    console.log(chalk.bold(`The subscriptions url is ${subscriptionsUrl}\n`));
});
