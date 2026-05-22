<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto } from "$app/navigation";
  import { fly, slide, fade } from "svelte/transition";

  let form = $state({ username: "", password: "" });
  let error = $state("");
  let submitting = $state(false);
  let showPassword = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    submitting = true;

    try {
      const data = await api.post<{ token: string }>("/api/auth/login", form);
      auth.login(data.token, form.username, "");

      const me = await api.get<{ username: string; role: string }>("/me");
      auth.login(data.token, me.username, me.role);

      goto("/");
    } catch (err) {
      error = (err as Error).message;
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Masuk - Balonggarut Quiz</title>
</svelte:head>

<div
  in:fly={{ y: 20, duration: 400 }}
  out:fade={{ duration: 200 }}
  class="mx-auto max-w-md space-y-6 rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-6 mt-14"
>
  <h1 class="text-3xl font-bold text-white">Login</h1>

  {#if error}
    <div
      transition:slide={{ duration: 300 }}
      class="rounded-md bg-red-900/30 border border-red-500/30 p-3 text-sm text-red-400"
    >
      {error}
    </div>
  {/if}

  <form onsubmit={handleSubmit} class="space-y-4" autocomplete="off">
    <div>
      <label
        for="username"
        class="mb-1 block text-sm font-medium text-slate-300">Username</label
      >
      <input
        id="username"
        type="text"
        bind:value={form.username}
        required
        class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 focus:outline-none"
      />
    </div>

    <div>
      <label
        for="password"
        class="mb-1 block text-sm font-medium text-slate-300">Password</label
      >
      <div class="relative">
        <input
          id="password"
          type={showPassword ? 'text' : 'password'}
          bind:value={form.password}
          autocomplete="off"
          required
          class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 pr-10 text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 focus:outline-none"
        />
        <button
          type="button"
          onclick={() => (showPassword = !showPassword)}
          class="absolute inset-y-0 right-0 flex items-center px-3 text-slate-400 hover:text-slate-200 transition-colors"
          aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
        >
          {#if showPassword}
            <!-- Eye-off icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/>
              <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/>
              <line x1="1" y1="1" x2="23" y2="23"/>
            </svg>
          {:else}
            <!-- Eye icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
          {/if}
        </button>
      </div>
    </div>

    <button
      type="submit"
      disabled={submitting}
      class="w-full rounded-xl bg-indigo-500 px-4 py-2.5 text-white font-medium hover:bg-indigo-400 transition disabled:opacity-50 shadow-lg shadow-indigo-500/25"
    >
      {submitting ? "Logging in..." : "Login"}
    </button>
  </form>

  <p class="text-center text-sm text-slate-500">
    Belum punya akun? <a
      href="/register"
      class="text-indigo-400 hover:text-indigo-300">Register</a
    >
  </p>
</div>
