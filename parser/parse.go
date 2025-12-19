package parser

import (
	"regexp"
	"strings"
)

type Transaction struct {
	Database string
	Table    string
	Type     string
	Ts       uint64
	Data     map[string]any
	Old      map[string]any
	RawQuery string
}

func ParseTransactions(binlog []byte) []string {
	results := strings.Split(string(binlog), "START TRANSACTION")
	size := len(results)
	transactions := make([]string, 0)
	for i := 1; i < size; i++ {
		transactions = append(transactions, results[i])
	}
	return transactions
}

func ParseTransactionQuery(transaction string) Transaction {
	transactionParts := strings.Split(transaction, "#Q> ")

	queryParts := strings.SplitN(transactionParts[1], "#", 2)
	query := queryParts[0]

	t := Transaction{
		Data:     make(map[string]any),
		Old:      make(map[string]any),
		RawQuery: query,
	}

	ParseStatementLog(queryParts[1], &t)
	ParseStatement(query, &t)
	return t
}

func ParseStatement(query string, t *Transaction) {
	switch t.Type {
	case "INSERT":
		parseInsert(query, t)
	case "UPDATE":
		parseUpdate(query, t)
	default:
		return
	}
}

func parseInsert(query string, t *Transaction) {
	splittedQuery := strings.Split(query, "(")
	fields := strings.Split(strings.Split(splittedQuery[1], ")")[0], ",")
	values := strings.Split(strings.Split(splittedQuery[2], ")")[0], ",")
	setData(fields, values, t)
}

func parseUpdate(query string, t *Transaction) {
	splittedQuery := strings.Split(query, "=")
	firstField := strings.Split(splittedQuery[0], " ")
	splittedQuery[0] = firstField[len(firstField)-2]

	delimiterPattern := "(?i)where"
	re := regexp.MustCompile(delimiterPattern)
	lastValue := re.Split(splittedQuery[1], -1)
	splittedQuery[len(splittedQuery)-2] = lastValue[0]

	fields := make([]string, 0)
	values := make([]string, 0)

	for index, data := range splittedQuery {
		if data == "" {
			continue
		}

		if index%2 == 0 {
			fields = append(fields, data)
			continue
		}

		values = append(values, data)
	}

	setData(fields, values, t)
}

func setData(fields, values []string, t *Transaction) {
	for index := range values {
		t.Data[strings.TrimSpace(fields[index])] = strings.TrimSpace(values[index])
	}
}

func ParseStatementLog(query string, t *Transaction) {
	statementLog := strings.Split(query, "### ")
	firstLine := strings.Split(statementLog[1], " ")
	getQueryType(firstLine[0], firstLine[1], t)

	setDbAndTable(firstLine, t)
	parseContext(statementLog, t)
}

func parseContext(statementLog []string, t *Transaction) {
	switch t.Type {
	case "INSERT":
		t.Data["id_"+t.Table] = strings.Split(strings.Split(statementLog[3], " /*")[0], "=")[1]
	case "UPDATE":
		t.Data["id_"+t.Table] = strings.Split(strings.Split(statementLog[3], " /*")[0], "=")[1]
	default:
		return
	}
}

func setDbAndTable(firstLine []string, t *Transaction) {
	switch t.Type {
	case "INSERT":
		dbData := strings.Split(firstLine[2], "`")
		t.Database = dbData[1]
		t.Table = dbData[3]
	case "UPDATE":
		dbData := strings.Split(firstLine[1], "`")
		t.Database = dbData[1]
		t.Table = dbData[3]
	default:
	}
}

func getQueryType(q1, q2 string, t *Transaction) {
	switch q1 {
	case "INSERT":
		t.Type = "INSERT"
	case "UPDATE":
		t.Type = "UPDATE"
	default:
		t.Type = q1 + "_" + q2
	}
}
