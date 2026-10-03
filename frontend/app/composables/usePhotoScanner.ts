import { computed, onMounted, ref } from "vue";

import {
  GetPhotoPreview,
  ScanPhotos,
} from "../../wailsjs/go/controller/PhotoController";
import { GetSetting } from "../../wailsjs/go/controller/SettingController";
import type { dto } from "../../wailsjs/go/models";
import type { ScanMode, ScanOption } from "../types/photo-scan";

export function usePhotoScanner() {
  const scanOptions: ScanOption[] = [
    {
      mode: "source",
      label: "Scan Source",
      description: "ค้นหารูปต้นทางและเลือกภาพที่ต้องการย้าย",
    },
    {
      mode: "destination",
      label: "Scan Destination",
      description: "ตรวจดูรูปที่จัดเก็บอยู่ในโฟลเดอร์ปลายทาง",
    },
  ];

  const source = ref("");
  const destination = ref("");
  const activeMode = ref<ScanMode>("source");
  const photosByMode = ref<Record<ScanMode, dto.ScannedPhoto[]>>({
    source: [],
    destination: [],
  });
  const hasScannedByMode = ref<Record<ScanMode, boolean>>({
    source: false,
    destination: false,
  });
  const selectedPaths = ref<string[]>([]);
  const loadingSetting = ref(true);
  const scanningMode = ref<ScanMode | null>(null);
  const errorMessage = ref("");
  const previewUrls = ref<Record<string, string>>({});
  const previewFinished = ref<Record<string, boolean>>({});

  const currentPath = computed(() =>
    activeMode.value === "source" ? source.value : destination.value,
  );
  const currentFolderLabel = computed(() =>
    activeMode.value === "source" ? "Source" : "Destination",
  );
  const currentPhotos = computed(() => photosByMode.value[activeMode.value]);
  const currentHasScanned = computed(
    () => hasScannedByMode.value[activeMode.value],
  );
  const isCurrentScanning = computed(
    () => scanningMode.value === activeMode.value,
  );
  const canScan = computed(
    () =>
      currentPath.value.trim() !== "" &&
      !loadingSetting.value &&
      scanningMode.value === null,
  );
  const allSelected = computed(() => {
    const sourcePhotos = photosByMode.value.source;

    return (
      sourcePhotos.length > 0 &&
      selectedPaths.value.length === sourcePhotos.length
    );
  });

  function getErrorMessage(error: unknown): string {
    return error instanceof Error ? error.message : String(error);
  }

  function setActiveMode(mode: ScanMode) {
    activeMode.value = mode;
    errorMessage.value = "";
  }

  async function loadSetting() {
    loadingSetting.value = true;

    try {
      const setting = await GetSetting();
      source.value = setting.Source ?? "";
      destination.value = setting.Dest ?? "";
    } catch (error) {
      errorMessage.value = getErrorMessage(error);
    } finally {
      loadingSetting.value = false;
    }
  }

  async function scan() {
    if (!canScan.value) return;

    const mode = activeMode.value;
    const root = currentPath.value;

    scanningMode.value = mode;
    hasScannedByMode.value[mode] = false;
    errorMessage.value = "";

    if (mode === "source") selectedPaths.value = [];

    try {
      const scannedPhotos = await ScanPhotos(root);

      photosByMode.value[mode] = scannedPhotos;
      hasScannedByMode.value[mode] = true;

      for (const photo of scannedPhotos) {
        previewUrls.value[photo.path] = "";
        previewFinished.value[photo.path] = false;
      }

      void loadPreviews(scannedPhotos, root);
    } catch (error) {
      photosByMode.value[mode] = [];
      errorMessage.value = getErrorMessage(error);
    } finally {
      scanningMode.value = null;
    }
  }

  function toggleSelectAll() {
    if (allSelected.value) {
      selectedPaths.value = [];
      return;
    }

    selectedPaths.value = photosByMode.value.source.map((photo) => photo.path);
  }

  async function loadPreviews(items: dto.ScannedPhoto[], root: string) {
    let nextIndex = 0;

    async function worker() {
      while (nextIndex < items.length) {
        const photo = items[nextIndex];
        nextIndex += 1;

        if (!photo) return;

        try {
          previewUrls.value[photo.path] = await GetPhotoPreview(root, photo.path);
        } catch {
          previewUrls.value[photo.path] = "";
        } finally {
          previewFinished.value[photo.path] = true;
        }
      }
    }

    const workerCount = Math.min(4, items.length);
    await Promise.all(Array.from({ length: workerCount }, () => worker()));
  }

  onMounted(loadSetting);

  return {
    scanOptions,
    activeMode,
    selectedPaths,
    loadingSetting,
    scanningMode,
    errorMessage,
    previewUrls,
    previewFinished,
    currentPath,
    currentFolderLabel,
    currentPhotos,
    currentHasScanned,
    isCurrentScanning,
    canScan,
    allSelected,
    setActiveMode,
    scan,
    toggleSelectAll,
  };
}
