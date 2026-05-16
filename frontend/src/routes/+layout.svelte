<script lang="ts">
  import favicon from "$lib/assets/favicon.svg";
  import "../app.css";
  import Navbar from "$lib/components/Navbar.svelte";
  import { auth } from "$lib/stores/auth.svelte";

  let { children } = $props();
  auth.init();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<Navbar />
<main class="mx-auto max-w-6xl px-4 py-8 pt-20 text-slate-200">
  {@render children()}
</main>
