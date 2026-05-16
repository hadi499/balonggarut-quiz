<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";

  let mobileMenuOpen = $state(false);

  async function handleLogout() {
    try {
      await api.post("/api/auth/logout", {});
    } catch {}
    auth.logout();
    goto("/");
  }
</script>

<nav
  class="fixed top-0 left-0 right-0 z-50 border-b border-slate-700/50 bg-slate-900/80 backdrop-blur-md"
>
  <div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
    <a href="/" class="text-xl font-bold text-white">Les Balonggarut</a>

    <button
      class="rounded-md p-2 text-slate-400 hover:bg-slate-800 md:hidden"
      aria-label="Toggle menu"
      onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
    >
      <svg
        class="h-6 w-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 12h16M4 18h16"
        />
      </svg>
    </button>

    <div class="hidden items-center gap-4 md:flex">
      <a
        href="/quizzes"
        class="text-md transition-colors {page.url.pathname.startsWith('/quizzes') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
        >Kuis</a
      >
      {#if auth.isLoggedIn}
        <a
          href="/dashboard"
          class="text-md transition-colors {page.url.pathname.startsWith('/dashboard') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
          >Dashboard</a
        >
        {#if auth.isTeacher}
          <a
            href="/admin"
            class="text-md transition-colors {page.url.pathname.startsWith('/admin') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
            >Admin</a
          >
        {/if}
        <span class="text-md text-slate-300"
          >{auth.username}
          <span class="text-xs text-slate-500">({auth.role})</span></span
        >
        <button
          onclick={handleLogout}
          class="rounded-md bg-red-500/80 px-4 py-1.5 text-md text-white hover:bg-red-500 transition"
        >
          Logout
        </button>
      {:else}
        <a
          href="/login"
          class="rounded-md bg-indigo-500 px-4 py-1.5 text-md text-white hover:bg-indigo-400 transition"
        >
          Login
        </a>
        <a
          href="/register"
          class="rounded-md border border-slate-600 px-4 py-1.5 text-md text-slate-300 hover:border-slate-500 hover:bg-slate-800 transition"
        >
          Register
        </a>
      {/if}
    </div>
  </div>

  {#if mobileMenuOpen}
    <div class="border-t border-slate-700/50 px-4 py-3 space-y-2 md:hidden">
      <a
        href="/quizzes"
        class="block text-md transition-colors {page.url.pathname.startsWith('/quizzes') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
        >Kuis</a
      >
      {#if auth.isLoggedIn}
        <a
          href="/dashboard"
          class="block text-md transition-colors {page.url.pathname.startsWith('/dashboard') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
          >Dashboard</a
        >
        {#if auth.isTeacher}
          <a
            href="/admin"
            class="block text-md transition-colors {page.url.pathname.startsWith('/admin') ? 'text-indigo-400 font-semibold' : 'text-slate-400 hover:text-white'}"
            >Admin</a
          >
        {/if}
        <p class="text-md text-slate-300">
          {auth.username}
          <span class="text-xs text-slate-500">({auth.role})</span>
        </p>
        <button
          onclick={handleLogout}
          class="block w-full rounded-md bg-red-500/80 py-2 text-md text-white hover:bg-red-500 transition"
        >
          Logout
        </button>
      {:else}
        <a
          href="/login"
          class="block w-full rounded-md bg-indigo-500 py-2 text-center text-md text-white hover:bg-indigo-400 transition"
        >
          Login
        </a>
        <a
          href="/register"
          class="mt-2 block w-full rounded-md border border-slate-600 py-2 text-center text-md text-slate-300 hover:border-slate-500 hover:bg-slate-800 transition"
        >
          Register
        </a>
      {/if}
    </div>
  {/if}
</nav>
