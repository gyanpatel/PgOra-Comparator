package main

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"
)

// compareOraPgTables compares Oracle and PG table data for given queries as input
// and prints the comparison result
func compareOraPgTables(selTableQueryList []TableQueryList) (*[]ComparisonResult, error) {
	var comparisonResult []ComparisonResult
	date := time.Now().Format("02-Jan-2006")
	for _, queryRec := range selTableQueryList {
		log.Println("Comparison started for table ", queryRec.TableName)
		var pgTableBytes []byte
		var pgTableBytesStr []string
		rowsPG, err := dbConPG.Query(queryRec.PgQuery)
		if err != nil {
			log.Println(" PG - Unable to query data ", err)
			return &comparisonResult, err
		}
		log.Println(" Fetching PG Data Table = ", queryRec.TableName)
		for rowsPG.Next() {
			var rowvalueStr sql.NullString
			err := rowsPG.Scan(&rowvalueStr)
			if err != nil {
				log.Println(" PG - Unable to scan query data ", err)
				return &comparisonResult, err
			}
			//rowvalueByte := []byte(rowvalueStr)
			pgTableBytesStr = append(pgTableBytesStr, rowvalueStr.String)
		}
		//pgTableBytes = []byte(pgTableBytesStr)
		log.Println(" Fetching PG Data - Completed Table = ", queryRec.TableName)
		var oraTableBytes []byte
		var oraTableBytesStr []string
		rowsOra, err := dbConOra.Query(queryRec.OraQuery)
		if err != nil {
			log.Println(" Oracle - Unable to query data ", err, queryRec.OraQuery)
			return &comparisonResult, err
		}
		log.Println(" Fetching Oracle Data Table = ", queryRec.TableName)
		for rowsOra.Next() {
			var rowvalueStr sql.NullString
			err := rowsOra.Scan(&rowvalueStr)
			if err != nil {
				log.Println("Oracle - Unable to scan query data ", err)
				return &comparisonResult, err
			}
			//rowvalueByte := []byte(rowvalueStr)
			//oraTableBytes = append(oraTableBytes, rowvalueByte...)
			oraTableBytesStr = append(oraTableBytesStr, rowvalueStr.String)
		}
		log.Println(" Fetching Oracle Data - Completed Table = ", queryRec.TableName)
		//pgTableBytesStr = sort.StringSlice(pgTableBytesStr)
		pgCount := len(pgTableBytesStr)
		oraCount := len(oraTableBytesStr)
		log.Println(" Sorting slices Data")
		sort.Strings(pgTableBytesStr)
		sort.Strings(oraTableBytesStr)
		log.Println("Converting slice of string to slice of byte")
		pgTableBytes = []byte(strings.Join(pgTableBytesStr, ""))
		oraTableBytes = []byte(strings.Join(oraTableBytesStr, ""))
		if len(pgTableBytes) == 0 && len(oraTableBytes) == 0 {
			log.Println("###################### Both Oracle and PG tables are empty ###################### Table = ", queryRec.TableName)
			comparisonResult = append(comparisonResult, ComparisonResult{Date: date, Result: "Matches", TableName: queryRec.TableName, PGCount: pgCount, OraCount: oraCount})
		} else {
			log.Println("Table comparison Started ", time.Now())

			res := bytes.Compare(pgTableBytes, oraTableBytes)
			log.Println("Table comparison Completed ", time.Now())
			if res == 0 {
				log.Println("###################### Oracle and PG Data match ###################### Table = ", queryRec.TableName)
				comparisonResult = append(comparisonResult, ComparisonResult{Date: date, Result: "Matches", TableName: queryRec.TableName, PGCount: pgCount, OraCount: oraCount})

			} else {
				log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!! Oracle and PG Data do NOT match !!!!!!!!!!!!!!!!!!!!!!!!!")
				dataDiff := getTableDiff(&oraTableBytesStr, &pgTableBytesStr)
				log.Printf(" ---- Table diff starts for %s---- \n%v\n ---- Table diff ends ----\n", queryRec.TableName, *dataDiff)
				comparisonResult = append(comparisonResult, ComparisonResult{Date: date, Result: "Differs", DataDiff: *dataDiff, TableName: queryRec.TableName, PGCount: pgCount, OraCount: oraCount})

			}
		}
	}
	log.Println("Comparison Completed")
	return &comparisonResult, nil
}

func getTableDiff(oracData, pgData *[]string) *[]DataDiff {
	var dataDiff []DataDiff
	dbName := "Oracle "
	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, set1 := range *oracData {
			found := false
			for _, set2 := range *pgData {
				if set1 == set2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				dataDiff = append(dataDiff, DataDiff{DbName: dbName, Data: set1})
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			oracData, pgData = pgData, oracData
			dbName = "PG "
		}
	}
	return &dataDiff
}

func getTableQuery() error {
	sqlQuery, err := ioutil.ReadFile("CompQuery.sql")
	if err != nil {
		log.Println("getTableQuery : UNABLE TO READ CompQuery.sql ,  FAILURE ", err)
		return err
	}
	sqlQueryStr := string(sqlQuery)
	rowsTableList, err := dbConOra.Query(sqlQueryStr)
	if err != nil {
		log.Println("getTableQuery : UNABLE TO QUERY TABLE LIST ,  FAILURE ", err)
		return err
	}
	for rowsTableList.Next() {
		var tableName, query sql.NullString
		err := rowsTableList.Scan(&tableName, &query)
		if err != nil {
			log.Println("getTableQuery : UNABLE TO SCAN TABLE LIST ,  FAILURE ", err)
			return err

		}
		//rowvalueByte := []byte(rowvalueStr)
		//oraTableBytes = append(oraTableBytes, rowvalueByte...)
		tableQueryList = append(tableQueryList, TableQueryList{TableName: tableName.String, PgQuery: query.String, OraQuery: query.String})
	}
	return nil
}
