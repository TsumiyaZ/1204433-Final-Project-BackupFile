package service

import (
	"context"
	"fmt"
	"sort"
)

type tagConcept struct {
	text  string
	label string
}

type scoredTag struct {
	label string
	score float64
}

// These concepts provide a small, stable vocabulary for similarity-based
// tagging. They are labels, not generated captions.
var tagConcepts = []tagConcept{
	{"a person or people", "คน"},
	{"a dog", "สุนัข"}, {"a cat", "แมว"}, {"a bird", "นก"}, {"a horse", "ม้า"},
	{"a sheep", "แกะ"}, {"a cow", "วัว"}, {"an elephant", "ช้าง"}, {"a bear", "หมี"},
	{"a zebra", "ม้าลาย"}, {"a giraffe", "ยีราฟ"}, {"a bicycle", "จักรยาน"},
	{"a car", "รถยนต์"}, {"a motorcycle", "รถจักรยานยนต์"}, {"an airplane", "เครื่องบิน"},
	{"a bus", "รถโดยสาร"}, {"a train", "รถไฟ"}, {"a truck", "รถบรรทุก"}, {"a boat", "เรือ"},
	{"a traffic light", "สัญญาณไฟจราจร"}, {"a street or road", "ถนน"}, {"a building", "อาคาร"},
	{"a beach", "ชายหาด"}, {"the ocean or sea", "ทะเล"}, {"a river or lake", "แม่น้ำหรือทะเลสาบ"},
	{"a mountain", "ภูเขา"}, {"a forest or trees", "ป่าและต้นไม้"}, {"flowers or plants", "ดอกไม้และพืช"},
	{"the sky or clouds", "ท้องฟ้าและเมฆ"}, {"food", "อาหาร"}, {"fruit", "ผลไม้"},
	{"a banana", "กล้วย"}, {"an apple", "แอปเปิล"}, {"pizza", "พิซซ่า"}, {"a cake or dessert", "ขนมหรือเค้ก"},
	{"a cup or glass", "แก้วน้ำ"}, {"a bottle", "ขวด"}, {"a dining table", "โต๊ะอาหาร"},
	{"a chair", "เก้าอี้"}, {"a couch or sofa", "โซฟา"}, {"a bed", "เตียง"},
	{"a television", "โทรทัศน์"}, {"a laptop or computer", "คอมพิวเตอร์"}, {"a mobile phone", "โทรศัพท์มือถือ"},
	{"a book", "หนังสือ"}, {"a clock", "นาฬิกา"}, {"sports or a game", "กีฬา"},
	{"a ball", "ลูกบอล"}, {"a kite", "ว่าว"}, {"a backpack", "กระเป๋าเป้"}, {"an umbrella", "ร่ม"},
	{"a kitchen", "ห้องครัว"}, {"a bathroom", "ห้องน้ำ"}, {"an indoor room", "ภายในอาคาร"},
	{"an outdoor scene", "กลางแจ้ง"}, {"a landscape", "ทิวทัศน์"},
}

const minimumTagSimilarity = 0.28
const maximumPhotoTags = 5

func (s *photoService) tagsFromEmbedding(ctx context.Context, imageVector []float32) ([]string, error) {
	s.tagMu.Lock()
	defer s.tagMu.Unlock()

	scores := make([]scoredTag, 0, len(tagConcepts))
	for _, concept := range tagConcepts {
		textVector, ok := s.tagVectors[concept.text]
		if !ok {
			var err error
			textVector, err = s.engine.EmbedText(ctx, concept.text)
			if err != nil {
				return nil, fmt.Errorf("cannot embed tag candidate %q: %w", concept.text, err)
			}
			s.tagVectors[concept.text] = textVector
		}
		score, err := cosineSimilarity(imageVector, textVector)
		if err != nil {
			return nil, err
		}
		if score >= minimumTagSimilarity {
			scores = append(scores, scoredTag{label: concept.label, score: score})
		}
	}

	sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
	if len(scores) > maximumPhotoTags {
		scores = scores[:maximumPhotoTags]
	}

	tags := make([]string, 0, len(scores))
	for _, item := range scores {
		tags = append(tags, item.label)
	}
	if len(tags) == 0 {
		tags = []string{"ทั่วไป"}
	}
	return tags, nil
}
