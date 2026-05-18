<script lang="ts">
  import { api } from "$lib/api";

  interface ScoreEntry {
    id: number;
    username: string;
    score: number;
    quiz_title: string;
    created_at: string;
  }

  let { scores, refresh } = $props<{
    scores: ScoreEntry[];
    refresh: () => Promise<void>;
  }>();

  let resetConfirm = $state(false);

  async function resetAllScores() {
    try {
      await api.delete("/api/admin/scores/reset");
      resetConfirm = false;
      await refresh();
    } catch (err) {
      console.error(err);
    }
  }
</script>

<div class="space-y-4">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold text-white">Semua Nilai</h2>
    {#if !resetConfirm}
      <button
        onclick={() => (resetConfirm = true)}
        class="rounded-lg bg-red-500/80 px-3 py-1.5 text-sm text-white hover:bg-red-500 transition"
        >Reset Semua Nilai</button
      >
    {:else}
      <div class="flex items-center gap-2">
        <span class="text-sm text-red-400">Yakin?</span>
        <button
          onclick={resetAllScores}
          class="rounded-lg bg-red-600 px-3 py-1.5 text-sm text-white hover:bg-red-500 transition"
          >Ya, Reset</button
        >
        <button
          onclick={() => (resetConfirm = false)}
          class="rounded-lg border border-slate-600 px-3 py-1.5 text-sm text-slate-300 hover:bg-slate-700 transition"
          >Batal</button
        >
      </div>
    {/if}
  </div>
  {#if scores.length === 0}
    <p class="text-sm text-slate-500">Belum ada nilai.</p>
  {:else}
    <div
      class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm"
    >
      <table class="w-full text-left text-sm">
        <thead class="border-b border-slate-700">
          <tr>
            <th class="px-4 py-3 font-medium text-slate-400">ID</th>
            <th class="px-4 py-3 font-medium text-slate-400">Username</th>
            <th class="px-4 py-3 font-medium text-slate-400">Kuis</th>
            <th class="px-4 py-3 font-medium text-slate-400">Skor</th>
            <th class="px-4 py-3 font-medium text-slate-400">Tanggal</th>
          </tr>
        </thead>
        <tbody>
          {#each scores as s}
            <tr class="border-b border-slate-700/50 last:border-0">
              <td class="px-4 py-3 text-slate-500">{s.id}</td>
              <td class="px-4 py-3 text-slate-300">{s.username}</td>
              <td class="px-4 py-3 text-slate-300">{s.quiz_title}</td>
              <td class="px-4 py-3 font-semibold text-indigo-400"
                >{s.score}</td
              >
              <td class="px-4 py-3 text-slate-500"
                >{new Date(s.created_at).toLocaleDateString("id-ID")}</td
              >
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
