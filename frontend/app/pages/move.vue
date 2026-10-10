<template>
  <main class="min-h-[calc(100dvh-4rem)] bg-base-200/70 px-4 py-10 sm:px-6">
    <section class="mx-auto max-w-6xl">
      <header class="mb-8">
        <span
          class="inline-flex rounded-full border border-primary/20 bg-primary/8 px-4 py-2 text-sm font-semibold text-primary"
        >
          Move photos
        </span>

        <h1 class="mt-5 text-3xl font-bold tracking-tight sm:text-4xl">
          ย้ายรูปไป Destination
        </h1>

        <p class="mt-3 max-w-3xl text-base leading-7 text-base-content/65">
          ระบบจะตรวจชนิดไฟล์ ข้ามรูปที่ซ้ำ และลบต้นฉบับเมื่อปลายทางตรวจสอบผ่านแล้วเท่านั้น
        </p>
      </header>

      <div v-if="errorMessage" role="alert" class="alert alert-error mb-6">
        {{ errorMessage }}
      </div>

      <div
        v-if="selectedPhotos.length === 0 && !result"
        role="status"
        class="rounded-[1.5rem] border-2 border-dashed border-primary/20 bg-base-100/60 p-12 text-center"
      >
        <h2 class="text-xl font-semibold">ยังไม่ได้เลือกรูป</h2>
        <p class="mt-2 text-base-content/70">
          ไปที่หน้าสแกนรูป เลือกรูปจาก Source แล้วกดย้ายรูป
        </p>
        <NuxtLink to="/scan" class="btn btn-primary mt-6 min-h-12 px-6">
          ไปหน้าสแกนรูป
        </NuxtLink>
      </div>

      <div v-else class="space-y-5">
        <MoveProgressPanel
          v-if="isMoving || result"
          :total="checkedPaths.length"
          :progress="progress"
          :result="result"
        />

        <div
          class="flex flex-col gap-4 rounded-[1.5rem] border border-base-300 bg-base-100 p-5 sm:flex-row sm:items-center sm:justify-between"
        >
          <div>
            <p class="text-lg font-bold">
              เลือกไว้ {{ checkedPaths.length }} จาก
              {{ photos.length }} รูป
            </p>
            <p class="text-sm text-base-content/70">
              ตรวจรายการอีกครั้งก่อนเริ่มย้าย
            </p>
          </div>

          <div class="flex flex-wrap gap-2">
            <template v-if="!result">
              <button
                type="button"
                class="btn btn-outline min-h-12 border-primary/35 text-primary"
                :disabled="isMoving"
                @click="toggleAll"
              >
                {{ allChecked ? "ยกเลิกทั้งหมด" : "เลือกทั้งหมด" }}
              </button>

              <button
                type="button"
                class="btn btn-primary min-h-12 px-6"
                :disabled="!canMove"
                @click="move"
              >
                <span v-if="isMoving" class="loading loading-spinner loading-sm" />
                {{ isMoving ? "กำลังย้าย..." : `ย้าย ${checkedPaths.length} รูป` }}
              </button>
            </template>

            <NuxtLink v-else to="/scan" class="btn btn-primary min-h-12 px-6">
              กลับไปสแกนรูป
            </NuxtLink>
          </div>
        </div>

        <MovePhotoGrid
          v-model:checked-paths="checkedPaths"
          :photos="photos"
          :preview-urls="previews"
          :disabled="isMoving || result !== null"
          :results="statusByPath"
          :status-of="statusOf"
        />
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import MovePhotoGrid from "../components/move/MovePhotoGrid.vue";
import MoveProgressPanel from "../components/move/MoveProgressPanel.vue";
import { usePhotoMover } from "../composables/usePhotoMover";

const {
  selectedPhotos,
  selectedPreviews,
  checkedPaths,
  statusByPath,
  progress,
  result,
  isMoving,
  errorMessage,
  canMove,
  allChecked,
  statusOf,
  toggleAll,
  move,
} = usePhotoMover();

// Keep showing the moved photos after the shared selection is cleared.
const photos = ref([...selectedPhotos.value]);
const previews = ref({ ...selectedPreviews.value });
</script>
