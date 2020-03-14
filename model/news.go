package model

import (
	"fmt"
	"time"

	db "bitbucket.org/michaelchandrag/kumparan-test/database"
)

type (
	News struct {
		ID				int 		`json:"id" db:"id"`
		Author			string 		`json:"author" db:"author"`
		Body			string 		`json:"body" db:"body"`
		Created 		string 		`json:"created" db:"created"`
		EsCreated 		string 		`json:"es_created"`
	}
)

func (this *News) Finds() ([]*News, error){
	query := `
		SELECT
			id,
			author,
			body,
			created
		FROM
			news`
	rows, err := db.Engine.Queryx(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var news []*News
	for rows.Next() {
		_n := &News{}
		err = rows.StructScan(_n)
		if err != nil {
			return nil, err
		}
		news = append(news, _n)
	}

	return news, nil
}

func (this *News) FindByID(id int) error {
	query := `
		SELECT
			id,
			author,
			body,
			created
		FROM
			news
		WHERE
			id = ?`
	if err := db.Engine.Get(this, query, id); err != nil {
		return err
	}

	return nil
}

func (this *News) Create(data News) (result News, err error) {
	currentTime := time.Now()

	data.Created = currentTime.Format("2006-01-02 15:04:05")
	data.EsCreated = currentTime.Format("2006/01/02 15:04:05") // https://www.elastic.co/guide/en/elasticsearch/reference/current/date.html
	fmt.Println(data.EsCreated)
	query := fmt.Sprintf(`
		INSERT INTO news (
			author, body, created
		) VALUES (
			?, ?, ?
		)
	`)
	resp, err := db.Engine.Exec(query,
		data.Author, data.Body, data.Created)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	lastID, _ := resp.LastInsertId()
	data.ID = int(lastID)
	result.ID = data.ID
	result.Author = data.Author
	result.Body = data.Body
	result.Created = data.Created
	result.EsCreated = data.EsCreated
	return result, nil
}