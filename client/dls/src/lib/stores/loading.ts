import { writable } from 'svelte/store';

export let loading = writable<boolean>(false);
