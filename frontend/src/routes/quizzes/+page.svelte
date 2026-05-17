<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api";

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  let quizzes = $state<Quiz[]>([]);
  let loading = $state(true);
  let searchQuery = $state("");
  let selectedCategory = $state("Semua");

  onMount(async () => {
    try {
      quizzes = (await api.get<Quiz[]>("/api/quizzes")) || [];
    } catch {}
    loading = false;
  });

  let categories = $derived([
    "Semua",
    ...new Set(quizzes.map((q) => q.category)),
  ]);

  let filteredQuizzes = $derived(
    quizzes.filter((q) => {
      const matchSearch =
        q.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        q.category.toLowerCase().includes(searchQuery.toLowerCase());
      const matchCategory =
        selectedCategory === "Semua" || q.category === selectedCategory;
      return matchSearch && matchCategory;
    }),
  );
</script>

<svelte:head>
  <title>Semua Kuis - Balonggarut Quiz</title>
</svelte:head>

<div class="min-h-[calc(100vh-4rem)] py-8">
  <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
    <!-- Search Bar -->
    <div class="mb-6 max-w-2xl mx-auto">
      <div class="relative w-full">
        <div
          class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4"
        >
          <svg
            class="h-5 w-5 text-slate-500"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <input
          type="text"
          bind:value={searchQuery}
          class="block w-full rounded-2xl border-0 bg-slate-800/80 py-3.5 pl-11 px-4 text-white ring-1 ring-inset ring-slate-700/50 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-indigo-500 shadow-sm backdrop-blur-sm transition-all text-base"
          placeholder="Cari kuis..."
        />
      </div>
    </div>

    <!-- Category Filters -->
    <div class="flex flex-wrap items-center justify-center gap-2 mb-12">
      {#each categories as category}
        <button
          onclick={() => (selectedCategory = category)}
          class="inline-flex items-center rounded-lg px-3 py-1.5 text-sm font-semibold transition-all duration-200 border {selectedCategory ===
          category
            ? 'bg-indigo-500/20 text-indigo-300 border-indigo-500/40 shadow-sm shadow-indigo-500/10'
            : 'bg-slate-800/40 text-slate-400 border-slate-700/50 hover:bg-slate-800 hover:text-slate-300 hover:border-slate-600'}"
        >
          {category}
        </button>
      {/each}
    </div>

    <!-- Quiz Grid -->
    {#if loading}
      <div class="flex justify-center py-20">
        <div
          class="h-10 w-10 animate-spin rounded-full border-4 border-slate-700 border-t-indigo-500"
        ></div>
      </div>
    {:else if filteredQuizzes.length === 0}
      <div
        class="text-center py-20 bg-slate-800/30 rounded-3xl border border-slate-700/50 border-dashed"
      >
        <p class="text-lg text-slate-400">Tidak ada kuis yang ditemukan.</p>
        <button
          onclick={() => {
            searchQuery = "";
            selectedCategory = "Semua";
          }}
          class="text-sm font-medium text-indigo-400 hover:text-indigo-300 mt-2 transition-colors"
        >
          Reset pencarian
        </button>
      </div>
    {:else}
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {#each filteredQuizzes as quiz}
          <a
            href="/quizzes/{quiz.id}"
            class="group relative flex flex-col justify-between rounded-3xl bg-slate-800/80 p-8 ring-1 ring-slate-700/50 transition-all duration-300 hover:ring-indigo-500/50 hover:bg-slate-800 hover:-translate-y-1 hover:shadow-2xl hover:shadow-indigo-500/10 overflow-hidden"
          >
            <div
              class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"
            ></div>

            <div class="relative z-10">
              <div class="flex items-start justify-between gap-x-4 mb-4">
                <span
                  class="inline-flex items-center rounded-lg bg-indigo-500/10 px-2.5 py-1 text-xs font-semibold text-indigo-400 border border-indigo-500/20"
                >
                  {quiz.category}
                </span>
              </div>
              <h3
                class="text-xl font-bold text-white group-hover:text-indigo-400 transition-colors mb-2"
              >
                {quiz.title}
              </h3>
              <p class="text-sm text-slate-400 line-clamp-2">
                Uji kemampuanmu dalam {quiz.category.toLowerCase()} dengan {quiz
                  .questions?.length ?? 0} pertanyaan menantang.
              </p>
            </div>

            <div
              class="mt-8 flex items-center justify-between relative z-10 pt-6 border-t border-slate-700/50 group-hover:border-slate-700 transition-colors"
            >
              <div class="flex items-center gap-4">
                <div
                  class="flex items-center gap-1.5 text-xs font-medium text-slate-300"
                >
                  <svg
                    class="w-4 h-4 text-slate-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  {quiz.questions?.length ?? 0} Soal
                </div>
                <div
                  class="flex items-center gap-1.5 text-xs font-medium text-slate-300"
                >
                  <svg
                    class="w-4 h-4 text-slate-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  {quiz.timeLimit}s
                </div>
              </div>
              <div
                class="w-8 h-8 rounded-full bg-slate-700/50 flex items-center justify-center group-hover:bg-indigo-500 group-hover:text-white transition-colors text-slate-400"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </div>
            </div>
          </a>
        {/each}
      </div>
    {/if}
  </div>
</div>
