package db

import (
	"database/sql"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func scanPost(row interface{ Scan(...any) error }) (models.Post, error) {
	var p models.Post
	var idParentPost sql.NullString
	err := row.Scan(&p.Id_post, &p.Body_post, &p.Id_author, &p.AuthorName, &idParentPost, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	if idParentPost.Valid {
		p.Id_parent_post = &idParentPost.String
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
		       t.created_at
		FROM TOPIC t
		LEFT JOIN USER u ON u.id_user = t.id_author
		LEFT JOIN TOPIC_POST tp ON tp.id_topic = t.id_topic
		`+where+`
		GROUP BY t.id_topic, t.title_topic, t.id_author, u.first_name, u.last_name, t.created_at
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(
			&t.Id_topic, &t.Title_topic, &t.Id_author, &t.AuthorName, &t.RepliesCount, &t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	return list, total, nil
}

func GetTopicByID(id string) (models.TopicDetail, error) {
	var t models.TopicDetail
	err := config.Conn.QueryRow(`
		SELECT t.id_topic, t.title_topic, t.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       t.created_at
		FROM TOPIC t
		LEFT JOIN USER u ON u.id_user = t.id_author
		WHERE t.id_topic = ? AND t.deleted_at IS NULL`, id,
	).Scan(&t.Id_topic, &t.Title_topic, &t.Id_author, &t.AuthorName, &t.CreatedAt)
	if err != nil {
		return t, err
	}

	rows, err := config.Conn.Query(`
		SELECT p.id_post, p.body_post, p.id_author,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS author_name,
		       p.id_parent_post,
		       p.created_at
		FROM POST p
		JOIN TOPIC_POST tp ON tp.id_post = p.id_post
		LEFT JOIN USER u ON u.id_user = p.id_author
		WHERE tp.id_topic = ? AND p.deleted_at IS NULL
		ORDER BY p.created_at ASC`, id)
	if err != nil {
		return t, err
	}
	defer rows.Close()

	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return t, err
		}
		t.Posts = append(t.Posts, p)
	}
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
