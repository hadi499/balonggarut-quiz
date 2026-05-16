<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto } from "$app/navigation";
  import { fly, slide, fade } from "svelte/transition";

  let form = $state({ username: "", password: "" });
  let error = $state("");
  let submitting = $state(false);

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
      <input
        id="password"
        type="password"
        bind:value={form.password}
        autocomplete="off"
        required
        class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 focus:outline-none"
      />
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
