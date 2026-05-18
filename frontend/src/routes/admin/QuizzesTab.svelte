<script lang="ts">
  import { api } from "$lib/api";

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  let { quizzes, refresh } = $props<{
    quizzes: Quiz[];
    refresh: () => Promise<void>;
  }>();

  let quizForm = $state({ title: "", category: "", timeLimit: 60 });
  let editingQuizId = $state<number | null>(null);
  let quizFormError = $state("");

  const inputClass =
    "w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500";

  async function saveQuiz() {
    quizFormError = "";
    const body = {
      title: quizForm.title,
      category: quizForm.category,
      timeLimit: quizForm.timeLimit,
    };
    try {
      if (editingQuizId) {
        await api.put(`/api/quizzes/${editingQuizId}`, body);
      } else {
        await api.post("/api/quizzes", body);
      }
      resetQuizForm();
      await refresh();
    } catch (err) {
      quizFormError = (err as Error).message;
    }
  }

  function editQuiz(quiz: Quiz) {
    editingQuizId = quiz.id;
    quizForm = {
      title: quiz.title,
      category: quiz.category,
      timeLimit: quiz.timeLimit,
    };
  }

  function resetQuizForm() {
    editingQuizId = null;
    quizForm = { title: "", category: "", timeLimit: 60 };
    quizFormError = "";
  }

  async function deleteQuiz(id: number) {
    if (!confirm("Yakin ingin menghapus kuis ini beserta semua soalnya?")) return;
    try {
      await api.delete(`/api/quizzes/${id}`);
      await refresh();
    } catch (err) {
      console.error(err);
    }
  }
</script>

<div class="grid gap-6 lg:grid-cols-2 items-start">
  <div
    class="min-w-0 h-fit rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5 lg:sticky lg:top-20"
  >
    <h2 class="mb-4 text-lg font-semibold text-white">
      {editingQuizId ? "Edit Kuis" : "Buat Kuis Baru"}
    </h2>
    {#if quizFormError}
      <p class="mb-3 text-sm text-red-400">{quizFormError}</p>
    {/if}
    <form
      onsubmit={(e) => {
        e.preventDefault();
        saveQuiz();
      }}
      class="space-y-3"
      autocomplete="off"
    >
      <input
        placeholder="Judul"
        bind:value={quizForm.title}
        required
        class={inputClass}
      />
      <input
        placeholder="Kategori"
        bind:value={quizForm.category}
        required
        class={inputClass}
      />
      <input
        type="number"
        placeholder="Waktu per soal (detik)"
        bind:value={quizForm.timeLimit}
        required
        min="1"
        class={inputClass}
      />
      <div class="flex gap-2">
        <button
          type="submit"
          class="rounded-lg bg-indigo-500 px-4 py-2 text-sm text-white hover:bg-indigo-400 transition"
        >
          {editingQuizId ? "Update" : "Buat"}
        </button>
        {#if editingQuizId}
          <button
            type="button"
            onclick={resetQuizForm}
            class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
            >Batal</button
          >
        {/if}
      </div>
    </form>
  </div>

  <div
    class="min-w-0 flex flex-col h-[600px] rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5"
  >
    <h2 class="flex-none mb-4 text-lg font-semibold text-white">
      Daftar Kuis
    </h2>
    {#if quizzes.length === 0}
      <p class="text-sm text-slate-500">Belum ada kuis.</p>
    {:else}
      <div class="space-y-2 overflow-y-auto pr-2 flex-1">
        {#each quizzes as quiz}
          <div
            class="flex items-center justify-between rounded-lg border border-slate-700 px-3 py-2"
          >
            <div>
              <span class="text-md font-medium text-slate-200"
                >{quiz.title}</span
              >
              <span class="ml-2 text-sm text-slate-500"
                >{quiz.questions?.length ?? 0} soal</span
              >
              <span class="ml-2 text-sm text-slate-500"
                >id: {quiz.id}</span
              >
            </div>
            <div class="flex gap-1">
              <button
                onclick={() => editQuiz(quiz)}
                class="rounded px-2 py-1 text-sm text-indigo-400 hover:bg-indigo-500/10 transition"
                >Edit</button
              >
              <button
                onclick={() => deleteQuiz(quiz.id)}
                class="rounded px-2 py-1 text-sm text-red-400 hover:bg-red-500/10 transition"
                >Hapus</button
              >
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
