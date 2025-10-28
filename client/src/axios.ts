import axios from "axios";

const serverUrl = "https://go-u-sh.onrender.com/api";

type MethodType = "GET" | "POST" | "DELETE" | "PUT";

export const makeRequest = async ({
	endpoint,
	method,
	dataOrQuery,
	jwtToken,
}: {
	endpoint: string;
	method: MethodType;
	dataOrQuery: unknown;
	jwtToken?: string | null;
}) => {
	const url = serverUrl + endpoint;

	const options = {
		data: {},
		header: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${jwtToken}`,
		},
		method,
		url,
	};

	if (typeof dataOrQuery === "string") {
		options.url = `${url}?id=${dataOrQuery}`;

		const { data } = await axios.request(options);
		return data;
	}

	options.data = dataOrQuery as Record<string, string | string[]>;

	const { data } = await axios.request(options);
	console.log({ data });
	return data;
};
