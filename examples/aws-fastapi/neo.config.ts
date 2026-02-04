/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
	app(input) {
		return {
			name: "aws-fastapi",
			removal: input?.stage === "production" ? "retain" : "remove",
			home: "aws",
			providers: {
				aws: true,
			},
		};
	},
	async run() {
		const linkableValue = new neo.Linkable("MyLinkableValue", {
			properties: {
				foo: "Hello World",
			},
		});

		const fastapi = new neo.aws.Function("FastAPI", {
			handler: "functions/src/functions/api.handler",
			runtime: "python3.11",
			url: true,
			link: [linkableValue],
		});

		return {
			fastapi: fastapi.url,
		};
	},
});
