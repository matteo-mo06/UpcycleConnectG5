package models

type Topic struct {
	IdTopic      string `json:"id"`
	TitleTopic   string `json:"title"`
	IdAuthor     string `json:"id_author"`
	AuthorName   string `json:"author_name"`
	RepliesCount int    `json:"replies_count"`
	CreatedAt    string `json:"created_at"`
}

type TopicDetail struct {
	IdTopic    string `json:"id"`
	TitleTopic string `json:"title"`
	IdAuthor   string `json:"id_author"`
	AuthorName string `json:"author_name"`
	CreatedAt  string `json:"created_at"`
	Posts      []Post `json:"posts"`
}

type Post struct {
	IdPost       string  `json:"id"`
	BodyPost     string  `json:"body"`
	IdAuthor     string  `json:"id_author"`
	AuthorName   string  `json:"author_name"`
	IdParentPost *string `json:"id_parent_post"`
	CreatedAt    string  `json:"created_at"`
}

type CreateTopicRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreatePostRequest struct {
	Body         string  `json:"body"`
	IdParentPost *string `json:"id_parent_post"`
}

type UpdateTopicRequest struct {
	Title string `json:"title"`
}

type UpdatePostRequest struct {
	Body string `json:"body"`
}
