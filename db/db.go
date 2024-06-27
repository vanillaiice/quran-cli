package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "modernc.org/sqlite"
)

type Surah struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Transliteration string  `json:"transliteration"`
	Translation     string  `json:"translation"`
	Type            string  `json:"type"`
	TotalVerses     int     `json:"total_verses"`
	Verses          []Verse `json:"verses"`
}

type Verse struct {
	Id          int    `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

type Conn struct {
	db *sql.DB
}

func New(path string) (*Conn, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS Quran(
			surah_id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			transliteration TEXT NOT NULL,
			translation TEXT NOT NULL,
			type TEXT NOT NULL,
			total_verses INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS Verses(
			id INTEGER PRIMARY KEY,
			surah_id INTEGER NOT NULL,
			verse_id INTEGER NOT NULL,
			text TEXT NOT NULL,
			translation TEXT NOT NULL,
			FOREIGN KEY (surah_id) REFERENCES Quran(surah_id)
		);
	`

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	if _, err = conn.Exec(stmt); err != nil {
		return nil, err
	}

	return &Conn{db: conn}, nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func (c *Conn) InitFromReader(r io.Reader) error {
	var s []*Surah
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return err
	}
	return initDb(s, c)
}

func (c *Conn) InitFromFile(file string) error {
	var surahs []*Surah

	f, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(f, &surahs); err != nil {
		return err
	}

	return initDb(surahs, c)
}

func (c *Conn) GetSurahById(id int) (*Surah, error) {
	stmt := `
		SELECT
			Quran.surah_id,
			Quran.name,
			Quran.transliteration,
			Quran.translation,
			Quran.type,
			Quran.total_verses,
			Verses.verse_id,
			Verses.text,
			Verses.translation
		FROM Quran
		LEFT JOIN verses
		ON Verses.surah_id = Quran.surah_id
		WHERE Quran.surah_id = ?`

	rows, err := c.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surah Surah

	for rows.Next() {
		var v Verse

		if err = rows.Scan(
			&surah.Id,
			&surah.Name,
			&surah.Transliteration,
			&surah.Translation,
			&surah.Type,
			&surah.TotalVerses,
			&v.Id, &v.Text,
			&v.Translation,
		); err != nil {
			return nil, err
		}

		surah.Verses = append(surah.Verses, v)
	}

	return &surah, nil
}

func (c *Conn) GetSurahByName(name string) (*Surah, error) {
	stmt := `
		SELECT
			Quran.surah_id,
			Quran.name,
			Quran.transliteration,
			Quran.translation,
			Quran.type,
			Quran.total_verses,
			Verses.verse_id,
			Verses.text,
			Verses.translation
		FROM Quran
		LEFT JOIN verses
		ON Verses.surah_id = Quran.surah_id
		WHERE Quran.transliteration = ?`

	rows, err := c.db.Query(stmt, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surah Surah

	for rows.Next() {
		var v Verse

		if err = rows.Scan(
			&surah.Id,
			&surah.Name,
			&surah.Transliteration,
			&surah.Translation,
			&surah.Type,
			&surah.TotalVerses,
			&v.Id,
			&v.Text,
			&v.Translation,
		); err != nil {
			return nil, err
		}

		surah.Verses = append(surah.Verses, v)
	}

	return &surah, nil
}

func (c *Conn) GetSurahByNameLike(name string) (*Surah, error) {
	stmt := `
		SELECT Quran.surah_id,
			Quran.name,
			Quran.transliteration,
			Quran.translation,
			Quran.type,
			Quran.total_verses,
			Verses.verse_id,
			Verses.text,
			Verses.translation
		FROM Quran
		LEFT JOIN verses
		ON Verses.surah_id = Quran.surah_id WHERE Quran.transliteration
		LIKE ?
		LIMIT 1`

	rows, err := c.db.Query(stmt, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surah Surah

	for rows.Next() {
		var v Verse

		if err = rows.Scan(
			&surah.Id,
			&surah.Name,
			&surah.Transliteration,
			&surah.Translation,
			&surah.Type,
			&surah.TotalVerses,
			&v.Id,
			&v.Text,
			&v.Translation,
		); err != nil {
			return nil, err
		}

		if v.Id <= len(surah.Verses) {
			break
		}

		surah.Verses = append(surah.Verses, v)
	}

	return &surah, nil
}

func initDb(surahs []*Surah, c *Conn) (err error) {
	var verseId int

	for _, s := range surahs {
		tx, err := c.db.Begin()
		if err != nil {
			return err
		}

		if _, err = tx.Exec(`INSERT INTO Quran VALUES (?, ?, ?, ?, ?, ?)`, s.Id, s.Name, s.Transliteration, s.Translation, s.Type, s.TotalVerses); err != nil {
			return err
		}

		for _, v := range s.Verses {
			verseId++

			if _, err = tx.Exec(`INSERT INTO Verses VALUES (?, ?, ?, ?, ?)`, verseId, s.Id, v.Id, v.Text, v.Translation); err != nil {
				return err
			}
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return
}
