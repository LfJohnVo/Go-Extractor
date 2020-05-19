package main

import (
	"./estacionamientos/est8"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

func main() {
	t := time.Now()
	//fmt.Println(t.String())
	fmt.Printf("Fecha hoy: ")
	fmt.Println(t.Format("2006-01-02"))
	fmt.Printf("Dia anterior: ")
	fecha2 := t.Add(365*24 - (24 * time.Hour))
	//fecha2 := "2019-11-13"
	fecha3 := fecha2.Format("2006/01/02")
	fmt.Println(fecha3)

	fmt.Println("Extracción de base de datos a facturador => ", t)

	//variable para array
	var estacionamientos []string

	//hostinger
	dbDestino := "XXXXXXX:XXXXX@tcp(XXXXXXX:3306)/XXXXXXX"

	fmt.Printf("Ingresa fecha de inicio:\n")
	fmt.Scanf("%s", &inicio)
	//time.Sleep(time.Second * 3)
	fmt.Printf("Ingresa fecha final:\n")
	fmt.Scanf("%s", &final)
	//time.Sleep(time.Second * 3)*/

	//inicio := "2020/03/03"
	//final := "2020/03/05"


	file, err := os.Open("lista.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println("Lei:")
		//fmt.Println(scanner.Text())
		estacionamientos = append(estacionamientos, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}


	//fmt.Printf("%v\n", estacionamientos)

	ArrLen := len(estacionamientos)

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	//Creación de log
	//errr := ioutil.WriteFile("logs/"+fecha2+".txt", []byte("/***************Comienza log***************/"), 0755)
	//if errr != nil {
	//	fmt.Println("Error al leer log: %v", err)
	//}
	//conexión a hostinger

	dbHostinger, err := sql.Open("mysql", dbDestino)

	defer Recuperacion(dbDestino)
	if err != nil {
		panic(err.Error())
		log.Print(err)
	}

	defer dbHostinger.Close()

	fmt.Println("Conectado a la base destino")

	for f := 0; f < ArrLen; f++ {
		//fmt.Print("datos de conexion: ")
		no_est := estacionamientos[f]
		// Lectura de valores, pueden ser representados como vacios
		DB_IP := cfg.Section(no_est).Key("DB_IP").String()
		DB_NAME := cfg.Section(no_est).Key("DB_NAME").String()
		DB_USER := cfg.Section(no_est).Key("DB_USER").String()
		DB_PASSWORD := cfg.Section(no_est).Key("DB_PASSWORD").String()

		//fmt.Println(DB_IP,DB_NAME,DB_USER, DB_PASSWORD)
		Conexion := DB_USER + ":" +DB_PASSWORD + "@tcp" + "(" +DB_IP + ":3306)/" +DB_NAME
		go est8.Ocho(dbHostinger, fecha3, Conexion, DB_IP)

	}

	defer elapsed("Robot")()
	time.Sleep(time.Minute * 70)
	fmt.Println("Finalizando aplicación....")
	time.Sleep(time.Second * 5)
	fmt.Println("Busque en la carpeta reportes su excel")
	time.Sleep(time.Second * 2)
	fmt.Println("El programa se cerrara automaticamente")
	time.Sleep(time.Second * 5)

} //termina main

//funcion de cronometro
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s Finalizo en: %v\n", what, time.Since(start))
	}
}

func Recuperacion(IP string) {
	recuperado := recover()
	if recuperado != nil {
		fmt.Println("Recuperación de: ", IP, recuperado)
	}
}