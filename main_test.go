package main

import (
	"bufio"
	"strings"
	"testing"
)

var TEST_STRING string = `String property = props.getProperty(key);
descr = (String)qdsSession.getProperty(connectorLog.getAbx(), WF_FWREQ);
buf.append(System.getProperty("line.separator"));
result = is.getPRoperty(arg1);
config.setUsername(p.getProperty("username"));
String reasonsString = (String)qdsSession.getProperty(abx, DEFERRED_NO_SCORE_REASON_PROP);
descr = (String)session.getProperty(connectorLog.getAbx(), "my_key_name");
if (StringUtils.equals(String)session.getProperty(abx, THIS_PROPERTY_1), "hello") && StringUtils.equals(String)session.getProperty(abx, THIS_PROPERTY_2), "hello") {
`

func checkKeys(values []string, scanner *bufio.Scanner, t *testing.T) {
	scanner.Scan()
	text := scanner.Text()
	actual := scanLine(text)
	num := len(values)

	if len(actual) != num {
		t.Errorf("Expected %d keys; found %d in '%s'", num, len(actual), text)
	}

	if len(actual) > 0 {
		for i := 0; i < len(values); i++ {
			if values[i] != actual[i] {
				t.Errorf("Expected %s; found %s in '%s'", values[i], actual[i], text)
			}
		}
	}
}

func TestScanLine(t *testing.T) {
	r := strings.NewReader(TEST_STRING)

	scanner := bufio.NewScanner(r)
	checkKeys([]string{}, scanner, t)
	checkKeys([]string{"WF_FWREQ"}, scanner, t)
	checkKeys([]string{}, scanner, t)
	checkKeys([]string{}, scanner, t)
	checkKeys([]string{}, scanner, t)
	checkKeys([]string{"DEFERRED_NO_SCORE_REASON_PROP"}, scanner, t)
	checkKeys([]string{"\"my_key_name\""}, scanner, t)
	checkKeys([]string{"THIS_PROPERTY_1", "THIS_PROPERTY_2"}, scanner, t)
}
