<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  $effect(() => {
    if (!auth.isLoggedIn) {
      goto("/login");
    }
  });

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  interface MyScore {
    id: number;
    username: string;
    quiz_id: number;
    score: number;
    created_at: string;
  }

  interface LeaderboardEntry {
    id: number;
    username: string;
    score: number;
    quiz_title: string;
    created_at: string;
  }

  let quizzes = $state<Quiz[]>([]);
  let myScores = $state<MyScore[]>([]);
  let leaderboard = $state<LeaderboardEntry[]>([]);
  let loading = $state(true);
  let showAllScoresModal = $state(false);
  let deleteAccountConfirm = $state(false);
  let deletingAccount = $state(false);

  async function deleteAccount() {
    deletingAccount = true;
    try {
      await api.delete("/api/users/me");
      auth.logout();
      goto("/");
    } catch (e) {
      console.error(e);
      deletingAccount = false;
      deleteAccountConfirm = false;
    }
  }

  onMount(async () => {
    try {
      const [q, s, l] = await Promise.all([
        api.get<Quiz[]>("/api/quizzes"),
        api.get<MyScore[]>("/api/scores"),
        api.get<LeaderboardEntry[]>("/api/leaderboard"),
      ]);
      quizzes = q || [];
      myScores = s || [];
      leaderboard = l || [];
    } catch {
    } finally {
      loading = false;
    }
  });

  function getQuizTitle(id: number): string {
    return quizzes.find((q) => q.id === id)?.title ?? `Quiz #${id}`;
  }
</script>

<svelte:head>
  <title>Dashboard - Balonggarut Quiz</title>
</svelte:head>

{#if !auth.isLoggedIn}
  <p class="text-slate-400">Redirecting...</p>
{:else if loading}
  <div class="flex justify-center py-12">
    <div
      class="h-8 w-8 animate-spin rounded-full border-2 border-indigo-500 border-t-transparent"
    ></div>
  </div>
{:else}
  <div class="space-y-8">
    <h1 class="text-3xl font-bold text-white">Dashboard</h1>

    <div class="grid gap-6 lg:grid-cols-2">
      <div
        class="rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-md p-5"
      >
        <h2 class="mb-4 text-lg font-semibold text-white">Riwayat Nilai</h2>
        {#if myScores.length === 0}
          <p class="text-md text-slate-500">
            Belum ada nilai. Kerjakan kuis untuk melihat skor kamu.
          </p>
        {:else}
          <div class="space-y-2">
            {#each myScores.slice(0, 5) as score}
              <div
                class="flex items-center justify-between rounded-lg border border-slate-700 bg-slate-800/50 px-3 py-2"
              >
                <div>
                  <p class="text-md font-medium text-slate-200">
                    {getQuizTitle(score.quiz_id)}
                  </p>
                  <p class="text-xs text-slate-500">
                    {new Date(score.created_at).toLocaleDateString("id-ID")}
                  </p>
                </div>
                <span class="text-lg font-bold text-indigo-400"
                  >{score.score}</span
                >
              </div>
            {/each}
          </div>
          {#if myScores.length > 5}
            <div class="mt-4 flex justify-center">
              <button
                onclick={() => (showAllScoresModal = true)}
                class="text-sm font-medium text-indigo-400 hover:text-indigo-300 hover:underline transition"
              >
                Lihat Semua Skor &rarr;
              </button>
            </div>
          {/if}
        {/if}
      </div>

      <div
        class="rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-md p-5"
      >
        <h2 class="mb-4 text-lg font-semibold text-white">Kuis Tersedia</h2>
        {#if quizzes.length === 0}
          <p class="text-md text-slate-500">Belum ada kuis.</p>
        {:else}
          <div class="space-y-2">
            {#each quizzes.slice(0, 5) as quiz}
              <a
                href="/quizzes/{quiz.id}"
                class="flex items-center justify-between rounded-lg border border-slate-700 px-3 py-2 hover:bg-slate-700/50 transition"
              >
                <div>
                  <p class="text-md font-medium text-slate-200">{quiz.title}</p>
                  <span class="text-xs text-slate-500"
                    >{quiz.category} &middot; {quiz.questions?.length ?? 0} soal</span
                  >
                </div>
                <span class="text-xs text-indigo-400">Kerjakan &rarr;</span>
              </a>
            {/each}
          </div>
          <div class="mt-4">
            <a
              href="/quizzes"
              class="flex w-full items-center justify-center rounded-lg border border-indigo-500/30 bg-indigo-500/10 px-4 py-2 text-sm font-medium text-indigo-400 transition hover:bg-indigo-500/20"
            >
              Lihat Semua Kuis &rarr;
            </a>
          </div>
        {/if}
      </div>
    </div>

    <div>
      <h2 class="mb-4 text-2xl font-bold text-white">Leaderboard</h2>
      {#if leaderboard.length === 0}
        <p class="text-slate-500">Belum ada skor.</p>
      {:else}
        <div
          class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-md"
        >
          <table class="w-full text-left text-md">
            <thead class="border-b border-slate-700">
              <tr>
                <th class="px-4 py-3 font-medium text-slate-400">#</th>
                <th class="px-4 py-3 font-medium text-slate-400">Username</th>
                <th class="px-4 py-3 font-medium text-slate-400">Kuis</th>
                <th class="px-4 py-3 font-medium text-slate-400">Skor</th>
              </tr>
            </thead>
            <tbody>
              {#each leaderboard as entry, i}
                <tr class="border-b border-slate-700/50 last:border-0">
                  <td class="px-4 py-3 font-medium text-slate-200">{i + 1}</td>
                  <td class="px-4 py-3 text-slate-300">{entry.username}</td>
                  <td class="px-4 py-3 text-slate-300">{entry.quiz_title}</td>
                  <td class="px-4 py-3 font-semibold text-indigo-400"
                    >{entry.score}</td
                  >
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>

    <!-- Pengaturan Akun -->
    <div class="mt-8 rounded-xl border border-red-500/30 bg-red-500/5 backdrop-blur-md p-6">
      <h2 class="text-xl font-bold text-red-400 mb-2">Pengaturan Akun</h2>
      <p class="text-slate-400 text-sm mb-4">
        Tindakan ini akan menghapus akun Anda beserta seluruh riwayat nilai kuis yang pernah Anda kerjakan secara permanen. Tindakan ini tidak dapat dibatalkan.
      </p>
      <button
        onclick={() => (deleteAccountConfirm = true)}
        class="rounded-lg bg-red-600/20 border border-red-500/50 px-4 py-2 text-sm font-medium text-red-400 transition hover:bg-red-600 hover:text-white"
      >
        Hapus Akun Saya
      </button>
    </div>
  </div>

  {#if showAllScoresModal}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4" onclick={() => (showAllScoresModal = false)}>
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div class="flex flex-col max-h-[80vh] w-full max-w-2xl rounded-xl border border-slate-600 bg-slate-800 p-6 shadow-2xl" onclick={(e) => e.stopPropagation()}>
        <div class="mb-4 flex items-center justify-between">
          <h3 class="text-xl font-bold text-white">Semua Riwayat Nilai</h3>
          <button onclick={() => (showAllScoresModal = false)} class="text-slate-400 hover:text-white">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="flex-1 overflow-y-auto pr-2 space-y-2">
          {#each myScores as score}
            <div class="flex items-center justify-between rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-3">
              <div>
                <p class="text-md font-medium text-slate-200">
                  {getQuizTitle(score.quiz_id)}
                </p>
                <p class="text-sm text-slate-500">
                  {new Date(score.created_at).toLocaleDateString("id-ID", { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                </p>
              </div>
              <span class="text-xl font-bold text-indigo-400">{score.score}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  {#if deleteAccountConfirm}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4" onclick={() => (deleteAccountConfirm = false)}>
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div class="rounded-xl border border-red-600/50 bg-slate-800 p-6 shadow-2xl w-full max-w-md" onclick={(e) => e.stopPropagation()}>
        <h3 class="text-xl font-bold text-white mb-2">Konfirmasi Hapus Akun</h3>
        <p class="text-sm text-slate-300 mb-6">
          Apakah Anda benar-benar yakin ingin menghapus akun Anda? <strong class="text-red-400">Seluruh data dan skor Anda akan hilang selamanya.</strong>
        </p>
        <div class="flex justify-end gap-3">
          <button
            onclick={() => (deleteAccountConfirm = false)}
            disabled={deletingAccount}
            class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={deleteAccount}
            disabled={deletingAccount}
            class="rounded-lg bg-red-600 px-4 py-2 text-sm text-white hover:bg-red-500 transition disabled:opacity-50 flex items-center gap-2"
          >
            {#if deletingAccount}
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
              Menghapus...
            {:else}
              Ya, Hapus Akun
            {/if}
          </button>
        </div>
      </div>
    </div>
  {/if}
{/if}
