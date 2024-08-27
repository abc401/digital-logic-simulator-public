import { userStore } from '@src/lib/stores/user.js';
import { apiURL } from '@src/ts/api/helpers.js';
import { getAuthHeader } from '@src/ts/helpers.js';

export interface ProjectMetaData {
	id: number;
	name: string;
}

export enum _Problem {
	NotLoggedIn,
	NotFound,
	BadRequest,
	None
}

export async function load(ctx) {
	userStore.init();
	const user = userStore.get();

	if (user == null) {
		return {
			projects: [],
			problem: _Problem.NotLoggedIn
		};
	}
	const res = await ctx.fetch(apiURL('project-list'), {
		credentials: 'include',
		headers: getAuthHeader()
	});
	const json = await res.json();
	return {
		projects: json.list as ProjectMetaData[],
		problem: _Problem.None
	};
}
