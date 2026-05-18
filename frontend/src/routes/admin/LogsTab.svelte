<script lang="ts">
  import { api } from "$lib/api";

  interface ActivityLog {
    id: number;
    username: string;
    action: string;
    timestamp: string;
  }

  let logs = $state<ActivityLog[]>([]);
  let searchLog = $state("");
  let logsPage = $state(1);
  let logsTotalPages = $state(1);

  async function fetchLogs() {
    try {
      const res = await api.get<{ data: ActivityLog[]; totalPages: number }>(
        `/api/admin/logs?page=${logsPage}&limit=25&search=${searchLog}`,
      );
      logs = res.data || [];
      logsTotalPages = res.totalPages || 1;
    } catch (e) {
      console.error("Failed to fetch logs", e);
    }
  }

  $effect(() => {
    fetchLogs();
  });
</script>

<div class="space-y-4">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <h2 class="text-lg font-semibold text-white">Log Aktivitas</h2>

    <form
      class="relative w-full sm:w-64"
      onsubmit={(e) => {
        e.preventDefault();
        logsPage = 1;
        fetchLogs();
      }}
    >
      <input
        type="text"
        bind:value={searchLog}
        placeholder="Cari log (Enter)..."
        class="w-full rounded-lg border border-slate-600 bg-slate-800 py-1.5 pl-3 pr-8 text-sm text-white placeholder-slate-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
      />
      {#if searchLog}
        <button
          type="button"
          onclick={() => {
            searchLog = "";
            logsPage = 1;
            fetchLogs();
          }}
          class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
          aria-label="Clear search"
        >
          <svg
            class="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      {/if}
    </form>
  </div>

  {#if logs.length === 0}
    <p class="text-sm text-slate-500">Log tidak ditemukan.</p>
  {:else}
    <div
      class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm max-h-[600px] overflow-y-auto custom-scrollbar"
    >
      <table class="w-full text-left text-sm relative">
        <thead
          class="border-b border-slate-700 bg-slate-800/90 backdrop-blur-sm sticky top-0"
        >
          <tr>
            <th class="px-4 py-3 font-medium text-slate-400">Waktu</th>
            <th class="px-4 py-3 font-medium text-slate-400">Username</th>
            <th class="px-4 py-3 font-medium text-slate-400">Aktivitas</th>
          </tr>
        </thead>
        <tbody>
          {#each logs as log}
            <tr class="border-b border-slate-700/50 last:border-0">
              <td class="px-4 py-3 text-slate-400 whitespace-nowrap">
                {new Date(log.timestamp).toLocaleDateString("id-ID", {
                  timeZone: "Asia/Jakarta",
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })} WIB
              </td>
              <td class="px-4 py-3 text-slate-300 font-medium"
                >{log.username}</td
              >
              <td class="px-4 py-3">
                {#if log.action === "REGISTER"}
                  <span
                    class="inline-flex items-center rounded-md bg-green-500/10 px-2 py-1 text-xs font-medium text-green-400 ring-1 ring-inset ring-green-500/20"
                    >Registrasi Akun Baru</span
                  >
                {:else if log.action === "DELETE_ACCOUNT"}
                  <span
                    class="inline-flex items-center rounded-md bg-red-500/10 px-2 py-1 text-xs font-medium text-red-400 ring-1 ring-inset ring-red-500/20"
                    >Hapus Akun Sendiri</span
                  >
                {:else if log.action === "DELETED_BY_ADMIN"}
                  <span
                    class="inline-flex items-center rounded-md bg-orange-500/10 px-2 py-1 text-xs font-medium text-orange-400 ring-1 ring-inset ring-orange-500/20"
                    >Dihapus Admin</span
                  >
                {:else}
                  <span
                    class="inline-flex items-center rounded-md bg-slate-500/10 px-2 py-1 text-xs font-medium text-slate-400 ring-1 ring-inset ring-slate-500/20"
                    >{log.action}</span
                  >
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-2"
    >
      <span class="text-sm text-slate-400 text-center sm:text-left"
        >Halaman {logsPage} dari {logsTotalPages}</span
      >
      <div class="flex gap-2">
        <button
          disabled={logsPage <= 1}
          onclick={() => {
            logsPage--;
            fetchLogs();
          }}
          class="px-3 py-1 rounded border border-slate-600 text-sm text-slate-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-700 transition"
        >
          Sebelumnya
        </button>
        <button
          disabled={logsPage >= logsTotalPages}
          onclick={() => {
            logsPage++;
            fetchLogs();
          }}
          class="px-3 py-1 rounded border border-slate-600 text-sm text-slate-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-700 transition"
        >
          Selanjutnya
        </button>
      </div>
    </div>
  {/if}
</div>
