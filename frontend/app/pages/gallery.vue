<template>
  <main class="gallery-page mx-auto max-w-[1500px] px-5 py-6 sm:px-8">
    <header class="gallery-header">

      <div class="gallery-search-row">
        <form class="search-shell" role="search" @submit.prevent>
          <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="search-icon">
            <circle cx="10.8" cy="10.8" r="6.8" stroke="currentColor" stroke-width="2" />
            <path d="m16 16 4.2 4.2" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
          <label class="sr-only" for="gallery-search">ค้นหารูปภาพ</label>
          <input
            id="gallery-search"
            v-model="searchText"
            type="text"
            autocomplete="off"
            placeholder="อธิบายสิ่งที่จำได้ เช่น เรือกลางทะเล หรือ หมา..."
          />
          <button v-if="searchText" type="button" class="clear-search" aria-label="ล้างคำค้นหา" @click="clearSearch">
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" class="h-5 w-5">
              <path d="m6 6 8 8m0-8-8 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </form>
        <button
          type="button"
          class="btn btn-square btn-outline refresh-button"
          :aria-label="loading ? 'กำลังโหลดรูปภาพ' : 'รีเฟรชรูปภาพ'"
          :aria-busy="loading"
          :title="loading ? 'กำลังโหลดรูปภาพ...' : 'รีเฟรชรูปภาพ'"
          :disabled="loading || deletingPhotoID !== null || deletingAllPhotos"
          @click="refreshPhoto"
        >
          <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" class="h-5 w-5" :class="{ 'refresh-icon-spinning': loading }">
            <path d="M20 7v5h-5M4 17v-5h5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
            <path d="M5.64 5.64A9 9 0 0 1 20 12M4 12a9 9 0 0 0 14.36 6.36" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
        </button>
      </div>

      <div class="gallery-filter-row" aria-label="ตัวกรองรูปภาพ">
        <button type="button" class="filter-chip" :class="{ 'filter-chip-active': activeFilter === 'all' }" :aria-pressed="activeFilter === 'all'" @click="activeFilter = 'all'">
          ทั้งหมด <span>{{ readyPhotos.length }}</span>
        </button>
        <button type="button" class="filter-chip" :class="{ 'filter-chip-active': activeFilter === 'waiting' }" :aria-pressed="activeFilter === 'waiting'" @click="activeFilter = 'waiting'">
          กำลังจัดทำดัชนี <span>{{ pendingPhotoCount + processingPhotoCount }}</span>
        </button>
        <button type="button" class="filter-chip" :class="{ 'filter-chip-active': activeFilter === 'failed' }" :aria-pressed="activeFilter === 'failed'" @click="activeFilter = 'failed'">
          ต้องตรวจสอบ <span>{{ failedPhotoCount }}</span>
        </button>
      </div>
      <p v-if="searching || searchPending" class="search-status" role="status" aria-live="polite">
        <span class="loading loading-spinner loading-xs text-primary" aria-hidden="true" />
        {{ searching ? "กำลังค้นหารูปภาพ..." : "กำลังรอค้นหา..." }}
      </p>
      <p v-if="searchError" role="alert" class="mt-2 text-sm text-error">{{ searchError }}</p>
    </header>

    <section v-if="photos.length > 0 && !loading && !errMessage" class="gallery-toolbar" aria-label="จัดการคลังรูปภาพ">
      <p class="text-sm text-base-content/65" aria-live="polite">
        <template v-if="activeFilter === 'all'">{{ displayedPhotos.length }} รูปภาพ</template>
        <template v-else-if="activeFilter === 'waiting'">{{ pendingPhotoCount + processingPhotoCount }} รูปรอจัดทำดัชนี</template>
        <template v-else>{{ failedPhotoCount }} รูปต้องตรวจสอบ</template>
      </p>
      <button type="button" class="btn btn-sm btn-ghost text-error hover:bg-error/10" :disabled="loading || deletingAllPhotos || deletingPhotoID !== null || processingPhotoCount > 0 || photos.length === 0" @click="openDeleteAllDialog">
        {{ deletingAllPhotos ? `กำลังลบ ${deleteAllProgress.completed}/${deleteAllProgress.total}...` : "ลบทั้งหมด" }}
      </button>
    </section>

    <div v-if="loading" class="gallery-message" aria-live="polite">
      <span class="loading loading-spinner loading-md text-primary" aria-hidden="true"></span>
      <span>กำลังโหลดรูปภาพ...</span>
    </div>
    <div v-else-if="errMessage" role="alert" class="gallery-message error-message">
      <p>{{ errMessage }}</p>
      <button type="button" class="btn btn-sm btn-outline" @click="refreshPhoto">ลองอีกครั้ง</button>
    </div>
    <div v-else-if="photos.length === 0" class="gallery-empty-state">
      <div class="empty-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none" class="size-9"><rect x="3" y="4" width="18" height="16" rx="3" stroke="currentColor" stroke-width="1.6"/><circle cx="8.5" cy="9" r="1.5" fill="currentColor"/><path d="m5 17 5-5 3 3 2-2 4 4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </div>
      <h2 class="text-xl font-semibold">คลังรูปภาพยังว่างอยู่</h2>
      <p class="max-w-md text-sm leading-6 text-base-content/65">ย้ายรูปไปยัง Destination ก่อน แล้วรูปจะปรากฏที่นี่เพื่อให้ค้นหาด้วยคำอธิบายได้</p>
      <NuxtLink to="/scan" class="btn btn-primary mt-2 min-h-11">ไปสแกนรูปภาพ</NuxtLink>
    </div>
    <section v-else-if="activeFilter === 'waiting' || activeFilter === 'failed'" class="status-list" :aria-label="activeFilter === 'waiting' ? 'รูปที่รอจัดทำดัชนี' : 'รูปที่ต้องตรวจสอบ'">
      <template v-if="activeFilter === 'waiting'">
        <div v-if="waitingPhotos.length === 0" class="gallery-message">ไม่มีรูปที่รอจัดทำดัชนี</div>
        <article v-for="photo in waitingPhotos" :key="photo.id" class="status-row">
          <span class="status-dot status-dot-waiting" aria-hidden="true" />
          <div class="min-w-0 flex-1"><p class="truncate font-medium">{{ photo.filename }}</p><p class="mt-1 text-xs text-base-content/60">{{ photo.analysisStatus === 'processing' ? 'กำลังสร้าง Embedding' : 'รอสร้าง Embedding' }}</p></div>
        </article>
      </template>
      <template v-else>
        <div v-if="failedPhotos.length === 0" class="gallery-message">ไม่มีรูปที่ต้องตรวจสอบ</div>
        <article v-for="photo in failedPhotos" :key="photo.id" class="status-row">
          <span class="status-dot status-dot-failed" aria-hidden="true" />
          <div class="min-w-0 flex-1"><p class="truncate font-medium">{{ photo.filename }}</p><p class="mt-1 break-words text-xs text-error">{{ photo.analysisError || 'สร้าง Embedding ไม่สำเร็จ' }}</p></div>
        </article>
      </template>
    </section>
    <div v-else-if="displayedPhotos.length === 0" class="gallery-message">
      <p class="font-semibold">{{ hasSearched && searchText.trim() ? `ไม่พบรูปที่ตรงกับ “${searchText}”` : "ยังไม่มีรูปที่สร้าง Embedding เสร็จ" }}</p>
      <p v-if="pendingPhotoCount || processingPhotoCount">รูปจะแสดงใน Gallery เมื่อสร้าง Embedding เสร็จ</p>
      <p v-if="failedPhotoCount" class="text-error">มี {{ failedPhotoCount }} รูปที่สร้าง Embedding ไม่สำเร็จ กรุณารีเฟรชเพื่อตรวจสอบสถานะ</p>
      <button v-if="searchText" type="button" class="btn btn-sm btn-outline" @click="clearSearch">ล้างคำค้นหา</button>
    </div>

    <section v-else class="thumbnail-grid" aria-label="ผลลัพธ์รูปภาพ">
      <article v-for="photo in displayedPhotos" :key="photo.id" class="photo-tile">
        <div class="thumbnail-frame">
          <button type="button" class="photo-open-trigger" :aria-label="`เปิดรายละเอียดรูป ${photo.filename}`" @click="openPhotoDetails(photo)">
            <img v-if="previewUrls[photo.id]" :src="previewUrls[photo.id]" :alt="photo.filename" loading="lazy" decoding="async" />
            <span v-else class="thumbnail-placeholder">
              {{ photo.status === "missing" ? "ไม่พบไฟล์รูปภาพ" : "กำลังโหลดตัวอย่าง..." }}
            </span>
          </button>
          <span class="photo-status status-complete"><span aria-hidden="true">✓</span> พร้อมค้นหา</span>
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

    <dialog ref="photoDetailsDialog" class="modal photo-details-dialog" aria-labelledby="photo-details-title" @close="selectedPhoto = null">
      <section v-if="selectedPhoto" class="modal-box photo-details-box">
        <header class="photo-details-header">
          <div class="min-w-0">
            <p class="text-xs font-semibold uppercase tracking-[0.16em] text-base-content/55">Photo details</p>
            <h2 id="photo-details-title" class="mt-1 truncate text-lg font-semibold" :title="selectedPhoto.filename">{{ selectedPhoto.filename }}</h2>
          </div>
          <button type="button" class="btn btn-square btn-ghost min-h-11 min-w-11" aria-label="ปิดรายละเอียดรูป" @click="photoDetailsDialog?.close()">
            <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" class="size-5"><path d="m5 5 10 10M15 5 5 15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" /></svg>
          </button>
        </header>

        <div class="photo-details-preview">
          <img v-if="previewUrls[selectedPhoto.id]" :src="previewUrls[selectedPhoto.id]" :alt="selectedPhoto.filename" @load="readDetailImageSize" />
          <div v-else class="thumbnail-placeholder">
            <span v-if="photoPreviewPending" class="loading loading-spinner loading-md text-primary" aria-hidden="true" />
            {{ photoPreviewPending ? "กำลังโหลดตัวอย่าง..." : "ไม่สามารถแสดงตัวอย่างรูปนี้ได้" }}
          </div>
        </div>

        <div class="photo-details-meta">
          <p class="break-all font-semibold">{{ selectedPhoto.filename }}</p>
          <p class="mt-1 break-all text-xs leading-5 text-base-content/60">{{ selectedPhoto.path }}</p>
          <p class="mt-2 text-xs text-base-content/65">
            {{ selectedPhoto.filename.split('.').pop()?.toUpperCase() }}
            <span v-if="selectedPhoto.modifiedAt"> · {{ formatModifiedAt(selectedPhoto.modifiedAt) }}</span>
            <span v-if="detailImageSize.width && detailImageSize.height"> · {{ detailImageSize.width }} × {{ detailImageSize.height }}</span>
            <span v-if="selectedPhoto.size"> · {{ formatFileSize(selectedPhoto.size) }}</span>
          </p>
        </div>

        <p v-if="photoActionError" class="mt-3 text-sm text-error" role="alert">{{ photoActionError }}</p>
        <footer class="photo-details-actions">
          <button type="button" class="btn btn-primary min-h-11 min-w-24" :disabled="photoActionPending !== null || selectedPhoto.status !== 'active'" @click="openSelectedPhoto">
            <span v-if="photoActionPending === 'open'" class="loading loading-spinner loading-xs" />
            Open
          </button>
          <button type="button" class="btn btn-outline min-h-11" :disabled="photoActionPending !== null || selectedPhoto.status !== 'active'" @click="showSelectedPhotoInFolder">
            <span v-if="photoActionPending === 'reveal'" class="loading loading-spinner loading-xs" />
            Show in Folder
          </button>
        </footer>
      </section>
      <form method="dialog" class="modal-backdrop"><button aria-label="ปิดรายละเอียดรูป">close</button></form>
    </dialog>

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
  OpenGalleryPhoto,
  ShowGalleryPhotoInFolder,
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
const photoDetailsDialog = ref<HTMLDialogElement | null>(null);
const selectedPhoto = ref<dto.GalleryPhoto | null>(null);
const photoActionPending = ref<"open" | "reveal" | null>(null);
const photoActionError = ref("");
const photoPreviewPending = ref(false);
const detailImageSize = ref({ width: 0, height: 0 });
const deletingPhotoID = ref<number | null>(null);
const deleteDialog = ref<HTMLDialogElement | null>(null);
const deleteAllDialog = ref<HTMLDialogElement | null>(null);
const pendingDelete = ref<dto.GalleryPhoto | null>(null);
const deletingAllPhotos = ref(false);
const deleteAllProgress = ref({ completed: 0, total: 0 });
const deleteAllError = ref("");
const queuePaused = ref(false);
const queueActionPending = ref(false);
const activeFilter = ref<"all" | "waiting" | "failed">("all");

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
const waitingPhotos = computed(() =>
  photos.value.filter((photo) => photo.analysisStatus === "pending" || photo.analysisStatus === "processing"),
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

const displayedPhotos = computed(() =>
  activeFilter.value === "all" ? searchMatchedPhotos.value : [],
);

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

function openPhotoDetails(photo: dto.GalleryPhoto) {
  selectedPhoto.value = photo;
  photoActionPending.value = null;
  photoActionError.value = "";
  photoPreviewPending.value = false;
  detailImageSize.value = { width: 0, height: 0 };
  photoDetailsDialog.value?.showModal();
  if (!previewUrls.value[photo.id] && photo.status === "active") {
    photoPreviewPending.value = true;
    void GetGalleryPhotoPreview(photo.id)
      .then((preview) => {
        previewUrls.value[photo.id] = preview;
      })
      .catch((err) => {
        if (selectedPhoto.value?.id === photo.id) {
          photoActionError.value = err instanceof Error ? err.message : String(err);
        }
      })
      .finally(() => {
        if (selectedPhoto.value?.id === photo.id) photoPreviewPending.value = false;
      });
  }
}

function readDetailImageSize(event: Event) {
  const image = event.currentTarget;
  if (!(image instanceof HTMLImageElement)) return;
  detailImageSize.value = { width: image.naturalWidth, height: image.naturalHeight };
}

function formatFileSize(size: number) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(0)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
}

function formatModifiedAt(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("en-US", { month: "short", day: "numeric", year: "numeric" }).format(date);
}

async function runPhotoAction(action: "open" | "reveal") {
  const photo = selectedPhoto.value;
  if (!photo || photoActionPending.value) return;

  photoActionPending.value = action;
  photoActionError.value = "";
  try {
    if (action === "open") await OpenGalleryPhoto(photo.id);
    else await ShowGalleryPhotoInFolder(photo.id);
  } catch (err) {
    photoActionError.value = err instanceof Error ? err.message : String(err);
  } finally {
    photoActionPending.value = null;
  }
}

function openSelectedPhoto() {
  void runPhotoAction("open");
}

function showSelectedPhotoInFolder() {
  void runPhotoAction("reveal");
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
.gallery-header {
  padding-bottom: 18px;
  border-bottom: 1px solid color-mix(in srgb, var(--color-base-content) 12%, transparent);
}

.gallery-search-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.search-shell {
  display: flex;
  min-width: 0;
  flex: 1 1 auto;
  min-height: 52px;
  align-items: center;
  gap: 13px;
  border: 1px solid var(--color-base-300);
  border-radius: 13px;
  background: #2d2e32;
  padding: 4px 7px 4px 15px;
  transition: border-color 160ms ease, box-shadow 160ms ease, background-color 160ms ease;
}

.search-shell:focus-within {
  border-color: #1687f8;
  box-shadow: 0 0 0 1px #1687f8, 0 0 0 3px rgb(22 135 248 / 18%);
  background: #303136;
}

.search-icon {
  width: 22px;
  height: 22px;
  flex: 0 0 auto;
  color: var(--color-accent);
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

.refresh-button {
  width: 48px;
  min-width: 48px;
  height: 48px;
  min-height: 48px;
  padding: 0;
  color: color-mix(in srgb, var(--color-base-content) 82%, transparent);
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.refresh-button:hover:not(:disabled) {
  border-color: var(--color-primary);
  background: color-mix(in srgb, var(--color-primary) 12%, transparent);
  color: var(--color-primary);
}

.refresh-button:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

.refresh-icon-spinning {
  animation: refresh-spin 900ms linear infinite;
}

@keyframes refresh-spin {
  to { transform: rotate(360deg); }
}

.clear-search {
  display: grid;
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  place-items: center;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: color-mix(in srgb, var(--color-base-content) 70%, transparent);
  cursor: pointer;
}

.clear-search:hover {
  background: color-mix(in srgb, var(--color-base-content) 12%, transparent);
}

.gallery-filter-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
}

.filter-chip {
  display: inline-flex;
  min-height: 34px;
  align-items: center;
  gap: 7px;
  border: 1px solid transparent;
  border-radius: 999px;
  background: #34363a;
  padding: 0 12px;
  color: #e4e4e6;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 160ms ease, color 160ms ease, border-color 160ms ease;
}

.filter-chip:hover {
  background: #42444a;
}

.filter-chip span {
  color: #b5b6bb;
  font-variant-numeric: tabular-nums;
}

.filter-chip-active {
  border-color: #1687f8;
  background: #1687f8;
  color: #fff;
}

.filter-chip-active:hover {
  background: #0879e8;
}

.filter-chip-active span {
  color: rgb(255 255 255 / 85%);
}

.search-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
  color: color-mix(in srgb, var(--color-base-content) 72%, transparent);
  font-size: 0.8125rem;
}

.gallery-toolbar {
  display: flex;
  min-height: 54px;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 2px 10px;
}

.thumbnail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 190px), 1fr));
  gap: 16px;
}

.photo-tile {
  min-width: 0;
  overflow: hidden;
  border: 1px solid #383a3f;
  border-radius: 13px;
  background: var(--color-base-100);
  transition: border-color 160ms ease, box-shadow 160ms ease;
}

.photo-tile:hover {
  border-color: #1687f8;
  box-shadow: 0 8px 24px rgb(0 0 0 / 22%);
}

.thumbnail-frame {
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
  background: #1d1e21;
}

.photo-open-trigger {
  position: absolute;
  inset: 0;
  display: grid;
  width: 100%;
  height: 100%;
  place-items: center;
  border: 0;
  background: transparent;
  padding: 0;
  cursor: pointer;
}

.photo-open-trigger:focus-visible {
  z-index: 1;
  outline: 3px solid #1687f8;
  outline-offset: -3px;
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
  border: 1px solid rgb(22 24 28 / 12%);
  border-radius: 999px;
  padding: 4px 9px;
  background: rgb(240 245 250 / 94%);
  color: #225d42;
  font-size: 0.7rem;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-complete {
  color: #276749;
}

.photo-tile-body {
  padding: 11px 12px 12px;
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
  border: 1px dashed #45474d;
  border-radius: 14px;
  padding: 24px;
  color: color-mix(in srgb, var(--color-base-content) 82%, transparent);
  text-align: center;
}

.gallery-empty-state {
  display: flex;
  min-height: min(52vh, 460px);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 36px 20px;
  text-align: center;
}

.empty-icon {
  display: grid;
  width: 72px;
  height: 72px;
  place-items: center;
  margin-bottom: 4px;
  border: 1px solid #46484d;
  border-radius: 21px;
  background: #292a2e;
  color: var(--color-accent);
  box-shadow: 0 10px 40px rgb(0 0 0 / 20%);
}

.status-list {
  display: grid;
  gap: 8px;
  max-width: 860px;
}

.status-row {
  display: flex;
  min-height: 62px;
  align-items: center;
  gap: 12px;
  border: 1px solid var(--color-base-300);
  border-radius: 12px;
  background: var(--color-base-100);
  padding: 10px 14px;
}

.status-dot {
  width: 9px;
  height: 9px;
  flex: 0 0 auto;
  border-radius: 50%;
}

.status-dot-waiting { background: #e9a34b; }
.status-dot-failed { background: #ee7777; }

.error-message {
  border-color: color-mix(in srgb, var(--color-error) 48%, transparent);
  color: var(--color-error);
}

.photo-details-dialog::backdrop {
  background: rgb(0 0 0 / 62%);
}

.photo-details-box {
  width: min(720px, calc(100vw - 32px));
  max-width: 720px;
  max-height: 92vh;
  overflow-y: auto;
  border: 1px solid #45474d;
  border-radius: 16px;
  background: #252629;
  padding: 0;
  color: #f1f1f2;
  box-shadow: 0 24px 90px rgb(0 0 0 / 48%);
}

.photo-details-header {
  display: flex;
  min-height: 64px;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px 10px 20px;
}

.photo-details-preview {
  display: grid;
  min-height: 180px;
  max-height: 56vh;
  place-items: center;
  overflow: hidden;
  background: #191a1c;
}

.photo-details-preview img {
  display: block;
  width: 100%;
  max-height: 56vh;
  object-fit: contain;
}

.photo-details-meta {
  padding: 16px 20px 0;
}

.photo-details-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 16px 20px 20px;
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

  .photo-details-actions .btn {
    flex: 1 1 auto;
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
