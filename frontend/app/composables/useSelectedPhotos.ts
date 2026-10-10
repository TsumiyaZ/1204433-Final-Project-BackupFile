import type { dto } from "../../wailsjs/go/models";

// Shares the photos picked on the Scan page with the Move page.
export function useSelectedPhotos() {
  const selectedPhotos = useState<dto.ScannedPhoto[]>(
    "selected-photos",
    () => [],
  );
  const selectedPreviews = useState<Record<string, string>>(
    "selected-previews",
    () => ({}),
  );

  function setSelection(
    photos: dto.ScannedPhoto[],
    previews: Record<string, string>,
  ) {
    selectedPhotos.value = photos;
    selectedPreviews.value = Object.fromEntries(
      photos.map((photo) => [photo.path, previews[photo.path] ?? ""]),
    );
  }

  function clearSelection() {
    selectedPhotos.value = [];
    selectedPreviews.value = {};
  }

  return { selectedPhotos, selectedPreviews, setSelection, clearSelection };
}
