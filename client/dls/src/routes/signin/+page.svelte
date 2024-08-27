<script lang="ts">
	import '@src/app.css';
	import { apiURL } from '@src/ts/api/helpers';
	import { goto } from '$app/navigation';
	import { userStore, type User } from '@src/lib/stores/user';
	import ErrorMsg from '@src/lib/components/ErrorMsg.svelte';
	import { loading } from '@src/lib/stores/loading';
	import { onMount } from 'svelte';
	import { makePath } from '@src/ts/helpers';

	let email: string;
	let pword: string;

	let emailError: string | undefined = undefined;
	let pwordError: string | undefined = undefined;
	let genericError: string | undefined = undefined;

	function setError(json: any) {
		emailError = undefined;
		pwordError = undefined;
		genericError = undefined;

		if (json['field'] === 'email') {
			emailError = json['error'];
		} else if (json['field'] === 'pword') {
			pwordError = json['error'];
		} else {
			genericError = 'Internal Server Error';
		}
	}

	function userFromJSON(json: any) {
		const user: User = {
			uname: json['uname'],
			email: json['email'],
			uid: json['uid']
		};
		return user;
	}

	$: console.log($userStore);
	onMount(() => {
		loading.set(false);
	});
</script>

<div class="grid min-h-screen items-center">
	<div class="mx-auto flex flex-col gap-8">
		<form
			on:submit={async (ev) => {
				ev.preventDefault();
				const url = apiURL('/signin');
				loading.set(true);
				const res = await fetch(url, {
					method: 'POST',
					body: JSON.stringify({
						pword: pword,
						email: email
					}),
					credentials: 'include'
				});

				const json = await res.json();

				localStorage.setItem('bearer', json['token']);
				console.log('response: ', res);
				console.log('body: ', json);

				if (res.ok) {
					userStore.set(userFromJSON(json));
					goto(makePath('/profile'));
				} else {
					setError(json);
					loading.set(false);
				}
			}}
			method="post"
			class=" flex w-fit flex-col gap-6"
		>
			<div class="grid w-fit grid-cols-[auto,minmax(2rem,1fr)] gap-x-4 gap-y-4">
				<label for="email" class="contents">
					<span class="text-end text-xl capitalize">email</span>
					<div class="flex flex-col gap-1">
						<input
							required
							class="min-w-20 rounded-md bg-neutral-700 px-5 py-1 text-base text-neutral-50 focus:bg-neutral-50 focus:text-neutral-950"
							bind:value={email}
							type="email"
							name="email"
							id="email"
						/>
						<!-- {#if emailError != undefined} -->
						<ErrorMsg msg={emailError} />
						<!-- <span class="text-red-600">{emailError}</span> -->
						<!-- {/if} -->
					</div>
				</label>
				<label for="pword" class="contents">
					<span class="text-end text-xl capitalize">Password</span>

					<div class="flex flex-col gap-1 text-base">
						<input
							required
							class="min-w-20 rounded-md bg-neutral-700 px-5 py-1 text-neutral-50 focus:bg-neutral-50 focus:text-neutral-950"
							bind:value={pword}
							type="password"
							name="pword"
							id="pword"
						/>
						<!-- {#if pwordError != undefined} -->
						<ErrorMsg msg={pwordError} />
						<!-- {/if} -->
					</div>
				</label>
			</div>
			<button
				class="rounded-full border-2 border-accent px-5 py-3 text-xl capitalize transition-colors hover:bg-accent active:bg-accent"
				>signin</button
			>
		</form>
		<div class="text-xs">
			Don't have an account? <a href={makePath('./signup')}>Create one</a>
		</div>
		<ErrorMsg msg={genericError} />
	</div>
</div>
