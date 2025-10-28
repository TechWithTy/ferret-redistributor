package instagram

type createMediaParams struct {
    Caption             string
    ImageURL            string
    VideoURL            string
    IsCarousel          bool
    Children            []string
    IsReel              bool
    IsStory             bool
    ThumbOffsetSeconds  int
    DisableComments     bool
    CoverURL            string
}

type publishParams struct {
    CreationID string
}

