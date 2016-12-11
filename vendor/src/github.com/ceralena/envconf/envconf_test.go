package envconf

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestInvalidConfig(t *testing.T) {
	tests := []struct {
		v        interface{}
		errmatch string
	}{
		{make(map[string]string), "Invalid kind for config: "},
		{[]string{}, "Invalid kind for config: "},
		{
			struct {
				M map[string]string `required:"true"`
			}{
				make(map[string]string),
			}, "Invalid kind for config field",
		},
	}
	tm := mapgetter{"M": "hi"}

	for _, test := range tests {
		err := ReadConfig(test.v, tm.get)
		if err == nil {
			t.Errorf("Expected an error for config val: %v", test.v)
			t.Fail()
		} else if err != nil && !strings.Contains(err.Error(), test.errmatch) {
			t.Errorf("Expected a different error for ReadConfig(): looking for '%s' in '%s'", test.errmatch, err.Error())
			t.Fail()
		}
	}
}

func TestConfig(t *testing.T) {
	type MyConf struct {
		Foo      string `required:"true"`
		Bar      int
		On       bool
		Def      string `default:"somedefault"`
		Some     []string
		SomeInt  []int
		SomeBool []bool
		ignored  bool
	}
	tests := []struct {
		vals     mapgetter
		valid    bool
		errmatch string
	}{
		{mapgetter{"FOO": "hehe", "BAr": "3", "on": "TRUE"}, true, ""},
		{mapgetter{"FOO": "hehe", "BAr": "3", "on": "TRUE", "SOME": "yes,no"}, true, ""},

		{mapgetter{"FOO": "hehe", "BAR": "3"}, true, ""},

		// missing a required field
		{mapgetter{"BAR": "3", "ON": "true"}, false, "Missing config fields: "},

		// invalid int
		{mapgetter{"FOO": "hehe", "BAR": "sup", "ON": "true"}, false, "strconv.ParseInt: "},

		// invalid int list
		{mapgetter{"FOO": "hehe", "BAr": "3", "on": "TRUE", "SOMEINT": "yes,no"}, false, "strconv.ParseInt: "},

		// invalid bool list
		{mapgetter{"FOO": "hehe", "BAr": "3", "on": "TRUE", "SOMEBOOL": "yes,no"}, false, "strconv.ParseBool: "},

		// invalid bool
		{mapgetter{"FOO": "hehe", "BAR": "3", "ON": "damn"}, false, "strconv.ParseBool: "},

		// ignore unexported
		{mapgetter{"FOO": "hehe", "BAR": "3", "ON": "true", "ignored": "fdjhkl"}, true, ""},
	}

	for _, test := range tests {
		c := MyConf{}
		err := ReadConfig(&c, test.vals.get)
		if err != nil && test.valid {
			t.Errorf("Unexpected error with '%v': %v", test.vals, err)
			t.Fail()
		} else if err == nil && !test.valid {
			t.Errorf("Expected an error with: %v", test.vals)
			t.Fail()
		} else if err != nil && !strings.Contains(err.Error(), test.errmatch) {
			t.Errorf("Error strings did not match for err '%v': looking for '%s'", err, test.errmatch)
			t.Fail()
		}
	}
}

func TestConfigDefaults(t *testing.T) {
	var myConf struct {
		K string `default:"foo"`
	}
	err := ReadConfig(&myConf, make(mapgetter).get)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.Fail()
	}
	if myConf.K != "foo" {
		t.Errorf("ReadConfig(): Expected default of 'foo'; got '%s'", myConf.K)
		t.Fail()
	}
}


func TestConfigMap(t *testing.T) {
	var myConf struct {
		K string
	}
	m := map[string]string{"K": "ah"}
	err := ReadConfigMap(&myConf, m)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.Fail()
	}
	if myConf.K != "ah" {
		t.Errorf("ReadConfig(): Expected default of \"ah\"; got %q'", myConf.K)
		t.Fail()
	}
}

func TestConfigBadSlice(t *testing.T) {
	var myConf struct {
		Hi []struct{} `required:"true"`
	}
	input := mapgetter{"HI": "a,b,c"}
	match := "[]struct {}"

	if err := ReadConfig(&myConf, input.get); err == nil || !strings.Contains(err.Error(), match) {
		t.Errorf("ReadConfig(): expected an error matching '%s', got '%v'", match, err)
		t.Fail()
	}
}

// Test a config object with slice values.
func TestConfigSlice(t *testing.T) {
	var myConf struct {
		Ints    []int
		Bools   []bool
		Strings []string
	}
	var input = mapgetter{
		"INTS":    "1,2,3,4",
		"BOOLS":   "true,false,true",
		"STRINGS": "hello,yes,hi,lol,ok",
	}
	var (
		expectInts    = []int{1, 2, 3, 4}
		expectBools   = []bool{true, false, true}
		expectStrings = []string{"hello", "yes", "hi", "lol", "ok"}
	)

	err := ReadConfig(&myConf, input.get)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.Fail()
	}

	if len(myConf.Ints) != len(expectInts) {
		t.Errorf("Wrong length for Ints: wanted %d, got %d", len(expectInts), len(myConf.Ints))
		t.FailNow()
	}
	if len(myConf.Bools) != len(expectBools) {
		t.Errorf("Wrong length for Bools: wanted %d, got %d", len(expectBools), len(myConf.Bools))
		t.FailNow()
	}
	if len(myConf.Strings) != len(expectStrings) {
		t.Errorf("Wrong length for Strings: wanted %d, got %d", len(expectStrings), len(myConf.Strings))
		t.FailNow()
	}

	for i, iv := range expectInts {
		if eiv := myConf.Ints[i]; eiv != iv {
			t.Errorf("Ints[%d]: expected %d, got %d", i, iv, eiv)
			t.Fail()
		}
	}
	for i, bv := range expectBools {
		if ebv := myConf.Bools[i]; ebv != bv {
			t.Errorf("Bools[%d]: expected %d, got %d", i, bv, ebv)
			t.Fail()
		}
	}
	for i, sv := range expectStrings {
		if esv := myConf.Strings[i]; esv != sv {
			t.Errorf("Strings[%d]: expected %d, got %d", i, sv, esv)
			t.Fail()
		}
	}
}

func TestConfigEnv(t *testing.T) {
	// Test of real environment
	os.Setenv("ENVCONFTEST1", "foo")
	defer os.Setenv("ENVCONFTEST1", "")
	var conf struct {
		ENVCONFTEST1 string
	}
	if err := ReadConfigEnv(&conf); err != nil {
		t.Errorf("Unexpected error %v", err)
		t.FailNow()
	}
	if v := conf.ENVCONFTEST1; v != "foo" {
		t.Errorf("ReadConfigEnv: got '%s', wanted 'foo'", v)
		t.Fail()
	}
}

func TestConfigEnvPrefix(t *testing.T) {
	// Test of real environment
	os.Setenv("FOO_ENVCONFTEST1", "foo")
	defer os.Setenv("FOO_ENVCONFTEST1", "")
	var conf struct {
		ENVCONFTEST1 string
	}
	if err := ReadConfigEnvPrefix("FOO_", &conf); err != nil {
		t.Errorf("Unexpected error %v", err)
		t.FailNow()
	}
	if v := conf.ENVCONFTEST1; v != "foo" {
		t.Errorf("ReadConfigEnvPrefix: got '%s', wanted 'foo'", v)
		t.Fail()
	}
}

func ExampleReadConfigEnv() {
	os.Setenv("FOO", "hi")
	os.Setenv("BAR", "yes")

	defer os.Setenv("FOO", "")
	defer os.Setenv("BAR", "")

	var conf struct {
		Foo string
		Bar string
	}

	if err := ReadConfigEnv(&conf); err != nil {
		panic(err)
	}

	fmt.Println(conf.Foo)
	fmt.Println(conf.Bar)
	// Output:
	// hi
	// yes
}
