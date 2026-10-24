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
          :disabled="loading || deletingPhotoID !== null || deletingAllPhotos"
          @click="refreshPhoto"
        >
          <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="h-5 w-5">
            <path d="M20 7v5h-5M4 17v-5h5m10-1a7 7 0 0 0-12.3-4.6L4 12m16 0-2.7 5.6A7 7 0 0 1 5 17" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <span class="hidden sm:inline">{{ loading ? "กำลังโหลด..." : "รีเฟรช" }}</span>
        </button>
      </div>

      <form class="search-shell" role="search" @submit.prevent>
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
        />
        <button v-if="searchText" type="button" class="clear-search" aria-label="ล้างคำค้นหา" @click="clearSearch">
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" class="h-5 w-5">
            <path d="m6 6 8 8m0-8-8 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
          </svg>
        </button>
      </form>
      <p v-if="searching || searchPending" class="mt-2 text-sm text-base-content/65" role="status" aria-live="polite">
        {{ searching ? "กำลังค้นหารูปภาพ..." : "กำลังรอค้นหา..." }}
      </p>
      <p v-if="searchError" role="alert" class="mt-2 text-sm text-error">{{ searchError }}</p>

    </header>

    <section v-if="photos.length > 0 && !loading && !errMessage" class="selection-toolbar mb-5" aria-label="สถานะรูปภาพ">
      <p class="text-sm text-base-content/75" aria-live="polite">
        <strong>{{ readyPhotos.length }}</strong> รูปพร้อมค้นหา
        <span v-if="pendingPhotoCount"> · รอ {{ pendingPhotoCount }} รูป</span>
        <span v-if="processingPhotoCount"> · กำลังสร้าง {{ processingPhotoCount }} รูป</span>
        <span v-if="failedPhotoCount" class="text-error"> · ล้มเหลว {{ failedPhotoCount }} รูป</span>
      </p>
      <div class="flex flex-wrap items-center gap-2">
        <button type="button" class="btn btn-sm btn-outline btn-error" :disabled="loading || deletingAllPhotos || deletingPhotoID !== null || processingPhotoCount > 0 || photos.length === 0" @click="openDeleteAllDialog">
          {{ deletingAllPhotos ? `กำลังลบ ${deleteAllProgress.completed}/${deleteAllProgress.total}...` : "ลบทั้งหมด" }}
        </button>
      </div>
    </section>
    <section v-if="failedPhotos.length > 0 && !loading" class="mb-5 rounded-xl border border-error/30 bg-base-100 p-4" aria-label="Embedding ที่ล้มเหลว">
      <h2 class="font-semibold text-error">มีรูปที่สร้าง Embedding ไม่สำเร็จ {{ failedPhotos.length }} รูป</h2>
      <ul class="mt-2 space-y-1 text-sm">
        <li v-for="photo in failedPhotos" :key="photo.id" class="break-words">
          <strong>{{ photo.filename }}</strong><span v-if="photo.analysisError"> — {{ photo.analysisError }}</span>
        </li>
      </ul>
    </section>

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
      <p class="font-semibold">{{ hasSearched && searchText.trim() ? `ไม่พบรูปที่ตรงกับ “${searchText}”` : "ยังไม่มีรูปที่สร้าง Embedding เสร็จ" }}</p>
      <p v-if="pendingPhotoCount || processingPhotoCount">รูปจะแสดงใน Gallery เมื่อสร้าง Embedding เสร็จ</p>
      <p v-if="failedPhotoCount" class="text-error">มี {{ failedPhotoCount }} รูปที่สร้าง Embedding ไม่สำเร็จ กรุณารีเฟรชเพื่อตรวจสอบสถานะ</p>
      <button v-if="searchText" type="button" class="btn btn-sm btn-outline" @click="clearSearch">ล้างคำค้นหา</button>
    </div>

    <section v-else class="thumbnail-grid" aria-label="ผลลัพธ์รูปภาพ">
      <article v-for="photo in displayedPhotos" :key="photo.id" class="photo-tile">
        <div class="thumbnail-frame">
          <img v-if="previewUrls[photo.id]" :src="previewUrls[photo.id]" :alt="photo.filename" loading="lazy" decoding="async" />
          <div v-else class="thumbnail-placeholder">
            {{ photo.status === "missing" ? "ไม่พบไฟล์รูปภาพ" : "กำลังโหลดตัวอย่าง..." }}
          </div>
          <!-- <span class="photo-status status-complete">พร้อมค้นหา</span> -->
        </div>
        <div class="photo-tile-body">
          <h2 :title="photo.filename">{{ photo.filename }}</h2>
          <p class="file-state" :class="{ 'file-missing': photo.status !== 'active' }">
            {{ photo.status === "active" ? "ไฟล์พร้อมใช้งาน" : "ไม่พบไฟล์ในโฟลเดอร์" }}
          </p>
          <div class="photo-actions">
            <button type="button" class="btn btn-xs btn-ghost text-error" :disabled="loading || deletingPhotoID !== null || deletingAllPhotos" @click="openDeleteDialog(photo)">
              {{ deletingPhotoID === photo.id ? "กำลังลบ..." : "ลบ" }}
            </button>
          </div>
        </div>
      </article>
    </section>

    <aside
      v-if="embeddingProgress.isRunning && embeddingProgress.total > 0"
      class="fixed bottom-3 right-3 z-40 w-[min(25rem,calc(100%-1.5rem))] rounded-lg border border-neutral-content/10 bg-neutral px-3 py-2 text-neutral-content shadow-lg sm:bottom-4 sm:right-4"
      role="status"
      aria-live="polite"
      aria-label="ความคืบหน้าการสร้าง Embedding"
    >
      <div class="flex items-center gap-2.5">
        <div class="min-w-0 flex-1">
          <div class="mb-1 flex items-center justify-between gap-2 text-xs">
            <span class="truncate font-semibold">Embedding · {{ embeddingProgress.finished }}/{{ embeddingProgress.total }}</span>
            <span class="shrink-0 text-[11px] text-neutral-content/65">{{ embeddingProgress.percent }}%</span>
          </div>
          <progress
            class="progress progress-success h-1.5 w-full bg-neutral-content/15"
            :value="embeddingProgress.finished"
            :max="embeddingProgress.total"
            :aria-label="`เสร็จแล้ว ${embeddingProgress.finished} จาก ${embeddingProgress.total} รูป`"
          ></progress>
          <p class="mt-1 text-[10px] leading-tight text-neutral-content/70">
            ค้นหาได้ {{ readyPhotos.length }} · รอ {{ pendingPhotoCount }}<span v-if="processingPhotoCount"> · ทำ {{ processingPhotoCount }}</span><span v-if="failedPhotoCount"> · พลาด {{ failedPhotoCount }}</span>
          </p>
        </div>
        <button
          type="button"
          class="btn btn-xs h-7 min-h-7 shrink-0 border-neutral-content/20 bg-transparent px-2 text-[11px] text-neutral-content hover:bg-neutral-content/10"
          :disabled="queueActionPending"
          @click="toggleEmbeddingQueue"
        >
          {{ queueActionPending ? "กำลังอัปเดต..." : queuePaused ? "ทำต่อ" : "พัก" }}
        </button>
      </div>
    </aside>

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
          การลบนี้รวมรูปที่ยังสร้าง Embedding ไม่เสร็จด้วย และจะลบทั้งไฟล์จริงกับข้อมูลในฐานข้อมูล ไม่สามารถย้อนกลับผ่านแอปได้
        </p>
        <p v-if="deletingAllPhotos" class="mt-4 text-sm" aria-live="polite">
          ดำเนินการแล้ว {{ deleteAllProgress.completed }}/{{ deleteAllProgress.total }} รูป
        </p>
        <p v-if="deleteAllError" role="alert" class="mt-4 whitespace-pre-line text-sm text-error">{{ deleteAllError }}</p>
        <div class="modal-action">
          <button type="button" class="btn" :disabled="deletingAllPhotos" @click="deleteAllDialog?.close()">{{ deleteAllError ? "ปิด" : "ยกเลิก" }}</button>
          <button type="button" class="btn btn-error" :disabled="deletingAllPhotos || processingPhotoCount > 0 || photos.length === 0" @click="deleteAllPhotos">
            {{ deletingAllPhotos ? "กำลังลบ..." : deleteAllError ? "ลองลบที่เหลืออีกครั้ง" : `ยืนยันลบ ${deleteAllProgress.total || photos.length} รูป` }}
          </button>
        </div>
      </div>
    </dialog>
  </main>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import {
  GetGalleryPhotos,
  GetGalleryPhotoPreview,
  DeleteGalleryPhoto,
  SearchGalleryPhotos,
  PauseEmbeddingQueue,
  ResumeEmbeddingQueue,
  IsEmbeddingQueuePaused,
} from "../../wailsjs/go/controller/GalleryController";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import type { dto } from "../../wailsjs/go/models";

type EmbeddingStatusUpdate = {
  photoId: number;
  status: string;
  error: string;
};

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
const queuePaused = ref(false);
const queueActionPending = ref(false);

const searchText = ref("");
const searchResults = ref<dto.PhotoSearchResult[]>([]);
const hasSearched = ref(false);
const searching = ref(false);
const searchPending = ref(false);
const searchError = ref("");
let searchDebounceTimer: ReturnType<typeof setTimeout> | undefined;
let searchRequestID = 0;

const readyPhotos = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "completed"),
);
const pendingPhotoCount = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "pending").length,
);
const processingPhotoCount = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "processing").length,
);
const failedPhotoCount = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "failed").length,
);
const failedPhotos = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "failed"),
);
const embeddingProgress = computed(() => {
  const total = photos.value.length;
  const finished = readyPhotos.value.length + failedPhotoCount.value;
  return {
    total,
    finished,
    percent: total === 0 ? 0 : Math.round((finished / total) * 100),
    isRunning: pendingPhotoCount.value > 0 || processingPhotoCount.value > 0,
  };
});

async function toggleEmbeddingQueue() {
  if (queueActionPending.value) return;

  queueActionPending.value = true;
  try {
    if (queuePaused.value) {
      await ResumeEmbeddingQueue();
      queuePaused.value = false;
    } else {
      await PauseEmbeddingQueue();
      queuePaused.value = true;
    }
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    searchError.value = `เปลี่ยนสถานะคิว Embedding ไม่สำเร็จ: ${message}`;
  } finally {
    queueActionPending.value = false;
  }
}

async function searchGallery(text: string, requestID: number) {
  if (!text) {
    hasSearched.value = false;
    searchResults.value = [];
    searchError.value = "";
    return;
  }

  searching.value = true;
  searchPending.value = false;
  searchError.value = "";

  try {
    const results = await SearchGalleryPhotos(text);
    if (requestID === searchRequestID) {
      searchResults.value = results;
      hasSearched.value = true;
    }
  } catch (err) {
    if (requestID === searchRequestID) {
      searchError.value = err instanceof Error ? err.message : String(err);
    }
  } finally {
    if (requestID === searchRequestID) searching.value = false;
  }
}

watch(searchText, (value) => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  const requestID = ++searchRequestID;
  const text = value.trim();
  searchError.value = "";

  if (!text) {
    searchPending.value = false;
    searching.value = false;
    hasSearched.value = false;
    searchResults.value = [];
    return;
  }

  searchPending.value = true;
  searchDebounceTimer = setTimeout(() => {
    void searchGallery(text, requestID);
  }, 400);
});

function clearSearch() {
  searchText.value = "";
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchRequestID += 1;
  searchPending.value = false;
  searching.value = false;
  hasSearched.value = false;
  searchResults.value = [];
  searchError.value = "";
}

const searchMatchedPhotos = computed(() => {
  if (!hasSearched.value || !searchText.value.trim()) return readyPhotos.value;

  const photosByID = new Map(readyPhotos.value.map((photo) => [photo.id, photo]));
  return searchResults.value
    .map((result) => photosByID.get(result.id))
    .filter((photo): photo is dto.GalleryPhoto => photo !== undefined);
});

const displayedPhotos = searchMatchedPhotos;

async function loadCompletedPreview(photo: dto.GalleryPhoto) {
  if (photo.status !== "active" || previewUrls.value[photo.id]) return;
  try {
    previewUrls.value[photo.id] = await GetGalleryPhotoPreview(photo.id);
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    errMessage.value = `${photo.filename}: ${message}`;
  }
}

function onEmbeddingStatus(update: EmbeddingStatusUpdate) {
  const photo = photos.value.find((item) => item.id === update.photoId);
  if (!photo) return;

  photo.analysisStatus = update.status;
  photo.analysisError = update.error;
  if (update.status === "completed") void loadCompletedPreview(photo);
}

async function loadPhotos() {
  loading.value = true;
  errMessage.value = "";

  try {
    photos.value = await GetGalleryPhotos();

    previewUrls.value = {};
    await loadPreviews(readyPhotos.value);
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
  if (photos.value.length === 0 || loading.value || deletingAllPhotos.value || processingPhotoCount.value > 0) return;

  deleteAllError.value = "";
  deleteAllProgress.value = { completed: 0, total: photos.value.length };
  deleteAllDialog.value?.showModal();
}

async function deleteAllPhotos() {
  if (
    deletingAllPhotos.value ||
    deletingPhotoID.value !== null ||
    processingPhotoCount.value > 0
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
      deleteAllProgress.value = { completed: 0, total: 0 };
    }
  } finally {
    deletingAllPhotos.value = false;
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

let stopEmbeddingListener: (() => void) | undefined;
onMounted(() => {
  stopEmbeddingListener = EventsOn("embedding:status", onEmbeddingStatus);
  void IsEmbeddingQueuePaused().then((paused) => {
    queuePaused.value = paused;
  }).catch((err) => console.error("อ่านสถานะคิว Embedding ไม่สำเร็จ", err));
  void loadPhotos();
});
onBeforeUnmount(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchRequestID += 1;
  stopEmbeddingListener?.();
});
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
