package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const masterPassword = "123456"
const dbPath = "C:/Users/micael.santos/Documents/password-generator/password.db"

var db *sql.DB

func main() {
	senha := getSenha("Insira sua senha master:")
	if senha != masterPassword {
		fmt.Println("Senha Invalida...")
		os.Exit(1)
	}

	var erro error
	db, erro = sql.Open("sqlite3", dbPath)
	if erro != nil {
		log.Fatal(erro)
	}

	defer db.Close()

	createTable()

	for {
		menu()
		op := getInput("O que deseja fazer?")
		switch op {
		case "i":
			inserirSenha()
		case "l":
			mostrarServicos()
		case "r":
			recuperarSenha()
		case "s":
			os.Exit(0)
		default:
			fmt.Println("Opção inválida")
		}
	}
}

func getSenha(prompt string) string {
	var senha string
	fmt.Print(prompt)
	fmt.Scan(&senha)
	return senha
}

func getInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scan(&input)
	return input
}

func createTable() {
	_, erro := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			service TEXT NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if erro != nil {
		log.Fatal(erro)
	}
}

func menu() {
	fmt.Println("********************************")
	fmt.Println("** i : inserir nova senha     **")
	fmt.Println("** l : listar serviços salvos **")
	fmt.Println("** r : recuperar uma senha    **")
	fmt.Println("** s : sair                   **")
	fmt.Println("********************************")
}

func recuperarSenha() {
	service := getInput("Qual o serviço para o qual quer a senha?")
	rows, erro := db.Query("SELECT username, password FROM users WHERE service = ?", service)
	if erro != nil {
		log.Fatal(erro)
	}
	defer rows.Close()

	if !rows.Next() {
		fmt.Println("Serviço não cadastrado (use 'l' para verificar os serviços)")
		return
	}

	var username, password string
	erro = rows.Scan(&username, &password)
	if erro != nil {
		log.Fatal(erro)
	}

	fmt.Printf("Username: %s\nPasswords: %s\n", username, password)
}

func inserirSenha() {
	service := getInput("Qual o nome do serviço?")
	username := getInput("Qual o nome de usuário?")
	password := getInput("Qual a senha?")

	_, err := db.Exec("INSERT INTO users (service, username, password) VALUES (?, ?, ?)", service, username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Senha inserida com sucesso!")
}

func mostrarServicos() {
	rows, erro := db.Query("SELECT service FROM users")
	if erro != nil {
		log.Fatal(erro)
	}

	defer rows.Close()

	var service string
	for rows.Next() {
		erro := rows.Scan(&service)
		if erro != nil {
			log.Fatal(erro)
		}

		fmt.Println(service)
	}
}
