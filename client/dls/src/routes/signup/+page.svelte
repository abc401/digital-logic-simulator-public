<script lang="ts">
	import '@src/app.css';
	import { apiURL } from '@src/ts/api/helpers';
	import { goto } from '$app/navigation';
	import ErrorMsg from '@src/lib/components/ErrorMsg.svelte';
	import { loading } from '@src/lib/stores/loading';
	import { makePath } from '@src/ts/helpers';
	let uname: string;
	let pword: string;
	let cPword: string;
	let email: string;

	let emailError: string | undefined = undefined;
	let unameError: string | undefined = undefined;
	let pwordError: string | undefined = undefined;
	let cpwordError: string | undefined = undefined;
	let genericError: string | undefined = undefined;

	function setError(json: any) {
		emailError = undefined;
		unameError = undefined;
		pwordError = undefined;
		cpwordError = undefined;
		genericError = undefined;

		if (json['field'] === 'email') {
			emailError = json['error'];
		} else if (json['field'] === 'uname') {
			unameError = json['error'];
		} else if (json['field'] === 'pword') {
			pwordError = json['error'];
		} else if (json['field'] === 'cpword') {
			cpwordError = json['error'];
		} else {
			genericError = 'Internal Server Error';
		}
	}
</script>

<div class="grid min-h-screen items-center">
	<div class="mx-auto flex flex-col gap-8">
		<form
			on:submit={async (ev) => {
				ev.preventDefault();
				const url = apiURL('/signup');
				loading.set(true);
				const res = await fetch(url, {
					method: 'POST',
					body: JSON.stringify({
						uname: uname,
						pword: pword,
						cpword: cPword,
						email: email
					})
				});

				const json = await res.json();
				console.log('response: ', res);
				console.log('body: ', json);

				setError(json);

				if (res.ok) {
					await goto(makePath('/signin'));
					return;
				}
				loading.set(false);
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
						<ErrorMsg msg={emailError} />
					</div>
				</label>
				<label for="uname" class="contents">
					<span class="text-end text-xl capitalize">username</span>
					<div class="flex flex-col gap-1">
						<input
							required
							class="min-w-20 rounded-md bg-neutral-700 px-5 py-1 text-base text-neutral-50 focus:bg-neutral-50 focus:text-neutral-950"
							bind:value={uname}
							type="uname"
							name="uname"
							id="uname"
						/>
						<ErrorMsg msg={unameError} />
					</div>
				</label>
				<label for="pword" class="contents">
					<span class="text-end text-xl capitalize">Password</span>

					<div class="flex flex-col gap-1">
						<input
							required
							class="min-w-20 rounded-md bg-neutral-700 px-5 py-1 text-base text-neutral-50 focus:bg-neutral-50 focus:text-neutral-950"
							bind:value={pword}
							type="password"
							name="pword"
							id="pword"
						/>
						<ErrorMsg msg={pwordError} />
					</div>
				</label>
				<label for="cpword" class="contents">
					<span class="text-end text-xl capitalize">Confirm Password</span>
					<div class="flex flex-col gap-1">
						<input
							required
							class="min-w-20 rounded-md bg-neutral-700 px-5 py-1 text-base text-neutral-50 focus:bg-neutral-50 focus:text-neutral-950"
							bind:value={cPword}
							type="password"
							name="cpword"
							id="cpword"
						/>
						<ErrorMsg msg={cpwordError} />
					</div>
				</label>
			</div>
			<button
				class="rounded-full border-2 border-accent px-5 py-3 text-xl capitalize transition-colors hover:bg-accent active:bg-accent"
				>signup</button
			>
		</form>
		<div class="text-xs">Already have an account? <a href={makePath('./signin')}>Sign in</a></div>
		<ErrorMsg msg={genericError} />
	</div>
</div>
