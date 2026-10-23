<template>
  <main class="gallery-page mx-auto max-w-7xl px-5 py-8 sm:px-8">
    <header class="mb-6">
      <div class="mb-5 flex items-center justify-between gap-4">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.16em] text-primary/80">Photo Backup</p>
          <h1 class="mt-1 text-3xl font-bold tracking-tight">Gallery</h1>
        </div>
        <button
          type="button"
          class="btn btn-outline refresh-button"
          aria-label="รีเฟรชรูปภาพ"
          :disabled="loading || deletingPhotoID !== null || deletingAllPhotos || batchEmbedding"
          @click="refreshPhoto"
        >
          <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="h-5 w-5">
            <path d="M20 7v5h-5M4 17v-5h5m10-1a7 7 0 0 0-12.3-4.6L4 12m16 0-2.7 5.6A7 7 0 0 1 5 17" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <span class="hidden sm:inline">{{ loading ? "กำลังโหลด..." : "รีเฟรช" }}</span>
        </button>
      </div>

      <form class="search-shell" role="search" @submit.prevent="searchGallery">
        <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="search-icon">
          <circle cx="10.8" cy="10.8" r="6.8" stroke="currentColor" stroke-width="2" />
          <path d="m16 16 4.2 4.2" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
        </svg>
        <label class="sr-only" for="gallery-search">ค้นหารูปภาพ</label>
        <input
          id="gallery-search"
          v-model="searchText"
          type="search"
          autocomplete="off"
          placeholder="ค้นหารูปภาพ เช่น เรือกลางทะเล หรือ หมา"
          :disabled="batchEmbedding"
        />
        <button v-if="searchText" type="button" class="clear-search" aria-label="ล้างคำค้นหา" @click="clearSearch">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" class="h-5 w-5">
            <path d="m6 6 8 8m0-8-8 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
          </svg>
        </button>
        <button class="btn btn-primary search-submit" type="submit" :disabled="searching || batchEmbedding">
          {{ searching ? "กำลังค้นหา..." : "ค้นหา" }}
        </button>
      </form>
      <p v-if="searchError" role="alert" class="mt-2 text-sm text-error">{{ searchError }}</p>

      <nav class="filter-row mt-4" aria-label="ตัวกรองรูปภาพ">
        <button
          v-for="filter in filters"
          :key="filter.id"
          type="button"
          class="filter-chip"
          :class="{ 'filter-chip-active': activeFilter === filter.id }"
          :aria-pressed="activeFilter === filter.id"
          @click="activeFilter = filter.id"
        >
          {{ filter.label }} <span class="filter-count">{{ filter.count }}</span>
        </button>
      </nav>
    </header>

    <section v-if="photos.length > 0 && !loading && !errMessage" class="selection-toolbar mb-5" aria-label="การเลือกและสร้าง Embedding">
      <p class="text-sm text-base-content/75" aria-live="polite">
        <strong>{{ displayedPhotos.length }}</strong> รูป
        <span v-if="selectedPhotoIDs.length"> · เลือกแล้ว {{ selectedPhotoIDs.length }} รูป</span>
        <span v-if="batchEmbedding"> · กำลังทำ {{ batchProgress.completed }}/{{ batchProgress.total }}</span>
      </p>
      <div class="flex flex-wrap items-center gap-2">
        <button type="button" class="btn btn-sm btn-outline" :disabled="selectablePhotos.length === 0 || deletingAllPhotos || batchEmbedding || embeddingPhotoIDs.length > 0" @click="toggleSelectAll">
          {{ allSelectableSelected ? "ยกเลิกเลือก" : "เลือกทั้งหมด" }}
        </button>
        <button v-if="selectedPhotoIDs.length" type="button" class="btn btn-sm btn-ghost" :disabled="deletingAllPhotos || batchEmbedding || embeddingPhotoIDs.length > 0" @click="selectedPhotoIDs = []">
          ล้างที่เลือก
        </button>
        <button type="button" class="btn btn-sm btn-secondary" :disabled="selectedPhotoIDs.length === 0 || deletingAllPhotos || batchEmbedding || embeddingPhotoIDs.length > 0" @click="generateSelectedEmbeddings">
          {{ batchEmbedding ? "กำลังสร้าง..." : `สร้าง Embedding (${selectedPhotoIDs.length})` }}
        </button>
        <button type="button" class="btn btn-sm btn-outline btn-error" :disabled="loading || deletingAllPhotos || deletingPhotoID !== null || batchEmbedding || photos.length === 0" @click="openDeleteAllDialog">
          {{ deletingAllPhotos ? `กำลังลบ ${deleteAllProgress.completed}/${deleteAllProgress.total}...` : "ลบทั้งหมด" }}
        </button>
      </div>
    </section>

    <section v-if="batchEmbedding" class="mb-5 rounded-xl border border-base-300 bg-base-100 p-4" aria-label="ความคืบหน้าการสร้าง Embedding" role="status">
      <div class="mb-2 flex flex-wrap items-center justify-between gap-2 text-sm">
        <strong>กำลังสร้าง Embedding {{ batchProgress.completed }} / {{ batchProgress.total }} รูป</strong>
        <span>{{ batchPercent }}%</span>
      </div>
      <progress class="progress progress-primary w-full" :value="batchProgress.completed" :max="batchProgress.total" :aria-label="`ทำแล้ว ${batchProgress.completed} จาก ${batchProgress.total} รูป`"></progress>
      <p class="mt-2 text-sm text-base-content/75">
        ใช้เวลา {{ formatDuration(batchProgress.elapsedMs) }}
        <span v-if="batchEtaMs !== null"> · เหลือประมาณ {{ formatDuration(batchEtaMs) }}</span>
        <span v-else> · กำลังคำนวณเวลาที่เหลือ</span>
        <span v-if="batchProgress.failed > 0" class="text-error"> · ล้มเหลว {{ batchProgress.failed }} รูป</span>
      </p>
    </section>
    <p v-else-if="batchResult" class="mb-5 rounded-xl border border-base-300 bg-base-100 px-4 py-3 text-sm" role="status" aria-live="polite">
      {{ batchResult }}
    </p>

    <div v-if="loading" class="gallery-message" aria-live="polite">
      <span class="loading loading-spinner loading-md text-primary" aria-hidden="true"></span>
      <span>กำลังโหลดรูปภาพ...</span>
    </div>
    <div v-else-if="errMessage" role="alert" class="gallery-message error-message">
      <p>{{ errMessage }}</p>
      <button type="button" class="btn btn-sm btn-outline" @click="refreshPhoto">ลองอีกครั้ง</button>
    </div>
    <div v-else-if="photos.length === 0" class="gallery-message">
      ยังไม่มีรูปที่ย้ายมาใน Destination ปัจจุบัน
    </div>
    <div v-else-if="displayedPhotos.length === 0" class="gallery-message">
      <p class="font-semibold">{{ hasSearched && searchText.trim() ? `ไม่พบรูปที่ตรงกับ “${searchText}”` : "ไม่มีรูปในหมวดนี้" }}</p>
      <button v-if="activeFilter !== 'all' || searchText" type="button" class="btn btn-sm btn-outline" @click="resetGalleryView">แสดงรูปทั้งหมด</button>
    </div>

    <section v-else class="thumbnail-grid" aria-label="ผลลัพธ์รูปภาพ">
      <article v-for="photo in displayedPhotos" :key="photo.id" class="photo-tile" :class="{ 'photo-tile-selected': selectedPhotoIDs.includes(photo.id) }">
        <div class="thumbnail-frame">
          <img v-if="previewUrls[photo.id]" :src="previewUrls[photo.id]" :alt="photo.filename" loading="lazy" decoding="async" />
          <div v-else class="thumbnail-placeholder">
            {{ photo.status === "missing" ? "ไม่พบไฟล์รูปภาพ" : "กำลังโหลดตัวอย่าง..." }}
          </div>
          <label v-if="photo.status === 'active' && photo.analysisStatus !== 'completed'" class="photo-select" :aria-label="`เลือกรูป ${photo.filename}`" :title="`เลือกรูป ${photo.filename}`">
            <input v-model="selectedPhotoIDs" type="checkbox" class="checkbox checkbox-primary" :value="photo.id" :disabled="deletingAllPhotos || batchEmbedding || embeddingPhotoIDs.length > 0" />
          </label>
          <span class="photo-status" :class="photo.analysisStatus === 'completed' ? 'status-complete' : 'status-pending'">
            {{ photo.analysisStatus === "completed" ? "พร้อมค้นหา" : "รอ Embedding" }}
          </span>
        </div>
        <div class="photo-tile-body">
          <h2 :title="photo.filename">{{ photo.filename }}</h2>
          <p class="file-state" :class="{ 'file-missing': photo.status !== 'active' }">
            {{ photo.status === "active" ? "ไฟล์พร้อมใช้งาน" : "ไม่พบไฟล์ในโฟลเดอร์" }}
          </p>
          <p v-if="embeddingErrors[photo.id]" role="alert" class="mt-2 text-xs text-error">{{ embeddingErrors[photo.id] }}</p>
          <div class="photo-actions">
            <button type="button" class="btn btn-xs btn-secondary" :disabled="photo.status !== 'active' || photo.analysisStatus === 'completed' || embeddingPhotoIDs.length > 0 || deletingAllPhotos || batchEmbedding" @click="generateEmbedding(photo)">
              {{ embeddingPhotoIDs.includes(photo.id) ? "กำลังทำ..." : photo.analysisStatus === "completed" ? "ทำแล้ว" : "สร้าง Embedding" }}
            </button>
            <button type="button" class="btn btn-xs btn-ghost text-error" :disabled="loading || deletingPhotoID !== null || deletingAllPhotos || batchEmbedding" @click="openDeleteDialog(photo)">
              {{ deletingPhotoID === photo.id ? "กำลังลบ..." : "ลบ" }}
            </button>
          </div>
        </div>
      </article>
    </section>

    <dialog ref="deleteDialog" class="modal" aria-labelledby="delete-dialog-title" @cancel="deletingPhotoID !== null && $event.preventDefault()">
      <div class="modal-box bg-base-100">
        <h2 id="delete-dialog-title" class="text-xl font-bold">ยืนยันการลบรูปภาพ</h2>
        <p class="mt-4 break-all">ต้องการลบ “{{ pendingDelete?.filename }}” หรือไม่?</p>
        <p class="mt-2 text-sm text-base-content/70">ไฟล์จริงและข้อมูลในฐานข้อมูลจะถูกลบ การลบนี้ไม่สามารถย้อนกลับผ่านแอปได้</p>
        <p v-if="errMessage" role="alert" class="mt-4 text-error">{{ errMessage }}</p>
        <div class="modal-action">
          <button type="button" class="btn" :disabled="deletingPhotoID !== null" @click="deleteDialog?.close()">ยกเลิก</button>
          <button type="button" class="btn btn-error" :disabled="deletingPhotoID !== null || !pendingDelete" @click="pendingDelete && deletePhoto(pendingDelete)">
            {{ deletingPhotoID !== null ? "กำลังลบ..." : "ยืนยันลบ" }}
          </button>
        </div>
      </div>
    </dialog>

    <dialog ref="deleteAllDialog" class="modal" aria-labelledby="delete-all-dialog-title" @cancel="deletingAllPhotos && $event.preventDefault()">
      <div class="modal-box bg-base-100">
        <h2 id="delete-all-dialog-title" class="text-xl font-bold text-error">ยืนยันการลบรูปทั้งหมด</h2>
        <p class="mt-4">
          ต้องการลบรูปทั้งหมด <strong>{{ deleteAllProgress.total || photos.length }} รูป</strong> จาก Gallery ใช่ไหม?
        </p>
        <p class="mt-2 text-sm text-base-content/75">
          การทำงานนี้ไม่ขึ้นกับคำค้นหาหรือตัวกรอง และจะลบทั้งไฟล์จริงกับข้อมูลในฐานข้อมูล ไม่สามารถย้อนกลับผ่านแอปได้
        </p>
        <p v-if="deletingAllPhotos" class="mt-4 text-sm" aria-live="polite">
          ดำเนินการแล้ว {{ deleteAllProgress.completed }}/{{ deleteAllProgress.total }} รูป
        </p>
        <p v-if="deleteAllError" role="alert" class="mt-4 whitespace-pre-line text-sm text-error">{{ deleteAllError }}</p>
        <div class="modal-action">
          <button type="button" class="btn" :disabled="deletingAllPhotos" @click="deleteAllDialog?.close()">{{ deleteAllError ? "ปิด" : "ยกเลิก" }}</button>
          <button type="button" class="btn btn-error" :disabled="deletingAllPhotos || photos.length === 0" @click="deleteAllPhotos">
            {{ deletingAllPhotos ? "กำลังลบ..." : deleteAllError ? "ลองลบที่เหลืออีกครั้ง" : `ยืนยันลบ ${deleteAllProgress.total || photos.length} รูป` }}
          </button>
        </div>
      </div>
    </dialog>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  GetGalleryPhotos,
  GetGalleryPhotoPreview,
  GenerateGalleryPhotoEmbedding,
  DeleteGalleryPhoto,
  SearchGalleryPhotos,
} from "../../wailsjs/go/controller/GalleryController";
import type { dto } from "../../wailsjs/go/models";

const photos = ref<dto.GalleryPhoto[]>([]);
const loading = ref(true);
const errMessage = ref("");
const previewUrls = ref<Record<number, string>>({});
const deletingPhotoID = ref<number | null>(null);
const deleteDialog = ref<HTMLDialogElement | null>(null);
const deleteAllDialog = ref<HTMLDialogElement | null>(null);
const pendingDelete = ref<dto.GalleryPhoto | null>(null);
const deletingAllPhotos = ref(false);
const deleteAllProgress = ref({ completed: 0, total: 0 });
const deleteAllError = ref("");
const embeddingPhotoIDs = ref<number[]>([]);
const embeddingErrors = ref<Record<number, string>>({});

const searchText = ref("");
const searchResults = ref<dto.PhotoSearchResult[]>([]);
const hasSearched = ref(false);
const searching = ref(false);
const searchError = ref("");
const selectedPhotoIDs = ref<number[]>([]);
const batchEmbedding = ref(false);
const batchProgress = ref({ completed: 0, total: 0, failed: 0, elapsedMs: 0 });
const batchResult = ref("");
const activeFilter = ref<"all" | "pending" | "completed">("all");

const batchPercent = computed(() =>
  batchProgress.value.total > 0
    ? Math.floor((batchProgress.value.completed / batchProgress.value.total) * 100)
    : 0,
);
const batchEtaMs = computed(() => {
  const { completed, total, elapsedMs } = batchProgress.value;
  if (completed === 0 || total === 0) return null;
  return Math.max(0, (elapsedMs / completed) * (total - completed));
});

function formatDuration(milliseconds: number): string {
  const totalSeconds = Math.floor(milliseconds / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (hours > 0) return `${hours} ชม. ${minutes} นาที`;
  if (minutes > 0) return seconds > 0 ? `${minutes} นาที ${seconds} วินาที` : `${minutes} นาที`;
  return `${seconds} วินาที`;
}

async function searchGallery() {
  const text = searchText.value.trim();
  if (!text) {
    hasSearched.value = false;
    searchResults.value = [];
    searchError.value = "";
    activeFilter.value = "all";
    return;
  }

  searching.value = true;
  searchError.value = "";

  try {
    searchResults.value = await SearchGalleryPhotos(text);
    hasSearched.value = true;
  } catch (err) {
    searchError.value = err instanceof Error ? err.message : String(err);
  } finally {
    searching.value = false;
  }
}

function clearSearch() {
  searchText.value = "";
  hasSearched.value = false;
  searchResults.value = [];
  searchError.value = "";
  activeFilter.value = "all";
}

function resetGalleryView() {
  activeFilter.value = "all";
  clearSearch();
}

const searchMatchedPhotos = computed(() => {
  if (!hasSearched.value || !searchText.value.trim()) return photos.value;

  const photosByID = new Map(photos.value.map((photo) => [photo.id, photo]));
  return searchResults.value
    .map((result) => photosByID.get(result.id))
    .filter((photo): photo is dto.GalleryPhoto => photo !== undefined);
});

const displayedPhotos = computed(() => {
  if (activeFilter.value === "pending") {
    return searchMatchedPhotos.value.filter(
      (photo) => photo.status === "active" && photo.analysisStatus !== "completed",
    );
  }
  if (activeFilter.value === "completed") {
    return searchMatchedPhotos.value.filter((photo) => photo.analysisStatus === "completed");
  }
  return searchMatchedPhotos.value;
});

const filters = computed(() => [
  { id: "all" as const, label: "ทั้งหมด", count: searchMatchedPhotos.value.length },
  {
    id: "pending" as const,
    label: "รอสร้าง Embedding",
    count: searchMatchedPhotos.value.filter(
      (photo) => photo.status === "active" && photo.analysisStatus !== "completed",
    ).length,
  },
  {
    id: "completed" as const,
    label: "สร้างแล้ว",
    count: searchMatchedPhotos.value.filter((photo) => photo.analysisStatus === "completed").length,
  },
]);

const selectablePhotos = computed(() =>
  displayedPhotos.value.filter(
    (photo) => photo.status === "active" && photo.analysisStatus !== "completed",
  ),
);

const allSelectableSelected = computed(
  () =>
    selectablePhotos.value.length > 0 &&
    selectablePhotos.value.every((photo) => selectedPhotoIDs.value.includes(photo.id)),
);

function toggleSelectAll() {
  const visibleIDs = new Set(selectablePhotos.value.map((photo) => photo.id));

  if (allSelectableSelected.value) {
    selectedPhotoIDs.value = selectedPhotoIDs.value.filter((id) => !visibleIDs.has(id));
    return;
  }

  selectedPhotoIDs.value = [
    ...new Set([...selectedPhotoIDs.value, ...visibleIDs]),
  ];
}

async function loadPhotos() {
  loading.value = true;
  errMessage.value = "";

  try {
    photos.value = await GetGalleryPhotos();

    previewUrls.value = {};
    await loadPreviews(photos.value);
  } catch (err) {
    errMessage.value = err instanceof Error ? err.message : String(err);
  } finally {
    loading.value = false;
  }
}

async function loadPreviews(items: dto.GalleryPhoto[]) {
  for (const photo of items) {
    if (photo.status !== "active") continue;

    try {
      previewUrls.value[photo.id] = await GetGalleryPhotoPreview(photo.id);
    } catch (err) {
      previewUrls.value[photo.id] = "";

      const message = err instanceof Error ? err.message : String(err);
      errMessage.value = `${photo.filename}: ${message}`;
      console.error("โหลด Preview ไม่สำเร็จ", photo.id, err);
    }
  }
}

function openDeleteDialog(photo: dto.GalleryPhoto) {
  pendingDelete.value = photo;
  errMessage.value = "";
  deleteDialog.value?.showModal();
}

function openDeleteAllDialog() {
  if (photos.value.length === 0 || loading.value || deletingAllPhotos.value) return;

  deleteAllError.value = "";
  deleteAllProgress.value = { completed: 0, total: photos.value.length };
  deleteAllDialog.value?.showModal();
}

async function deleteAllPhotos() {
  if (
    deletingAllPhotos.value ||
    deletingPhotoID.value !== null ||
    embeddingPhotoIDs.value.length > 0 ||
    batchEmbedding.value
  ) return;

  const targets = [...photos.value];
  if (targets.length === 0) return;

  deletingAllPhotos.value = true;
  deleteAllError.value = "";
  deleteAllProgress.value = { completed: 0, total: targets.length };
  const failures: string[] = [];

  try {
    for (const photo of targets) {
      try {
        await DeleteGalleryPhoto(photo.id);
        photos.value = photos.value.filter((item) => item.id !== photo.id);
        delete previewUrls.value[photo.id];
        delete embeddingErrors.value[photo.id];
        selectedPhotoIDs.value = selectedPhotoIDs.value.filter((id) => id !== photo.id);
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        failures.push(`${photo.filename}: ${message}`);
      } finally {
        deleteAllProgress.value.completed += 1;
      }
    }

    if (failures.length > 0) {
      deleteAllError.value = `ลบสำเร็จ ${targets.length - failures.length}/${targets.length} รูป; ลบไม่สำเร็จ ${failures.length} รูป:\n${failures.join("\n")}`;
      deleteAllProgress.value = { completed: 0, total: photos.value.length };
    } else {
      deleteAllDialog.value?.close();
      selectedPhotoIDs.value = [];
      deleteAllProgress.value = { completed: 0, total: 0 };
    }
  } finally {
    deletingAllPhotos.value = false;
  }
}

async function generateEmbedding(photo: dto.GalleryPhoto) {
  if (
    photo.status !== "active" ||
    photo.analysisStatus === "completed" ||
    embeddingPhotoIDs.value.length > 0 ||
    batchEmbedding.value
  ) {
    return;
  }

  await requestPhotoEmbedding(photo);
}

async function requestPhotoEmbedding(photo: dto.GalleryPhoto): Promise<boolean> {
  embeddingPhotoIDs.value = [...embeddingPhotoIDs.value, photo.id];
  delete embeddingErrors.value[photo.id];

  try {
    await GenerateGalleryPhotoEmbedding(photo.id);
    photo.analysisStatus = "completed";
    selectedPhotoIDs.value = selectedPhotoIDs.value.filter((id) => id !== photo.id);
    return true;
  } catch (err) {
    embeddingErrors.value[photo.id] =
      err instanceof Error ? err.message : String(err);
    return false;
  } finally {
    embeddingPhotoIDs.value = embeddingPhotoIDs.value.filter((id) => id !== photo.id);
  }
}

async function generateSelectedEmbeddings() {
  if (batchEmbedding.value || embeddingPhotoIDs.value.length > 0) return;

  const selectedIDs = new Set(selectedPhotoIDs.value);
  const targets = photos.value.filter(
    (photo) =>
      selectedIDs.has(photo.id) &&
      photo.status === "active" &&
      photo.analysisStatus !== "completed",
  );

  if (targets.length === 0) return;

  batchEmbedding.value = true;
  batchResult.value = "";
  batchProgress.value = { completed: 0, total: targets.length, failed: 0, elapsedMs: 0 };
  const startedAt = Date.now();
  const timer = setInterval(() => {
    batchProgress.value.elapsedMs = Date.now() - startedAt;
  }, 1000);

  try {
    let nextIndex = 0;

    async function worker() {
      while (true) {
        const photo = targets[nextIndex];
        nextIndex += 1;

        if (!photo) return;

        const succeeded = await requestPhotoEmbedding(photo);
        if (!succeeded) batchProgress.value.failed += 1;
        batchProgress.value.completed += 1;
      }
    }

    const workerCount = Math.min(2, targets.length);
    await Promise.all(Array.from({ length: workerCount }, () => worker()));
  } finally {
    clearInterval(timer);
    batchProgress.value.elapsedMs = Date.now() - startedAt;
    batchEmbedding.value = false;
    const failed = batchProgress.value.failed;
    const processed = batchProgress.value.completed;
    batchResult.value = `ประมวลผลเสร็จ ${processed}/${targets.length} รูป ใช้เวลา ${formatDuration(batchProgress.value.elapsedMs)}${failed > 0 ? ` · ไม่สำเร็จ ${failed} รูป (ดูข้อความผิดพลาดในการ์ดรูป)` : ""}`;
  }
}

async function deletePhoto(photo: dto.GalleryPhoto) {
  if (loading.value || deletingPhotoID.value !== null) {
    return;
  }

  deletingPhotoID.value = photo.id;
  errMessage.value = "";

  try {
    await DeleteGalleryPhoto(photo.id);

    photos.value = photos.value.filter((item) => item.id != photo.id);

    delete previewUrls.value[photo.id];

    deleteDialog.value?.close();
    pendingDelete.value = null;
  } catch (err) {
    errMessage.value = err instanceof Error ? err.message : String(err);
  } finally {
    deletingPhotoID.value = null;
  }
}

function refreshPhoto() {
  loadPhotos();
}

onMounted(loadPhotos);
</script>

<style scoped>
.search-shell {
  display: flex;
  min-height: 54px;
  align-items: center;
  gap: 12px;
  border: 1px solid color-mix(in srgb, var(--color-base-content) 18%, transparent);
  border-radius: 16px;
  background: var(--color-base-100);
  padding: 5px 6px 5px 16px;
  box-shadow: 0 2px 10px rgb(55 39 27 / 5%);
  transition: border-color 160ms ease, box-shadow 160ms ease;
}

.search-shell:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 22%, transparent);
}

.search-icon {
  width: 22px;
  height: 22px;
  flex: 0 0 auto;
  color: var(--color-primary);
}

.search-shell input {
  min-width: 0;
  flex: 1;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--color-base-content);
  font: inherit;
}

.search-shell input::placeholder {
  color: color-mix(in srgb, var(--color-base-content) 58%, transparent);
}

.search-submit {
  min-height: 42px;
  border-radius: 12px;
}

.clear-search {
  display: grid;
  width: 40px;
  height: 40px;
  flex: 0 0 auto;
  place-items: center;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: color-mix(in srgb, var(--color-base-content) 65%, transparent);
  cursor: pointer;
}

.clear-search:hover {
  background: color-mix(in srgb, var(--color-base-content) 8%, transparent);
}

.filter-row {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 0 5px;
  scrollbar-width: thin;
}

.filter-chip {
  display: inline-flex;
  min-height: 40px;
  flex: 0 0 auto;
  align-items: center;
  gap: 8px;
  padding: 0 14px;
  border: 1px solid color-mix(in srgb, var(--color-base-content) 16%, transparent);
  border-radius: 999px;
  background: var(--color-base-100);
  color: color-mix(in srgb, var(--color-base-content) 88%, transparent);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 160ms ease, color 160ms ease, border-color 160ms ease;
}

.filter-chip:hover {
  border-color: color-mix(in srgb, var(--color-primary) 62%, transparent);
}

.filter-chip-active {
  border-color: var(--color-primary);
  background: var(--color-primary);
  color: var(--color-primary-content);
}

.filter-count {
  font-size: 0.75rem;
  opacity: 0.78;
}

.selection-toolbar {
  display: flex;
  min-height: 56px;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid color-mix(in srgb, var(--color-base-content) 14%, transparent);
  padding: 0 2px 12px;
}

.thumbnail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 190px), 1fr));
  gap: 18px;
}

.photo-tile {
  min-width: 0;
  overflow: hidden;
  border: 1px solid color-mix(in srgb, var(--color-base-content) 17%, transparent);
  border-radius: 14px;
  background: var(--color-base-100);
  transition: border-color 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.photo-tile:hover {
  transform: translateY(-2px);
  border-color: color-mix(in srgb, var(--color-primary) 62%, transparent);
  box-shadow: 0 8px 22px rgb(55 39 27 / 9%);
}

.photo-tile-selected {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-primary) 28%, transparent);
}

.thumbnail-frame {
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
  background: var(--color-base-200);
}

.thumbnail-frame img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  display: grid;
  width: 100%;
  height: 100%;
  place-items: center;
  padding: 16px;
  color: color-mix(in srgb, var(--color-base-content) 68%, transparent);
  font-size: 0.8125rem;
  text-align: center;
}

.photo-select {
  position: absolute;
  top: 9px;
  left: 9px;
  display: grid;
  min-width: 42px;
  min-height: 42px;
  place-items: center;
  border: 1px solid rgb(255 255 255 / 75%);
  border-radius: 10px;
  background: rgb(255 255 255 / 90%);
  cursor: pointer;
}

.photo-status {
  position: absolute;
  right: 8px;
  bottom: 8px;
  max-width: calc(100% - 16px);
  overflow: hidden;
  border-radius: 999px;
  padding: 4px 9px;
  background: rgb(255 255 255 / 94%);
  color: #43372f;
  font-size: 0.7rem;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-complete {
  color: #276749;
}

.status-pending {
  color: #81540e;
}

.photo-tile-body {
  padding: 10px 11px 11px;
}

.photo-tile-body h2 {
  overflow: hidden;
  color: var(--color-base-content);
  font-size: 0.875rem;
  font-weight: 700;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-state {
  margin-top: 5px;
  color: color-mix(in srgb, var(--color-base-content) 72%, transparent);
  font-size: 0.75rem;
}

.file-missing {
  color: var(--color-error);
}

.photo-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 5px;
  margin-top: 9px;
}

.gallery-message {
  display: flex;
  min-height: 160px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  border: 1px dashed color-mix(in srgb, var(--color-base-content) 24%, transparent);
  border-radius: 16px;
  padding: 24px;
  color: color-mix(in srgb, var(--color-base-content) 82%, transparent);
  text-align: center;
}

.error-message {
  border-color: color-mix(in srgb, var(--color-error) 48%, transparent);
  color: var(--color-error);
}

@media (max-width: 480px) {
  .gallery-page {
    padding-top: 20px;
  }

  .search-shell {
    gap: 7px;
    padding-left: 11px;
  }

  .search-shell input {
    font-size: 0.9rem;
  }

  .search-submit {
    padding-inline: 12px;
  }

  .thumbnail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
  }

  .photo-tile-body {
    padding: 8px;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    scroll-behavior: auto !important;
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
  }
}
</style>
