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
	transactions := strings.Split(string(binlog), "START TRANSACTION")
	return transactions[1:]
}

func ParseTransactionQuery(transaction string) Transaction {
	transactionParts := strings.Split(transaction, "#Q> ")

	totalParts := len(transactionParts)
	var query string
	for i := 1; i < totalParts-1; i++ {
		query = query + transactionParts[i]
	}
	queryParts := strings.SplitN(transactionParts[totalParts-1], "#", 2)
	query = query + queryParts[0]
	query = strings.ReplaceAll(strings.ReplaceAll(query, "\n", " "), "`", "")

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
	delimiterPattern := "(?i)set"
	re := regexp.MustCompile(delimiterPattern)
	setStatement := re.Split(query, -1)[1]

	delimiterPattern = "(?i)where"
	re = regexp.MustCompile(delimiterPattern)
	fieldsAndValues := strings.Split(re.Split(setStatement, -1)[0], ",")

	fields := make([]string, 0)
	values := make([]string, 0)

	for _, data := range fieldsAndValues {
		if data == "" {
			continue
		}

		fieldValue := strings.Split(strings.Trim(data, " "), "=")
		fields = append(fields, string(fieldValue[0]))

		values = append(values, string(fieldValue[1]))
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
