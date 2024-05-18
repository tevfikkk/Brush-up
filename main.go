package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var query *pgx.Conn

func main() {
	dbUrl := "postgres://postgres:184812@localhost:5432/recordings"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to the database %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	query = conn // Assigning the connection to the global variable

	fmt.Println("Connected to the database")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		fmt.Errorf("error while searching for the artist: %v", err)
	} else {
		for _, album := range albums {
			fmt.Printf("ID: %d, Title: %s, Artist: %s, Price: %.2f\n", album.ID, album.Title, album.Artist, album.Price)
		}
	}

	fmt.Println("---------------------------------")

	newAlbum := Album{
		Title:  "skibid",
		Artist: "Goon on ya",
		Price:  31.69,
	}
	albumID, err := addAlbum(newAlbum)
	if err != nil {
		log.Fatalf("err happened: %v", err)
	}
	fmt.Println(albumID)
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := query.Query(context.Background(), "SELECT id, title, artist, price FROM album WHERE artist=$1", name)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	for rows.Next() {
		var album Album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Title)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", rows.Err())
	}

	return albums, nil
}

func addAlbum(album Album) (int64, error) {
	var id int64
	err := query.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", album.Title, album.Artist, album.Price).Scan(&id)

	if err != nil {
		return 0, err
	}

	log.Println("Testing1... and id: ", id)

	
	return id, nil
}