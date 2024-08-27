<script>
	import '@src/app.css';
	import './style.css';
	import { makePath } from '@src/ts/helpers';
</script>

<svelte:head>
	<title>Digital Logic Simulator</title>
</svelte:head>

<section class="header">
	<nav class="landing-page-nav">
		<div class="nav-links">
			<ul>
				<li><a data-sveltekit-reload href={makePath('./signin')}>Login</a></li>
				<li><a data-sveltekit-reload href={makePath('./signup')}>Sign up</a></li>
				<li><a data-sveltekit-reload href={makePath('./tutorials/index')}>Tutorials</a></li>
			</ul>
		</div>
	</nav>
	<div class="text-box">
		<h1 class="m-10 text-2xl leading-tight">Digital Logic Simulator</h1>
		<a href={makePath('/profile')} class="hero-btn custom-blur">Start Creating</a>
	</div>
</section>
