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

  // ── Mode toggle ──────────────────────────────────────────────
  let inputMode = $state<"form" | "json">("form");

  // ── Form mode state ──────────────────────────────────────────
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

  // ── JSON mode state ──────────────────────────────────────────
  const JSON_PLACEHOLDER = `[
  {
    "quiz_id": 1,
    "question": "Ibukota Indonesia adalah?",
    "options": ["Bandung", "Jakarta", "Surabaya", "Medan"],
    "answer": 1
  },
  {
    "quiz_id": 1,
    "question": "2 + 2 = ?",
    "options": ["3", "4", "5", "6"],
    "answer": 1
  }
]`;

  let jsonInput = $state("");
  let jsonError = $state("");
  let jsonPreview = $state<{ quiz_id: number; question: string; options: string[]; answer: number }[]>([]);
  let jsonSubmitting = $state(false);
  let jsonSuccess = $state("");

  // Parse JSON realtime untuk preview
  $effect(() => {
    if (inputMode !== "json" || !jsonInput.trim()) {
      jsonPreview = [];
      jsonError = "";
      return;
    }
    try {
      const parsed = JSON.parse(jsonInput);
      if (!Array.isArray(parsed)) {
        jsonError = "JSON harus berupa array [ ... ]";
        jsonPreview = [];
        return;
      }
      // Batas maksimal import sekaligus
      if (parsed.length > 10) {
        jsonError = `Maksimal 10 soal per import (ditemukan ${parsed.length} soal).`;
        jsonPreview = [];
        return;
      }
      for (let i = 0; i < parsed.length; i++) {
        const q = parsed[i];
        if (!q.quiz_id || typeof q.quiz_id !== "number") {
          jsonError = `Soal #${i + 1}: "quiz_id" wajib berupa angka`;
          jsonPreview = [];
          return;
        }
        if (!quizzes.find((qz: Quiz) => qz.id === q.quiz_id)) {
          jsonError = `Soal #${i + 1}: quiz_id ${q.quiz_id} tidak ditemukan. Gunakan ID dari daftar kuis yang ada.`;
          jsonPreview = [];
          return;
        }
        if (!q.question || typeof q.question !== "string") {
          jsonError = `Soal #${i + 1}: "question" wajib berupa string`;
          jsonPreview = [];
          return;
        }
        if (q.question.trim().length > 500) {
          jsonError = `Soal #${i + 1}: "question" terlalu panjang (maks 500 karakter)`;
          jsonPreview = [];
          return;
        }
        if (!Array.isArray(q.options) || q.options.length !== 4) {
          jsonError = `Soal #${i + 1}: "options" wajib berupa array dengan 4 elemen`;
          jsonPreview = [];
          return;
        }
        for (let j = 0; j < q.options.length; j++) {
          if (typeof q.options[j] !== "string" || q.options[j].trim() === "") {
            jsonError = `Soal #${i + 1}: opsi ${["A","B","C","D"][j]} tidak boleh kosong`;
            jsonPreview = [];
            return;
          }
          if (q.options[j].length > 200) {
            jsonError = `Soal #${i + 1}: opsi ${["A","B","C","D"][j]} terlalu panjang (maks 200 karakter)`;
            jsonPreview = [];
            return;
          }
        }
        if (typeof q.answer !== "number" || q.answer < 0 || q.answer > 3) {
          jsonError = `Soal #${i + 1}: "answer" wajib angka 0–3 (0=A, 1=B, 2=C, 3=D)`;
          jsonPreview = [];
          return;
        }
      }
      jsonError = "";
      jsonPreview = parsed;
    } catch {
      jsonError = "JSON tidak valid. Periksa sintaks.";
      jsonPreview = [];
    }
  });

  async function submitJson() {
    if (jsonPreview.length === 0) return;
    jsonSubmitting = true;
    jsonSuccess = "";
    jsonError = "";
    let successCount = 0;
    const errors: string[] = [];
    for (let i = 0; i < jsonPreview.length; i++) {
      try {
        await api.post("/api/questions", jsonPreview[i]);
        successCount++;
      } catch (err) {
        errors.push(`Soal #${i + 1}: ${(err as Error).message}`);
      }
    }
    jsonSubmitting = false;
    if (errors.length > 0) {
      jsonError = errors.join(" | ");
    }
    if (successCount > 0) {
      jsonSuccess = `${successCount} soal berhasil ditambahkan!`;
      jsonInput = "";
      await refresh();
    }
  }

  // ── Shared ───────────────────────────────────────────────────
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
    inputMode = "form";
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
  <!-- ── Left Panel: Form / JSON ─────────────────────────── -->
  <div
    class="min-w-0 h-fit rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5 lg:sticky lg:top-20"
  >
    <!-- Header + Mode Toggle -->
    <div class="mb-4 flex items-center justify-between gap-3">
      <h2 class="text-lg font-semibold text-white">
        {editingQuestionId ? "Edit Soal" : "Buat Soal Baru"}
      </h2>
      <!-- Toggle tabs -->
      {#if !editingQuestionId}
        <div class="flex rounded-lg border border-slate-600 bg-slate-900/60 p-0.5 text-xs font-medium shrink-0">
          <button
            type="button"
            onclick={() => { inputMode = "form"; jsonError = ""; jsonSuccess = ""; }}
            class="rounded-md px-3 py-1.5 transition-all {inputMode === 'form'
              ? 'bg-indigo-500 text-white shadow'
              : 'text-slate-400 hover:text-white'}"
          >
            Form
          </button>
          <button
            type="button"
            onclick={() => { inputMode = "json"; questionFormError = ""; }}
            class="rounded-md px-3 py-1.5 transition-all {inputMode === 'json'
              ? 'bg-indigo-500 text-white shadow'
              : 'text-slate-400 hover:text-white'}"
          >
            JSON
          </button>
        </div>
      {/if}
    </div>

    <!-- ── FORM MODE ── -->
    {#if inputMode === "form"}
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

    <!-- ── JSON MODE ── -->
    {:else}
      <div class="space-y-3">
        <!-- Hint / format -->
        <div class="rounded-lg border border-slate-700 bg-slate-900/50 px-3 py-2.5 text-xs text-slate-400 leading-relaxed">
          <p class="font-semibold text-slate-300 mb-1">Format JSON (array soal):</p>
          <p>• <code class="text-indigo-400">quiz_id</code> — ID kuis (angka)</p>
          <p>• <code class="text-indigo-400">question</code> — teks pertanyaan</p>
          <p>• <code class="text-indigo-400">options</code> — array 4 pilihan [A, B, C, D]</p>
          <p>• <code class="text-indigo-400">answer</code> — indeks jawaban benar (0=A, 1=B, 2=C, 3=D)</p>
          <button
            type="button"
            onclick={() => (jsonInput = JSON_PLACEHOLDER)}
            class="mt-2 text-indigo-400 hover:text-indigo-300 underline underline-offset-2 transition"
          >
            Isi dengan contoh template →
          </button>
        </div>

        <!-- Textarea -->
        <textarea
          bind:value={jsonInput}
          placeholder={JSON_PLACEHOLDER}
          rows={12}
          class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-mono text-white placeholder-slate-600 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 resize-y"
        ></textarea>

        <!-- Error -->
        {#if jsonError}
          <p class="text-sm text-red-400">{jsonError}</p>
        {/if}

        <!-- Success -->
        {#if jsonSuccess}
          <p class="text-sm text-emerald-400">{jsonSuccess}</p>
        {/if}

        <!-- Preview -->
        {#if jsonPreview.length > 0}
          <div class="rounded-lg border border-emerald-500/30 bg-emerald-500/5 px-3 py-2.5">
            <p class="mb-2 text-xs font-semibold text-emerald-400">
              ✓ {jsonPreview.length} soal siap diimpor
            </p>
            <ul class="space-y-1 max-h-32 overflow-y-auto pr-1">
              {#each jsonPreview as q, i}
                <li class="text-xs text-slate-300 truncate">
                  <span class="text-slate-500 mr-1">#{i + 1}</span>
                  {q.question}
                  <span class="text-slate-500 ml-1">(Quiz #{q.quiz_id}, Jwb: {["A","B","C","D"][q.answer]})</span>
                </li>
              {/each}
            </ul>
          </div>
        {/if}

        <!-- Submit -->
        <button
          type="button"
          onclick={submitJson}
          disabled={jsonPreview.length === 0 || jsonSubmitting}
          class="w-full rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-400 transition disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {#if jsonSubmitting}
            <span class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
            Mengimpor...
          {:else}
            Import {jsonPreview.length > 0 ? `${jsonPreview.length} Soal` : "Soal"}
          {/if}
        </button>
      </div>
    {/if}
  </div>

  <!-- ── Right Panel: Question List ─────────────────────────── -->
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
