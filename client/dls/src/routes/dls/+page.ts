import { fetchProject } from '@src/ts/api/helpers.js';
import { combineHeaders, getAuthHeader, getProjectIDHeader } from '@src/ts/helpers.js';

export async function load(ctx) {
	const res = await fetchProject((url) =>
		ctx.fetch(url, { method: 'GET', credentials: 'include', headers: combineHeaders(getAuthHeader(), getProjectIDHeader()) })
	);
	return {
		res
	};
}
