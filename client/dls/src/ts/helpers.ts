import { base } from '$app/paths';
import path from 'path-browserify';
import { AUTH_TOKEN_LOCAL_STORAGE_NAME, PROJECT_ID_LOCAL_STORAGE_KEY } from './config';

export function exportToFile(blob: Blob, filename: string) {
	const a = document.createElement('a');
	const url = URL.createObjectURL(blob);
	a.href = url;
	a.download = filename;

	document.body.appendChild(a);
	a.click();

	setTimeout(function () {
		document.body.removeChild(a);
		window.URL.revokeObjectURL(url);
	}, 0);
}

export function makePath(targetUrl: string) {
	return path.join(base, targetUrl);

	// return `${base}/${path}`;
}

export function setProjectID(id: number) {
	localStorage.setItem(PROJECT_ID_LOCAL_STORAGE_KEY, id.toString());
}

export function getProjectID() {
	try {
		const id = localStorage.getItem(PROJECT_ID_LOCAL_STORAGE_KEY);
		if (id == null) {
			return -1;
		}
		return +id;
	} catch {
		return -1;
	}
}
export function setAuthToken(token: string) {
	localStorage.setItem(AUTH_TOKEN_LOCAL_STORAGE_NAME, token);
}
export function getAuthToken() {
	try {
		const token = localStorage.getItem(AUTH_TOKEN_LOCAL_STORAGE_NAME);
		if (token == null) {
			return 'undefined';
		}
		return token;
	} catch {
		return '';
	}
}

export function getAuthHeader() {
	const token = getAuthToken();
	return {
		Authorization: `Bearer ${token}`
	};
}
export function getProjectIDHeader() {
	const id = getProjectID();
	return {
		ProjectID: id.toString()
	};
}

export function combineHeaders(...headerContainers: { [key: string]: string }[]) {
	const res: { [key: string]: string } = {};

	for (const headers of headerContainers) {
		for (const [headerName, value] of Object.entries(headers)) {
			res[headerName] = value;
		}
	}
	return res;
}
