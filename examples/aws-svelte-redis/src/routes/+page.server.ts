import { Resource } from "neo";
import { Cluster } from "ioredis";

const redis = new Cluster(
	[{ host: Resource.MyRedis.host, port: Resource.MyRedis.port }],
	{
		dnsLookup: (address, callback) => callback(null, address),
		redisOptions: {
			tls: {},
			username: Resource.MyRedis.username,
			password: Resource.MyRedis.password,
		},
	}
);

/** @type {import('./$types').PageServerLoad} */
export async function load() {
	const counter = await redis.incr("counter");

	return { counter };
}
