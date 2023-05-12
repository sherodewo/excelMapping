package model

import "time"

type Header struct {
	HeaderID           uint      `json:"header_id"`
	CatalogueHeaderID  uint      `json:"catalogue_header_id"`
	PublisherID        uint      `json:"publisher_id"`
	MessageVersion     string    `json:"message_version"`
	Profile            string    `json:"profile"`
	ProfileVersion     string    `json:"profile_version"`
	MessageId          string    `json:"message_id"`
	DateTimeCreated    time.Time `json:"date_time_created"`
	FileNumber         int       `json:"file_number"`
	NumberOfFiles      int       `json:"number_of_files"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	SenderPartyId      string    `json:"sender_party_id"`
	SenderName         string    `json:"sender_name"`
	ServiceDescription string    `json:"service_description"`
	RecipientName      string    `json:"recipient_name"`
	TotalUsages        float64   `json:"total_usages"`
	NetRevenue         float64   `json:"net_revenue"`
	CommercialModel    string    `json:"commercial_model"`
	UseType            string    `json:"use_type"`
	Subscribers        int       `json:"subscribers"`
}

type SongUsage struct {
	BlockId                  int      `json:"block_id"`                   //block_id_dsrf
	ResourceReference        int      `json:"resource_reference"`         //resource_id_dsrf
	ResourceWamiReference    string   `json:"resource_wami_reference"`    //resource_id_dsrf_wami
	ResourceYoutubeReference string   `json:"resource_youtube_reference"` //resource_youtube_reference
	DspResourceId            string   `json:"dsp_resource_id"`
	ISRC                     string   `json:"isrc"`
	Title                    string   `json:"title"`
	SubTitle                 string   `json:"sub_title"`
	DisplayArtistName        string   `json:"display_artist_name"`
	Duration                 string   `json:"duration"`
	ResourceType             string   `json:"resource_type"`
	ComposerAuthor           []string `json:"composer_author"`
	IsRoyaltyBearing         bool     `json:"is_royalty_bearing"`
	NumberOfStreams          float64  `json:"number_of_streams"` //total_listener
	DSPReleaseID             string   `json:"dsp_release_id"`
	Offer                    string   `json:"offer"`
	IndexNumber              string   `json:"index_number"`
	Revenue                  float64  `json:"revenue"`
}

type ServiceMeta struct {
	SummaryRecordID    int     `json:"summary_record_id"`
	CommercialMode     string  `json:"commercial_mode"`
	UseType            string  `json:"use_type"`
	ServiceDescription string  `json:"service_description"`
	Usages             float64 `json:"usages"`
	NetRevenue         float64 `json:"net_revenue"`
}
