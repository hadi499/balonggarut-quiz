<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto, afterNavigate } from "$app/navigation";
  import { page } from "$app/state";

  let mobileMenuOpen = $state(false);

  afterNavigate(() => {
    mobileMenuOpen = false;
  });

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
    <!-- Brand & Navigation Links -->
    <div class="flex items-center gap-8">
      <a href="/" class="text-xl font-bold text-white">Les Balonggarut</a>
      
      <div class="hidden items-center gap-6 md:flex">
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
        {/if}
      </div>
    </div>

    <!-- Mobile Toggle -->
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

    <!-- Auth Actions -->
    <div class="hidden items-center gap-6 md:flex">
      {#if auth.isLoggedIn}
        <span class="text-md text-slate-300"
          >{auth.username}
          <span class="text-xs text-slate-500">({auth.role})</span></span
        >
        <button
          onclick={handleLogout}
          class="rounded-md bg-red-500/10 border border-red-500/20 px-4 py-1.5 text-sm font-medium text-red-400 hover:bg-red-500 hover:text-white transition-all"
        >
          Logout
        </button>
      {:else}
        <a
          href="/login"
          class="text-sm font-medium text-slate-300 hover:text-indigo-400 transition-colors"
        >
          Masuk
        </a>
        <a
          href="/register"
          class="rounded-full bg-indigo-500 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-indigo-500/20 hover:bg-indigo-400 hover:shadow-indigo-500/40 hover:-translate-y-0.5 transition-all"
        >
          Daftar
        </a>
      {/if}
    </div>
  </div>

  {#if mobileMenuOpen}
    <div class="border-t border-slate-700/50 px-4 py-4 space-y-4 md:hidden">
      <div class="flex flex-col space-y-3">
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
        {/if}
      </div>

      <div class="pt-4 border-t border-slate-700/50 flex flex-col space-y-3">
        {#if auth.isLoggedIn}
          <p class="text-md text-slate-300">
            {auth.username}
            <span class="text-xs text-slate-500">({auth.role})</span>
          </p>
          <button
            onclick={handleLogout}
            class="block w-full rounded-md bg-red-500/10 border border-red-500/20 py-2 text-center text-sm font-medium text-red-400 hover:bg-red-500 hover:text-white transition-all"
          >
            Logout
          </button>
        {:else}
          <a
            href="/login"
            class="block w-full rounded-full border border-slate-600 bg-transparent py-2.5 text-center text-sm font-medium text-slate-300 hover:bg-slate-800 transition-colors"
          >
            Masuk
          </a>
          <a
            href="/register"
            class="block w-full rounded-full bg-indigo-500 py-2.5 text-center text-sm font-semibold text-white shadow-lg shadow-indigo-500/20 hover:bg-indigo-400 transition-all"
          >
            Daftar Gratis
          </a>
        {/if}
      </div>
    </div>
  {/if}
</nav>
