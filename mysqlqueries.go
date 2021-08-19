package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type branchStruct struct {
	BranchId   string `json:"branchId"`
	ConnId     string `json:"connId"`
	BranchName string `json:"branchName"`
	Script     string `json:"script"`
}

type connStuct struct {
	ConnId         string
	Host           string
	Address        string
	PrivateKeyPath string
	ConnName       string
}

func MySQL_GetBranches(connId string) (int, string, [][]string) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	rows, err := msDB.Query("SELECT BranchId, ConnId, BranchName, BranchScript FROM Devina.Branches WHERE ConnId=?", connId)
	if err != nil {
		log.Println("Error in SQL Branch Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Branch Fetch : " + string(err.Error())
		return request_code, request_msg, nil
	}

	defer rows.Close()

	var BranchArr [][]string
	for rows.Next() {
		var tempBranch branchStruct
		err = rows.Scan(&tempBranch.BranchId, &tempBranch.ConnId, &tempBranch.BranchName, &tempBranch.Script)
		if err != nil {
			request_code = 0
			request_msg = "Error in SQL Rows Scan : " + string(err.Error())
			return request_code, request_msg, nil
		}

		tempArr := []string{tempBranch.BranchId, tempBranch.ConnId, tempBranch.BranchName, tempBranch.Script}

		BranchArr = append(BranchArr, tempArr)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, nil
	}

	return 1, "Success", BranchArr
}

func MySQL_GetBranchName(BranchId string) (int, string, string) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	var branchName string
	err := msDB.QueryRow("SELECT BranchName FROM Devina.Branches WHERE BranchId=?", BranchId).Scan(&branchName)
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, branchName
	}
	return 1, "Success", branchName
}

func MySQL_GetIntegrations() (int, string, [][]string) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	rows, err := msDB.Query("SELECT * FROM Devina.Integrations")
	if err != nil {
		log.Println("Error in SQL Branch Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Branch Fetch : " + string(err.Error())
		return request_code, request_msg, nil
	}

	defer rows.Close()

	var BranchArr [][]string
	for rows.Next() {
		var tempBranch IntegrationStruct
		err = rows.Scan(&tempBranch.IntegrationId, &tempBranch.Connid, &tempBranch.IntegrationName, &tempBranch.IbranchId, &tempBranch.FbranchId, &tempBranch.Iscript, &tempBranch.Fscript, &tempBranch.SSHscript)
		if err != nil {
			request_code = 0
			request_msg = "Error in SQL Rows Scan : " + string(err.Error())
			return request_code, request_msg, nil
		}

		req_code, req_msg, ibranchName := MySQL_GetBranchName(tempBranch.IbranchId)
		if req_code == 0 {
			return req_code, req_msg, nil
		}

		req_code, req_msg, fbranchName := MySQL_GetBranchName(tempBranch.FbranchId)
		if req_code == 0 {
			return req_code, req_msg, nil
		}

		tempArr := []string{tempBranch.IntegrationId, tempBranch.Connid, tempBranch.IntegrationName, ibranchName, fbranchName, tempBranch.Iscript, tempBranch.Fscript, tempBranch.SSHscript}
		BranchArr = append(BranchArr, tempArr)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, nil
	}

	return 1, "Success", BranchArr
}

func MySQL_GetScripts(IntegrationId string) (int, string, []string) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	type branchStruct struct {
		branchId int
		initial  string
		final    string
		typez    int
	}

	var IScript string
	var FScript string
	err := msDB.QueryRow("SELECT InitialScript, FinalScript FROM Devina.Integrations WHERE IntegrationId=?", &IntegrationId).Scan(&IScript, &FScript)
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, nil
	}

	return 1, "Success", []string{IScript, FScript}
}

func MYSQL_GetSingleConnectionParams(ConnId string) (int, string, connStuct) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	var tempStruct connStuct
	err := msDB.QueryRow("SELECT ConnId, Host, Address, PrivateKeyPath, ConnName FROM Devina.ConnectionParam WHERE ConnId=?", ConnId).Scan(&tempStruct.ConnId, &tempStruct.Host, &tempStruct.Address, &tempStruct.PrivateKeyPath, &tempStruct.ConnName)
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, tempStruct
	}

	return 1, "Success", tempStruct
}

func MYSQL_GetAllConnections() (int, string, []connStuct) {

	var request_code int
	var request_msg string

	msDB := MySQLConnect()
	defer msDB.Close()

	var tempStructArr []connStuct
	rows, err := msDB.Query("SELECT ConnId, Host, Address, ConnName FROM Devina.ConnectionParam")
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, tempStructArr
	}

	defer rows.Close()

	for rows.Next() {
		var tempStruct connStuct
		err = rows.Scan(&tempStruct.ConnId, &tempStruct.Host, &tempStruct.Address, &tempStruct.ConnName)
		if err != nil {
			request_code = 0
			request_msg = "Error in SQL Rows Scan : " + string(err.Error())
			return request_code, request_msg, tempStructArr
		}

		tempStructArr = append(tempStructArr, tempStruct)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error in SQL Rows Fetch : ", err)
		request_code = 0
		request_msg = "Error in SQL Rows Fetch : " + string(err.Error())
		return request_code, request_msg, tempStructArr
	}

	return 1, "Success", tempStructArr
}

func AddBranch(initial string, final string, typez int) int {

	msDB := MySQLConnect()
	defer msDB.Close()

	err := msDB.QueryRow("INSERT INTO `Devina`.`Branches` (`Initial`, `Final`, `Type`) VALUES (?, ?, ?)", initial, final, typez)
	if err != nil {
		log.Println("Error in SQL Branch Insertion : ", err)
		return 0
	}

	return 1
}

func MYSQL_AddIntegration(integrationId string, connid string, integrationName string, ibranchId string, fbranchId string, iscript string, fscript string, SSHscript string) (int, string) {

	msDB := MySQLConnect()
	defer msDB.Close()

	err := msDB.QueryRow("INSERT INTO `Devina`.`Integrations` (`IntegrationId`, `ConnId`, `IntegrationName`, `InitialBranchId`, `FinalBranch`, `InitialScript`, `FinalScript`, `SSHScript`) VALUES (?,?,?,?,?,?,?,?);", integrationId, connid, integrationName, ibranchId, fbranchId, iscript, fscript, SSHscript)
	if err.Err() != nil {
		log.Println("Error in SQL Script Insertion : ", err.Err())
		return 0, err.Err().Error()
	}

	return 0, err.Err().Error()
}

func AddConnParameters(Host string, Address string, PrivateKeyPath string, ConnName string, ConnId string) error {

	msDB := MySQLConnect()
	defer msDB.Close()

	err := msDB.QueryRow("INSERT INTO `Devina`.`ConnectionParam` (`ConnId`, `Host`, `Address`, `PrivateKeyPath`, `ConnName`) VALUES (?, ?, ?, ?, ?)", &ConnId, &Host, &Address, &PrivateKeyPath, &ConnName)
	if err.Err() != nil {
		log.Println("Error in SQL Script Insertion : ", err.Err())
		return err.Err()
	}

	return nil
}

func EditConnParameters(Host string, Address string, PrivateKeyPath string, ConnName string, ConnId string) error {

	msDB := MySQLConnect()
	defer msDB.Close()

	log.Println("Conn Id SQL", ConnId)

	err := msDB.QueryRow("UPDATE `Devina`.`ConnectionParam` SET `Host` = ?, `Address` = ?, `PrivateKeyPath` = ?, `ConnName` = ? WHERE (`ConnId` = ?)", &Host, &Address, &PrivateKeyPath, &ConnName, &ConnId)
	if err.Err() != nil {
		log.Println("Error in SQL Script Insertion : ", err.Err())
		return err.Err()
	}

	return nil
}

func MySQL_PostAddBranch(BranchInfo branchStruct) (int, string) {

	msDB := MySQLConnect()
	defer msDB.Close()

	log.Println("BranchInfo", BranchInfo)

	err := msDB.QueryRow("INSERT INTO `Devina`.`Branches` (`BranchId`, `ConnId`, `BranchName`, `BranchScript`)	VALUES (?, ?, ?, ?)", &BranchInfo.BranchId, &BranchInfo.ConnId, &BranchInfo.BranchName, &BranchInfo.Script)
	if err.Err() != nil {
		log.Println("Error in SQL Script Insertion : ", err.Err())
		return 0, string(err.Err().Error())
	}

	return 1, "Success"
}

func MySQLConnect() *sql.DB {

	var MYSQLHOST = os.Getenv("MYSQL_HOST")
	var MYSQLPORT = os.Getenv("MYSQL_PORT")
	var MYSQLUSER = os.Getenv("MYSQL_USER")
	var MYSQLPASS = os.Getenv("MYSQL_PASS")
	var MYSQLDB = os.Getenv("MYSQL_AUTHDB")

	mySQLdb, err := sql.Open("mysql", MYSQLUSER+":"+MYSQLPASS+"@tcp("+MYSQLHOST+":"+MYSQLPORT+")/"+MYSQLDB+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// defer mySQLdb.Close() // Causes connection problem (sql server not running.)
	return mySQLdb
}
