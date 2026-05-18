<script lang="ts">
  import { api } from "$lib/api";

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  interface Question {
    id: number;
    quiz_id: number;
    question: string;
    options: string[];
    answer: number;
  }

  let { questions, quizzes, refresh } = $props<{
    questions: Question[];
    quizzes: Quiz[];
    refresh: () => Promise<void>;
  }>();

  let questionForm = $state({
    quiz_id: 0,
    question: "",
    optionA: "",
    optionB: "",
    optionC: "",
    optionD: "",
    answer: 0,
  });
  let editingQuestionId = $state<number | null>(null);
  let questionFormError = $state("");
  let searchQuestion = $state("");

  function getQuizTitle(id: number): string {
    return quizzes.find((q: Quiz) => q.id === id)?.title ?? `Quiz #${id}`;
  }

  let filteredQuestions = $derived(
    questions.filter(
      (q: Question) =>
        q.question.toLowerCase().includes(searchQuestion.toLowerCase()) ||
        getQuizTitle(q.quiz_id)
          .toLowerCase()
          .includes(searchQuestion.toLowerCase()),
    ),
  );

  const inputClass =
    "w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500";

  async function saveQuestion() {
    questionFormError = "";
    const body = {
      quiz_id: questionForm.quiz_id,
      question: questionForm.question,
      options: [
        questionForm.optionA,
        questionForm.optionB,
        questionForm.optionC,
        questionForm.optionD,
      ],
      answer: questionForm.answer,
    };
    try {
      if (editingQuestionId) {
        await api.put(`/api/questions/${editingQuestionId}`, body);
      } else {
        await api.post("/api/questions", body);
      }
      resetQuestionForm();
      await refresh();
    } catch (err) {
      questionFormError = (err as Error).message;
    }
  }

  function editQuestion(q: Question) {
    editingQuestionId = q.id;
    questionForm = {
      quiz_id: q.quiz_id,
      question: q.question,
      optionA: q.options[0] || "",
      optionB: q.options[1] || "",
      optionC: q.options[2] || "",
      optionD: q.options[3] || "",
      answer: q.answer,
    };
  }

  function resetQuestionForm() {
    editingQuestionId = null;
    questionForm = {
      quiz_id: 0,
      question: "",
      optionA: "",
      optionB: "",
      optionC: "",
      optionD: "",
      answer: 0,
    };
    questionFormError = "";
  }

  async function deleteQuestion(id: number) {
    if (!confirm("Yakin ingin menghapus soal ini?")) return;
    try {
      await api.delete(`/api/questions/${id}`);
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
      {editingQuestionId ? "Edit Soal" : "Buat Soal Baru"}
    </h2>
    {#if questionFormError}
      <p class="mb-3 text-sm text-red-400">{questionFormError}</p>
    {/if}
    <form
      onsubmit={(e) => {
        e.preventDefault();
        saveQuestion();
      }}
      class="space-y-3"
      autocomplete="off"
    >
      <select bind:value={questionForm.quiz_id} class={inputClass}>
        <option value={0} disabled>Pilih Kuis</option>
        {#each quizzes as quiz}
          <option value={quiz.id}>{quiz.title}</option>
        {/each}
      </select>
      <input
        placeholder="Pertanyaan"
        bind:value={questionForm.question}
        required
        class={inputClass}
      />
      <input
        placeholder="A. Opsi jawaban"
        bind:value={questionForm.optionA}
        required
        class={inputClass}
      />
      <input
        placeholder="B. Opsi jawaban"
        bind:value={questionForm.optionB}
        required
        class={inputClass}
      />
      <input
        placeholder="C. Opsi jawaban"
        bind:value={questionForm.optionC}
        required
        class={inputClass}
      />
      <input
        placeholder="D. Opsi jawaban"
        bind:value={questionForm.optionD}
        required
        class={inputClass}
      />
      <select bind:value={questionForm.answer} class={inputClass}>
        <option value={0}>Jawaban: A</option>
        <option value={1}>Jawaban: B</option>
        <option value={2}>Jawaban: C</option>
        <option value={3}>Jawaban: D</option>
      </select>
      <div class="flex gap-2">
        <button
          type="submit"
          class="rounded-lg bg-indigo-500 px-4 py-2 text-sm text-white hover:bg-indigo-400 transition"
        >
          {editingQuestionId ? "Update" : "Buat"}
        </button>
        {#if editingQuestionId}
          <button
            type="button"
            onclick={resetQuestionForm}
            class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
            >Batal</button
          >
        {/if}
      </div>
    </form>
  </div>

  <div
    class="min-w-0 flex flex-col h-[650px] rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5"
  >
    <div
      class="flex-none flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4"
    >
      <h2 class="text-lg font-semibold text-white">Daftar Soal</h2>
      <div class="relative w-full sm:w-64 flex items-center">
        <svg
          class="absolute left-2.5 h-4 w-4 text-slate-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        <input
          type="text"
          placeholder="Cari soal atau kuis..."
          bind:value={searchQuestion}
          class="w-full rounded-lg border border-slate-600 bg-slate-800 py-1.5 pl-9 pr-8 text-sm text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        {#if searchQuestion}
          <button
            type="button"
            onclick={() => (searchQuestion = "")}
            class="absolute right-2.5 text-slate-400 hover:text-white"
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
      </div>
    </div>
    {#if questions.length === 0}
      <p class="text-sm text-slate-500">Belum ada soal.</p>
    {:else if filteredQuestions.length === 0}
      <p class="text-sm text-slate-500">Soal tidak ditemukan.</p>
    {:else}
      <div class="space-y-2 overflow-y-auto pr-2 flex-1">
        {#each filteredQuestions as question}
          <div
            class="flex items-center justify-between rounded-lg border border-slate-700 px-3 py-2"
          >
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-slate-200">
                {question.question}
              </p>
              <span class="text-xs text-slate-500"
                >{getQuizTitle(question.quiz_id)}</span
              >
            </div>
            <div class="ml-2 flex shrink-0 gap-1">
              <button
                onclick={() => editQuestion(question)}
                class="rounded px-2 py-1 text-sm text-indigo-400 hover:bg-indigo-500/10 transition"
                >Edit</button
              >
              <button
                onclick={() => deleteQuestion(question.id)}
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
