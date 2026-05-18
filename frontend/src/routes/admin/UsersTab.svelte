<script lang="ts">
  import { api } from "$lib/api";
  import { auth } from "$lib/stores/auth.svelte";

  interface User {
    id: number;
    username: string;
    role: string;
  }

  let { users, refresh } = $props<{
    users: User[];
    refresh: () => Promise<void>;
  }>();

  let deleteUserId = $state<number | null>(null);

  async function deleteUser() {
    if (!deleteUserId) return;
    try {
      await api.delete(`/api/admin/users/${deleteUserId}`);
      deleteUserId = null;
      await refresh();
    } catch (err) {
      console.error(err);
    }
  }

  async function changeUserRole(id: number, newRole: string) {
    try {
      await api.put(`/api/admin/users/${id}/role`, { role: newRole });
      await refresh();
    } catch (err) {
      console.error(err);
    }
  }
</script>

<div class="space-y-4">
  <h2 class="text-lg font-semibold text-white">Semua Pengguna</h2>
  {#if users.length === 0}
    <p class="text-sm text-slate-500">Belum ada pengguna.</p>
  {:else}
    <div
      class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm"
    >
      <table class="w-full text-left text-sm">
        <thead class="border-b border-slate-700">
          <tr>
            <th class="px-4 py-3 font-medium text-slate-400">ID</th>
            <th class="px-4 py-3 font-medium text-slate-400">Username</th>
            <th class="px-4 py-3 font-medium text-slate-400">Role</th>
            <th class="px-4 py-3 font-medium text-slate-400">Aksi</th>
          </tr>
        </thead>
        <tbody>
          {#each users as u}
            <tr class="border-b border-slate-700/50 last:border-0">
              <td class="px-4 py-3 text-slate-500">{u.id}</td>
              <td class="px-4 py-3 text-slate-300 font-medium"
                >{u.username}</td
              >
              <td class="px-4 py-3">
                <select
                  class="rounded-md border border-slate-600 bg-slate-800 px-2 py-1 text-xs text-white focus:border-indigo-500 focus:outline-none"
                  value={u.role}
                  onchange={(e) =>
                    changeUserRole(u.id, e.currentTarget.value)}
                  disabled={u.username === auth.username}
                >
                  <option value="student">Student</option>
                  <option value="teacher">Teacher</option>
                </select>
              </td>
              <td class="px-4 py-3">
                {#if u.username !== auth.username}
                  <button
                    onclick={() => (deleteUserId = u.id)}
                    class="rounded px-2 py-1 text-xs text-red-400 hover:bg-red-500/10 transition"
                    >Hapus</button
                  >
                {:else}
                  <span class="text-xs text-slate-500 italic">Anda</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if deleteUserId}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    onclick={() => (deleteUserId = null)}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div
      class="rounded-xl border border-slate-600 bg-slate-800 p-6 shadow-2xl mx-4"
      onclick={(e) => e.stopPropagation()}
    >
      <h3 class="text-lg font-semibold text-white">
        Konfirmasi Hapus Pengguna
      </h3>
      <p class="mt-2 text-sm text-slate-400 max-w-md">
        Yakin menghapus pengguna ini? <strong class="text-red-400"
          >Semua nilai skor dari pengguna ini juga akan ikut terhapus
          selamanya.</strong
        >
      </p>
      <div class="mt-6 flex justify-end gap-3">
        <button
          onclick={() => (deleteUserId = null)}
          class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
          >Batal</button
        >
        <button
          onclick={deleteUser}
          class="rounded-lg bg-red-500 px-4 py-2 text-sm text-white hover:bg-red-400 transition"
          >Hapus</button
        >
      </div>
    </div>
  </div>
{/if}
