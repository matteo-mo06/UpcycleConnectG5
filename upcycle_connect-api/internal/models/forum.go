package models

type Topic struct {
	Id_topic     string `json:"id"`
	Title_topic  string `json:"title"`
	Id_author    string `json:"id_author"`
	AuthorName   string `json:"author_name"`
	RepliesCount int    `json:"replies_count"`
	CreatedAt    string `json:"created_at"`
}

type TopicDetail struct {
	Id_topic    string `json:"id"`
	Title_topic string `json:"title"`
	Id_author   string `json:"id_author"`
	AuthorName  string `json:"author_name"`
	CreatedAt   string `json:"created_at"`
	Posts       []Post `json:"posts"`
}

type Post struct {
	Id_post        string  `json:"id"`
	Body_post      string  `json:"body"`
	Id_author      string  `json:"id_author"`
	AuthorName     string  `json:"author_name"`
	Id_parent_post *string `json:"id_parent_post"`
	CreatedAt      string  `json:"created_at"`
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
