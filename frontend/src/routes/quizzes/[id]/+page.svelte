<script lang="ts">
  import { page } from "$app/state";
  import { api } from "$lib/api";
  import { auth } from "$lib/stores/auth.svelte";
  import { onMount } from "svelte";

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: Question[];
  }

  interface Question {
    id: number;
    quiz_id: number;
    question: string;
    options: string[];
    answer: number;
  }

  let quiz = $state<Quiz | null>(null);
  let loading = $state(true);
  let error = $state("");

  let quizState = $state<"preview" | "playing" | "finished">("preview");
  let currentIndex = $state(0);
  let selectedAnswers = $state<number[]>([]);
  let score = $state(0);
  let submittingScore = $state(false);
  let scoreSubmitted = $state(false);
  let scoreError = $state("");
  let timeLeft = $state(0);
  let timerInterval: ReturnType<typeof setInterval> | null = null;

  $effect(() => {
    return () => {
      if (timerInterval) clearInterval(timerInterval);
    };
  });

  onMount(async () => {
    try {
      const id = page.params.id;
      quiz = await api.get<Quiz>(`/api/quizzes/${id}`);
      if (quiz) {
        selectedAnswers = new Array(quiz.questions.length).fill(-1);
      }
    } catch {
      error = "Kuis tidak ditemukan.";
    } finally {
      loading = false;
    }
  });

  function resetTimer() {
    if (timerInterval) clearInterval(timerInterval);
    if (!quiz) return;
    timeLeft = quiz.timeLimit;
    const q = quiz;
    timerInterval = setInterval(() => {
      if (timeLeft > 0) {
        timeLeft--;
      } else {
        if (currentIndex < q.questions.length - 1) {
          currentIndex++;
          resetTimer();
        } else {
          finishQuiz();
        }
      }
    }, 1000);
  }

  function startQuiz() {
    if (!quiz) return;
    quizState = "playing";
    currentIndex = 0;
    resetTimer();
  }

  function selectAnswer(index: number) {
    selectedAnswers[currentIndex] = index;
    if (quiz && currentIndex < quiz.questions.length - 1) {
      setTimeout(() => {
        currentIndex++;
        resetTimer();
      }, 100);
    } else {
      setTimeout(() => finishQuiz(), 100);
    }
  }

  function finishQuiz() {
    if (timerInterval) clearInterval(timerInterval);
    if (!quiz) return;

    let correct = 0;
    quiz.questions.forEach((q, i) => {
      if (selectedAnswers[i] === q.answer) correct++;
    });
    score = correct * 10;
    quizState = "finished";

    if (auth.isLoggedIn) {
      submitScore();
    }
  }

  async function submitScore() {
    if (!quiz) return;
    submittingScore = true;
    scoreError = "";
    try {
      await api.post("/api/scores", { quiz_id: quiz.id, score });
      scoreSubmitted = true;
    } catch (err) {
      scoreError = (err as Error).message;
    } finally {
      submittingScore = false;
    }
  }

  function formatTime(seconds: number): string {
    return `${seconds}s`;
  }

  let dotCount = $derived(quiz ? [...Array(quiz.questions.length)] : []);
</script>

<svelte:head>
  <title
    >{quiz
      ? quiz.title + " - Balonggarut Quiz"
      : "Kuis - Balonggarut Quiz"}</title
  >
</svelte:head>

{#if loading}
  <div class="flex justify-center py-12">
    <div
      class="h-8 w-8 animate-spin rounded-full border-2 border-indigo-500 border-t-transparent"
    ></div>
  </div>
{:else if error || !quiz}
  <div
    class="rounded-md bg-red-900/30 border border-red-500/30 p-4 text-red-400"
  >
    {error || "Kuis tidak ditemukan."}
  </div>
{:else if quizState === "preview"}
  <div class="mx-auto max-w-2xl space-y-6 mt-12">
    <div
      class="rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-6"
    >
      <h1 class="text-2xl font-bold text-white">{quiz.title}</h1>
      <span
        class="mt-2 inline-block rounded-lg bg-indigo-500/10 border border-indigo-500/20 px-3 py-1 text-sm text-indigo-400"
        >{quiz.category}</span
      >
      <div class="mt-4 space-y-2 text-slate-300">
        <p>Waktu: {quiz.timeLimit} detik / soal</p>
        <p>Jumlah Soal: {quiz.questions.length}</p>
      </div>
      <button
        onclick={startQuiz}
        class="mt-6 w-full rounded-xl bg-indigo-500 px-4 py-3 text-white font-medium hover:bg-indigo-400 transition shadow-lg shadow-indigo-500/25"
      >
        Mulai Kuis
      </button>
    </div>
  </div>
{:else if quizState === "playing"}
  <div class="mx-auto max-w-2xl space-y-6">
    <div class="flex items-center justify-between">
      <span class="text-lg text-slate-400"
        >Soal {currentIndex + 1} / {quiz.questions.length}</span
      >
      <span class="text-lg font-semibold text-indigo-400"
        >{formatTime(timeLeft)}</span
      >
    </div>

    <div class="flex justify-center gap-1.5">
      {#each dotCount as _, i}
        <span
          class="h-2 w-2 rounded-full transition-all duration-300 {i ===
          currentIndex
            ? 'bg-indigo-400 scale-125'
            : i < currentIndex
              ? 'bg-slate-500'
              : 'bg-slate-700'}"
        ></span>
      {/each}
    </div>

    <div class="h-2 w-full overflow-hidden rounded-full bg-slate-700">
      <div
        class="h-full rounded-full bg-indigo-500 transition-all duration-1000 ease-linear"
        style="width: {(timeLeft / quiz.timeLimit) * 100}%"
      ></div>
    </div>

    <div
      class="rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-6"
    >
      <h2 class="mb-4 text-lg font-semibold text-white">
        {quiz.questions[currentIndex].question}
      </h2>

      <div class="space-y-2">
        {#each quiz.questions[currentIndex].options as option, i}
          <button
            onclick={() => selectAnswer(i)}
            class="w-full rounded-lg border px-4 py-3 text-left transition {selectedAnswers[
              currentIndex
            ] === i
              ? 'border-indigo-500 bg-indigo-500/10 text-indigo-300'
              : 'border-slate-600 hover:border-slate-500 text-slate-300'}"
          >
            <span class="font-medium">{String.fromCharCode(65 + i)}.</span>
            {option}
          </button>
        {/each}
      </div>
    </div>
  </div>
{:else if quizState === "finished"}
  <div class="mx-auto max-w-2xl space-y-6">
    <div
      class="rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-6 text-center"
    >
      <h2 class="text-2xl font-bold text-white">Kuis Selesai!</h2>
      <p class="mt-2 text-lg text-slate-300">
        Skor kamu: <span class="font-bold text-indigo-400">{score}</span> / {quiz
          .questions.length * 10}
      </p>

      {#if auth.isLoggedIn}
        {#if submittingScore}
          <p class="mt-3 text-slate-400">Menyimpan skor...</p>
        {:else if scoreSubmitted}
          <p class="mt-3 text-green-400">Skor berhasil disimpan!</p>
        {:else if scoreError}
          <p class="mt-3 text-red-400">{scoreError}</p>
          <button
            onclick={submitScore}
            class="mt-2 rounded-xl bg-indigo-500 px-4 py-2 text-white hover:bg-indigo-400 transition"
          >
            Coba Lagi
          </button>
        {/if}
      {:else}
        <p class="mt-3 text-sm text-slate-500">Login untuk menyimpan skor.</p>
      {/if}
    </div>

    <div class="space-y-4">
      <h3 class="text-xl font-bold text-white">Review Jawaban</h3>
      {#each quiz.questions as q, i}
        <div
          class="rounded-lg border p-4 {selectedAnswers[i] === q.answer
            ? 'border-green-500/30 bg-green-500/5'
            : 'border-red-500/30 bg-red-500/5'}"
        >
          <div class="flex items-start gap-2">
            <span
              class="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-bold text-white {selectedAnswers[
                i
              ] === q.answer
                ? 'bg-green-500'
                : 'bg-red-500'}"
            >
              {selectedAnswers[i] === q.answer ? "✓" : "✗"}
            </span>
            <div class="flex-1">
              <p class="font-semibold text-white">{i + 1}. {q.question}</p>
              <div class="mt-2 space-y-1">
                {#each q.options as option, j}
                  <div
                    class="rounded-md px-3 py-1.5 text-sm {j === q.answer
                      ? 'bg-green-500/10 text-green-400 font-medium'
                      : j === selectedAnswers[i] && j !== q.answer
                        ? 'bg-red-500/10 text-red-400'
                        : 'text-slate-500'}"
                  >
                    {String.fromCharCode(65 + j)}. {option}
                    {#if j === q.answer && j === selectedAnswers[i]}
                      <span class="ml-1 text-xs">(Jawaban kamu)</span>
                    {:else if j === q.answer}
                      <span class="ml-1 text-xs">(Jawaban benar)</span>
                    {:else if j === selectedAnswers[i]}
                      <span class="ml-1 text-xs">(Jawaban kamu)</span>
                    {/if}
                  </div>
                {/each}
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <div class="mt-8 text-center">
      <a
        href="/quizzes"
        class="inline-flex items-center gap-2 rounded-xl border border-slate-600 bg-slate-800 px-6 py-3 text-sm font-medium text-slate-300 transition hover:border-slate-500 hover:bg-slate-700 hover:text-white"
      >
        <svg
          class="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 19l-7-7m0 0l7-7m-7 7h18"
          />
        </svg>
        Kembali ke Daftar Kuis
      </a>
    </div>
  </div>
{/if}
