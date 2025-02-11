package domain

type Link struct {
	ShortLink   string `yaml:"hash_id" json:"hash_id"`
	OriginalURL string `yaml:"original_url" json:"original_url"`
}
