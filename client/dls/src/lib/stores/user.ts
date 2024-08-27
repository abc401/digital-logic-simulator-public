import { writable } from 'svelte/store';
import Cookies from 'js-cookie';
import { jwtDecode, type JwtPayload } from 'jwt-decode';
import { getAuthToken } from '@src/ts/helpers';

const AUTH_COOKIE_NAME = 'bearer';

export interface User {
	email: string;
	uname: string;
	uid: number;
}

const { set, subscribe } = writable<User | undefined>(undefined);

export const userStore = {
	set,
	subscribe,

	init() {
		const authToken = getAuthToken();
		if (authToken == null) {
			console.log('Auth token not present');
			return;
		}

		// console.log(authCookie);
		let payload: JwtPayload;
		try {
			payload = jwtDecode(authToken);
		} catch {
			console.log('Invalid auth token');
			return;
		}
		const exp = payload.exp;
		if (exp == null) {
			console.log('Expiration time not provided on jwt');
			return;
		}
		if (exp * 1000 < Date.now()) {
			console.log('Jwt expired');
			return;
		}
		const uname: string | undefined = payload['username'];
		if (uname == null) {
			return;
		}
		const email: string | undefined = payload['email'];
		if (email == null) {
			return;
		}
		const uid: number | undefined = payload['uid'];
		if (uid == null) {
			return;
		}
		userStore.set({ uname, email, uid });
	},
	get() {
		let user: User | undefined;
		const unsub = subscribe((value) => {
			user = value;
		});
		unsub();
		return user;
	}
};
