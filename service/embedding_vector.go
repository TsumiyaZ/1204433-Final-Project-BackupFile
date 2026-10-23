package service

import (
	"encoding/binary"
	"fmt"
	"math"
)

const embeddingDimensions = 768

func encodeEmbedding(vector []float32) ([]byte, error) {
	if len(vector) != embeddingDimensions {
		return nil, fmt.Errorf(
			"expected %d values, got %d",
			embeddingDimensions,
			len(vector),
		)
	}

	data := make([]byte, embeddingDimensions*4)
	for i, value := range vector {
		binary.LittleEndian.PutUint32(
			data[i*4:],
			math.Float32bits(value),
		)
	}

	return data, nil
}

func decodeEmbedding(data []byte) ([]float32, error) {
	if len(data) != embeddingDimensions*4 {
		return nil, fmt.Errorf("invalid embedding data size: %d", len(data))
	}

	vector := make([]float32, embeddingDimensions)
	for i := range vector {
		bits := binary.LittleEndian.Uint32(data[i*4:])
		vector[i] = math.Float32frombits(bits)
	}

	return vector, nil
}

func cosineSimilarity(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector dimensions do not match")
	}
	if len(a) == 0 {
		return 0, fmt.Errorf("vectors cannot be empty")
	}

	var dotProduct, normA, normB float64

	for i := range a {
		x := float64(a[i])
		y := float64(b[i])

		dotProduct += x * y
		normA += x * x
		normB += y * y
	}

	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("cannot compare a zero vector")
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
