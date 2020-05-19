package est8

import (
	"database/sql"
	"fmt"
	"log"
)

//funcion a estacionamiento
func Ocho(dbHostinger *sql.DB, inicio string, conexion string, IP string) {

	const (
		USE_MYMYSQL = false // En caso de no funcionar un driver utiliza otro de mysql :3
	)

	driver := ""
	connstr := ""
	if USE_MYMYSQL {
		driver = "mymysql"
		connstr = conexion
		defer Recuperacion("producción")
	} else {
		driver = "mysql"
		connstr = conexion
		defer Recuperacion("Produccion")
	}

	//estacionamiento
	db, err := sql.Open(driver, connstr)

	defer Recuperacion(IP)
	if err != nil {
		//panic(err.Error())
		fmt.Printf("Error al conectar a la IP: " + IP)
		log.Print(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(0)

	defer db.Close()

	fmt.Println("Conectado a la base estacionamiento " + IP + " usando el driver: " + driver)

	var count int

	//query de dos dias
	//error := db.QueryRow("SELECT COUNT(*) FROM ecentral.ope_operacion AS oope JOIN ecentral.ope_pago AS opag ON (oope.ID = opag.ID_OPERACION) WHERE DATE(oope.FECHA) BETWEEN ? AND ? AND oope.folio IS NOT NULL", inicio, final).Scan(&count)
	//query de un dia
	error := db.QueryRow("SELECT COUNT(*) FROM ecentral.ope_operacion AS oope JOIN ecentral.ope_pago AS opag ON (oope.ID = opag.ID_OPERACION) WHERE DATE(oope.FECHA) = ? AND oope.folio IS NOT NULL", inicio).Scan(&count)
	switch {
	case error != nil:
		//log.Fatal("Error presentado en: " +IP+ " error")
		//panic(err.Error())
		fmt.Println("Error al consultar a la IP: " + IP)
		defer Recuperacion(IP)
		panic(err.Error())
		log.Print(err)
	default:
		fmt.Printf("Cantidad de registros de "+IP+": %v\n", count)
	}

	//Query para consulta de estacionamientos

	consulta, err := db.Query("SELECT oope.FECHA, oope.folio, opag.ID_ESTACIONAMIENTO, opag.monto_por_cobrar_local FROM ecentral.ope_operacion AS oope JOIN ecentral.ope_pago AS opag ON (oope.ID = opag.ID_OPERACION) WHERE DATE(oope.FECHA) = ? AND oope.folio IS NOT NULL", inicio)
	defer Recuperacion(IP)
	if err != nil {
		panic(err.Error())
		log.Print(err)
	}

	fmt.Println("Comienza insert estacionamiento " + IP)

	for consulta.Next() {
		//var folio int
		//var monto float32
		var id, folio, monto, Fecha string
		//var Fecha uint8

		//strconv.Itoa(int(Fecha))
		defer Recuperacion(IP)
		if err := consulta.Scan(&id, &monto, &folio, &Fecha)
			err != nil {
			log.Fatal(err.Error())
		} //termina err

		//fmt.Printf("datos encontrados = %s - %s - %v - %v\n", id, monto, folio, Fecha)

		//comienza insert

		tx, err := dbHostinger.Begin()
		defer Recuperacion(IP)
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()

		if folio == "" {
			continue
		} else {
			stmt, err := tx.Prepare("INSERT INTO park_temps (fecha, folio,monto, park) VALUES (?,?,?,?)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close() // danger!

			for i := 0; i < 1; i++ {
				_, err = stmt.Exec(id, monto, Fecha, folio)
				defer Recuperacion(IP)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		err = tx.Commit()
		defer Recuperacion(IP)
		if err != nil {
			log.Fatal(err)
		}
		// stmt.Close() runs here!
		//time.Sleep(500 * time.Millisecond)

		//fmt.Printf("Insertado\n")
	} //termina for
	//loader(consulta)
	fmt.Printf("Termino de insertar estacionamiento " + IP + "\n")

} //termina funcion estacionamiento

func Recuperacion(IP string) {
	recuperado := recover()
	if recuperado != nil {
		fmt.Println("Recuperación de: ", IP, recuperado)
	}
}