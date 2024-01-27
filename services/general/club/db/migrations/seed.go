package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	conf := config.LoadConfig()
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.User, conf.DB.Password, conf.DB.Url, conf.DB.Port, conf.DB.Database)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db")
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db")
	}

	seed(ctx, conn)
}

func seed(ctx context.Context, conn *pgx.Conn) {
	clubs := []types.ClubCreate{
		{Name: "Speedy Stallions", Website: "https://www.speedystallions.net"},
		{Name: "Whinny Warriors", Website: "https://www.whinnywarriors.club"},
		{Name: "Buckle Up", Website: "https://www.buckleup.club"},
	}
	var ids = make([]string, len(clubs))
	for i, club := range clubs {
		var id uint64
		err := conn.QueryRow(ctx, "INSERT INTO club.clubs (name, website) VALUES ($1, $2) RETURNING id;", club.Name, club.Website).Scan(&id)
		if err != nil {
			fmt.Printf("error inserting club %d: %v\n", i, err)
		}
		ids[i] = fmt.Sprintf("%d", id)
	}

	members := []types.MemberCreate{
		{FirstName: "John", MiddleName: "J.", LastName: "Wick", Email: "john.wick@movies.com", Phone: "(468) 7514434", Height: 186, Weight: 80, ClubID: ids[0]},
		{FirstName: "The", MiddleName: "", LastName: "Doctor", Email: "thedoctor@gailfray.gf", Phone: "(763) 3543091", Height: 185, Weight: 73, ClubID: ids[0]},
		{FirstName: "Clara", MiddleName: "Oswin", LastName: "Oswald", Email: "coo@impossile.gl", Phone: "(850) 8504498", Height: 157, Weight: 50, ClubID: ids[0]},
		{FirstName: "Anthony", MiddleName: "J.", LastName: "Crowley", Email: "anthony.j.crowly@666.hell", Phone: "(666) 6666666", Height: 185, Weight: 73, ClubID: ids[0]},
		{FirstName: "Aziraphale", LastName: "Angel", Email: "cesa@heaven.com", Phone: "(636) 1466262", Height: 163, Weight: 59, ClubID: ids[0]},
		{FirstName: "Samuel", LastName: "Winchester", Email: "sam.winchester@winchesters.com", Phone: "(989) 2771892", Height: 193, Weight: 95, ClubID: ids[0]},
		{FirstName: "Dean", LastName: "Winchester", Email: "dean.winchester@winchesters.com", Phone: "(268) 1086400", Height: 186, Weight: 82, ClubID: ids[0]},
		{FirstName: "Walter", MiddleName: "Hartwell", LastName: "White", Email: "walter.white@danger.co", Phone: "(813) 3397491", Height: 179, Weight: 83, ClubID: ids[0]},
		{FirstName: "Jesse", MiddleName: "Bruce", LastName: "Pinkman", Email: "jesse.pinkman@bitch.com", Phone: "(370) 4704697", Height: 173, Weight: 69, ClubID: ids[0]},
		{FirstName: "Gustavo", LastName: "Fring", Email: "gustavo.fring@los-pollos-hermanos.me", Phone: "(559) 2604877", Height: 173, Weight: 71, ClubID: ids[0]},
		{FirstName: "James", MiddleName: "Morgan", LastName: "McGill", Email: "better.call.saul@saulgoodman.usa", Phone: "(124) 2283271", Height: 175, Weight: 82, ClubID: ids[0]},
		{FirstName: "James", MiddleName: "Morgan", LastName: "McGill", Email: "better.call.saul@saulgoodman.usa", Phone: "(124) 2283271", Height: 175, Weight: 82, ClubID: ids[1]},
	}
	for i, member := range members {
		_, err := conn.Exec(ctx, "INSERT INTO club.members (first_name, middle_name, last_name, email, phone, height, weight, club_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);", member.FirstName, member.MiddleName, member.LastName, member.Email, member.Phone, member.Height, member.Weight, member.ClubID)
		if err != nil {
			fmt.Printf("error inserting member %d: %v\n", i, err)
		}
	}
}
