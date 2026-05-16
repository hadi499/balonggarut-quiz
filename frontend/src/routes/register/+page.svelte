<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto } from "$app/navigation";
  import { fly, slide, fade } from "svelte/transition";

  let form = $state({ username: "", password: "" });
  let error = $state("");
  let message = $state("");
  let submitting = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    message = "";

    const pwd = form.password;
    if (pwd.length < 8) {
      error = "Password minimal 8 karakter.";
      return;
    }
    if (!/[A-Z]/.test(pwd)) {
      error = "Password harus mengandung minimal satu huruf besar.";
      return;
    }
    if (!/[a-z]/.test(pwd)) {
      error = "Password harus mengandung minimal satu huruf kecil.";
      return;
    }
    if (!/[!@#$%^&*(),.?":{}|<>\-\_=\+\[\]\\\/`~]/.test(pwd)) {
      error = "Password harus mengandung minimal satu karakter spesial.";
      return;
    }

    submitting = true;

    try {
      const data = await api.post<{ message: string }>(
        "/api/auth/register",
        form,
      );
      message = data.message || "Registration successful! Silakan login.";
      form = { username: "", password: "" };
      setTimeout(() => goto("/login"), 1500);
    } catch (err) {
      error = (err as Error).message;
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Daftar - Balonggarut Quiz</title>
</svelte:head>

<div
  in:fly={{ y: 20, duration: 400 }}
  out:fade={{ duration: 200 }}
  class="mx-auto max-w-md space-y-6 rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-6 mt-14"
>
  <h1 class="text-3xl font-bold text-white">Register</h1>

  {#if error}
    <div
      transition:slide={{ duration: 300 }}
      class="rounded-md bg-red-900/30 border border-red-500/30 p-3 text-sm text-red-400"
    >
      {error}
    </div>
  {/if}
  {#if message}
    <div
      transition:slide={{ duration: 300 }}
      class="rounded-md bg-green-900/30 border border-green-500/30 p-3 text-sm text-green-400"
    >
      {message}
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
        required
        class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-white placeholder-slate-500 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 focus:outline-none"
      />
    </div>

    <button
      type="submit"
      disabled={submitting}
      class="w-full rounded-xl bg-indigo-500 px-4 py-2.5 text-white font-medium hover:bg-indigo-400 transition disabled:opacity-50 shadow-lg shadow-indigo-500/25"
    >
      {submitting ? "Registering..." : "Register"}
    </button>
  </form>

  <p class="text-center text-sm text-slate-500">
    Sudah punya akun? <a
      href="/login"
      class="text-indigo-400 hover:text-indigo-300">Login</a
    >
  </p>
</div>
