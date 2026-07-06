package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func scanPost(row interface{ Scan(...any) error }) (models.Post, error) {
	var p models.Post
	var idParentPost, moderatedBy, moderationMessage sql.NullString
	err := row.Scan(&p.IdPost, &p.BodyPost, &p.IdAuthor, &p.AuthorName, &idParentPost, &moderatedBy, &moderationMessage, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	if idParentPost.Valid {
		p.IdParentPost = &idParentPost.String
	}
	if moderatedBy.Valid {
		p.ModeratedBy = &moderatedBy.String
	}
	if moderationMessage.Valid {
		p.ModerationMessage = &moderationMessage.String
	}
	return p, nil
}

func GetTopics(search string, limit, offset int) ([]models.Topic, int, error) {
	where := `WHERE t.deleted_at IS NULL`
	var args []any

	if search != "" {
		where += " AND t.title_topic LIKE ?"
		args = append(args, "%"+search+"%")
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow(`SELECT COUNT(*) FROM TOPIC t `+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT t.id_topic, t.title_topic, t.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       GREATEST(COUNT(tp.id_post) - 1, 0) AS replies_count,
		       t.moderated_by, t.moderation_message,
		       t.created_at
		FROM TOPIC t
		LEFT JOIN USER u ON u.id_user = t.id_author
		LEFT JOIN TOPIC_POST tp ON tp.id_topic = t.id_topic
		`+where+`
		GROUP BY t.id_topic, t.title_topic, t.id_author, u.first_name, u.last_name, t.moderated_by, t.moderation_message, t.created_at
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Topic
	for rows.Next() {
		var t models.Topic
		var moderatedBy, moderationMessage sql.NullString
		if err := rows.Scan(
			&t.IdTopic, &t.TitleTopic, &t.IdAuthor, &t.AuthorName, &t.RepliesCount, &moderatedBy, &moderationMessage, &t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		if moderatedBy.Valid {
			t.ModeratedBy = &moderatedBy.String
		}
		if moderationMessage.Valid {
			t.ModerationMessage = &moderationMessage.String
		}
		list = append(list, t)
	}
	return list, total, nil
}

func loadTopicPosts(topicID string) ([]models.Post, error) {
	rows, err := config.Conn.Query(`
		SELECT p.id_post, p.body_post, p.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       p.id_parent_post,
		       p.moderated_by, p.moderation_message,
		       p.created_at
		FROM POST p
		JOIN TOPIC_POST tp ON tp.id_post = p.id_post
		LEFT JOIN USER u ON u.id_user = p.id_author
		WHERE tp.id_topic = ? AND p.deleted_at IS NULL
		ORDER BY p.created_at ASC`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		if p.ModerationMessage != nil {
			p.BodyPost = *p.ModerationMessage
			p.AuthorName = "Modéré"
			p.IdAuthor = ""
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetDeletedTopics() ([]models.DeletedTopic, error) {
	rows, err := config.Conn.Query(`
		SELECT t.id_topic, t.title_topic, t.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       t.deleted_at
		FROM TOPIC t
		LEFT JOIN USER u ON u.id_user = t.id_author
		WHERE t.deleted_at IS NOT NULL
		ORDER BY t.deleted_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DeletedTopic
	for rows.Next() {
		var t models.DeletedTopic
		if err := rows.Scan(&t.IdTopic, &t.TitleTopic, &t.IdAuthor, &t.AuthorName, &t.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}

	for i := range list {
		posts, err := loadTopicPosts(list[i].IdTopic)
		if err != nil {
			return nil, err
		}
		list[i].Posts = posts
	}
	return list, nil
}

func GetTopicByID(id string) (models.TopicDetail, error) {
	var t models.TopicDetail
	var moderatedBy, moderationMessage sql.NullString
	err := config.Conn.QueryRow(`
		SELECT t.id_topic, t.title_topic, t.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       t.moderated_by, t.moderation_message,
		       t.created_at
		FROM TOPIC t
		LEFT JOIN USER u ON u.id_user = t.id_author
		WHERE t.id_topic = ? AND t.deleted_at IS NULL`, id,
	).Scan(&t.IdTopic, &t.TitleTopic, &t.IdAuthor, &t.AuthorName, &moderatedBy, &moderationMessage, &t.CreatedAt)
	if err != nil {
		return t, err
	}
	if moderatedBy.Valid {
		t.ModeratedBy = &moderatedBy.String
	}
	if moderationMessage.Valid {
		t.ModerationMessage = &moderationMessage.String
	}

	posts, err := loadTopicPosts(id)
	if err != nil {
		return t, err
	}
	t.Posts = posts
	return t, nil
}

func CreateTopic(authorID, title, body string) (string, error) {
	topicID := uuid.New().String()
	postID := uuid.New().String()

	if _, err := config.Conn.Exec(
		`INSERT INTO TOPIC (id_topic, title_topic, id_author) VALUES (?, ?, ?)`,
		topicID, title, authorID,
	); err != nil {
		return "", err
	}

	if _, err := config.Conn.Exec(
		`INSERT INTO POST (id_post, body_post, id_author) VALUES (?, ?, ?)`,
		postID, body, authorID,
	); err != nil {
		return "", err
	}

	if _, err := config.Conn.Exec(
		`INSERT INTO TOPIC_POST (id_topic, id_post) VALUES (?, ?)`,
		topicID, postID,
	); err != nil {
		return "", err
	}

	if _, err := config.Conn.Exec(
		`INSERT INTO USER_TOPIC (id_user, id_topic) VALUES (?, ?)`,
		authorID, topicID,
	); err != nil {
		return "", err
	}

	return topicID, nil
}

func CreatePost(topicID, authorID, body string, parentPostID *string) error {
	postID := uuid.New().String()

	var parent interface{}
	if parentPostID != nil && *parentPostID != "" {
		parent = *parentPostID
	}

	if _, err := config.Conn.Exec(
		`INSERT INTO POST (id_post, body_post, id_author, id_parent_post) VALUES (?, ?, ?, ?)`,
		postID, body, authorID, parent,
	); err != nil {
		return err
	}

	_, err := config.Conn.Exec(
		`INSERT INTO TOPIC_POST (id_topic, id_post) VALUES (?, ?)`,
		topicID, postID,
	)
	return err
}

func GetTopicAuthor(id string) (string, error) {
	var authorID string
	err := config.Conn.QueryRow(
		`SELECT id_author FROM TOPIC WHERE id_topic = ? AND deleted_at IS NULL`, id,
	).Scan(&authorID)
	return authorID, err
}

func GetPostAuthor(id string) (string, error) {
	var authorID string
	err := config.Conn.QueryRow(
		`SELECT id_author FROM POST WHERE id_post = ? AND deleted_at IS NULL`, id,
	).Scan(&authorID)
	return authorID, err
}

func UpdateTopicTitle(id, title string) error {
	_, err := config.Conn.Exec(
		`UPDATE TOPIC SET title_topic = ? WHERE id_topic = ? AND deleted_at IS NULL`, title, id,
	)
	return err
}

func UpdatePostBody(id, body string) error {
	_, err := config.Conn.Exec(
		`UPDATE POST SET body_post = ? WHERE id_post = ? AND deleted_at IS NULL`, body, id,
	)
	return err
}

func ModerateContent(contentType, contentID, moderatorID, message string) error {
	var table, idCol string
	switch contentType {
	case "topic":
		table = "TOPIC"
		idCol = "id_topic"
	case "post":
		table = "POST"
		idCol = "id_post"
	default:
		return fmt.Errorf("type de contenu inconnu: %s", contentType)
	}
	_, err := config.Conn.Exec(
		"UPDATE "+table+" SET moderated_by = ?, moderation_message = ? WHERE "+idCol+" = ?",
		moderatorID, message, contentID,
	)
	return err
}
