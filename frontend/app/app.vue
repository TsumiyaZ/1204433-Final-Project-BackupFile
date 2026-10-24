<template>
  <main
    v-if="embeddingStatus.state !== 'ready'"
    class="fixed inset-0 z-[100] grid place-items-center bg-base-200 px-6 py-10"
    :role="embeddingStatus.state === 'error' ? 'alert' : 'status'"
    :aria-live="embeddingStatus.state === 'error' ? 'assertive' : 'polite'"
  >
    <section
      class="w-full max-w-lg rounded-3xl border border-base-300 bg-base-100 px-8 py-10 text-center shadow-xl sm:px-12"
    >
      <div
        class="mx-auto mb-6 grid size-16 place-items-center rounded-2xl bg-primary/10 text-primary"
        aria-hidden="true"
      >
        <svg
          v-if="embeddingStatus.state === 'error'"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          class="size-8"
        >
          <path d="M12 3 2.8 19h18.4L12 3Z" />
          <path d="M12 9v4m0 3h.01" />
        </svg>
        <span v-else class="loading loading-spinner loading-lg text-primary" />
      </div>

      <h1 class="text-2xl font-bold text-base-content">
        {{ embeddingStatus.state === "error" ? "เตรียมโมเดลไม่สำเร็จ" : "กำลังเตรียมโมเดลค้นหารูป" }}
      </h1>

      <p class="mx-auto mt-3 max-w-md text-base leading-7 text-base-content/75">
        {{ embeddingStatus.message }}
      </p>

      <div v-if="embeddingStatus.state !== 'error'" class="mx-auto mt-5 max-w-md text-left">
        <div v-if="embeddingStatus.progressVisible" class="mb-2 flex justify-end text-xs font-medium text-base-content/65">
          {{ embeddingStatus.progress }}%
        </div>
        <div
          class="h-2 overflow-hidden rounded-full bg-base-300"
          role="progressbar"
          :aria-valuenow="embeddingStatus.progressVisible ? embeddingStatus.progress : undefined"
          :aria-valuemin="embeddingStatus.progressVisible ? 0 : undefined"
          :aria-valuemax="embeddingStatus.progressVisible ? 100 : undefined"
          :aria-label="embeddingStatus.message"
        >
          <div
            v-if="embeddingStatus.progressVisible"
            class="h-full rounded-full bg-primary transition-[width] duration-300 ease-out"
            :style="{ width: `${embeddingStatus.progress}%` }"
          />
          <div v-else class="h-full w-1/3 animate-pulse rounded-full bg-primary/80" />
        </div>
      </div>

      <p
        v-if="embeddingStatus.state !== 'error'"
        class="mx-auto mt-5 max-w-md text-sm leading-6 text-base-content/60"
      >
        การเปิดครั้งแรกต้องเชื่อมต่ออินเทอร์เน็ตเพื่อดาวน์โหลดโมเดล หลังจากเตรียมเสร็จแล้ว
        โมเดลจะอยู่ในเครื่องและใช้ครั้งต่อไปได้
      </p>

      <p
        v-else
        class="mt-5 rounded-xl border border-error/25 bg-error/5 p-4 text-sm leading-6 text-base-content/75"
      >
        ตรวจสอบอินเทอร์เน็ตและตำแหน่ง llama-server จากนั้นปิดแล้วเปิดแอปใหม่
      </p>
    </section>
  </main>

  <div v-else>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { GetEmbeddingStatus } from "../wailsjs/go/main/App";

type EmbeddingStatus = {
  state: string;
  message: string;
  progress: number;
  progressVisible: boolean;
};

const embeddingStatus = ref<EmbeddingStatus>({
  state: "loading",
  message: "กำลังเชื่อมต่อระบบโมเดล...",
  progress: 0,
  progressVisible: false,
});

let statusTimer: ReturnType<typeof setInterval> | undefined;

async function refreshEmbeddingStatus() {
  try {
    embeddingStatus.value = await GetEmbeddingStatus();
  } catch {
    embeddingStatus.value = {
      state: "error",
      message: "ไม่สามารถตรวจสอบสถานะโมเดลจากแอปได้",
      progress: 0,
      progressVisible: false,
    };
  }
}

onMounted(() => {
  void refreshEmbeddingStatus();
  statusTimer = setInterval(() => void refreshEmbeddingStatus(), 1000);
});

onUnmounted(() => {
  if (statusTimer) clearInterval(statusTimer);
});
</script>
