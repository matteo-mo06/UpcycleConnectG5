package models

type Topic struct {
	IdTopic           string  `json:"id"`
	TitleTopic        string  `json:"title"`
	IdAuthor          string  `json:"id_author"`
	AuthorName        string  `json:"author_name"`
	RepliesCount      int     `json:"replies_count"`
	ModeratedBy       *string `json:"moderated_by"`
	ModerationMessage *string `json:"moderation_message"`
	CreatedAt         string  `json:"created_at"`
}

type DeletedTopic struct {
	IdTopic    string `json:"id"`
	TitleTopic string `json:"title"`
	IdAuthor   string `json:"id_author"`
	AuthorName string `json:"author_name"`
	DeletedAt  string `json:"deleted_at"`
	Posts      []Post `json:"posts"`
}

type TopicDetail struct {
	IdTopic           string  `json:"id"`
	TitleTopic        string  `json:"title"`
	IdAuthor          string  `json:"id_author"`
	AuthorName        string  `json:"author_name"`
	ModeratedBy       *string `json:"moderated_by"`
	ModerationMessage *string `json:"moderation_message"`
	CreatedAt         string  `json:"created_at"`
	Posts             []Post  `json:"posts"`
}

type Post struct {
	IdPost            string  `json:"id"`
	BodyPost          string  `json:"body"`
	IdAuthor          string  `json:"id_author"`
	AuthorName        string  `json:"author_name"`
	IdParentPost      *string `json:"id_parent_post"`
	ModeratedBy       *string `json:"moderated_by"`
	ModerationMessage *string `json:"moderation_message"`
	CreatedAt         string  `json:"created_at"`
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

type ModerateContentRequest struct {
	Message string `json:"message"`
}
